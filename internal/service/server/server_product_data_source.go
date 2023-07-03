package server

import (
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/zone"
)

func DataSourceNcloudServerProduct() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudServerProductRead,

		Schema: map[string]*schema.Schema{
			"server_image_product_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"product_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Deprecated
			"internet_line_type_code": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: verify.ToDiagFunc(validation.StringInSlice([]string{"PUBLC", "GLBL"}, false)),
				Deprecated:       "This parameter is no longer used.",
			},
			"filter": DataSourceFiltersSchema(),

			"product_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"infra_resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"base_block_storage_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"generation_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Deprecated
			"product_name_regex": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateDiagFunc: verify.ToDiagFunc(validation.StringIsValidRegExp),
				Deprecated:       "use filter instead",
			},
			"exclusion_product_code": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "This field is no longer support",
			},
		},
	}
}

func dataSourceNcloudServerProductRead(d *schema.ResourceData, meta interface{}) error {
	resources, err := getServerProductListFiltered(d, meta.(*conn.ProviderConfig))
	if err != nil {
		return err
	}

	if err := verify.ValidateOneResult(len(resources)); err != nil {

		return err
	}

	SetSingularResourceDataFromMap(d, resources[0])

	return nil
}

func getServerProductListFiltered(d *schema.ResourceData, config *conn.ProviderConfig) ([]map[string]interface{}, error) {
	var resources []map[string]interface{}
	var err error

	if config.SupportVPC == true {
		resources, err = getVpcServerProductList(d, config)
	} else {
		resources, err = getClassicServerProductList(d, config)
	}

	if err != nil {
		return nil, err
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudServerProduct().Schema)
	}

	return resources, nil
}

func getClassicServerProductList(d *schema.ResourceData, config *conn.ProviderConfig) ([]map[string]interface{}, error) {
	client := config.Client
	regionNo := config.RegionNo

	zoneNo, err := zone.ParseZoneNoParameter(config, d)
	if err != nil {
		return nil, err
	}
	reqParams := &server.GetServerProductListRequest{
		ExclusionProductCode:   StringPtrOrNil(d.GetOk("exclusion_product_code")),
		ServerImageProductCode: ncloud.String(d.Get("server_image_product_code").(string)),
		ProductCode:            StringPtrOrNil(d.GetOk("product_code")),
		RegionNo:               &regionNo,
		ZoneNo:                 zoneNo,
	}

	LogCommonRequest("getClassicServerProductList", reqParams)
	resp, err := client.Server.V2Api.GetServerProductList(reqParams)
	if err != nil {
		LogErrorResponse("getClassicServerProductList", err, reqParams)
		return nil, err
	}
	LogResponse("getClassicServerProductList", resp)

	var resources []map[string]interface{}

	for _, r := range resp.ProductList {
		instance := map[string]interface{}{
			"id":                      *r.ProductCode,
			"product_code":            *r.ProductCode,
			"product_name":            *r.ProductName,
			"product_type":            *r.ProductType.Code,
			"product_description":     *r.ProductDescription,
			"infra_resource_type":     *r.InfraResourceType.Code,
			"cpu_count":               *r.CpuCount,
			"memory_size":             fmt.Sprintf("%dGB", *r.MemorySize/GIGABYTE),
			"base_block_storage_size": fmt.Sprintf("%dGB", *r.BaseBlockStorageSize/GIGABYTE),
			"disk_type":               *r.DiskType.Code,
			"generation_code":         *r.GenerationCode,
		}

		resources = append(resources, instance)
	}

	return resources, nil
}

func getVpcServerProductList(d *schema.ResourceData, config *conn.ProviderConfig) ([]map[string]interface{}, error) {
	client := config.Client
	regionCode := config.RegionCode

	reqParams := &vserver.GetServerProductListRequest{
		ExclusionProductCode:   StringPtrOrNil(d.GetOk("exclusion_product_code")),
		ServerImageProductCode: ncloud.String(d.Get("server_image_product_code").(string)),
		ProductCode:            StringPtrOrNil(d.GetOk("product_code")),
		RegionCode:             &regionCode,
		ZoneCode:               StringPtrOrNil(d.GetOk("zone")),
	}

	LogCommonRequest("getVpcServerProductList", reqParams)
	resp, err := client.Vserver.V2Api.GetServerProductList(reqParams)
	if err != nil {
		LogErrorResponse("getVpcServerProductList", err, reqParams)
		return nil, err
	}
	LogResponse("getVpcServerProductList", resp)

	var resources []map[string]interface{}

	for _, r := range resp.ProductList {
		instance := map[string]interface{}{
			"id":                      *r.ProductCode,
			"product_code":            *r.ProductCode,
			"product_name":            *r.ProductName,
			"product_type":            *r.ProductType.Code,
			"product_description":     *r.ProductDescription,
			"infra_resource_type":     *r.InfraResourceType.Code,
			"cpu_count":               *r.CpuCount,
			"memory_size":             fmt.Sprintf("%dGB", *r.MemorySize/GIGABYTE),
			"base_block_storage_size": fmt.Sprintf("%dGB", *r.BaseBlockStorageSize/GIGABYTE),
			"disk_type":               *r.DiskType.Code,
			"generation_code":         *r.GenerationCode,
		}

		resources = append(resources, instance)
	}

	return resources, nil
}
