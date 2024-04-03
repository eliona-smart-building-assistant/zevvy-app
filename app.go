//  This file is part of the eliona project.
//  Copyright © 2024 LEICOM iTEC AG. All Rights Reserved.
//  ______ _ _
// |  ____| (_)
// | |__  | |_  ___  _ __   __ _
// |  __| | | |/ _ \| '_ \ / _` |
// | |____| | | (_) | | | | (_| |
// |______|_|_|\___/|_| |_|\__,_|
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
//  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
//  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"context"
	"fmt"
	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-eliona/app"
	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"github.com/eliona-smart-building-assistant/go-eliona/dashboard"
	"github.com/eliona-smart-building-assistant/go-eliona/frontend"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/db"
	utilshttp "github.com/eliona-smart-building-assistant/go-utils/http"
	"github.com/eliona-smart-building-assistant/go-utils/log"
	"net/http"
	"sync"
	"time"
	"zevvy/apiserver"
	"zevvy/apiservices"
	"zevvy/appdb"
	"zevvy/conf"
	"zevvy/eliona"
	"zevvy/model"
	"zevvy/zevvy"
)

func initialization() {
	ctx := context.Background()

	// Necessary to close used init resources
	conn := db.NewInitConnectionWithContextAndApplicationName(ctx, app.AppName())
	defer conn.Close(ctx)

	// Init the app before the first run.
	app.Init(conn, app.AppName(),
		app.ExecSqlFile("conf/init.sql"),
		asset.InitAssetTypeFiles("resources/asset-types/*.json"),
		dashboard.InitWidgetTypeFiles("resources/widget-types/*.json"),
	)
}

var once sync.Once

func sendData() {
	dbConfigs, err := conf.GetDbConfigs(context.Background())
	if err != nil {
		log.Fatal("conf", "Couldn't read configs from DB: %v", err)
		return
	}
	if len(dbConfigs) == 0 {
		once.Do(func() {
			log.Info("conf", "No configs in DB. Please configure the app in Eliona.")
		})
		return
	}

	for _, dbConfig := range dbConfigs {

		if conf.IsDbConfigEnabled(dbConfig) {

			if !conf.IsDbConfigActive(dbConfig) {
				_, _ = conf.SetDbConfigActiveState(context.Background(), dbConfig.ID, false)
			}

			_, _ = conf.SetDbConfigActiveState(context.Background(), dbConfig.ID, true)
			log.Trace("conf", "Collecting initialized with Configuration %d:\n"+
				"Enable: %t\nRoot URL: %s\nClient ID: %s\nRefresh Interval: %d\nRequest Timeout: %d\n",
				dbConfig.ID, dbConfig.Enable.Bool, dbConfig.APIRootURL, dbConfig.ClientID, dbConfig.RefreshInterval, dbConfig.RequestTimeout)

			// Check for Login process
			if conf.IsLoginNeeded(dbConfig) {
				common.RunOnceWithParam(func(config appdb.Configuration) {
					startLoginProcess(&config)
					time.Sleep(time.Second * time.Duration(config.VerificationInterval.Int32))
				}, *dbConfig, dbConfig.ID)
				continue
			}

			// start working
			common.RunOnceWithParam(func(config appdb.Configuration) {
				log.Info("main", "Collecting %d started.", config.ID)

				// Check for valid access token
				if !conf.IsAccessTokenIsValid(&config) {
					refreshTokens(&config)
				}

				// Send data to Zevvy
				if err := collectData(&config); err != nil {
					return // Error is handled in the method itself.
				}

				log.Info("main", "Collecting %d finished.", config.ID)
				time.Sleep(time.Second * time.Duration(config.RefreshInterval))
			}, *dbConfig, dbConfig.ID)

		}
	}
}

func collectData(dbConfig *appdb.Configuration) error {

	ctx := context.Background()
	dbAssetAttributes, err := conf.GetDbAssetAttributes(ctx, dbConfig.ID)
	if err != nil {
		log.Error("app", "Cannot get asset attributes: %v", err)
		return err
	}

	for _, dbAssetAttribute := range dbAssetAttributes {

		apiDataTrends, err := eliona.GetDataTrends(dbAssetAttribute)
		if err != nil {
			log.Error("ELiona", "Cannot get asset attributes: %v", err)
			return err
		}

		if len(apiDataTrends) > 1 {
			log.Debug("main", "Sending for attribute %d %s %s.", dbAssetAttribute.AssetID, dbAssetAttribute.Subtype, dbAssetAttribute.AttributeName)
		}

		// convert trend data to measurements
		var measurements []model.Measurement
		var latestTimestamp = dbAssetAttribute.LatestTS
		for _, apiDataTrend := range apiDataTrends {
			if apiDataTrend.Timestamp.IsSet() {
				timestamp := common.Val(apiDataTrend.Timestamp.Get())

				// check if data is already sent
				if !timestamp.After(dbAssetAttribute.LatestTS) {
					continue
				}

				// convert data trend to measurement
				measurement := measurementFromTrend(timestamp, apiDataTrend, dbAssetAttribute)
				if measurement.Value != nil {
					log.Debug("main", "Sending data for attribute %s: %d", measurement.ReadAt, *measurement.Value)
					measurements = append(measurements, measurement)
				}

				// remember latest timestamp
				if latestTimestamp.Before(timestamp) {
					latestTimestamp = timestamp
				}

			}
		}

		// send measurements to Zevvy
		if len(measurements) > 0 {
			err := zevvy.SendMeasurements(dbConfig, dbAssetAttribute, measurements)
			if err != nil {
				log.Error("Zevvy", "Cannot send measurements to Zevvy: %v", err)
				return err
			}
		}

		// Store latest timestamp
		err = conf.UpdateAssetAttributeLatestTimestamp(ctx, dbAssetAttribute, latestTimestamp)
		if err != nil {
			log.Error("Conf", "Cannot update latest timestamp: %v", err)
			return err
		}
	}

	return nil
}

func measurementFromTrend(timestamp time.Time, dataTrend api.Data, dbAssetAttribute *appdb.AssetAttribute) model.Measurement {
	const outputFormat = "2006-01-02T15:04:05.000Z"
	measurement := model.Measurement{
		ReadAt: timestamp.Format(outputFormat),
	}
	if value, ok := dataTrend.Data[dbAssetAttribute.AttributeName]; ok {
		if intValue, ok := value.(int); ok {
			measurement.Value = &intValue
		}
		if intValue, ok := value.(int32); ok {
			measurement.Value = common.Ptr(int(intValue))
		}
		if intValue, ok := value.(int64); ok {
			measurement.Value = common.Ptr(int(intValue))
		}
		if intValue, ok := value.(float32); ok {
			measurement.Value = common.Ptr(int(intValue))
		}
		if intValue, ok := value.(float64); ok {
			measurement.Value = common.Ptr(int(intValue))
		}
	}
	return measurement
}

func refreshTokens(dbConfig *appdb.Configuration) {
	log.Info("zevvy", "Get new access token for configuration %d", dbConfig.ID)
	token, err := zevvy.RefreshTokens(dbConfig)
	if err != nil {
		log.Error("zevvy", "Cannot get new token: %v", err)
		return
	}

	log.Info("zevvy", "Update new access and refresh token %d", dbConfig.ID)
	err = conf.UpdateToken(dbConfig, token)
	if err != nil || !conf.IsAccessTokenIsValid(dbConfig) {
		log.Error("zevvy", "Cannot update token in configuration: %v", err)
		return
	}
}

func startLoginProcess(dbConfig *appdb.Configuration) {

	// Get verification URI
	if !conf.IsVerificationUriIsValid(dbConfig) {

		// Get verification URL
		log.Info("zevvy", "Start authentication process for configuration %d", dbConfig.ID)
		log.Info("zevvy", "Get new verification URL for authentication process for configuration %d", dbConfig.ID)
		verification, err := zevvy.GetVerification(dbConfig)
		if err != nil {
			log.Error("zevvy", "Cannot get verification: %v", err)
			return
		}

		err = conf.UpdateVerification(dbConfig, verification)
		if err != nil || !conf.IsVerificationUriIsValid(dbConfig) {
			log.Error("zevvy", "Cannot update verification in configuration: %v", err)
			return
		}

		// Notify user
		log.Info("zevvy", "Notify user about verification URL for authentication process for configuration %d", dbConfig.ID)
		err = eliona.NotifyUser(dbConfig.UserID.String, dbConfig.ProjectID.String, api.Translation{
			De: common.Ptr(fmt.Sprintf("Sie haben die Zevvy-App kürzlich eingerichtet. Um der App den Zugriff auf die Zevvy-API zu ermöglichen, müssen Sie Ihre Anmeldung verifizieren: %s", dbConfig.VerificationURI.String)),
			En: common.Ptr(fmt.Sprintf("You recently set up the Zevvy app. To enable the app's access to the Zevvy API, you must verify your login: %s", dbConfig.VerificationURI.String)),
		})
		if err != nil {
			log.Error("eliona", "Cannot notify user about verification: %v", err)
			return
		}
	}

	// Check if verification is done
	token, err := zevvy.GetTokens(dbConfig)
	if err != nil {
		log.Error("zevvy", "Cannot check verification: %v", err)
		return
	}

	log.Info("zevvy", "Update new access and refresh token %d", dbConfig.ID)
	err = conf.UpdateToken(dbConfig, token)
	if err != nil || !conf.IsAccessTokenIsValid(dbConfig) {
		log.Error("zevvy", "Cannot update token in configuration: %v", err)
		return
	} else {
		// Notify user
		log.Info("zevvy", "Notify user about successful authentication process for configuration %d", dbConfig.ID)
		err = eliona.NotifyUser(dbConfig.UserID.String, dbConfig.ProjectID.String, api.Translation{
			De: common.Ptr(fmt.Sprintf("Zevvy App wurde erfolgreich verifiziert.")),
			En: common.Ptr(fmt.Sprintf("Zevvy app was successful verfied.")),
		})
		if err != nil {
			log.Error("eliona", "Cannot notify user about verification: %v", err)
			return
		}
	}

	return
}

// listenApi starts the API server and listen for requests
func listenApi() {
	err := http.ListenAndServe(":"+common.Getenv("API_SERVER_PORT", "3000"),
		frontend.NewEnvironmentHandler(
			utilshttp.NewCORSEnabledHandler(
				apiserver.NewRouter(
					apiserver.NewConfigurationAPIController(apiservices.NewConfigurationAPIService()),
					apiserver.NewVersionAPIController(apiservices.NewVersionAPIService()),
					apiserver.NewAssetAttributeAPIController(apiservices.NewAssetAttributeAPIService()),
				))))
	log.Fatal("main", "API server: %v", err)
}
