package ncloud

import (
	"fmt"
	"regexp"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNcloudServerImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudServerImageRead,

		Schema: map[string]*schema.Schema{
			"product_name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRegexp,
			},
			"exclusion_product_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"product_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"product_type_code": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"platform_type_code_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"block_storage_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"region_no": {
				Type:     schema.TypeString,
				Optional: true,
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

func dataSourceNcloudServerImageRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	reqParams := &sdk.RequestGetServerImageProductList{
		ExclusionProductCode: d.Get("exclusion_product_code").(string),
		ProductCode:          d.Get("product_code").(string),
		PlatformTypeCodeList: StringList(d.Get("platform_type_code_list").([]interface{})),
		RegionNo:             parseRegionNoParameter(conn, d),
	}

	resp, err := conn.GetServerImageProductList(reqParams)
	if err != nil {
		logErrorResponse("GetServerImageProductList", err, reqParams)
		return err
	}
	logCommonResponse("GetServerImageProductList", reqParams, resp.CommonResponse)

	allServerImages := resp.Product
	var serverImage sdk.Product
	var filteredServerImages []sdk.Product

	nameRegex, nameRegexOk := d.GetOk("product_name_regex")
	productTypeCode, productTypeCodeOk := d.GetOk("product_type_code")

	var r *regexp.Regexp
	if nameRegexOk {
		r = regexp.MustCompile(nameRegex.(string))
	}

	if !nameRegexOk && !productTypeCodeOk {
		filteredServerImages = allServerImages[:]
	} else {
		for _, serverImage := range allServerImages {
			if nameRegexOk && r.MatchString(serverImage.ProductName) {
				filteredServerImages = append(filteredServerImages, serverImage)
				break
			} else if productTypeCodeOk && productTypeCode == serverImage.ProductType.Code {
				filteredServerImages = append(filteredServerImages, serverImage)
				break
			}
		}
	}

	if len(filteredServerImages) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	serverImage = filteredServerImages[0]

	return serverImageAttributes(d, serverImage)
}

func serverImageAttributes(d *schema.ResourceData, serverImage sdk.Product) error {
	d.Set("product_code", serverImage.ProductCode)
	d.Set("product_name", serverImage.ProductName)
	d.Set("product_type", setCommonCode(serverImage.ProductType))
	d.Set("product_description", serverImage.ProductDescription)
	d.Set("infra_resource_type", setCommonCode(serverImage.InfraResourceType))
	d.Set("cpu_count", serverImage.CPUCount)
	d.Set("memory_size", serverImage.MemorySize)
	d.Set("base_block_storage_size", serverImage.BaseBlockStorageSize)
	d.Set("platform_type", setCommonCode(serverImage.PlatformType))
	d.Set("os_information", serverImage.OsInformation)
	d.Set("add_block_storage_size", serverImage.AddBlockStroageSize)
	d.SetId(serverImage.ProductCode)

	return nil
}
