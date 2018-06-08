package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"time"
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
			},
			"volume_size_gb": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Enter the nas volume size to be created. You can enter in GB units.",
				ValidateFunc: validateIntegerInRange(500, 10000),
			},
			"volume_allotment_protocol_type_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(1, 5),
			},
			"server_instance_no_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"custom_ip_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"cifs_user_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cifs_user_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"nas_volume_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(1, 1000),
			},
			"region_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone_no": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"volume_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nas_volume_instance_status": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_allotment_protocol_type": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"volume_total_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"volume_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"volume_use_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"volume_use_ratio": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"snapshot_volume_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"snapshot_volume_use_ratio": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"is_snapshot_configuration": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_event_configuration": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"nas_volume_instance_custom_ip_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"zone": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_no": {
							Type: schema.TypeString,
						},
						"zone_name": {
							Type: schema.TypeString,
						},
						"zone_description": {
							Type: schema.TypeString,
						},
					},
				},
			},
			"region": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_no": {
							Type: schema.TypeString,
						},
						"region_code": {
							Type: schema.TypeString,
						},
						"region_name": {
							Type: schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func resourceNcloudNasVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceNcloudNasVolumeCreate")
	conn := meta.(*NcloudSdk).conn

	reqParams := buildCreateNasVolumeInstanceParams(d)
	resp, err := conn.CreateNasVolumeInstance(reqParams)
	if err != nil {
		logErrorResponse("CreateNasVolumeInstance", err, reqParams)
		return err
	}
	logCommonResponse("CreateNasVolumeInstance", reqParams, resp.CommonResponse)

	nasVolumeInstance := &resp.NasVolumeInstanceList[0]
	d.SetId(nasVolumeInstance.NasVolumeInstanceNo)

	if err := waitForNasVolumeInstance(conn, nasVolumeInstance.NasVolumeInstanceNo, "CREAT", DefaultCreateTimeout); err != nil {
		return err
	}
	return resourceNcloudNasVolumeRead(d, meta)
}

func resourceNcloudNasVolumeRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceNcloudNasVolumeRead")
	conn := meta.(*NcloudSdk).conn

	nasVolume, err := getNasVolumeInstance(conn, d.Id())
	if err != nil {
		return err
	}
	if nasVolume != nil {
		d.Set("nas_volume_instance_status", map[string]interface{}{
			"code":      nasVolume.NasVolumeInstanceStatus.Code,
			"code_name": nasVolume.NasVolumeInstanceStatus.CodeName,
		})
		d.Set("create_date", nasVolume.CreateDate)
		d.Set("nas_volume_description", nasVolume.NasVolumeInstanceDescription)
		d.Set("volume_allotment_protocol_type", map[string]interface{}{
			"code":      nasVolume.VolumeAllotmentProtocolType.Code,
			"code_name": nasVolume.VolumeAllotmentProtocolType.CodeName,
		})
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
		d.Set("nas_volume_instance_custom_ip_list", nasVolume.NasVolumeInstanceCustomIpList)
		d.Set("zone", map[string]interface{}{
			"zone_no":          nasVolume.Zone.ZoneNo,
			"zone_name":        nasVolume.Zone.ZoneName,
			"zone_description": nasVolume.Zone.ZoneDescription,
		})
		d.Set("region", map[string]interface{}{
			"region_no":   nasVolume.Region.RegionNo,
			"region_code": nasVolume.Region.RegionCode,
			"region_name": nasVolume.Region.RegionName,
		})
	}

	return nil
}

func resourceNcloudNasVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceNcloudNasVolumeDelete")
	conn := meta.(*NcloudSdk).conn
	return deleteNasVolumeInstance(conn, d.Id())
}

func resourceNcloudNasVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceNcloudNasVolumeUpdate")
	return resourceNcloudNasVolumeRead(d, meta)
}

func buildCreateNasVolumeInstanceParams(d *schema.ResourceData) *sdk.RequestCreateNasVolumeInstance {
	reqParams := &sdk.RequestCreateNasVolumeInstance{
		VolumeName:                      d.Get("volume_name_postfix").(string),
		VolumeSize:                      d.Get("volume_size_gb").(int),
		VolumeAllotmentProtocolTypeCode: d.Get("volume_allotment_protocol_type_code").(string),
		ServerInstanceNoList:            StringList(d.Get("server_instance_no_list").([]interface{})),
		CustomIpList:                    StringList(d.Get("custom_ip_list").([]interface{})),
		CifsUserName:                    d.Get("cifs_user_name").(string),
		CifsUserPassword:                d.Get("cifs_user_password").(string),
		NasVolumeDescription:            d.Get("nas_volume_description").(string),
		RegionNo:                        d.Get("region_no").(string),
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
			log.Printf("[DEBUG] %s NasVolumeInstanceNo: %s, VolumeName: %s", "GetNasVolumeInstanceList", inst.NasVolumeInstanceNo, inst.VolumeName)
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

	if err := waitForNasVolumeInstance(conn, nasVolumeInstanceNo, "TERMT", DefaultTimeout); err != nil {
		return err
	}

	return nil
}

func waitForNasVolumeInstance(conn *sdk.Conn, id string, status string, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultWaitForInterval
	}
	for {
		instance, err := getNasVolumeInstance(conn, id)
		if err != nil {
			return err
		}
		if instance == nil || instance.NasVolumeInstanceStatus.Code == status {
			break
		}
		timeout = timeout - DefaultWaitForInterval
		if timeout <= 0 {
			return fmt.Errorf("error: Timeout: %d", timeout)
		}
		time.Sleep(DefaultWaitForInterval * time.Second)
	}
	return nil
}
