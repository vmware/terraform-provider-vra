package tango

import (
	"errors"
	"regexp"
)

const (
	machineEndpoint      = "/iaas/api/machines"
	blockDeviceEndpoint  = "/iaas/api/block-devices"
	networkEndpoint      = "/iaas/api/networks"
	loadBalancerEndpoint = "/iaas/api/load-balancers"
)

var (
	machineEndpointDisksR = regexp.MustCompile("^/iaas/api/machines/[a-zA-Z0-9]+/disks")
	machineEndpointR      = regexp.MustCompile("^/iaas/api/machines")
	networkEndpointR      = regexp.MustCompile("^/iaas/api/networks")
	blockDeviceEndpointR  = regexp.MustCompile("^/iaas/api/block-devices")
	loadBalancerEndpointR = regexp.MustCompile("^/iaas/api/load-balancers")
)

// Tag consists of key, value pair.
type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Constraint represents the condition that one resource may have on other resources.
type Constraint struct {
	Mandatory  bool   `json:"mandatory"`
	Expression string `json:"expression"`
}

// Nic represents the machine network interface configuration.
type Nic struct {
	Name             string            `json:"name,omitempty"`
	Description      string            `json:"description,omitempty"`
	DeviceIndex      int               `json:"deviceIndex,omitempty"`
	NetworkID        string            `json:"networkId"`
	Addresses        []string          `json:"addresses,omitempty"`
	SecurityGroupIDs []string          `json:"securityGroupIds,omitempty"`
	CustomProperties map[string]string `json:"customProperties,omitempty"`
}

// Disk represents the disk configuration.
type Disk struct {
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	BlockDeviceID string `json:"blockDeviceId"`
}

// Route represents the configuration for routing incoming requests to the back-end instances.
// A load balancer may support multiple such configurations.
type Route struct {
	Protocol       string                    `json:"protocol"`
	Port           string                    `json:"port"`
	MemberProtocol string                    `json:"memberProtocol"`
	MemberPort     string                    `json:"memberPort"`
	HCC            *HealthCheckConfiguration `json:"healthCheckConfiguration,omitempty"`
}

// HealthCheckConfiguration represents a load balancer configuration for checking the health of
// the load-balanced back-end instances.
type HealthCheckConfiguration struct {
	Protocol          string `json:"protocol"`
	Port              string `json:"port"`
	URLPath           string `json:"urlPath,omitempty"`
	IntervalSeconds   int    `json:"intervalSeconds,omitempty"`
	TimeoutSeconds    int    `json:"timeoutSeconds,omitempty"`
	UnhealthThreshold int    `json:"unhealthyThreshold,omitempty"`
	HealthThreshold   int    `json:"healthyThreshold,omitempty"`
}

// LoginRequest consits of the refresh token used to login to CAS.
type LoginRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// LoginResponse consists of a token and its type after the request is successfully authenticated.
type LoginResponse struct {
	TokenType string `json:"tokenType"`
	Token     string `json:"token"`
}

// LoadBalancerSpecification represents the cas_load_balancer resource configuration.
type LoadBalancerSpecification struct {
	Name             string            `json:"name,omitempty"`
	Description      string            `json:"description,omitempty"`
	DeploymentID     string            `json:"deploymentId,omitempty"`
	ProjectID        string            `json:"projectId"`
	InternetFacing   bool              `json:"internetFacing,omitempty"`
	CustomProperties map[string]string `json:"customProperties,omitempty"`
	Tags             []Tag             `json:"tags,omitempty"`
	TargetLinks      []string          `json:"targetLinks,omitempty"`
	Nics             []Nic             `json:"nics"`
	Routes           []Route           `json:"routes"`
}

// MachineSpecification represents the cas_machine resource configuration.
type MachineSpecification struct {
	Name             string            `json:"name,omitempty"`
	Description      string            `json:"description,omitempty"`
	DeploymentID     string            `json:"deploymentId,omitempty"`
	Image            string            `json:"image"`
	Flavor           string            `json:"flavor"`
	ProjectID        string            `json:"projectId"`
	MachineCount     int               `json:"machineCount,omitempty"`
	Constraints      []Constraint      `json:"constraints,omitempty"`
	CustomProperties map[string]string `json:"customProperties,omitempty"`
	Tags             []Tag             `json:"tags,omitempty"`
	Nics             []Nic             `json:"nics"`
	Disks            []Disk            `json:"disks,omitempty"`
	BootConfig       map[string]string `json:"bootConfig,omitempty"`
}

// NetworkSpecification represents the cas_network resource configuration.
type NetworkSpecification struct {
	Name             string            `json:"name,omitempty"`
	Description      string            `json:"description,omitempty"`
	DeploymentID     string            `json:"deploymentId,omitempty"`
	ProjectID        string            `json:"projectId"`
	Constraints      []Constraint      `json:"constraints,omitempty"`
	CustomProperties map[string]string `json:"customProperties,omitempty"`
	Tags             []Tag             `json:"tags,omitempty"`
	OutboundAccess   bool              `json:"outboundAccess,omitempty"`
}

// BlockDeviceSpecification represents the cas_block_device resource configuration.
type BlockDeviceSpecification struct {
	Name              string            `json:"name,omitempty"`
	Description       string            `json:"description,omitempty"`
	DeploymentID      string            `json:"deploymentId,omitempty"`
	ProjectID         string            `json:"projectId"`
	CapacityInGB      int               `json:"capacityInGB"`
	Encrypted         bool              `json:"encrypted,omitempty"`
	SourceReference   string            `json:"sourceReference,omitempty"`
	DiskContentBase64 string            `json:"diskContentBase64,omitempty"`
	CustomProperties  map[string]string `json:"customProperties,omitempty"`
	Constraints       []Constraint      `json:"constraints,omitempty"`
	Tags              []Tag             `json:"tags,omitempty"`
}

// CreateResourceEndpointResponse consists of id, name, status, selfLink and the progress.
type CreateResourceEndpointResponse struct {
	Progress int8   `json:"progress"`
	Status   string `json:"status"`
	Name     string `json:"name"`
	ID       string `json:"id"`
	SelfLink string `json:"selfLink"`
}

// RequestTrackerResponse is the response received for the on-going tracker requests.
type RequestTrackerResponse struct {
	Progress  int8     `json:"progress"`
	Message   string   `json:"message"`
	Status    string   `json:"status"`
	Resource  string   `json:"resource"`
	Resources []string `json:"resources"`
	Name      string   `json:"name"`
	ID        string   `json:"id"`
	SelfLink  string   `json:"selfLink"`
}

// HypertextReference is used to represent the link(s) that one resource might have to other resources.
type HypertextReference struct {
	Href  string   `json:"href,omitempty"`
	Hrefs []string `json:"hrefs,omitempty"`
}

// Machine is used to represent a CAS provisioned machine state.
type Machine struct {
	PowerState       string                        `json:"powerState"`
	Address          string                        `json:"address"`
	CustomProperties map[string]string             `json:"customProperties"`
	ProjectID        string                        `json:"projectId"`
	ExternalZoneID   string                        `json:"externalZoneId"`
	ExternalRegionID string                        `json:"externalRegionId"`
	ExternalID       string                        `json:"externalId"`
	Name             string                        `json:"name"`
	Description      string                        `json:"description"`
	ID               string                        `json:"id"`
	SelfLink         string                        `json:"selfLink"`
	CreatedAt        string                        `json:"createdAt"`
	UpdatedAt        string                        `json:"updatedAt"`
	Owner            string                        `json:"owner"`
	OrganizationID   string                        `json:"organizationId"`
	Links            map[string]HypertextReference `json:"_links"`
	Tags             []Tag                         `json:"tags"`
}

// Network is used to represent a CAS provisioned network resource state.
type Network struct {
	CIDR             string                        `json:"cidr"`
	CustomProperties map[string]string             `json:"customProperties"`
	ProjectID        string                        `json:"projectId"`
	ExternalZoneID   string                        `json:"externalZoneId"`
	ExternalID       string                        `json:"externalId"`
	Name             string                        `json:"name"`
	Description      string                        `json:"description"`
	ID               string                        `json:"id"`
	SelfLink         string                        `json:"selfLink"`
	UpdatedAt        string                        `json:"updatedAt"`
	Owner            string                        `json:"owner"`
	OrganizationID   string                        `json:"organizationId"`
	Links            map[string]HypertextReference `json:"_links"`
	Tags             []Tag                         `json:"tags"`
}

// BlockDevice is used to represent a CAS provisioned block device resource state.
type BlockDevice struct {
	CapacityInGB     int                           `json:"capacityInGB"`
	Status           string                        `json:"status"`
	CustomProperties map[string]string             `json:"customProperties"`
	ProjectID        string                        `json:"projectId"`
	ExternalZoneID   string                        `json:"externalZoneId"`
	ExternalRegionID string                        `json:"externalRegionId"`
	ExternalID       string                        `json:"externalId"`
	Name             string                        `json:"name"`
	Description      string                        `json:"description"`
	ID               string                        `json:"id"`
	SelfLink         string                        `json:"selfLink"`
	CreatedAt        string                        `json:"createdAt"`
	UpdatedAt        string                        `json:"updatedAt"`
	Owner            string                        `json:"owner"`
	OrganizationID   string                        `json:"organizationId"`
	Links            map[string]HypertextReference `json:"_links"`
	Tags             []Tag                         `json:"tags"`
}

// LoadBalancer is used to represent a CAS provisioned load balancer resource state.
type LoadBalancer struct {
	ID               string                        `json:"id"`
	SelfLink         string                        `json:"selfLink"`
	CreatedAt        string                        `json:"createdAt"`
	UpdatedAt        string                        `json:"updatedAt"`
	Owner            string                        `json:"owner"`
	OrganizationID   string                        `json:"organizationId"`
	Links            map[string]HypertextReference `json:"_links"`
	Name             string                        `json:"name"`
	Description      string                        `json:"description"`
	ExternalID       string                        `json:"externalId"`
	ProjectID        string                        `json:"projectId"`
	ExternalZoneID   string                        `json:"externalZoneId"`
	ExternalRegionID string                        `json:"externalRegionId"`
	CustomProperties map[string]string             `json:"customProperties"`
	Tags             []Tag                         `json:"tags"`
	Routes           []Route                       `json:"routes"`
	Address          string                        `json:"address"`
}

// MachineAttachedDisks represent the disks attached to a machine resource.
type MachineAttachedDisks struct {
	Content       []BlockDevice `json:"content"`
	TotalElements int           `json:"totalElements"`
}

// NewResourceObject returns the new resource object for a given resource endpoint.
// Works only for supported resource endpoints.
func NewResourceObject(endpoint string) (interface{}, error) {
	switch {
	case machineEndpointDisksR.MatchString(endpoint):
		return &MachineAttachedDisks{}, nil
	case machineEndpointR.MatchString(endpoint):
		return &Machine{}, nil
	case networkEndpointR.MatchString(endpoint):
		return &Network{}, nil
	case blockDeviceEndpointR.MatchString(endpoint):
		return &BlockDevice{}, nil
	case loadBalancerEndpointR.MatchString(endpoint):
		return &LoadBalancer{}, nil
	default:
		return nil, errors.New("endpoint: " + endpoint + " not supported")
	}
}

// ResourceSpecification tracks the endpoint for CAS resources.
type ResourceSpecification interface {
	GetEndpoint() string
}

// GetEndpoint returns the machine endpoint.
func (machineSpecification MachineSpecification) GetEndpoint() string {
	return machineEndpoint
}

// GetEndpoint returns network endpoint.
func (networkSpecification NetworkSpecification) GetEndpoint() string {
	return networkEndpoint
}

// GetEndpoint returns the block device endpoint.
func (blockDeviceSpecification BlockDeviceSpecification) GetEndpoint() string {
	return blockDeviceEndpoint
}

// GetEndpoint returns the load balancer endpoint.
func (loadBalancerSpecification LoadBalancerSpecification) GetEndpoint() string {
	return loadBalancerEndpoint
}
