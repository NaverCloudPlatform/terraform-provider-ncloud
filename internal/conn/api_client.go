package conn

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/ncloudsdk"
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

// API docs: https://api.ncloud-docs.com/docs/platform-region-getregionlist
// Common object storage docs; https://api.ncloud-docs.com/docs/storage-objectstorage
func genEndpointWithCode(region, site string) string {
	var s3Endpoint string
	switch site {
	case "gov":
		s3Endpoint = fmt.Sprintf("https://%[1]s.object.gov-ncloudstorage.com", strings.ToLower(region))
	case "fin":
		s3Endpoint = "https://kr.object.fin-ncloudstorage.com"
	default:
		s3Endpoint = fmt.Sprintf("https://%[1]s.object.ncloudstorage.com", strings.ToLower(region[:2]))
	}

	return s3Endpoint
}

func NewApigwClient(site string, api *ncloud.APIKey) *ncloudsdk.NClient {
	var baseURL string

	switch site {
	case "gov":
		baseURL = "https://apigateway.apigw.gov-ntruss.com/api/v1"
	case "fin":
		baseURL = "https://apigateway.apigw.fin-ntruss.com/api/v1"
	default:
		baseURL = "https://apigateway.apigw.ntruss.com/api/v1"
	}

	return &ncloudsdk.NClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
		AccessKey:  api.AccessKey,
		SecretKey:  api.SecretKey,
	}
}
