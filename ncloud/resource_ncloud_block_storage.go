package ncloud

import (
	"fmt"
	"time"

	"log"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
			"size": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Enter the block storage size to be created. You can enter in GB units, and you can only enter up to 1000 GB.",
				ValidateFunc: validation.IntBetween(10, 1000),
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Block storage name. default: Assigned by Ncloud",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Block storage description",
			},
			"disk_detail_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "You can choose a disk detail type code of HDD and SSD. default: HDD",
			},

			"instance_no": {
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
			"device_name": { // TODO: response check
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_operation": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_status_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNcloudBlockStorageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	reqParams := buildRequestBlockStorageInstance(d)

	logCommonRequest("CreateBlockStorageInstance", reqParams)

	resp, err := client.server.V2Api.CreateBlockStorageInstance(reqParams)
	if err != nil {
		logErrorResponse("CreateBlockStorageInstance", err, reqParams)
		return err
	}
	logCommonResponse("CreateBlockStorageInstance", GetCommonResponse(resp))

	blockStorageInstance := resp.BlockStorageInstanceList[0]
	blockStorageInstanceNo := ncloud.StringValue(blockStorageInstance.BlockStorageInstanceNo)
	d.SetId(blockStorageInstanceNo)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"INIT", "CREAT"},
		Target:  []string{"ATTAC"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getBlockStorageInstance(client, blockStorageInstanceNo)
			if err != nil {
				return 0, "", err
			}
			return instance, ncloud.StringValue(instance.BlockStorageInstanceStatus.Code), nil
		},
		Timeout:    DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for BlockStorageInstance state to be \"ATTAC\": %s", err)
	}

	return resourceNcloudBlockStorageRead(d, meta)
}

func resourceNcloudBlockStorageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client
	storage, err := getBlockStorageInstance(client, d.Id())
	if err != nil {
		return err
	}

	if storage != nil {
		d.Set("instance_no", storage.BlockStorageInstanceNo)
		d.Set("server_instance_no", storage.ServerInstanceNo)
		d.Set("size", ncloud.Int64Value(storage.BlockStorageSize)/GIGABYTE)
		d.Set("name", storage.BlockStorageName)
		d.Set("server_name", storage.ServerName)
		d.Set("device_name", storage.DeviceName)
		d.Set("product_code", storage.BlockStorageProductCode)
		d.Set("instance_status_name", storage.BlockStorageInstanceStatusName)
		d.Set("description", storage.BlockStorageInstanceDescription)

		if blockStorageType := flattenCommonCode(storage.BlockStorageType); blockStorageType["code"] != nil {
			d.Set("type", blockStorageType["code"])
		}

		if instanceStatus := flattenCommonCode(storage.BlockStorageInstanceStatus); instanceStatus["code"] != nil {
			d.Set("instance_status", instanceStatus["code"])
		}

		if instanceOperation := flattenCommonCode(storage.BlockStorageInstanceOperation); instanceOperation["code"] != nil {
			d.Set("instance_operation", instanceOperation["code"])
		}

		if diskType := flattenCommonCode(storage.DiskType); diskType["code"] != nil {
			d.Set("disk_type", diskType["code"])
		}

		if diskDetailType := flattenCommonCode(storage.DiskDetailType); diskDetailType["code"] != nil {
			d.Set("disk_detail_type", diskDetailType["code"])
		}
	} else {
		log.Printf("unable to find resource: %s", d.Id())
		d.SetId("") // resource not found
	}

	return nil
}

func resourceNcloudBlockStorageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client
	blockStorageInstanceNo := d.Get("instance_no").(string)
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
	reqParams := &server.CreateBlockStorageInstanceRequest{
		ServerInstanceNo: ncloud.String(d.Get("server_instance_no").(string)),
		BlockStorageSize: ncloud.Int64(int64(d.Get("size").(int))),
	}

	if blockStorageName, ok := d.GetOk("name"); ok {
		reqParams.BlockStorageName = ncloud.String(blockStorageName.(string))
	}

	if blockStorageDescription, ok := d.GetOk("description"); ok {
		reqParams.BlockStorageDescription = ncloud.String(blockStorageDescription.(string))
	}

	if diskDetailTypeCode, ok := d.GetOk("disk_detail_type"); ok {
		reqParams.DiskDetailTypeCode = ncloud.String(diskDetailTypeCode.(string))
	}

	return reqParams
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

		stateConf := &resource.StateChangeConf{
			Pending: []string{"INIT"},
			Target:  []string{"TERMINATED"},
			Refresh: func() (interface{}, string, error) {
				instance, err := getBlockStorageInstance(client, *blockStorageId)
				if err != nil {
					return 0, "", err
				}
				if instance == nil { // Instance is terminated.
					return instance, "TERMINATED", nil
				}
				return instance, ncloud.StringValue(instance.BlockStorageInstanceStatus.Code), nil
			},
			Timeout:    DefaultTimeout,
			Delay:      2 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for BlockStorageInstance state to be \"CREAT\": %s", err)
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

		stateConf := &resource.StateChangeConf{
			Pending: []string{"INIT"},
			Target:  []string{"CREAT"},
			Refresh: func() (interface{}, string, error) {
				instance, err := getBlockStorageInstance(client, blockStorageId)
				if err != nil {
					return 0, "", err
				}
				return instance, ncloud.StringValue(instance.BlockStorageInstanceStatus.Code), nil
			},
			Timeout:    DefaultTimeout,
			Delay:      2 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for BlockStorageInstance state to be \"CREAT\": %s", err)
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
