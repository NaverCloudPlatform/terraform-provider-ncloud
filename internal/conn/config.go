package conn

import (
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vhadoop"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vautoscaling"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vcdss"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vloadbalancer"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnks"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vredis"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vses2"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnas"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/autoscaling"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/cdn"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/clouddb"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/loadbalancer"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/sourcebuild"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/sourcecommit"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/sourcepipeline"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmongodb"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmssql"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmysql"
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
const s3Endpoint = "https://kr.object.ncloudstorage.com/"

var version = ""

type Config struct {
	AccessKey string
	SecretKey string
	Region    string
	Endpoint  string
}

type NcloudAPIClient struct {
	Server          *server.APIClient
	Autoscaling     *autoscaling.APIClient
	Loadbalancer    *loadbalancer.APIClient
	Cdn             *cdn.APIClient
	Clouddb         *clouddb.APIClient
	Vpc             *vpc.APIClient
	Vserver         *vserver.APIClient
	Vnas            *vnas.APIClient
	Vautoscaling    *vautoscaling.APIClient
	Vloadbalancer   *vloadbalancer.APIClient
	Vnks            *vnks.APIClient
	Sourcecommit    *sourcecommit.APIClient
	Sourcebuild     *sourcebuild.APIClient
	Sourcepipeline  *sourcepipeline.APIClient
	Vsourcepipeline *vsourcepipeline.APIClient
	Vsourcedeploy   *vsourcedeploy.APIClient
	Vses            *vses2.APIClient
	Vcdss           *vcdss.APIClient
	Vmysql          *vmysql.APIClient
	Vmongodb        *vmongodb.APIClient
	Vmssql          *vmssql.APIClient
	Vhadoop         *vhadoop.APIClient
	Vredis          *vredis.APIClient
	ObjectStorage   *s3.Client
}

func (c *Config) Client() (*NcloudAPIClient, error) {
	apiKey := &ncloud.APIKey{
		AccessKey: c.AccessKey,
		SecretKey: c.SecretKey,
	}

	return &NcloudAPIClient{
		Server:          server.NewAPIClient(server.NewConfiguration(apiKey)),
		Autoscaling:     autoscaling.NewAPIClient(autoscaling.NewConfiguration(apiKey)),
		Loadbalancer:    loadbalancer.NewAPIClient(loadbalancer.NewConfiguration(apiKey)),
		Cdn:             cdn.NewAPIClient(cdn.NewConfiguration(apiKey)),
		Clouddb:         clouddb.NewAPIClient(clouddb.NewConfiguration(apiKey)),
		Vpc:             vpc.NewAPIClient(vpc.NewConfiguration(apiKey)),
		Vserver:         vserver.NewAPIClient(vserver.NewConfiguration(apiKey)),
		Vnas:            vnas.NewAPIClient(vnas.NewConfiguration(apiKey)),
		Vautoscaling:    vautoscaling.NewAPIClient(vautoscaling.NewConfiguration(apiKey)),
		Vloadbalancer:   vloadbalancer.NewAPIClient(vloadbalancer.NewConfiguration(apiKey)),
		Vnks:            vnks.NewAPIClient(vnks.NewConfigurationWithUserAgent(c.Region, fmt.Sprintf("Ncloud Terraform Provider/%s", version), apiKey)),
		Sourcecommit:    sourcecommit.NewAPIClient(sourcecommit.NewConfiguration(c.Region, apiKey)),
		Sourcebuild:     sourcebuild.NewAPIClient((sourcebuild.NewConfiguration(c.Region, apiKey))),
		Sourcepipeline:  sourcepipeline.NewAPIClient(sourcepipeline.NewConfiguration(c.Region, apiKey)),
		Vsourcedeploy:   vsourcedeploy.NewAPIClient(vsourcedeploy.NewConfiguration(c.Region, apiKey)),
		Vsourcepipeline: vsourcepipeline.NewAPIClient(vsourcepipeline.NewConfiguration(c.Region, apiKey)),
		Vses:            vses2.NewAPIClient(vses2.NewConfiguration(c.Region, apiKey)),
		Vcdss:           vcdss.NewAPIClient(vcdss.NewConfiguration(c.Region, apiKey)),
		Vmysql:          vmysql.NewAPIClient(vmysql.NewConfiguration(apiKey)),
		Vmongodb:        vmongodb.NewAPIClient(vmongodb.NewConfiguration(apiKey)),
		Vmssql:          vmssql.NewAPIClient(vmssql.NewConfiguration(apiKey)),
		Vhadoop:         vhadoop.NewAPIClient(vhadoop.NewConfiguration(apiKey)),
		Vredis:          vredis.NewAPIClient(vredis.NewConfiguration(apiKey)),
		ObjectStorage:   NewS3Client(c.Region, apiKey, s3Endpoint),
	}, nil
}

type ProviderConfig struct {
	Site       string
	SupportVPC bool
	RegionCode string
	RegionNo   string
	Client     *NcloudAPIClient
}
