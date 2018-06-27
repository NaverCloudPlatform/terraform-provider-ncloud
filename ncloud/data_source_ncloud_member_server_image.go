package ncloud

import (
	"fmt"
	"regexp"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNcloudMemberServerImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudMemberServerImageRead,

		Schema: map[string]*schema.Schema{
			"member_server_image_name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				// ForceNew:     true,
				ValidateFunc: validateRegexp,
				Description:  "A regex string to apply to the server image list returned by ncloud",
			},
			"member_server_image_no_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "",
			},
			"platform_type_code_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "",
			},
			"region_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"most_recent": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    true,
				Description: "",
			},

			"member_server_image_no": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			"member_server_image_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			"member_server_image_description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			"original_server_instance_no": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			"original_server_product_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			"original_server_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			"original_base_block_storage_disk_type": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        commonCodeSchemaResource,
				Description: "",
			},
			"original_server_image_product_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			"original_os_information": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			"original_server_image_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			"member_server_image_status_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			"member_server_image_status": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        commonCodeSchemaResource,
				Description: "",
			},
			"member_server_image_operation": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        commonCodeSchemaResource,
				Description: "",
			},
			"member_server_image_platform_type": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        commonCodeSchemaResource,
				Description: "",
			},
			"create_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			"region": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        regionSchemaResource,
				Description: "",
			},
			"member_server_image_block_storage_total_rows": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "",
			},
			"member_server_image_block_storage_total_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "",
			},
		},
	}
}

func dataSourceNcloudMemberServerImageRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	reqParams := &sdk.RequestServerImageList{
		MemberServerImageNoList: StringList(d.Get("member_server_image_no_list").([]interface{})),
		PlatformTypeCodeList:    StringList(d.Get("platform_type_code_list").([]interface{})),
		RegionNo:                parseRegionNoParameter(conn, d),
	}

	resp, err := conn.GetMemberServerImageList(reqParams)
	if err != nil {
		logErrorResponse("GetMemberServerImageList", err, reqParams)
		return err
	}
	logCommonResponse("GetMemberServerImageList", reqParams, resp.CommonResponse)

	var memberServerImage sdk.ServerImage

	allMemberServerImages := resp.MemberServerImageList
	var filteredMemberServerImages []sdk.ServerImage
	nameRegex, nameRegexOk := d.GetOk("member_server_image_name_regex")
	if nameRegexOk {
		r := regexp.MustCompile(nameRegex.(string))
		for _, memberServerImage := range allMemberServerImages {
			if r.MatchString(memberServerImage.MemberServerImageName) {
				filteredMemberServerImages = append(filteredMemberServerImages, memberServerImage)
			}
		}
	} else {
		filteredMemberServerImages = allMemberServerImages[:]
	}

	if len(filteredMemberServerImages) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	if len(filteredMemberServerImages) > 1 && d.Get("most_recent").(bool) {
		// Query returned single result.
		memberServerImage = mostRecentServerImage(filteredMemberServerImages)
	} else {
		memberServerImage = filteredMemberServerImages[0]
	}

	return memberServerImageAttributes(d, memberServerImage)
}

func memberServerImageAttributes(d *schema.ResourceData, m sdk.ServerImage) error {
	d.Set("member_server_image_no", m.MemberServerImageNo)
	d.Set("member_server_image_name", m.MemberServerImageName)
	d.Set("member_server_image_description", m.MemberServerImageDescription)
	d.Set("original_server_instance_no", m.OriginalServerInstanceNo)
	d.Set("original_server_product_code", m.OriginalServerProductCode)
	d.Set("original_server_name", m.OriginalServerName)
	d.Set("original_base_block_storage_disk_type", setCommonCode(m.OriginalBaseBlockStorageDiskType))
	d.Set("original_server_image_product_code", m.OriginalServerImageProductCode)
	d.Set("original_os_information", m.OriginalOsInformation)
	d.Set("original_server_image_name", m.OriginalServerImageName)
	d.Set("member_server_image_status_name", m.MemberServerImageStatusName)
	d.Set("member_server_image_status", setCommonCode(m.MemberServerImageStatus))
	d.Set("member_server_image_operation", setCommonCode(m.MemberServerImageOperation))
	d.Set("member_server_image_platform_type", setCommonCode(m.MemberServerImagePlatformType))
	d.Set("create_date", m.CreateDate)
	d.Set("region", setRegion(m.Region))
	d.Set("member_server_image_block_storage_total_rows", m.MemberServerImageBlockStorageTotalRows)
	d.Set("member_server_image_block_storage_total_size", m.MemberServerImageBlockStorageTotalSize)

	d.SetId(m.MemberServerImageNo)

	return nil
}
