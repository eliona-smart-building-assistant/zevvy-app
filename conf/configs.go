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

package conf

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"zevvy-app/apiserver"
	"zevvy-app/appdb"
	"zevvy-app/model"

	"github.com/eliona-smart-building-assistant/go-eliona/frontend"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var ErrNotFound = errors.New("not found")

func InsertConfig(ctx context.Context, config apiserver.Configuration) (apiserver.Configuration, error) {
	dbConfig, err := dbConfigFromApiConfig(ctx, config)
	if err != nil {
		return apiserver.Configuration{}, fmt.Errorf("creating DB config from API config: %v", err)
	}
	if err := dbConfig.InsertG(ctx, boil.Infer()); err != nil {
		return apiserver.Configuration{}, fmt.Errorf("inserting DB config: %v", err)
	}
	return config, nil
}

func UpsertConfig(ctx context.Context, config apiserver.Configuration) (apiserver.Configuration, error) {
	dbConfig, err := dbConfigFromApiConfig(ctx, config)
	if err != nil {
		return apiserver.Configuration{}, fmt.Errorf("creating DB config from API config: %v", err)
	}
	if err := dbConfig.UpsertG(ctx, true, []string{"id"}, boil.Blacklist("id"), boil.Infer()); err != nil {
		return apiserver.Configuration{}, fmt.Errorf("inserting DB config: %v", err)
	}
	return config, nil
}

func GetConfig(ctx context.Context, configID int64) (apiserver.Configuration, error) {
	dbConfig, err := GetDbConfig(ctx, configID)
	if err != nil {
		return apiserver.Configuration{}, err
	}
	apiConfig, err := apiConfigFromDbConfig(dbConfig)
	if err != nil {
		return apiserver.Configuration{}, fmt.Errorf("creating API config from DB config: %v", err)
	}
	return apiConfig, nil
}

func GetDbConfig(ctx context.Context, configID int64) (*appdb.Configuration, error) {
	dbConfig, err := appdb.Configurations(
		appdb.ConfigurationWhere.ID.EQ(configID),
	).OneG(ctx)
	if errors.Is(err, sql.ErrNoRows) || dbConfig == nil {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("fetching config from database: %v", err)
	}
	return dbConfig, nil
}

func DeleteConfig(ctx context.Context, configID int64) error {
	count, err := appdb.Configurations(
		appdb.ConfigurationWhere.ID.EQ(configID),
	).DeleteAllG(ctx)
	if err != nil {
		return fmt.Errorf("deleting config from database: %v", err)
	}
	if count > 1 {
		return fmt.Errorf("shouldn't happen: deleted more (%v) configs by ID", count)
	}
	if count == 0 {
		return ErrNotFound
	}
	return nil
}

func dbConfigFromApiConfig(ctx context.Context, apiConfig apiserver.Configuration) (dbConfig appdb.Configuration, err error) {
	dbConfig.RootURL = apiConfig.RootUrl
	dbConfig.AuthURLPath = apiConfig.AuthUrlPath
	dbConfig.ClientID = apiConfig.ClientId
	dbConfig.ClientSecret = apiConfig.ClientSecret
	dbConfig.RefreshToken = null.StringFromPtr(apiConfig.RefreshToken)
	dbConfig.ID = null.Int64FromPtr(apiConfig.Id).Int64
	dbConfig.Enable = null.BoolFromPtr(apiConfig.Enable)
	dbConfig.RefreshInterval = apiConfig.RefreshInterval
	if apiConfig.RequestTimeout != nil {
		dbConfig.RequestTimeout = *apiConfig.RequestTimeout
	}
	dbConfig.Active = null.BoolFromPtr(apiConfig.Active)
	env := frontend.GetEnvironment(ctx)
	if env != nil {
		dbConfig.UserID = null.StringFrom(env.UserId)
		dbConfig.ProjectID = null.StringFrom(env.ProjId)
	}
	return dbConfig, nil
}

func apiConfigFromDbConfig(dbConfig *appdb.Configuration) (apiConfig apiserver.Configuration, err error) {
	apiConfig.RootUrl = dbConfig.RootURL
	apiConfig.AuthUrlPath = dbConfig.AuthURLPath
	apiConfig.ClientId = dbConfig.ClientID
	apiConfig.VerificationUri = dbConfig.VerificationURI.Ptr()
	apiConfig.ClientSecret = maskSecret(dbConfig.ClientSecret)
	apiConfig.RefreshToken = common.Ptr(maskSecret(dbConfig.RefreshToken.String))
	apiConfig.Id = &dbConfig.ID
	apiConfig.Enable = dbConfig.Enable.Ptr()
	apiConfig.RefreshInterval = dbConfig.RefreshInterval
	apiConfig.RequestTimeout = &dbConfig.RequestTimeout
	apiConfig.Active = dbConfig.Active.Ptr()
	apiConfig.UserId = dbConfig.UserID.Ptr()
	apiConfig.ProjectId = dbConfig.ProjectID.Ptr()
	return apiConfig, nil
}

func maskSecret(input string) string {
	runes := []rune(input)
	length := len(runes)
	cutoff := length / 3
	for i := cutoff; i < length; i++ {
		runes[i] = '*'
	}
	return string(runes)
}

func IsVerificationUriIsValid(dbConfig *appdb.Configuration) bool {
	return dbConfig.VerificationURI.Valid && len(dbConfig.VerificationURI.String) > 0 && dbConfig.VerificationURIExpire.Time.After(time.Now())
}

func IsAccessTokenIsValid(dbConfig *appdb.Configuration) bool {
	return dbConfig.AccessToken.Valid && len(dbConfig.AccessToken.String) > 0 && dbConfig.AccessTokenExpire.Time.After(time.Now())
}

func GetConfigs(ctx context.Context) ([]apiserver.Configuration, error) {
	dbConfigs, err := appdb.Configurations().AllG(ctx)
	if err != nil {
		return nil, err
	}
	var apiConfigs []apiserver.Configuration
	for _, dbConfig := range dbConfigs {
		ac, err := apiConfigFromDbConfig(dbConfig)
		if err != nil {
			return nil, fmt.Errorf("creating API config from DB config: %v", err)
		}
		apiConfigs = append(apiConfigs, ac)
	}
	return apiConfigs, nil
}

func GetDbConfigs(ctx context.Context) ([]*appdb.Configuration, error) {
	return appdb.Configurations().AllG(ctx)
}

func SetDbConfigActiveState(ctx context.Context, configId int64, state bool) (int64, error) {
	return appdb.Configurations(
		appdb.ConfigurationWhere.ID.EQ(configId),
	).UpdateAllG(ctx, appdb.M{
		appdb.ConfigurationColumns.Active: state,
	})
}

func IsDbConfigActive(dbConfig *appdb.Configuration) bool {
	return !dbConfig.Active.Valid || dbConfig.Active.Bool
}

func IsDbConfigEnabled(config *appdb.Configuration) bool {
	return config.Enable.Valid && config.Enable.Bool
}

func IsLoginNeeded(config *appdb.Configuration) bool {
	return !config.RefreshToken.Valid || len(config.RefreshToken.String) == 0
}

func UpdateVerification(dbConfig *appdb.Configuration, verification *model.Verification) error {
	dbConfig.DeviceCode = null.StringFrom(verification.DeviceCode)
	dbConfig.VerificationURI = null.StringFrom(verification.VerificationUriComplete)
	dbConfig.VerificationURIExpire = null.TimeFrom(time.Now().Add(time.Second * time.Duration(verification.ExpiresIn)))
	dbConfig.VerificationInterval = null.Int32From(verification.Interval)
	_, err := dbConfig.UpdateG(context.Background(), boil.Infer())
	if err != nil {
		return fmt.Errorf("error updating validation information in config %d: %w", dbConfig.ID, err)
	}
	return nil
}

func UpdateToken(dbConfig *appdb.Configuration, token *model.Token) error {
	dbConfig.AccessToken = null.StringFrom(token.AccessToken)
	dbConfig.AccessTokenExpire = null.TimeFrom(time.Now().Add(time.Second * time.Duration(token.ExpiresIn)))
	dbConfig.RefreshToken = null.StringFrom(token.RefreshToken)
	_, err := dbConfig.UpdateG(context.Background(), boil.Infer())
	if err != nil {
		return fmt.Errorf("error updating token information in config %d: %w", dbConfig.ID, err)
	}
	return nil
}
