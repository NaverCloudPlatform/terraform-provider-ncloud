package server

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

func DataSourceNcloudMemberServerImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudMemberServerImageRead,

		Schema: map[string]*schema.Schema{
			"no_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of member server images to view",
			},
			"platform_type_code_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of platform codes of server images to view",
			},
			"filter": DataSourceFiltersSchema(),
			"name_regex": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsValidRegExp),
				Description:      "A regex string to apply to the member server image list returned by ncloud",
				Deprecated:       "use filter instead",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region code. Get available values using the `data ncloud_regions`.",
				Deprecated:  "use region attribute of provider instead",
			},

			"no": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Member server image no",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Member server image name",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Member server image description",
			},
			"original_server_instance_no": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Original server instance no",
			},
			"original_server_product_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Original server product code",
			},
			"original_server_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Original server name",
			},
			"original_base_block_storage_disk_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Original base block storage disk type",
			},
			"original_server_image_product_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Original server image product code",
			},
			"original_os_information": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Original os information",
			},
			"original_server_image_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Original server image name",
			},
			"platform_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Member server image platform type",
			},
			"block_storage_total_rows": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Member server image block storage total rows",
			},
			"block_storage_total_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Member server image block storage total size",
			},
		},
	}
}

func dataSourceNcloudMemberServerImageRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)
	var resources []map[string]interface{}
	var err error

	if config.SupportVPC {
		resources, err = getVpcMemberServerImage(d, config)
	} else {
		resources, err = getClassicMemberServerImage(d, config)
	}

	if err != nil {
		return err
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudMemberServerImage().Schema)
	}

	if err := verify.ValidateOneResult(len(resources)); err != nil {
		return err
	}

	SetSingularResourceDataFromMap(d, resources[0])

	return nil
}

func getClassicMemberServerImage(d *schema.ResourceData, config *conn.ProviderConfig) ([]map[string]interface{}, error) {
	client := config.Client
	regionNo := config.RegionNo

	reqParams := &server.GetMemberServerImageListRequest{
		RegionNo: &regionNo,
	}

	if noList, ok := d.GetOk("no_list"); ok {
		reqParams.MemberServerImageNoList = ExpandStringInterfaceList(noList.([]interface{}))
	}

	if platformTypeCodeList, ok := d.GetOk("platform_type_code_list"); ok {
		reqParams.PlatformTypeCodeList = ExpandStringInterfaceList(platformTypeCodeList.([]interface{}))
	}

	LogCommonRequest("getClassicMemberServerImage", reqParams)

	resp, err := client.Server.V2Api.GetMemberServerImageList(reqParams)
	if err != nil {
		LogErrorResponse("getClassicMemberServerImage", err, reqParams)
		return nil, err
	}
	LogCommonResponse("getClassicMemberServerImage", GetCommonResponse(resp))

	resources := []map[string]interface{}{}

	for _, r := range resp.MemberServerImageList {
		instance := map[string]interface{}{
			"id":                                    *r.MemberServerImageNo,
			"no":                                    *r.MemberServerImageNo,
			"name":                                  *r.MemberServerImageName,
			"description":                           *r.MemberServerImageDescription,
			"original_server_instance_no":           *r.OriginalServerInstanceNo,
			"original_server_product_code":          *r.OriginalServerProductCode,
			"original_server_name":                  *r.OriginalServerName,
			"original_base_block_storage_disk_type": *r.OriginalBaseBlockStorageDiskType.Code,
			"original_server_image_product_code":    *r.OriginalServerImageProductCode,
			"original_os_information":               *r.OriginalOsInformation,
			"original_server_image_name":            *r.OriginalServerImageName,
			"platform_type":                         *r.MemberServerImagePlatformType.Code,
		}

		if r.MemberServerImageBlockStorageTotalRows != nil {
			instance["block_storage_total_rows"] = *r.MemberServerImageBlockStorageTotalRows
		}

		if r.MemberServerImageBlockStorageTotalSize != nil {
			instance["block_storage_total_size"] = *r.MemberServerImageBlockStorageTotalSize
		}

		resources = append(resources, instance)
	}

	return resources, nil
}

func getVpcMemberServerImage(d *schema.ResourceData, config *conn.ProviderConfig) ([]map[string]interface{}, error) {
	client := config.Client
	regionCode := config.RegionCode

	reqParams := &vserver.GetMemberServerImageInstanceListRequest{
		RegionCode: &regionCode,
	}

	if noList, ok := d.GetOk("no_list"); ok {
		reqParams.MemberServerImageInstanceNoList = ExpandStringInterfaceList(noList.([]interface{}))
	}

	if platformTypeCodeList, ok := d.GetOk("platform_type_code_list"); ok {
		reqParams.PlatformTypeCodeList = ExpandStringInterfaceList(platformTypeCodeList.([]interface{}))
	}

	LogCommonRequest("getVpcMemberServerImage", reqParams)

	resp, err := client.Vserver.V2Api.GetMemberServerImageInstanceList(reqParams)
	if err != nil {
		LogErrorResponse("getVpcMemberServerImage", err, reqParams)
		return nil, err
	}
	LogCommonResponse("getVpcMemberServerImage", GetCommonResponse(resp))

	resources := []map[string]interface{}{}

	for _, r := range resp.MemberServerImageInstanceList {
		instance := map[string]interface{}{
			"id":                                 *r.MemberServerImageInstanceNo,
			"no":                                 *r.MemberServerImageInstanceNo,
			"name":                               *r.MemberServerImageName,
			"description":                        *r.MemberServerImageDescription,
			"original_server_instance_no":        *r.OriginalServerInstanceNo,
			"original_server_image_product_code": *r.OriginalServerImageProductCode,
		}

		if r.MemberServerImageBlockStorageTotalRows != nil {
			instance["block_storage_total_rows"] = *r.MemberServerImageBlockStorageTotalRows
		}

		if r.MemberServerImageBlockStorageTotalSize != nil {
			instance["block_storage_total_size"] = *r.MemberServerImageBlockStorageTotalSize
		}

		resources = append(resources, instance)
	}

	return resources, nil
}
