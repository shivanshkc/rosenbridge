package cluster

import (
	"context"

	"github.com/shivanshkc/rosenbridge/src/core/models"
)

// Intercom implements deps.Intercom interface using HTTP.
type Intercom struct{}

// NewIntercom is a constructor for *Intercom.
func NewIntercom() *Intercom {
	return nil
}

func (i *Intercom) PostMessageInternal(ctx context.Context, nodeAddr string, params *models.OutgoingMessageInternalReq,
) (*models.OutgoingMessageInternalRes, error) {
	panic("implement me")
}
