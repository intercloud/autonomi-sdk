package autonomisdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/intercloud/autonomi-sdk/models"
)

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
