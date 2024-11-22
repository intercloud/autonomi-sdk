package autonomisdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/intercloud/autonomi-sdk/models"
)

func (c *Client) CreateUser(ctx context.Context, payload models.CreateUser) (*models.User, error) {
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(&payload)
	if err != nil {
		return nil, err
	}

	if errV := c.validate.StructCtx(ctx, payload); errV != nil {
		return nil, errV
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/accounts/%s/users", c.hostURL, c.accountID), body)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	user := models.UserResponse{}
	if err := json.Unmarshal(resp, &user); err != nil {
		return nil, err
	}

	return &user.Data, nil
}

func (c *Client) ListUsers(ctx context.Context, accountID uuid.UUID) (models.Users, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/accounts/%s/users", c.hostURL, accountID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	users := models.Users{}
	if err = json.Unmarshal(resp, &users); err != nil {
		return nil, err
	}

	return users, nil
}
