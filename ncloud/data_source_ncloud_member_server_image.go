package ncloud

import (
	"fmt"
	"regexp"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNcloudMemberServerImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudMemberServerImageRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				// ForceNew:     true,
				ValidateFunc: validateRegexp,
				Description:  "A regex string to apply to the member server image list returned by ncloud",
			},
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
			"status_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Member server image status name",
			},
			"status": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        commonCodeSchemaResource,
				Description: "Member server image status",
			},
			"operation": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        commonCodeSchemaResource,
				Description: "Member server image operation",
			},
			"platform_type": {
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
	client := meta.(*NcloudAPIClient)

	regionNo, err := parseRegionNoParameter(client, d)
	if err != nil {
		return err
	}
	reqParams := &server.GetMemberServerImageListRequest{
		MemberServerImageNoList: expandStringInterfaceList(d.Get("no_list").([]interface{})),
		PlatformTypeCodeList:    expandStringInterfaceList(d.Get("platform_type_code_list").([]interface{})),
		RegionNo:                regionNo,
	}

	logCommonRequest("GetMemberServerImageList", reqParams)

	resp, err := client.server.V2Api.GetMemberServerImageList(reqParams)
	if err != nil {
		logErrorResponse("GetMemberServerImageList", err, reqParams)
		return err
	}
	logCommonResponse("GetMemberServerImageList", GetCommonResponse(resp))

	var memberServerImage *server.MemberServerImage

	allMemberServerImages := resp.MemberServerImageList
	var filteredMemberServerImages []*server.MemberServerImage
	nameRegex, nameRegexOk := d.GetOk("name_regex")
	if nameRegexOk {
		r := regexp.MustCompile(nameRegex.(string))
		for _, memberServerImage := range allMemberServerImages {
			if r.MatchString(*memberServerImage.MemberServerImageName) {
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
		memberServerImage = mostRecentMemberServerImage(filteredMemberServerImages)
	} else {
		memberServerImage = filteredMemberServerImages[0]
	}

	return memberServerImageAttributes(d, memberServerImage)
}

func memberServerImageAttributes(d *schema.ResourceData, m *server.MemberServerImage) error {
	d.Set("no", m.MemberServerImageNo)
	d.Set("name", m.MemberServerImageName)
	d.Set("description", m.MemberServerImageDescription)
	d.Set("original_server_instance_no", m.OriginalServerInstanceNo)
	d.Set("original_server_product_code", m.OriginalServerProductCode)
	d.Set("original_server_name", m.OriginalServerName)
	d.Set("original_server_image_product_code", m.OriginalServerImageProductCode)
	d.Set("original_os_information", m.OriginalOsInformation)
	d.Set("original_server_image_name", m.OriginalServerImageName)
	d.Set("status_name", m.MemberServerImageStatusName)
	d.Set("create_date", m.CreateDate)
	d.Set("block_storage_total_rows", m.MemberServerImageBlockStorageTotalRows)
	d.Set("block_storage_total_size", m.MemberServerImageBlockStorageTotalSize)

	if err := d.Set("original_base_block_storage_disk_type", flattenCommonCode(m.OriginalBaseBlockStorageDiskType)); err != nil {
		return err
	}

	if err := d.Set("status", flattenCommonCode(m.MemberServerImageStatus)); err != nil {
		return err
	}

	if err := d.Set("operation", flattenCommonCode(m.MemberServerImageOperation)); err != nil {
		return err
	}

	if err := d.Set("platform_type", flattenCommonCode(m.MemberServerImagePlatformType)); err != nil {
		return err
	}

	if err := d.Set("region", flattenRegion(m.Region)); err != nil {
		return err
	}

	d.SetId(*m.MemberServerImageNo)

	return nil
}
