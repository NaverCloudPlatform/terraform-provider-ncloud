package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
)

// Interval for checking status in WaitForXXX method
const DefaultWaitForInterval = 10

// Default timeout
const DefaultTimeout = 60
const DefaultCreateTimeout = 15 * 60
const DefaultStopTimeout = 5 * 60

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
