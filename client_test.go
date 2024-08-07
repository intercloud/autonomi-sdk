package autonomisdk

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/intercloud/autonomi-sdk/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

const (
	personalAccessToken = "123456"
	accountId           = "c208a91f-72f4-4e2b-94f9-d311f369538a"
)

func TestNewClientSuccess(t *testing.T) {
	RegisterFailHandler(Fail)
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
			Timeout: 2 * time.Minute,
		}),
		WithPersonalAccessToken(personalAccessToken),
	)

	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(cli.personalAccessToken).To(Equal(personalAccessToken))
	g.Expect(cli.accountID.String()).To(Equal(accountId))
	g.Expect(cli.httpClient.Timeout).To(Equal(timeout))
}

func TestNewClientMissingPAT(t *testing.T) {
	RegisterFailHandler(Fail)
	g := NewWithT(t)
	server := ghttp.NewServer()
	defer server.Close()

	serverURL, err := url.Parse(server.URL())
	g.Expect(err).ShouldNot(HaveOccurred())

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
	)

	g.Expect(err).To(Equal(ErrPersonalAccessTokenRequired))
	g.Expect(cli).To(BeNil())
}

func TestNewClientMissingHostURL(t *testing.T) {
	RegisterFailHandler(Fail)
	g := NewWithT(t)
	server := ghttp.NewServer()
	defer server.Close()

	cli, err := NewClient(
		true,
		WithHTTPClient(&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, //nolint:gosec //No
				},
			},
		}),
		WithPersonalAccessToken(personalAccessToken),
	)

	g.Expect(err).To(Equal(ErrHostURLRequired))
	g.Expect(cli).To(BeNil())
}

func TestNewClientMissingTermsAndConditions(t *testing.T) {
	RegisterFailHandler(Fail)
	g := NewWithT(t)

	server := ghttp.NewServer()
	defer server.Close()

	serverURL, err := url.Parse(server.URL())
	g.Expect(err).ShouldNot(HaveOccurred())

	cli, err := NewClient(
		false,
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

	g.Expect(err).To(Equal(ErrTermsAndConditionsRequired))
	g.Expect(cli).To(BeNil())
}
