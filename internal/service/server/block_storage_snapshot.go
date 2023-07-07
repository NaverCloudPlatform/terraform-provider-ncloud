package server

import (
	"fmt"
	"time"

	"log"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
				Optional:    true,
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
				Optional:    true,
				Computed:    true,
				Description: "Block Storage Snapshot Instance Status Name",
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
	client := meta.(*conn.ProviderConfig).Client

	config := meta.(*conn.ProviderConfig)

	if config.SupportVPC {
		return NotSupportVpc("resource `ncloud_block_storage_snapshot`")
	}

	reqParams := buildRequestBlockStorageSnapshotInstance(d)
	LogCommonRequest("CreateBlockStorageSnapshotInstance", reqParams)

	resp, err := client.Server.V2Api.CreateBlockStorageSnapshotInstance(reqParams)
	if err != nil {
		LogErrorResponse("CreateBlockStorageSnapshotInstance", err, reqParams)
		return err
	}
	LogCommonResponse("CreateBlockStorageSnapshotInstance", GetCommonResponse(resp))

	blockStorageSnapshotInstance := resp.BlockStorageSnapshotInstanceList[0]
	blockStorageSnapshotInstanceNo := ncloud.StringValue(blockStorageSnapshotInstance.BlockStorageSnapshotInstanceNo)
	d.SetId(blockStorageSnapshotInstanceNo)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"INIT"},
		Target:  []string{"CREAT"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetBlockStorageSnapshotInstance(client, blockStorageSnapshotInstanceNo)
			if err != nil {
				return 0, "", err
			}
			return instance, ncloud.StringValue(instance.BlockStorageSnapshotInstanceStatus.Code), nil
		},
		Timeout:    conn.DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for BlockStorageSnapshotInstance state to be \"CREAT\": %s", err)
	}

	return resourceNcloudBlockStorageSnapshotRead(d, meta)
}

func resourceNcloudBlockStorageSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*conn.ProviderConfig).Client
	snapshot, err := GetBlockStorageSnapshotInstance(client, d.Id())
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
		d.Set("server_image_product_code", snapshot.ServerImageProductCode)
		d.Set("os_information", snapshot.OsInformation)

		if instanceStatus := FlattenCommonCode(snapshot.BlockStorageSnapshotInstanceStatus); instanceStatus["code"] != nil {
			d.Set("instance_status", instanceStatus["code"])
		}

		if instanceOperation := FlattenCommonCode(snapshot.BlockStorageSnapshotInstanceOperation); instanceOperation["code"] != nil {
			d.Set("instance_operation", instanceOperation["code"])
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
	client := meta.(*conn.ProviderConfig).Client
	blockStorageSnapshotInstanceNo := d.Get("instance_no").(string)
	if err := deleteBlockStorageSnapshotInstance(client, blockStorageSnapshotInstanceNo); err != nil {
		return err
	}
	d.SetId("")
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

func getBlockStorageSnapshotInstanceList(client *conn.NcloudAPIClient, blockStorageSnapshotInstanceNo string) ([]*server.BlockStorageSnapshotInstance, error) {
	reqParams := &server.GetBlockStorageSnapshotInstanceListRequest{
		BlockStorageSnapshotInstanceNoList: []*string{ncloud.String(blockStorageSnapshotInstanceNo)},
	}

	LogCommonRequest("getBlockStorageSnapshotInstanceList", reqParams)

	resp, err := client.Server.V2Api.GetBlockStorageSnapshotInstanceList(reqParams)
	if err != nil {
		LogErrorResponse("getBlockStorageSnapshotInstanceList", err, reqParams)
		return nil, err
	}
	LogCommonResponse("getBlockStorageSnapshotInstanceList", GetCommonResponse(resp))
	return resp.BlockStorageSnapshotInstanceList, nil
}

func GetBlockStorageSnapshotInstance(client *conn.NcloudAPIClient, blockStorageSnapshotInstanceNo string) (*server.BlockStorageSnapshotInstance, error) {
	snapshots, err := getBlockStorageSnapshotInstanceList(client, blockStorageSnapshotInstanceNo)
	if err != nil {
		LogErrorResponse("getBlockStorageSnapshotInstanceList", err, []*string{ncloud.String(blockStorageSnapshotInstanceNo)})
		return nil, err
	}
	if len(snapshots) > 0 {
		inst := snapshots[0]
		return inst, nil
	}
	return nil, nil
}

func deleteBlockStorageSnapshotInstance(client *conn.NcloudAPIClient, blockStorageSnapshotInstanceNo string) error {
	reqParams := server.DeleteBlockStorageSnapshotInstancesRequest{
		BlockStorageSnapshotInstanceNoList: []*string{ncloud.String(blockStorageSnapshotInstanceNo)},
	}

	LogCommonRequest("DeleteBlockStorageSnapshotInstances", reqParams)

	resp, err := client.Server.V2Api.DeleteBlockStorageSnapshotInstances(&reqParams)
	if err != nil {
		LogErrorResponse("DeleteBlockStorageSnapshotInstances", err, []*string{ncloud.String(blockStorageSnapshotInstanceNo)})
		return err
	}
	var commonResponse = &CommonResponse{}
	if resp != nil {
		commonResponse = GetCommonResponse(resp)
	}

	LogCommonResponse("DeleteBlockStorageSnapshotInstances", commonResponse)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"CREAT"},
		Target:  []string{"TERMT"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetBlockStorageSnapshotInstance(client, blockStorageSnapshotInstanceNo)
			if err != nil {
				return 0, "", err
			}
			if instance == nil { // Instance is terminated.
				return instance, "TERMT", nil
			}
			return instance, ncloud.StringValue(instance.BlockStorageSnapshotInstanceStatus.Code), nil
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for BlockStorageSnapshotInstance state to be \"TERMT\": %s", err)
	}

	return nil
}
