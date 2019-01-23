package ncloud

import (
	"fmt"
	"regexp"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceNcloudServerImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudServerImagesRead,

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
			"region_code": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Region code. Get available values using the `data ncloud_regions`.",
				ConflictsWith: []string{"region_no"},
			},
			"region_no": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Region number. Get available values using the `data ncloud_regions`.",
				ConflictsWith: []string{"region_code"},
			},
			"infra_resource_detail_type_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "infra resource detail type code.",
			},
			"server_images": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"product_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Product Code",
						},
						"product_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Product name",
						},
						"product_type": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        commonCodeSchemaResource,
							Description: "Product type",
						},
						"product_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Product description",
						},
						"infra_resource_type": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        commonCodeSchemaResource,
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
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        commonCodeSchemaResource,
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
							Description: "additional block storage size",
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceNcloudServerImagesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	regionNo, err := parseRegionNoParameter(client, d)
	if err != nil {
		return err
	}
	reqParams := &server.GetServerImageProductListRequest{
		ExclusionProductCode:        StringPtrOrNil(d.GetOk("exclusion_product_code")),
		ProductCode:                 StringPtrOrNil(d.GetOk("product_code")),
		PlatformTypeCodeList:        expandStringInterfaceList(d.Get("platform_type_code_list").([]interface{})),
		RegionNo:                    regionNo,
		InfraResourceDetailTypeCode: StringPtrOrNil(d.GetOk("infra_resource_detail_type_code")),
	}

	logCommonRequest("GetServerImageProductList", reqParams)

	resp, err := client.server.V2Api.GetServerImageProductList(reqParams)
	if err != nil {
		logErrorResponse("GetServerImageProductList", err, reqParams)
		return err
	}
	logCommonResponse("GetServerImageProductList", GetCommonResponse(resp))

	allServerImages := resp.ProductList
	var filteredServerImages []*server.Product
	nameRegex, nameRegexOk := d.GetOk("product_name_regex")
	if nameRegexOk {
		r := regexp.MustCompile(nameRegex.(string))
		for _, serverImage := range allServerImages {
			if r.MatchString(*serverImage.ProductName) {
				filteredServerImages = append(filteredServerImages, serverImage)
				break
			}
		}
	} else {
		filteredServerImages = allServerImages[:]
	}

	if len(filteredServerImages) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	return serverImagesAttributes(d, filteredServerImages)
}

func serverImagesAttributes(d *schema.ResourceData, serverImages []*server.Product) error {
	var ids []string

	for _, product := range serverImages {
		ids = append(ids, *product.ProductCode)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("server_images", flattenServerImages(serverImages)); err != nil {
		return err
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), d.Get("server_images"))
	}

	return nil
}
