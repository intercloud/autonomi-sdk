package autonomisdk

import (
	"context"
	"time"

	"github.com/intercloud/autonomi-sdk/models"
)

type elementOptions struct {
	administrativeState models.AdministrativeState
}

type OptionElement func(*elementOptions)

// WithAdministrativeState allows setting a specific administrative state. The node will be returned when it reaches the specified administrative state.
func WithAdministrativeState(administrativeState models.AdministrativeState) OptionElement {
	return func(c *elementOptions) {
		c.administrativeState = administrativeState
	}
}

type WaitForAdministrativeStateSignature func(ctx context.Context, c *Client, workspaceID, elementID string, state models.AdministrativeState) bool

func (c *Client) WaitForAdministrativeState(ctx context.Context, workspaceID, elementID string, state models.AdministrativeState, funcToCall WaitForAdministrativeStateSignature) bool {
	for i := 0; i < c.maxRetry; i++ {
		// check if element is in required administrative state, if not the loop continues until it reaches it or it timeout
		if funcToCall(ctx, c, workspaceID, elementID, state) {
			return true
		}

		time.Sleep(c.retryInterval)
	}

	return false
}
