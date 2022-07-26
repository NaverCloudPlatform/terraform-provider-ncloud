package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vautoscaling"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vcdss"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vloadbalancer"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnks"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vses2"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnas"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/autoscaling"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/cdn"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/clouddb"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/loadbalancer"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/monitoring"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/sourcebuild"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/sourcecommit"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/sourcepipeline"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vsourcedeploy"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vsourcepipeline"
)

// DefaultWaitForInterval is Interval for checking status in WaitForXXX method
const DefaultWaitForInterval = 10

// Default timeout
const DefaultTimeout = 5 * time.Minute
const DefaultCreateTimeout = 1 * time.Hour
const DefaultUpdateTimeout = 10 * time.Minute
const DefaultStopTimeout = 5 * time.Minute

type Config struct {
	AccessKey string
	SecretKey string
	Region    string
}

type NcloudAPIClient struct {
	server          *server.APIClient
	autoscaling     *autoscaling.APIClient
	loadbalancer    *loadbalancer.APIClient
	cdn             *cdn.APIClient
	clouddb         *clouddb.APIClient
	monitoring      *monitoring.APIClient
	vpc             *vpc.APIClient
	vserver         *vserver.APIClient
	vnas            *vnas.APIClient
	vautoscaling    *vautoscaling.APIClient
	vloadbalancer   *vloadbalancer.APIClient
	vnks            *vnks.APIClient
	sourcecommit    *sourcecommit.APIClient
	sourcebuild     *sourcebuild.APIClient
	sourcepipeline  *sourcepipeline.APIClient
	vsourcepipeline *vsourcepipeline.APIClient
	vsourcedeploy   *vsourcedeploy.APIClient
	vses            *vses2.APIClient
	vcdss           *vcdss.APIClient
}

func (c *Config) Client() (*NcloudAPIClient, error) {
	apiKey := &ncloud.APIKey{
		AccessKey: c.AccessKey,
		SecretKey: c.SecretKey,
	}
	return &NcloudAPIClient{
		server:          server.NewAPIClient(server.NewConfiguration(apiKey)),
		autoscaling:     autoscaling.NewAPIClient(autoscaling.NewConfiguration(apiKey)),
		loadbalancer:    loadbalancer.NewAPIClient(loadbalancer.NewConfiguration(apiKey)),
		cdn:             cdn.NewAPIClient(cdn.NewConfiguration(apiKey)),
		clouddb:         clouddb.NewAPIClient(clouddb.NewConfiguration(apiKey)),
		monitoring:      monitoring.NewAPIClient(monitoring.NewConfiguration(apiKey)),
		vpc:             vpc.NewAPIClient(vpc.NewConfiguration(apiKey)),
		vserver:         vserver.NewAPIClient(vserver.NewConfiguration(apiKey)),
		vnas:            vnas.NewAPIClient(vnas.NewConfiguration(apiKey)),
		vautoscaling:    vautoscaling.NewAPIClient(vautoscaling.NewConfiguration(apiKey)),
		vloadbalancer:   vloadbalancer.NewAPIClient(vloadbalancer.NewConfiguration(apiKey)),
		vnks:            vnks.NewAPIClient(vnks.NewConfiguration(c.Region, apiKey)),
		sourcecommit:    sourcecommit.NewAPIClient(sourcecommit.NewConfiguration(c.Region, apiKey)),
		sourcebuild:     sourcebuild.NewAPIClient((sourcebuild.NewConfiguration(c.Region, apiKey))),
		sourcepipeline:  sourcepipeline.NewAPIClient(sourcepipeline.NewConfiguration(c.Region, apiKey)),
		vsourcedeploy:   vsourcedeploy.NewAPIClient(vsourcedeploy.NewConfiguration(c.Region, apiKey)),
		vsourcepipeline: vsourcepipeline.NewAPIClient(vsourcepipeline.NewConfiguration(c.Region, apiKey)),
		vses:            vses2.NewAPIClient(vses2.NewConfiguration(c.Region, apiKey)),
		vcdss:           vcdss.NewAPIClient(vcdss.NewConfiguration(c.Region, apiKey)),
	}, nil
}

type ProviderConfig struct {
	Site       string
	SupportVPC bool
	RegionCode string
	RegionNo   string
	Client     *NcloudAPIClient
}
