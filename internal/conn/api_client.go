package conn

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewS3Client(region string, api *ncloud.APIKey, site, endpointFromEnv string) *s3.Client {
	var endpoint string
	if endpointFromEnv != "" {
		endpoint = endpointFromEnv
	} else {
		endpoint = genEndpointWithCode(region, site)
	}

	if api.AccessKey == "" || api.SecretKey == "" {
		log.Fatal("AccessKey and SecretKey must not be empty")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(api.AccessKey, api.SecretKey, "")),
		config.WithRegion(region),
	)

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	newClient := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = ncloud.String(endpoint)
	})

	return newClient
}

func genEndpointWithCode(region, site string) string {
	var s3Endpoint string
	switch site {
	case "gov":
		s3Endpoint = fmt.Sprintf("https://%[1]s.object.gov-ncloudstorage.com", strings.ToLower(region))
	case "fin":
		s3Endpoint = fmt.Sprintf("https://%[1]s.object.fin-ncloudstorage.com", strings.ToLower(region))
	default:
		s3Endpoint = fmt.Sprintf("https://%[1]s.object.ncloudstorage.com", strings.ToLower(region))
	}

	return s3Endpoint
}
