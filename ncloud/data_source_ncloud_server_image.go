package ncloud

import (
	"regexp"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
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
			},
			"exclusion_product_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
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
				Optional:    true,
				Computed:    true,
				Description: "Product type code",
			},
			"platform_type_code_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Values required for identifying platforms in list-type.",
			},
			"block_storage_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Block storage size.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region code. Get available values using the `data ncloud_regions`.",
			},
			"infra_resource_detail_type_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "infra resource detail type code.",
			},

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
		},
	}
}

func dataSourceNcloudServerImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	regionNo, err := parseRegionNoParameter(client, d)
	if err != nil {
		return err
	}
	reqParams := &server.GetServerImageProductListRequest{
		ExclusionProductCode:        StringPtrOrNil(d.GetOk("exclusion_product_code")),
		ProductCode:                 StringPtrOrNil(d.GetOk("product_code")),
		RegionNo:                    regionNo,
		InfraResourceDetailTypeCode: StringPtrOrNil(d.GetOk("infra_resource_detail_type_code")),
	}

	if platformTypeCodeList, ok := d.GetOk("platform_type_code_list"); ok {
		reqParams.PlatformTypeCodeList = expandStringInterfaceList(platformTypeCodeList.([]interface{}))
	}

	logCommonRequest("GetServerImageProductList", reqParams)

	resp, err := client.server.V2Api.GetServerImageProductList(reqParams)
	if err != nil {
		logErrorResponse("GetServerImageProductList", err, reqParams)
		return err
	}
	logCommonResponse("GetServerImageProductList", GetCommonResponse(resp))

	allServerImages := resp.ProductList
	var serverImage *server.Product
	var filteredServerImages []*server.Product

	nameRegex, nameRegexOk := d.GetOk("product_name_regex")
	productTypeCode, productTypeCodeOk := d.GetOk("product_type")

	var r *regexp.Regexp
	if nameRegexOk {
		r = regexp.MustCompile(nameRegex.(string))
	}

	if !nameRegexOk && !productTypeCodeOk {
		filteredServerImages = allServerImages[:]
	} else {
		for _, serverImage := range allServerImages {
			if nameRegexOk && r.MatchString(ncloud.StringValue(serverImage.ProductName)) {
				filteredServerImages = append(filteredServerImages, serverImage)
				break
			} else if productTypeCodeOk && productTypeCode == ncloud.StringValue(serverImage.ProductType.Code) {
				filteredServerImages = append(filteredServerImages, serverImage)
				break
			}
		}
	}

	if err := validateOneResult(len(filteredServerImages)); err != nil {
		return err
	}
	serverImage = filteredServerImages[0]

	return serverImageAttributes(d, serverImage)
}

func serverImageAttributes(d *schema.ResourceData, serverImage *server.Product) error {
	d.Set("product_code", serverImage.ProductCode)
	d.Set("product_name", serverImage.ProductName)
	d.Set("product_description", serverImage.ProductDescription)
	d.Set("cpu_count", serverImage.CpuCount)
	d.Set("memory_size", serverImage.MemorySize)
	d.Set("base_block_storage_size", serverImage.BaseBlockStorageSize)
	d.Set("os_information", serverImage.OsInformation)
	d.Set("add_block_storage_size", serverImage.AddBlockStorageSize)

	if productType := flattenCommonCode(serverImage.ProductType); productType["code"] != nil {
		d.Set("product_type", productType["code"])
	}

	if infraResourceType := flattenCommonCode(serverImage.InfraResourceType); infraResourceType["code"] != nil {
		d.Set("infra_resource_type", infraResourceType["code"])
	}

	if platformType := flattenCommonCode(serverImage.PlatformType); platformType["code"] != nil {
		d.Set("platform_type", platformType["code"])
	}

	d.SetId(*serverImage.ProductCode)

	return nil
}
