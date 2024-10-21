package server

import (
	"fmt"
	"regexp"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"

	"log"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

const (
	BlockStorageStatusCodeCreate     = "CREAT"
	BlockStorageStatusCodeInit       = "INIT"
	BlockStorageStatusCodeAttach     = "ATTAC"
	BlockStorageStatusNameInit       = "initialized"
	BlockStorageStatusNameCreating   = "creating"
	BlockStorageStatusNameOptimizing = "optimizing"
	BlockStorageStatusNameAttaching  = "attaching"
	BlockStorageStatusNameAttach     = "attached"
	BlockStorageStatusNameDetach     = "detached"
	BlockStorageVolumeTypeHdd        = "HDD"
	BlockStorageVolumeTypeSsd        = "SSD"
	BlockStorageVolumeTypeFb1        = "FB1"
	BlockStorageVolumeTypeCb1        = "CB1"
	BlockStorageHypervisorTypeXen    = "XEN"
	BlockStorageHypervisorTypeKvm    = "KVM"
)

func ResourceNcloudBlockStorage() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudBlockStorageCreate,
		Read:   resourceNcloudBlockStorageRead,
		Update: resourceNcloudBlockStorageUpdate,
		Delete: resourceNcloudBlockStorageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(conn.DefaultCreateTimeout),
			Delete: schema.DefaultTimeout(conn.DefaultTimeout),
		},

		Schema: map[string]*schema.Schema{
			"server_instance_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"size": {
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(10)),
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(
					validation.StringLenBetween(3, 30),
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z][A-Za-z0-9-_]+$`), "Allows only alphabets, numbers, hyphen (-) and underbar (_). Must start with an alphabetic character"),
				)),
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 1000)),
			},
			"disk_detail_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ConflictsWith:    []string{"volume_type"},
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{BlockStorageVolumeTypeHdd, BlockStorageVolumeTypeSsd}, false)),
			},
			"hypervisor_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				RequiredWith:     []string{"volume_type"},
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{BlockStorageHypervisorTypeXen, BlockStorageHypervisorTypeKvm}, false)),
			},
			"volume_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ConflictsWith:    []string{"disk_detail_type"},
				RequiredWith:     []string{"hypervisor_type"},
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{BlockStorageVolumeTypeHdd, BlockStorageVolumeTypeSsd, BlockStorageVolumeTypeFb1, BlockStorageVolumeTypeCb1}, false)),
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
			"stop_instance_before_detaching": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"return_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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
			"max_iops": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"encrypted_volume": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceNcloudBlockStorageCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

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
	config := meta.(*conn.ProviderConfig)

	r, err := GetBlockStorage(config, d.Id())
	if err != nil {
		return err
	}

	if r == nil {
		d.SetId("")
		return nil
	}

	instance := ConvertToMap(r)

	SetSingularResourceDataFromMapSchema(ResourceNcloudBlockStorage(), d, instance)

	if err := d.Set("server_instance_no", r.ServerInstanceNo); err != nil {
		return err
	}

	return nil
}

func resourceNcloudBlockStorageDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

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
	config := meta.(*conn.ProviderConfig)

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

			if err := detachThenWaitServerInstance(config, o.(string)); err != nil {
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

			if err := detachThenWaitServerInstance(config, d.Get("server_instance_no").(string)); err != nil {
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

	if d.HasChange("return_protection") {
		if !config.SupportVPC {
			return fmt.Errorf("`return_protection` only available in VPC environments")
		}

		if err := changeVpcBlockStorageReturnProtection(d, config); err != nil {
			return err
		}
	}

	return resourceNcloudBlockStorageRead(d, meta)
}

func createBlockStorage(d *schema.ResourceData, config *conn.ProviderConfig) (*string, error) {
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

	return id, nil
}

func createClassicBlockStorage(d *schema.ResourceData, config *conn.ProviderConfig) (*string, error) {
	reqParams := &server.CreateBlockStorageInstanceRequest{
		ServerInstanceNo:        ncloud.String(d.Get("server_instance_no").(string)),
		BlockStorageSize:        ncloud.Int64(int64(d.Get("size").(int))),
		BlockStorageName:        StringPtrOrNil(d.GetOk("name")),
		BlockStorageDescription: StringPtrOrNil(d.GetOk("description")),
		DiskDetailTypeCode:      StringPtrOrNil(d.GetOk("disk_detail_type")),
	}

	LogCommonRequest("createClassicBlockStorage", reqParams)

	resp, err := config.Client.Server.V2Api.CreateBlockStorageInstance(reqParams)
	if err != nil {
		LogErrorResponse("createClassicBlockStorage", err, reqParams)
		return nil, err
	}
	LogResponse("createClassicBlockStorage", resp)

	instance := resp.BlockStorageInstanceList[0]
	if err := waitForBlockStorageAttachment(config, *instance.BlockStorageInstanceNo); err != nil {
		return nil, err
	}

	return instance.BlockStorageInstanceNo, nil
}

func createVpcBlockStorage(d *schema.ResourceData, config *conn.ProviderConfig) (*string, error) {
	reqParams := &vserver.CreateBlockStorageInstanceRequest{
		RegionCode:                     &config.RegionCode,
		BlockStorageSize:               ncloud.Int32(int32(d.Get("size").(int))),
		BlockStorageName:               StringPtrOrNil(d.GetOk("name")),
		BlockStorageDescription:        StringPtrOrNil(d.GetOk("description")),
		BlockStorageDiskDetailTypeCode: StringPtrOrNil(d.GetOk("disk_detail_type")),
		BlockStorageSnapshotInstanceNo: StringPtrOrNil(d.GetOk("snapshot_no")),
		BlockStorageVolumeTypeCode:     StringPtrOrNil(d.GetOk("volume_type")),
		ZoneCode:                       StringPtrOrNil(d.GetOk("zone")),
		IsReturnProtection:             BoolPtrOrNil(d.GetOk("return_protection")),
	}

	hypervisorType, hypervisorTypeOk := d.GetOk("hypervisor_type")
	_, diskTypeOk := d.GetOk("disk_detail_type")
	volumeType := d.Get("volume_type").(string)

	if (!hypervisorTypeOk && !diskTypeOk) || diskTypeOk || (hypervisorType == BlockStorageHypervisorTypeXen) {
		reqParams.ServerInstanceNo = ncloud.String(d.Get("server_instance_no").(string))
	}

	if (hypervisorType == BlockStorageHypervisorTypeXen) && ((volumeType == BlockStorageVolumeTypeFb1) || (volumeType == BlockStorageVolumeTypeCb1)) {
		err := fmt.Errorf("Only `%s` and `%s` can be entered as `%s` hypervisor type", BlockStorageVolumeTypeSsd, BlockStorageVolumeTypeHdd, BlockStorageHypervisorTypeXen)
		LogErrorResponse("createVpcBlockStorage", err, reqParams)
		return nil, err
	}

	if (hypervisorType == BlockStorageHypervisorTypeKvm) && ((volumeType == BlockStorageVolumeTypeHdd) || (volumeType == BlockStorageVolumeTypeSsd)) {
		err := fmt.Errorf("Only `%s` and `%s` can be entered as `%s` hypervisor type", BlockStorageVolumeTypeCb1, BlockStorageVolumeTypeFb1, BlockStorageHypervisorTypeKvm)
		LogErrorResponse("createVpcBlockStorage", err, reqParams)
		return nil, err
	}

	if hypervisorType == BlockStorageHypervisorTypeKvm {
		zone := d.Get("zone").(string)
		if len(zone) == 0 {
			err := fmt.Errorf("`zone` is required for KVM type")
			LogErrorResponse("createVpcBlockStorage", err, reqParams)
			return nil, err
		}

		server, err := GetServerInstance(config, d.Get("server_instance_no").(string))
		if err == nil && server == nil {
			err = fmt.Errorf("fail to get serverInstance")
		}
		if err != nil {
			LogErrorResponse("createVpcBlockStorage", err, reqParams)
			return nil, err
		}

		if *server.Zone != zone {
			err := fmt.Errorf("Different from the server's zone code %s", *server.Zone)
			LogErrorResponse("createVpcBlockStorage", err, reqParams)
			return nil, err
		}
	}

	LogCommonRequest("createVpcBlockStorage", reqParams)

	resp, err := config.Client.Vserver.V2Api.CreateBlockStorageInstance(reqParams)
	if err != nil {
		LogErrorResponse("createVpcBlockStorage", err, reqParams)
		return nil, err
	}
	LogResponse("createVpcBlockStorage", resp)

	if resp == nil || len(resp.BlockStorageInstanceList) < 1 {
		err := fmt.Errorf("response invalid")
		LogErrorResponse("createVpcBlockStorage", err, reqParams)
		return nil, err
	}

	instance := resp.BlockStorageInstanceList[0]
	output, err := waitForBlockStorageCreation(config, *instance.BlockStorageInstanceNo)
	if err != nil {
		LogErrorResponse("createVpcBlockStorage", err, reqParams)
		return nil, err
	}

	if *output.StatusName == BlockStorageStatusNameDetach {
		d.SetId(*instance.BlockStorageInstanceNo)
		if err := attachBlockStorage(d, config); err != nil {
			return nil, err
		}
	}

	return instance.BlockStorageInstanceNo, nil
}

func GetBlockStorage(config *conn.ProviderConfig, id string) (*BlockStorage, error) {
	if config.SupportVPC {
		return getVpcBlockStorage(config, id)
	}

	return getClassicBlockStorage(config, id)
}

func getClassicBlockStorage(config *conn.ProviderConfig, id string) (*BlockStorage, error) {
	reqParams := &server.GetBlockStorageInstanceListRequest{
		BlockStorageInstanceNoList: ncloud.StringList([]string{id}),
	}

	LogCommonRequest("getClassicBlockStorage", reqParams)

	resp, err := config.Client.Server.V2Api.GetBlockStorageInstanceList(reqParams)
	if err != nil {
		LogErrorResponse("getClassicBlockStorage", err, reqParams)
		return nil, err
	}
	LogResponse("getClassicBlockStorage", resp)

	if len(resp.BlockStorageInstanceList) > 0 {
		inst := resp.BlockStorageInstanceList[0]

		return &BlockStorage{
			BlockStorageInstanceNo:  inst.BlockStorageInstanceNo,
			ServerInstanceNo:        inst.ServerInstanceNo,
			ServerName:              inst.ServerName,
			BlockStorageType:        common.GetCodePtrByCommonCode(inst.BlockStorageType),
			BlockStorageName:        inst.BlockStorageName,
			BlockStorageSize:        ncloud.Int64(*inst.BlockStorageSize / GIGABYTE),
			DeviceName:              inst.DeviceName,
			BlockStorageProductCode: inst.BlockStorageProductCode,
			Status:                  common.GetCodePtrByCommonCode(inst.BlockStorageInstanceStatus),
			Operation:               common.GetCodePtrByCommonCode(inst.BlockStorageInstanceOperation),
			Description:             inst.BlockStorageInstanceDescription,
			DiskType:                common.GetCodePtrByCommonCode(inst.DiskType),
			DiskDetailType:          common.GetCodePtrByCommonCode(inst.DiskDetailType),
		}, nil
	}

	return nil, nil
}

func getVpcBlockStorage(config *conn.ProviderConfig, id string) (*BlockStorage, error) {
	reqParams := &vserver.GetBlockStorageInstanceDetailRequest{
		RegionCode:             &config.RegionCode,
		BlockStorageInstanceNo: ncloud.String(id),
	}

	LogCommonRequest("getVpcBlockStorage", reqParams)

	resp, err := config.Client.Vserver.V2Api.GetBlockStorageInstanceDetail(reqParams)
	if err != nil {
		LogErrorResponse("getVpcBlockStorage", err, reqParams)
		return nil, err
	}
	LogResponse("getVpcBlockStorage", resp)

	if len(resp.BlockStorageInstanceList) > 0 {
		inst := resp.BlockStorageInstanceList[0]

		blockStorage := BlockStorage{
			BlockStorageInstanceNo:  inst.BlockStorageInstanceNo,
			ServerInstanceNo:        inst.ServerInstanceNo,
			BlockStorageType:        inst.BlockStorageType.Code,
			BlockStorageName:        inst.BlockStorageName,
			BlockStorageSize:        ncloud.Int64(*inst.BlockStorageSize / GIGABYTE),
			DeviceName:              inst.DeviceName,
			BlockStorageProductCode: inst.BlockStorageProductCode,
			Status:                  common.GetCodePtrByCommonCode(inst.BlockStorageInstanceStatus),
			Operation:               common.GetCodePtrByCommonCode(inst.BlockStorageInstanceOperation),
			StatusName:              inst.BlockStorageInstanceStatusName,
			Description:             inst.BlockStorageDescription,
			DiskType:                common.GetCodePtrByCommonCode(inst.BlockStorageDiskType),
			DiskDetailType:          common.GetCodePtrByCommonCode(inst.BlockStorageDiskDetailType),
			ZoneCode:                inst.ZoneCode,
			MaxIops:                 inst.MaxIopsThroughput,
			EncryptedVolume:         inst.IsEncryptedVolume,
			ReturnProtection:        inst.IsReturnProtection,
			VolumeType:              common.GetCodePtrByCommonCode(inst.BlockStorageVolumeType),
			HypervisorType:          common.GetCodePtrByCommonCode(inst.HypervisorType),
		}

		return &blockStorage, nil
	}

	return nil, nil
}

func deleteBlockStorage(d *schema.ResourceData, config *conn.ProviderConfig, id string) error {

	var err error

	if config.SupportVPC {
		err = deleteVpcBlockStorage(d, config, id)
	} else {
		err = deleteClassicBlockStorage(d, config, id)
	}

	if err != nil {
		return err
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{BlockStorageStatusCodeCreate, BlockStorageStatusCodeInit, BlockStorageStatusCodeAttach},
		Target:  []string{"TERMINATED"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetBlockStorage(config, id)
			if err != nil {
				return 0, "", err
			}
			if instance == nil { // Instance is terminated.
				return instance, "TERMINATED", nil
			}
			return instance, ncloud.StringValue(instance.Status), nil
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for BlockStorageInstance state to be \"TERMINATED\": %s", err)
	}

	return nil
}

func deleteClassicBlockStorage(d *schema.ResourceData, config *conn.ProviderConfig, id string) error {
	reqParams := server.DeleteBlockStorageInstancesRequest{
		BlockStorageInstanceNoList: []*string{ncloud.String(id)},
	}

	LogCommonRequest("deleteClassicBlockStorage", reqParams)

	resp, err := config.Client.Server.V2Api.DeleteBlockStorageInstances(&reqParams)

	if err != nil {
		LogErrorResponse("deleteClassicBlockStorage", err, reqParams)
		return err
	}
	LogResponse("deleteClassicBlockStorage", resp)

	return nil
}

func deleteVpcBlockStorage(d *schema.ResourceData, config *conn.ProviderConfig, id string) error {
	reqParams := vserver.DeleteBlockStorageInstancesRequest{
		RegionCode:                 &config.RegionCode,
		BlockStorageInstanceNoList: []*string{ncloud.String(id)},
	}

	LogCommonRequest("deleteVpcBlockStorage", reqParams)

	resp, err := config.Client.Vserver.V2Api.DeleteBlockStorageInstances(&reqParams)

	if err != nil {
		LogErrorResponse("deleteVpcBlockStorage", err, reqParams)
		return err
	}
	LogResponse("deleteVpcBlockStorage", resp)

	return nil
}

func detachBlockStorage(config *conn.ProviderConfig, id string) error {
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

func detachClassicBlockStorage(config *conn.ProviderConfig, id string) error {
	reqParams := &server.DetachBlockStorageInstancesRequest{
		BlockStorageInstanceNoList: []*string{ncloud.String(id)},
	}

	LogCommonRequest("detachClassicBlockStorage", reqParams)

	resp, err := config.Client.Server.V2Api.DetachBlockStorageInstances(reqParams)
	if err != nil {
		LogErrorResponse("detachClassicBlockStorage", err, reqParams)
		return err
	}
	LogResponse("detachClassicBlockStorage", resp)

	return nil
}

func detachVpcBlockStorage(config *conn.ProviderConfig, id string) error {
	reqParams := &vserver.DetachBlockStorageInstancesRequest{
		BlockStorageInstanceNoList: []*string{ncloud.String(id)},
	}

	LogCommonRequest("detachVpcBlockStorage", reqParams)

	resp, err := config.Client.Vserver.V2Api.DetachBlockStorageInstances(reqParams)
	if err != nil {
		LogErrorResponse("detachVpcBlockStorage", err, reqParams)
		return err
	}
	LogResponse("detachVpcBlockStorage", resp)

	return nil
}

func waitForBlockStorageDetachment(config *conn.ProviderConfig, id string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{BlockStorageStatusCodeAttach},
		Target:  []string{BlockStorageStatusCodeCreate},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetBlockStorage(config, id)
			if err != nil {
				return 0, "", err
			}

			if instance == nil {
				return 0, "", fmt.Errorf("GetBlockStorage is nil")
			}
			return instance, ncloud.StringValue(instance.Status), nil
		},
		Timeout:    conn.DefaultUpdateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for BlockStorageInstance state to be \"CREAT\": %s", err)
	}

	return nil
}

func attachBlockStorage(d *schema.ResourceData, config *conn.ProviderConfig) error {
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

func attachClassicBlockStorage(d *schema.ResourceData, config *conn.ProviderConfig) error {
	reqParams := &server.AttachBlockStorageInstanceRequest{
		ServerInstanceNo:       ncloud.String(d.Get("server_instance_no").(string)),
		BlockStorageInstanceNo: ncloud.String(d.Id()),
	}

	LogCommonRequest("attachClassicBlockStorage", reqParams)

	resp, err := config.Client.Server.V2Api.AttachBlockStorageInstance(reqParams)
	if err != nil {
		LogErrorResponse("attachClassicBlockStorage", err, reqParams)
		return err
	}
	LogResponse("attachClassicBlockStorage", resp)

	return nil
}

func attachVpcBlockStorage(d *schema.ResourceData, config *conn.ProviderConfig) error {
	reqParams := &vserver.AttachBlockStorageInstanceRequest{
		ServerInstanceNo:       ncloud.String(d.Get("server_instance_no").(string)),
		BlockStorageInstanceNo: ncloud.String(d.Id()),
	}

	LogCommonRequest("attachVpcBlockStorage", reqParams)

	resp, err := config.Client.Vserver.V2Api.AttachBlockStorageInstance(reqParams)
	if err != nil {
		LogErrorResponse("attachVpcBlockStorage", err, reqParams)
		return err
	}
	LogResponse("attachVpcBlockStorage", resp)

	return nil
}

func waitForBlockStorageCreation(config *conn.ProviderConfig, id string) (*BlockStorage, error) {
	var blockStorageInstance *BlockStorage
	stateConf := &retry.StateChangeConf{
		Pending: []string{BlockStorageStatusNameInit, BlockStorageStatusNameCreating, BlockStorageStatusNameAttaching},
		Target:  []string{BlockStorageStatusNameAttach, BlockStorageStatusNameDetach},
		Refresh: func() (interface{}, string, error) {
			resp, err := GetBlockStorage(config, id)
			if err != nil {
				return 0, "", err
			}

			if resp == nil {
				return 0, "", fmt.Errorf("GetBlockStorage is nil")
			}
			blockStorageInstance = resp

			return resp, *resp.StatusName, nil
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return nil, fmt.Errorf("error waiting for BlockStorageInstance create: %s", err)
	}

	return blockStorageInstance, nil
}

func waitForBlockStorageAttachment(config *conn.ProviderConfig, id string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{BlockStorageStatusCodeInit, BlockStorageStatusCodeCreate},
		Target:  []string{BlockStorageStatusCodeAttach},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetBlockStorage(config, id)
			if err != nil {
				return 0, "", err
			}

			if instance == nil {
				return 0, "", fmt.Errorf("GetBlockStorage is nil")
			}
			return instance, ncloud.StringValue(instance.Status), nil
		},
		Timeout:    conn.DefaultUpdateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for BlockStorageInstance state to be \"ATTAC\": %s", err)
	}

	return nil
}

func changeBlockStorageSize(d *schema.ResourceData, config *conn.ProviderConfig) error {
	var err error
	if config.SupportVPC {
		if d.Get("hypervisor_type").(string) == BlockStorageHypervisorTypeXen {
			err = changeVpcBlockStorageVolumeSize(d, config)
		} else {
			err = changeVpcBlockStorageInstance(d, config)
		}
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

func changeVpcBlockStorageVolumeSize(d *schema.ResourceData, config *conn.ProviderConfig) error {
	reqParams := &vserver.ChangeBlockStorageVolumeSizeRequest{
		RegionCode:             &config.RegionCode,
		BlockStorageInstanceNo: ncloud.String(d.Id()),
		BlockStorageSize:       ncloud.Int32(int32(d.Get("size").(int))),
	}

	LogCommonRequest("changeVpcBlockStorageVolumeSize", reqParams)
	resp, err := config.Client.Vserver.V2Api.ChangeBlockStorageVolumeSize(reqParams)
	if err != nil {
		LogErrorResponse("changeVpcBlockStorageVolumeSize", err, reqParams)
		return err
	}
	LogResponse("changeVpcBlockStorageVolumeSize", resp)

	return nil
}

func changeVpcBlockStorageInstance(d *schema.ResourceData, config *conn.ProviderConfig) error {
	reqParams := &vserver.ChangeBlockStorageInstanceRequest{
		RegionCode:             &config.RegionCode,
		BlockStorageInstanceNo: ncloud.String(d.Id()),
		BlockStorageSize:       ncloud.Int32(int32(d.Get("size").(int))),
	}

	LogCommonRequest("changeVpcBlockStorageInstance", reqParams)
	resp, err := config.Client.Vserver.V2Api.ChangeBlockStorageInstance(reqParams)
	if err != nil {
		LogErrorResponse("changeVpcBlockStorageInstance", err, reqParams)
		return err
	}
	LogResponse("changeVpcBlockStorageInstance", resp)

	return nil
}

func changeClassicBlockStorageSize(d *schema.ResourceData, config *conn.ProviderConfig) error {
	reqParams := &server.ChangeBlockStorageVolumeSizeRequest{
		BlockStorageInstanceNo: ncloud.String(d.Id()),
		BlockStorageSize:       ncloud.Int64(int64(d.Get("size").(int))),
	}

	LogCommonRequest("changeClassicBlockStorageSize", reqParams)
	resp, err := config.Client.Server.V2Api.ChangeBlockStorageVolumeSize(reqParams)
	if err != nil {
		LogErrorResponse("changeClassicBlockStorageSize", err, reqParams)
		return err
	}
	LogResponse("changeClassicBlockStorageSize", resp)

	return nil
}

func waitForBlockStorageOperationIsNull(config *conn.ProviderConfig, id string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"CHNG"},
		Target:  []string{"NULL"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetBlockStorage(config, id)
			if err != nil {
				return 0, "", err
			}

			if instance == nil {
				return 0, "", fmt.Errorf("GetBlockStorage is nil")
			}
			return instance, ncloud.StringValue(instance.Operation), nil
		},
		Timeout:    conn.DefaultUpdateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for BlockStorageInstance operation to be \"NULL\": %s", err)
	}

	return nil
}

func changeVpcBlockStorageReturnProtection(d *schema.ResourceData, config *conn.ProviderConfig) error {
	reqParams := &vserver.SetBlockStorageReturnProtectionRequest{
		RegionCode:             &config.RegionCode,
		BlockStorageInstanceNo: ncloud.String(d.Id()),
		IsReturnProtection:     ncloud.Bool(d.Get("return_protection").(bool)),
	}

	LogCommonRequest("changeVpcBlockStorageReturnProtection", reqParams)
	resp, err := config.Client.Vserver.V2Api.SetBlockStorageReturnProtection(reqParams)
	if err != nil {
		LogErrorResponse("changeVpcBlockStorageReturnProtection", err, reqParams)
		return err
	}
	LogResponse("changeVpcBlockStorageReturnProtection", resp)

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
	StatusName              *string `json:"status_name,omitempty"`
	Description             *string `json:"description,omitempty"`
	DiskType                *string `json:"disk_type,omitempty"`
	DiskDetailType          *string `json:"disk_detail_type,omitempty"`
	ZoneCode                *string `json:"zone,omitempty"`
	MaxIops                 *int32  `json:"max_iops,omitempty"`
	EncryptedVolume         *bool   `json:"encrypted_volume,omitempty"`
	ReturnProtection        *bool   `json:"return_protection,omitempty"`
	VolumeType              *string `json:"volume_type,omitempty"`
	HypervisorType          *string `json:"hypervisor_type,omitempty"`
}
