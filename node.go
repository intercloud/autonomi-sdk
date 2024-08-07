package autonomisdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/intercloud/autonomi-sdk/models"
)

func (c *Client) CreateNode(ctx context.Context, payload models.CreateNode, workspaceID string) (*models.Node, error) {
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(&payload)
	if err != nil {
		return nil, err
	}

	if errV := c.validate.StructCtx(ctx, payload); errV != nil {
		return nil, errV
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

	return &node.Data, err
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

func (c *Client) UpdateNode(ctx context.Context, payload models.UpdateNode, workspaceID, nodeID string) (*models.Node, error) {
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
