package ncloud

import (
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NCLOUD_ACCESS_KEY", os.Getenv("NCLOUD_ACCESS_KEY")),
				Description: descriptions["access_key"],
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NCLOUD_SECRET_KEY", os.Getenv("NCLOUD_SECRET_KEY")),
				Description: descriptions["secret_key"],
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NCLOUD_REGION", os.Getenv("NCLOUD_REGION")),
				Description: descriptions["region"],
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ncloud_regions":       dataSourceNcloudRegions(),
			"ncloud_server_images": dataSourceNcloudServerImages(),
			"ncloud_zone":          dataSourceNcloudZones(),
			//"ncloud_instances":     dataSourceNcloudInstances(),

		},
		ResourcesMap: map[string]*schema.Resource{
			"ncloud_instance":      resourceNcloudInstance(),
			"ncloud_block_storage": resourceNcloudBlockStorage(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		AccessKey: d.Get("access_key").(string),
		SecretKey: d.Get("secret_key").(string),
	}

	sdk, err := config.Client()
	if err != nil {
		return nil, err
	}
	return sdk, nil
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"access_key": "Access key of ncloud",
		"secret_key": "Secret key of ncloud",
		"region":     "Region of ncloud",
	}
}
