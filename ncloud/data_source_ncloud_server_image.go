package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceNcloudServerImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudServerImageRead,

		Schema: map[string]*schema.Schema{
			"product_name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
				Description:  "A regex string to apply to the server image list returned by ncloud.",
				Deprecated:   "use filter instead",
			},
			"exclusion_product_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Product code you want to exclude from the list.",
			},
			"product_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Product code you want to view on the list. Use this when searching for 1 product.",
			},
			"product_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Product type code",
			},
			"platform_type_code_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Values required for identifying platforms in list-type.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region code. Get available values using the `data ncloud_regions`.",
				Deprecated:  "use region attribute of provider instead",
			},
			"infra_resource_detail_type_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "infra resource detail type code.",
			},
			"filter": dataSourceFiltersSchema(),

			"product_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Product name",
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

func dataSourceNcloudServerImageRead(d *schema.ResourceData, meta interface{}) error {
	var resources []map[string]interface{}
	var err error

	if meta.(*ProviderConfig).SupportVPC == true || meta.(*ProviderConfig).Site == "fin" {
		resources, err = getVpcServerImageProductList(d, meta.(*ProviderConfig))
	} else {
		resources, err = getClassicServerImageProductList(d, meta.(*ProviderConfig))
	}

	if err != nil {
		return err
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudServerImage().Schema)
	}

	if err := validateOneResult(len(resources)); err != nil {
		return err
	}

	SetSingularResourceDataFromMap(d, resources[0])

	return nil
}

func getClassicServerImageProductList(d *schema.ResourceData, config *ProviderConfig) ([]map[string]interface{}, error) {
	client := config.Client
	regionNo := config.RegionNo

	reqParams := &server.GetServerImageProductListRequest{
		ExclusionProductCode:        StringPtrOrNil(d.GetOk("exclusion_product_code")),
		ProductCode:                 StringPtrOrNil(d.GetOk("product_code")),
		RegionNo:                    &regionNo,
		InfraResourceDetailTypeCode: StringPtrOrNil(d.GetOk("infra_resource_detail_type_code")),
	}

	if platformTypeCodeList, ok := d.GetOk("platform_type_code_list"); ok {
		reqParams.PlatformTypeCodeList = expandStringInterfaceList(platformTypeCodeList.([]interface{}))
	}

	logCommonRequest("GetServerImageProductList", reqParams)
	resp, err := client.server.V2Api.GetServerImageProductList(reqParams)
	if err != nil {
		logErrorResponse("GetServerImageProductList", err, reqParams)
		return nil, err
	}
	logResponse("GetServerImageProductList", resp)

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
			"platform_type":           *r.PlatformType.Code,
			"os_information":          *r.OsInformation,
			"add_block_storage_size":  *r.AddBlockStorageSize,
		}

		if r.InfraResourceDetailType != nil {
			instance["infra_resource_detail_type_code"] = *r.InfraResourceDetailType.Code
		}

		resources = append(resources, instance)
	}

	return resources, nil
}

func getVpcServerImageProductList(d *schema.ResourceData, config *ProviderConfig) ([]map[string]interface{}, error) {
	client := config.Client
	regionCode := config.RegionCode

	reqParams := &vserver.GetServerImageProductListRequest{
		ExclusionProductCode: StringPtrOrNil(d.GetOk("exclusion_product_code")),
		ProductCode:          StringPtrOrNil(d.GetOk("product_code")),
		RegionCode:           &regionCode,
	}

	if platformTypeCodeList, ok := d.GetOk("platform_type_code_list"); ok {
		reqParams.PlatformTypeCodeList = expandStringInterfaceList(platformTypeCodeList.([]interface{}))
	}

	logCommonRequest("GetServerImageProductList", reqParams)
	resp, err := client.vserver.V2Api.GetServerImageProductList(reqParams)
	if err != nil {
		logErrorResponse("GetServerImageProductList", err, reqParams)
		return nil, err
	}
	logCommonResponse("GetServerImageProductList", GetCommonResponse(resp))

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
			"platform_type":           *r.PlatformType.Code,
			"os_information":          *r.OsInformation,
			"add_block_storage_size":  *r.AddBlockStorageSize,
			"generation_code":         *r.GenerationCode,
		}

		if r.InfraResourceDetailType != nil {
			instance["infra_resource_detail_type_code"] = *r.InfraResourceDetailType.Code
		}

		resources = append(resources, instance)
	}

	return resources, nil
}
