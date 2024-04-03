/*
 * App template API
 *
 * API to access and configure the app template
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package apiserver

// Configuration - Each configuration defines access to provider's API.
type Configuration struct {

	// Internal identifier for the configured API (created automatically).
	Id *int64 `json:"id,omitempty"`

	// Root URL for the authentication process
	AuthRootUrl string `json:"authRootUrl,omitempty"`

	// Root URL for the API access
	ApiRootUrl string `json:"apiRootUrl,omitempty"`

	// Client ID for API access
	ClientId string `json:"clientId"`

	// Set the client secret for API access
	ClientSecret string `json:"clientSecret"`

	// Login-URL to verify the access to the API
	VerificationUri *string `json:"verificationUri,omitempty"`

	// Optionally set the refresh token. If not provided, it will be automatically assigned during the authentication login process managed by the app.
	RefreshToken *string `json:"refreshToken,omitempty"`

	// Flag to enable or disable fetching from this API
	Enable *bool `json:"enable,omitempty"`

	// Interval in seconds for collecting data from API
	RefreshInterval int32 `json:"refreshInterval,omitempty"`

	// Timeout in seconds
	RequestTimeout *int32 `json:"requestTimeout,omitempty"`

	// Set to `true` by the app when running and to `false` when app is stopped
	Active *bool `json:"active,omitempty"`

	// ID of the last Eliona user who created or updated the configuration
	UserId *string `json:"userId,omitempty"`

	// ID of the project the Eliona user created or updated the configuration
	ProjectId *string `json:"projectId,omitempty"`
}

// AssertConfigurationRequired checks if the required fields are not zero-ed
func AssertConfigurationRequired(obj Configuration) error {
	elements := map[string]interface{}{
		"clientId":     obj.ClientId,
		"clientSecret": obj.ClientSecret,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertConfigurationConstraints checks if the values respects the defined constraints
func AssertConfigurationConstraints(obj Configuration) error {
	return nil
}
