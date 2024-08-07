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

const (
	workspaceID = "f6bee4c6-f7f0-493c-8e3c-5b6c6af6c2de"
)

var (
	nodeID = uuid.MustParse("20ea29c9-f892-4aaf-8907-de79fa83e7bb")

	nodeCreateResponse = models.NodeResponse{
		Data: models.Node{
			BaseModel: models.BaseModel{
				ID: nodeID,
			},
			AccountID:   accountId,
			WorkspaceID: workspaceID,
			Name:        "node_name",
			State:       models.AdministrativeStateCreationPending,
			Type:        models.NodeTypeCloud,
			Product: models.NodeProduct{
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
				CSPName:         "AWS",
				CSPCity:         "Frankfurt",
				CSPRegion:       "eu-central-1",
				CSPNameUnderlay: "AWS",
			},
			ProviderConfig: &models.ProviderCloudConfig{
				AccountID: "456789",
			},
		},
	}

	nodeUpdateResponse = models.NodeResponse{
		Data: models.Node{
			BaseModel: models.BaseModel{
				ID: nodeID,
			},
			AccountID:   accountId,
			WorkspaceID: workspaceID,
			Name:        "node_updated_name",
			State:       models.AdministrativeStateDeployed,
			Type:        models.NodeTypeCloud,
			Product: models.NodeProduct{
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
				CSPName:         "AWS",
				CSPCity:         "Frankfurt",
				CSPRegion:       "eu-central-1",
				CSPNameUnderlay: "AWS",
			},
			ProviderConfig: &models.ProviderCloudConfig{
				AccountID: "456789",
			},
		},
	}

	nodeDeleteResponse = models.NodeResponse{
		Data: models.Node{
			BaseModel: models.BaseModel{
				ID: nodeID,
			},
			AccountID:   accountId,
			WorkspaceID: workspaceID,
			Name:        "node_updated_name",
			State:       models.AdministrativeStateDeletePending,
			Type:        models.NodeTypeCloud,
			Product: models.NodeProduct{
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
				CSPName:         "AWS",
				CSPCity:         "Frankfurt",
				CSPRegion:       "eu-central-1",
				CSPNameUnderlay: "AWS",
			},
			ProviderConfig: &models.ProviderCloudConfig{
				AccountID: "456789",
			},
		},
	}
)

func TestCreateNodeSuccessfully(t *testing.T) {
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

	result := nodeCreateResponse

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodPost, fmt.Sprintf("/accounts/%s/workspaces/%s/nodes", accountId, workspaceID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusAccepted, nodeCreateResponse),
		),
	)

	data, err := cli.CreateNode(
		context.Background(),
		models.CreateNode{
			Name: "node_name",
			Type: models.NodeTypeCloud,
			Product: models.AddProduct{
				SKU: "CEQUFR5100AWS",
			},
			ProviderConfig: &models.ProviderCloudConfig{
				AccountID: "456789",
			},
		},
		workspaceID,
	)

	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(*data).Should(Equal(result.Data))
}

func TestCreateNodeForbidden(t *testing.T) {
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
			gh.VerifyRequest(http.MethodPost, fmt.Sprintf("/accounts/%s/workspaces/%s/nodes", accountId, workspaceID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusForbidden, nil),
		),
	)

	data, err := cli.CreateNode(
		context.Background(),
		models.CreateNode{
			Name: "node_name",
			Type: models.NodeTypeCloud,
			Product: models.AddProduct{
				SKU: "CEQUFR5100AWS",
			},
			ProviderConfig: &models.ProviderCloudConfig{
				AccountID: "456789",
			},
		},
		workspaceID,
	)

	g.Expect(err).ShouldNot(BeNil())
	g.Expect(data).Should(BeNil())
}

func TestCreateNodeFailedValidator(t *testing.T) {
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

	data, err := cli.CreateNode(
		context.Background(),
		models.CreateNode{
			Type: models.NodeTypeCloud,
			Product: models.AddProduct{
				SKU: "CEQUFR5100AWS",
			},
			ProviderConfig: &models.ProviderCloudConfig{
				AccountID: "456789",
			},
		},
		workspaceID,
	)

	g.Expect(err.Error()).Should(Equal("Key: 'CreateNode.Name' Error:Field validation for 'Name' failed on the 'required' tag"))
	g.Expect(data).Should(BeNil())

	data, err = cli.CreateNode(
		context.Background(),
		models.CreateNode{
			Name: "node_name",
			Product: models.AddProduct{
				SKU: "CEQUFR5100AWS",
			},
			ProviderConfig: &models.ProviderCloudConfig{
				AccountID: "456789",
			},
		},
		workspaceID,
	)

	g.Expect(err.Error()).Should(Equal("Key: 'CreateNode.Type' Error:Field validation for 'Type' failed on the 'required' tag"))
	g.Expect(data).Should(BeNil())

	data, err = cli.CreateNode(
		context.Background(),
		models.CreateNode{
			Name:    "node_name",
			Type:    models.NodeTypeCloud,
			Product: models.AddProduct{},
			ProviderConfig: &models.ProviderCloudConfig{
				AccountID: "456789",
			},
		},
		workspaceID,
	)

	g.Expect(err.Error()).Should(Equal("Key: 'CreateNode.Product.SKU' Error:Field validation for 'SKU' failed on the 'required' tag"))
	g.Expect(data).Should(BeNil())

	data, err = cli.CreateNode(
		context.Background(),
		models.CreateNode{
			Name: "node_name",
			Type: models.NodeTypeCloud,
			Product: models.AddProduct{
				SKU: "CEQUFR5100AWS",
			},
		},
		workspaceID,
	)

	g.Expect(err.Error()).Should(Equal("Key: 'CreateNode.ProviderConfig' Error:Field validation for 'ProviderConfig' failed on the 'required_if' tag"))
	g.Expect(data).Should(BeNil())
}

func TestGetNodeSuccessfully(t *testing.T) {
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

	result := nodeCreateResponse

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/workspaces/%s/nodes/%s", accountId, workspaceID, nodeID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusOK, nodeCreateResponse),
		),
	)

	data, err := cli.GetNode(
		context.Background(),
		workspaceID,
		nodeID.String(),
	)

	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(*data).Should(Equal(result.Data))
}

func TestGetNodeNotFound(t *testing.T) {
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
			gh.VerifyRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/workspaces/%s/nodes/%s", accountId, workspaceID, nodeID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusNotFound, nil),
		),
	)

	data, err := cli.GetNode(
		context.Background(),
		workspaceID,
		nodeID.String(),
	)

	g.Expect(err).ShouldNot(BeNil())
	g.Expect(data).Should(BeNil())
}

func TestUpdateNodeSuccessfully(t *testing.T) {
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

	result := nodeUpdateResponse

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodPatch, fmt.Sprintf("/accounts/%s/workspaces/%s/nodes/%s", accountId, workspaceID, nodeID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusAccepted, nodeUpdateResponse),
		),
	)

	data, err := cli.UpdateNode(
		context.Background(),
		models.UpdateNode{
			Name: "node_updated_name",
		},
		workspaceID,
		nodeID.String(),
	)

	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(*data).Should(Equal(result.Data))
}

func TestUpdateNodeNotFound(t *testing.T) {
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
			gh.VerifyRequest(http.MethodPatch, fmt.Sprintf("/accounts/%s/workspaces/%s/nodes/%s", accountId, workspaceID, nodeID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusNotFound, nil),
		),
	)

	data, err := cli.UpdateNode(
		context.Background(),
		models.UpdateNode{
			Name: "node_updated_name",
		},
		workspaceID,
		nodeID.String(),
	)

	g.Expect(err).ShouldNot(BeNil())
	g.Expect(data).Should(BeNil())
}

func TestDeleteNodeSuccessfully(t *testing.T) {
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

	result := nodeDeleteResponse

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodDelete, fmt.Sprintf("/accounts/%s/workspaces/%s/nodes/%s", accountId, workspaceID, nodeID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusAccepted, nodeDeleteResponse),
		),
	)

	data, err := cli.DeleteNode(
		context.Background(),
		workspaceID,
		nodeID.String(),
	)

	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(*data).Should(Equal(result.Data))
}

func TestDeleteNodeForbidden(t *testing.T) {
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
			gh.VerifyRequest(http.MethodDelete, fmt.Sprintf("/accounts/%s/workspaces/%s/nodes/%s", accountId, workspaceID, nodeID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusForbidden, nil),
		),
	)

	data, err := cli.DeleteNode(
		context.Background(),
		workspaceID,
		nodeID.String(),
	)

	g.Expect(err).ShouldNot(BeNil())
	g.Expect(data).Should(BeNil())
}
