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
				Description:  "A regex string to apply to the member server image list returned by ncloud",
			},
			"member_server_image_no_list": {
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
			"region_code": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Region code.",
				ConflictsWith: []string{"region_no"},
			},
			"region_no": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Region number.",
				ConflictsWith: []string{"region_code"},
			},
			"most_recent": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    true,
				Description: "If more than one result is returned, get the most recent created member server image.",
			},

			"member_server_image_no": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Member server image no",
			},
			"member_server_image_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Member server image name",
			},
			"member_server_image_description": {
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
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        commonCodeSchemaResource,
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
			"member_server_image_status_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Member server image status name",
			},
			"member_server_image_status": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        commonCodeSchemaResource,
				Description: "Member server image status",
			},
			"member_server_image_operation": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        commonCodeSchemaResource,
				Description: "Member server image operation",
			},
			"member_server_image_platform_type": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        commonCodeSchemaResource,
				Description: "Member server image platform type",
			},
			"create_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation date of the member server image",
			},
			"region": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        regionSchemaResource,
				Description: "Region info",
			},
			"member_server_image_block_storage_total_rows": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Member server image block storage total rows",
			},
			"member_server_image_block_storage_total_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Member server image block storage total size",
			},
		},
	}
}

func dataSourceNcloudMemberServerImageRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	regionNo, err := parseRegionNoParameter(conn, d)
	if err != nil {
		return err
	}
	reqParams := &sdk.RequestServerImageList{
		MemberServerImageNoList: StringList(d.Get("member_server_image_no_list").([]interface{})),
		PlatformTypeCodeList:    StringList(d.Get("platform_type_code_list").([]interface{})),
		RegionNo:                regionNo,
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
