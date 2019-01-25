package ncloud

import (
	"fmt"
	"log"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceNcloudPublicIpInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudPublicIpCreate,
		Read:   resourceNcloudPublicIpRead,
		Delete: resourceNcloudPublicIpDelete,
		Update: resourceNcloudPublicIpUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"server_instance_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Server instance No. to assign after creating a public IP. You can get one by calling getPublicIpTargetServerInstanceList.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 10000),
				Description:  "Public IP description.",
			},
			"internet_line_type_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PUBLC", "GLBL"}, false),
				Description:  "Internet line code. PUBLC(Public), GLBL(Global)",
			},
			"region_code": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Region code. Get available values using the `data ncloud_regions`.",
				ConflictsWith: []string{"region_no"},
			},
			"region_no": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Region number. Get available values using the `data ncloud_regions`.",
				ConflictsWith: []string{"region_code"},
			},
			"zone_code": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Zone code. You can determine the ZONE where the server will be created. It can be obtained through the getZoneList action. Default : Assigned by NAVER Cloud Platform.",
				ConflictsWith: []string{"zone_no"},
			},
			"zone_no": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Zone number. You can determine the ZONE where the server will be created. It can be obtained through the getZoneList action. Default : Assigned by NAVER Cloud Platform.",
				ConflictsWith: []string{"zone_code"},
			},

			"instance_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internet_line_type": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"instance_status_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_status": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"instance_operation": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"kind_type": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
		},
	}
}

func resourceNcloudPublicIpCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	reqParams, err := buildCreatePublicIpInstanceReqParams(client, d)
	if err != nil {
		return err
	}
	logCommonRequest("CreatePublicIpInstance", reqParams)

	resp, err := client.server.V2Api.CreatePublicIpInstance(reqParams)
	if err != nil {
		logErrorResponse("CreatePublicIpInstance", err, reqParams)
		return err
	}
	logCommonResponse("CreatePublicIpInstance", GetCommonResponse(resp))

	publicIPInstance := resp.PublicIpInstanceList[0]
	d.SetId(ncloud.StringValue(publicIPInstance.PublicIpInstanceNo))

	if err := waitPublicIpInstance(client, ncloud.StringValue(publicIPInstance.PublicIpInstanceNo), "USED"); err != nil {
		return err
	}

	return resourceNcloudPublicIpRead(d, meta)
}

func resourceNcloudPublicIpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	instance, err := getPublicIpInstance(client, d.Id())
	if err != nil {
		logErrorResponse("Create Public IP Instance", err, d.Id())
		return err
	}

	if instance != nil {
		d.Set("instance_no", instance.PublicIpInstanceNo)
		d.Set("public_ip", instance.PublicIp)
		d.Set("description", instance.PublicIpDescription)
		d.Set("instance_status_name", instance.PublicIpInstanceStatusName)

		if err := d.Set("internet_line_type", flattenCommonCode(instance.InternetLineType)); err != nil {
			return err
		}
		if err := d.Set("instance_status", flattenCommonCode(instance.PublicIpInstanceStatus)); err != nil {
			return err
		}
		if err := d.Set("instance_operation", flattenCommonCode(instance.PublicIpInstanceOperation)); err != nil {
			return err
		}
		if err := d.Set("kind_type", flattenCommonCode(instance.PublicIpKindType)); err != nil {
			return err
		}
	} else {
		log.Printf("unable to find resource: %s", d.Id())
		d.SetId("") // resource not found
	}

	return nil
}

func resourceNcloudPublicIpDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	// Check associated public ip
	if associated, err := checkAssociatedPublicIp(client, d.Id()); associated {
		// if associated public ip, disassociated the public ip
		disassociatedPublicIp(client, d.Id())
	} else if err != nil {
		return err
	}

	// Step 3 : public ip 삭제
	reqParams := &server.DeletePublicIpInstancesRequest{
		PublicIpInstanceNoList: ncloud.StringList([]string{d.Id()}),
	}
	logCommonRequest("DeletePublicIpInstances", reqParams)
	resp, err := client.server.V2Api.DeletePublicIpInstances(reqParams)
	logCommonResponse("DeletePublicIpInstances", GetCommonResponse(resp))
	if err != nil {
		logErrorResponse("Delete Public IP Instance", err, reqParams)
		return err
	}
	if err := waitDeletePublicIpInstance(client, d.Id()); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceNcloudPublicIpUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceNcloudPublicIpRead(d, meta)
}

func buildCreatePublicIpInstanceReqParams(client *NcloudAPIClient, d *schema.ResourceData) (*server.CreatePublicIpInstanceRequest, error) {
	regionNo, err := parseRegionNoParameter(client, d)
	if err != nil {
		return nil, err
	}
	zoneNo, err := parseZoneNoParameter(client, d)
	if err != nil {
		return nil, err
	}

	reqParams := &server.CreatePublicIpInstanceRequest{
		InternetLineTypeCode: StringPtrOrNil(d.GetOk("internet_line_type_code")),
		RegionNo:             regionNo,
		ZoneNo:               zoneNo,
	}

	if serverInstanceNo, ok := d.GetOk("server_instance_no"); ok {
		reqParams.ServerInstanceNo = ncloud.String(serverInstanceNo.(string))
	}

	if description, ok := d.GetOk("description"); ok {
		reqParams.PublicIpDescription = ncloud.String(description.(string))
	}

	return reqParams, nil
}

func getPublicIpInstance(client *NcloudAPIClient, publicIPInstanceNo string) (*server.PublicIpInstance, error) {
	reqParams := new(server.GetPublicIpInstanceListRequest)
	reqParams.PublicIpInstanceNoList = ncloud.StringList([]string{publicIPInstanceNo})

	logCommonRequest("GetPublicIpInstanceList", reqParams)

	resp, err := client.server.V2Api.GetPublicIpInstanceList(reqParams)
	if err != nil {
		logErrorResponse("GetPublicIpInstanceList", err, reqParams)
		return nil, err
	}
	logCommonResponse("GetPublicIpInstanceList", GetCommonResponse(resp))
	if len(resp.PublicIpInstanceList) > 0 {
		inst := resp.PublicIpInstanceList[0]
		return inst, nil
	}
	return nil, nil
}

func checkAssociatedPublicIp(client *NcloudAPIClient, publicIPInstanceNo string) (bool, error) {
	reqParams := new(server.GetPublicIpInstanceListRequest)
	reqParams.IsAssociated = ncloud.Bool(true)
	reqParams.PublicIpInstanceNoList = ncloud.StringList([]string{publicIPInstanceNo})

	logCommonRequest("GetPublicIpInstanceList", reqParams)

	resp, err := client.server.V2Api.GetPublicIpInstanceList(reqParams)
	if err != nil {
		logErrorResponse("GetPublicIpInstanceList", err, reqParams)
		return false, err
	}
	logCommonResponse("GetPublicIpInstanceList", GetCommonResponse(resp))

	if *resp.TotalRows == 0 {
		return false, nil
	}

	return true, nil
}

func disassociatedPublicIp(client *NcloudAPIClient, publicIpInstanceNo string) error {
	reqParams := &server.DisassociatePublicIpFromServerInstanceRequest{PublicIpInstanceNo: ncloud.String(publicIpInstanceNo)}

	logCommonRequest("DisassociatePublicIpFromServerInstance", reqParams)

	resp, err := client.server.V2Api.DisassociatePublicIpFromServerInstance(reqParams)
	if err != nil {
		logErrorResponse("DisassociatePublicIpFromServerInstance", err, publicIpInstanceNo)
		return err
	}
	logCommonResponse("DisassociatePublicIpFromServerInstance", GetCommonResponse(resp))

	return waitDisassociatePublicIp(client, publicIpInstanceNo)
}

func waitDisassociatePublicIp(client *NcloudAPIClient, publicIPInstanceNo string) error {
	reqParams := new(server.GetPublicIpInstanceListRequest)
	reqParams.PublicIpInstanceNoList = ncloud.StringList([]string{publicIPInstanceNo})

	c1 := make(chan error, 1)

	go func() {
		for {
			resp, err := client.server.V2Api.GetPublicIpInstanceList(reqParams)

			if err != nil {
				c1 <- err
				return
			}

			if *resp.TotalRows == 0 {
				c1 <- nil
				return
			}

			log.Printf("[DEBUG] Wait disassociate public ip [%s] ", publicIPInstanceNo)
			time.Sleep(time.Second * 3)
		}
	}()

	select {
	case res := <-c1:
		return res
	case <-time.After(DefaultTimeout):
		return fmt.Errorf("TIMEOUT : disassociate public ip[%s] ", publicIPInstanceNo)
	}
}

func waitPublicIpInstance(client *NcloudAPIClient, publicIPInstanceNo string, status string) error {
	reqParams := new(server.GetPublicIpInstanceListRequest)
	reqParams.PublicIpInstanceNoList = ncloud.StringList([]string{publicIPInstanceNo})

	c1 := make(chan error, 1)

	go func() {
		for {
			resp, err := client.server.V2Api.GetPublicIpInstanceList(reqParams)

			if err != nil {
				c1 <- err
				return
			}

			if ncloud.StringValue(resp.PublicIpInstanceList[0].PublicIpInstanceStatus.Code) == status {
				c1 <- nil
				return
			}

			log.Printf("[DEBUG] Wait public ip(%s) status(%s)", publicIPInstanceNo, status)
			time.Sleep(time.Second * 1)
		}
	}()

	select {
	case res := <-c1:
		return res
	case <-time.After(DefaultTimeout):
		return fmt.Errorf("TIMEOUT : Wait public ip(%s) status(%s)", publicIPInstanceNo, status)
	}
}

func waitDeletePublicIpInstance(client *NcloudAPIClient, publicIPInstanceNo string) error {
	reqParams := new(server.GetPublicIpInstanceListRequest)
	reqParams.PublicIpInstanceNoList = ncloud.StringList([]string{publicIPInstanceNo})

	c1 := make(chan error, 1)

	go func() {
		for {
			resp, err := client.server.V2Api.GetPublicIpInstanceList(reqParams)

			if err != nil {
				c1 <- err
				return
			}

			if ncloud.Int32Value(resp.TotalRows) == 0 {
				c1 <- nil
				return
			}

			log.Printf("[DEBUG] Wait to delete public ip(%s)", publicIPInstanceNo)
			time.Sleep(time.Second * 1)
		}
	}()

	select {
	case res := <-c1:
		return res
	case <-time.After(DefaultTimeout):
		return fmt.Errorf("TIMEOUT : Wait to delete public ip(%s)", publicIPInstanceNo)
	}
}
