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

package zevvy

import (
	"fmt"
	utilshttp "github.com/eliona-smart-building-assistant/go-utils/http"
	"github.com/eliona-smart-building-assistant/go-utils/log"
	"net/http"
	"net/url"
	"time"
	"zevvy/appdb"
	"zevvy/model"
)

func GetVerification(dbConfig *appdb.Configuration) (*model.Verification, error) {
	fullUrl := dbConfig.AuthRootURL + "/protocol/openid-connect/auth/device"
	request, err := utilshttp.NewPostFormRequestWithHeaders(
		fullUrl,
		map[string][]string{
			"client_id":     {dbConfig.ClientID},
			"client_secret": {dbConfig.ClientSecret},
			"scope":         {"offline_access measurement register device"},
		}, map[string]string{},
	)
	if err != nil {
		return nil, err
	}
	verification, statusCode, err := utilshttp.ReadWithStatusCode[*model.Verification](request, time.Duration(dbConfig.RequestTimeout)*time.Second, true)
	if err != nil || statusCode != http.StatusOK {
		return nil, fmt.Errorf("error reading request for %s: %d %w", fullUrl, statusCode, err)
	}

	return verification, nil
}

func GetTokens(dbConfig *appdb.Configuration) (*model.Token, error) {
	fullUrl := dbConfig.AuthRootURL + "/protocol/openid-connect/token"
	request, err := utilshttp.NewPostFormRequestWithHeaders(
		fullUrl,
		map[string][]string{
			"client_id":     {dbConfig.ClientID},
			"client_secret": {dbConfig.ClientSecret},
			"device_code":   {dbConfig.DeviceCode.String},
			"grant_type":    {"urn:ietf:params:oauth:grant-type:device_code"},
		}, map[string]string{},
	)
	if err != nil {
		return nil, err
	}
	token, statusCode, err := utilshttp.ReadWithStatusCode[*model.Token](request, time.Duration(dbConfig.RequestTimeout)*time.Second, true)
	if token.Error != nil {
		return nil, fmt.Errorf("status reading request for config %d: %s %s", dbConfig.ID, *token.Error, token.ErrorDescription)
	}
	if err != nil || statusCode != http.StatusOK {
		return nil, fmt.Errorf("error reading request for %s: %d %w", fullUrl, statusCode, err)
	}

	return token, nil
}

func RefreshTokens(dbConfig *appdb.Configuration) (*model.Token, error) {
	fullUrl := dbConfig.AuthRootURL + "/protocol/openid-connect/token"
	request, err := utilshttp.NewPostFormRequestWithHeaders(
		fullUrl,
		map[string][]string{
			"client_id":     {dbConfig.ClientID},
			"client_secret": {dbConfig.ClientSecret},
			"device_code":   {dbConfig.DeviceCode.String},
			"grant_type":    {"refresh_token"},
			"refresh_token": {dbConfig.RefreshToken.String},
		}, map[string]string{},
	)
	if err != nil {
		return nil, err
	}
	token, statusCode, err := utilshttp.ReadWithStatusCode[*model.Token](request, time.Duration(dbConfig.RequestTimeout)*time.Second, true)
	if token != nil && token.Error != nil {
		return nil, fmt.Errorf("status reading request for config %d: %s %s", dbConfig.ID, *token.Error, token.ErrorDescription)
	}
	if err != nil || statusCode != http.StatusOK {
		return nil, fmt.Errorf("error reading request for %s: %d %w", fullUrl, statusCode, err)
	}

	return token, nil
}

func SendMeasurements(dbConfig *appdb.Configuration, dbAssetAttribute *appdb.AssetAttribute, measurements []model.Measurement) error {
	fullUrl := dbConfig.APIRootURL + fmt.Sprintf("/deviceRef/%s/registerRef/%s/measurements/_bulk_create", url.PathEscape(dbAssetAttribute.DeviceReference), url.PathEscape(dbAssetAttribute.RegisterReference))
	request, err := utilshttp.NewPostRequestWithBearer(fullUrl, measurements, dbConfig.AccessToken.String)
	if err != nil {
		return err
	}
	_, statusCode, err := utilshttp.ReadWithStatusCode[any](request, time.Duration(dbConfig.RequestTimeout)*time.Second, true)
	if err != nil || (statusCode != http.StatusCreated && statusCode != http.StatusConflict) {
		return fmt.Errorf("error reading request for %s: %d %w", fullUrl, statusCode, err)
	}
	if statusCode == http.StatusConflict {
		log.Warn("Zevvy", "Conflict for %s", fullUrl)
	}
	return nil
}
