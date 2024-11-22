package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	BaseModel
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	Activated       bool       `json:"activated"`
	AccountID       uuid.UUID  `json:"accountId"`
	CGUAcceptedDate *time.Time `json:"cguAcceptedDate,omitempty"`
	LastConnection  *time.Time `json:"lastConnection,omitempty"`
	IsAdmin         bool       `json:"isAdmin"`
}

type Users []User

type CreateUser struct {
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required"`
	ProfileId string `json:"profileId"`
}

type UserResponse struct {
	Data User `json:"data"`
}
