package autonomisdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/intercloud/autonomi-sdk/models"
)

func checkTransportFinishedTask(ctx context.Context, c *Client, workspaceID, transportID string, waiterOptionState models.AdministrativeState) (*models.Transport, bool) {
	transport, err := c.GetTransport(ctx, workspaceID, transportID)
	if err != nil {
		// if wanted state is deleted and the attachment is in this state, api has returned 404
		if waiterOptionState == models.AdministrativeStateDeleted {
			if strings.Contains(err.Error(), "status: 404") {
				return nil, true
			}
		}
		log.Printf("an error occurs when getting node, err: %s" + err.Error())
		return nil, false
	}

	return transport, waiterOptionState == transport.State
}

// CreateTransport creates asynchronously a transport. The transport returned will depend of the passed option.
// If none is passed the transport will be returned once created in database with administrative state creation_pending.
// If the option WithWaitUntilElementDeployed() is passed, the transport will be returned when its state reach deployed or creation_error.
// If the option WithWaitUntilElementUndeployed() is passed, it will not be considered hence the transport returned will be in state creation_pending
func (c *Client) CreateTransport(ctx context.Context, payload models.CreateTransport, workspaceID string, options ...OptionElement) (*models.Transport, error) {
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(&payload)
	if err != nil {
		return nil, err
	}

	if errV := c.validate.StructCtx(ctx, payload); errV != nil {
		return nil, errV
	}

	transportOptions := &elementOptions{}
	for _, o := range options {
		o(transportOptions)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/accounts/%s/workspaces/%s/transports", c.hostURL, c.accountID, workspaceID), body)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	transport := models.TransportResponse{}
	err = json.Unmarshal(resp, &transport)
	if err != nil {
		return nil, err
	}

	var transportPolled = &transport.Data
	if transportOptions.waitUntilElementDeployed {
		var success bool
		transportPolled, success = WaitUntilFinishedTask(ctx, c, workspaceID, transport.Data.ID.String(), models.AdministrativeStateDeployed, checkTransportFinishedTask)
		if !success {
			return nil, fmt.Errorf("Transport did not reach '%s' state in time.", models.AdministrativeStateDeployed)
		}
	}

	return transportPolled, nil
}

func (c *Client) GetTransport(ctx context.Context, workspaceID, transportID string) (*models.Transport, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/accounts/%s/workspaces/%s/transports/%s", c.hostURL, c.accountID, workspaceID, transportID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	transport := models.TransportResponse{}
	err = json.Unmarshal(resp, &transport)
	if err != nil {
		return nil, err
	}

	return &transport.Data, err
}

func (c *Client) UpdateTransport(ctx context.Context, payload models.UpdateElement, workspaceID, transportID string) (*models.Transport, error) {
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(&payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/accounts/%s/workspaces/%s/transports/%s", c.hostURL, c.accountID, workspaceID, transportID), body)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	transport := models.TransportResponse{}
	err = json.Unmarshal(resp, &transport)
	if err != nil {
		return nil, err
	}

	return &transport.Data, err
}

// DeleteTransport deletes asynchronously a transport. The transport returned will depend of the option passed.
// If none is passed the transport will be returned once the request accepted, its state will be delete_pending
// If the option WithWaitUntilElementUndeployed() is passed, the transport won't be returned as it would have been deleted. However, if an error is triggered, an object could be returned with a delete_error state.
// If the option WithWaitUntilElementDeployed() is passed, it will not be considered hence the transport returned will be in state delete_pending
func (c *Client) DeleteTransport(ctx context.Context, workspaceID, transportID string, options ...OptionElement) (*models.Transport, error) {
	transportOptions := &elementOptions{}
	for _, o := range options {
		o(transportOptions)
	}

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/accounts/%s/workspaces/%s/transports/%s", c.hostURL, c.accountID, workspaceID, transportID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	transport := models.TransportResponse{}
	err = json.Unmarshal(resp, &transport)
	if err != nil {
		return nil, err
	}

	var transportPolled = &transport.Data
	if transportOptions.waitUntilElementUndeployed {
		var success bool
		transportPolled, success = WaitUntilFinishedTask(ctx, c, workspaceID, transport.Data.ID.String(), models.AdministrativeStateDeleted, checkTransportFinishedTask)
		if !success {
			return nil, fmt.Errorf("Transport did not reach '%s' state in time.", models.AdministrativeStateDeleted)
		}
	}

	// If the transport was deleted and we were waiting for the "deleted" state,
	// transportPolled will be nil. To prevent a panic when dereferencing, we
	// assign an empty structure.
	if transportPolled == nil {
		transportPolled = &models.Transport{}
	}

	return transportPolled, nil
}
