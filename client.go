package autonomisdk

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type pollElement struct {
	retryInterval time.Duration
	maxRetry      int
}

type Client struct {
	hostURL             *url.URL
	httpClient          *http.Client
	personalAccessToken string
	accountID           uuid.UUID

	validate *validator.Validate

	poll pollElement
}

type OptionClient func(*Client)

const timeout = 30 * time.Second

var (
	ErrTermsAndConditionsRequired  = errors.New("terms and conditions must be accepted")
	ErrHostURLRequired             = errors.New("host url must be set, please use the option WithHostURL()")
	ErrPersonalAccessTokenRequired = errors.New("personal acess token must be set, please use the option WithPersonalAccessToken()")
)

func WithHostURL(url *url.URL) OptionClient {
	return func(a *Client) {
		a.hostURL = url
	}
}

// WithHTTPClient init http client. If its timeout is above 30s it will be overide to be equal to 30 seconds
func WithHTTPClient(client *http.Client) OptionClient {
	return func(a *Client) {
		a.httpClient = client

		if a.httpClient.Timeout > timeout {
			a.httpClient.Timeout = timeout
		}
	}
}

func WithPersonalAccessToken(token string) OptionClient {
	return func(a *Client) {
		a.personalAccessToken = token
	}
}

func initClient(opts ...OptionClient) *Client {
	client := &Client{
		httpClient: &http.Client{},
		hostURL:    &url.URL{},
		poll: pollElement{
			retryInterval: 20 * time.Second,
			maxRetry:      30,
		},
	}

	for _, o := range opts {
		o(client)
	}

	return client
}

// NewClient - Init and return an http client
func NewClient(termsAndConditions bool, opts ...OptionClient) (*Client, error) {
	if !termsAndConditions {
		return nil, ErrTermsAndConditionsRequired
	}

	client := initClient(opts...)

	if client.hostURL != nil && client.hostURL.Host == "" {
		return nil, ErrHostURLRequired
	}

	if client.personalAccessToken == "" {
		return nil, ErrPersonalAccessTokenRequired
	}

	accountID, err := client.GetSelf()
	if err != nil {
		return nil, err
	}

	client.accountID = accountID

	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return nil, errors.New("not a valid binding validator")
	}
	client.validate = validate

	return client, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.personalAccessToken))
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
