package autonomisdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/intercloud/autonomi-sdk/models"
)

// CreatePhysicalPort creates a physical port in Autonomi platform.
func (c *Client) CreatePhysicalPort(ctx context.Context, payload models.CreatePhysicalPort) (*models.PhysicalPort, error) {
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(&payload)
	if err != nil {
		return nil, err
	}

	if errV := c.validate.StructCtx(ctx, payload); errV != nil {
		return nil, errV
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/accounts/%s/ports", c.hostURL, c.accountID), body)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	physicalPort := models.PhysicalPortSingleResponse{}
	err = json.Unmarshal(resp, &physicalPort)
	if err != nil {
		return nil, err
	}

	return &physicalPort.Data, nil
}

func (c *Client) GetPhysicalPort(ctx context.Context, portID string) (*models.PhysicalPort, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/accounts/%s/ports/%s", c.hostURL, c.accountID, portID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	physicalPort := models.PhysicalPortSingleResponse{}
	err = json.Unmarshal(resp, &physicalPort)
	if err != nil {
		return nil, err
	}

	return &physicalPort.Data, err
}

func (c *Client) ListPort(options ...OptionElement) (*[]models.PhysicalPort, error) {

	// retrieve options from request
	portOptions := &elementOptions{}
	for _, o := range options {
		o(portOptions)
	}

	// run request
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/accounts/%s/ports", c.hostURL, c.accountID), nil)
	if err != nil {
		return nil, err
	}

	// add query param if needed
	if portOptions.administrativeState != "" {
		q := req.URL.Query()
		q.Add("state", portOptions.administrativeState.String())
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	ports := models.PhysicalPortListResponse{}
	err = json.Unmarshal(resp, &ports)
	if err != nil {
		return nil, err
	}

	return &ports.Data, err
}

// DeletePhysicalPort creates a physical port in Autonomi platform. As
func (c *Client) DeletePhysicalPort(ctx context.Context, physicalPortID string) error {

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/accounts/%s/ports/%s", c.hostURL, c.accountID, physicalPortID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
