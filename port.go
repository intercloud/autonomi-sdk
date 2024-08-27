package autonomisdk

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/intercloud/autonomi-sdk/models"
)

func (c *Client) ListPort(options ...OptionElement) (*[]models.PhysicalPort, error) {

	// retrieve options from request
	portOptions := &elementOptions{}
	for _, o := range options {
		o(portOptions)
	}

	// run request
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/accounts/%s/ports", c.hostURL, c.accountID), nil)
	if err != nil {
		return nil, err
	}

	// add query param if needed
	if portOptions.administrativeState != "" {
		q := req.URL.Query()
		q.Add("state", portOptions.administrativeState.String())
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	ports := models.PhysicalPortResponse{}
	err = json.Unmarshal(resp, &ports)
	if err != nil {
		return nil, err
	}

	return &ports.Data, err
}
