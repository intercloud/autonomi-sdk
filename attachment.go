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

func checkAttachmentFinishedTask(ctx context.Context, c *Client, workspaceID, attachmentID string, waiterOptionState models.AdministrativeState) (*models.Attachment, bool) {
	attachment, err := c.GetAttachment(ctx, workspaceID, attachmentID)
	if err != nil {
		// if wanted state is deleted and the attachment is in this state, api has returned 404
		if waiterOptionState == models.AdministrativeStateDeleted {
			if strings.Contains(err.Error(), "status: 404") {
				return nil, true
			}
		}
		log.Printf("an error occurs when getting attachment, err: %s" + err.Error())
		return nil, false
	}

	return attachment, waiterOptionState == attachment.State
}

// CreateAttachment creates asynchronously an attachment. The attachment returned will depend of the passed option.
// If none is passed the attachment will be returned once created in database with administrative state creation_pending.
// If the option WithWaitUntilElementDeployed() is passed, the attachment will be returned when its state reach deployed or creation_error.
// If the option WithWaitUntilElementUndeployed() is passed, it will not be considered hence the attachment returned will be in state creation_pending
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

	var attachmentPolled = &attachment.Data
	if attachmentOptions.waitUntilElementDeployed {
		var success bool
		attachmentPolled, success = WaitUntilFinishedTask(ctx, c, workspaceID, attachment.Data.ID.String(), models.AdministrativeStateDeployed, checkAttachmentFinishedTask)
		if !success {
			return nil, fmt.Errorf("Attachment did not reach '%s' state in time.", models.AdministrativeStateDeployed)
		}
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

// DeleteAttachment deletes asynchronously an attachment. The attachment returned will depend of the option passed.
// If none is passed the attachment will be returned once the request accepted, its state will be delete_pending
// If the option WithWaitUntilElementUndeployed() is passed, the attachment won't be returned as it would have been deleted. However, if an error is triggered, an object could be returned with a delete_error state.
// If the option WithWaitUntilElementDeployed() is passed, it will not be considered hence the attachment returned will be in state delete_pending
func (c *Client) DeleteAttachment(ctx context.Context, workspaceID, attachmentID string, options ...OptionElement) (*models.Attachment, error) {
	attachmentOptions := &elementOptions{}
	for _, o := range options {
		o(attachmentOptions)
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

	var attachmentPolled = &attachment.Data
	if attachmentOptions.waitUntilElementUndeployed {
		var success bool
		attachmentPolled, success = WaitUntilFinishedTask(ctx, c, workspaceID, attachment.Data.ID.String(), models.AdministrativeStateDeleted, checkAttachmentFinishedTask)
		if !success {
			return nil, fmt.Errorf("Attachment did not reach '%s' state in time.", models.AdministrativeStateDeleted)
		}
	}

	// If the attachment was deleted and we were waiting for the "deleted" state,
	// attachmentPolled will be nil. To prevent a panic when dereferencing, we
	// assign an empty structure.
	if attachmentPolled == nil {
		attachmentPolled = &models.Attachment{}
	}

	return attachmentPolled, nil
}
