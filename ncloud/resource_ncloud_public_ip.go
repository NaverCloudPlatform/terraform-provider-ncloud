package ncloud

import (
	"fmt"
	"log"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/resource"
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
			"internet_line_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"PUBLC", "GLBL"}, false),
				Description:  "Internet line code. PUBLC(Public), GLBL(Global)",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Zone code. You can determine the ZONE where the server will be created. It can be obtained through the getZoneList action. Default : Assigned by NAVER Cloud Platform.",
			},

			"instance_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_status_name": {
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
			"kind_type": {
				Type:     schema.TypeString,
				Computed: true,
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

	if err := waitPublicIPInstance(client, ncloud.StringValue(publicIPInstance.PublicIpInstanceNo)); err != nil {
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

		if lineType := flattenCommonCode(instance.InternetLineType); lineType["code"] != nil {
			d.Set("internet_line_type", lineType["code"])
		}

		if instanceStatus := flattenCommonCode(instance.PublicIpInstanceStatus); instanceStatus["code"] != nil {
			d.Set("instance_status", instanceStatus["code"])
		}

		if instanceOperation := flattenCommonCode(instance.PublicIpInstanceOperation); instanceOperation["code"] != nil {
			d.Set("instance_operation", instanceOperation["code"])
		}

		if kindType := flattenCommonCode(instance.PublicIpKindType); kindType["code"] != nil {
			d.Set("kind_type", kindType["code"])
		}

		if zone := flattenZone(instance.Zone); zone["zone_code"] != nil {
			d.Set("zone", zone["zone_code"])
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
	if err := waitDeletePublicIPInstance(client, d.Id()); err != nil {
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
		InternetLineTypeCode: StringPtrOrNil(d.GetOk("internet_line_type")),
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

	return waitDisassociatePublicIP(client, publicIpInstanceNo)
}

func waitDisassociatePublicIP(client *NcloudAPIClient, publicIPInstanceNo string) error {
	reqParams := new(server.GetPublicIpInstanceListRequest)
	reqParams.PublicIpInstanceNoList = ncloud.StringList([]string{publicIPInstanceNo})

	stateConf := &resource.StateChangeConf{
		Pending: []string{"NOT OK"},
		Target:  []string{"OK"},
		Refresh: func() (interface{}, string, error) {
			resp, err := client.server.V2Api.GetPublicIpInstanceList(reqParams)
			if err != nil {
				return 0, "", err
			}

			if resp.PublicIpInstanceList[0].ServerInstanceAssociatedWithPublicIp.PublicIp == nil {
				return resp, "OK", nil
			}

			return resp, "NOT OK", nil
		},
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for waitDisassociatePublicIp: %s", err)
	}

	return nil
}

func waitPublicIPInstance(client *NcloudAPIClient, publicIPInstanceNo string) error {
	reqParams := new(server.GetPublicIpInstanceListRequest)
	reqParams.PublicIpInstanceNoList = ncloud.StringList([]string{publicIPInstanceNo})

	stateConf := &resource.StateChangeConf{
		Pending: []string{"CREAT"},
		Target:  []string{"USED"},
		Refresh: func() (interface{}, string, error) {
			resp, err := client.server.V2Api.GetPublicIpInstanceList(reqParams)
			if err != nil {
				return 0, "", err
			}

			return resp, ncloud.StringValue(resp.PublicIpInstanceList[0].PublicIpInstanceStatus.Code), nil
		},
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for waitPublicIPInstance: %s", err)
	}

	return nil
}

func waitDeletePublicIPInstance(client *NcloudAPIClient, publicIPInstanceNo string) error {
	reqParams := new(server.GetPublicIpInstanceListRequest)
	reqParams.PublicIpInstanceNoList = ncloud.StringList([]string{publicIPInstanceNo})

	stateConf := &resource.StateChangeConf{
		Pending: []string{"NOT OK"},
		Target:  []string{"OK"},
		Refresh: func() (interface{}, string, error) {
			resp, err := client.server.V2Api.GetPublicIpInstanceList(reqParams)
			if err != nil {
				return 0, "", err
			}

			if ncloud.Int32Value(resp.TotalRows) == 0 {
				return resp, "OK", nil
			}
			return resp, "NOT OK", nil
		},
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for waitDeletePublicIPInstance: %s", err)
	}

	return nil
}
