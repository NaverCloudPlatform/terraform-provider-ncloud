package ncloud

import (
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNcloudBlockStorageSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudBlockStorageSnapshotCreate,
		Read:   resourceNcloudBlockStorageSnapshotRead,
		Delete: resourceNcloudBlockStorageSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DefaultCreateTimeout),
			Delete: schema.DefaultTimeout(DefaultTimeout),
		},

		Schema: map[string]*schema.Schema{
			"block_storage_instance_no": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Block storage instance No for creating snapshot.",
			},
			"block_storage_snapshot_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Block storage snapshot name to create. default : Ncloud assigns default values.",
			},
			"block_storage_snapshot_description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Descriptions on a snapshot to create",
			},

			"block_storage_snapshot_instance_no": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Block Storage Snapshot Instance Number",
			},
			"block_storage_snapshot_volume_size": {
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
			"block_storage_snapshot_instance_status": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        commonCodeSchemaResource,
				Description: "Block Storage Snapshot Instance Status",
			},
			"block_storage_snapshot_instance_operation": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        commonCodeSchemaResource,
				Description: "Block Storage Snapshot Instance Operation",
			},
			"block_storage_snapshot_instance_status_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Block Storage Snapshot Instance Status Name",
			},
			"create_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation date of the Block Storage Snapshot Instance",
			},
			"block_storage_snapshot_instance_description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Block Storage Snapshot Instance Description",
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
		},
	}
}

func resourceNcloudBlockStorageSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	reqParams := buildRequestBlockStorageSnapshotInstance(d)
	resp, err := conn.CreateBlockStorageSnapshotInstance(reqParams)
	if err != nil {
		logErrorResponse("CreateBlockStorageSnapshotInstance", err, reqParams)
		return err
	}
	logCommonResponse("CreateBlockStorageSnapshotInstance", reqParams, resp.CommonResponse)

	blockStorageSnapshotInstance := &resp.BlockStorageSnapshotInstanceList[0]
	d.SetId(blockStorageSnapshotInstance.BlockStorageSnapshotInstanceNo)

	if err := waitForBlockStorageSnapshotInstance(conn, blockStorageSnapshotInstance.BlockStorageSnapshotInstanceNo, "CREAT"); err != nil {
		return err
	}
	return resourceNcloudBlockStorageRead(d, meta)
}

func resourceNcloudBlockStorageSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn
	snapshot, err := getBlockStorageSnapshotInstance(conn, d.Id())
	if err != nil {
		return err
	}
	if snapshot != nil {
		d.Set("block_storage_snapshot_instance_no", snapshot.BlockStorageSnapshotInstanceNo)
		d.Set("block_storage_snapshot_name", snapshot.BlockStorageSnapshotName)
		d.Set("block_storage_snapshot_volume_size", snapshot.BlockStorageSnapshotVolumeSize)
		d.Set("original_block_storage_instance_no", snapshot.OriginalBlockStorageInstanceNo)
		d.Set("original_block_storage_name", snapshot.OriginalBlockStorageName)
		d.Set("block_storage_snapshot_instance_status", setCommonCode(snapshot.BlockStorageSnapshotInstanceStatus))
		d.Set("block_storage_snapshot_instance_operation", setCommonCode(snapshot.BlockStorageSnapshotInstanceOperation))
		d.Set("block_storage_snapshot_instance_status_name", snapshot.BlockStorageSnapshotInstanceStatusName)
		d.Set("create_date", snapshot.CreateDate)
		d.Set("server_image_product_code", snapshot.ServerImageProductCode)
		d.Set("os_information", snapshot.OsInformation)
	}

	return nil
}

func resourceNcloudBlockStorageSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn
	blockStorageSnapshotInstanceNo := d.Get("block_storage_snapshot_instance_no").(string)
	return deleteBlockStorageSnapshotInstance(conn, blockStorageSnapshotInstanceNo)
}

func buildRequestBlockStorageSnapshotInstance(d *schema.ResourceData) *sdk.RequestCreateBlockStorageSnapshotInstance {
	return &sdk.RequestCreateBlockStorageSnapshotInstance{
		BlockStorageInstanceNo:          d.Get("block_storage_instance_no").(string),
		BlockStorageSnapshotName:        d.Get("block_storage_snapshot_name").(string),
		BlockStorageSnapshotDescription: d.Get("block_storage_snapshot_description").(string),
	}
}

func getBlockStorageSnapshotInstanceList(conn *sdk.Conn, blockStorageSnapshotInstanceNo string) ([]sdk.BlockStorageSnapshotInstance, error) {
	reqParams := &sdk.RequestGetBlockStorageSnapshotInstanceList{
		BlockStorageSnapshotInstanceNoList: []string{blockStorageSnapshotInstanceNo},
	}
	resp, err := conn.GetBlockStorageSnapshotInstanceList(reqParams)
	if err != nil {
		logErrorResponse("GetBlockStorageSnapshotInstanceList", err, reqParams)
		return nil, err
	}
	logCommonResponse("GetBlockStorageSnapshotInstanceList", reqParams, resp.CommonResponse)
	return resp.BlockStorageSnapshotInstanceList, nil
}

func getBlockStorageSnapshotInstance(conn *sdk.Conn, blockStorageSnapshotInstanceNo string) (*sdk.BlockStorageSnapshotInstance, error) {
	snapshots, err := getBlockStorageSnapshotInstanceList(conn, blockStorageSnapshotInstanceNo)
	if err != nil {
		return nil, err
	}
	if len(snapshots) > 0 {
		inst := &snapshots[0]
		return inst, nil
	}
	return nil, nil
}

func deleteBlockStorageSnapshotInstance(conn *sdk.Conn, blockStorageSnapshotInstanceNo string) error {
	resp, err := conn.DeleteBlockStorageSnapshotInstances([]string{blockStorageSnapshotInstanceNo})
	if err != nil {
		logErrorResponse("DeleteBlockStorageSnapshotInstances", err, []string{blockStorageSnapshotInstanceNo})
		return err
	}
	var commonResponse = common.CommonResponse{}
	if resp != nil {
		commonResponse = resp.CommonResponse
	}
	logCommonResponse("DeleteBlockStorageSnapshotInstances", blockStorageSnapshotInstanceNo, commonResponse)

	if err := waitForBlockStorageSnapshotInstance(conn, blockStorageSnapshotInstanceNo, "TERMT"); err != nil {
		return err
	}
	return nil
}

func waitForBlockStorageSnapshotInstance(conn *sdk.Conn, id string, status string) error {

	c1 := make(chan error, 1)

	go func() {
		for {
			snapshot, err := getBlockStorageSnapshotInstance(conn, id)

			if err != nil {
				c1 <- err
				return
			}
			if snapshot == nil || snapshot.BlockStorageSnapshotInstanceStatus.Code == status {
				c1 <- nil
				return
			}
			time.Sleep(time.Second * 1)
		}
	}()

	select {
	case res := <-c1:
		return res
	case <-time.After(DefaultTimeout):
		return fmt.Errorf("TIMEOUT : Wait to block storage snapshot instance  (%s)", id)

	}
}
