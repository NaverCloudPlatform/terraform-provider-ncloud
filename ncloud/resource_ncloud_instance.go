package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
	"time"
)

// Interval for checking status in WaitForXXX method
const DefaultWaitForInterval = 5

// Default timeout
const DefaultTimeout = 60
const DefaultCreateTimeout = 15 * 60

func resourceNcloudInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudInstanceCreate,
		Read:   resourceNcloudInstanceRead,
		Delete: resourceNcloudInstanceDelete,
		//Update: resourceNcloudInstanceUpdate,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DefaultCreateTimeout),
			Delete: schema.DefaultTimeout(DefaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			"server_image_product_code": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"server_product_code": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"member_server_image_no": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"server_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"server_description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"login_key_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"is_protect_server_termination": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"server_create_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
				ForceNew: true,
			},
			"server_create_start_no": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"internet_line_type_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInternetLineTypeCode,
				ForceNew:     true,
			},
			"fee_system_type_code": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"zone_no": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"access_control_group_configuration_no_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				MinItems: 1,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNcloudInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	reqParams := buildCreateServerInstanceReqParams(d)
	resp, err := conn.CreateServerInstances(reqParams)

	if err != nil {
		return err
	}
	logCommonResponse("CreateServerInstances", resp.CommonResponse)

	serverInstance := &resp.ServerInstanceList[0]
	d.SetId(serverInstance.ServerInstanceNo)

	if err := WaitForInstance(meta, serverInstance.ServerInstanceNo, "RUN", DefaultCreateTimeout); err != nil {
		return err
	}
	return nil
}

func resourceNcloudInstanceRead(d *schema.ResourceData, meta interface{}) error {
	serverInstance, err := getServerInstance(meta, d.Id())
	if err != nil {
		return err
	}

	if serverInstance != nil {
		// TODO:
		d.Set("server_name", serverInstance.ServerName)
		d.Set("server_image_product_code", serverInstance.ServerImageProductCode)
	}

	return nil
}

func resourceNcloudInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	reqParams := &sdk.RequestTerminateServerInstances{
		ServerInstanceNoList: []string{d.Id()},
	}
	if _, err := conn.TerminateServerInstances(reqParams); err != nil {
		return err
	}

	return nil
}

func WaitForInstance(meta interface{}, instanceId string, status string, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultWaitForInterval
	}
	for {
		instance, err := getServerInstance(meta, instanceId)
		if err != nil {
			return err
		}
		if instance.ServerInstanceStatus.Code == status {
			//TODO
			//Sleep one more time for timing issues
			//time.Sleep(DefaultWaitForInterval * time.Second)
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

func buildCreateServerInstanceReqParams(d *schema.ResourceData) *sdk.RequestCreateServerInstance {

	var paramAccessControlGroupConfigurationNoList []string
	if param, ok := d.GetOk("access_control_group_configuration_no_list"); ok {
		paramAccessControlGroupConfigurationNoList = convertToStringList(param.(*schema.Set).List())
	}

	reqParams := &sdk.RequestCreateServerInstance{
		ServerImageProductCode:     d.Get("server_image_product_code").(string),
		ServerProductCode:          d.Get("server_product_code").(string),
		MemberServerImageNo:        d.Get("member_server_image_no").(string),
		ServerName:                 d.Get("server_name").(string),
		ServerDescription:          d.Get("server_description").(string),
		LoginKeyName:               d.Get("login_key_name").(string),
		IsProtectServerTermination: d.Get("is_protect_server_termination").(bool),
		ServerCreateCount:          d.Get("server_create_count").(int),
		ServerCreateStartNo:        d.Get("server_create_start_no").(int),
		InternetLineTypeCode:       d.Get("internet_line_type_code").(string),
		FeeSystemTypeCode:          d.Get("fee_system_type_code").(string),
		ZoneNo:                     d.Get("zone_no").(string),
		AccessControlGroupConfigurationNoList: paramAccessControlGroupConfigurationNoList,
		UserData: d.Get("user_data").(string),
	}
	return reqParams
}

func getServerInstance(meta interface{}, serverInstanceNo string) (*sdk.ServerInstance, error) {
	conn := meta.(*NcloudSdk).conn
	reqParams := new(sdk.RequestGetServerInstanceList)
	reqParams.ServerInstanceNoList = []string{serverInstanceNo}
	resp, err := conn.GetServerInstanceList(reqParams)

	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("error: GetServerInstanceList response error")
	}
	if len(resp.ServerInstanceList) > 0 {
		return &resp.ServerInstanceList[0], nil
	}
	return nil, nil
}
