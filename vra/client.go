package vra

import (
	"fmt"
	"log"
	"net/http"
	neturl "net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/vmware/vra-sdk-go/pkg/client"
	"github.com/vmware/vra-sdk-go/pkg/client/login"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

// API Versions
const (
	CatalogAPIVersion     = "2019-01-15"
	DeploymentsAPIVersion = "2019-01-15"
)

const IncreasedTimeOut = 60 * time.Second

type ReauthTimeout struct {
	mu      sync.Mutex
	seconds time.Duration
	f       func()
	timer   *time.Timer
	reload  bool
}

func InitializeTimeout(d time.Duration) *ReauthTimeout {
	t := ReauthTimeout{seconds: d, reload: false}

	if d.Seconds() != 0 {
		t.f = func() {
			t.mu.Lock()
			t.reload = true
			t.mu.Unlock()
		}
		t.timer = time.AfterFunc(d, t.f)
	}

	return &t
}

func (t *ReauthTimeout) ShouldReload() bool {
	t.mu.Lock()
	reload := t.reload
	if reload {
		// reset the timer
		t.reload = false
		t.timer = time.AfterFunc(t.seconds, t.f)
	}
	t.mu.Unlock()
	return reload
}

type ReauthorizeRuntime struct {
	origClient   httptransport.Runtime
	url          string
	refreshToken string
	insecure     bool
	reauthtimer  *ReauthTimeout
}

// Submit implements the ClientTransport interface as a wrapper to retry a 401 with a new token.
func (r *ReauthorizeRuntime) Submit(operation *runtime.ClientOperation) (interface{}, error) {
	if r.reauthtimer.ShouldReload() {
		log.Printf("Reauthorize timer expired, generating a new access token")
		token, tokenErr := getToken(r.url, r.refreshToken, r.insecure)
		if tokenErr != nil {
			return nil, tokenErr
		}

		// Fix up the Authorization header with the new token
		r.origClient.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Bearer "+token)
	}

	result, err := r.origClient.Submit(operation)
	if err == nil {
		return result, err
	}

	// Check if an error but not a 401, then return results. Errors strings checked are if 401 is implemented in the swagger API or not.
	if !(strings.Contains(err.Error(), "[401]") || strings.Contains(err.Error(), "unknown error (status 401)")) {
		return result, err
	}

	// We have a 401 with a refresh token, let's try refreshing once and try again
	log.Printf("Response back was a 401, trying again with new access token")
	token, tokenErr := getToken(r.url, r.refreshToken, r.insecure)
	if tokenErr != nil {
		return result, err
	}

	// Fix up the Authorization header with the new token and resubmit the request
	r.origClient.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Bearer "+token)
	result, err = r.origClient.Submit(operation)
	return result, err
}

// Client the VRA Client
type Client struct {
	url       string
	apiClient *client.MulticloudIaaS
}

// NewClientFromRefreshToken configures and returns a VRA "Client" struct using "refresh_token" from provider config
func NewClientFromRefreshToken(url, refreshToken string, insecure bool, reauth string) (interface{}, error) {
	token, err := getToken(url, refreshToken, insecure)
	if err != nil {
		return "", err
	}
	apiClient, err := getAPIClient(url, token, insecure)
	if err != nil {
		return "", err
	}

	t := apiClient.Transport.(*httptransport.Runtime)
	reautDuration, err := time.ParseDuration(reauth)

	if err != nil {
		return "", err
	}
	apiClient.SetTransport(&ReauthorizeRuntime{*t, url, refreshToken, insecure, InitializeTimeout(reautDuration)})

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
	transport := httptransport.New(parsedURL.Host, parsedURL.Path, nil)
	transport.SetDebug(false)
	transport.Transport, err = createTransport(insecure)
	if err != nil {
		return "", err
	}
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

func createTransport(insecure bool) (http.RoundTripper, error) {
	cfg, err := httptransport.TLSClientAuth(httptransport.TLSClientOptions{
		InsecureSkipVerify: insecure,
	})
	if err != nil {
		return nil, err
	}

	return &http.Transport{
		TLSClientConfig: cfg,
		Proxy:           http.ProxyFromEnvironment,
	}, nil
}

func getAPIClient(url string, token string, insecure bool) (*client.MulticloudIaaS, error) {
	parsedURL, err := neturl.Parse(url)
	if err != nil {
		return nil, err
	}
	t := httptransport.New(parsedURL.Host, parsedURL.Path, nil)
	t.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Bearer "+token)
	newTransport, err := createTransport(insecure)
	if err != nil {
		return nil, err
	}

	// Setup logging through the terraform helper
	t.Transport = logging.NewTransport("VRA", newTransport)
	t.SetDebug(true)
	t.SetLogger(SwaggerLogger{})
	apiclient := client.New(t, strfmt.Default)
	return apiclient, nil
}
