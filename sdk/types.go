package tango

import (
	"errors"
	"regexp"
)

const (
	machineEndpoint      = "/iaas/machines"
	blockDeviceEndpoint  = "/iaas/block-devices"
	networkEndpoint      = "/iaas/networks"
	loadBalancerEndpoint = "/iaas/load-balancers"
)

var (
	machineEndpointDisksR = regexp.MustCompile("^/iaas/machines/[a-zA-Z0-9]+/disks")
	machineEndpointR      = regexp.MustCompile("^/iaas/machines")
	networkEndpointR      = regexp.MustCompile("^/iaas/networks")
	blockDeviceEndpointR  = regexp.MustCompile("^/iaas/block-devices")
	loadBalancerEndpointR = regexp.MustCompile("^/iaas/load-balancers")
)

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Constraint struct {
	Mandatory  bool   `json:"mandatory"`
	Expression string `json:"expression"`
}

type Nic struct {
	Name             string            `json:"name,omitempty"`
	Description      string            `json:"description,omitempty"`
	DeviceIndex      int               `json:"deviceIndex,omitempty"`
	NetworkID        string            `json:"networkId"`
	Addresses        []string          `json:"addresses,omitempty"`
	SecurityGroupIDs []string          `json:"securityGroupIds,omitempty"`
	CustomProperties map[string]string `json:"customProperties,omitempty"`
}

type Disk struct {
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	BlockDeviceID string `json:"blockDeviceId"`
}

type Route struct {
	Protocol       string                    `json:"protocol"`
	Port           string                    `json:"port"`
	MemberProtocol string                    `json:"memberProtocol"`
	MemberPort     string                    `json:"memberPort"`
	HCC            *HealthCheckConfiguration `json:"healthCheckConfiguration,omitempty"`
}

type HealthCheckConfiguration struct {
	Protocol          string `json:"protocol"`
	Port              string `json:"port"`
	URLPath           string `json:"urlPath,omitempty"`
	IntervalSeconds   int    `json:"intervalSeconds,omitempty"`
	TimeoutSeconds    int    `json:"timeoutSeconds,omitempty"`
	UnhealthThreshold int    `json:"unhealthyThreshold,omitempty"`
	HealthThreshold   int    `json:"healthyThreshold,omitempty"`
}

type LoginRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type LoginResponse struct {
	TokenType string `json:"tokenType"`
	Token     string `json:"token"`
}

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

type CreateResourceEndpointResponse struct {
	Progress int8   `json:"progress"`
	Status   string `json:"status"`
	Name     string `json:"name"`
	ID       string `json:"id"`
	SelfLink string `json:"selfLink"`
}

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

type HypertextReference struct {
	Href  string   `json:"href,omitempty"`
	Hrefs []string `json:"hrefs,omitempty"`
}

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

type MachineAttachedDisks struct {
	Content       []BlockDevice `json:"content"`
	TotalElements int           `json:"totalElements"`
}

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

type ResourceSpecification interface {
	GetEndpoint() string
}

func (machineSpecification MachineSpecification) GetEndpoint() string {
	return machineEndpoint
}

func (networkSpecification NetworkSpecification) GetEndpoint() string {
	return networkEndpoint
}

func (blockDeviceSpecification BlockDeviceSpecification) GetEndpoint() string {
	return blockDeviceEndpoint
}

func (loadBalancerSpecification LoadBalancerSpecification) GetEndpoint() string {
	return loadBalancerEndpoint
}
