package autonomisdk

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

var timeout = 30 * time.Second

var (
	ErrTermsAndConditionsRequired  = errors.New("terms and conditions must be accepted")
	ErrHostURLRequired             = errors.New("host url must be set, please use the option WithHostURL()")
	ErrPersonalAccessTokenRequired = errors.New("personal acess token must be set, please use the option WithPersonalAccessToken()")
)

type Client struct {
	hostURL                  *url.URL
	httpClient               *http.Client
	personalAccessTokentoken string
	accountID                uuid.UUID
}

type OptionClient func(*Client)

func WithHostURL(url *url.URL) OptionClient {
	return func(a *Client) {
		a.hostURL = url
	}
}

func WithHTTPClient(client *http.Client) OptionClient {
	return func(a *Client) {
		a.httpClient = client
	}
}

func WithPersonalAccessToken(token string) OptionClient {
	return func(a *Client) {
		a.personalAccessTokentoken = token
	}
}

// NewClient - Init and return an http client
func NewClient(termsAndConditions bool, opts ...OptionClient) (*Client, error) {
	if !termsAndConditions {
		return nil, ErrTermsAndConditionsRequired
	}

	// by default requests timeout after 30s
	client := &Client{
		httpClient: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, //nolint:gosec //No
				},
			},
		},
		hostURL: &url.URL{},
	}

	for _, o := range opts {
		o(client)
	}

	if client.hostURL == nil {
		return nil, ErrHostURLRequired
	}

	if client.personalAccessTokentoken == "" {
		return nil, ErrPersonalAccessTokenRequired
	}

	accountID, err := client.GetSelf()
	if err != nil {
		return nil, err
	}

	client.accountID = accountID

	return client, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.personalAccessTokentoken))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
