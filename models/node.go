package models

import (
	"time"
)

type NodeType string

const (
	NodeTypeAccess = "access"
	NodeTypeCloud  = "cloud"
	NodeTypeBridge = "bridge"
	NodeTypeRouter = "router"
)

func (nt NodeType) String() string {
	return string(nt)
}

type NodeProduct struct {
	Product
	CSPName         string `json:"cspName,omitempty"`
	CSPNameUnderlay string `json:"cspNameUnderlay,omitempty"`
	CSPCity         string `json:"cspCity,omitempty"`
	CSPRegion       string `json:"cspRegion,omitempty"`
}

type Port struct {
	ID              string `json:"id"`
	LocationID      string `json:"locationId"`
	CSPName         string `json:"cspName"`
	CSPNameUnderlay string `json:"cspNameUnderlay"`
}

type ProviderCloudConfig struct {
	PairingKey string `json:"pairingKey,omitempty"`
	AccountID  string `json:"accountId,omitempty"`
	ServiceKey string `json:"serviceKey,omitempty"`
}

type Node struct {
	BaseModel
	AccountID      string               `json:"accountId"`
	WorkspaceID    string               `json:"workspaceId"`
	Name           string               `json:"name"`
	State          AdministrativeState  `json:"administrativeState"`
	DeployedAt     *time.Time           `json:"deployedAt,omitempty"`
	Product        NodeProduct          `json:"product,omitempty"`
	Type           NodeType             `json:"type,omitempty"`
	ConnectionID   string               `json:"connectionId,omitempty"`
	Port           *Port                `json:"port,omitempty"`
	ProviderConfig *ProviderCloudConfig `json:"providerConfig,omitempty"`
	Vlan           int64                `json:"vlan,omitempty"`
	DxconID        string               `json:"dxconId,omitempty"`
}

type NodeResponse struct {
	Data Node `json:"data"`
}

type AddProduct struct {
	SKU string `json:"sku" binding:"required"`
}
type CreateNode struct {
	Name           string               `json:"name" binding:"required"`
	Type           string               `json:"type" binding:"required"`
	Product        AddProduct           `json:"product" binding:"required"`
	ProviderConfig *ProviderCloudConfig `json:"providerConfig" binding:"required_if=Type cloud"`
}

type UpdateNode struct {
	Name string `json:"name"`
}
