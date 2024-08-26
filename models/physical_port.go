package models

type PhysicalPort struct {
	BaseModel
}

type PhysicalPortResponse struct {
	Data []PhysicalPort `json:"data"`
}
