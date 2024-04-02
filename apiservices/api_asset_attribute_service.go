/*
 * App template API
 *
 * API to access and configure the app template
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package apiservices

import (
	"context"
	"errors"
	"net/http"
	"zevvy-app/apiserver"
	"zevvy-app/conf"
)

// AssetAttributeAPIService is a service that implements the logic for the AssetAttributeAPIServicer
// This service should implement the business logic for every endpoint for the AssetAttributeAPI API.
// Include any external packages or services that will be required by this service.
type AssetAttributeAPIService struct {
}

// NewAssetAttributeAPIService creates a default api service
func NewAssetAttributeAPIService() apiserver.AssetAttributeAPIServicer {
	return &AssetAttributeAPIService{}
}

// DeleteAssetAttributes - Deletes configured asset attributes
func (s *AssetAttributeAPIService) DeleteAssetAttributes(ctx context.Context, configId int32, assetId int32, subtype string, attributeName string) (apiserver.ImplResponse, error) {
	err := conf.DeleteAssetAttributes(ctx, configId, assetId, subtype, attributeName)
	if err != nil {
		return apiserver.ImplResponse{Code: http.StatusInternalServerError}, err
	}
	return apiserver.ImplResponse{Code: http.StatusNoContent}, nil
}

// GetAssetAttributes - Get configured asset attributes
func (s *AssetAttributeAPIService) GetAssetAttributes(ctx context.Context, configId int32, assetId int32, subtype string, attributeName string) (apiserver.ImplResponse, error) {
	configs, err := conf.GetAssetAttributes(ctx, configId, assetId, subtype, attributeName)
	if err != nil {
		return apiserver.ImplResponse{Code: http.StatusInternalServerError}, err
	}
	return apiserver.Response(http.StatusOK, configs), nil
}

// PutAssetAttribute - Creates or updates a configured asset attribute
func (s *AssetAttributeAPIService) PutAssetAttribute(ctx context.Context, assetAttribute apiserver.AssetAttribute) (apiserver.ImplResponse, error) {
	upserted, err := conf.UpsertAssetAttribute(ctx, assetAttribute)
	if errors.Is(err, conf.ErrNotFound) {
		return apiserver.ImplResponse{Code: http.StatusNotFound}, nil
	}
	if err != nil {
		return apiserver.ImplResponse{Code: http.StatusInternalServerError}, err
	}
	return apiserver.Response(http.StatusOK, upserted), nil
}
