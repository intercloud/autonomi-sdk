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
	State              string              `json:"administrativeState"`
	UsedVLANs          []int64             `json:"usedVlans"`
}

type PhysicalPortResponse struct {
	Data []PhysicalPort `json:"data"`
}
