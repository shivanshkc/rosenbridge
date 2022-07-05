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
)

// Comm is the interface to communicate with the other cluster nodes.
type Comm struct {
	// httpClients persists the http clients for nodes. So, we can reuse existing TCP connections.
	httpClients map[string]*http.Client
	// httpClientsMutex makes the httpClients map thread-safe to use.
	httpClientsMutex *sync.RWMutex
}

// NewComm provides a new Comm instance.
func NewComm() *Comm {
	return &Comm{
		httpClients:      map[string]*http.Client{},
		httpClientsMutex: &sync.RWMutex{},
	}
}

func (c *Comm) PostMessageInternal(ctx context.Context, nodeAddr string, input *core.PostMessageInternalParams,
) (*core.OutgoingMessageRes, error) {
	// Prerequisites.
	conf := configs.Get()

	// Locking the httpClients map for read-write operations.
	c.httpClientsMutex.Lock()
	defer c.httpClientsMutex.Unlock()

	// Checking if there's already a http client for this node, otherwise creating a new one.
	httpClient, exists := c.httpClients[nodeAddr]
	if !exists {
		httpClient = &http.Client{}
		c.httpClients[nodeAddr] = httpClient
	}

	// Converting the body into an io.Reader.
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("error in json.Marshal call: %w", err)
	}
	inputReader := bytes.NewReader(inputBytes)

	// Forming the request endpoint.
	endpoint := fmt.Sprintf("%s://%s/api/internal/message", conf.HTTPServer.DiscoveryProtocol, nodeAddr)

	// Creating the HTTP request.
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, inputReader)
	if err != nil {
		return nil, fmt.Errorf("error in http.NewRequestWithContext call: %w", err)
	}
	// Adding the required headers.
	request.Header.Set("x-request-id", input.RequestID)
	request.Header.Set("x-client-id", input.ClientID)
	// Setting the cluster basic auth params.
	request.SetBasicAuth(conf.Auth.ClusterUsername, conf.Auth.ClusterPassword)

	// Executing the request.
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error in httpClient.Do call: %w", err)
	}
	defer func() { _ = response.Body.Close() }()

	// Unmarshalling the response body into *core.OutgoingMessageRes
	outgoingMessageRes := &core.OutgoingMessageRes{}
	if err := json.NewDecoder(response.Body).Decode(outgoingMessageRes); err != nil {
		return nil, fmt.Errorf("error in json.NewDecoder(...).Decode call: %w", err)
	}

	return outgoingMessageRes, nil
}
