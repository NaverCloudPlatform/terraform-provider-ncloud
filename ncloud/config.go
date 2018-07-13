package ncloud

import (
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
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

type NcloudSdk struct {
	conn *sdk.Conn
}

func (c *Config) Client() (*NcloudSdk, error) {
	return &NcloudSdk{
		conn: sdk.NewConnection(c.AccessKey, c.SecretKey),
	}, nil
}
