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

func checkAttachmentAdministrativeState(ctx context.Context, c *Client, workspaceID, attachmentID string, state models.AdministrativeState) (*models.Attachment, bool) {
	attachment, err := c.GetAttachment(ctx, workspaceID, attachmentID)
	if err != nil {
		// if wanted state is deleted and the attachment is in this state, api has returned 404
		if state == models.AdministrativeStateDeleted && strings.Contains(err.Error(), "status: 404") {
			return nil, true
		}
		log.Printf("an error occurs when getting attachment, err: %s" + err.Error())
		return nil, false
	}

	return attachment, attachment.State == state
}

// CreateAttachment creates asynchronously an attachment. The attachment returned will depend of the administrative state passed in options.
// If none is passed models.AdministrativeStateDeployed will be set by default.
// The valid administrative states options for a attachment creation are [models.AdministrativeStateCreationPending, models.AdministrativeStateCreationProceed, models.AdministrativeStateCreationError, models.AdministrativeStateDeployed]
func (c *Client) CreateAttachment(ctx context.Context, payload models.CreateAttachment, workspaceID string, options ...OptionElement) (*models.Attachment, error) {
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(&payload)
	if err != nil {
		return nil, err
	}

	if errV := c.validate.StructCtx(ctx, payload); errV != nil {
		return nil, errV
	}

	attachmentOptions := &elementOptions{}
	for _, o := range options {
		o(attachmentOptions)
	}

	if attachmentOptions.administrativeState == "" {
		attachmentOptions.administrativeState = models.AdministrativeStateDeployed
	}

	if _, ok := validCreationAdministrativeStates[attachmentOptions.administrativeState]; !ok {
		return nil, ErrCreationAdministrativeState
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/accounts/%s/workspaces/%s/attachments", c.hostURL, c.accountID, workspaceID), body)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	attachment := models.AttachmentResponse{}
	err = json.Unmarshal(resp, &attachment)
	if err != nil {
		return nil, err
	}

	success, attachmentPolled := WaitForAdministrativeState(ctx, c, workspaceID, attachment.Data.ID.String(), attachmentOptions.administrativeState, checkAttachmentAdministrativeState)
	if !success {
		return nil, fmt.Errorf("Attachment did not reach '%s' state in time.", attachmentOptions.administrativeState)
	}

	return attachmentPolled, nil
}

func (c *Client) GetAttachment(ctx context.Context, workspaceID, attachmentID string) (*models.Attachment, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/accounts/%s/workspaces/%s/attachments/%s", c.hostURL, c.accountID, workspaceID, attachmentID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	attachment := models.AttachmentResponse{}
	err = json.Unmarshal(resp, &attachment)
	if err != nil {
		return nil, err
	}

	return &attachment.Data, err
}

// DeleteAttachment deletes asynchronously a cloud attachment. The attachment returned will depend of the administrative state passed in options.
// If none is passed models.AdministrativeStateDeleted will be set by default.
// The valid administrative states options for a attachment creation are [ models.AdministrativeStateDeletePending, models.AdministrativeStateDeleteProceed, models.AdministrativeStateDeleteError, models.AdministrativeStateDeleted]
func (c *Client) DeleteAttachment(ctx context.Context, workspaceID, attachmentID string, options ...OptionElement) (*models.Attachment, error) {
	attachmentOptions := &elementOptions{}
	for _, o := range options {
		o(attachmentOptions)
	}

	if attachmentOptions.administrativeState == "" {
		attachmentOptions.administrativeState = models.AdministrativeStateDeleted
	}

	if _, ok := validDeletionAdministrativeStates[attachmentOptions.administrativeState]; !ok {
		return nil, ErrDeletionAdministrativeState
	}

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/accounts/%s/workspaces/%s/attachments/%s", c.hostURL, c.accountID, workspaceID, attachmentID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	attachment := models.AttachmentResponse{}
	err = json.Unmarshal(resp, &attachment)
	if err != nil {
		return nil, err
	}

	success, attachmentPolled := WaitForAdministrativeState(ctx, c, workspaceID, attachment.Data.ID.String(), attachmentOptions.administrativeState, checkAttachmentAdministrativeState)
	if !success {
		return nil, fmt.Errorf("Attachment did not reach '%s' state in time.", attachmentOptions.administrativeState)
	}

	return attachmentPolled, nil
}
