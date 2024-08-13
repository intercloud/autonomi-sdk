package autonomisdk

import (
	"context"
	"errors"
	"time"

	"github.com/intercloud/autonomi-sdk/models"
)

type elementOptions struct {
	administrativeState models.AdministrativeState
}

type OptionElement func(*elementOptions)

var validCreationAdministrativeStates = map[models.AdministrativeState]struct{}{
	models.AdministrativeStateCreationPending: {},
	models.AdministrativeStateCreationProceed: {},
	models.AdministrativeStateCreationError:   {},
	models.AdministrativeStateDeployed:        {},
}

var validDeletionAdministrativeStates = map[models.AdministrativeState]struct{}{
	models.AdministrativeStateDeletePending: {},
	models.AdministrativeStateDeleteProceed: {},
	models.AdministrativeStateDeleteError:   {},
	models.AdministrativeStateDeleted:       {},
}

var (
	ErrCreationAdministrativeState = errors.New("not a valid administrative state for element creation. Try [AdministrativeStateCreationPending, AdministrativeStateCreationProceed, AdministrativeStateCreationError, AdministrativeStateDeployed ]")
	ErrDeletionAdministrativeState = errors.New("not a valid administrative state for element deletion. Try [AdministrativeStateDeletePending, AdministrativeStateDeleteProceed, AdministrativeStateDeleteError, AdministrativeStateDeleted ]")
)

// WithAdministrativeState allows setting a specific administrative state. The node will be returned when it reaches the specified administrative state.
func WithAdministrativeState(administrativeState models.AdministrativeState) OptionElement {
	return func(c *elementOptions) {
		c.administrativeState = administrativeState
	}
}

type CheckForAdministrativeStateSignature func(ctx context.Context, c *Client, workspaceID, elementID string, state models.AdministrativeState) bool

func (c *Client) WaitForAdministrativeState(ctx context.Context, workspaceID, elementID string, state models.AdministrativeState, funcToCall CheckForAdministrativeStateSignature) bool {
	for i := 0; i < c.poll.maxRetry; i++ {
		// check if element is in required administrative state, if not the loop continues until it reaches it or it timeout
		if funcToCall(ctx, c, workspaceID, elementID, state) {
			return true
		}

		time.Sleep(c.poll.retryInterval)
	}

	return false
}
