package ncloud

import (
	"fmt"
	"time"

	"log"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNcloudBlockStorage() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudBlockStorageCreate,
		Read:   resourceNcloudBlockStorageRead,
		Delete: resourceNcloudBlockStorageDelete,
		Update: resourceNcloudBlockStorageUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DefaultCreateTimeout),
			Delete: schema.DefaultTimeout(DefaultTimeout),
		},

		Schema: map[string]*schema.Schema{
			"server_instance_no": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Server instance number to attach. Required value. The server instance number can be obtained through the getServerInstanceList action.",
			},
			"block_storage_size_gb": {
				// note : value of block_storage_size is different from the parameter and response value.
				// 	 change the parameter name to block_storage_size_gb.
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Enter the block storage size to be created. You can enter in GB units, and you can only enter up to 1000 GB.",
			},
			"block_storage_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Block storage name. default: Assigned by Ncloud",
			},
			"block_storage_description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Block storage description",
			},
			"disk_detail_type_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "You can choose a disk detail type code of HDD and SSD. default: HDD",
			},

			"block_storage_instance_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"block_storage_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"server_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"block_storage_type": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"device_name": { // TODO: response check
				Type:     schema.TypeString,
				Computed: true,
			},
			"block_storage_product_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"block_storage_instance_status": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"block_storage_instance_operation": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"block_storage_instance_status_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_type": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"disk_detail_type": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
		},
	}
}

func resourceNcloudBlockStorageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	reqParams := buildRequestBlockStorageInstance(d)

	logCommonRequest("CreateBlockStorageInstance", reqParams)

	resp, err := client.server.V2Api.CreateBlockStorageInstance(reqParams)
	if err != nil {
		logErrorResponse("CreateBlockStorageInstance", err, reqParams)
		return err
	}
	logCommonResponse("CreateBlockStorageInstance", GetCommonResponse(resp))

	blockStorageInstance := resp.BlockStorageInstanceList[0]
	d.SetId(*blockStorageInstance.BlockStorageInstanceNo)

	if err := waitForBlockStorageInstance(client, *blockStorageInstance.BlockStorageInstanceNo, "ATTAC"); err != nil {
		return err
	}
	return resourceNcloudBlockStorageRead(d, meta)
}

func resourceNcloudBlockStorageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)
	storage, err := getBlockStorageInstance(client, d.Id())
	if err != nil {
		return err
	}

	if storage != nil {
		d.Set("block_storage_instance_no", storage.BlockStorageInstanceNo)
		d.Set("server_instance_no", storage.ServerInstanceNo)
		d.Set("block_storage_size", storage.BlockStorageSize)
		d.Set("block_storage_name", storage.BlockStorageName)
		d.Set("server_name", storage.ServerName)
		d.Set("device_name", storage.DeviceName)
		d.Set("block_storage_product_code", storage.BlockStorageProductCode)
		d.Set("block_storage_instance_status_name", storage.BlockStorageInstanceStatusName)
		d.Set("create_date", storage.CreateDate)
		d.Set("block_storage_description", storage.BlockStorageInstanceDescription)

		if err := d.Set("block_storage_type", flattenCommonCode(storage.BlockStorageType)); err != nil {
			return err
		}
		if err := d.Set("block_storage_instance_status", flattenCommonCode(storage.BlockStorageInstanceStatus)); err != nil {
			return err
		}
		if err := d.Set("block_storage_instance_operation", flattenCommonCode(storage.BlockStorageInstanceOperation)); err != nil {
			return err
		}
		if err := d.Set("disk_type", flattenCommonCode(storage.DiskType)); err != nil {
			return err
		}
		if err := d.Set("disk_detail_type", flattenCommonCode(storage.DiskDetailType)); err != nil {
			return err
		}
	}

	return nil
}

func resourceNcloudBlockStorageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)
	blockStorageInstanceNo := d.Get("block_storage_instance_no").(string)
	err := detachBlockStorage(d, client, []string{blockStorageInstanceNo})
	if err != nil {
		log.Printf("[ERROR] detachBlockStorage %#v", err)
		return err
	}
	if err := deleteBlockStorage(client, []*string{ncloud.String(blockStorageInstanceNo)}); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceNcloudBlockStorageUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceNcloudBlockStorageRead(d, meta)
}

func buildRequestBlockStorageInstance(d *schema.ResourceData) *server.CreateBlockStorageInstanceRequest {
	return &server.CreateBlockStorageInstanceRequest{
		ServerInstanceNo:        ncloud.String(d.Get("server_instance_no").(string)),
		BlockStorageSize:        ncloud.Int64(int64(d.Get("block_storage_size_gb").(int))),
		BlockStorageName:        ncloud.String(d.Get("block_storage_name").(string)),
		BlockStorageDescription: ncloud.String(d.Get("block_storage_description").(string)),
		DiskDetailTypeCode:      ncloud.String(d.Get("disk_detail_type_code").(string)),
	}
}

func getBlockStorageInstanceList(client *NcloudAPIClient, serverInstanceNo string) ([]*server.BlockStorageInstance, error) {
	reqParams := &server.GetBlockStorageInstanceListRequest{
		ServerInstanceNo: ncloud.String(serverInstanceNo),
	}

	logCommonRequest("GetBlockStorageInstanceList", reqParams)

	resp, err := client.server.V2Api.GetBlockStorageInstanceList(reqParams)
	if err != nil {
		logErrorResponse("GetBlockStorageInstanceList", err, reqParams)
		return nil, err
	}
	logCommonResponse("GetBlockStorageInstanceList", GetCommonResponse(resp))
	return resp.BlockStorageInstanceList, nil
}

func getBlockStorageInstance(client *NcloudAPIClient, blockStorageInstanceNo string) (*server.BlockStorageInstance, error) {
	reqParams := &server.GetBlockStorageInstanceListRequest{
		BlockStorageInstanceNoList: ncloud.StringList([]string{blockStorageInstanceNo}),
	}

	logCommonRequest("GetBlockStorageInstance", reqParams)

	resp, err := client.server.V2Api.GetBlockStorageInstanceList(reqParams)
	if err != nil {
		logErrorResponse("GetBlockStorageInstance", err, reqParams)
		return nil, err
	}
	logCommonResponse("GetBlockStorageInstance", GetCommonResponse(resp))

	if len(resp.BlockStorageInstanceList) > 0 {
		inst := resp.BlockStorageInstanceList[0]
		return inst, nil
	}
	return nil, nil
}

func deleteBlockStorage(client *NcloudAPIClient, blockStorageIds []*string) error {
	for _, blockStorageId := range blockStorageIds {
		reqParams := server.DeleteBlockStorageInstancesRequest{
			BlockStorageInstanceNoList: []*string{blockStorageId},
		}
		logCommonRequest("DeleteBlockStorageInstances", reqParams)

		resp, err := client.server.V2Api.DeleteBlockStorageInstances(&reqParams)
		if err != nil {
			logErrorResponse("DeleteBlockStorageInstances", err, []*string{blockStorageId})
			return err
		}
		var commonResponse = &CommonResponse{}
		if resp != nil {
			commonResponse = GetCommonResponse(resp)
		}
		logCommonResponse("DeleteBlockStorageInstances", commonResponse)

		if err := waitForBlockStorageInstance(client, *blockStorageId, "CREAT"); err != nil {
			return err
		}
	}
	return nil
}

func deleteBlockStorageByServerInstanceNo(client *NcloudAPIClient, serverInstanceNo string) error {
	blockStorageInstanceList, _ := getBlockStorageInstanceList(client, serverInstanceNo)
	if len(blockStorageInstanceList) < 1 {
		return nil
	}
	var ids []*string
	for _, bs := range blockStorageInstanceList {
		if *bs.BlockStorageType.Code != "BASIC" { // ignore basic storage
			ids = append(ids, bs.BlockStorageInstanceNo)
		}
	}
	return deleteBlockStorage(client, ids)
}

func detachBlockStorage(d *schema.ResourceData, client *NcloudAPIClient, blockStorageIds []string) error {
	var resp *server.DetachBlockStorageInstancesResponse
	for _, blockStorageId := range blockStorageIds {
		reqParams := &server.DetachBlockStorageInstancesRequest{
			BlockStorageInstanceNoList: []*string{ncloud.String(blockStorageId)},
		}
		err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			var err error

			logCommonRequest("DetachBlockStorageInstances", reqParams)

			resp, err = client.server.V2Api.DetachBlockStorageInstances(reqParams)
			if err == nil && resp == nil {
				return resource.NonRetryableError(err)
			}
			if resp != nil && isRetryableErr(GetCommonResponse(resp), []string{ApiErrorUnknown, ApiErrorDetachingMountedStorage}) {
				logErrorResponse("retry DetachBlockStorageInstances", err, reqParams)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		})

		if err != nil {
			logErrorResponse("DetachBlockStorageInstances", err, reqParams)
			return err
		}
		logCommonResponse("DetachBlockStorageInstances", GetCommonResponse(resp))

		if err := waitForBlockStorageInstance(client, blockStorageId, "CREAT"); err != nil {
			return err
		}
	}
	return nil
}

func detachBlockStorageByServerInstanceNo(d *schema.ResourceData, client *NcloudAPIClient, serverInstanceNo string) error {
	blockStorageInstanceList, _ := getBlockStorageInstanceList(client, serverInstanceNo)
	if len(blockStorageInstanceList) < 1 {
		return nil
	}
	var ids []string
	for _, bs := range blockStorageInstanceList {
		if *bs.BlockStorageType.Code != "BASIC" { // ignore basic storage
			ids = append(ids, ncloud.StringValue(bs.BlockStorageInstanceNo))
		}
	}
	return detachBlockStorage(d, client, ids)
}

func waitForBlockStorageInstance(client *NcloudAPIClient, id string, status string) error {

	c1 := make(chan error, 1)

	go func() {
		for {
			instance, err := getBlockStorageInstance(client, id)

			if err != nil {
				c1 <- err
				return
			}
			if instance == nil || ncloud.StringValue(instance.BlockStorageInstanceStatus.Code) == status {
				c1 <- nil
				return
			}
			log.Printf("[DEBUG] Wait block storage instance [%s] status [%s] to be [%s]", id, ncloud.StringValue(instance.BlockStorageInstanceStatus.Code), status)
			time.Sleep(time.Second * 1)
		}
	}()

	select {
	case res := <-c1:
		return res
	case <-time.After(DefaultTimeout):
		return fmt.Errorf("TIMEOUT : Wait to block storage instance  (%s)", id)

	}
}
