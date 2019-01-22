package ncloud

import (
	"fmt"
	"time"

	"log"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNcloudBlockStorageSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudBlockStorageSnapshotCreate,
		Read:   resourceNcloudBlockStorageSnapshotRead,
		Update: resourceNcloudBlockStorageSnapshotUpdate,
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
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        commonCodeSchemaResource,
				Description: "Block Storage Snapshot Instance Status",
			},
			"instance_operation": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        commonCodeSchemaResource,
				Description: "Block Storage Snapshot Instance Operation",
			},
			"instance_status_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Block Storage Snapshot Instance Status Name",
			},
			"create_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation date of the Block Storage Snapshot Instance",
			},
			"instance_description": {
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
	client := meta.(*NcloudAPIClient)

	reqParams := buildRequestBlockStorageSnapshotInstance(d)
	logCommonRequest("CreateBlockStorageSnapshotInstance", reqParams)

	resp, err := client.server.V2Api.CreateBlockStorageSnapshotInstance(reqParams)
	if err != nil {
		logErrorResponse("CreateBlockStorageSnapshotInstance", err, reqParams)
		return err
	}
	logCommonResponse("CreateBlockStorageSnapshotInstance", GetCommonResponse(resp))

	blockStorageSnapshotInstance := resp.BlockStorageSnapshotInstanceList[0]
	d.SetId(ncloud.StringValue(blockStorageSnapshotInstance.BlockStorageSnapshotInstanceNo))

	if err := waitForBlockStorageSnapshotInstance(client, ncloud.StringValue(blockStorageSnapshotInstance.BlockStorageSnapshotInstanceNo), "CREAT"); err != nil {
		return err
	}
	return resourceNcloudBlockStorageRead(d, meta)
}

func resourceNcloudBlockStorageSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)
	snapshot, err := getBlockStorageSnapshotInstance(client, d.Id())
	if err != nil {
		return err
	}

	if snapshot != nil {
		d.Set("instance_no", snapshot.BlockStorageSnapshotInstanceNo)
		d.Set("name", snapshot.BlockStorageSnapshotName)
		d.Set("volume_size", snapshot.BlockStorageSnapshotVolumeSize)
		d.Set("original_block_storage_instance_no", snapshot.OriginalBlockStorageInstanceNo)
		d.Set("original_block_storage_name", snapshot.OriginalBlockStorageName)
		d.Set("instance_status_name", snapshot.BlockStorageSnapshotInstanceStatusName)
		d.Set("create_date", snapshot.CreateDate)
		d.Set("server_image_product_code", snapshot.ServerImageProductCode)
		d.Set("os_information", snapshot.OsInformation)

		if err := d.Set("instance_status", flattenCommonCode(snapshot.BlockStorageSnapshotInstanceStatus)); err != nil {
			return err
		}
		if err := d.Set("instance_operation", flattenCommonCode(snapshot.BlockStorageSnapshotInstanceOperation)); err != nil {
			return err
		}
	} else {
		log.Printf("unable to find resource: %s", d.Id())
		d.SetId("") // resource not found
	}

	return nil
}

func resourceNcloudBlockStorageSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceNcloudBlockStorageSnapshotRead(d, meta)
}

func resourceNcloudBlockStorageSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)
	blockStorageSnapshotInstanceNo := d.Get("instance_no").(string)
	if err := deleteBlockStorageSnapshotInstance(client, blockStorageSnapshotInstanceNo); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func buildRequestBlockStorageSnapshotInstance(d *schema.ResourceData) *server.CreateBlockStorageSnapshotInstanceRequest {
	return &server.CreateBlockStorageSnapshotInstanceRequest{
		BlockStorageInstanceNo:          ncloud.String(d.Get("block_storage_instance_no").(string)),
		BlockStorageSnapshotName:        ncloud.String(d.Get("name").(string)),
		BlockStorageSnapshotDescription: ncloud.String(d.Get("description").(string)),
	}
}

func getBlockStorageSnapshotInstanceList(client *NcloudAPIClient, blockStorageSnapshotInstanceNo string) ([]*server.BlockStorageSnapshotInstance, error) {
	reqParams := &server.GetBlockStorageSnapshotInstanceListRequest{
		BlockStorageSnapshotInstanceNoList: []*string{ncloud.String(blockStorageSnapshotInstanceNo)},
	}

	logCommonRequest("GetBlockStorageSnapshotInstanceList", reqParams)

	resp, err := client.server.V2Api.GetBlockStorageSnapshotInstanceList(reqParams)
	if err != nil {
		logErrorResponse("GetBlockStorageSnapshotInstanceList", err, reqParams)
		return nil, err
	}
	logCommonResponse("GetBlockStorageSnapshotInstanceList", GetCommonResponse(resp))
	return resp.BlockStorageSnapshotInstanceList, nil
}

func getBlockStorageSnapshotInstance(client *NcloudAPIClient, blockStorageSnapshotInstanceNo string) (*server.BlockStorageSnapshotInstance, error) {
	snapshots, err := getBlockStorageSnapshotInstanceList(client, blockStorageSnapshotInstanceNo)
	if err != nil {
		logErrorResponse("getBlockStorageSnapshotInstanceList", err, []*string{ncloud.String(blockStorageSnapshotInstanceNo)})
		return nil, err
	}
	if len(snapshots) > 0 {
		inst := snapshots[0]
		return inst, nil
	}
	return nil, nil
}

func deleteBlockStorageSnapshotInstance(client *NcloudAPIClient, blockStorageSnapshotInstanceNo string) error {
	reqParams := server.DeleteBlockStorageSnapshotInstancesRequest{
		BlockStorageSnapshotInstanceNoList: []*string{ncloud.String(blockStorageSnapshotInstanceNo)},
	}

	logCommonRequest("DeleteBlockStorageSnapshotInstances", reqParams)

	resp, err := client.server.V2Api.DeleteBlockStorageSnapshotInstances(&reqParams)
	if err != nil {
		logErrorResponse("DeleteBlockStorageSnapshotInstances", err, []*string{ncloud.String(blockStorageSnapshotInstanceNo)})
		return err
	}
	var commonResponse = &CommonResponse{}
	if resp != nil {
		commonResponse = GetCommonResponse(resp)
	}

	logCommonResponse("DeleteBlockStorageSnapshotInstances", commonResponse)

	if err := waitForBlockStorageSnapshotInstance(client, blockStorageSnapshotInstanceNo, "TERMT"); err != nil {
		return err
	}
	return nil
}

func waitForBlockStorageSnapshotInstance(client *NcloudAPIClient, id string, status string) error {

	c1 := make(chan error, 1)

	go func() {
		for {
			snapshot, err := getBlockStorageSnapshotInstance(client, id)

			if err != nil {
				c1 <- err
				return
			}
			if snapshot == nil || ncloud.StringValue(snapshot.BlockStorageSnapshotInstanceStatus.Code) == status {
				c1 <- nil
				return
			}
			log.Printf("[DEBUG] Wait block storage snapshot instance [%s] status [%s] to be [%s]", id, ncloud.StringValue(snapshot.BlockStorageSnapshotInstanceStatus.Code), status)
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
