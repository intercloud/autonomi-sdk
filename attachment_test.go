package autonomisdk

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/intercloud/autonomi-sdk/models"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

const (
	aSide = "A"
	zSide = "Z"
)

var (
	attachmentID = uuid.MustParse("20ea29c9-f892-4aaf-8907-de79fa83e7bb")

	attachmentCreateResponse = models.AttachmentResponse{
		Data: models.Attachment{
			BaseModel: models.BaseModel{
				ID: attachmentID,
			},
			WorkspaceID: workspaceID,
			NodeID:      nodeID.String(),
			TransportID: transportID.String(),
			State:       models.AdministrativeStateCreationPending,
			Side:        aSide,
		},
	}

	attachmentCreationErrorResponse = models.AttachmentResponse{
		Data: models.Attachment{
			BaseModel: models.BaseModel{
				ID: attachmentID,
			},
			WorkspaceID: workspaceID,
			NodeID:      nodeID.String(),
			TransportID: transportID.String(),
			State:       models.AdministrativeStateCreationError,
			Error: &models.SupportError{
				Code: "ERR_INTERNAL",
				Msg:  "an internal error occured",
			},
		},
	}

	attachmentDeployedResponse = models.AttachmentResponse{
		Data: models.Attachment{
			BaseModel: models.BaseModel{
				ID: attachmentID,
			},
			WorkspaceID: workspaceID,
			NodeID:      nodeID.String(),
			TransportID: transportID.String(),
			State:       models.AdministrativeStateDeployed,
			Side:        aSide,
		},
	}

	attachmentDeletePendingResponse = models.AttachmentResponse{
		Data: models.Attachment{
			BaseModel: models.BaseModel{
				ID: attachmentID,
			},
			WorkspaceID: workspaceID,
			NodeID:      nodeID.String(),
			TransportID: transportID.String(),
			State:       models.AdministrativeStateDeletePending,
		},
	}

	attachmentDeleteErrorResponse = models.AttachmentResponse{
		Data: models.Attachment{
			BaseModel: models.BaseModel{
				ID: attachmentID,
			},
			WorkspaceID: workspaceID,
			NodeID:      nodeID.String(),
			TransportID: transportID.String(),
			State:       models.AdministrativeStateDeleteError,
		},
	}

	attachmentDeleteResponse = models.AttachmentResponse{
		Data: models.Attachment{},
	}
)

func TestCreateAttachmentSuccessfully(t *testing.T) {
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

	result := attachmentCreateResponse

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodPost, fmt.Sprintf("/accounts/%s/workspaces/%s/attachments", accountId, workspaceID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusAccepted, attachmentCreateResponse),
		),
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/workspaces/%s/attachments/%s", accountId, workspaceID, attachmentID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusOK, attachmentCreateResponse),
		),
	)

	data, err := cli.CreateAttachment(
		context.Background(),
		models.CreateAttachment{
			NodeID:      nodeID.String(),
			TransportID: transportID.String(),
		},
		workspaceID,
	)

	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(*data).Should(Equal(result.Data))
}

func TestCreateAttachmentWaitForStateDeployed(t *testing.T) {
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

	cli.poll.maxRetry = 2
	cli.poll.retryInterval = 1 * time.Second

	result := attachmentDeployedResponse

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodPost, fmt.Sprintf("/accounts/%s/workspaces/%s/attachments", accountId, workspaceID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusAccepted, attachmentCreateResponse),
		),
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/workspaces/%s/attachments/%s", accountId, workspaceID, attachmentID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusOK, attachmentCreateResponse),
		),
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/workspaces/%s/attachments/%s", accountId, workspaceID, attachmentID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusOK, attachmentDeployedResponse),
		),
	)

	data, err := cli.CreateAttachment(
		context.Background(),
		models.CreateAttachment{
			NodeID:      nodeID.String(),
			TransportID: transportID.String(),
		},
		workspaceID,
		WithWaitUntilElementDeployed(),
	)

	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(*data).Should(Equal(result.Data))
}

func TestCreateAttachmentCreationError(t *testing.T) {
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

	cli.poll.maxRetry = 2
	cli.poll.retryInterval = 1 * time.Second

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodPost, fmt.Sprintf("/accounts/%s/workspaces/%s/attachments", accountId, workspaceID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusAccepted, attachmentCreateResponse),
		),
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/workspaces/%s/attachments/%s", accountId, workspaceID, attachmentID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusOK, attachmentCreateResponse),
		),
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/workspaces/%s/attachments/%s", accountId, workspaceID, attachmentID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusOK, attachmentCreationErrorResponse),
		),
	)

	data, err := cli.CreateAttachment(
		context.Background(),
		models.CreateAttachment{
			NodeID:      nodeID.String(),
			TransportID: transportID.String(),
		},
		workspaceID,
		WithWaitUntilElementDeployed(),
	)

	g.Expect(err).ShouldNot(BeNil())
	g.Expect(data).Should(BeNil())
}

func TestCreateAttachmentWaitForStateTimeout(t *testing.T) {
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
	cli.poll.maxRetry = 1
	cli.poll.retryInterval = 1 * time.Second

	g.Expect(err).ShouldNot(HaveOccurred())

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodPost, fmt.Sprintf("/accounts/%s/workspaces/%s/attachments", accountId, workspaceID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusAccepted, attachmentCreateResponse),
		),
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/workspaces/%s/attachments/%s", accountId, workspaceID, attachmentID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusOK, attachmentCreateResponse),
		),
	)

	_, err = cli.CreateAttachment(
		context.Background(),
		models.CreateAttachment{
			NodeID:      nodeID.String(),
			TransportID: transportID.String(),
		},
		workspaceID,
		WithWaitUntilElementDeployed(),
	)

	g.Expect(err).ShouldNot(BeNil())
}

func TestCreateAttachmentPollNotFound(t *testing.T) {
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
	cli.poll.maxRetry = 1
	cli.poll.retryInterval = 1 * time.Second

	g.Expect(err).ShouldNot(HaveOccurred())

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodPost, fmt.Sprintf("/accounts/%s/workspaces/%s/attachments", accountId, workspaceID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusAccepted, attachmentCreateResponse),
		),
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/workspaces/%s/attachments/%s", accountId, workspaceID, attachmentID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusNotFound, nil),
		),
	)

	_, err = cli.CreateAttachment(
		context.Background(),
		models.CreateAttachment{
			NodeID:      nodeID.String(),
			TransportID: transportID.String(),
		},
		workspaceID,
		WithWaitUntilElementDeployed(),
	)

	g.Expect(err).ShouldNot(BeNil())
}

func TestCreateAttachmentForbidden(t *testing.T) {
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
			gh.VerifyRequest(http.MethodPost, fmt.Sprintf("/accounts/%s/workspaces/%s/attachments", accountId, workspaceID)),
			gh.RespondWithJSONEncoded(http.StatusForbidden, nil),
		),
	)

	data, err := cli.CreateAttachment(
		context.Background(),
		models.CreateAttachment{
			NodeID:      nodeID.String(),
			TransportID: transportID.String(),
		},
		workspaceID,
	)

	g.Expect(err).ShouldNot(BeNil())
	g.Expect(data).Should(BeNil())
}

func TestCreateAttachmentFailedValidator(t *testing.T) {
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

	data, err := cli.CreateAttachment(
		context.Background(),
		models.CreateAttachment{
			NodeID: nodeID.String(),
		},
		workspaceID,
	)

	g.Expect(err.Error()).Should(Equal("Key: 'CreateAttachment.TransportID' Error:Field validation for 'TransportID' failed on the 'required' tag"))
	g.Expect(data).Should(BeNil())

	data, err = cli.CreateAttachment(
		context.Background(),
		models.CreateAttachment{
			TransportID: transportID.String(),
		},
		workspaceID,
	)

	g.Expect(err.Error()).Should(Equal("Key: 'CreateAttachment.NodeID' Error:Field validation for 'NodeID' failed on the 'required' tag"))
	g.Expect(data).Should(BeNil())
}

func TestGetAttachmentSuccessfully(t *testing.T) {
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

	result := attachmentCreateResponse

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/workspaces/%s/attachments/%s", accountId, workspaceID, attachmentID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusOK, attachmentCreateResponse),
		),
	)

	data, err := cli.GetAttachment(
		context.Background(),
		workspaceID,
		attachmentID.String(),
	)

	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(*data).Should(Equal(result.Data))
}

func TestGetAttachmentCreationError(t *testing.T) {
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

	result := attachmentCreationErrorResponse

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/workspaces/%s/attachments/%s", accountId, workspaceID, attachmentID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusOK, attachmentCreationErrorResponse),
		),
	)

	data, err := cli.GetAttachment(
		context.Background(),
		workspaceID,
		attachmentID.String(),
	)

	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(*data).Should(Equal(result.Data))
}

func TestGetAttachmentNotFound(t *testing.T) {
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
			gh.VerifyRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/workspaces/%s/attachments/%s", accountId, workspaceID, attachmentID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusNotFound, nil),
		),
	)

	data, err := cli.GetAttachment(
		context.Background(),
		workspaceID,
		attachmentID.String(),
	)

	g.Expect(err).ShouldNot(BeNil())
	g.Expect(data).Should(BeNil())
}

func TestDeleteAttachmentSuccessfully(t *testing.T) {
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

	result := attachmentDeletePendingResponse

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodDelete, fmt.Sprintf("/accounts/%s/workspaces/%s/attachments/%s", accountId, workspaceID, attachmentID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusAccepted, attachmentDeletePendingResponse),
		),
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/workspaces/%s/attachments/%s", accountId, workspaceID, attachmentID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusOK, attachmentDeletePendingResponse),
		),
	)

	data, err := cli.DeleteAttachment(
		context.Background(),
		workspaceID,
		attachmentID.String(),
	)

	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(*data).Should(Equal(result.Data))
}

func TestDeleteAttachmentWaitForStateDeleted(t *testing.T) {
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

	cli.poll.maxRetry = 1
	cli.poll.retryInterval = 1 * time.Second

	result := attachmentDeleteResponse

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodDelete, fmt.Sprintf("/accounts/%s/workspaces/%s/attachments/%s", accountId, workspaceID, attachmentID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusAccepted, attachmentDeletePendingResponse),
		),
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/workspaces/%s/attachments/%s", accountId, workspaceID, attachmentID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusNotFound, nil),
		),
	)

	data, err := cli.DeleteAttachment(
		context.Background(),
		workspaceID,
		attachmentID.String(),
		WithWaitUntilElementUndeployed(),
	)

	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(*data).Should(Equal(result.Data))
}

func TestDeleteAttachmentDeleteError(t *testing.T) {
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

	cli.poll.maxRetry = 1
	cli.poll.retryInterval = 1 * time.Second

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodDelete, fmt.Sprintf("/accounts/%s/workspaces/%s/attachments/%s", accountId, workspaceID, attachmentID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusAccepted, attachmentDeletePendingResponse),
		),
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/workspaces/%s/attachments/%s", accountId, workspaceID, attachmentID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusOK, attachmentDeleteErrorResponse),
		),
	)

	data, err := cli.DeleteAttachment(
		context.Background(),
		workspaceID,
		attachmentID.String(),
		WithWaitUntilElementUndeployed(),
	)

	g.Expect(err).ShouldNot(BeNil())
	g.Expect(data).Should(BeNil())
}

func TestDeleteAttachmentWaitForStateTimeout(t *testing.T) {
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

	cli.poll.maxRetry = 1
	cli.poll.retryInterval = 1 * time.Second

	server.AppendHandlers(
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodDelete, fmt.Sprintf("/accounts/%s/workspaces/%s/attachments/%s", accountId, workspaceID, attachmentID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusAccepted, attachmentDeletePendingResponse),
		),
		ghttp.CombineHandlers(
			gh.VerifyRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/workspaces/%s/attachments/%s", accountId, workspaceID, attachmentID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusOK, attachmentDeletePendingResponse),
		),
	)

	_, err = cli.DeleteAttachment(
		context.Background(),
		workspaceID,
		attachmentID.String(),
		WithWaitUntilElementUndeployed(),
	)

	g.Expect(err).ShouldNot(BeNil())
}

func TestDeleteAttachmentForbidden(t *testing.T) {
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
			gh.VerifyRequest(http.MethodDelete, fmt.Sprintf("/accounts/%s/workspaces/%s/attachments/%s", accountId, workspaceID, attachmentID)),
			gh.VerifyHeaderKV("Authorization", "Bearer "+personalAccessToken), //nolint
			gh.RespondWithJSONEncoded(http.StatusForbidden, nil),
		),
	)

	data, err := cli.DeleteAttachment(
		context.Background(),
		workspaceID,
		attachmentID.String(),
	)

	g.Expect(err).ShouldNot(BeNil())
	g.Expect(data).Should(BeNil())
}
