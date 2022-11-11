package ncloud

import (
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"

	"log"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	RegisterResource("ncloud_block_storage", resourceNcloudBlockStorage())
}

const (
	BlockStorageStatusCodeCreate = "CREAT"
	BlockStorageStatusCodeInit   = "INIT"
	BlockStorageStatusCodeAttach = "ATTAC"
)

func resourceNcloudBlockStorage() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudBlockStorageCreate,
		Read:   resourceNcloudBlockStorageRead,
		Update: resourceNcloudBlockStorageUpdate,
		Delete: resourceNcloudBlockStorageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DefaultCreateTimeout),
			Delete: schema.DefaultTimeout(DefaultTimeout),
		},

		Schema: map[string]*schema.Schema{
			"server_instance_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"size": {
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: ToDiagFunc(validation.IntBetween(10, 2000)),
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"disk_detail_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"snapshot_no": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"block_storage_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"device_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stop_instance_before_detaching": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceNcloudBlockStorageCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if len(d.Get("server_instance_no").(string)) == 0 {
		return fmt.Errorf("'server_instance_no' has to be present when ncloud_block_storage is first created.")
	}

	id, err := createBlockStorage(d, config)
	if err != nil {
		return err
	}

	d.SetId(ncloud.StringValue(id))
	log.Printf("[INFO] Block Storage ID: %s", d.Id())

	return resourceNcloudBlockStorageRead(d, meta)
}

func resourceNcloudBlockStorageRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	r, err := getBlockStorage(config, d.Id())
	if err != nil {
		return err
	}

	if r == nil {
		d.SetId("")
		return nil
	}

	instance := ConvertToMap(r)

	SetSingularResourceDataFromMapSchema(resourceNcloudBlockStorage(), d, instance)

	if err := d.Set("server_instance_no", r.ServerInstanceNo); err != nil {
		return err
	}

	return nil
}

func resourceNcloudBlockStorageDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if d.Get("stop_instance_before_detaching").(bool) {
		log.Printf("[INFO] Stopping Instance %s for destroying block storage", d.Get("server_instance_no").(string))
		if err := stopThenWaitServerInstance(config, d.Get("server_instance_no").(string)); err != nil {
			return err
		}
	}

	if err := deleteBlockStorage(d, config, d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceNcloudBlockStorageUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if d.HasChange("server_instance_no") {
		o, n := d.GetChange("server_instance_no")

		// If server instance attached block storage, detach first
		if len(o.(string)) > 0 {
			if d.Get("stop_instance_before_detaching").(bool) {
				log.Printf("[INFO] Start Instance %s after detaching block storage", o.(string))
				if err := stopThenWaitServerInstance(config, o.(string)); err != nil {
					return err
				}
			}

			if err := detachBlockStorage(config, d.Id()); err != nil {
				return err
			}
		}

		if len(n.(string)) > 0 {
			if err := attachBlockStorage(d, config); err != nil {
				return err
			}
		}
	}

	if d.HasChange("size") {
		o, n := d.GetChange("size")

		if o.(int) >= n.(int) {
			return fmt.Errorf("The storage size is only expandable, not shrinking. new size(%d) must be greater than the existing size(%d)", n, o)
		}

		// If server instance attached block storage, detach first
		if len(d.Get("server_instance_no").(string)) > 0 {
			if d.Get("stop_instance_before_detaching").(bool) {
				log.Printf("[INFO] Start Instance %s after detaching block storage", d.Get("server_instance_no").(string))
				if err := stopThenWaitServerInstance(config, d.Get("server_instance_no").(string)); err != nil {
					return err
				}
			}

			if err := detachBlockStorage(config, d.Id()); err != nil {
				return err
			}
		}

		if err := changeBlockStorageSize(d, config); err != nil {
			return err
		}

		if len(d.Get("server_instance_no").(string)) > 0 {
			if err := attachBlockStorage(d, config); err != nil {
				return err
			}
		}
	}

	return resourceNcloudBlockStorageRead(d, meta)
}

func createBlockStorage(d *schema.ResourceData, config *ProviderConfig) (*string, error) {
	var id *string
	var err error

	if config.SupportVPC {
		id, err = createVpcBlockStorage(d, config)
	} else {
		id, err = createClassicBlockStorage(d, config)
	}

	if err != nil {
		return nil, err
	}

	if err := waitForBlockStorageAttachment(config, *id); err != nil {
		return nil, err
	}

	return id, nil
}

func createClassicBlockStorage(d *schema.ResourceData, config *ProviderConfig) (*string, error) {
	reqParams := &server.CreateBlockStorageInstanceRequest{
		ServerInstanceNo:        ncloud.String(d.Get("server_instance_no").(string)),
		BlockStorageSize:        ncloud.Int64(int64(d.Get("size").(int))),
		BlockStorageName:        StringPtrOrNil(d.GetOk("name")),
		BlockStorageDescription: StringPtrOrNil(d.GetOk("description")),
		DiskDetailTypeCode:      StringPtrOrNil(d.GetOk("disk_detail_type")),
	}

	logCommonRequest("createClassicBlockStorage", reqParams)

	resp, err := config.Client.server.V2Api.CreateBlockStorageInstance(reqParams)
	if err != nil {
		logErrorResponse("createClassicBlockStorage", err, reqParams)
		return nil, err
	}
	logResponse("createClassicBlockStorage", resp)

	instance := resp.BlockStorageInstanceList[0]

	return instance.BlockStorageInstanceNo, nil
}

func createVpcBlockStorage(d *schema.ResourceData, config *ProviderConfig) (*string, error) {
	reqParams := &vserver.CreateBlockStorageInstanceRequest{
		RegionCode:                     &config.RegionCode,
		BlockStorageSize:               ncloud.Int32(int32(d.Get("size").(int))),
		ServerInstanceNo:               ncloud.String(d.Get("server_instance_no").(string)),
		BlockStorageName:               StringPtrOrNil(d.GetOk("name")),
		BlockStorageDescription:        StringPtrOrNil(d.GetOk("description")),
		BlockStorageDiskDetailTypeCode: StringPtrOrNil(d.GetOk("disk_detail_type")),
		BlockStorageSnapshotInstanceNo: StringPtrOrNil(d.GetOk("snapshot_no")),
		ZoneCode:                       StringPtrOrNil(d.GetOk("zone")),
	}

	logCommonRequest("createVpcBlockStorage", reqParams)

	resp, err := config.Client.vserver.V2Api.CreateBlockStorageInstance(reqParams)
	if err != nil {
		logErrorResponse("createVpcBlockStorage", err, reqParams)
		return nil, err
	}
	logResponse("createVpcBlockStorage", resp)

	instance := resp.BlockStorageInstanceList[0]

	return instance.BlockStorageInstanceNo, nil
}

func getBlockStorage(config *ProviderConfig, id string) (*BlockStorage, error) {
	if config.SupportVPC {
		return getVpcBlockStorage(config, id)
	}

	return getClassicBlockStorage(config, id)
}

func getClassicBlockStorage(config *ProviderConfig, id string) (*BlockStorage, error) {
	reqParams := &server.GetBlockStorageInstanceListRequest{
		BlockStorageInstanceNoList: ncloud.StringList([]string{id}),
	}

	logCommonRequest("getClassicBlockStorage", reqParams)

	resp, err := config.Client.server.V2Api.GetBlockStorageInstanceList(reqParams)
	if err != nil {
		logErrorResponse("getClassicBlockStorage", err, reqParams)
		return nil, err
	}
	logResponse("getClassicBlockStorage", resp)

	if len(resp.BlockStorageInstanceList) > 0 {
		inst := resp.BlockStorageInstanceList[0]

		return &BlockStorage{
			BlockStorageInstanceNo:  inst.BlockStorageInstanceNo,
			ServerInstanceNo:        inst.ServerInstanceNo,
			ServerName:              inst.ServerName,
			BlockStorageType:        inst.BlockStorageType.Code,
			BlockStorageName:        inst.BlockStorageName,
			BlockStorageSize:        ncloud.Int64(*inst.BlockStorageSize / GIGABYTE),
			DeviceName:              inst.DeviceName,
			BlockStorageProductCode: inst.BlockStorageProductCode,
			Status:                  inst.BlockStorageInstanceStatus.Code,
			Operation:               inst.BlockStorageInstanceOperation.Code,
			Description:             inst.BlockStorageInstanceDescription,
			DiskType:                inst.DiskType.Code,
			DiskDetailType:          inst.DiskDetailType.Code,
		}, nil
	}

	return nil, nil
}

func getVpcBlockStorage(config *ProviderConfig, id string) (*BlockStorage, error) {
	reqParams := &vserver.GetBlockStorageInstanceDetailRequest{
		RegionCode:             &config.RegionCode,
		BlockStorageInstanceNo: ncloud.String(id),
	}

	logCommonRequest("getVpcBlockStorage", reqParams)

	resp, err := config.Client.vserver.V2Api.GetBlockStorageInstanceDetail(reqParams)
	if err != nil {
		logErrorResponse("getVpcBlockStorage", err, reqParams)
		return nil, err
	}
	logResponse("getVpcBlockStorage", resp)

	if len(resp.BlockStorageInstanceList) > 0 {
		inst := resp.BlockStorageInstanceList[0]

		return &BlockStorage{
			BlockStorageInstanceNo:  inst.BlockStorageInstanceNo,
			ServerInstanceNo:        inst.ServerInstanceNo,
			BlockStorageType:        inst.BlockStorageType.Code,
			BlockStorageName:        inst.BlockStorageName,
			BlockStorageSize:        ncloud.Int64(*inst.BlockStorageSize / GIGABYTE),
			DeviceName:              inst.DeviceName,
			BlockStorageProductCode: inst.BlockStorageProductCode,
			Status:                  inst.BlockStorageInstanceStatus.Code,
			Operation:               inst.BlockStorageInstanceOperation.Code,
			Description:             inst.BlockStorageDescription,
			DiskType:                inst.BlockStorageDiskType.Code,
			DiskDetailType:          inst.BlockStorageDiskDetailType.Code,
			ZoneCode:                inst.ZoneCode,
		}, nil
	}

	return nil, nil
}

func deleteBlockStorage(d *schema.ResourceData, config *ProviderConfig, id string) error {

	var err error

	if config.SupportVPC {
		err = deleteVpcBlockStorage(d, config, id)
	} else {
		err = deleteClassicBlockStorage(d, config, id)
	}

	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{BlockStorageStatusCodeInit, BlockStorageStatusCodeAttach},
		Target:  []string{"TERMINATED"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getBlockStorage(config, id)
			if err != nil {
				return 0, "", err
			}
			if instance == nil { // Instance is terminated.
				return instance, "TERMINATED", nil
			}
			return instance, ncloud.StringValue(instance.Status), nil
		},
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for BlockStorageInstance state to be \"TERMINATED\": %s", err)
	}

	return nil
}

func deleteClassicBlockStorage(d *schema.ResourceData, config *ProviderConfig, id string) error {
	reqParams := server.DeleteBlockStorageInstancesRequest{
		BlockStorageInstanceNoList: []*string{ncloud.String(id)},
	}

	resp, err := config.Client.server.V2Api.DeleteBlockStorageInstances(&reqParams)

	if err != nil {
		logErrorResponse("deleteClassicBlockStorage", err, reqParams)
		return err
	}
	logResponse("deleteClassicBlockStorage", resp)

	return nil
}

func deleteVpcBlockStorage(d *schema.ResourceData, config *ProviderConfig, id string) error {
	reqParams := vserver.DeleteBlockStorageInstancesRequest{
		RegionCode:                 &config.RegionCode,
		BlockStorageInstanceNoList: []*string{ncloud.String(id)},
	}

	resp, err := config.Client.vserver.V2Api.DeleteBlockStorageInstances(&reqParams)

	if err != nil {
		logErrorResponse("deleteVpcBlockStorage", err, reqParams)
		return err
	}
	logResponse("deleteVpcBlockStorage", resp)

	return nil
}

func detachBlockStorage(config *ProviderConfig, id string) error {
	var err error

	if config.SupportVPC {
		err = detachVpcBlockStorage(config, id)
	} else {
		err = detachClassicBlockStorage(config, id)
	}

	if err != nil {
		return err
	}

	if err = waitForBlockStorageDetachment(config, id); err != nil {
		return err
	}

	return nil
}

func detachClassicBlockStorage(config *ProviderConfig, id string) error {
	reqParams := &server.DetachBlockStorageInstancesRequest{
		BlockStorageInstanceNoList: []*string{ncloud.String(id)},
	}

	logCommonRequest("detachClassicBlockStorage", reqParams)

	resp, err := config.Client.server.V2Api.DetachBlockStorageInstances(reqParams)
	if err != nil {
		logErrorResponse("detachClassicBlockStorage", err, reqParams)
		return err
	}
	logResponse("detachClassicBlockStorage", resp)

	return nil
}

func detachVpcBlockStorage(config *ProviderConfig, id string) error {
	reqParams := &vserver.DetachBlockStorageInstancesRequest{
		BlockStorageInstanceNoList: []*string{ncloud.String(id)},
	}

	logCommonRequest("detachVpcBlockStorage", reqParams)

	resp, err := config.Client.vserver.V2Api.DetachBlockStorageInstances(reqParams)
	if err != nil {
		logErrorResponse("detachVpcBlockStorage", err, reqParams)
		return err
	}
	logResponse("detachVpcBlockStorage", resp)

	return nil
}

func waitForBlockStorageDetachment(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{BlockStorageStatusCodeAttach},
		Target:  []string{BlockStorageStatusCodeCreate},
		Refresh: func() (interface{}, string, error) {
			instance, err := getBlockStorage(config, id)
			if err != nil {
				return 0, "", err
			}
			return instance, ncloud.StringValue(instance.Status), nil
		},
		Timeout:    DefaultUpdateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for BlockStorageInstance state to be \"CREAT\": %s", err)
	}

	return nil
}

func attachBlockStorage(d *schema.ResourceData, config *ProviderConfig) error {
	var err error
	if config.SupportVPC {
		err = attachVpcBlockStorage(d, config)
	} else {
		err = attachClassicBlockStorage(d, config)
	}

	if err != nil {
		return err
	}

	if err = waitForBlockStorageAttachment(config, d.Id()); err != nil {
		return err
	}

	return nil
}

func attachClassicBlockStorage(d *schema.ResourceData, config *ProviderConfig) error {
	reqParams := &server.AttachBlockStorageInstanceRequest{
		ServerInstanceNo:       ncloud.String(d.Get("server_instance_no").(string)),
		BlockStorageInstanceNo: ncloud.String(d.Id()),
	}

	logCommonRequest("attachClassicBlockStorage", reqParams)

	resp, err := config.Client.server.V2Api.AttachBlockStorageInstance(reqParams)
	if err != nil {
		logErrorResponse("attachClassicBlockStorage", err, reqParams)
		return err
	}
	logResponse("attachClassicBlockStorage", resp)

	return nil
}

func attachVpcBlockStorage(d *schema.ResourceData, config *ProviderConfig) error {
	reqParams := &vserver.AttachBlockStorageInstanceRequest{
		ServerInstanceNo:       ncloud.String(d.Get("server_instance_no").(string)),
		BlockStorageInstanceNo: ncloud.String(d.Id()),
	}

	logCommonRequest("attachVpcBlockStorage", reqParams)

	resp, err := config.Client.vserver.V2Api.AttachBlockStorageInstance(reqParams)
	if err != nil {
		logErrorResponse("attachVpcBlockStorage", err, reqParams)
		return err
	}
	logResponse("attachVpcBlockStorage", resp)

	return nil
}

func waitForBlockStorageAttachment(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{BlockStorageStatusCodeInit, BlockStorageStatusCodeCreate},
		Target:  []string{BlockStorageStatusCodeAttach},
		Refresh: func() (interface{}, string, error) {
			instance, err := getBlockStorage(config, id)
			if err != nil {
				return 0, "", err
			}
			return instance, ncloud.StringValue(instance.Status), nil
		},
		Timeout:    DefaultUpdateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for BlockStorageInstance state to be \"ATTAC\": %s", err)
	}

	return nil
}

func changeBlockStorageSize(d *schema.ResourceData, config *ProviderConfig) error {
	var err error
	if config.SupportVPC {
		err = changeVpcBlockStorageSize(d, config)
	} else {
		err = changeClassicBlockStorageSize(d, config)
	}

	if err != nil {
		return err
	}

	if err = waitForBlockStorageOperationIsNull(config, d.Id()); err != nil {
		return err
	}

	return nil
}

func changeVpcBlockStorageSize(d *schema.ResourceData, config *ProviderConfig) error {
	reqParams := &vserver.ChangeBlockStorageVolumeSizeRequest{
		RegionCode:             &config.RegionCode,
		BlockStorageInstanceNo: ncloud.String(d.Id()),
		BlockStorageSize:       ncloud.Int32(int32(d.Get("size").(int))),
	}

	logCommonRequest("changeVpcBlockStorageSize", reqParams)
	resp, err := config.Client.vserver.V2Api.ChangeBlockStorageVolumeSize(reqParams)
	if err != nil {
		logErrorResponse("changeVpcBlockStorageSize", err, reqParams)
		return err
	}
	logResponse("changeVpcBlockStorageSize", resp)

	return nil
}

func changeClassicBlockStorageSize(d *schema.ResourceData, config *ProviderConfig) error {
	reqParams := &server.ChangeBlockStorageVolumeSizeRequest{
		BlockStorageInstanceNo: ncloud.String(d.Id()),
		BlockStorageSize:       ncloud.Int64(int64(d.Get("size").(int))),
	}

	logCommonRequest("changeClassicBlockStorageSize", reqParams)
	resp, err := config.Client.server.V2Api.ChangeBlockStorageVolumeSize(reqParams)
	if err != nil {
		logErrorResponse("changeClassicBlockStorageSize", err, reqParams)
		return err
	}
	logResponse("changeClassicBlockStorageSize", resp)

	return nil
}

func waitForBlockStorageOperationIsNull(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"CHNG"},
		Target:  []string{"NULL"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getBlockStorage(config, id)
			if err != nil {
				return 0, "", err
			}
			return instance, ncloud.StringValue(instance.Operation), nil
		},
		Timeout:    DefaultUpdateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for BlockStorageInstance operation to be \"NULL\": %s", err)
	}

	return nil
}

// BlockStorage Dto for block storage
type BlockStorage struct {
	BlockStorageInstanceNo  *string `json:"block_storage_no,omitempty"`
	ServerInstanceNo        *string `json:"server_instance_no,omitempty"`
	ServerName              *string `json:"server_name,omitempty"`
	BlockStorageType        *string `json:"type,omitempty"`
	BlockStorageName        *string `json:"name,omitempty"`
	BlockStorageSize        *int64  `json:"size,omitempty"`
	DeviceName              *string `json:"device_name,omitempty"`
	BlockStorageProductCode *string `json:"product_code,omitempty"`
	Status                  *string `json:"status,omitempty"`
	Operation               *string `json:"operation,omitempty"`
	Description             *string `json:"description,omitempty"`
	DiskType                *string `json:"disk_type,omitempty"`
	DiskDetailType          *string `json:"disk_detail_type,omitempty"`
	ZoneCode                *string `json:"zone,omitempty"`
}
