package autonomisdk

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/intercloud/autonomi-sdk/models"
)

func (c *Client) GetSelf() (uuid.UUID, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/users/self", c.hostURL), nil)
	if err != nil {
		return uuid.Nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return uuid.Nil, err
	}

	self := models.Self{}
	err = json.Unmarshal(resp, &self)
	if err != nil {
		return uuid.Nil, err
	}

	return self.AccountID, nil
}
