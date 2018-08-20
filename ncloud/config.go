package ncloud

import (
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/autoscaling"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/cdn"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/clouddb"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/loadbalancer"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/monitoring"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
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
}

type NcloudAPIClient struct {
	server       *server.APIClient
	autoscaling  *autoscaling.APIClient
	loadbalancer *loadbalancer.APIClient
	cdn          *cdn.APIClient
	clouddb      *clouddb.APIClient
	monitoring   *monitoring.APIClient
}

func (c *Config) Client() (*NcloudAPIClient, error) {
	return &NcloudAPIClient{
		server: server.NewAPIClient(server.NewConfiguration(&server.APIKey{
			AccessKey: c.AccessKey,
			SecretKey: c.SecretKey,
		})),
		autoscaling: autoscaling.NewAPIClient(autoscaling.NewConfiguration(&autoscaling.APIKey{
			AccessKey: c.AccessKey,
			SecretKey: c.SecretKey,
		})),
		loadbalancer: loadbalancer.NewAPIClient(loadbalancer.NewConfiguration(&loadbalancer.APIKey{
			AccessKey: c.AccessKey,
			SecretKey: c.SecretKey,
		})),
		cdn: cdn.NewAPIClient(cdn.NewConfiguration(&cdn.APIKey{
			AccessKey: c.AccessKey,
			SecretKey: c.SecretKey,
		})),
		clouddb: clouddb.NewAPIClient(clouddb.NewConfiguration(&clouddb.APIKey{
			AccessKey: c.AccessKey,
			SecretKey: c.SecretKey,
		})),
		monitoring: monitoring.NewAPIClient(monitoring.NewConfiguration(&monitoring.APIKey{
			AccessKey: c.AccessKey,
			SecretKey: c.SecretKey,
		})),
	}, nil
}
