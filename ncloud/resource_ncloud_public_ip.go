package ncloud

import (
	"fmt"
	"log"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNcloudPublicIPInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudPublicIPCreate,
		Read:   resourceNcloudPublicIPRead,
		Delete: resourceNcloudPublicIPDelete,
		Update: resourceNcloudPublicIPUpdate,
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
			"region_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "1",
				Description: "You can reach a state in which inout is possible by calling `data ncloud_regions`",
			},
			"zone_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "You can determine the ZONE where the server will be created. It can be obtained through the getZoneList action. Default : Assigned by NAVER Cloud Platform.",
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

func resourceNcloudPublicIPCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	reqParams := buildCreatePublicIPInstanceReqParams(conn, d)
	resp, err := conn.CreatePublicIPInstance(reqParams)
	logCommonResponse("Create Public IP Instance", reqParams, resp.CommonResponse)

	if err != nil {
		logErrorResponse("Create Public IP Instance", err, reqParams)
		return err
	}

	publicIPInstance := &resp.PublicIPInstanceList[0]
	d.SetId(publicIPInstance.PublicIPInstanceNo)

	if err := waitPublicIPInstance(conn, publicIPInstance.PublicIPInstanceNo, "USED"); err != nil {
		return err
	}

	return resourceNcloudPublicIPRead(d, meta)
}

func resourceNcloudPublicIPRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	instance, err := getPublicIPInstance(conn, d.Id())
	if err != nil {
		logErrorResponse("Create Public IP Instance", err, d.Id())
		return err
	}

	if instance != nil {
		d.Set("public_ip_instance_no", instance.PublicIPInstanceNo)
		d.Set("public_ip", instance.PublicIP)
		d.Set("public_ip_description", instance.PublicIPDescription)
		d.Set("create_date", instance.CreateDate)
		d.Set("internet_line_type", setCommonCode(instance.InternetLineType))
		d.Set("public_ip_instance_status_name", instance.PublicIPInstanceStatusName)
		d.Set("public_ip_instance_status", setCommonCode(instance.PublicIPInstanceStatus))
		d.Set("public_ip_instance_operation", setCommonCode(instance.PublicIPInstanceOperation))
		d.Set("public_ip_kind_type", setCommonCode(instance.PublicIPKindType))
	}

	return nil
}

func resourceNcloudPublicIPDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	// Check associated public ip
	if associated, err := checkAssociatedPublicIP(conn, d.Id()); associated {
		// if associated public ip, disassociated the public ip
		disassociatedPublicIP(conn, d.Id())
	} else if err != nil {
		return err
	}

	// Step 3 : public ip 삭제
	reqParams := &sdk.RequestDeletePublicIPInstances{
		PublicIPInstanceNoList: []string{d.Id()},
	}
	resp, err := conn.DeletePublicIPInstances(reqParams)
	logCommonResponse("Delete Public IP Instance", reqParams, resp.CommonResponse)

	waitDeletePublicIPInstance(conn, d.Id())

	return err
}

func resourceNcloudPublicIPUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceNcloudPublicIPRead(d, meta)
}

func buildCreatePublicIPInstanceReqParams(conn *sdk.Conn, d *schema.ResourceData) *sdk.RequestCreatePublicIPInstance {
	reqParams := &sdk.RequestCreatePublicIPInstance{
		ServerInstanceNo:     d.Get("server_instance_no").(string),
		PublicIPDescription:  d.Get("public_ip_description").(string),
		InternetLineTypeCode: d.Get("internet_line_type_code").(string),
		RegionNo:             parseRegionNoParameter(conn, d),
		ZoneNo:               d.Get("zone_no").(string),
	}
	return reqParams
}

func getPublicIPInstance(conn *sdk.Conn, publicIPInstanceNo string) (*sdk.PublicIPInstance, error) {
	reqParams := new(sdk.RequestPublicIPInstanceList)
	reqParams.PublicIPInstanceNoList = []string{publicIPInstanceNo}
	resp, err := conn.GetPublicIPInstanceList(reqParams)

	if err != nil {
		logErrorResponse("Get Public IP Instance", err, reqParams)
		return nil, err
	}
	logCommonResponse("Get Public IP Instance", reqParams, resp.CommonResponse)
	if len(resp.PublicIPInstanceList) > 0 {
		inst := &resp.PublicIPInstanceList[0]
		return inst, nil
	}
	return nil, nil
}

func checkAssociatedPublicIP(conn *sdk.Conn, publicIPInstanceNo string) (bool, error) {
	reqParams := new(sdk.RequestPublicIPInstanceList)
	reqParams.IsAssociated = "true"
	reqParams.PublicIPInstanceNoList = []string{publicIPInstanceNo}
	resp, err := conn.GetPublicIPInstanceList(reqParams)

	if err != nil {
		logErrorResponse("Check Associated Public IP Instance", err, reqParams)
		return false, err
	}

	logCommonResponse("Check Associated Public IP Instance", reqParams, resp.CommonResponse)

	if resp.TotalRows == 0 {
		return false, nil
	}

	return true, nil
}

func disassociatedPublicIP(conn *sdk.Conn, publicIPInstanceNo string) error {
	resp, err := conn.DisassociatePublicIP(publicIPInstanceNo)

	if err != nil {
		logErrorResponse("Dissociated Public IP Instance", err, publicIPInstanceNo)
		return err
	}

	logCommonResponse("Dissociated Public IP Instance", publicIPInstanceNo, resp.CommonResponse)

	return waitDiassociatePublicIP(conn, publicIPInstanceNo)
}

func waitDiassociatePublicIP(conn *sdk.Conn, publicIPInstanceNo string) error {
	reqParams := new(sdk.RequestPublicIPInstanceList)
	reqParams.PublicIPInstanceNoList = []string{publicIPInstanceNo}

	c1 := make(chan error, 1)

	go func() {
		for {
			resp, err := conn.GetPublicIPInstanceList(reqParams)

			if err != nil {
				c1 <- err
				return
			}

			if resp.TotalRows == 0 {
				c1 <- nil
				return
			}

			log.Printf("[DEBUG] Wait disssociate public ip(%s) ", publicIPInstanceNo)
			time.Sleep(time.Second * 1)
		}
	}()

	select {
	case res := <-c1:
		return res
	case <-time.After(DefaultTimeout):
		return fmt.Errorf("TIMEOUT : diassociation public ip[%s] ", publicIPInstanceNo)
	}
}

func waitPublicIPInstance(conn *sdk.Conn, publicIPInstanceNo string, status string) error {
	reqParams := new(sdk.RequestPublicIPInstanceList)
	reqParams.PublicIPInstanceNoList = []string{publicIPInstanceNo}

	c1 := make(chan error, 1)

	go func() {
		for {
			resp, err := conn.GetPublicIPInstanceList(reqParams)

			if err != nil {
				c1 <- err
				return
			}

			if resp.PublicIPInstanceList[0].PublicIPInstanceStatus.Code == status {
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

func waitDeletePublicIPInstance(conn *sdk.Conn, publicIPInstanceNo string) error {
	reqParams := new(sdk.RequestPublicIPInstanceList)
	reqParams.PublicIPInstanceNoList = []string{publicIPInstanceNo}

	c1 := make(chan error, 1)

	go func() {
		for {
			resp, err := conn.GetPublicIPInstanceList(reqParams)

			if err != nil {
				c1 <- err
				return
			}

			if resp.TotalRows == 0 {
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
