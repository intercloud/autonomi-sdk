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
	portListResponse = models.PhysicalPortListResponse{
		Data: []models.PhysicalPort{
			{
				BaseModel: models.BaseModel{
					ID: uuid.New(),
				},
			},
			{
				BaseModel: models.BaseModel{
					ID: uuid.New(),
				},
			},
		},
	}
	portCreatedListResponse = models.PhysicalPortListResponse{
		Data: []models.PhysicalPort{
			{
				BaseModel: models.BaseModel{
					ID: uuid.New(),
				},
			},
		},
	}
	portCreatedSingleResponse = models.PhysicalPortSingleResponse{
		Data: models.PhysicalPort{
			BaseModel: models.BaseModel{
				ID: uuid.New(),
			},
		},
	}
)

func TestCreatePhysicalPortSuccessfully(t *testing.T) {
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

	result := portCreatedSingleResponse

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodPost, fmt.Sprintf("/accounts/%s/ports", accountId)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusAccepted, portCreatedSingleResponse),
		),
	)

	data, err := cli.CreatePhysicalPort(
		context.Background(),
		models.CreatePhysicalPort{
			Name: "physical_port_name",
			Product: models.AddProduct{
				SKU: "physical_port_sku",
			},
		},
	)

	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(*data).Should(Equal(result.Data))
}

func TestCreatePhysicalPortForbidden(t *testing.T) {
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

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodPost, fmt.Sprintf("/accounts/%s/ports", accountId)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusForbidden, nil),
		),
	)

	data, err := cli.CreatePhysicalPort(
		context.Background(),
		models.CreatePhysicalPort{
			Name: "physical_port_name",
			Product: models.AddProduct{
				SKU: "physical_port_sku",
			},
		},
	)

	g.Expect(err).ShouldNot(BeNil())
	g.Expect(data).Should(BeNil())
}

func TestCreatePhysicalPortFailedValidator(t *testing.T) {
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

	data, err := cli.CreatePhysicalPort(
		context.Background(),
		models.CreatePhysicalPort{
			Product: models.AddProduct{
				SKU: "physical_port_sku",
			},
		},
	)

	g.Expect(err.Error()).Should(Equal("Key: 'CreatePhysicalPort.Name' Error:Field validation for 'Name' failed on the 'required' tag"))
	g.Expect(data).Should(BeNil())

	data, err = cli.CreatePhysicalPort(
		context.Background(),
		models.CreatePhysicalPort{
			Name: "physical_port_name",
		},
	)

	g.Expect(err.Error()).Should(Equal("Key: 'CreatePhysicalPort.Product.SKU' Error:Field validation for 'SKU' failed on the 'required' tag"))
	g.Expect(data).Should(BeNil())
}

func TestGetPhysicalPortSuccessfully(t *testing.T) {
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

	result := portCreatedSingleResponse

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/ports/%s", accountId, physicalPortId)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusOK, portCreatedSingleResponse),
		),
	)

	data, err := cli.GetPhysicalPort(
		context.Background(),
		physicalPortId.String(),
	)

	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(*data).Should(Equal(result.Data))
}

func TestGetPhysicalPortNotFound(t *testing.T) {
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

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/ports/%s", accountId, physicalPortId)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusNotFound, nil),
		),
	)

	data, err := cli.GetPhysicalPort(
		context.Background(),
		physicalPortId.String(),
	)

	g.Expect(err).ShouldNot(BeNil())
	g.Expect(data).Should(BeNil())
}

func TestListPhysicalPort(t *testing.T) {
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
	result := portListResponse
	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/ports", accountId)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusOK, portListResponse),
		),
	)

	// run target function through testing framework
	data, err := cli.ListPort()

	// test results
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(*data).Should(Equal(result.Data))
}

func TestListPhysicalPortWithState(t *testing.T) {
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
	result := portCreatedListResponse
	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/ports", accountId), "state=created"),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusOK, portCreatedListResponse),
		),
	)

	// run target function through testing framework
	data, err := cli.ListPort(WithAdministrativeState(models.AdministrativeStateCreated))

	// test results
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(*data).Should(Equal(result.Data))
}

func TestDeletePhysicalPortSuccessfully(t *testing.T) {
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

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodDelete, fmt.Sprintf("/accounts/%s/ports/%s", accountId, physicalPortId)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusNoContent, nil),
		),
	)

	err = cli.DeletePhysicalPort(
		context.Background(),
		physicalPortId.String(),
	)

	g.Expect(err).ShouldNot(HaveOccurred())
}

func TestDeletePhysicalPortForbidden(t *testing.T) {
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

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodDelete, fmt.Sprintf("/accounts/%s/ports/%s", accountId, physicalPortId)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusForbidden, nil),
		),
	)

	err = cli.DeletePhysicalPort(
		context.Background(),
		physicalPortId.String(),
	)

	g.Expect(err).ShouldNot(BeNil())
}
