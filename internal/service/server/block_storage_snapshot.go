package server

import (
	"fmt"
	"time"

	"log"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func ResourceNcloudBlockStorageSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudBlockStorageSnapshotCreate,
		Read:   resourceNcloudBlockStorageSnapshotRead,
		Update: resourceNcloudBlockStorageSnapshotUpdate,
		Delete: resourceNcloudBlockStorageSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(conn.DefaultCreateTimeout),
			Delete: schema.DefaultTimeout(conn.DefaultTimeout),
		},

		Schema: map[string]*schema.Schema{
			"block_storage_instance_no": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Block storage instance No for creating snapshot.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Block storage snapshot name to create. default : Ncloud assigns default values.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Descriptions on a snapshot to create",
			},

			"instance_no": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Block Storage Snapshot Instance Number",
			},
			"volume_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Block Storage Snapshot Volume Size",
			},
			"original_block_storage_instance_no": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Original Block Storage Instance Number",
			},
			"original_block_storage_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Original Block Storage Name",
			},
			"instance_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Block Storage Snapshot Instance Status",
			},
			"instance_operation": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Block Storage Snapshot Instance Operation",
			},
			"instance_status_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Block Storage Snapshot Instance Status Name",
			},
			"server_image_product_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server Image Product Code",
			},
			"os_information": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OS Information",
			},
			"hypervisor_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Hypervisor Type",
			},
		},
	}
}

func resourceNcloudBlockStorageSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	var err error
	config := meta.(*conn.ProviderConfig)

	if config.SupportVPC {
		err = createVpcBlockStorageSnapshot(d, config)
	} else {
		err = createClassicBlockStorageSnapshot(d, config)
	}

	if err != nil {
		return err
	}

	return resourceNcloudBlockStorageSnapshotRead(d, meta)
}

func resourceNcloudBlockStorageSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var r *BlockStorageSnapshot
	config := meta.(*conn.ProviderConfig)

	if config.SupportVPC {
		r, err = GetVpcBlockStorageSnapshotDetail(config, d.Id())
	} else {
		r, err = GetClassicBlockStorageSnapshotInstance(config, d.Id())
	}

	if err != nil {
		return err
	}

	if r == nil {
		log.Printf("unable to find resource: %s", d.Id())
		d.SetId("")
		return nil
	}

	d.SetId(*r.BlockStorageSnapshotInstanceNo)
	d.Set("block_storage_instance_no", *r.BlockStorageInstanceNo)
	d.Set("name", *r.BlockStorageSnapshotName)
	d.Set("description", *r.Description)
	d.Set("instance_no", *r.BlockStorageSnapshotInstanceNo)
	d.Set("volume_size", *r.BlockStorageSnapshotVolumeSize)
	d.Set("instance_status", *r.Status)
	d.Set("instance_operation", *r.Operation)
	d.Set("instance_status_name", *r.StatusName)

	if r.OriginalBlockStorageInstanceNo != nil {
		d.Set("original_block_storage_instance_no", *r.OriginalBlockStorageInstanceNo)
	}
	if r.OriginalBlockStorageName != nil {
		d.Set("original_block_storage_name", *r.OriginalBlockStorageName)
	}
	if r.ServerImageProductCode != nil {
		d.Set("server_image_product_code", *r.ServerImageProductCode)
	}
	if r.OsInformation != nil {
		d.Set("os_information", *r.OsInformation)
	}
	if r.HypervisorType != nil {
		d.Set("hypervisor_type", *r.HypervisorType)
	}

	return nil
}

func resourceNcloudBlockStorageSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceNcloudBlockStorageSnapshotRead(d, meta)
}

func resourceNcloudBlockStorageSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	var err error
	config := meta.(*conn.ProviderConfig)

	if config.SupportVPC {
		err = deleteVpcBlockStorageSnapshot(config, d.Id())
	} else {
		err = deleteClassicBlockStorageSnapshot(config, d.Id())
	}

	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func createVpcBlockStorageSnapshot(d *schema.ResourceData, config *conn.ProviderConfig) error {
	reqParams := &vserver.CreateBlockStorageSnapshotInstanceRequest{
		RegionCode:                      &config.RegionCode,
		OriginalBlockStorageInstanceNo:  ncloud.String(d.Get("block_storage_instance_no").(string)),
		BlockStorageSnapshotName:        StringPtrOrNil(d.GetOk("name")),
		BlockStorageSnapshotDescription: StringPtrOrNil(d.GetOk("description")),
	}

	LogCommonRequest("createVpcBlockStorageSnapshot", reqParams)

	resp, err := config.Client.Vserver.V2Api.CreateBlockStorageSnapshotInstance(reqParams)
	if err != nil {
		LogErrorResponse("createVpcBlockStorageSnapshot", err, reqParams)
		return err
	}
	LogResponse("createVpcBlockStorageSnapshot", resp)

	if resp == nil || len(resp.BlockStorageSnapshotInstanceList) < 1 {
		err := fmt.Errorf("response invalid")
		LogErrorResponse("createVpcBlockStorageSnapshot", err, reqParams)
		return err
	}

	instance := resp.BlockStorageSnapshotInstanceList[0]
	err = waitForBlockStorageSnapshotCreation(config, *instance.BlockStorageSnapshotInstanceNo)
	if err != nil {
		LogErrorResponse("createVpcBlockStorageSnapshot", err, reqParams)
		return err
	}
	d.SetId(*instance.BlockStorageSnapshotInstanceNo)

	return nil
}

func waitForBlockStorageSnapshotCreation(config *conn.ProviderConfig, id string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"INIT"},
		Target:  []string{"CREAT"},
		Refresh: func() (interface{}, string, error) {
			resp, err := GetVpcBlockStorageSnapshotDetail(config, id)
			if err != nil {
				return 0, "", err
			}

			if resp == nil {
				return 0, "", fmt.Errorf("GetVpcBlockStorageSnapshotDetail is nil")
			}

			return resp, *resp.Status, nil
		},
		Timeout:    conn.DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for BlockStorageSnapshotInstance state to be \"CREAT\": %s", err)
	}

	return nil
}

func createClassicBlockStorageSnapshot(d *schema.ResourceData, config *conn.ProviderConfig) error {
	reqParams := buildRequestBlockStorageSnapshotInstance(d)
	LogCommonRequest("createClassicBlockStorageSnapshot", reqParams)

	resp, err := config.Client.Server.V2Api.CreateBlockStorageSnapshotInstance(reqParams)
	if err != nil {
		LogErrorResponse("createClassicBlockStorageSnapshot", err, reqParams)
		return err
	}
	LogResponse("createClassicBlockStorageSnapshot", resp)

	blockStorageSnapshotInstance := resp.BlockStorageSnapshotInstanceList[0]
	blockStorageSnapshotInstanceNo := ncloud.StringValue(blockStorageSnapshotInstance.BlockStorageSnapshotInstanceNo)

	stateConf := &retry.StateChangeConf{
		Pending: []string{"INIT"},
		Target:  []string{"CREAT"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetClassicBlockStorageSnapshotInstance(config, blockStorageSnapshotInstanceNo)
			if err != nil {
				return 0, "", err
			}
			return instance, *instance.Status, nil
		},
		Timeout:    conn.DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for BlockStorageSnapshotInstance state to be \"CREAT\": %s", err)
	}
	d.SetId(blockStorageSnapshotInstanceNo)

	return nil
}

func buildRequestBlockStorageSnapshotInstance(d *schema.ResourceData) *server.CreateBlockStorageSnapshotInstanceRequest {
	reqParams := &server.CreateBlockStorageSnapshotInstanceRequest{
		BlockStorageInstanceNo: ncloud.String(d.Get("block_storage_instance_no").(string)),
	}

	if blockStorageSnapshotName, ok := d.GetOk("name"); ok {
		reqParams.BlockStorageSnapshotName = ncloud.String(blockStorageSnapshotName.(string))
	}

	if blockStorageSnapshotDescription, ok := d.GetOk("description"); ok {
		reqParams.BlockStorageSnapshotDescription = ncloud.String(blockStorageSnapshotDescription.(string))
	}

	return reqParams
}

func GetClassicBlockStorageSnapshotInstance(config *conn.ProviderConfig, blockStorageSnapshotInstanceNo string) (*BlockStorageSnapshot, error) {
	reqParams := &server.GetBlockStorageSnapshotInstanceListRequest{
		BlockStorageSnapshotInstanceNoList: []*string{ncloud.String(blockStorageSnapshotInstanceNo)},
	}

	LogCommonRequest("getClassicBlockStorageSnapshotInstanceList", reqParams)

	resp, err := config.Client.Server.V2Api.GetBlockStorageSnapshotInstanceList(reqParams)
	if err != nil {
		LogErrorResponse("getClassicBlockStorageSnapshotInstanceList", err, reqParams)
		return nil, err
	}
	LogResponse("getClassicBlockStorageSnapshotInstanceList", resp)

	if len(resp.BlockStorageSnapshotInstanceList) > 0 {
		inst := resp.BlockStorageSnapshotInstanceList[0]

		blockStorageSnapshot := BlockStorageSnapshot{
			BlockStorageInstanceNo:         inst.OriginalBlockStorageInstanceNo,
			BlockStorageSnapshotName:       inst.BlockStorageSnapshotName,
			Description:                    inst.BlockStorageSnapshotInstanceDescription,
			BlockStorageSnapshotInstanceNo: inst.BlockStorageSnapshotInstanceNo,
			BlockStorageSnapshotVolumeSize: inst.BlockStorageSnapshotVolumeSize,
			Status:                         inst.BlockStorageSnapshotInstanceStatus.Code,
			Operation:                      inst.BlockStorageSnapshotInstanceOperation.Code,
			StatusName:                     inst.BlockStorageSnapshotInstanceStatusName,
			OriginalBlockStorageInstanceNo: inst.OriginalBlockStorageInstanceNo,
			OriginalBlockStorageName:       inst.OriginalBlockStorageName,
			ServerImageProductCode:         inst.ServerImageProductCode,
			OsInformation:                  inst.OsInformation,
		}

		return &blockStorageSnapshot, nil
	}

	return nil, nil
}

func GetVpcBlockStorageSnapshotDetail(config *conn.ProviderConfig, blockStorageSnapshotInstanceNo string) (*BlockStorageSnapshot, error) {
	reqParams := &vserver.GetBlockStorageSnapshotInstanceDetailRequest{
		BlockStorageSnapshotInstanceNo: ncloud.String(blockStorageSnapshotInstanceNo),
	}

	LogCommonRequest("GetVpcBlockStorageSnapshotDetail", reqParams)

	resp, err := config.Client.Vserver.V2Api.GetBlockStorageSnapshotInstanceDetail(reqParams)
	if err != nil {
		LogErrorResponse("GetVpcBlockStorageSnapshotDetail", err, reqParams)
		return nil, err
	}
	LogResponse("GetVpcBlockStorageSnapshotDetail", resp)

	if len(resp.BlockStorageSnapshotInstanceList) > 0 {
		inst := resp.BlockStorageSnapshotInstanceList[0]

		blockStorageSnapshot := BlockStorageSnapshot{
			BlockStorageInstanceNo:         inst.OriginalBlockStorageInstanceNo,
			BlockStorageSnapshotName:       inst.BlockStorageSnapshotName,
			Description:                    inst.BlockStorageSnapshotDescription,
			BlockStorageSnapshotInstanceNo: inst.BlockStorageSnapshotInstanceNo,
			BlockStorageSnapshotVolumeSize: inst.BlockStorageSnapshotVolumeSize,
			Status:                         inst.BlockStorageSnapshotInstanceStatus.Code,
			Operation:                      inst.BlockStorageSnapshotInstanceOperation.Code,
			StatusName:                     inst.BlockStorageSnapshotInstanceStatusName,
			HypervisorType:                 inst.HypervisorType.Code,
		}

		return &blockStorageSnapshot, nil
	}

	return nil, nil
}

func deleteVpcBlockStorageSnapshot(config *conn.ProviderConfig, id string) error {
	reqParams := &vserver.DeleteBlockStorageSnapshotInstancesRequest{
		RegionCode:                         &config.RegionCode,
		BlockStorageSnapshotInstanceNoList: []*string{ncloud.String(id)},
	}

	LogCommonRequest("deleteVpcBlockStorageSnapshot", reqParams)

	resp, err := config.Client.Vserver.V2Api.DeleteBlockStorageSnapshotInstances(reqParams)
	if err != nil {
		LogErrorResponse("deleteVpcBlockStorageSnapshot", err, reqParams)
		return err
	}
	LogResponse("deleteVpcBlockStorageSnapshot", resp)

	err = waitForBlockStorageSnapshotDelete(config, id)
	if err != nil {
		LogErrorResponse("deleteVpcBlockStorageSnapshot", err, reqParams)
		return err
	}

	return nil
}

func waitForBlockStorageSnapshotDelete(config *conn.ProviderConfig, id string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREAT"},
		Target:  []string{"TERMINATED"},
		Refresh: func() (interface{}, string, error) {
			resp, err := GetVpcBlockStorageSnapshotDetail(config, id)
			if err != nil {
				return 0, "", err
			}

			if resp == nil {
				return resp, "TERMINATED", nil
			}

			return resp, *resp.Status, nil
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for BlockStorageSnapshotInstance state to be \"TERMINATED\": %s", err)
	}

	return nil
}

func deleteClassicBlockStorageSnapshot(config *conn.ProviderConfig, id string) error {
	reqParams := server.DeleteBlockStorageSnapshotInstancesRequest{
		BlockStorageSnapshotInstanceNoList: []*string{ncloud.String(id)},
	}

	LogCommonRequest("DeleteBlockStorageSnapshotInstances", reqParams)

	resp, err := config.Client.Server.V2Api.DeleteBlockStorageSnapshotInstances(&reqParams)
	if err != nil {
		LogErrorResponse("DeleteBlockStorageSnapshotInstances", err, []*string{ncloud.String(id)})
		return err
	}
	LogResponse("DeleteBlockStorageSnapshotInstances", resp)

	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREAT"},
		Target:  []string{"TERMINATED"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetClassicBlockStorageSnapshotInstance(config, id)
			if err != nil {
				return 0, "", err
			}
			if instance == nil { // Instance is terminated.
				return instance, "TERMINATED", nil
			}
			return instance, *instance.Status, nil
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for BlockStorageSnapshotInstance state to be \"TERMINATED\": %s", err)
	}

	return nil
}

// BlockStorage Dto for block storage
type BlockStorageSnapshot struct {
	BlockStorageInstanceNo         *string `json:"block_storage_instance_no,omitempty"`
	BlockStorageSnapshotName       *string `json:"name,omitempty"`
	Description                    *string `json:"description,omitempty"`
	BlockStorageSnapshotInstanceNo *string `json:"instance_no,omitempty"`
	BlockStorageSnapshotVolumeSize *int64  `json:"volume_size,omitempty"`
	Status                         *string `json:"instance_status,omitempty"`
	Operation                      *string `json:"instance_operation,omitempty"`
	StatusName                     *string `json:"instance_status_name,omitempty"`
	// CLASSIC only
	OriginalBlockStorageInstanceNo *string `json:"original_block_storage_instance_no,omitempty"`
	OriginalBlockStorageName       *string `json:"original_block_storage_name,omitempty"`
	ServerImageProductCode         *string `json:"server_image_product_code,omitempty"`
	OsInformation                  *string `json:"os_information,omitempty"`
	// VPC only
	HypervisorType *string `json:"hypervisor_type,omitempty"`
	// for DataSource
	SnapshotNo     *string `json:"snapshot_no,omitempty"`
	BlockStorageNo *string `json:"block_storage_no,omitempty"`
}
