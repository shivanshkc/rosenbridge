package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/shivanshkc/rosenbridge/src/utils/httputils"

	"golang.org/x/sync/errgroup"
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

func (r *ResolverCloudRun) getProjectID(ctx context.Context) (string, error) {
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

func (r *ResolverCloudRun) getRegion(ctx context.Context) (string, error) {
	endpoint := fmt.Sprintf("%s%s", gcpMetadataBaseURL, gcpRegionURL)
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
	regionBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error in ioutil.Readall call: %w", err)
	}

	// Region is of the form: "projects/948115683669/regions/us-central1".
	// So, we split it by "/" character and then take the last element.
	regionElements := strings.Split(string(regionBytes), "/")
	region := regionElements[len(regionElements)-1]

	return region, nil
}

func (r *ResolverCloudRun) getToken(ctx context.Context) (string, error) {
	endpoint := fmt.Sprintf("%s%s", gcpMetadataBaseURL, gcpTokenURL)
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

	// This is the expected response body structure.
	bodyStruct := &struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}{}

	// Decoding the response body into the intended struct.
	if err := json.NewDecoder(response.Body).Decode(bodyStruct); err != nil {
		return "", fmt.Errorf("error in json.NewDecoder(...).Decode call: %w", err)
	}

	// Returning the access token.
	return bodyStruct.AccessToken, nil
}

func (r *ResolverCloudRun) GetAddress(ctx context.Context) (string, error) {
	kService := os.Getenv("K_SERVICE")
	if kService == "" {
		return "", errors.New("K_SERVICE env var is empty")
	}

	// Creating an err-group to manage multiple goroutines.
	eGroup, eCtx := errgroup.WithContext(ctx)

	// Creating a channel to receive the project ID.
	projectIDChan := make(chan string, 1)
	defer close(projectIDChan)

	// Project ID goroutine.
	eGroup.Go(func() error {
		projectID, err := r.getProjectID(eCtx)
		if err != nil {
			return fmt.Errorf("error in getProjectID call: %w", err)
		}

		projectIDChan <- projectID
		return nil
	})

	// Creating a channel to receive the region.
	regionChan := make(chan string, 1)
	defer close(regionChan)

	// Region goroutine.
	eGroup.Go(func() error {
		region, err := r.getRegion(eCtx)
		if err != nil {
			return fmt.Errorf("error in getRegion call: %w", err)
		}

		regionChan <- region
		return nil
	})

	// Creating a channel to receive the token.
	tokenChan := make(chan string, 1)
	defer close(tokenChan)

	// Region goroutine.
	eGroup.Go(func() error {
		token, err := r.getToken(eCtx)
		if err != nil {
			return fmt.Errorf("error in getToken call: %w", err)
		}

		tokenChan <- token
		return nil
	})

	// Awaiting goroutine completions.
	if err := eGroup.Wait(); err != nil {
		return "", fmt.Errorf("error in eGroup.Wait call: %w", err)
	}

	// Retrieving the values from channels.
	projectID, region, token := <-projectIDChan, <-regionChan, <-tokenChan

	// Forming the API route.
	endpoint := fmt.Sprintf("https://%s-run.googleapis.com/apis/serving.knative.dev/v1/namespaces/%s/services/%s",
		region, projectID, kService)

	// Forming the HTTP request.
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("error in http.NewRequestWithContext call: %w", err)
	}

	// Setting the auth header.
	request.Header.Set("authorization", fmt.Sprintf("Bearer %s", token))

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

	// This is the expected response body structure.
	bodyStruct := &struct {
		Status *struct {
			URL string `json:"url"`
		} `json:"status"`
	}{}

	// Decoding the response body into the intended struct.
	if err := json.NewDecoder(response.Body).Decode(bodyStruct); err != nil {
		return "", fmt.Errorf("error in json.NewDecoder(...).Decode call: %w", err)
	}

	// Returning the access token.
	return bodyStruct.Status.URL, nil
}
