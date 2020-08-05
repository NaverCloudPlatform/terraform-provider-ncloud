package ncloud

import (
	"regexp"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
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
			"add_block_storage_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Additional block storage size",
			},
		},
	}
}

func dataSourceNcloudServerProductRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	regionNo, err := parseRegionNoParameter(client, d)
	if err != nil {
		return err
	}
	zoneNo, err := parseZoneNoParameter(client, d)
	if err != nil {
		return err
	}
	reqParams := &server.GetServerProductListRequest{
		ServerImageProductCode: ncloud.String(d.Get("server_image_product_code").(string)),
		RegionNo:               regionNo,
		ZoneNo:                 zoneNo,
		InternetLineTypeCode:   StringPtrOrNil(d.GetOk("internet_line_type_code")),
	}

	if exclusionProductCode, ok := d.GetOk("exclusion_product_code"); ok {
		reqParams.ExclusionProductCode = ncloud.String(exclusionProductCode.(string))
	}

	if productCode, ok := d.GetOk("product_code"); ok {
		reqParams.ProductCode = ncloud.String(productCode.(string))
	}

	logCommonRequest("GetServerProductList", reqParams)

	resp, err := client.server.V2Api.GetServerProductList(reqParams)
	if err != nil {
		logErrorResponse("GetServerProductList", err, reqParams)
		return err
	}
	logCommonResponse("GetServerProductList", GetCommonResponse(resp))

	var serverProduct *server.Product
	allServerProducts := resp.ProductList
	var filteredServerProducts []*server.Product
	nameRegex, nameRegexOk := d.GetOk("product_name_regex")
	if nameRegexOk {
		r := regexp.MustCompile(nameRegex.(string))
		for _, serverProduct := range allServerProducts {
			if r.MatchString(ncloud.StringValue(serverProduct.ProductName)) {
				filteredServerProducts = append(filteredServerProducts, serverProduct)
			}
		}
	} else {
		filteredServerProducts = allServerProducts[:]
	}

	if err := validateOneResult(len(filteredServerProducts)); err != nil {
		return err
	}
	serverProduct = filteredServerProducts[0]

	return serverProductAttributes(d, serverProduct)
}

func serverProductAttributes(d *schema.ResourceData, product *server.Product) error {
	d.Set("product_code", product.ProductCode)
	d.Set("product_name", product.ProductName)
	d.Set("product_description", product.ProductDescription)
	d.Set("cpu_count", product.CpuCount)
	d.Set("memory_size", product.MemorySize)
	d.Set("base_block_storage_size", product.BaseBlockStorageSize)
	d.Set("os_information", product.OsInformation)
	d.Set("add_block_storage_size", product.AddBlockStorageSize)

	if productType := flattenCommonCode(product.ProductType); productType["code"] != nil {
		d.Set("product_type", productType["code"])
	}

	if infraResourceType := flattenCommonCode(product.InfraResourceType); infraResourceType["code"] != nil {
		d.Set("infra_resource_type", infraResourceType["code"])
	}

	if platformType := flattenCommonCode(product.PlatformType); platformType["code"] != nil {
		d.Set("platform_type", platformType["code"])
	}

	d.SetId(ncloud.StringValue(product.ProductCode))

	return nil
}
