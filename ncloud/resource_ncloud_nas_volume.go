package ncloud

import (
	"fmt"
	"log"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceNcloudNasVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudNasVolumeCreate,
		Read:   resourceNcloudNasVolumeRead,
		Delete: resourceNcloudNasVolumeDelete,
		Update: resourceNcloudNasVolumeUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DefaultCreateTimeout),
			Delete: schema.DefaultTimeout(DefaultTimeout),
		},

		Schema: map[string]*schema.Schema{
			"volume_name_postfix": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 30),
				Description:  "Name of a NAS volume to create. Enter a volume name that can be 3-20 characters in length after the name already entered for user identification.",
			},
			"volume_size_gb": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(500, 10000),
				Description:  "Enter the nas volume size to be created. You can enter in GB units.",
			},
			"volume_allotment_protocol_type_code": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 5),
				Description:  "Volume allotment protocol type code. `NFS` | `CIFS`",
			},
			"server_instance_no_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of server instance numbers for which access to NFS is to be controlled",
			},
			"custom_ip_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "To add a server of another account to the NAS volume, enter a private IP address of the server.",
			},
			"cifs_user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CIFS user name. The ID must contain a combination of English alphabet and numbers, which can be 6-20 characters in length.",
			},
			"cifs_user_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CIFS user password. The password must contain a combination of at least 2 English letters, numbers and special characters, which can be 8-14 characters in length.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
				Description:  "NAS volume description",
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
			"zone_code": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Zone code. Zone in which you want to create a NAS volume.",
				ConflictsWith: []string{"zone_no"},
			},
			"zone_no": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Zone number. Zone in which you want to create a NAS volume.",
				ConflictsWith: []string{"zone_code"},
			},

			"volume_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "NAS volume name.",
			},
			"instance_status": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        commonCodeSchemaResource,
				Description: "NAS Volume instance status",
			},
			"volume_allotment_protocol_type": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        commonCodeSchemaResource,
				Description: "Volume allotment protocol type.",
			},
			"volume_total_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Volume total size",
			},
			"volume_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Volume size",
			},
			"volume_use_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Volume use size",
			},
			"volume_use_ratio": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Volume use ratio",
			},
			"snapshot_volume_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Snapshot volume size",
			},
			"snapshot_volume_use_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Snapshot volume use size",
			},
			"snapshot_volume_use_ratio": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Snapshot volume use ratio",
			},
			"is_snapshot_configuration": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether a snapshot volume is set.",
			},
			"is_event_configuration": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the event is set.",
			},
			"instance_custom_ip_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "NAS volume instance custom IP list",
			},
			"zone": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        zoneSchemaResource,
				Description: "Zone info",
			},
			"region": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        regionSchemaResource,
				Description: "Region info",
			},
		},
	}
}

func resourceNcloudNasVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	reqParams, err := buildCreateNasVolumeInstanceParams(client, d)
	if err != nil {
		return nil
	}
	logCommonRequest("CreateNasVolumeInstance", reqParams)

	resp, err := client.server.V2Api.CreateNasVolumeInstance(reqParams)
	if err != nil {
		logErrorResponse("CreateNasVolumeInstance", err, reqParams)
		return err
	}
	logCommonResponse("CreateNasVolumeInstance", GetCommonResponse(resp))

	nasVolumeInstance := resp.NasVolumeInstanceList[0]
	d.SetId(ncloud.StringValue(nasVolumeInstance.NasVolumeInstanceNo))

	stateConf := &resource.StateChangeConf{
		Pending: []string{"INIT"},
		Target:  []string{"CREAT"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getNasVolumeInstance(client, ncloud.StringValue(nasVolumeInstance.NasVolumeInstanceNo))

			if err != nil {
				return 0, "", err
			}
			return instance, ncloud.StringValue(instance.NasVolumeInstanceStatus.Code), nil
		},
		Timeout:    DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for NasVolumeInstance state to be \"CREAT\": %s", err)
	}

	return resourceNcloudNasVolumeRead(d, meta)
}

func resourceNcloudNasVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	nasVolume, err := getNasVolumeInstance(client, d.Id())
	if err != nil {
		return err
	}

	if nasVolume != nil {
		d.Set("description", nasVolume.NasVolumeInstanceDescription)
		d.Set("volume_name", nasVolume.VolumeName)
		d.Set("volume_total_size", nasVolume.VolumeTotalSize)
		d.Set("volume_size", nasVolume.VolumeSize)
		d.Set("volume_use_size", nasVolume.VolumeUseSize)
		d.Set("volume_use_ratio", nasVolume.VolumeUseRatio)
		d.Set("snapshot_volume_size", nasVolume.SnapshotVolumeSize)
		d.Set("snapshot_volume_use_size", nasVolume.SnapshotVolumeUseSize)
		d.Set("snapshot_volume_use_ratio", nasVolume.SnapshotVolumeUseRatio)
		d.Set("is_snapshot_configuration", nasVolume.IsSnapshotConfiguration)
		d.Set("is_event_configuration", nasVolume.IsEventConfiguration)
		d.Set("instance_custom_ip_list", nasVolume.NasVolumeInstanceCustomIpList)

		if err := d.Set("instance_status", flattenCommonCode(nasVolume.NasVolumeInstanceStatus)); err != nil {
			return err
		}
		if err := d.Set("volume_allotment_protocol_type", flattenCommonCode(nasVolume.VolumeAllotmentProtocolType)); err != nil {
			return err
		}
		if err := d.Set("zone", flattenZone(nasVolume.Zone)); err != nil {
			return err
		}
		if err := d.Set("region", flattenRegion(nasVolume.Region)); err != nil {
			return err
		}
	} else {
		log.Printf("unable to find resource: %s", d.Id())
		d.SetId("") // resource not found
	}

	return nil
}

func resourceNcloudNasVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)
	if err := deleteNasVolumeInstance(client, d.Id()); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceNcloudNasVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	if d.HasChange("volume_size_gb") {
		reqParams := new(server.ChangeNasVolumeSizeRequest)
		reqParams.NasVolumeInstanceNo = ncloud.String(d.Id())
		if volumeSizeGb, ok := d.GetOk("volume_size_gb"); ok {
			reqParams.VolumeSize = ncloud.Int32(int32(volumeSizeGb.(int)))
		}

		logCommonRequest("ChangeNasVolumeSize", reqParams)

		resp, err := client.server.V2Api.ChangeNasVolumeSize(reqParams)
		if err != nil {
			logErrorResponse("ChangeNasVolumeSize", err, reqParams)
			return err
		}
		logCommonResponse("ChangeNasVolumeSize", GetCommonResponse(resp))
	}

	if d.HasChange("server_instance_no_list") || d.HasChange("custom_ip_list") {
		reqParams := &server.SetNasVolumeAccessControlRequest{
			NasVolumeInstanceNo:  ncloud.String(d.Id()),
			ServerInstanceNoList: expandStringInterfaceList(d.Get("server_instance_no_list").([]interface{})),
			CustomIpList:         expandStringInterfaceList(d.Get("custom_ip_list").([]interface{})),
		}

		logCommonRequest("SetNasVolumeAccessControl", reqParams)

		resp, err := client.server.V2Api.SetNasVolumeAccessControl(reqParams)
		if err != nil {
			logErrorResponse("SetNasVolumeAccessControl", err, reqParams)
			return err
		}
		logCommonResponse("SetNasVolumeAccessControl", GetCommonResponse(resp))
	}

	return resourceNcloudNasVolumeRead(d, meta)
}

func buildCreateNasVolumeInstanceParams(client *NcloudAPIClient, d *schema.ResourceData) (*server.CreateNasVolumeInstanceRequest, error) {
	regionNo, err := parseRegionNoParameter(client, d)
	if err != nil {
		return nil, err
	}
	zoneNo, err := parseZoneNoParameter(client, d)
	if err != nil {
		return nil, err
	}

	reqParams := &server.CreateNasVolumeInstanceRequest{
		VolumeName:                      ncloud.String(d.Get("volume_name_postfix").(string)),
		VolumeSize:                      ncloud.Int32(int32(d.Get("volume_size_gb").(int))),
		VolumeAllotmentProtocolTypeCode: ncloud.String(d.Get("volume_allotment_protocol_type_code").(string)),
		RegionNo:                        regionNo,
		ZoneNo:                          zoneNo,
	}

	if serverInstanceNoList, ok := d.GetOk("server_instance_no_list"); ok {
		reqParams.ServerInstanceNoList = expandStringInterfaceList(serverInstanceNoList.([]interface{}))
	}

	if customIPList, ok := d.GetOk("custom_ip_list"); ok {
		reqParams.CustomIpList = expandStringInterfaceList(customIPList.([]interface{}))
	}

	if cifsUserName, ok := d.GetOk("cifs_user_name"); ok {
		reqParams.CifsUserName = ncloud.String(cifsUserName.(string))
	}

	if cifsUserPassword, ok := d.GetOk("cifs_user_password"); ok {
		reqParams.CifsUserPassword = ncloud.String(cifsUserPassword.(string))
	}

	if nasVolumeDescription, ok := d.GetOk("description"); ok {
		reqParams.NasVolumeDescription = ncloud.String(nasVolumeDescription.(string))
	}

	return reqParams, nil
}

func getNasVolumeInstance(client *NcloudAPIClient, nasVolumeInstanceNo string) (*server.NasVolumeInstance, error) {
	reqParams := &server.GetNasVolumeInstanceListRequest{
		NasVolumeInstanceNoList: []*string{ncloud.String(nasVolumeInstanceNo)},
	}

	logCommonRequest("GetNasVolumeInstanceList", reqParams)

	resp, err := client.server.V2Api.GetNasVolumeInstanceList(reqParams)
	if err != nil {
		logErrorResponse("GetNasVolumeInstanceList", err, reqParams)
		return nil, err
	}
	logCommonResponse("GetNasVolumeInstanceList", GetCommonResponse(resp))

	for _, inst := range resp.NasVolumeInstanceList {
		if nasVolumeInstanceNo == ncloud.StringValue(inst.NasVolumeInstanceNo) {
			return inst, nil
		}
	}
	return nil, nil
}

func deleteNasVolumeInstance(client *NcloudAPIClient, nasVolumeInstanceNo string) error {
	reqParams := &server.DeleteNasVolumeInstanceRequest{NasVolumeInstanceNo: ncloud.String(nasVolumeInstanceNo)}
	logCommonRequest("DeleteNasVolumeInstance", reqParams)

	resp, err := client.server.V2Api.DeleteNasVolumeInstance(reqParams)
	if err != nil {
		logErrorResponse("DeleteNasVolumeInstance", err, nasVolumeInstanceNo)
		return err
	}
	var commonResponse = &CommonResponse{}
	if resp != nil {
		commonResponse = GetCommonResponse(resp)
	}
	logCommonResponse("DeleteNasVolumeInstance", commonResponse)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"INIT", "CREAT"},
		Target:  []string{"TERMT"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getNasVolumeInstance(client, nasVolumeInstanceNo)

			if err != nil {
				return 0, "", err
			}

			if instance == nil { // Instance is terminated.
				return instance, "TERMT", nil
			}

			return instance, ncloud.StringValue(instance.NasVolumeInstanceStatus.Code), nil
		},
		Timeout:    DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for NasVolumeInstance state to be \"TERMT\": %s", err)
	}

	return nil
}
