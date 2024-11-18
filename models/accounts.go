package models

type Account struct {
	BaseModel
	Name             string `json:"name" binding:"required"`
	Address          string `json:"address" binding:"required"`
	ZipCode          string `json:"zipCode" binding:"required"`
	City             string `json:"city" binding:"required"`
	Country          string `json:"country" binding:"required"`
	FinancialContact string `json:"financialContact,omitempty"`
	TechnicalContact string `json:"technicalContact,omitempty"`
}

type Accounts []Account
