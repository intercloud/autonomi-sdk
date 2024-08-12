package autonomisdk

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/uuid"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/intercloud/autonomi-sdk/models"
)

var (
	transportID = uuid.MustParse("20ea29c9-f892-4aaf-8907-de79fa83e7bb")

	transportCreateResponse = models.TransportResponse{
		Data: models.Transport{
			BaseModel: models.BaseModel{
				ID: transportID,
			},
			WorkspaceID: workspaceID,
			Name:        "transport_name",
			State:       models.AdministrativeStateCreationPending,
			Product: models.TransportProduct{
				Product: models.Product{
					Provider:  "EQUINIX",
					Duration:  0,
					Location:  "EQUINIX FR5",
					Bandwidth: 100,
					PriceNRC:  0,
					PriceMRC:  0,
					CostNRC:   0,
					CostMRC:   0,
					SKU:       "CEQUFR5100AWS",
				},
				LocationTo: "EQUINIX LD5",
			},
		},
	}

	transportCreateErrorResponse = models.TransportResponse{
		Data: models.Transport{
			BaseModel: models.BaseModel{
				ID: transportID,
			},
			WorkspaceID: workspaceID,
			Name:        "transport_name_error",
			State:       models.AdministrativeStateCreationError,
			Product: models.TransportProduct{
				Product: models.Product{
					Provider:  "EQUINIX",
					Duration:  0,
					Location:  "EQUINIX FR5",
					Bandwidth: 100,
					PriceNRC:  0,
					PriceMRC:  0,
					CostNRC:   0,
					CostMRC:   0,
					SKU:       "CEQUFR5100AWS",
				},
				LocationTo: "EQUINIX LD5",
			},
			Error: &models.SupportError{
				Code: "ERR_INTERNAL",
				Msg:  "an internal error occured",
			},
		},
	}

	transportUpdateResponse = models.TransportResponse{
		Data: models.Transport{
			BaseModel: models.BaseModel{
				ID: transportID,
			},
			WorkspaceID: workspaceID,
			Name:        "transport_updated_name",
			State:       models.AdministrativeStateDeployed,
			TransportVlans: models.TransportVlans{
				AVlan: 19,
				ZVlan: 19,
			},
			Product: models.TransportProduct{
				Product: models.Product{
					Provider:  "EQUINIX",
					Duration:  0,
					Location:  "EQUINIX FR5",
					Bandwidth: 100,
					PriceNRC:  0,
					PriceMRC:  0,
					CostNRC:   0,
					CostMRC:   0,
					SKU:       "CEQUFR5100AWS",
				},
				LocationTo: "EQUINIX LD5",
			},
			ConnectionID: "3091af46-3586-4cd1-bdbf-b569d2219823",
		},
	}

	transportDeleteResponse = models.TransportResponse{
		Data: models.Transport{
			BaseModel: models.BaseModel{
				ID: transportID,
			},
			WorkspaceID: workspaceID,
			Name:        "transport_name",
			State:       models.AdministrativeStateDeletePending,
			TransportVlans: models.TransportVlans{
				AVlan: 19,
				ZVlan: 19,
			},
			Product: models.TransportProduct{
				Product: models.Product{
					Provider:  "EQUINIX",
					Duration:  0,
					Location:  "EQUINIX FR5",
					Bandwidth: 100,
					PriceNRC:  0,
					PriceMRC:  0,
					CostNRC:   0,
					CostMRC:   0,
					SKU:       "CEQUFR5100AWS",
				},
				LocationTo: "EQUINIX LD5",
			},
			ConnectionID: "3091af46-3586-4cd1-bdbf-b569d2219823",
		},
	}
)

func TestCreateTransportSuccessfully(t *testing.T) {
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

	result := transportCreateResponse

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodPost, fmt.Sprintf("/accounts/%s/workspaces/%s/transports", accountId, workspaceID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusAccepted, transportCreateResponse),
		),
	)

	data, err := cli.CreateTransport(
		context.Background(),
		models.CreateTransport{
			Name: "transport_name",
			Product: models.AddProduct{
				SKU: "CEQUFR5100AWS",
			},
		},
		workspaceID,
	)

	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(*data).Should(Equal(result.Data))
}

func TestCreateTransportForbidden(t *testing.T) {
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
			gh.VerifyRequest(http.MethodPost, fmt.Sprintf("/accounts/%s/workspaces/%s/transports", accountId, workspaceID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusForbidden, nil),
		),
	)

	data, err := cli.CreateTransport(
		context.Background(),
		models.CreateTransport{
			Name: "transport_name",
			Product: models.AddProduct{
				SKU: "CEQUFR5100AWS",
			},
		},
		workspaceID,
	)

	g.Expect(err).ShouldNot(BeNil())
	g.Expect(data).Should(BeNil())
}

func TestCreateTransportFailedValidator(t *testing.T) {
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

	data, err := cli.CreateTransport(
		context.Background(),
		models.CreateTransport{
			Product: models.AddProduct{
				SKU: "CEQUFR5100AWS",
			},
		},
		workspaceID,
	)

	g.Expect(err.Error()).Should(Equal("Key: 'CreateTransport.Name' Error:Field validation for 'Name' failed on the 'required' tag"))
	g.Expect(data).Should(BeNil())

	data, err = cli.CreateTransport(
		context.Background(),
		models.CreateTransport{
			Name: "node_name",
		},
		workspaceID,
	)

	g.Expect(err.Error()).Should(Equal("Key: 'CreateTransport.Product.SKU' Error:Field validation for 'SKU' failed on the 'required' tag"))
	g.Expect(data).Should(BeNil())
}

func TestGetTransportSuccessfully(t *testing.T) {
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

	result := transportCreateResponse

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/workspaces/%s/transports/%s", accountId, workspaceID, transportID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusOK, transportCreateResponse),
		),
	)

	data, err := cli.GetTransport(
		context.Background(),
		workspaceID,
		transportID.String(),
	)

	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(*data).Should(Equal(result.Data))
}

func TestGetTransportCreationError(t *testing.T) {
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

	result := transportCreateErrorResponse

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/workspaces/%s/transports/%s", accountId, workspaceID, transportID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusOK, transportCreateErrorResponse),
		),
	)

	data, err := cli.GetTransport(
		context.Background(),
		workspaceID,
		transportID.String(),
	)

	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(*data).Should(Equal(result.Data))
}

func TestGetTransportNotFound(t *testing.T) {
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
			gh.VerifyRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/workspaces/%s/transports/%s", accountId, workspaceID, transportID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusNotFound, nil),
		),
	)

	data, err := cli.GetTransport(
		context.Background(),
		workspaceID,
		nodeID.String(),
	)

	g.Expect(err).ShouldNot(BeNil())
	g.Expect(data).Should(BeNil())
}

func TestUpdateTransportSuccessfully(t *testing.T) {
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

	result := transportUpdateResponse

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodPatch, fmt.Sprintf("/accounts/%s/workspaces/%s/transports/%s", accountId, workspaceID, transportID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusAccepted, transportUpdateResponse),
		),
	)

	data, err := cli.UpdateTransport(
		context.Background(),
		models.UpdateElement{
			Name: "transport_updated_name",
		},
		workspaceID,
		transportID.String(),
	)

	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(*data).Should(Equal(result.Data))
}

func TestUpdateTransportNotFound(t *testing.T) {
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
			gh.VerifyRequest(http.MethodPatch, fmt.Sprintf("/accounts/%s/workspaces/%s/transports/%s", accountId, workspaceID, transportID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusNotFound, nil),
		),
	)

	data, err := cli.UpdateTransport(
		context.Background(),
		models.UpdateElement{
			Name: "transport_updated_name",
		},
		workspaceID,
		transportID.String(),
	)

	g.Expect(err).ShouldNot(BeNil())
	g.Expect(data).Should(BeNil())
}

func TestDeleteTransportSuccessfully(t *testing.T) {
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

	result := transportDeleteResponse

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodDelete, fmt.Sprintf("/accounts/%s/workspaces/%s/transports/%s", accountId, workspaceID, transportID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusAccepted, transportDeleteResponse),
		),
	)

	data, err := cli.DeleteTransport(
		context.Background(),
		workspaceID,
		transportID.String(),
	)

	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(*data).Should(Equal(result.Data))
}

func TestDeleteTransportForbidden(t *testing.T) {
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
			gh.VerifyRequest(http.MethodDelete, fmt.Sprintf("/accounts/%s/workspaces/%s/transports/%s", accountId, workspaceID, transportID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusForbidden, nil),
		),
	)

	data, err := cli.DeleteTransport(
		context.Background(),
		workspaceID,
		transportID.String(),
	)

	g.Expect(err).ShouldNot(BeNil())
	g.Expect(data).Should(BeNil())
}
