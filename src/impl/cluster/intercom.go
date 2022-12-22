package cluster

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/shivanshkc/rosenbridge/src/configs"
	"github.com/shivanshkc/rosenbridge/src/core"
	"github.com/shivanshkc/rosenbridge/src/utils/httputils"
)

// Intercom implements core.IntercomService interface using HTTP.
type Intercom struct {
	// httpClients persists the http clients for nodes. So, we can reuse existing TCP connections.
	httpClients map[string]*http.Client
	// httpClientsMutex makes the httpClients map thread-safe to use.
	httpClientsMutex *sync.RWMutex
}

// NewIntercom is a constructor for *Intercom.
func NewIntercom() core.IntercomService {
	return &Intercom{
		httpClients:      map[string]*http.Client{},
		httpClientsMutex: &sync.RWMutex{},
	}
}

func (i *Intercom) SendMessageInternal(ctx context.Context, nodeAddr string, params *core.OutgoingMessageInternalReq,
) (*core.OutgoingMessageInternalRes, error) {
	// Prerequisites.
	conf := configs.Get()

	// Obtaining request ID from the provided context.
	requestID := httputils.GetReqCtx(ctx).ID

	// Converting the body into a byte slice.
	reqBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("error in json.Marshal call: %w", err)
	}
	// Converting the byte slice into an io.Reader.
	reqReader := bytes.NewReader(reqBytes)

	// Forming the request endpoint.
	endpoint := fmt.Sprintf("%s/api/internal/message", nodeAddr)
	// Creating the HTTP request.
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, reqReader)
	if err != nil {
		return nil, fmt.Errorf("error in http.NewRequestWithContext call: %w", err)
	}
	// Adding the required headers.
	request.Header.Set("x-request-id", requestID)
	// Setting the cluster basic auth params.
	request.SetBasicAuth(conf.Auth.InternalUsername, conf.Auth.InternalPassword)

	// Locking the httpClients map for read-write operations.
	i.httpClientsMutex.Lock()

	// Checking if there's already a http client for this node, otherwise creating a new one.
	httpClient, exists := i.httpClients[nodeAddr]
	if !exists {
		httpClient = &http.Client{}
		i.httpClients[nodeAddr] = httpClient
	}

	// Unlocking the httpClients map.
	i.httpClientsMutex.Unlock()

	// Executing the request.
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error in httpClient.Do call: %w", err)
	}
	// Closing the response body upon function return.
	defer func() { _ = response.Body.Close() }()

	// Unmarshalling the response body into OutgoingMessageInternalRes.
	responseBody := &core.OutgoingMessageInternalRes{}
	if err := json.NewDecoder(response.Body).Decode(responseBody); err != nil {
		return nil, fmt.Errorf("error in json.NewDecoder(...).Decode call: %w", err)
	}

	return responseBody, nil
}
