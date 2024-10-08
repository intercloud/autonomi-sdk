package models

type PhysicalPortProduct struct {
	Product
}

type PhysicalPort struct {
	BaseModel
	Name               string              `json:"name"`
	AccountID          string              `json:"accountId"`
	Product            PhysicalPortProduct `json:"product"`
	AvailableBandwidth int                 `json:"availableBandwidth"`
	State              AdministrativeState `json:"administrativeState"`
	UsedVLANs          []int64             `json:"usedVlans"`
	LOAAccessURL       string              `json:"loaAccessUrl"`
}

type PhysicalPortSingleResponse struct {
	Data PhysicalPort `json:"data"`
}

type PhysicalPortListResponse struct {
	Data []PhysicalPort `json:"data"`
}

type CreatePhysicalPort struct {
	Name    string     `json:"name" binding:"required"`
	Product AddProduct `json:"product" binding:"required"`
}
