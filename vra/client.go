// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	neturl "net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/vmware/vra-sdk-go/pkg/client"
	"github.com/vmware/vra-sdk-go/pkg/client/login"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

// API Versions
const (
	IaaSAPIVersion        = "2021-07-15"
	CatalogAPIVersion     = "2020-08-25"
	DeploymentsAPIVersion = "2020-08-25"
)

const DefaultDollarTop = 1000

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
	organization string
	refreshToken string
	insecure     bool
	reauthtimer  *ReauthTimeout
}

// Submit implements the ClientTransport interface as a wrapper to retry a 401 with a new token.
func (r *ReauthorizeRuntime) Submit(operation *runtime.ClientOperation) (interface{}, error) {
	if r.reauthtimer.ShouldReload() {
		log.Printf("Reauthorize timer expired, generating a new access token")
		token, tokenErr := getToken(r.url, r.organization, r.refreshToken, r.insecure)
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
	token, tokenErr := getToken(r.url, r.organization, r.refreshToken, r.insecure)
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
	apiClient *client.API
}

// NewClientFromRefreshToken configures and returns a VRA "Client" struct using "refresh_token" from provider config
func NewClientFromRefreshToken(url, organization, refreshToken string, insecure bool, reauth string, apiTimeout int) (interface{}, error) {
	token, err := getToken(url, organization, refreshToken, insecure)
	if err != nil {
		return "", err
	}
	apiClient, err := getAPIClient(url, token, insecure, apiTimeout)
	if err != nil {
		return "", err
	}

	t := apiClient.Transport.(*httptransport.Runtime)
	reautDuration, err := time.ParseDuration(reauth)

	if err != nil {
		return "", err
	}
	apiClient.SetTransport(&ReauthorizeRuntime{*t, url, organization, refreshToken, insecure, InitializeTimeout(reautDuration)})

	return &Client{url, apiClient}, nil
}

// NewClientFromAccessToken configures and returns a VRA "Client" struct using "access_token" from provider config
func NewClientFromAccessToken(url, accessToken string, insecure bool, apiTimeout int) (interface{}, error) {
	apiClient, err := getAPIClient(url, accessToken, insecure, apiTimeout)
	if err != nil {
		return "", err
	}
	return &Client{url, apiClient}, nil
}

func getToken(url, organization, refreshToken string, insecure bool) (string, error) {
	isVCFA, err := isVCFA(url, insecure)
	if err != nil {
		return "", fmt.Errorf("error determining whether vRA or VCFA: %s", err)
	}
	if isVCFA {
		if organization == "" {
			return "", errors.New("organization is required for VCFA")
		} else if organization == "system" {
			return "", errors.New("system organization is not allowed")
		}
		return getVCFAToken(url, organization, refreshToken, insecure)
	}
	return getVRAToken(url, refreshToken, insecure)
}

func isVCFA(url string, insecure bool) (bool, error) {
	parsedURL, err := neturl.Parse(url)
	if err != nil {
		return false, fmt.Errorf("error parsing the URL %s: %s", url, err)
	}
	transport, err := createTransport(insecure)
	if err != nil {
		return false, fmt.Errorf("error creating an http transport: %s", err)
	}
	client := &http.Client{Transport: transport}
	response, err := client.Get(fmt.Sprintf("%s://%s/automation/config.json", parsedURL.Scheme, parsedURL.Host))
	if err != nil {
		return false, fmt.Errorf("error retrieving the configuration from url %s: %s", url, err)
	}
	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return false, nil
		}
		return false, fmt.Errorf("error retrieving the configuration from url %s, http response code is %d", url, response.StatusCode)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return false, fmt.Errorf("error reading the http response body from url %s: %s", url, err)
	}
	var config map[string]interface{}
	if err := json.Unmarshal(body, &config); err != nil {
		return false, fmt.Errorf("error unmarshalling the configuration from url %s: %s", url, err)
	}
	applicationVersion, ok := config["applicationVersion"]
	if ok {
		productName, err := base64.StdEncoding.DecodeString(applicationVersion.(string))
		if err != nil {
			return false, fmt.Errorf("error decoding the application version %s: %s", applicationVersion, err)
		}
		re := regexp.MustCompile(`v?(\d+\.\d+\.\d+\.\d+)`)
		match := re.FindStringSubmatch(string(productName))
		if match != nil {
			productVersion, err := version.NewVersion(match[1])
			if err != nil {
				return false, fmt.Errorf("error parsing the application version %s: %s", match[1], err)
			}
			vcfaVersion, _ := version.NewVersion("9.0.0.0")
			if productVersion.GreaterThanOrEqual(vcfaVersion) {
				return true, nil
			}
		}
	}
	return false, nil
}

// Retrieve the access token for vRA 8.x instances
func getVRAToken(url, refreshToken string, insecure bool) (string, error) {
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

// Retrieve the access token for VCFA 9.x instances
func getVCFAToken(url, org string, refreshToken string, insecure bool) (string, error) {
	parsedURL, err := neturl.Parse(url)
	if err != nil {
		return "", fmt.Errorf("error parsing the URL %s: %s", url, err)
	}
	transport, err := createTransport(insecure)
	if err != nil {
		return "", fmt.Errorf("error creating an http transport: %s", err)
	}
	client := &http.Client{Transport: transport}

	data := neturl.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	tenant := "tenant/" + org
	if strings.EqualFold(org, "system") {
		tenant = "provider"
	}
	response, err := client.Post(fmt.Sprintf("%s://%s/tm/oauth/%s/token", parsedURL.Scheme, parsedURL.Host, tenant), "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("error retrieving the auth token: %s", err)
	}
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error retrieving the auth token, http response code is %d", response.StatusCode)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error reading the http response body: %s", err)
	}
	var tokenData map[string]interface{}
	if err := json.Unmarshal(body, &tokenData); err != nil {
		return "", fmt.Errorf("error unmarshalling the token data: %s", err)
	}
	accessToken, ok := tokenData["access_token"]
	if ok {
		return accessToken.(string), nil
	}
	return "", errors.New("Unable to obtain an access token")
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

func getAPIClient(url string, token string, insecure bool, apiTimeout int) (*client.API, error) {
	parsedURL, err := neturl.Parse(url)
	if err != nil {
		return nil, err
	}

	if apiTimeout > 0 {
		httptransport.DefaultTimeout = time.Duration(apiTimeout) * time.Second
	}

	t := httptransport.New(parsedURL.Host, parsedURL.Path, nil)
	t.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Bearer "+token)
	newTransport, err := createTransport(insecure)
	if err != nil {
		return nil, err
	}

	// Setup logging through the terraform helper
	t.Transport = logging.NewSubsystemLoggingHTTPTransport("VRA", newTransport)
	t.SetDebug(true)
	t.SetLogger(SwaggerLogger{})
	apiclient := client.New(t, strfmt.Default)
	return apiclient, nil
}
