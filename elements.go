package autonomisdk

import (
	"context"
	"time"

	"github.com/intercloud/autonomi-sdk/models"
)

type elementOptions struct {
	waitUntilElementDeployed   bool
	waitUntilElementUndeployed bool
}
type OptionElement func(*elementOptions)

var (
	waitUntilElementDeployed   = []models.AdministrativeState{models.AdministrativeStateCreationError, models.AdministrativeStateDeployed}
	waitUntilElementUndeployed = []models.AdministrativeState{models.AdministrativeStateDeleteError, models.AdministrativeStateDeleted}
)

func WithWaitUntilElementDeployed() OptionElement {
	return func(e *elementOptions) {
		e.waitUntilElementDeployed = true
	}
}

func WithWaitUntilElementUndeployed() OptionElement {
	return func(e *elementOptions) {
		e.waitUntilElementUndeployed = true
	}
}

type Element interface {
	*models.Node | *models.Transport | *models.Attachment
}

func WaitUntilFinishedTask[T Element](ctx context.Context, client *Client, workspaceID, elementID string, waiterOptions []models.AdministrativeState, funcToCall func(context.Context, *Client, string, string, []models.AdministrativeState) (T, bool)) (T, bool) {
	var lastElement T
	for i := 0; i < client.poll.maxRetry; i++ {
		// Retrieve the element and check if it is in the required administrative state
		element, isInDesiredState := funcToCall(ctx, client, workspaceID, elementID, waiterOptions)
		lastElement = element

		if isInDesiredState {
			return lastElement, true
		}
		time.Sleep(client.poll.retryInterval)
	}

	return lastElement, false
}
