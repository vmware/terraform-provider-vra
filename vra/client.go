package vra

import (
	"fmt"
	"log"
	neturl "net/url"
	"os"
	"strings"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/vmware/vra-sdk-go/pkg/client"
	"github.com/vmware/vra-sdk-go/pkg/client/login"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

// Client the VRA Client
type Client struct {
	url       string
	apiClient *client.MulticloudIaaS
}

// NewClientFromRefreshToken configures and returns a VRA "Client" struct using "refresh_token" from provider config
func NewClientFromRefreshToken(url, refreshToken string, insecure bool) (interface{}, error) {
	token, err := getToken(url, refreshToken, insecure)
	if err != nil {
		return "", err
	}
	apiClient, err := getAPIClient(url, token, insecure)
	if err != nil {
		return "", err
	}
	return &Client{url, apiClient}, nil
}

// NewClientFromAccessToken configures and returns a VRA "Client" struct using "access_token" from provider config
func NewClientFromAccessToken(url, accessToken string, insecure bool) (interface{}, error) {
	apiClient, err := getAPIClient(url, accessToken, insecure)
	if err != nil {
		return "", err
	}
	return &Client{url, apiClient}, nil
}

func getToken(url, refreshToken string, insecure bool) (string, error) {
	parsedURL, err := neturl.Parse(url)
	if err != nil {
		return "", err
	}
	transport := httptransport.New(parsedURL.Host, "", nil)
	newTransport, err := httptransport.TLSTransport(httptransport.TLSClientOptions{
		InsecureSkipVerify: insecure,
	})
	if err != nil {
		return "", err
	}
	transport.Transport = newTransport
	transport.SetDebug(false)
	apiclient := client.New(transport, strfmt.Default)

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

func getAPIClient(url string, token string, insecure bool) (*client.MulticloudIaaS, error) {
	debug := false
	if os.Getenv("VRA_DEBUG") != "" {
		debug = true
	}

	parsedURL, err := neturl.Parse(url)
	if err != nil {
		return nil, err
	}
	transport := httptransport.New(parsedURL.Host, "", nil)
	transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Bearer "+token)
	newTransport, err := httptransport.TLSTransport(httptransport.TLSClientOptions{
		InsecureSkipVerify: insecure,
	})
	if err != nil {
		return nil, err
	}
	transport.Transport = newTransport
	if debug {
		transport.SetDebug(debug)
		transport.SetLogger(SwaggerLogger{})
	}
	apiclient := client.New(transport, strfmt.Default)
	return apiclient, nil
}
