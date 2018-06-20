package ncloud

import (
	"fmt"
	"regexp"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNcloudServerProduct() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudServerProductRead,

		Schema: map[string]*schema.Schema{
			"product_name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRegexp,
			},
			"exclusion_product_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Product code to exclude",
			},
			"product_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Product code to search",
			},
			"server_image_product_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Server image product code",
			},
			"cpu_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"region_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"internet_line_type_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInternetLineTypeCode,
				Description:  "Internet line identification code. PUBLC(Public), GLBL(Global). default : PUBLC(Public)",
			},
			"product_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_type": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"product_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"infra_resource_type": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"memory_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"base_block_storage_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"platform_type": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"os_information": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"add_block_storage_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceNcloudServerProductRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	reqParams := &sdk.RequestGetServerProductList{
		ExclusionProductCode:   d.Get("exclusion_product_code").(string),
		ProductCode:            d.Get("product_code").(string),
		ServerImageProductCode: d.Get("server_image_product_code").(string),
		RegionNo:               parseRegionNoParameter(conn, d),
		//ZoneNo:                 d.Get("zone_no").(string),
		//InternetLineTypeCode:   d.Get("internet_line_type_code").(string),
	}

	resp, err := conn.GetServerProductList(reqParams)
	if err != nil {
		logErrorResponse("GetServerProductList", err, reqParams)
		return err
	}
	logCommonResponse("GetServerProductList", reqParams, resp.CommonResponse)

	var serverProduct sdk.Product
	allServerProducts := resp.Product
	var filteredServerProducts []sdk.Product
	nameRegex, nameRegexOk := d.GetOk("product_name_regex")
	if nameRegexOk {
		r := regexp.MustCompile(nameRegex.(string))
		for _, serverProduct := range allServerProducts {
			if r.MatchString(serverProduct.ProductName) {
				filteredServerProducts = append(filteredServerProducts, serverProduct)
			}
		}
	} else {
		filteredServerProducts = allServerProducts[:]
	}

	if len(filteredServerProducts) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	serverProduct = filteredServerProducts[0]

	return serverProductAttributes(d, serverProduct)
}

func serverProductAttributes(d *schema.ResourceData, product sdk.Product) error {
	d.Set("product_code", product.ProductCode)
	d.Set("product_name", product.ProductName)
	d.Set("product_type", setCommonCode(product.ProductType))
	d.Set("product_description", product.ProductDescription)
	d.Set("infra_resource_type", setCommonCode(product.InfraResourceType))
	d.Set("cpu_count", product.CPUCount)
	d.Set("memory_size", product.MemorySize)
	d.Set("base_block_storage_size", product.BaseBlockStorageSize)
	d.Set("platform_type", setCommonCode(product.PlatformType))
	d.Set("os_information", product.OsInformation)
	d.Set("add_block_storage_size", product.AddBlockStroageSize)

	d.SetId(product.ProductCode)

	return nil
}
