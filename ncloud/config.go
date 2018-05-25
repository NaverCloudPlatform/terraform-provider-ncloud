package ncloud

import "github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"

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
