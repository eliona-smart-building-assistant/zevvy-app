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
	"fmt"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"time"
	"zevvy-app/apiserver"
	"zevvy-app/appdb"
	"zevvy-app/eliona"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

func UpsertAssetAttribute(ctx context.Context, apiAssetAttribute apiserver.AssetAttribute) (apiserver.AssetAttribute, error) {
	dbAssetAttribute := dbAssetAttributeFromApiAssetAttribute(&apiAssetAttribute)
	apiAsset, err := eliona.GetAsset(dbAssetAttribute)
	if err != nil {
		return apiAssetAttribute, fmt.Errorf("Erro getting asset from Eliona: %w", err)
	}
	if apiAsset == nil {
		return apiAssetAttribute, ErrNotFound
	}
	if len(dbAssetAttribute.DeviceReference) == 0 {
		dbAssetAttribute.DeviceReference = apiAsset.GlobalAssetIdentifier
	}
	if len(dbAssetAttribute.RegisterReference) == 0 {
		dbAssetAttribute.RegisterReference = apiAssetAttribute.AttributeName
	}

	err = dbAssetAttribute.UpsertG(ctx, true,
		[]string{
			appdb.AssetAttributeColumns.ConfigID,
			appdb.AssetAttributeColumns.AssetID,
			appdb.AssetAttributeColumns.Subtype,
			appdb.AssetAttributeColumns.AttributeName,
		},
		boil.Whitelist(
			appdb.AssetAttributeColumns.DeviceReference,
			appdb.AssetAttributeColumns.RegisterReference,
			appdb.AssetAttributeColumns.LatestTS,
		),
		boil.Whitelist(
			appdb.AssetAttributeColumns.ConfigID,
			appdb.AssetAttributeColumns.AssetID,
			appdb.AssetAttributeColumns.Subtype,
			appdb.AssetAttributeColumns.AttributeName,
			appdb.AssetAttributeColumns.DeviceReference,
			appdb.AssetAttributeColumns.RegisterReference,
			appdb.AssetAttributeColumns.LatestTS,
		),
	)
	if err != nil {
		return apiAssetAttribute, err
	}
	return apiAssetAttribute, nil
}

func UpdateAssetAttributeLatestTimestamp(ctx context.Context, dbAssetAttribute *appdb.AssetAttribute, latestTimestamp time.Time) error {
	dbAssetAttribute.LatestTS = latestTimestamp
	_, err := dbAssetAttribute.UpdateG(ctx, boil.Whitelist(appdb.AssetAttributeColumns.LatestTS))
	return err
}

func GetDbAssetAttributes(ctx context.Context, configId int64) (dbAssetAttributes []*appdb.AssetAttribute, err error) {
	return appdb.AssetAttributes(appdb.AssetAttributeWhere.ConfigID.EQ(int32(configId))).AllG(ctx)
}

func GetAssetAttributes(ctx context.Context, configId int32, assetId int32, subtype string, attributeName string) (apiAssetAttributes []*apiserver.AssetAttribute, err error) {
	mods := selectAssetAttributesMods(configId, assetId, subtype, attributeName)
	dbAssetAttributes, err := appdb.AssetAttributes(mods...).AllG(ctx)
	if err != nil {
		return nil, err
	}
	for _, dbAssetAttribute := range dbAssetAttributes {
		apiAssetAttributes = append(apiAssetAttributes, apiAssetAttributeFromDbAssetAttribute(dbAssetAttribute))
	}
	return apiAssetAttributes, nil
}

func selectAssetAttributesMods(configId int32, assetId int32, subtype string, attributeName string) []qm.QueryMod {
	var mods []qm.QueryMod
	if configId > 0 {
		mods = append(mods, appdb.AssetAttributeWhere.ConfigID.EQ(configId))
	}
	if assetId > 0 {
		mods = append(mods, appdb.AssetAttributeWhere.AssetID.EQ(assetId))
	}
	if len(subtype) > 0 {
		mods = append(mods, appdb.AssetAttributeWhere.Subtype.EQ(subtype))
	}
	if len(attributeName) > 0 {
		mods = append(mods, appdb.AssetAttributeWhere.AttributeName.EQ(attributeName))
	}
	return mods
}

func DeleteAssetAttributes(ctx context.Context, configId int32, assetId int32, subtype string, attributeName string) error {
	mods := selectAssetAttributesMods(configId, assetId, subtype, attributeName)
	_, err := appdb.AssetAttributes(mods...).DeleteAllG(ctx)
	if err != nil {
		return err
	}
	return nil
}

func dbAssetAttributeFromApiAssetAttribute(apiAssetAttribute *apiserver.AssetAttribute) (dbAssetAttribute *appdb.AssetAttribute) {
	if apiAssetAttribute != nil {
		dbAssetAttribute.ConfigID = apiAssetAttribute.ConfigId
		dbAssetAttribute.AssetID = apiAssetAttribute.AssetId
		dbAssetAttribute.Subtype = apiAssetAttribute.Subtype
		dbAssetAttribute.AttributeName = apiAssetAttribute.AttributeName
		dbAssetAttribute.DeviceReference = common.Val(apiAssetAttribute.DeviceReference)
		dbAssetAttribute.RegisterReference = common.Val(apiAssetAttribute.RegisterReference)
	}
	return dbAssetAttribute
}

func apiAssetAttributeFromDbAssetAttribute(dbAssetAttribute *appdb.AssetAttribute) (apiAssetAttribute *apiserver.AssetAttribute) {
	if dbAssetAttribute != nil {
		apiAssetAttribute.ConfigId = dbAssetAttribute.ConfigID
		apiAssetAttribute.AssetId = dbAssetAttribute.AssetID
		apiAssetAttribute.Subtype = dbAssetAttribute.Subtype
		apiAssetAttribute.AttributeName = dbAssetAttribute.AttributeName
		apiAssetAttribute.DeviceReference = common.Ptr(dbAssetAttribute.DeviceReference)
		apiAssetAttribute.RegisterReference = common.Ptr(dbAssetAttribute.RegisterReference)
		apiAssetAttribute.LatestTimestamp = dbAssetAttribute.LatestTS
	}
	return apiAssetAttribute
}
