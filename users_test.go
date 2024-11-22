package autonomisdk

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/intercloud/autonomi-sdk/models"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var (
	userId = uuid.MustParse("14bb42ef-571d-403e-a62f-d4fb9717ac6c")

	user = models.User{
		BaseModel: models.BaseModel{
			ID: userId,
		},
		Name:      "name",
		Email:     "email@gmail.com",
		Activated: false,
		AccountID: accountID,
		IsAdmin:   true,
	}

	userCreateResponse = models.UserResponse{
		Data: models.User{
			BaseModel: models.BaseModel{
				ID: userId,
			},
			Name:      "name",
			Email:     "email@gmail.com",
			Activated: false,
			AccountID: accountID,
			IsAdmin:   false,
		},
	}

	usersListResponse = models.Users{
		user,
	}
)

func TestListUsers(t *testing.T) {
	// init testing framework
	g := NewWithT(t)
	gh := ghttp.NewGHTTPWithGomega(g)

	// init server
	server := ghttp.NewServer()
	defer server.Close()

	serverURL, err := url.Parse(server.URL())
	g.Expect(err).ShouldNot(HaveOccurred())

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodGet, "/users/self"),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusOK, models.Self{
				AccountID: uuid.MustParse(accountId),
			}),
		),
	)
	// init testing http client
	cli, err := NewClient(
		true,
		WithHostURL(serverURL),
		WithHTTPClient(&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, //nolint:gosec //No
				},
			},
		}),
		WithPersonalAccessToken(personalAccessToken),
	)
	g.Expect(err).ShouldNot(HaveOccurred())

	// mock testing response
	result := usersListResponse
	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/users", accountId)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusOK, usersListResponse),
		),
	)

	// run target function through testing framework
	data, err := cli.ListUsers(
		context.Background(),
		uuid.MustParse(accountId),
	)

	// test results
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(data).Should(Equal(result))
}

func TestListForbidden(t *testing.T) {
	// init testing framework
	g := NewWithT(t)
	gh := ghttp.NewGHTTPWithGomega(g)

	// init server
	server := ghttp.NewServer()
	defer server.Close()

	serverURL, err := url.Parse(server.URL())
	g.Expect(err).ShouldNot(HaveOccurred())

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodGet, "/users/self"),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusOK, models.Self{
				AccountID: uuid.MustParse(accountId),
			}),
		),
	)
	// init testing http client
	cli, err := NewClient(
		true,
		WithHostURL(serverURL),
		WithHTTPClient(&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, //nolint:gosec //No
				},
			},
		}),
		WithPersonalAccessToken(personalAccessToken),
	)
	g.Expect(err).ShouldNot(HaveOccurred())

	// mock testing response
	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/users", accountId)),
			gh.RespondWithJSONEncoded(http.StatusForbidden, nil),
		),
	)

	// run target function through testing framework
	data, err := cli.ListUsers(
		context.Background(),
		uuid.MustParse(accountId),
	)

	// test results
	g.Expect(err).ShouldNot(BeNil())
	g.Expect(data).Should(BeNil())
}

func TestCreateUserSuccessfully(t *testing.T) {
	g := NewWithT(t)
	gh := ghttp.NewGHTTPWithGomega(g)

	server := ghttp.NewServer()
	defer server.Close()

	serverURL, err := url.Parse(server.URL())
	g.Expect(err).ShouldNot(HaveOccurred())

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodGet, "/users/self"),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusOK, models.Self{
				AccountID: uuid.MustParse(accountId),
			}),
		),
	)

	cli, err := NewClient(
		true,
		WithHostURL(serverURL),
		WithHTTPClient(&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, //nolint:gosec //No
				},
			},
		}),
		WithPersonalAccessToken(personalAccessToken),
	)
	g.Expect(err).ShouldNot(HaveOccurred())

	result := userCreateResponse

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodPost, fmt.Sprintf("/accounts/%s/users", accountId)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusOK, userCreateResponse),
		),
	)

	data, err := cli.CreateUser(
		context.Background(),
		models.CreateUser{
			Name:  "name",
			Email: "email@gmail.com",
		},
	)

	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(*data).Should(Equal(result.Data))
}
