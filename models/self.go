package models

import "github.com/google/uuid"

type Self struct {
	AccountID uuid.UUID `json:"accountId"`
}
