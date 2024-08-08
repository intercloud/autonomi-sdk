package models

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"-"`
}

type AdministrativeState string

const (
	AdministrativeStateCreationPending AdministrativeState = "creation_pending"
	AdministrativeStateCreationProceed AdministrativeState = "creation_proceed"
	AdministrativeStateCreationError   AdministrativeState = "creation_error"
	AdministrativeStateDeployed        AdministrativeState = "deployed"
	AdministrativeStateDeletePending   AdministrativeState = "delete_pending"
	AdministrativeStateDeleteProceed   AdministrativeState = "delete_proceed"
	AdministrativeStateDeleteError     AdministrativeState = "delete_error"
	AdministrativeStateDeleted         AdministrativeState = "deleted"
)

func (as AdministrativeState) String() string {
	return string(as)
}
