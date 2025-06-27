package nasvolume

import (
	"fmt"
	"log"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnas"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func ResourceNcloudNasVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudNasVolumeCreate,
		Read:   resourceNcloudNasVolumeRead,
		Update: resourceNcloudNasVolumeUpdate,
		Delete: resourceNcloudNasVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(conn.DefaultCreateTimeout),
			Delete: schema.DefaultTimeout(conn.DefaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			"volume_name_postfix": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(3, 20)),
			},
			"volume_size": {
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(500, 10000)),
			},
			"volume_allotment_protocol_type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 5)),
			},
			"server_instance_no_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
				},
			},
			"cifs_user_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cifs_user_password": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1000)),
			},
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"is_encrypted_volume": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"is_return_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"nas_volume_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_total_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"snapshot_volume_size": {
				Type:     schema.TypeInt,
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
			"mount_information": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNcloudNasVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	id, err := createNasVolume(d, config)
	if err != nil {
		return err
	}

	d.SetId(ncloud.StringValue(id))
	log.Printf("[INFO] NAS Volume ID: %s", d.Id())

	return resourceNcloudNasVolumeRead(d, meta)
}

func resourceNcloudNasVolumeRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	r, err := GetNasVolume(config, d.Id())
	if err != nil {
		return err
	}

	if r == nil {
		d.SetId("")
		return nil
	}

	instance := ConvertToMap(r)

	SetSingularResourceDataFromMapSchema(ResourceNcloudNasVolume(), d, instance)

	return nil
}

func resourceNcloudNasVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	if err := deleteNasVolume(d, config, d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceNcloudNasVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	if d.HasChange("volume_size") {
		if err := changeNasVolumeSize(d, config); err != nil {
			return err
		}
	}

	if d.HasChange("server_instance_no_list") {
		if err := setNasVolumeAccessControl(d, config); err != nil {
			return err
		}
	}

	return resourceNcloudNasVolumeRead(d, meta)
}

func GetNasVolume(config *conn.ProviderConfig, id string) (*NasVolume, error) {
	reqParams := &vnas.GetNasVolumeInstanceDetailRequest{
		RegionCode:          &config.RegionCode,
		NasVolumeInstanceNo: ncloud.String(id),
	}

	LogCommonRequest("getVpcNasVolume", reqParams)
	resp, err := config.Client.Vnas.V2Api.GetNasVolumeInstanceDetail(reqParams)
	if err != nil {
		LogErrorResponse("getVpcNasVolume", err, reqParams)
		return nil, err
	}
	LogResponse("getVpcNasVolume", resp)

	if len(resp.NasVolumeInstanceList) > 0 {
		return convertVpcNasVolume(resp.NasVolumeInstanceList[0]), nil
	}

	return nil, nil
}

func convertVpcNasVolume(inst *vnas.NasVolumeInstance) *NasVolume {
	if inst == nil {
		return nil
	}

	return &NasVolume{
		NasVolumeInstanceNo:          inst.NasVolumeInstanceNo,
		Status:                       inst.NasVolumeInstanceStatus.Code,
		NasVolumeInstanceDescription: inst.NasVolumeDescription,
		VolumeAllotmentProtocolType:  inst.VolumeAllotmentProtocolType.Code,
		VolumeName:                   inst.VolumeName,
		VolumeTotalSize:              ncloud.Int64(*inst.VolumeTotalSize / GIGABYTE),
		VolumeSize:                   ncloud.Int64(*inst.VolumeSize / GIGABYTE),
		SnapshotVolumeSize:           ncloud.Int64(*inst.SnapshotVolumeSize / GIGABYTE),
		IsSnapshotConfiguration:      inst.IsSnapshotConfiguration,
		IsEventConfiguration:         inst.IsEventConfiguration,
		Zone:                         inst.ZoneCode,
		IsEncryptedVolume:            inst.IsEncryptedVolume,
		ServerInstanceNoList:         inst.NasVolumeServerInstanceNoList,
		MountInformation:             inst.MountInformation,
		IsReturnProtection:           inst.IsReturnProtection,
	}
}

func createNasVolume(d *schema.ResourceData, config *conn.ProviderConfig) (*string, error) {
	var id *string
	var err error

	id, err = createVpcNasVolume(d, config)
	if err != nil {
		return nil, err
	}

	if err := waitForNasVolumeCreation(d, config, *id); err != nil {
		return nil, err
	}

	return id, nil
}

func createVpcNasVolume(d *schema.ResourceData, config *conn.ProviderConfig) (*string, error) {
	reqParams := &vnas.CreateNasVolumeInstanceRequest{
		RegionCode:                      &config.RegionCode,
		ZoneCode:                        StringPtrOrNil(d.GetOk("zone")),
		AccessControlRuleList:           makeVpcNasAclParams(d),
		VolumeName:                      ncloud.String(d.Get("volume_name_postfix").(string)),
		VolumeSize:                      ncloud.Int32(int32(d.Get("volume_size").(int))),
		VolumeAllotmentProtocolTypeCode: ncloud.String(d.Get("volume_allotment_protocol_type").(string)),
		CifsUserName:                    StringPtrOrNil(d.GetOk("cifs_user_name")),
		CifsUserPassword:                StringPtrOrNil(d.GetOk("cifs_user_password")),
		NasVolumeDescription:            StringPtrOrNil(d.GetOk("description")),
		IsEncryptedVolume:               BoolPtrOrNil(d.GetOk("is_encrypted_volume")),
		IsReturnProtection:              BoolPtrOrNil(d.GetOk("is_return_protection")),
	}

	resp, err := config.Client.Vnas.V2Api.CreateNasVolumeInstance(reqParams)
	if err != nil {
		LogErrorResponse("createVpcNasVolume", err, reqParams)
		return nil, err
	}
	LogResponse("createVpcNasVolume", resp)

	return resp.NasVolumeInstanceList[0].NasVolumeInstanceNo, nil
}

func waitForNasVolumeCreation(d *schema.ResourceData, config *conn.ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"INIT"},
		Target:  []string{"CREAT"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetNasVolume(config, id)

			if err != nil {
				return 0, "", err
			}

			if instance == nil {
				return instance, "INIT", nil
			}

			return instance, ncloud.StringValue(instance.Status), nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for NasVolumeInstance state to be \"CREAT\": %s", err)
	}

	return nil
}

func deleteNasVolume(d *schema.ResourceData, config *conn.ProviderConfig, id string) error {
	err := deleteVpcNasVolume(config, id)
	if err != nil {
		return err
	}

	if err := waitForNasVolumeDeletion(d, config, id); err != nil {
		return err
	}

	return nil
}

func deleteVpcNasVolume(config *conn.ProviderConfig, id string) error {
	reqParams := &vnas.DeleteNasVolumeInstancesRequest{
		RegionCode:              &config.RegionCode,
		NasVolumeInstanceNoList: []*string{ncloud.String(id)},
	}
	LogCommonRequest("deleteVpcNasVolume", reqParams)

	resp, err := config.Client.Vnas.V2Api.DeleteNasVolumeInstances(reqParams)
	if err != nil {
		LogErrorResponse("deleteVpcNasVolume", err, id)
		return err
	}
	LogResponse("deleteVpcNasVolume", resp)

	return nil
}

func waitForNasVolumeDeletion(d *schema.ResourceData, config *conn.ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"INIT", "CREAT"},
		Target:  []string{"TERMT"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetNasVolume(config, id)

			if err != nil {
				return 0, "", err
			}

			if instance == nil { // Instance is terminated.
				return instance, "TERMT", nil
			}

			return instance, ncloud.StringValue(instance.Status), nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for NasVolumeInstance state to be \"TERMT\": %s", err)
	}

	return nil
}

func changeNasVolumeSize(d *schema.ResourceData, config *conn.ProviderConfig) error {
	reqParams := &vnas.ChangeNasVolumeSizeRequest{
		RegionCode:          &config.RegionCode,
		NasVolumeInstanceNo: ncloud.String(d.Id()),
		VolumeSize:          Int32PtrOrNil(d.GetOk("volume_size")),
	}
	LogCommonRequest("changeVpcNasVolumeSize", reqParams)

	resp, err := config.Client.Vnas.V2Api.ChangeNasVolumeSize(reqParams)
	if err != nil {
		LogErrorResponse("changeVpcNasVolumeSize", err, reqParams)
		return err
	}
	LogResponse("changeVpcNasVolumeSize", resp)

	return nil
}

func setNasVolumeAccessControl(d *schema.ResourceData, config *conn.ProviderConfig) error {
	reqParams := &vnas.SetNasVolumeAccessControlRequest{
		RegionCode:            &config.RegionCode,
		NasVolumeInstanceNo:   ncloud.String(d.Id()),
		AccessControlRuleList: makeVpcNasAclParams(d),
	}

	LogCommonRequest("setVpcNasVolumeAccessControl", reqParams)

	resp, err := config.Client.Vnas.V2Api.SetNasVolumeAccessControl(reqParams)
	if err != nil {
		LogErrorResponse("setVpcNasVolumeAccessControl", err, reqParams)
		return err
	}
	LogResponse("setVpcNasVolumeAccessControl", resp)

	return nil
}

func makeVpcNasAclParams(d *schema.ResourceData) []*vnas.AccessControlRuleParameter {
	var aclParams []*vnas.AccessControlRuleParameter
	var serverList []*string

	if serverInstanceNoList, ok := d.GetOk("server_instance_no_list"); ok {
		serverList = ExpandStringInterfaceList(serverInstanceNoList.([]interface{}))

		for _, v := range serverList {
			aclParams = append(aclParams, &vnas.AccessControlRuleParameter{
				ServerInstanceNo: v,
			})
		}
	}

	return aclParams
}

// NasVolume Dto for NAS
type NasVolume struct {
	NasVolumeInstanceNo          *string   `json:"nas_volume_no,omitempty"`
	NasVolumeInstanceDescription *string   `json:"description,omitempty"`
	VolumeAllotmentProtocolType  *string   `json:"volume_allotment_protocol_type,omitempty"`
	VolumeName                   *string   `json:"name,omitempty"`
	VolumeTotalSize              *int64    `json:"volume_total_size,omitempty"`
	VolumeSize                   *int64    `json:"volume_size,omitempty"`
	SnapshotVolumeSize           *int64    `json:"snapshot_volume_size,omitempty"`
	IsSnapshotConfiguration      *bool     `json:"is_snapshot_configuration,omitempty"`
	IsEventConfiguration         *bool     `json:"is_event_configuration,omitempty"`
	Zone                         *string   `json:"zone,omitempty"`
	ServerInstanceNoList         []*string `json:"server_instance_no_list"`
	IsEncryptedVolume            *bool     `json:"is_encrypted_volume,omitempty"`
	Status                       *string   `json:"-"`
	MountInformation             *string   `json:"mount_information,omitempty"`
	IsReturnProtection           *bool     `json:"is_return_protection,omitempty"`
}
