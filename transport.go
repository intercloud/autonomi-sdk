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

func checkTransportAdministrativeState(ctx context.Context, c *Client, workspaceID, transportID string, state models.AdministrativeState) (*models.Transport, bool) {
	transport, err := c.GetTransport(ctx, workspaceID, transportID)
	if err != nil {
		// if wanted state is deleted and the transport is in this state, api has returned 404
		if state == models.AdministrativeStateDeleted && strings.Contains(err.Error(), "status: 404") {
			return nil, true
		}
		log.Printf("an error occurs when getting transport, err: %s" + err.Error())
		return nil, false
	}
	return transport, transport.State == state
}

// CreateTransport creates asynchronously a transport. The transport returned will depend of the administrative state passed in options.
// If none is passed models.AdministrativeStateDeployed will be set by default.
// The valid administrative states options for a transport creation are [models.AdministrativeStateCreationPending, models.AdministrativeStateCreationProceed, models.AdministrativeStateCreationError, models.AdministrativeStateDeployed]
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

	if transportOptions.administrativeState == "" {
		transportOptions.administrativeState = models.AdministrativeStateDeployed
	}

	if _, ok := validCreationAdministrativeStates[transportOptions.administrativeState]; !ok {
		return nil, ErrCreationAdministrativeState
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

	transportPolled, success := WaitForAdministrativeState(ctx, c, workspaceID, transport.Data.ID.String(), transportOptions.administrativeState, checkTransportAdministrativeState)
	if !success {
		return nil, fmt.Errorf("Transport did not reach '%s' state in time.", transportOptions.administrativeState)
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

func (c *Client) DeleteTransport(ctx context.Context, workspaceID, transportID string, options ...OptionElement) (*models.Transport, error) {
	transportOptions := &elementOptions{}
	for _, o := range options {
		o(transportOptions)
	}

	if transportOptions.administrativeState == "" {
		transportOptions.administrativeState = models.AdministrativeStateDeleted
	}

	if _, ok := validDeletionAdministrativeStates[transportOptions.administrativeState]; !ok {
		return nil, ErrDeletionAdministrativeState
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

	transportPolled, success := WaitForAdministrativeState(ctx, c, workspaceID, transport.Data.ID.String(), transportOptions.administrativeState, checkTransportAdministrativeState)
	if !success {
		return nil, fmt.Errorf("Transport did not reach '%s' state in time.", transportOptions.administrativeState)
	}

	// If the transport was deleted and we were waiting for the "deleted" state,
	// transportPolled will be nil. To prevent a panic when dereferencing, we
	// assign an empty structure.
	if transportPolled == nil {
		transportPolled = &models.Transport{}
	}

	return transportPolled, nil
}
