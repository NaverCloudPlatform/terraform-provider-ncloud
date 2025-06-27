package server

import (
	"fmt"
	"time"

	"log"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

const (
	BlockStorageSnapshotStatusCodeCreate     = "CREAT"
	BlockStorageSnapshotStatusCodeInit       = "INIT"
	BlockStorageSnapshotStatusCodeTerminated = "TERMINATED"
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

	err = createVpcBlockStorageSnapshot(d, config)
	if err != nil {
		return err
	}

	return resourceNcloudBlockStorageSnapshotRead(d, meta)
}

func resourceNcloudBlockStorageSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var r *BlockStorageSnapshot
	config := meta.(*conn.ProviderConfig)

	r, err = GetVpcBlockStorageSnapshotDetail(config, d.Id())
	if err != nil {
		return err
	}

	if r == nil {
		log.Printf("unable to find resource: %s", d.Id())
		d.SetId("")
		return nil
	}

	d.SetId(ncloud.StringValue(r.BlockStorageSnapshotInstanceNo))
	d.Set("block_storage_instance_no", ncloud.StringValue(r.BlockStorageInstanceNo))
	d.Set("name", ncloud.StringValue(r.BlockStorageSnapshotName))
	d.Set("description", ncloud.StringValue(r.Description))
	d.Set("instance_no", ncloud.StringValue(r.BlockStorageSnapshotInstanceNo))
	d.Set("volume_size", ncloud.Int64Value(r.BlockStorageSnapshotVolumeSize))
	d.Set("instance_status", ncloud.StringValue(r.Status))
	d.Set("instance_operation", ncloud.StringValue(r.Operation))
	d.Set("instance_status_name", ncloud.StringValue(r.StatusName))
	d.Set("hypervisor_type", ncloud.StringValue(r.HypervisorType))

	return nil
}

func resourceNcloudBlockStorageSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceNcloudBlockStorageSnapshotRead(d, meta)
}

func resourceNcloudBlockStorageSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	var err error
	config := meta.(*conn.ProviderConfig)

	err = deleteVpcBlockStorageSnapshot(config, d.Id())
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
		Pending: []string{BlockStorageSnapshotStatusCodeInit},
		Target:  []string{BlockStorageSnapshotStatusCodeCreate},
		Refresh: func() (interface{}, string, error) {
			resp, err := GetVpcBlockStorageSnapshotDetail(config, id)
			if err != nil {
				return 0, "", err
			}

			if resp == nil {
				return 0, "", fmt.Errorf("fail to get VPC BlockStorageSnapshot instance, %s doesn't exist", id)
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
			Status:                         common.GetCodePtrByCommonCode(inst.BlockStorageSnapshotInstanceStatus),
			Operation:                      common.GetCodePtrByCommonCode(inst.BlockStorageSnapshotInstanceOperation),
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
		Pending: []string{BlockStorageSnapshotStatusCodeCreate},
		Target:  []string{BlockStorageSnapshotStatusCodeTerminated},
		Refresh: func() (interface{}, string, error) {
			resp, err := GetVpcBlockStorageSnapshotDetail(config, id)
			if err != nil {
				return 0, "", err
			}

			if resp == nil {
				return resp, BlockStorageSnapshotStatusCodeTerminated, nil
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
	// VPC only
	HypervisorType *string `json:"hypervisor_type,omitempty"`
	// for DataSource
	SnapshotNo     *string `json:"snapshot_no,omitempty"`
	BlockStorageNo *string `json:"block_storage_no,omitempty"`
}
