package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"time"
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
	log.Println("[DEBUG] resourceNcloudBlockStorageCreate")
	conn := meta.(*NcloudSdk).conn

	reqParams := buildRequestBlockStorageInstance(d)
	resp, err := conn.CreateBlockStorageInstance(reqParams)
	logCommonResponse("CreateBlockStorageInstance", err, reqParams, resp.CommonResponse)
	if err != nil {
		return err
	}

	blockStorageInstance := &resp.BlockStorageInstance[0]
	d.SetId(blockStorageInstance.BlockStorageInstanceNo)

	if err := waitForBlockStorageInstance(conn, blockStorageInstance.BlockStorageInstanceNo, "ATTAC", DefaultCreateTimeout); err != nil {
		return err
	}
	return resourceNcloudBlockStorageRead(d, meta)
}

func resourceNcloudBlockStorageRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceNcloudBlockStorageRead")
	conn := meta.(*NcloudSdk).conn
	storage, err := getBlockStorageInstance(conn, d.Id())
	if err != nil {
		return err
	}
	if storage != nil {
		//d.Set("block_storage_size_gb", String(d.Get("block_storage_size_gb").(int)))
		d.Set("block_storage_instance_no", storage.BlockStorageInstanceNo)
		d.Set("server_instance_no", storage.ServerInstanceNo)
		d.Set("block_storage_size", storage.BlockStorageSize)
		d.Set("block_storage_name", storage.BlockStorageName)
		d.Set("disk_detail_type_code", storage.DiskDetailType)
		d.Set("server_name", storage.ServerName)
		d.Set("block_storage_type", storage.BlockStorageType)
		d.Set("block_storage_type", map[string]interface{}{
			"code":      storage.BlockStorageType.Code,
			"code_name": storage.BlockStorageType.CodeName,
		})
		d.Set("device_name", storage.DeviceName)
		d.Set("block_storage_product_code", storage.BlockStorageProductCode)
		d.Set("block_storage_instance_status", map[string]interface{}{
			"code":      storage.BlockStorageInstanceStatus.Code,
			"code_name": storage.BlockStorageInstanceStatus.CodeName,
		})
		d.Set("block_storage_instance_operation", map[string]interface{}{
			"code":      storage.BlockStorageInstanceOperation.Code,
			"code_name": storage.BlockStorageInstanceOperation.CodeName,
		})
		d.Set("block_storage_instance_status_name", storage.BlockStorageInstanceStatusName)
		d.Set("create_date", storage.CreateDate)
		d.Set("block_storage_description", storage.BlockStorageInstanceDescription)
		d.Set("disk_type", map[string]interface{}{
			"code":      storage.DiskType.Code,
			"code_name": storage.DiskType.CodeName,
		})
		d.Set("disk_detail_type", map[string]interface{}{
			"code":      storage.DiskDetailType.Code,
			"code_name": storage.DiskDetailType.CodeName,
		})
	}

	return nil
}

func resourceNcloudBlockStorageDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceNcloudBlockStorageDelete")
	conn := meta.(*NcloudSdk).conn
	blockStorageInstanceNo := d.Get("block_storage_instance_no").(string)
	return deleteBlockStorage(conn, []string{blockStorageInstanceNo})
}

func resourceNcloudBlockStorageUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceNcloudBlockStorageRead(d, meta)
}

func buildRequestBlockStorageInstance(d *schema.ResourceData) *sdk.RequestBlockStorageInstance {
	return &sdk.RequestBlockStorageInstance{
		ServerInstanceNo:        d.Get("server_instance_no").(string),
		BlockStorageSize:        d.Get("block_storage_size_gb").(int),
		BlockStorageName:        d.Get("block_storage_name").(string),
		BlockStorageDescription: d.Get("block_storage_description").(string),
		DiskDetailTypeCode:      d.Get("disk_detail_type_code").(string),
	}
}

func getBlockStorageInstanceList(conn *sdk.Conn, serverInstanceNo string) ([]sdk.BlockStorageInstance, error) {
	reqParams := &sdk.RequestBlockStorageInstanceList{
		ServerInstanceNo: serverInstanceNo,
	}
	resp, err := conn.GetBlockStorageInstance(reqParams)
	if err != nil {
		return nil, err
	}
	logCommonResponse("GetBlockStorageInstanceList", err, reqParams, resp.CommonResponse)
	return resp.BlockStorageInstance, nil
}

func getBlockStorageInstance(conn *sdk.Conn, blockStorageInstanceNo string) (*sdk.BlockStorageInstance, error) {
	reqParams := &sdk.RequestBlockStorageInstanceList{
		BlockStorageInstanceNoList: []string{blockStorageInstanceNo},
	}
	resp, err := conn.GetBlockStorageInstance(reqParams)
	if err != nil {
		return nil, err
	}
	logCommonResponse("GetBlockStorageInstance", err, reqParams, resp.CommonResponse)
	log.Printf("[DEBUG] GetBlockStorageInstance TotalRows: %d, BlockStorageInstance: %#v", resp.TotalRows, resp.BlockStorageInstance)
	if len(resp.BlockStorageInstance) > 0 {
		inst := &resp.BlockStorageInstance[0]
		log.Printf("[DEBUG] %s BlockStorageName: %s, Status: %s", "GetBlockStorageInstance", inst.BlockStorageName, inst.BlockStorageInstanceStatusName)
		return inst, nil
	}
	return nil, nil
}

func deleteBlockStorage(conn *sdk.Conn, blockStorageIds []string) error {
	for _, blockStorageId := range blockStorageIds {
		resp, err := conn.DeleteBlockStorageInstances([]string{blockStorageId})
		logCommonResponse("DeleteBlockStorageInstances", err, blockStorageIds, getCommonResponse(resp))

		if err := waitForBlockStorageInstance(conn, blockStorageId, "CREAT", DefaultTimeout); err != nil {
			return err
		}
	}
	return nil
}

func deleteBlockStorageByServerInstanceNo(conn *sdk.Conn, serverInstanceNo string) error {
	blockStorageInstanceList, _ := getBlockStorageInstanceList(conn, serverInstanceNo)
	if len(blockStorageInstanceList) < 1 {
		return nil
	}
	var ids []string
	for _, bs := range blockStorageInstanceList {
		if bs.BlockStorageType.Code != "BASIC" { // ignore basic storage
			log.Printf("[DEBUG] deleteBlockStorageByServerInstanceNo blockStorageInstance: %#v", bs)
			ids = append(ids, bs.BlockStorageInstanceNo)
		}
	}
	return deleteBlockStorage(conn, ids)
}

func getCommonResponse(resp *sdk.BlockStorageInstanceList) common.CommonResponse {
	var commonResponse common.CommonResponse
	if resp != nil {
		commonResponse = resp.CommonResponse
	} else {
		commonResponse = common.CommonResponse{}
	}
	return commonResponse
}

func waitForBlockStorageInstance(conn *sdk.Conn, id string, status string, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultWaitForInterval
	}
	for {
		instance, err := getBlockStorageInstance(conn, id)
		if err != nil {
			return err
		}
		if instance == nil || instance.BlockStorageInstanceStatus.Code == status {
			break
		}
		timeout = timeout - DefaultWaitForInterval
		if timeout <= 0 {
			return fmt.Errorf("error: Timeout: %d", timeout)
		}
		time.Sleep(DefaultWaitForInterval * time.Second)
	}
	return nil
}
