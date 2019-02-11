package ncloud

import (
	"fmt"
	"regexp"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceNcloudServerProducts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudServerProductsRead,

		Schema: map[string]*schema.Schema{
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
			"server_image_product_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "You can get one from `data ncloud_server_images`. This is a required value, and each available server's specification varies depending on the server image product.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region code. Get available values using the `data ncloud_regions`.",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Zone code. You can decide a zone where servers are created. You can decide which zone the product list will be requested at. You can get one by calling `data ncloud_zones`. default : Select the first Zone in the specific region",
			},
			"internet_line_type_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PUBLC", "GLBL"}, false),
				Description:  "Internet line identification code. PUBLC(Public), GLBL(Global). default : PUBLC(Public)",
			},
			"server_products": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of server product code.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceNcloudServerProductsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

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

	if len(filteredServerProducts) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	return serverProductsAttributes(d, filteredServerProducts)
}

func serverProductsAttributes(d *schema.ResourceData, serverProduct []*server.Product) error {
	var ids []string

	for _, product := range serverProduct {
		ids = append(ids, ncloud.StringValue(product.ProductCode))
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("server_products", flattenServerProducts(serverProduct)); err != nil {
		return err
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), d.Get("server_products"))
	}

	return nil
}
