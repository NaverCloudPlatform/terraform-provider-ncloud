package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	RegisterDataSource("ncloud_server_image", dataSourceNcloudServerImage())
}

func dataSourceNcloudServerImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudServerImageRead,

		Schema: map[string]*schema.Schema{
			"product_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"platform_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"infra_resource_detail_type_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"base_block_storage_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_information": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Deprecated
			"product_name_regex": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringIsValidRegExp),
				Deprecated:       "use `filter` instead",
			},
			"exclusion_product_code": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "This field no longer support",
			},
			"platform_type_code_list": {
				Type:       schema.TypeList,
				Optional:   true,
				Elem:       &schema.Schema{Type: schema.TypeString},
				Deprecated: "use `filter` or `platform_type` instead",
			},
		},
	}
}

func dataSourceNcloudServerImageRead(d *schema.ResourceData, meta interface{}) error {
	resources, err := getServerImageProductListFiltered(d, meta.(*ProviderConfig))

	if err != nil {
		return err
	}

	if err := validateOneResult(len(resources)); err != nil {
		return err
	}

	SetSingularResourceDataFromMap(d, resources[0])

	return nil
}

func getServerImageProductListFiltered(d *schema.ResourceData, config *ProviderConfig) ([]map[string]interface{}, error) {
	var resources []map[string]interface{}
	var err error

	if config.SupportVPC == true {
		resources, err = getVpcServerImageProductList(d, config)
	} else {
		resources, err = getClassicServerImageProductList(d, config)
	}

	if err != nil {
		return nil, err
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudServerImage().Schema)
	}

	return resources, nil
}

func getClassicServerImageProductList(d *schema.ResourceData, config *ProviderConfig) ([]map[string]interface{}, error) {
	client := config.Client
	regionNo := config.RegionNo

	reqParams := &server.GetServerImageProductListRequest{
		ProductCode:                 StringPtrOrNil(d.GetOk("product_code")),
		RegionNo:                    &regionNo,
		InfraResourceDetailTypeCode: StringPtrOrNil(d.GetOk("infra_resource_detail_type_code")),
	}

	if v, ok := d.GetOk("platform_type"); ok {
		reqParams.PlatformTypeCodeList = []*string{ncloud.String(v.(string))}
	}

	logCommonRequest("GetServerImageProductList", reqParams)
	resp, err := client.server.V2Api.GetServerImageProductList(reqParams)
	if err != nil {
		logErrorResponse("GetServerImageProductList", err, reqParams)
		return nil, err
	}
	logResponse("GetServerImageProductList", resp)

	var resources []map[string]interface{}

	for _, r := range resp.ProductList {
		instance := map[string]interface{}{
			"id":                      *r.ProductCode,
			"product_code":            *r.ProductCode,
			"product_name":            *r.ProductName,
			"product_type":            *r.ProductType.Code,
			"product_description":     *r.ProductDescription,
			"infra_resource_type":     *r.InfraResourceType.Code,
			"base_block_storage_size": fmt.Sprintf("%dGB", *r.BaseBlockStorageSize/GIGABYTE),
			"platform_type":           *r.PlatformType.Code,
			"os_information":          *r.OsInformation,
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
		ProductCode:                 StringPtrOrNil(d.GetOk("product_code")),
		RegionCode:                  &regionCode,
		InfraResourceDetailTypeCode: StringPtrOrNil(d.GetOk("infra_resource_detail_type_code")),
	}

	if v, ok := d.GetOk("platform_type"); ok {
		reqParams.PlatformTypeCodeList = []*string{ncloud.String(v.(string))}
	}

	logCommonRequest("GetServerImageProductList", reqParams)
	resp, err := client.vserver.V2Api.GetServerImageProductList(reqParams)
	if err != nil {
		logErrorResponse("GetServerImageProductList", err, reqParams)
		return nil, err
	}
	logResponse("GetServerImageProductList", resp)

	var resources []map[string]interface{}

	for _, r := range resp.ProductList {
		instance := map[string]interface{}{
			"id":                      *r.ProductCode,
			"product_code":            *r.ProductCode,
			"product_name":            *r.ProductName,
			"product_type":            *r.ProductType.Code,
			"product_description":     *r.ProductDescription,
			"infra_resource_type":     *r.InfraResourceType.Code,
			"base_block_storage_size": fmt.Sprintf("%dGB", *r.BaseBlockStorageSize/GIGABYTE),
			"platform_type":           *r.PlatformType.Code,
			"os_information":          *r.OsInformation,
		}

		if r.InfraResourceDetailType != nil {
			instance["infra_resource_detail_type_code"] = *r.InfraResourceDetailType.Code
		}
		resources = append(resources, instance)
	}

	return resources, nil
}
