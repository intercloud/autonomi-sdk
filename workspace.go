package autonomisdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/intercloud/autonomi-sdk/models"
)

func (c *Client) CreateWorkspace(ctx context.Context, payload models.CreateWorkspace) (*models.Workspace, error) {
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(&payload)
	if err != nil {
		return nil, err
	}

	if errV := c.validate.StructCtx(ctx, payload); errV != nil {
		return nil, errV
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/accounts/%s/workspaces", c.hostURL, c.accountID), body)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	workspace := models.WorkspaceResponse{}
	if err := json.Unmarshal(resp, &workspace); err != nil {
		return nil, err
	}

	return &workspace.Data, nil
}

func (c *Client) ListWorkspaces(ctx context.Context, accountID uuid.UUID) ([]models.Workspace, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/accounts/%s/workspaces", c.hostURL, accountID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	workspaces := models.WorkspacesResponse{}
	err = json.Unmarshal(resp, &workspaces)
	if err != nil {
		return nil, err
	}

	return workspaces.Data, nil
}

func (c *Client) GetWorkspace(ctx context.Context, workspaceID string) (*models.Workspace, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/accounts/%s/workspaces/%s", c.hostURL, c.accountID, workspaceID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	workspace := models.WorkspaceResponse{}
	err = json.Unmarshal(resp, &workspace)
	if err != nil {
		return nil, err
	}

	return &workspace.Data, nil
}

func (c *Client) UpdateWorkspace(ctx context.Context, payload models.UpdateWorkspace, workspaceID string) (*models.Workspace, error) {
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(&payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/accounts/%s/workspaces/%s", c.hostURL, c.accountID, workspaceID), body)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	workspace := models.WorkspaceResponse{}
	if err := json.Unmarshal(resp, &workspace); err != nil {
		return nil, err
	}

	return &workspace.Data, nil
}

func (c *Client) DeleteWorkspace(ctx context.Context, workspaceID string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/accounts/%s/workspaces/%s", c.hostURL, c.accountID, workspaceID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
