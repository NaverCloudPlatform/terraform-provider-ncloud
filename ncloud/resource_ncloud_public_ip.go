package ncloud

import (
	"fmt"
	"log"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/schema"
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
			"public_ip_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(1, 10000),
				Description:  "Public IP description.",
			},
			"internet_line_type_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInternetLineTypeCode,
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

			"public_ip_instance_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internet_line_type": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"public_ip_instance_status_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip_instance_status": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"public_ip_instance_operation": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"public_ip_kind_type": {
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
	resp, err := client.server.V2Api.CreatePublicIpInstance(reqParams)
	if err != nil {
		logErrorResponse("Create Public IP Instance", err, reqParams)
		return err
	}
	logCommonResponse("Create Public IP Instance", reqParams, GetCommonResponse(resp))

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
		d.Set("public_ip_instance_no", instance.PublicIpInstanceNo)
		d.Set("public_ip", instance.PublicIp)
		d.Set("public_ip_description", instance.PublicIpDescription)
		d.Set("create_date", instance.CreateDate)
		d.Set("internet_line_type", setCommonCode(instance.InternetLineType))
		d.Set("public_ip_instance_status_name", instance.PublicIpInstanceStatusName)
		d.Set("public_ip_instance_status", setCommonCode(instance.PublicIpInstanceStatus))
		d.Set("public_ip_instance_operation", setCommonCode(instance.PublicIpInstanceOperation))
		d.Set("public_ip_kind_type", setCommonCode(instance.PublicIpKindType))
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
	resp, err := client.server.V2Api.DeletePublicIpInstances(reqParams)
	logCommonResponse("Delete Public IP Instance", reqParams, GetCommonResponse(resp))

	waitDeletePublicIpInstance(client, d.Id())

	return err
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
		ServerInstanceNo:     ncloud.String(d.Get("server_instance_no").(string)),
		PublicIpDescription:  ncloud.String(d.Get("public_ip_description").(string)),
		InternetLineTypeCode: StringPtrOrNil(d.GetOk("internet_line_type_code")),
		RegionNo:             regionNo,
		ZoneNo:               zoneNo,
	}
	return reqParams, nil
}

func getPublicIpInstance(client *NcloudAPIClient, publicIPInstanceNo string) (*server.PublicIpInstance, error) {
	reqParams := new(server.GetPublicIpInstanceListRequest)
	reqParams.PublicIpInstanceNoList = ncloud.StringList([]string{publicIPInstanceNo})
	resp, err := client.server.V2Api.GetPublicIpInstanceList(reqParams)

	if err != nil {
		logErrorResponse("Get Public IP Instance", err, reqParams)
		return nil, err
	}
	logCommonResponse("Get Public IP Instance", reqParams, GetCommonResponse(resp))
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
	resp, err := client.server.V2Api.GetPublicIpInstanceList(reqParams)

	if err != nil {
		logErrorResponse("Check Associated Public IP Instance", err, reqParams)
		return false, err
	}

	logCommonResponse("Check Associated Public IP Instance", reqParams, GetCommonResponse(resp))

	if *resp.TotalRows == 0 {
		return false, nil
	}

	return true, nil
}

func disassociatedPublicIp(client *NcloudAPIClient, publicIpInstanceNo string) error {
	resp, err := client.server.V2Api.DisassociatePublicIpFromServerInstance(&server.DisassociatePublicIpFromServerInstanceRequest{PublicIpInstanceNo: ncloud.String(publicIpInstanceNo)})

	if err != nil {
		logErrorResponse("Dissociated Public IP Instance", err, publicIpInstanceNo)
		return err
	}

	logCommonResponse("Dissociated Public IP Instance", publicIpInstanceNo, GetCommonResponse(resp))

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
