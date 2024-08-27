package autonomisdk

import (
	"context"
	"time"

	"github.com/intercloud/autonomi-sdk/models"
)

type elementOptions struct {
	waitUntilElementDeployed   bool
	waitUntilElementUndeployed bool
	administrativeState        models.AdministrativeState
}
type OptionElement func(*elementOptions)

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

// WithAdministrativeState allows setting a specific administrative state.
func WithAdministrativeState(administrativeState models.AdministrativeState) OptionElement {
	return func(c *elementOptions) {
		c.administrativeState = administrativeState
	}
}

type Element interface {
	*models.Node | *models.Transport | *models.Attachment

	GetState() models.AdministrativeState
}

func WaitUntilFinishedTask[T Element](ctx context.Context, client *Client, workspaceID, elementID string, waiterOptions models.AdministrativeState, funcToCall func(context.Context, *Client, string, string, models.AdministrativeState) (T, bool)) (T, bool) {
	var lastElement T
	for i := 0; i < client.poll.maxRetry; i++ {
		// Retrieve the element and check if it is in the required administrative state
		element, finishedTask := funcToCall(ctx, client, workspaceID, elementID, waiterOptions)
		lastElement = element
		if finishedTask {
			return lastElement, true
		}
		if lastElement != nil && (lastElement.GetState() == models.AdministrativeStateCreationError || lastElement.GetState() == models.AdministrativeStateDeleteError) {
			return lastElement, false
		}
		time.Sleep(client.poll.retryInterval)
	}

	return lastElement, false
}
