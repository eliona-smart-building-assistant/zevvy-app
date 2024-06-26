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

package eliona

import (
	"fmt"
	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-eliona/client"
	"zevvy/appdb"
)

func GetDataList(dbAssetAttribute *appdb.AssetAttribute) ([]api.Data, error) {
	dataList, response, err := client.NewClient().DataAPI.GetData(client.AuthenticationContext()).
		AssetId(dbAssetAttribute.AssetID).
		DataSubtype(dbAssetAttribute.Subtype).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("error fetching data trends from ELiona API %d: %w", response.StatusCode, err)
	}
	return dataList, nil
}

func GetAsset(dbAssetAttribute *appdb.AssetAttribute) (*api.Asset, error) {
	asset, _, err := client.NewClient().AssetsAPI.GetAssetById(client.AuthenticationContext(), dbAssetAttribute.AssetID).Execute()
	return asset, err
}
