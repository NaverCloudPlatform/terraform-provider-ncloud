package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceNcloudServerProduct() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudServerProductRead,

		Schema: map[string]*schema.Schema{
			"server_image_product_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "You can get one from `data ncloud_server_images`. This is a required value, and each available server's specification varies depending on the server image product.",
			},
			"product_name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
				Description:  "A regex string to apply to the Server Product list returned.",
				Deprecated:   "use filter instead",
			},
			"exclusion_product_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Enter a product code to exclude from the list.",
			},
			"product_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Enter a product code to search from the list. Use it for a single search.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region code. Get available values using the `data ncloud_regions`.",
				Deprecated:  "use region attribute of provider instead",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Zone code. Get available values using the `data ncloud_zones`.",
			},
			"internet_line_type_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PUBLC", "GLBL"}, false),
				Description:  "Internet line identification code. PUBLC(Public), GLBL(Global). default : PUBLC(Public)",
			},
			"filter": dataSourceFiltersSchema(),

			"product_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Product name",
			},
			"product_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Product type",
			},
			"product_description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Product description",
			},
			"infra_resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Infra resource type",
			},
			"cpu_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "CPU count",
			},
			"memory_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Memory size",
			},
			"base_block_storage_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Base block storage size",
			},
			"platform_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Platform type",
			},
			"os_information": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OS Information",
			},
			"disk_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"add_block_storage_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Additional block storage size",
			},
			"generation_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNcloudServerProductRead(d *schema.ResourceData, meta interface{}) error {
	var resources []map[string]interface{}
	var err error

	if meta.(*ProviderConfig).SupportVPC == true || meta.(*ProviderConfig).Site == "fin" {
		resources, err = getVpcServerProductList(d, meta.(*ProviderConfig))
	} else {
		resources, err = getClassicServerProductList(d, meta.(*ProviderConfig))
	}

	if err != nil {
		return err
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudServerProduct().Schema)
	}

	if err := validateOneResult(len(resources)); err != nil {
		return err
	}

	SetSingularResourceDataFromMap(d, resources[0])

	return nil
}

func getClassicServerProductList(d *schema.ResourceData, config *ProviderConfig) ([]map[string]interface{}, error) {
	client := config.Client
	regionNo := config.RegionNo

	zoneNo, err := parseZoneNoParameter(client, d)
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
	logCommonResponse("getClassicServerProductList", GetCommonResponse(resp))

	resources := []map[string]interface{}{}

	for _, r := range resp.ProductList {
		instance := map[string]interface{}{
			"id":                      *r.ProductCode,
			"product_code":            *r.ProductCode,
			"product_name":            *r.ProductName,
			"product_type":            *r.ProductType.Code,
			"product_description":     *r.ProductDescription,
			"infra_resource_type":     *r.InfraResourceType.Code,
			"cpu_count":               *r.CpuCount,
			"memory_size":             *r.MemorySize,
			"base_block_storage_size": *r.BaseBlockStorageSize,
			"os_information":          *r.OsInformation,
			"disk_type":               *r.DiskType.Code,
			"add_block_storage_size":  *r.AddBlockStorageSize,
		}

		if r.InfraResourceDetailType != nil {
			instance["infra_resource_detail_type_code"] = *r.InfraResourceDetailType.Code
		}
		if r.PlatformType != nil {
			instance["platform_type"] = *r.PlatformType.Code
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
	logCommonResponse("getVpcServerProductList", GetCommonResponse(resp))

	resources := []map[string]interface{}{}

	for _, r := range resp.ProductList {
		instance := map[string]interface{}{
			"id":                      *r.ProductCode,
			"product_code":            *r.ProductCode,
			"product_name":            *r.ProductName,
			"product_type":            *r.ProductType.Code,
			"product_description":     *r.ProductDescription,
			"infra_resource_type":     *r.InfraResourceType.Code,
			"cpu_count":               *r.CpuCount,
			"memory_size":             *r.MemorySize,
			"base_block_storage_size": *r.BaseBlockStorageSize,
			"os_information":          *r.OsInformation,
			"disk_type":               *r.DiskType.Code,
			"add_block_storage_size":  *r.AddBlockStorageSize,
			"generation_code":         *r.GenerationCode,
		}

		if r.InfraResourceDetailType != nil {
			instance["infra_resource_detail_type_code"] = *r.InfraResourceDetailType.Code
		}
		if r.PlatformType != nil {
			instance["platform_type"] = *r.PlatformType.Code
		}

		resources = append(resources, instance)
	}

	return resources, nil
}
