package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/shivanshkc/rosenbridge/src/core"
	"github.com/shivanshkc/rosenbridge/src/utils/httputils"

	"golang.org/x/sync/errgroup"
)

// ResolverCloudRun implements the core.DiscoveryAddressResolver interface assuming the service is running on Cloud Run.
type ResolverCloudRun struct {
	// discoveryAddr is the resolved discovery address.
	discoveryAddr string
	// discoveryOnce ensures that the address is resolved only once.
	discoveryOnce *sync.Once
	// httpClient is used to make all the required calls (all of them are to GCP APIs).
	httpClient *http.Client
}

// NewResolverCloudRun is a constructor for *ResolverCloudRun.
func NewResolverCloudRun() core.DiscoveryAddressResolver {
	return &ResolverCloudRun{
		discoveryAddr: "",
		discoveryOnce: &sync.Once{},
		httpClient:    &http.Client{},
	}
}

// Read returns the persisted discovery address.
func (r *ResolverCloudRun) Read(ctx context.Context) (string, error) {
	var err error

	// Address will be resolved upon the first call.
	r.discoveryOnce.Do(func() {
		err = r.resolve(ctx)
	})

	return r.discoveryAddr, err
}

// Resolve calls the required GCP APIs to figure out the service's address and persists it.
//
//nolint:funlen // This function makes a number of parallel API calls. So, its length is acceptable.
func (r *ResolverCloudRun) resolve(ctx context.Context) error {
	// The K_SERVICE env var is set automatically inside the GCP machine.
	kService := os.Getenv("K_SERVICE") //nolint:revive // Leading k is acceptable here.
	if kService == "" {
		return errors.New("K_SERVICE env var is empty")
	}

	// Creating an err-group to manage multiple goroutines.
	eGroup, eCtx := errgroup.WithContext(ctx)

	// Creating vars to store the network call results.
	var projectID, region, token string

	// Project ID goroutine.
	eGroup.Go(func() error {
		pID, err := r.getProjectID(eCtx)
		if err != nil {
			return fmt.Errorf("error in getProjectID call: %w", err)
		}

		projectID = pID
		return nil
	})

	// Region goroutine.
	eGroup.Go(func() error {
		reg, err := r.getRegion(eCtx)
		if err != nil {
			return fmt.Errorf("error in getRegion call: %w", err)
		}

		region = reg
		return nil
	})

	// Region goroutine.
	eGroup.Go(func() error {
		tkn, err := r.getToken(eCtx)
		if err != nil {
			return fmt.Errorf("error in getToken call: %w", err)
		}

		token = tkn
		return nil
	})

	// Awaiting goroutine completions.
	if err := eGroup.Wait(); err != nil {
		return fmt.Errorf("error in eGroup.Wait call: %w", err)
	}

	// Forming the API route.
	endpoint := fmt.Sprintf("https://%s-run.googleapis.com/apis/serving.knative.dev/v1/namespaces/%s/services/%s",
		region, projectID, kService)

	// Forming the HTTP request.
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("error in http.NewRequestWithContext call: %w", err)
	}

	// Setting the auth header.
	request.Header.Set("authorization", fmt.Sprintf("Bearer %s", token))

	// Executing the request.
	response, err := r.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("error in client.Do call: %w", err)
	}

	// Closing response body upon function return.
	defer func() { _ = response.Body.Close() }()

	// Handling unsuccessful responses.
	if httputils.Is2xx(response.StatusCode) {
		return fmt.Errorf("response has unsuccessful status: %d", response.StatusCode)
	}

	// This is the expected response body structure.
	bodyStruct := &getServiceResponse{}
	// Decoding the response body into the intended struct.
	if err := json.NewDecoder(response.Body).Decode(bodyStruct); err != nil {
		return fmt.Errorf("error in json.NewDecoder(...).Decode call: %w", err)
	}

	// Setting the discovery address.
	r.discoveryAddr = bodyStruct.Status.URL

	return nil
}

// getProjectID queries the GCP VM metadata API to get the project ID of this Cloud Run instance.
func (r *ResolverCloudRun) getProjectID(ctx context.Context) (string, error) {
	// Forming the API route.
	endpoint := fmt.Sprintf("%s%s", gcpBaseURL, gcpProjectIDURL)

	// Forming the HTTP request.
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("error in http.NewRequestWithContext call: %w", err)
	}
	// Setting GCP headers.
	request.Header.Set(gcpHeaderKey, gcpHeaderValue)

	// Executing the request.
	response, err := r.httpClient.Do(request)
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
	projectIDBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error in ioutil.Readall call: %w", err)
	}

	return string(projectIDBytes), nil
}

// getRegion queries the GCP VM metadata API to get the region of this Cloud Run instance.
func (r *ResolverCloudRun) getRegion(ctx context.Context) (string, error) {
	// Forming the API route.
	endpoint := fmt.Sprintf("%s%s", gcpBaseURL, gcpRegionURL)

	// Forming the HTTP request.
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("error in http.NewRequestWithContext call: %w", err)
	}
	// Setting GCP headers.
	request.Header.Set(gcpHeaderKey, gcpHeaderValue)

	// Executing the request.
	response, err := r.httpClient.Do(request)
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
	regionBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error in ioutil.Readall call: %w", err)
	}

	// Region is of the form: "projects/948115683669/regions/us-central1".
	// So, we split it by "/" character and then take the last element.
	regionElements := strings.Split(string(regionBytes), "/")
	// Retuning the last element of the whole region string.
	return regionElements[len(regionElements)-1], nil
}

// getToken queries the GCP VM metadata API to get the access token that's required to hit the service status API.
func (r *ResolverCloudRun) getToken(ctx context.Context) (string, error) {
	// Forming the API route.
	endpoint := fmt.Sprintf("%s%s", gcpBaseURL, gcpTokenURL)

	// Forming the HTTP request.
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("error in http.NewRequestWithContext call: %w", err)
	}
	// Setting GCP headers.
	request.Header.Set(gcpHeaderKey, gcpHeaderValue)

	// Executing the request.
	response, err := r.httpClient.Do(request)
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
	bodyStruct := &tokenResponse{}
	// Decoding the response body into the intended struct.
	if err := json.NewDecoder(response.Body).Decode(bodyStruct); err != nil {
		return "", fmt.Errorf("error in json.NewDecoder(...).Decode call: %w", err)
	}

	// Returning the access token.
	return bodyStruct.AccessToken, nil
}
