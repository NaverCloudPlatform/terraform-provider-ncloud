package ncloud

import (
	"log"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	RegisterResource("ncloud_placement_group", resourceNcloudPlacementGroup())
}

func resourceNcloudPlacementGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudPlacementGroupCreate,
		Read:   resourceNcloudPlacementGroupRead,
		Delete: resourceNcloudPlacementGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: ToDiagFunc(validateInstanceName),
			},
			"placement_group_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"AA"}, false)),
			},
			"placement_group_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNcloudPlacementGroupCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if !config.SupportVPC {
		return NotSupportClassic("resource `ncloud_placement_group`")
	}

	reqParams := &vserver.CreatePlacementGroupRequest{
		RegionCode: &config.RegionCode,
	}

	if v, ok := d.GetOk("name"); ok {
		reqParams.PlacementGroupName = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("placement_group_type"); ok {
		reqParams.PlacementGroupTypeCode = ncloud.String(v.(string))
	}

	logCommonRequest("CreatePlacementGroup", reqParams)
	resp, err := config.Client.vserver.V2Api.CreatePlacementGroup(reqParams)
	if err != nil {
		logErrorResponse("CreatePlacementGroup", err, reqParams)
		return err
	}

	logResponse("CreatePlacementGroup", resp)

	instance := resp.PlacementGroupList[0]
	d.SetId(*instance.PlacementGroupNo)

	log.Printf("[INFO] Placement Group ID: %s", d.Id())

	return resourceNcloudPlacementGroupRead(d, meta)
}

func resourceNcloudPlacementGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	instance, err := getPlacementGroupInstance(config, d.Id())
	if err != nil {
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*instance.PlacementGroupNo)
	d.Set("placement_group_no", instance.PlacementGroupNo)
	d.Set("placement_group_type", instance.PlacementGroupType.Code)
	d.Set("name", instance.PlacementGroupName)

	return nil
}

func resourceNcloudPlacementGroupDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	reqParams := &vserver.DeletePlacementGroupRequest{
		RegionCode:       &config.RegionCode,
		PlacementGroupNo: ncloud.String(d.Get("placement_group_no").(string)),
	}

	logCommonRequest("DeletePlacementGroup", reqParams)
	resp, err := config.Client.vserver.V2Api.DeletePlacementGroup(reqParams)
	if err != nil {
		logErrorResponse("DeletePlacementGroup", err, reqParams)
		return err
	}

	logResponse("DeletePlacementGroup", resp)

	return nil
}

func getPlacementGroupInstance(config *ProviderConfig, id string) (*vserver.PlacementGroup, error) {
	reqParams := &vserver.GetPlacementGroupDetailRequest{
		RegionCode:       &config.RegionCode,
		PlacementGroupNo: ncloud.String(id),
	}

	logCommonRequest("GetPlacementGroupDetail", reqParams)
	resp, err := config.Client.vserver.V2Api.GetPlacementGroupDetail(reqParams)
	if err != nil {
		logErrorResponse("GetPlacementGroupDetail", err, reqParams)
		return nil, err
	}
	logResponse("GetPlacementGroupDetail", resp)

	if len(resp.PlacementGroupList) > 0 {
		return resp.PlacementGroupList[0], nil
	}

	return nil, nil
}
