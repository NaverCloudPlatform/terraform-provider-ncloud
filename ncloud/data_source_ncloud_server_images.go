package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func dataSourceNcloudServerImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudServerImagesRead,

		Schema: map[string]*schema.Schema{
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
			"server_images": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"product_code": {
							Type:     schema.TypeString,
							Computed: true,
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
						"cpu_count": {
							Type:     schema.TypeInt,
							Computed: true,
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
				},
			},
		},
	}
}

func dataSourceNcloudServerImagesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	reqParams := new(sdk.RequestGetServerImageProductList)
	// TODO: reqParams
	resp, err := conn.GetServerImageProductList(reqParams)
	logCommonResponse("GetServerImageProductList", err, reqParams, resp.CommonResponse)
	// log.Printf("[DEBUG] GetServerImageProductList: %v", resp.Product)
	if err != nil {
		return err
	}

	allServerImages := resp.Product
	var filteredServerImages []sdk.Product
	for _, serverImage := range allServerImages {
		if v, ok := d.GetOk("product_code"); ok && v.(string) != "" && serverImage.ProductCode != v.(string) {
			continue
		}
		filteredServerImages = append(filteredServerImages, serverImage)
	}

	if len(filteredServerImages) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	log.Printf("[DEBUG] ncloud_server_images - Server Images found: %#v", allServerImages)

	return serverImagesAttributes(d, filteredServerImages)
}

func serverImagesAttributes(d *schema.ResourceData, serverImages []sdk.Product) error {
	var ids []string
	var s []map[string]interface{}
	for _, product := range serverImages {
		mapping := map[string]interface{}{
			"product_code": product.ProductCode,
			"product_name": product.ProductName,
			"product_type": map[string]interface{}{
				"code":      product.ProductType.Code,
				"code_name": product.ProductType.CodeName,
			},
			"product_description": product.ProductDescription,
			"infra_resource_type": map[string]interface{}{
				"code":      product.InfraResourceType.Code,
				"code_name": product.InfraResourceType.CodeName,
			},
			"cpu_count":               product.CPUCount,
			"memory_size":             product.MemorySize,
			"base_block_storage_size": product.BaseBlockStorageSize,
			"platform_type": map[string]interface{}{
				"code":      product.PlatformType.Code,
				"code_name": product.PlatformType.CodeName,
			},
			"os_information":         product.OsInformation,
			"add_block_storage_size": product.AddBlockStroageSize,
		}

		ids = append(ids, product.ProductCode)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("server_images", s); err != nil {
		return err
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
