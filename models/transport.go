package models

import "time"

type TransportProduct struct {
	Product
	LocationTo string `json:"locationTo"`
}

type TransportVlans struct {
	AVlan int64 `json:"aVlan,omitempty"`
	ZVlan int64 `json:"zVlan,omitempty"`
}

type Transport struct {
	BaseModel
	WorkspaceID    string              `json:"workspaceId"`
	Name           string              `json:"name"`
	State          AdministrativeState `json:"administrativeState"`
	DeployedAt     *time.Time          `json:"deployedAt,omitempty"`
	Error          *SupportError       `json:"error,omitempty"`
	TransportVlans TransportVlans      `json:"vlans,omitempty"`
	IsLocal        bool                `json:"isLocal"`
	Product        TransportProduct    `json:"product,omitempty"`
	ConnectionID   string              `json:"connectionId,omitempty"`
}

func (t *Transport) GetState() AdministrativeState {
	return t.State
}

type CreateTransport struct {
	Name    string     `json:"name" binding:"required"`
	Product AddProduct `json:"product" binding:"required"`
}

type TransportResponse struct {
	Data Transport `json:"data"`
}
