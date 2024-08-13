package autonomisdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/intercloud/autonomi-sdk/models"
)

func waitForNodeAdministrativeState(ctx context.Context, c *Client, workspaceID, nodeID string, state models.AdministrativeState) bool {
	node, err := c.GetNode(ctx, workspaceID, nodeID)
	if err != nil {
		log.Printf("an error occurs when getting node, err: %s" + err.Error())
		return false
	}
	return node.State == state
}

func (c *Client) CreateNode(ctx context.Context, payload models.CreateNode, workspaceID string, options ...OptionElement) (*models.Node, error) {
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(&payload)
	if err != nil {
		return nil, err
	}

	if errV := c.validate.StructCtx(ctx, payload); errV != nil {
		return nil, errV
	}

	cloudOptions := &elementOptions{}
	for _, o := range options {
		o(cloudOptions)
	}

	if cloudOptions.administrativeState == "" {
		cloudOptions.administrativeState = models.AdministrativeStateDeployed
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/accounts/%s/workspaces/%s/nodes", c.hostURL, c.accountID, workspaceID), body)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	node := models.NodeResponse{}
	err = json.Unmarshal(resp, &node)
	if err != nil {
		return nil, err
	}

	if !c.WaitForAdministrativeState(ctx, workspaceID, node.Data.ID.String(), cloudOptions.administrativeState, waitForNodeAdministrativeState) {
		return nil, fmt.Errorf("Node did not reach '%s' state in time.", cloudOptions.administrativeState)
	}

	return &node.Data, nil
}

func (c *Client) GetNode(ctx context.Context, workspaceID, nodeID string) (*models.Node, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/accounts/%s/workspaces/%s/nodes/%s", c.hostURL, c.accountID, workspaceID, nodeID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	node := models.NodeResponse{}
	err = json.Unmarshal(resp, &node)
	if err != nil {
		return nil, err
	}

	return &node.Data, err
}

func (c *Client) UpdateNode(ctx context.Context, payload models.UpdateElement, workspaceID, nodeID string) (*models.Node, error) {
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(&payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/accounts/%s/workspaces/%s/nodes/%s", c.hostURL, c.accountID, workspaceID, nodeID), body)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	node := models.NodeResponse{}
	err = json.Unmarshal(resp, &node)
	if err != nil {
		return nil, err
	}

	return &node.Data, err
}

func (c *Client) DeleteNode(ctx context.Context, workspaceID, nodeID string) (*models.Node, error) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/accounts/%s/workspaces/%s/nodes/%s", c.hostURL, c.accountID, workspaceID, nodeID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	node := models.NodeResponse{}
	err = json.Unmarshal(resp, &node)
	if err != nil {
		return nil, err
	}

	return &node.Data, err
}
