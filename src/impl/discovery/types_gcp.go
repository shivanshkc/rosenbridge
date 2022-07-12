package discovery

const (
	gcpBaseURL = "http://metadata.google.internal/computeMetadata/v1"

	gcpProjectIDURL = "/project/project-id"
	gcpRegionURL    = "/instance/region"
	gcpTokenURL     = "/instance/service-accounts/default/token"

	gcpHeaderKey   = "Metadata-Flavor"
	gcpHeaderValue = "Google"
)

// tokenResponse is the schema of response body of the GCP access token endpoint.
type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// getServiceResponse is the schema of response body of the GCP namespaces.services.get API.
type getServiceResponse struct {
	Status *struct {
		URL string `json:"url"`
	} `json:"status"`
}
