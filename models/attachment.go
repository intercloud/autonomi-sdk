package models

import (
	"time"
)

type Attachment struct {
	BaseModel
	TransportID string              `json:"transportId"`
	NodeID      string              `json:"nodeId"`
	State       AdministrativeState `json:"administrativeState"`
	Side        string              `json:"side"`
	DeployedAt  *time.Time          `json:"deployedAt,omitempty"`
	Error       *SupportError       `json:"error,omitempty"`
	WorkspaceID string              `json:"workspaceId"`
}

type AttachmentResponse struct {
	Data Attachment `json:"data"`
}

type CreateAttachment struct {
	NodeID      string `json:"nodeId" binding:"required"`
	TransportID string `json:"transportId" binding:"required"`
}

func (a *Attachment) GetState() AdministrativeState {
	return a.State
}
