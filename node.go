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

func checkNodeAdministrativeState(ctx context.Context, c *Client, workspaceID, nodeID string, state models.AdministrativeState) (*models.Node, bool) {
	node, err := c.GetNode(ctx, workspaceID, nodeID)
	if err != nil {
		// if wanted state is deleted and the node is in this state, api has returned 404
		if state == models.AdministrativeStateDeleted && strings.Contains(err.Error(), "status: 404") {
			return nil, true
		}
		log.Printf("an error occurs when getting node, err: %s" + err.Error())
		return nil, false
	}

	return node, node.State == state
}

// CreateNode creates asynchronously a cloud node. The node returned will depend of the administrative state passed in options.
// If none is passed models.AdministrativeStateDeployed will be set by default.
// The valid administrative states options for a node creation are [models.AdministrativeStateCreationPending, models.AdministrativeStateCreationProceed, models.AdministrativeStateCreationError, models.AdministrativeStateDeployed]
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

	if _, ok := validCreationAdministrativeStates[cloudOptions.administrativeState]; !ok {
		return nil, ErrCreationAdministrativeState
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

	success, nodePoled := WaitForAdministrativeState(ctx, c, workspaceID, node.Data.ID.String(), cloudOptions.administrativeState, checkNodeAdministrativeState)
	if !success {
		return nil, fmt.Errorf("Node did not reach '%s' state in time.", cloudOptions.administrativeState)
	}

	return nodePoled, nil
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

// DeleteNode deletes asynchronously a cloud node. The node returned will depend of the administrative state passed in options.
// If none is passed models.AdministrativeStateDeleted will be set by default.
// The valid administrative states options for a node creation are [ models.AdministrativeStateDeletePending, models.AdministrativeStateDeleteProceed, models.AdministrativeStateDeleteError, models.AdministrativeStateDeleted]
func (c *Client) DeleteNode(ctx context.Context, workspaceID, nodeID string, options ...OptionElement) (*models.Node, error) {
	cloudOptions := &elementOptions{}
	for _, o := range options {
		o(cloudOptions)
	}

	if cloudOptions.administrativeState == "" {
		cloudOptions.administrativeState = models.AdministrativeStateDeleted
	}

	if _, ok := validDeletionAdministrativeStates[cloudOptions.administrativeState]; !ok {
		return nil, ErrDeletionAdministrativeState
	}

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

	success, nodePoled := WaitForAdministrativeState(ctx, c, workspaceID, node.Data.ID.String(), cloudOptions.administrativeState, checkNodeAdministrativeState)
	if !success {
		return nil, fmt.Errorf("Node did not reach '%s' state in time.", cloudOptions.administrativeState)
	}

	return nodePoled, nil
}
