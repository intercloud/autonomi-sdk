package autonomisdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/intercloud/autonomi-sdk/models"
)

func waitForTransportAdministrativeState(ctx context.Context, c *Client, workspaceID, transportID string, state models.AdministrativeState) bool {
	transport, err := c.GetTransport(ctx, workspaceID, transportID)
	if err != nil {
		return false
	}
	return transport.State == state
}

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

	if !c.WaitForAdministrativeState(ctx, workspaceID, transport.Data.ID.String(), transportOptions.administrativeState, waitForTransportAdministrativeState) {
		return nil, fmt.Errorf("Node did not reach '%s' state in time.", transportOptions.administrativeState)
	}

	return &transport.Data, err
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

func (c *Client) DeleteTransport(ctx context.Context, workspaceID, transportID string) (*models.Transport, error) {
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

	return &transport.Data, err
}
