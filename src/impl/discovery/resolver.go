package discovery

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/shivanshkc/rosenbridge/src/utils/httputils"
)

const (
	gcpMetadataBaseURL = "http://metadata.google.internal/computeMetadata/v1"

	gcpProjectIDURL = "/project/project-id"
	gcpRegionURL    = "/instance/region"
	gcpTokenURL     = "/instance/service-accounts/default/token"

	gcpHeaderKey   = "Metadata-Flavor"
	gcpHeaderValue = "Google"
)

// ResolverCloudRun implements the deps.DiscoveryAddressResolver interface assuming the service is running on Cloud Run.
type ResolverCloudRun struct {
	// discoveryAddr is the resolved discovery address.
	discoveryAddr string
}

// NewResolver is a constructor for *Resolver.
func NewResolver() *ResolverCloudRun {
	return nil
}

func (r *ResolverCloudRun) Resolve() string {
	return r.discoveryAddr
}

func (r *ResolverCloudRun) GetProjectID(ctx context.Context) (string, error) {
	endpoint := fmt.Sprintf("%s%s", gcpMetadataBaseURL, gcpProjectIDURL)
	// Forming the HTTP request.
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("error in http.NewRequestWithContext call: %w", err)
	}

	// Setting GCP headers.
	request.Header.Set(gcpHeaderKey, gcpHeaderValue)
	// Creating the HTTP client.
	client := &http.Client{}
	// Executing the request.
	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("error in client.Do call: %w", err)
	}
	// Closing response body upon function return.
	defer func() { _ = response.Body.Close() }()
	// Handling unsuccessful responses.
	if httputils.Is2xx(response.StatusCode) {
		return "", fmt.Errorf("response has unsuccessful status: %d", response.StatusCode)
	}

	// Reading the response body into a byte slice.
	projectIDBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error in ioutil.Readall call: %w", err)
	}

	return string(projectIDBytes), nil
}
