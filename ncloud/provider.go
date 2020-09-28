package ncloud

import (
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var NcloudResources map[string]*schema.Resource
var NcloudDatasources map[string]*schema.Resource

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NCLOUD_ACCESS_KEY", nil),
				Description: descriptions["access_key"],
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NCLOUD_SECRET_KEY", nil),
				Description: descriptions["secret_key"],
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NCLOUD_REGION", nil),
				Description: descriptions["region"],
			},
			"site": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NCLOUD_SITE", nil),
				Description: descriptions["site"],
			},
			"support_vpc": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NCLOUD_SUPPORT_VPC", nil),
				Description: descriptions["support_vpc"],
			},
		},
		DataSourcesMap: DataSourcesMap(),
		ResourcesMap:   ResourcesMap(),
		ConfigureFunc:  providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var providerConfig ProviderConfig

	config := Config{
		AccessKey: d.Get("access_key").(string),
		SecretKey: d.Get("secret_key").(string),
	}

	if region, ok := d.GetOk("region"); ok {
		os.Setenv("NCLOUD_REGION", region.(string))
		providerConfig.RegionCode = region.(string)
	}

	if site, ok := d.GetOk("site"); ok {
		providerConfig.Site = site.(string)

		switch site {
		case "gov":
			os.Setenv("NCLOUD_API_GW", "https://ncloud.apigw.gov-ntruss.com")
		case "fin":
			os.Setenv("NCLOUD_API_GW", "https://ncloud.apigw.fin-ntruss.com")
		}
	}

	if supportVpc, ok := d.GetOk("support_vpc"); ok {
		providerConfig.SupportVPC = supportVpc.(bool)
	}

	client, err := config.Client()
	if err != nil {
		return nil, err
	}

	// Fin only supports VPC
	if providerConfig.Site == "fin" {
		providerConfig.SupportVPC = true
	}

	if providerConfig.SupportVPC == false {
		if regionNo, err := parseRegionNoParameter(client, d); err != nil {
			return nil, err
		} else {
			providerConfig.RegionNo = *regionNo
		}
	}

	providerConfig.Client = client

	return &providerConfig, nil
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"access_key":  "Access key of ncloud",
		"secret_key":  "Secret key of ncloud",
		"region":      "Region of ncloud",
		"site":        "Site of ncloud (public / gov / fin)",
		"support_vpc": "Support VPC platform",
	}
}

//RegisterDatasource Register data sources terraform for NAVER CLOUD PLATFORM.
func RegisterDatasource(name string, datasourceSchema *schema.Resource) {
	if NcloudDatasources == nil {
		NcloudDatasources = make(map[string]*schema.Resource)
	}
	NcloudDatasources[name] = datasourceSchema
}

//RegisterResource Register resources terraform for NAVER CLOUD PLATFORM.
func RegisterResource(name string, resourceSchema *schema.Resource) {
	if NcloudResources == nil {
		NcloudResources = make(map[string]*schema.Resource)
	}
	NcloudResources[name] = resourceSchema
}

//DataSourcesMap This returns a map of all data sources to register with Terraform
func DataSourcesMap() map[string]*schema.Resource {
	return NcloudDatasources
}

//ResourcesMap This returns a map of all resources to register with Terraform
func ResourcesMap() map[string]*schema.Resource {
	return NcloudResources
}
