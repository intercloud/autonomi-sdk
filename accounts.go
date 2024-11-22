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

func (c *Client) CreateAccount(ctx context.Context, payload models.Account) (*models.Account, error) {
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(&payload)
	if err != nil {
		return nil, err
	}

	if errV := c.validate.StructCtx(ctx, payload); errV != nil {
		return nil, errV
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/accounts", c.hostURL), body)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	account := models.Account{}
	if err = json.Unmarshal(resp, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

func (c *Client) ListAccounts(ctx context.Context) (models.Accounts, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/accounts", c.hostURL), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	accounts := models.Accounts{}
	if err = json.Unmarshal(resp, &accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (c *Client) DeleteAccount(ctx context.Context, accountID uuid.UUID) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/accounts/%s", c.hostURL, accountID), nil)
	if err != nil {
		return err
	}

	if _, err = c.doRequest(req); err != nil {
		return err
	}

	return nil
}
