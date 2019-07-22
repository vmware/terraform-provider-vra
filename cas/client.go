package cas

import (
	"fmt"
	"log"
	neturl "net/url"
	"os"
	"strings"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/vmware/cas-sdk-go/pkg/client"
	"github.com/vmware/cas-sdk-go/pkg/client/login"
	"github.com/vmware/cas-sdk-go/pkg/models"
)

// Client the CAS Client
type Client struct {
	url       string
	apiClient *client.MulticloudIaaS
}

// NewClientFromRefreshToken configures and returns a CAS "Client" struct using "refresh_token" from provider config
func NewClientFromRefreshToken(url, refreshToken string) (interface{}, error) {
	token, err := getToken(url, refreshToken)
	if err != nil {
		return "", err
	}
	return &Client{url, getAPIClient(url, token)}, nil
}

// NewClientFromAccessToken configures and returns a CAS "Client" struct using "access_token" from provider config
func NewClientFromAccessToken(url, accessToken string) (interface{}, error) {
	return &Client{url, getAPIClient(url, accessToken)}, nil
}

func getToken(url, refreshToken string) (string, error) {
	parsedURL, err := neturl.Parse(url)
	if err != nil {
		return "", err
	}
	transport := httptransport.New(parsedURL.Host, "", nil)
	transport.SetDebug(true)
	fmt.Printf("transport: %+v\n", transport)
	apiclient := client.New(transport, strfmt.Default)

	fmt.Printf("apiclient: %+v\n", apiclient)
	fmt.Printf("transport: %+v\n", apiclient.Transport)
	fmt.Printf("Login: %+v\n", apiclient.Login)
	params := login.NewRetrieveAuthTokenParams().WithBody(
		&models.CspLoginSpecification{
			RefreshToken: &refreshToken,
		},
	)
	authTokenResponse, err := apiclient.Login.RetrieveAuthToken(params)
	if err != nil || !strings.EqualFold(*authTokenResponse.Payload.TokenType, "bearer") {
		return "", err
	}

	return *authTokenResponse.Payload.Token, nil
}

// SwaggerLogger is the interface into the swagger logging facility which logs http traffic
type SwaggerLogger struct{}

// Printf is a swagger debug Printf
func (SwaggerLogger) Printf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	// Handle and split mixed "\r\n" and "\n"
	lines := strings.Split(strings.Replace(s, "\r\n", "\n", -1), "\n")

	for _, l := range lines {
		log.Printf("%s\n", l)
	}
}

// Debugf is a swagger debug logger
func (SwaggerLogger) Debugf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	// Handle and split mixed "\r\n" and "\n"
	lines := strings.Split(strings.Replace(s, "\r\n", "\n", -1), "\n")

	for _, l := range lines {
		log.Printf("%s\n", l)
	}
}

func getAPIClient(url string, token string) *client.MulticloudIaaS {
	debug := false
	if os.Getenv("CAS_DEBUG") != "" {
		debug = true
	}

	parsedURL, err := neturl.Parse(url)
	if err != nil {
		return nil
	}
	transport := httptransport.New(parsedURL.Host, "", nil)
	transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Bearer "+token)
	if debug {
		transport.SetDebug(debug)
		transport.SetLogger(SwaggerLogger{})
	}
	apiclient := client.New(transport, strfmt.Default)
	return apiclient
}
