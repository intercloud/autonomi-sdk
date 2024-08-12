package autonomisdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/intercloud/autonomi-sdk/models"
)

func (c *Client) CreateTransport(ctx context.Context, payload models.CreateTransport, workspaceID string) (*models.Transport, error) {
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(&payload)
	if err != nil {
		return nil, err
	}

	if errV := c.validate.StructCtx(ctx, payload); errV != nil {
		return nil, errV
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/accounts/%s/workspaces/%s/transports", c.hostURL, c.accountID, workspaceID), body)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	transport := models.TransportResponse{}
	err = json.Unmarshal(resp, &transport)
	if err != nil {
		return nil, err
	}

	return &transport.Data, err
}