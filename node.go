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

func checkNodeFinishedTask(ctx context.Context, c *Client, workspaceID, nodeID string, waiterOptionState models.AdministrativeState) (*models.Node, bool) {
	node, err := c.GetNode(ctx, workspaceID, nodeID)
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

	return node, waiterOptionState == node.State
}

// CreateNode creates asynchronously a cloud node. The node returned will depend of the passed option.
// If none is passed the node will be returned once created in database with administrative state creation_pending.
// If the option WithWaitUntilElementDeployed() is passed, the node will be returned when its state reach deployed or creation_error.
// If the option WithWaitUntilElementUndeployed() is passed, it will not be considered hence the node returned will be in state creation_pending
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

	var nodePolled = &node.Data
	if cloudOptions.waitUntilElementDeployed {
		var success bool
		nodePolled, success = WaitUntilFinishedTask(ctx, c, workspaceID, node.Data.ID.String(), models.AdministrativeStateDeployed, checkNodeFinishedTask)
		if !success {
			return nil, fmt.Errorf("Node did not reach '%s' state in time.", models.AdministrativeStateDeployed)
		}
	}

	return nodePolled, nil
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

// DeleteNode deletes asynchronously a node. The attachment returned will depend of the option passed.
// If none is passed the node will be returned once the request accepted, its state will be delete_pending
// If the option WithWaitUntilElementUndeployed() is passed, the node won't be returned as it would have been deleted. However, if an error is triggered, an object could be returned with a delete_error state.
// If the option WithWaitUntilElementDeployed() is passed, it will not be considered hence the node returned will be in state delete_pending
func (c *Client) DeleteNode(ctx context.Context, workspaceID, nodeID string, options ...OptionElement) (*models.Node, error) {
	cloudOptions := &elementOptions{}
	for _, o := range options {
		o(cloudOptions)
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

	var nodePolled = &node.Data
	if cloudOptions.waitUntilElementUndeployed {
		var success bool
		nodePolled, success = WaitUntilFinishedTask(ctx, c, workspaceID, node.Data.ID.String(), models.AdministrativeStateDeleted, checkNodeFinishedTask)
		if !success {
			return nil, fmt.Errorf("Node did not reach '%s' state in time.", models.AdministrativeStateDeleted)
		}
	}

	// If the node was deleted and we were waiting for the "deleted" state,
	// nodePolled will be nil. To prevent a panic when dereferencing, we
	// assign an empty structure.
	if nodePolled == nil {
		nodePolled = &models.Node{}
	}

	return nodePolled, nil
}
