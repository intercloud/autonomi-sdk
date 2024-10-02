package models

import (
	"time"

	"github.com/google/uuid"
)

type NodeType string

const (
	NodeTypeAccess NodeType = "access"
	NodeTypeCloud  NodeType = "cloud"
	NodeTypeBridge NodeType = "bridge"
	NodeTypeRouter NodeType = "router"
)

func (nt NodeType) String() string {
	return string(nt)
}

type AccessProductType string

const (
	AccessProductTypePhysical AccessProductType = "PHYSICAL"
	AccessProductTypeVirtual  AccessProductType = "VIRTUAL"
)

func (at AccessProductType) String() string {
	return string(at)
}

type NodeProduct struct {
	Product
	CSPName         string            `json:"cspName,omitempty"`
	CSPNameUnderlay string            `json:"cspNameUnderlay,omitempty"`
	CSPCity         string            `json:"cspCity,omitempty"`
	CSPRegion       string            `json:"cspRegion,omitempty"`
	Type            AccessProductType `json:"type,omitempty"`
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

type ServiceKey struct {
	ID             string    `json:"id,omitempty"`
	ExpirationDate time.Time `json:"expirationDate,omitempty"`
	Name           string    `json:"name,omitempty"`
}

type Node struct {
	BaseModel
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
	Error          *SupportError        `json:"error,omitempty"`
	PhysicalPort   *PhysicalPort        `json:"physicalPort,omitempty"`
	ServiceKey     *ServiceKey          `json:"serviceKey,omitempty"`
}

func (n *Node) GetState() AdministrativeState {
	return n.State
}

type NodeResponse struct {
	Data Node `json:"data"`
}

type AddProduct struct {
	SKU string `json:"sku" binding:"required"`
}
type CreateNode struct {
	Name           string               `json:"name" binding:"required"`
	Type           NodeType             `json:"type" binding:"required"`
	Product        AddProduct           `json:"product" binding:"required"`
	ProviderConfig *ProviderCloudConfig `json:"providerConfig" binding:"required_if=Type cloud"`
	PhysicalPortID *uuid.UUID           `json:"physicalPortId,omitempty"`
	Vlan           int64                `json:"vlan,omitempty"`
}
