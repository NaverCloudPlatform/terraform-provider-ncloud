package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnas"
	"log"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	RegisterResource("ncloud_nas_volume", resourceNcloudNasVolume())
}

func resourceNcloudNasVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudNasVolumeCreate,
		Read:   resourceNcloudNasVolumeRead,
		Update: resourceNcloudNasVolumeUpdate,
		Delete: resourceNcloudNasVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DefaultCreateTimeout),
			Delete: schema.DefaultTimeout(DefaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			"volume_name_postfix": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(3, 30)),
			},
			"volume_size": {
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: ToDiagFunc(validation.IntBetween(100, 10000)),
			},
			"volume_allotment_protocol_type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(1, 5)),
			},
			"server_instance_no_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"custom_ip_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
				ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(1, 1000)),
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
	config := meta.(*ProviderConfig)

	id, err := createNasVolume(d, config)
	if err != nil {
		return err
	}

	d.SetId(ncloud.StringValue(id))
	log.Printf("[INFO] NAS Volume ID: %s", d.Id())

	return resourceNcloudNasVolumeRead(d, meta)
}

func resourceNcloudNasVolumeRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	r, err := getNasVolume(config, d.Id())
	if err != nil {
		return err
	}

	if r == nil {
		d.SetId("")
		return nil
	}

	instance := ConvertToMap(r)

	SetSingularResourceDataFromMapSchema(resourceNcloudNasVolume(), d, instance)

	return nil
}

func resourceNcloudNasVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if err := deleteNasVolume(d, config, d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceNcloudNasVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if d.HasChange("volume_size") {
		if err := changeNasVolumeSize(d, config); err != nil {
			return err
		}
	}

	if d.HasChange("server_instance_no_list") || d.HasChange("custom_ip_list") {
		if err := setNasVolumeAccessControl(d, config); err != nil {
			return err
		}
	}

	return resourceNcloudNasVolumeRead(d, meta)
}

func getNasVolume(config *ProviderConfig, id string) (*NasVolume, error) {
	if config.SupportVPC {
		return getVpcNasVolume(config, id)
	} else {
		return getClassicNasVolume(config, id)
	}
}

func getClassicNasVolume(config *ProviderConfig, id string) (*NasVolume, error) {
	reqParams := &server.GetNasVolumeInstanceListRequest{
		NasVolumeInstanceNoList: []*string{ncloud.String(id)},
	}

	logCommonRequest("getClassicNasVolume", reqParams)

	resp, err := config.Client.server.V2Api.GetNasVolumeInstanceList(reqParams)
	if err != nil {
		logErrorResponse("getClassicNasVolume", err, reqParams)
		return nil, err
	}
	logResponse("getClassicNasVolume", resp)

	if len(resp.NasVolumeInstanceList) > 0 {
		return convertClassicNasVolume(resp.NasVolumeInstanceList[0]), nil
	}

	return nil, nil
}

func convertClassicNasVolume(inst *server.NasVolumeInstance) *NasVolume {
	if inst == nil {
		return nil
	}

	return &NasVolume{
		NasVolumeInstanceNo:           inst.NasVolumeInstanceNo,
		Status:                        inst.NasVolumeInstanceStatus.Code,
		NasVolumeInstanceDescription:  inst.NasVolumeInstanceDescription,
		VolumeAllotmentProtocolType:   inst.VolumeAllotmentProtocolType.Code,
		VolumeName:                    inst.VolumeName,
		VolumeTotalSize:               ncloud.Int64(*inst.VolumeTotalSize / GIGABYTE),
		VolumeSize:                    ncloud.Int64(*inst.VolumeSize / GIGABYTE),
		SnapshotVolumeSize:            ncloud.Int64(*inst.SnapshotVolumeSize / GIGABYTE),
		IsSnapshotConfiguration:       inst.IsSnapshotConfiguration,
		IsEventConfiguration:          inst.IsEventConfiguration,
		Zone:                          inst.Zone.ZoneCode,
		NasVolumeInstanceCustomIpList: flattenArrayStructByKey(inst.NasVolumeInstanceCustomIpList, "customIp"),
		ServerInstanceNoList:          flattenArrayStructByKey(inst.NasVolumeServerInstanceList, "serverInstanceNo"),
		MountInformation:              inst.MountInformation,
	}
}

func getVpcNasVolume(config *ProviderConfig, id string) (*NasVolume, error) {
	reqParams := &vnas.GetNasVolumeInstanceDetailRequest{
		RegionCode:          &config.RegionCode,
		NasVolumeInstanceNo: ncloud.String(id),
	}

	logCommonRequest("getVpcNasVolume", reqParams)
	resp, err := config.Client.vnas.V2Api.GetNasVolumeInstanceDetail(reqParams)
	if err != nil {
		logErrorResponse("getVpcNasVolume", err, reqParams)
		return nil, err
	}
	logResponse("getVpcNasVolume", resp)

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
		NasVolumeInstanceNo:           inst.NasVolumeInstanceNo,
		Status:                        inst.NasVolumeInstanceStatus.Code,
		NasVolumeInstanceDescription:  inst.NasVolumeDescription,
		VolumeAllotmentProtocolType:   inst.VolumeAllotmentProtocolType.Code,
		VolumeName:                    inst.VolumeName,
		VolumeTotalSize:               ncloud.Int64(*inst.VolumeTotalSize / GIGABYTE),
		VolumeSize:                    ncloud.Int64(*inst.VolumeSize / GIGABYTE),
		SnapshotVolumeSize:            ncloud.Int64(*inst.SnapshotVolumeSize / GIGABYTE),
		IsSnapshotConfiguration:       inst.IsSnapshotConfiguration,
		IsEventConfiguration:          inst.IsEventConfiguration,
		Zone:                          inst.ZoneCode,
		IsEncryptedVolume:             inst.IsEncryptedVolume,
		ServerInstanceNoList:          inst.NasVolumeServerInstanceNoList,
		NasVolumeInstanceCustomIpList: []*string{},
		MountInformation:              inst.MountInformation,
	}
}

func createNasVolume(d *schema.ResourceData, config *ProviderConfig) (*string, error) {
	var id *string
	var err error

	if config.SupportVPC {
		id, err = createVpcNasVolume(d, config)
	} else {
		id, err = createClassicNasVolume(d, config)
	}

	if err != nil {
		return nil, err
	}

	if err := waitForNasVolumeCreation(d, config, *id); err != nil {
		return nil, err
	}

	return id, nil
}

func createClassicNasVolume(d *schema.ResourceData, config *ProviderConfig) (*string, error) {
	regionNo, err := parseRegionNoParameter(d)
	if err != nil {
		return nil, err
	}
	zoneNo, err := parseZoneNoParameter(config, d)
	if err != nil {
		return nil, err
	}

	reqParams := &server.CreateNasVolumeInstanceRequest{
		RegionNo:                        regionNo,
		ZoneNo:                          zoneNo,
		VolumeName:                      ncloud.String(d.Get("volume_name_postfix").(string)),
		VolumeSize:                      ncloud.Int32(int32(d.Get("volume_size").(int))),
		VolumeAllotmentProtocolTypeCode: ncloud.String(d.Get("volume_allotment_protocol_type").(string)),
		CifsUserName:                    StringPtrOrNil(d.GetOk("cifs_user_name")),
		CifsUserPassword:                StringPtrOrNil(d.GetOk("cifs_user_password")),
		NasVolumeDescription:            StringPtrOrNil(d.GetOk("description")),
	}

	if serverInstanceNoList, ok := d.GetOk("server_instance_no_list"); ok {
		reqParams.ServerInstanceNoList = expandStringInterfaceList(serverInstanceNoList.([]interface{}))
	}

	if customIPList, ok := d.GetOk("custom_ip_list"); ok {
		reqParams.CustomIpList = expandStringInterfaceList(customIPList.([]interface{}))
	}

	logCommonRequest("createClassicNasVolume", reqParams)

	resp, err := config.Client.server.V2Api.CreateNasVolumeInstance(reqParams)
	if err != nil {
		logErrorResponse("createClassicNasVolume", err, reqParams)
		return nil, err
	}
	logResponse("createClassicNasVolume", resp)

	return resp.NasVolumeInstanceList[0].NasVolumeInstanceNo, nil
}

func createVpcNasVolume(d *schema.ResourceData, config *ProviderConfig) (*string, error) {
	reqParams := &vnas.CreateNasVolumeInstanceRequest{
		RegionCode:                      &config.RegionCode,
		ZoneCode:                        StringPtrOrNil(d.GetOk("zone")),
		VolumeName:                      ncloud.String(d.Get("volume_name_postfix").(string)),
		VolumeSize:                      ncloud.Int32(int32(d.Get("volume_size").(int))),
		VolumeAllotmentProtocolTypeCode: ncloud.String(d.Get("volume_allotment_protocol_type").(string)),
		CifsUserName:                    StringPtrOrNil(d.GetOk("cifs_user_name")),
		CifsUserPassword:                StringPtrOrNil(d.GetOk("cifs_user_password")),
		NasVolumeDescription:            StringPtrOrNil(d.GetOk("description")),
		IsEncryptedVolume:               BoolPtrOrNil(d.GetOk("is_encrypted_volume")),
	}

	if serverInstanceNoList, ok := d.GetOk("server_instance_no_list"); ok {
		reqParams.ServerInstanceNoList = expandStringInterfaceList(serverInstanceNoList.([]interface{}))
	}

	logCommonRequest("createVpcNasVolume", reqParams)

	resp, err := config.Client.vnas.V2Api.CreateNasVolumeInstance(reqParams)
	if err != nil {
		logErrorResponse("createVpcNasVolume", err, reqParams)
		return nil, err
	}
	logResponse("createVpcNasVolume", resp)

	return resp.NasVolumeInstanceList[0].NasVolumeInstanceNo, nil
}

func waitForNasVolumeCreation(d *schema.ResourceData, config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"INIT"},
		Target:  []string{"CREAT"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getNasVolume(config, id)

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

func deleteNasVolume(d *schema.ResourceData, config *ProviderConfig, id string) error {
	var err error

	if config.SupportVPC {
		err = deleteVpcNasVolume(config, id)
	} else {
		err = deleteClassicNasVolume(config, id)
	}

	if err != nil {
		return err
	}

	if err := waitForNasVolumeDeletion(d, config, id); err != nil {
		return err
	}

	return nil
}

func deleteClassicNasVolume(config *ProviderConfig, id string) error {
	reqParams := &server.DeleteNasVolumeInstanceRequest{NasVolumeInstanceNo: ncloud.String(id)}
	logCommonRequest("deleteClassicNasVolume", reqParams)

	resp, err := config.Client.server.V2Api.DeleteNasVolumeInstance(reqParams)
	if err != nil {
		logErrorResponse("deleteClassicNasVolume", err, id)
		return err
	}
	logResponse("deleteClassicNasVolume", resp)

	return nil
}

func deleteVpcNasVolume(config *ProviderConfig, id string) error {
	reqParams := &vnas.DeleteNasVolumeInstancesRequest{
		RegionCode:              &config.RegionCode,
		NasVolumeInstanceNoList: []*string{ncloud.String(id)},
	}
	logCommonRequest("deleteVpcNasVolume", reqParams)

	resp, err := config.Client.vnas.V2Api.DeleteNasVolumeInstances(reqParams)
	if err != nil {
		logErrorResponse("deleteVpcNasVolume", err, id)
		return err
	}
	logResponse("deleteVpcNasVolume", resp)

	return nil
}

func waitForNasVolumeDeletion(d *schema.ResourceData, config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"INIT", "CREAT"},
		Target:  []string{"TERMT"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getNasVolume(config, id)

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

func changeNasVolumeSize(d *schema.ResourceData, config *ProviderConfig) error {
	if config.SupportVPC {
		return changeVpcNasVolumeSize(d, config)
	} else {
		return changeClassicNasVolumeSize(d, config)
	}
}

func changeClassicNasVolumeSize(d *schema.ResourceData, config *ProviderConfig) error {
	reqParams := &server.ChangeNasVolumeSizeRequest{
		NasVolumeInstanceNo: ncloud.String(d.Id()),
		VolumeSize:          Int32PtrOrNil(d.GetOk("volume_size")),
	}
	logCommonRequest("changeClassicNasVolumeSize", reqParams)

	resp, err := config.Client.server.V2Api.ChangeNasVolumeSize(reqParams)
	if err != nil {
		logErrorResponse("changeClassicNasVolumeSize", err, reqParams)
		return err
	}
	logResponse("changeClassicNasVolumeSize", resp)

	return nil
}

func changeVpcNasVolumeSize(d *schema.ResourceData, config *ProviderConfig) error {
	reqParams := &vnas.ChangeNasVolumeSizeRequest{
		NasVolumeInstanceNo: ncloud.String(d.Id()),
		VolumeSize:          Int32PtrOrNil(d.GetOk("volume_size")),
	}
	logCommonRequest("changeVpcNasVolumeSize", reqParams)

	resp, err := config.Client.vnas.V2Api.ChangeNasVolumeSize(reqParams)
	if err != nil {
		logErrorResponse("changeVpcNasVolumeSize", err, reqParams)
		return err
	}
	logResponse("changeVpcNasVolumeSize", resp)

	return nil
}

func setNasVolumeAccessControl(d *schema.ResourceData, config *ProviderConfig) error {
	if config.SupportVPC {
		if d.HasChange("server_instance_no_list") {
			return setVpcNasVolumeAccessControl(d, config)
		}
		return nil
	} else {
		return setClassicNasVolumeAccessControl(d, config)
	}
}

func setClassicNasVolumeAccessControl(d *schema.ResourceData, config *ProviderConfig) error {
	reqParams := &server.SetNasVolumeAccessControlRequest{
		NasVolumeInstanceNo:  ncloud.String(d.Id()),
		ServerInstanceNoList: expandStringInterfaceList(d.Get("server_instance_no_list").([]interface{})),
		CustomIpList:         expandStringInterfaceList(d.Get("custom_ip_list").([]interface{})),
	}

	logCommonRequest("setClassicNasVolumeAccessControl", reqParams)

	resp, err := config.Client.server.V2Api.SetNasVolumeAccessControl(reqParams)
	if err != nil {
		logErrorResponse("setClassicNasVolumeAccessControl", err, reqParams)
		return err
	}
	logResponse("setClassicNasVolumeAccessControl", resp)

	return nil
}

func setVpcNasVolumeAccessControl(d *schema.ResourceData, config *ProviderConfig) error {
	reqParams := &vnas.SetNasVolumeAccessControlRequest{
		NasVolumeInstanceNo:  ncloud.String(d.Id()),
		ServerInstanceNoList: expandStringInterfaceList(d.Get("server_instance_no_list").([]interface{})),
	}

	logCommonRequest("setVpcNasVolumeAccessControl", reqParams)

	resp, err := config.Client.vnas.V2Api.SetNasVolumeAccessControl(reqParams)
	if err != nil {
		logErrorResponse("setVpcNasVolumeAccessControl", err, reqParams)
		return err
	}
	logResponse("setVpcNasVolumeAccessControl", resp)

	return nil
}

// NasVolume Dto for NAS
type NasVolume struct {
	NasVolumeInstanceNo           *string   `json:"nas_volume_no,omitempty"`
	NasVolumeInstanceDescription  *string   `json:"description,omitempty"`
	VolumeAllotmentProtocolType   *string   `json:"volume_allotment_protocol_type,omitempty"`
	VolumeName                    *string   `json:"name,omitempty"`
	VolumeTotalSize               *int64    `json:"volume_total_size,omitempty"`
	VolumeSize                    *int64    `json:"volume_size,omitempty"`
	SnapshotVolumeSize            *int64    `json:"snapshot_volume_size,omitempty"`
	IsSnapshotConfiguration       *bool     `json:"is_snapshot_configuration,omitempty"`
	IsEventConfiguration          *bool     `json:"is_event_configuration,omitempty"`
	Zone                          *string   `json:"zone,omitempty"`
	NasVolumeInstanceCustomIpList []*string `json:"custom_ip_list"`
	ServerInstanceNoList          []*string `json:"server_instance_no_list"`
	IsEncryptedVolume             *bool     `json:"is_encrypted_volume,omitempty"`
	Status                        *string   `json:"-"`
	MountInformation              *string   `json:"mount_information,omitempty"`
}
