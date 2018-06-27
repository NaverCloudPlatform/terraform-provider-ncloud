package ncloud

import (
	"fmt"
	"log"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
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
				ValidateFunc: validateStringLengthInRange(3, 30),
				Description:  "Name of a NAS volume to create. Enter a volume name that can be 3-20 characters in length after the name already entered for user identification.",
			},
			"volume_size_gb": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(500, 10000),
				Description:  "Enter the nas volume size to be created. You can enter in GB units.",
			},
			"volume_allotment_protocol_type_code": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 5),
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
				Description: "CIFS user password. The password must contain a combination of at least 2 English letters, numbers and special characters, which can be 8-14 characters in length."
			},
			"nas_volume_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(1, 1000),
				Description:  "NAS volume description",
			},
			"region_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region number.",
			},
			"zone_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Zone number. Zone in which you want to create a NAS volume.",
			},

			"volume_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "NAS volume name.",
			},
			"nas_volume_instance_status": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        commonCodeSchemaResource,
				Description: "NAS Volume instance status",
			},
			"create_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation date of the NAS volume",
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
			"nas_volume_instance_custom_ip_list": {
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
	conn := meta.(*NcloudSdk).conn

	reqParams := buildCreateNasVolumeInstanceParams(conn, d)
	resp, err := conn.CreateNasVolumeInstance(reqParams)
	if err != nil {
		logErrorResponse("CreateNasVolumeInstance", err, reqParams)
		return err
	}
	logCommonResponse("CreateNasVolumeInstance", reqParams, resp.CommonResponse)

	nasVolumeInstance := &resp.NasVolumeInstanceList[0]
	d.SetId(nasVolumeInstance.NasVolumeInstanceNo)

	if err := waitForNasVolumeInstance(conn, nasVolumeInstance.NasVolumeInstanceNo, "CREAT"); err != nil {
		return err
	}
	return resourceNcloudNasVolumeRead(d, meta)
}

func resourceNcloudNasVolumeRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	nasVolume, err := getNasVolumeInstance(conn, d.Id())
	if err != nil {
		return err
	}
	if nasVolume != nil {
		d.Set("nas_volume_instance_status", setCommonCode(nasVolume.NasVolumeInstanceStatus))
		d.Set("create_date", nasVolume.CreateDate)
		d.Set("nas_volume_description", nasVolume.NasVolumeInstanceDescription)
		d.Set("volume_allotment_protocol_type", setCommonCode(nasVolume.VolumeAllotmentProtocolType))
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
		d.Set("nas_volume_instance_custom_ip_list", nasVolume.NasVolumeInstanceCustomIPList)
		d.Set("zone", setZone(nasVolume.Zone))
		d.Set("region", setRegion(nasVolume.Region))
	}

	return nil
}

func resourceNcloudNasVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn
	return deleteNasVolumeInstance(conn, d.Id())
}

func resourceNcloudNasVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	if d.HasChange("volume_size_gb") {
		reqParams := new(sdk.RequestChangeNasVolumeSize)
		reqParams.NasVolumeInstanceNo = d.Id()
		reqParams.VolumeSize = d.Get("volume_size_gb").(int)
		resp, err := conn.ChangeNasVolumeSize(reqParams)
		if err != nil {
			logErrorResponse("ChangeNasVolumeSize", err, reqParams)
			return err
		}
		logCommonResponse("ChangeNasVolumeSize", reqParams, resp.CommonResponse)
	}

	if d.HasChange("server_instance_no_list") || d.HasChange("custom_ip_list") {
		reqParams := &sdk.RequestNasVolumeAccessControl{
			NasVolumeInstanceNo:  d.Id(),
			ServerInstanceNoList: StringList(d.Get("server_instance_no_list").([]interface{})),
			CustomIPList:         StringList(d.Get("custom_ip_list").([]interface{})),
		}

		resp, err := conn.SetNasVolumeAccessControl(reqParams)
		if err != nil {
			logErrorResponse("SetNasVolumeAccessControl", err, reqParams)
			return err
		}
		logCommonResponse("SetNasVolumeAccessControl", reqParams, resp.CommonResponse)
	}

	return resourceNcloudNasVolumeRead(d, meta)
}

func buildCreateNasVolumeInstanceParams(conn *sdk.Conn, d *schema.ResourceData) *sdk.RequestCreateNasVolumeInstance {
	reqParams := &sdk.RequestCreateNasVolumeInstance{
		VolumeName:                      d.Get("volume_name_postfix").(string),
		VolumeSize:                      d.Get("volume_size_gb").(int),
		VolumeAllotmentProtocolTypeCode: d.Get("volume_allotment_protocol_type_code").(string),
		ServerInstanceNoList:            StringList(d.Get("server_instance_no_list").([]interface{})),
		CustomIpList:                    StringList(d.Get("custom_ip_list").([]interface{})),
		CifsUserName:                    d.Get("cifs_user_name").(string),
		CifsUserPassword:                d.Get("cifs_user_password").(string),
		NasVolumeDescription:            d.Get("nas_volume_description").(string),
		RegionNo:                        parseRegionNoParameter(conn, d),
		ZoneNo:                          d.Get("zone_no").(string),
	}
	return reqParams
}

func getNasVolumeInstance(conn *sdk.Conn, nasVolumeInstanceNo string) (*sdk.NasVolumeInstance, error) {
	reqParams := &sdk.RequestGetNasVolumeInstanceList{}
	resp, err := conn.GetNasVolumeInstanceList(reqParams)
	if err != nil {
		logErrorResponse("GetNasVolumeInstanceList", err, reqParams)
		return nil, err
	}
	logCommonResponse("GetNasVolumeInstanceList", reqParams, resp.CommonResponse)

	for _, inst := range resp.NasVolumeInstanceList {
		if nasVolumeInstanceNo == inst.NasVolumeInstanceNo {
			return &inst, nil
		}
	}
	return nil, nil
}

func deleteNasVolumeInstance(conn *sdk.Conn, nasVolumeInstanceNo string) error {
	resp, err := conn.DeleteNasVolumeInstance(nasVolumeInstanceNo)
	if err != nil {
		logErrorResponse("DeleteNasVolumeInstance", err, nasVolumeInstanceNo)
		return err
	}
	var commonResponse = common.CommonResponse{}
	if resp != nil {
		commonResponse = resp.CommonResponse
	}
	logCommonResponse("DeleteNasVolumeInstance", nasVolumeInstanceNo, commonResponse)

	if err := waitForNasVolumeInstance(conn, nasVolumeInstanceNo, "TERMT"); err != nil {
		return err
	}

	return nil
}

func waitForNasVolumeInstance(conn *sdk.Conn, id string, status string) error {

	c1 := make(chan error, 1)

	go func() {
		for {
			instance, err := getNasVolumeInstance(conn, id)

			if err != nil {
				c1 <- err
				return
			}
			if instance == nil || instance.NasVolumeInstanceStatus.Code == status {
				c1 <- nil
				return
			}
			log.Printf("[DEBUG] Wait to nas volume instance (%s)", id)
			time.Sleep(time.Second * 1)
		}
	}()

	select {
	case res := <-c1:
		return res
	case <-time.After(DefaultTimeout):
		return fmt.Errorf("TIMEOUT : Wait to nas volume instance (%s)", id)
	}
}
