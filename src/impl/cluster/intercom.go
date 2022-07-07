package cluster

import (
	"context"
	"net/http"
	"sync"

	"github.com/shivanshkc/rosenbridge/src/core/models"
)

// Intercom implements deps.Intercom interface using HTTP.
type Intercom struct {
	// httpClients persists the http clients for nodes. So, we can reuse existing TCP connections.
	httpClients map[string]*http.Client
	// httpClientsMutex makes the httpClients map thread-safe to use.
	httpClientsMutex *sync.RWMutex
}

// NewIntercom is a constructor for *Intercom.
func NewIntercom() *Intercom {
	return nil
}

func (i *Intercom) PostMessageInternal(ctx context.Context, nodeAddr string, params *models.OutgoingMessageInternalReq,
) (*models.OutgoingMessageInternalRes, error) {
	panic("implement me")
}
