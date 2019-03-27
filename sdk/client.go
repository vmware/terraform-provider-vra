package tango

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"strings"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/vmware/cas-sdk-go/pkg/client"
)

const (
	unsuccessfulRESTCall                        = "\nRequest: %s\nResponse status code: %d\nResponse body: %s"
	requestTrackerUnsuccessfulResourceOperation = "Request tracker returned: %s. Status: %s"

	cloudAccountsEndpoint = "/iaas/cloud-accounts"
	loginEndpoint         = "/iaas/login"
)

// Client used to store the Cloud Automation Services (cas) client and the provider configuration.
type Client struct {
	client       *http.Client
	base         string
	token        string
	projectID    string
	deploymentID string
	apiClient    *client.MulticloudIaaS
}

// GetProjectID returns the "project_id" from "cas" provider configuration
func (c *Client) GetProjectID() string {
	return c.projectID
}

// GetDeploymentID returns the "deployment_id" from "cas" provider configuration
func (c *Client) GetDeploymentID() string {
	return c.deploymentID
}

// NewClientFromRefreshToken configures and returns a CAS "Client" struct using "refresh_token" from provider config
func NewClientFromRefreshToken(url, refreshToken, projectID, deploymentID string) (interface{}, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 20,
	}

	loginRequest := LoginRequest{
		RefreshToken: refreshToken,
	}

	loginRequestBytes, err := json.Marshal(loginRequest)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url+loginEndpoint, bytes.NewReader(loginRequestBytes))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	response, err := netClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(unsuccessfulRESTCall, loginEndpoint+"\n"+formatRequest(request)+"\n"+
			string(loginRequestBytes), response.StatusCode, string(contents))
	}

	loginResponse := LoginResponse{}
	err = json.Unmarshal(contents, &loginResponse)
	if err != nil {
		return nil, err
	}

	return &Client{netClient, url, loginResponse.Token, projectID, deploymentID, getAPIClient(url, loginResponse.Token)}, nil
}

// NewClientFromAccessToken configures and returns a CAS "Client" struct using "access_token" from provider config
func NewClientFromAccessToken(url, accessToken, projectID, deploymentID string) (interface{}, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 20,
	}

	request, err := http.NewRequest(http.MethodGet, url+cloudAccountsEndpoint, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+accessToken)
	response, err := netClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(unsuccessfulRESTCall, cloudAccountsEndpoint+"\n"+formatRequest(request), response.StatusCode, string(contents))
	}

	return &Client{netClient, url, accessToken, projectID, deploymentID, getAPIClient(url, accessToken)}, nil
}

func getAPIClient(url string, token string) *client.MulticloudIaaS {
	parsedURL, err := neturl.Parse(url)
	if err != nil {
		return nil
	}
	transport := httptransport.New(parsedURL.Host, "", nil)
	transport.SetDebug(true)
	transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Bearer "+token)
	apiclient := client.New(transport, strfmt.Default)
	return apiclient
}

// GetAPIClient returns the api client needed for the new CAS sdk
func (c *Client) GetAPIClient() *client.MulticloudIaaS {
	return c.apiClient
}

// CreateResource creates a cas resource with the given resource specification.
func (c *Client) CreateResource(resourceSpecification ResourceSpecification) (interface{}, error) {
	resourceSpecificationBytes, err := json.Marshal(resourceSpecification)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, c.base+resourceSpecification.GetEndpoint(), bytes.NewReader(resourceSpecificationBytes))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+c.token)
	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(unsuccessfulRESTCall, resourceSpecification.GetEndpoint()+" create resource\n"+formatRequest(request)+"\n"+
			string(resourceSpecificationBytes), response.StatusCode, string(contents))
	}

	createResourceEndpointResponse := CreateResourceEndpointResponse{}
	err = json.Unmarshal(contents, &createResourceEndpointResponse)
	if err != nil {
		return nil, err
	}

	requestTrackerResponse, err := c.getFinishedRequestTrackerResponse(createResourceEndpointResponse.SelfLink)
	if err != nil {
		return nil, err
	}

	if len(requestTrackerResponse.Resources) == 0 {
		return nil, fmt.Errorf("%#v", requestTrackerResponse)
	}

	return c.ReadResource(requestTrackerResponse.Resources[0])
}

// ReadResource returns the resource description for the given resource endpoint
func (c *Client) ReadResource(resourceEndpoint string) (interface{}, error) {
	request, err := http.NewRequest(http.MethodGet, c.base+resourceEndpoint, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+c.token)
	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(unsuccessfulRESTCall, resourceEndpoint+" read resource\n"+formatRequest(request), response.StatusCode, string(contents))
	}

	resourceObject, err := NewResourceObject(resourceEndpoint)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(contents, &resourceObject)
	if err != nil {
		return nil, err
	}

	return resourceObject, nil
}

// DeleteResource deletes the given resource
func (c *Client) DeleteResource(resourceEndpoint string) error {
	request, err := http.NewRequest(http.MethodDelete, c.base+resourceEndpoint, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+c.token)
	response, err := c.client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(unsuccessfulRESTCall, resourceEndpoint+" delete resource\n"+formatRequest(request), response.StatusCode, string(contents))
	}

	createResourceEndpointResponse := CreateResourceEndpointResponse{}
	err = json.Unmarshal(contents, &createResourceEndpointResponse)
	if err != nil {
		return err
	}

	_, err = c.getFinishedRequestTrackerResponse(createResourceEndpointResponse.SelfLink)

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) getFinishedRequestTrackerResponse(endpoint string) (*RequestTrackerResponse, error) {
	attempts := 0
	requestTrackerResponse := RequestTrackerResponse{}
	for requestTrackerResponse.Status != "FINISHED" && requestTrackerResponse.Status != "FAILED" && attempts < 100 {
		time.Sleep(6 * time.Second)
		req, err := http.NewRequest(http.MethodGet, c.base+endpoint, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+c.token)

		res, err := c.client.Do(req)
		if err != nil {
			return nil, err
		}

		defer res.Body.Close()

		contents, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf(unsuccessfulRESTCall, endpoint+" request tracker\n"+formatRequest(req), res.StatusCode, string(contents))
		}

		err = json.Unmarshal(contents, &requestTrackerResponse)
		if err != nil {
			return nil, err
		}
		attempts++
	}

	if requestTrackerResponse.Progress != 100 || requestTrackerResponse.Status == "FAILED" {
		return nil, fmt.Errorf(requestTrackerUnsuccessfulResourceOperation, requestTrackerResponse.Message, requestTrackerResponse.Status)
	}

	return &requestTrackerResponse, nil
}

func formatRequest(r *http.Request) string {
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)

	var request []string
	request = append(request, url)
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	return strings.Join(request, "\n")
}
