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

// type CheckForAdministrativeState func(ctx context.Context, c *Client, workspaceID, elementID string, state models.AdministrativeState) bool

type Element interface {
	*models.Node | *models.Transport | *models.Attachment
}

func WaitForAdministrativeState[T Element](ctx context.Context, client *Client, workspaceID, elementID string, state models.AdministrativeState, funcToCall func(context.Context, *Client, string, string, models.AdministrativeState) (T, bool)) (T, bool) {
	var lastElement T
	for i := 0; i < client.poll.maxRetry; i++ {
		// Retrieve the element and check if it is in the required administrative state
		element, isInDesiredState := funcToCall(ctx, client, workspaceID, elementID, state)
		lastElement = element

		if isInDesiredState {
			return lastElement, true
		}
		time.Sleep(client.poll.retryInterval)
	}

	return lastElement, false
}
