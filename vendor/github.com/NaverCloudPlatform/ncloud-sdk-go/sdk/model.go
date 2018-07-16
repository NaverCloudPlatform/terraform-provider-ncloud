package sdk

import (
	common "github.com/NaverCloudPlatform/ncloud-sdk-go/common"
)

type Conn struct {
	accessKey string
	secretKey string
	apiURL    string
}

// ServerImage structures
type ServerImage struct {
	MemberServerImageNo                    string            `xml:"memberServerImageNo"`
	MemberServerImageName                  string            `xml:"memberServerImageName"`
	MemberServerImageDescription           string            `xml:"memberServerImageDescription"`
	OriginalServerInstanceNo               string            `xml:"originalServerInstanceNo"`
	OriginalServerProductCode              string            `xml:"originalServerProductCode"`
	OriginalServerName                     string            `xml:"originalServerName"`
	OriginalBaseBlockStorageDiskType       common.CommonCode `xml:"originalBaseBlockStorageDiskType"`
	OriginalServerImageProductCode         string            `xml:"originalServerImageProductCode"`
	OriginalOsInformation                  string            `xml:"originalOsInformation"`
	OriginalServerImageName                string            `xml:"originalServerImageName"`
	MemberServerImageStatusName            string            `xml:"memberServerImageStatusName"`
	MemberServerImageStatus                common.CommonCode `xml:"memberServerImageStatus"`
	MemberServerImageOperation             common.CommonCode `xml:"memberServerImageOperation"`
	MemberServerImagePlatformType          common.CommonCode `xml:"memberServerImagePlatformType"`
	CreateDate                             string            `xml:"createDate"`
	Zone                                   common.Zone       `xml:"zone"`
	Region                                 common.Region     `xml:"region"`
	MemberServerImageBlockStorageTotalRows int               `xml:"memberServerImageBlockStorageTotalRows"`
	MemberServerImageBlockStorageTotalSize int               `xml:"memberServerImageBlockStorageTotalSize"`
}

type MemberServerImageList struct {
	common.CommonResponse
	TotalRows             int           `xml:"totalRows"`
	MemberServerImageList []ServerImage `xml:"memberServerImageList>memberServerImage,omitempty"`
}

type RequestServerImageList struct {
	MemberServerImageNoList []string
	PlatformTypeCodeList    []string
	PageNo                  int
	PageSize                int
	RegionNo                string
	SortedBy                string
	SortingOrder            string
}

type RequestCreateServerImage struct {
	MemberServerImageName        string
	MemberServerImageDescription string
	ServerInstanceNo             string
}

type RequestGetServerImageProductList struct {
	ExclusionProductCode        string
	ProductCode                 string
	PlatformTypeCodeList        []string
	InfraResourceDetailTypeCode string
	BlockStorageSize            int
	RegionNo                    string
}

// ProductList : Response of server product list
type ProductList struct {
	common.CommonResponse
	TotalRows int       `xml:"totalRows"`
	Product   []Product `xml:"productList>product,omitempty"`
}

// Product : Product information of Server
type Product struct {
	ProductCode          string            `xml:"productCode"`
	ProductName          string            `xml:"productName"`
	ProductType          common.CommonCode `xml:"productType"`
	ProductDescription   string            `xml:"productDescription"`
	InfraResourceType    common.CommonCode `xml:"infraResourceType"`
	CPUCount             int               `xml:"cpuCount"`
	MemorySize           int               `xml:"memorySize"`
	BaseBlockStorageSize int               `xml:"baseBlockStorageSize"`
	PlatformType         common.CommonCode `xml:"platformType"`
	OsInformation        string            `xml:"osInformation"`
	AddBlockStroageSize  int               `xml:"addBlockStroageSize"`
}

// RequestCreateServerInstance is Server Instances structures
type RequestCreateServerInstance struct {
	ServerImageProductCode                string
	ServerProductCode                     string
	MemberServerImageNo                   string
	ServerName                            string
	ServerDescription                     string
	LoginKeyName                          string
	IsProtectServerTermination            string
	ServerCreateCount                     int
	ServerCreateStartNo                   int
	InternetLineTypeCode                  string
	FeeSystemTypeCode                     string
	UserData                              string
	ZoneNo                                string
	AccessControlGroupConfigurationNoList []string
	RaidTypeName                          string
}

type ServerInstanceList struct {
	common.CommonResponse
	TotalRows          int              `xml:"totalRows"`
	ServerInstanceList []ServerInstance `xml:"serverInstanceList>serverInstance,omitempty"`
}

type ServerInstance struct {
	ServerInstanceNo               string               `xml:"serverInstanceNo"`
	ServerName                     string               `xml:"serverName"`
	ServerDescription              string               `xml:"serverDescription"`
	CPUCount                       int                  `xml:"cpuCount"`
	MemorySize                     int                  `xml:"memorySize"`
	BaseBlockStorageSize           int                  `xml:"baseBlockStorageSize"`
	PlatformType                   common.CommonCode    `xml:"platformType"`
	LoginKeyName                   string               `xml:"loginKeyName"`
	IsFeeChargingMonitoring        bool                 `xml:"isFeeChargingMonitoring"`
	PublicIP                       string               `xml:"publicIp"`
	PrivateIP                      string               `xml:"privateIp"`
	ServerImageName                string               `xml:"serverImageName"`
	ServerInstanceStatus           common.CommonCode    `xml:"serverInstanceStatus"`
	ServerInstanceOperation        common.CommonCode    `xml:"serverInstanceOperation"`
	ServerInstanceStatusName       string               `xml:"serverInstanceStatusName"`
	CreateDate                     string               `xml:"createDate"`
	Uptime                         string               `xml:"uptime"`
	ServerImageProductCode         string               `xml:"serverImageProductCode"`
	ServerProductCode              string               `xml:"serverProductCode"`
	IsProtectServerTermination     bool                 `xml:"isProtectServerTermination"`
	PortForwardingPublicIP         string               `xml:"portForwardingPublicIp"`
	PortForwardingExternalPort     int                  `xml:"portForwardingExternalPort"`
	PortForwardingInternalPort     int                  `xml:"portForwardingInternalPort"`
	Zone                           common.Zone          `xml:"zone"`
	Region                         common.Region        `xml:"region"`
	BaseBlockStorageDiskType       common.CommonCode    `xml:"baseBlockStorageDiskType"`
	BaseBlockStroageDiskDetailType common.CommonCode    `xml:"baseBlockStroageDiskDetailType"`
	InternetLineType               common.CommonCode    `xml:"internetLineType"`
	ServerInstanceType             common.CommonCode    `xml:"serverInstanceType"`
	UserData                       string               `xml:"userData"`
	AccessControlGroupList         []AccessControlGroup `xml:"accessControlGroupList>accessControlGroup"`
}

type AccessControlGroup struct {
	AccessControlGroupConfigurationNo string `xml:"accessControlGroupConfigurationNo"`
	AccessControlGroupName            string `xml:"accessControlGroupName"`
	AccessControlGroupDescription     string `xml:"accessControlGroupDescription"`
	IsDefault                         bool   `xml:"isDefault"`
	CreateDate                        string `xml:"createDate"`
}

// RequestGetLoginKeyList is Login Key structures
type RequestGetLoginKeyList struct {
	KeyName  string
	PageNo   int
	PageSize int
}

type LoginKeyList struct {
	common.CommonResponse
	TotalRows    int        `xml:"totalRows"`
	LoginKeyList []LoginKey `xml:"loginKeyList>loginKey,omitempty"`
}

type LoginKey struct {
	Fingerprint string `xml:"fingerprint"`
	KeyName     string `xml:"keyName"`
	CreateDate  string `xml:"createDate"`
}

type PrivateKey struct {
	common.CommonResponse
	PrivateKey string `xml:"privateKey"`
}
type RequestCreatePublicIPInstance struct {
	ServerInstanceNo     string
	PublicIPDescription  string
	InternetLineTypeCode string
	RegionNo             string
	ZoneNo               string
}

type RequestAssociatePublicIP struct {
	ServerInstanceNo   string
	PublicIPInstanceNo string
}

type RequestPublicIPInstanceList struct {
	IsAssociated           string
	PublicIPInstanceNoList []string
	PublicIPList           []string
	SearchFilterName       string
	SearchFilterValue      string
	InternetLineTypeCode   string
	RegionNo               string
	ZoneNo                 string
	PageNo                 int
	PageSize               int
	SortedBy               string
	SortingOrder           string
}

type PublicIPInstanceList struct {
	common.CommonResponse
	TotalRows            int                `xml:"totalRows"`
	PublicIPInstanceList []PublicIPInstance `xml:"publicIpInstanceList>publicIpInstance,omitempty"`
}

type PublicIPInstance struct {
	PublicIPInstanceNo         string            `xml:"publicIpInstanceNo"`
	PublicIP                   string            `xml:"publicIp"`
	PublicIPDescription        string            `xml:"publicIpDescription"`
	CreateDate                 string            `xml:"createDate"`
	InternetLineType           common.CommonCode `xml:"internetLineType"`
	PublicIPInstanceStatusName string            `xml:"publicIpInstanceStatusName"`
	PublicIPInstanceStatus     common.CommonCode `xml:"publicIpInstanceStatus"`
	PublicIPInstanceOperation  common.CommonCode `xml:"publicIpInstanceOperation"`
	PublicIPKindType           common.CommonCode `xml:"publicIpKindType"`
	ServerInstance             ServerInstance    `xml:"serverInstanceAssociatedWithPublicIp"`
	Zone                       common.Zone       `xml:"zone"`
}

type RequestDeletePublicIPInstances struct {
	PublicIPInstanceNoList []string
}

// RequestGetServerInstanceList : Get Server Instance List
type RequestGetServerInstanceList struct {
	ServerInstanceNoList               []string
	SearchFilterName                   string
	SearchFilterValue                  string
	PageNo                             int
	PageSize                           int
	ServerInstanceStatusCode           string
	InternetLineTypeCode               string
	RegionNo                           string
	ZoneNo                             string
	BaseBlockStorageDiskTypeCode       string
	BaseBlockStorageDiskDetailTypeCode string
	SortedBy                           string
	SortingOrder                       string
	ServerInstanceTypeCodeList         []string
}

type RequestStopServerInstances struct {
	ServerInstanceNoList []string
}

type RequestStartServerInstances struct {
	ServerInstanceNoList []string
}

type RequestTerminateServerInstances struct {
	ServerInstanceNoList []string
}

type RequestRebootServerInstances struct {
	ServerInstanceNoList []string
}

type RequestChangeServerInstanceSpec struct {
	ServerInstanceNo  string
	ServerProductCode string
}

// RequestGetRootPassword : Request to get root password of the server
type RequestGetRootPassword struct {
	ServerInstanceNo string
	PrivateKey       string
}

// RootPassword : Response of getting root password of the server
type RootPassword struct {
	common.CommonResponse
	TotalRows    int    `xml:"totalRows"`
	RootPassword string `xml:"rootPassword"`
}

// RequestGetZoneList : Request to get zone list
type RequestGetZoneList struct {
	regionNo string
}

// ZoneList : Response of getting zone list
type ZoneList struct {
	common.CommonResponse
	TotalRows int           `xml:"totalRows"`
	Zone      []common.Zone `xml:"zoneList>zone"`
}

// RegionList : Response of getting region list
type RegionList struct {
	common.CommonResponse
	TotalRows  int             `xml:"totalRows"`
	RegionList []common.Region `xml:"regionList>region,omitempty"`
}

type RequestBlockStorageInstance struct {
	BlockStorageName        string
	BlockStorageSize        int
	BlockStorageDescription string
	ServerInstanceNo        string
	DiskDetailTypeCode      string
}

type RequestBlockStorageInstanceList struct {
	ServerInstanceNo               string
	BlockStorageInstanceNoList     []string
	SearchFilterName               string
	SearchFilterValue              string
	BlockStorageTypeCodeList       []string
	PageNo                         int
	PageSize                       int
	BlockStorageInstanceStatusCode string
	DiskTypeCode                   string
	DiskDetailTypeCode             string
	RegionNo                       string
	ZoneNo                         string
	SortedBy                       string
	SortingOrder                   string
}

type BlockStorageInstanceList struct {
	common.CommonResponse
	TotalRows            int                    `xml:"totalRows"`
	BlockStorageInstance []BlockStorageInstance `xml:"blockStorageInstanceList>blockStorageInstance,omitempty"`
}

type BlockStorageInstance struct {
	BlockStorageInstanceNo          string            `xml:"blockStorageInstanceNo"`
	ServerInstanceNo                string            `xml:"serverInstanceNo"`
	ServerName                      string            `xml:"serverName"`
	BlockStorageType                common.CommonCode `xml:"blockStorageType"`
	BlockStorageName                string            `xml:"blockStorageName"`
	BlockStorageSize                int               `xml:"blockStorageSize"`
	DeviceName                      string            `xml:"deviceName"`
	BlockStorageProductCode         string            `xml:"blockStorageProductCode"`
	BlockStorageInstanceStatus      common.CommonCode `xml:"blockStorageInstanceStatus"`
	BlockStorageInstanceOperation   common.CommonCode `xml:"blockStorageInstanceOperation"`
	BlockStorageInstanceStatusName  string            `xml:"blockStorageInstanceStatusName"`
	CreateDate                      string            `xml:"createDate"`
	BlockStorageInstanceDescription string            `xml:"blockStorageInstanceDescription"`
	DiskType                        common.CommonCode `xml:"diskType"`
	DiskDetailType                  common.CommonCode `xml:"diskDetailType"`
	Zone                            common.Zone       `xml:"zone"`
}

// RequestAttachBlockStorageInstance is request type to attach server instance
type RequestAttachBlockStorageInstance struct {
	ServerInstanceNo       string
	BlockStorageInstanceNo string
}

// RequestDetachBlockStorageInstance is request type to detach block storage instance from server instance
type RequestDetachBlockStorageInstance struct {
	BlockStorageInstanceNoList []string
}

// RequestGetServerProductList : Request to get server product list
type RequestGetServerProductList struct {
	ExclusionProductCode   string
	ProductCode            string
	ServerImageProductCode string
	ZoneNo                 string
	InternetLineTypeCode   string
	RegionNo               string
}

type RequestAccessControlGroupList struct {
	AccessControlGroupConfigurationNoList []string
	IsDefault                             string
	AccessControlGroupName                string
	PageNo                                int
	PageSize                              int
}

type AccessControlGroupList struct {
	common.CommonResponse
	TotalRows          int                  `xml:"totalRows"`
	AccessControlGroup []AccessControlGroup `xml:"accessControlGroupList>accessControlGroup,omitempty"`
}

type AccessControlRuleList struct {
	common.CommonResponse
	TotalRows             int                 `xml:"totalRows"`
	AccessControlRuleList []AccessControlRule `xml:"accessControlRuleList>accessControlRule,omitempty"`
}

type AccessControlRule struct {
	AccessControlRuleConfigurationNo       string            `xml:"accessControlRuleConfigurationNo"`
	AccessControlRuleDescription           string            `xml:"accessControlRuleDescription"`
	SourceAccessControlRuleConfigurationNo string            `xml:"sourceAccessControlRuleConfigurationNo"`
	SourceAccessControlRuleName            string            `xml:"sourceAccessControlRuleName"`
	ProtocolType                           common.CommonCode `xml:"protocolType"`
	SourceIP                               string            `xml:"sourceIp"`
	DestinationPort                        string            `xml:"destinationPort"`
}

type RequestCreateNasVolumeInstance struct {
	VolumeName                      string
	VolumeSize                      int
	VolumeAllotmentProtocolTypeCode string
	ServerInstanceNoList            []string
	CustomIpList                    []string
	CifsUserName                    string
	CifsUserPassword                string
	NasVolumeDescription            string
	RegionNo                        string
	ZoneNo                          string
}

type NasVolumeInstance struct {
	NasVolumeInstanceNo              string                      `xml:"nasVolumeInstanceNo"`
	NasVolumeInstanceStatus          common.CommonCode           `xml:"nasVolumeInstanceStatus"`
	NasVolumeInstanceOperation       common.CommonCode           `xml:"nasVolumeInstanceOperation"`
	NasVolumeInstanceStatusName      string                      `xml:"nasVolumeInstanceStatusName"`
	CreateDate                       string                      `xml:"createDate"`
	NasVolumeInstanceDescription     string                      `xml:"nasVolumeInstanceDescription"`
	MountInformation                 string                      `xml:"mountInformation"`
	VolumeAllotmentProtocolType      common.CommonCode           `xml:"volumeAllotmentProtocolType"`
	VolumeName                       string                      `xml:"volumeName"`
	VolumeTotalSize                  int                         `xml:"volumeTotalSize"`
	VolumeSize                       int                         `xml:"volumeSize"`
	VolumeUseSize                    int                         `xml:"volumeUseSize"`
	VolumeUseRatio                   float32                     `xml:"volumeUseRatio"`
	SnapshotVolumeConfigurationRatio float32                     `xml:"snapshotVolumeConfigurationRatio"`
	SnapshotVolumeSize               int                         `xml:"snapshotVolumeSize"`
	SnapshotVolumeUseSize            int                         `xml:"snapshotVolumeUseSize"`
	SnapshotVolumeUseRatio           float32                     `xml:"snapshotVolumeUseRatio"`
	IsSnapshotConfiguration          bool                        `xml:"isSnapshotConfiguration"`
	IsEventConfiguration             bool                        `xml:"isEventConfiguration"`
	Zone                             common.Zone                 `xml:"zone"`
	Region                           common.Region               `xml:"region"`
	NasVolumeInstanceCustomIPList    []NasVolumeInstanceCustomIp `xml:"nasVolumeInstanceCustomIpList>nasVolumeInstanceCustomIp,omitempty"`
	NasVolumeServerInstanceList      []ServerInstance            `xml:"nasVolumeServerInstanceList>serverInstance,omitempty"`
}

type NasVolumeInstanceList struct {
	common.CommonResponse
	TotalRows             int                 `xml:"totalRows"`
	NasVolumeInstanceList []NasVolumeInstance `xml:"nasVolumeInstanceList>nasVolumeInstance,omitempty"`
}

type NasVolumeInstanceCustomIp struct {
	CustomIP string `xml:"customIp"`
}

type RequestGetNasVolumeInstanceList struct {
	VolumeAllotmentProtocolTypeCode string
	IsEventConfiguration            string
	IsSnapshotConfiguration         string
	NasVolumeInstanceNoList         []string
	RegionNo                        string
	ZoneNo                          string
}

type PortForwardingRule struct {
	ServerInstanceNo           string `xml:"serverInstance>serverInstanceNo"`
	PortForwardingExternalPort string `xml:"portForwardingExternalPort"`
	PortForwardingInternalPort string `xml:"portForwardingInternalPort"`
	PortForwardingPublicIp     string `xml:"serverInstance>portForwardingPublicIp"`
}

type RequestAddPortForwardingRules struct {
	PortForwardingConfigurationNo string
	PortForwardingRuleList        []PortForwardingRule
}

type RequestDeletePortForwardingRules struct {
	PortForwardingConfigurationNo string
	PortForwardingRuleList        []PortForwardingRule
}

type PortForwardingRuleList struct {
	common.CommonResponse
	PortForwardingConfigurationNo int                  `xml:"portForwardingConfigurationNo"`
	PortForwardingPublicIp        string               `xml:"portForwardingPublicIp"`
	Zone                          common.Zone          `xml:"zone"`
	TotalRows                     int                  `xml:"totalRows"`
	PortForwardingRuleList        []PortForwardingRule `xml:"portForwardingRuleList>portForwardingRule,omitempty"`
}

type RequestPortForwardingRuleList struct {
	InternetLineTypeCode string
	RegionNo             string
	ZoneNo               string
}

// RequestLoadBalancerInstanceList is request type to get load balancer instance list
type RequestLoadBalancerInstanceList struct {
	LoadBalancerInstanceNoList []string
	InternetLineTypeCode       string
	NetworkUsageTypeCode       string
	RegionNo                   string
	PageNo                     int
	PageSize                   int
	SortedBy                   string
	SortingOrder               string
}

// LoadBalancerInstanceList is response type to return load balancer instance list
type LoadBalancerInstanceList struct {
	common.CommonResponse
	LoadBalancerInstanceList []LoadBalancerInstance `xml:"loadBalancerInstanceList>loadBalancerInstance,omitempty"`
	TotalRows                int                    `xml:"totalRows"`
}

// LoadBalancerInstance is struct for load balancer instance
type LoadBalancerInstance struct {
	LoadBalancerInstanceNo         string                       `xml:"loadBalancerInstanceNo"`
	VirtualIP                      string                       `xml:"virtualIp"`
	LoadBalancerName               string                       `xml:"loadBalancerName"`
	LoadBalancerAlgorithmType      common.CommonCode            `xml:"loadBalancerAlgorithmType"`
	LoadBalancerDescription        string                       `xml:"loadBalancerDescription"`
	CreateDate                     string                       `xml:"createDate"`
	DomainName                     string                       `xml:"domainName"`
	InternetLineType               common.CommonCode            `xml:"internetLineType"`
	LoadBalancerInstanceStatusName string                       `xml:"loadBalancerInstanceStatusName"`
	LoadBalancerInstanceStatus     common.CommonCode            `xml:"loadBalancerInstanceStatus"`
	LoadBalancerInstanceOperation  common.CommonCode            `xml:"loadBalancerInstanceOperation"`
	NetworkUsageType               common.CommonCode            `xml:"networkUsageType"`
	IsHTTPKeepAlive                bool                         `xml:"isHttpKeepAlive"`
	ConnectionTimeout              int                          `xml:"connectionTimeout"`
	CertificateName                string                       `xml:"certificateName"`
	LoadBalancerRuleList           []LoadBalancerRule           `xml:"loadBalancerRuleList>loadBalancerRule,omitempty"`
	LoadBalancedServerInstanceList []LoadBalancedServerInstance `xml:"loadBalancedServerInstanceList>loadBalancedServerInstance,omitempty"`
}

// LoadBalancerRule is struct for load balancer rule
type LoadBalancerRule struct {
	ProtocolType       common.CommonCode `xml:"protocolType"`
	LoadBalancerPort   int               `xml:"loadBalancerPort"`
	ServerPort         int               `xml:"serverPort"`
	L7HealthCheckPath  string            `xml:"l7HealthCheckPath"`
	CertificateName    string            `xml:"certificateName"`
	ProxyProtocolUseYn string            `xml:"proxyProtocolUseYn"`
}

// LoadBalancedServerInstance is struct for load balanced server instance
type LoadBalancedServerInstance struct {
	ServerInstanceList          []ServerInstance          `xml:"serverInstance,omitempty"`
	ServerHealthCheckStatusList []ServerHealthCheckStatus `xml:"serverHealthCheckStatusList>serverHealthCheckStatus,omitempty"`
}

// ServerHealthCheckStatus is struct for server health check status
type ServerHealthCheckStatus struct {
	ProtocolType       common.CommonCode `xml:"protocolType"`
	LoadBalancerPort   int               `xml:"loadBalancerPort"`
	ServerPort         int               `xml:"serverPort"`
	L7HealthCheckPath  string            `xml:"l7HealthCheckPath"`
	ProxyProtocolUseYn string            `xml:"proxyProtocolUseYn"`
	ServerStatus       bool              `xml:"serverStatus"`
}

// RequestCreateLoadBalancerInstance is request type to create load balancer instance
type RequestCreateLoadBalancerInstance struct {
	LoadBalancerName              string
	LoadBalancerAlgorithmTypeCode string
	LoadBalancerDescription       string
	LoadBalancerRuleList          []RequestLoadBalancerRule
	ServerInstanceNoList          []string
	InternetLineTypeCode          string
	NetworkUsageTypeCode          string
	RegionNo                      string
}

// RequestLoadBalancerRule is request type to create load balancer rule
type RequestLoadBalancerRule struct {
	ProtocolTypeCode   string
	LoadBalancerPort   int
	ServerPort         int
	L7HealthCheckPath  string
	CertificateName    string
	ProxyProtocolUseYn string
}

// RequestDeleteLoadBalancerInstances is request type to delete load balancer instances
type RequestDeleteLoadBalancerInstances struct {
	LoadBalancerInstanceNoList []string
}

// RequestChangeLoadBalancerInstanceConfiguration is request type to change load balancer instance configuration
type RequestChangeLoadBalancerInstanceConfiguration struct {
	LoadBalancerInstanceNo        string
	LoadBalancerAlgorithmTypeCode string
	LoadBalancerDescription       string
	LoadBalancerRuleList          []RequestLoadBalancerRule
}

// RequestChangeLoadBalancedServerInstances is request type to change load balanced server instances
type RequestChangeLoadBalancedServerInstances struct {
	LoadBalancerInstanceNo string
	ServerInstanceNoList   []string
}

// RequestGetLoadBalancerTargetServerInstanceList is request type to get load balancer target server instance list
type RequestGetLoadBalancerTargetServerInstanceList struct {
	InternetLineTypeCode string
	NetworkUsageTypeCode string
	RegionNo             string
}

// SslCertificateList is response type to return SSL Certificate list
type SslCertificateList struct {
	common.CommonResponse
	SslCertificateList []SslCertificate `xml:"sslCertificateList>sslCertificate,omitempty"`
	TotalRows          int              `xml:"totalRows"`
}

// SslCertificate is struct for SSL Certificate
type SslCertificate struct {
	CertificateName      string `xml:"certificateName"`
	PrivateKey           string `xml:"privateKey"`
	PublicKeyCertificate string `xml:"publicKeyCertificate"`
	CertificateChain     string `xml:"certificateChain"`
}

// RequestAddSslCertificate is request type to add SSL Certificate
type RequestAddSslCertificate struct {
	CertificateName      string
	PrivateKey           string
	PublicKeyCertificate string
	CertificateChain     string
}

type RequestChangeNasVolumeSize struct {
	NasVolumeInstanceNo string
	VolumeSize          int
}

// RequestNasVolumeAccessControl is request type for nas volume access control operations
type RequestNasVolumeAccessControl struct {
	NasVolumeInstanceNo  string
	ServerInstanceNoList []string
	CustomIPList         []string
}

type RequestRecreateServerInstance struct {
	ServerInstanceNo             string
	ServerInstanceName           string
	ChangeServerImageProductCode string
}

type RequestCreateBlockStorageSnapshotInstance struct {
	BlockStorageInstanceNo          string
	BlockStorageSnapshotName        string
	BlockStorageSnapshotDescription string
}

type BlockStorageSnapshotInstance struct {
	BlockStorageSnapshotInstanceNo          string            `xml:"blockStorageSnapshotInstanceNo"`
	BlockStorageSnapshotName                string            `xml:"blockStorageSnapshotName"`
	BlockStorageSnapshotVolumeSize          int               `xml:"blockStorageSnapshotVolumeSize"`
	OriginalBlockStorageInstanceNo          string            `xml:"originalBlockStorageInstanceNo"`
	OriginalBlockStorageName                string            `xml:"originalBlockStorageName"`
	BlockStorageSnapshotInstanceStatus      common.CommonCode `xml:"blockStorageSnapshotInstanceStatus"`
	BlockStorageSnapshotInstanceOperation   common.CommonCode `xml:"blockStorageSnapshotInstanceOperation"`
	BlockStorageSnapshotInstanceStatusName  string            `xml:"blockStorageSnapshotInstanceStatusName"`
	CreateDate                              string            `xml:"createDate"`
	BlockStorageSnapshotInstanceDescription string            `xml:"blockStorageSnapshotInstanceDescription"`
	ServerImageProductCode                  string            `xml:"serverImageProductCode"`
	OsInformation                           string            `xml:"osInformation"`
}

type BlockStorageSnapshotInstanceList struct {
	common.CommonResponse
	TotalRows                        int                            `xml:"totalRows"`
	BlockStorageSnapshotInstanceList []BlockStorageSnapshotInstance `xml:"blockStorageSnapshotInstanceList>blockStorageSnapshot,omitempty"`
}

type RequestGetBlockStorageSnapshotInstanceList struct {
	BlockStorageSnapshotInstanceNoList []string
	OriginalBlockStorageInstanceNoList []string
	RegionNo                           string
	PageNo                             int
	PageSize                           int
}

// RequestGetLaunchConfigurationList is request type for Launch Configuration List
type RequestGetLaunchConfigurationList struct {
	LaunchConfigurationNameList []string
	PageNo                      int
	PageSize                    int
	SortedBy                    string
	SortingOrder                string
}

type LaunchConfigurationList struct {
	common.CommonResponse
	TotalRows               int                   `xml:"totalRows"`
	LaunchConfigurationList []LaunchConfiguration `xml:"launchConfigurationList>launchConfiguration,omitempty"`
}

type LaunchConfiguration struct {
	LaunchConfigurationName string               `xml:"launchConfigurationName"`
	ServerImageProductCode  string               `xml:"serverImageProductCode"`
	ServerProductCode       string               `xml:"serverProductCode"`
	MemberServerImageNo     string               `xml:"memberServerImageNo"`
	LoginKeyName            string               `xml:"loginKeyName"`
	CreateDate              string               `xml:"createDate"`
	UserData                string               `xml:"userData"`
	AccessControlGroupList  []AccessControlGroup `xml:"accessControlGroupList>accessControlGroup,omitempty"`
}

type RequestCreateLaunchConfiguration struct {
	LaunchConfigurationName               string
	ServerImageProductCode                string
	ServerProductCode                     string
	MemberServerImageNo                   string
	AccessControlGroupConfigurationNoList []string
	LoginKeyName                          string
	UserData                              string
}
