package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func init() {
	RegisterDataSource("ncloud_server_product", dataSourceNcloudServerProduct())
}

func dataSourceNcloudServerProduct() *schema.Resource {
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
			"internet_line_type_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PUBLC", "GLBL"}, false),
			},
			"filter": dataSourceFiltersSchema(),

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
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
				Deprecated:   "use filter instead",
			},
			"exclusion_product_code": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "This field no longer support",
			},
		},
	}
}

func dataSourceNcloudServerProductRead(d *schema.ResourceData, meta interface{}) error {
	resources, err := getClassicServerProductListFiltered(d, meta.(*ProviderConfig))
	if err != nil {
		return err
	}

	if err := validateOneResult(len(resources)); err != nil {

		return err
	}

	SetSingularResourceDataFromMap(d, resources[0])

	return nil
}

func getClassicServerProductListFiltered(d *schema.ResourceData, config *ProviderConfig) ([]map[string]interface{}, error) {
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
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudServerProduct().Schema)
	}

	return resources, nil
}

func getClassicServerProductList(d *schema.ResourceData, config *ProviderConfig) ([]map[string]interface{}, error) {
	client := config.Client
	regionNo := config.RegionNo

	zoneNo, err := parseZoneNoParameter(config, d)
	if err != nil {
		return nil, err
	}
	reqParams := &server.GetServerProductListRequest{
		ExclusionProductCode:   StringPtrOrNil(d.GetOk("exclusion_product_code")),
		ServerImageProductCode: ncloud.String(d.Get("server_image_product_code").(string)),
		ProductCode:            StringPtrOrNil(d.GetOk("product_code")),
		RegionNo:               &regionNo,
		ZoneNo:                 zoneNo,
		InternetLineTypeCode:   StringPtrOrNil(d.GetOk("internet_line_type_code")),
	}

	logCommonRequest("getClassicServerProductList", reqParams)
	resp, err := client.server.V2Api.GetServerProductList(reqParams)
	if err != nil {
		logErrorResponse("getClassicServerProductList", err, reqParams)
		return nil, err
	}
	logResponse("getClassicServerProductList", resp)

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

func getVpcServerProductList(d *schema.ResourceData, config *ProviderConfig) ([]map[string]interface{}, error) {
	client := config.Client
	regionCode := config.RegionCode

	reqParams := &vserver.GetServerProductListRequest{
		ExclusionProductCode:   StringPtrOrNil(d.GetOk("exclusion_product_code")),
		ServerImageProductCode: ncloud.String(d.Get("server_image_product_code").(string)),
		ProductCode:            StringPtrOrNil(d.GetOk("product_code")),
		RegionCode:             &regionCode,
		ZoneCode:               StringPtrOrNil(d.GetOk("zone")),
	}

	logCommonRequest("getVpcServerProductList", reqParams)
	resp, err := client.vserver.V2Api.GetServerProductList(reqParams)
	if err != nil {
		logErrorResponse("getVpcServerProductList", err, reqParams)
		return nil, err
	}
	logResponse("getVpcServerProductList", resp)

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
