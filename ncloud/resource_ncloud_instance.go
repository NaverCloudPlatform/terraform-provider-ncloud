package ncloud

import (
	"fmt"
	"log"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNcloudInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudInstanceCreate,
		Read:   resourceNcloudInstanceRead,
		Delete: resourceNcloudInstanceDelete,
		Update: resourceNcloudInstanceUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DefaultCreateTimeout),
			Delete: schema.DefaultTimeout(DefaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			"server_image_product_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Server image product code to determine which server image to create. It can be obtained through getServerImageProductList. You are required to select one among two parameters: server image product code (server_image_product_code) and member server image number(member_server_image_no).",
			},
			"server_product_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Server product code to determine the server specification to create. It can be obtained through the getServerProductList action. Default : Selected as minimum specification. The minimum standards are 1. memory 2. CPU 3. basic block storage size 4. disk type (NET,LOCAL)",
			},
			"member_server_image_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Required value when creating a server from a manually created server image. It can be obtained through the getMemberServerImageList action.",
			},
			"server_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateServerName,
				Description:  "Server name to create. default: Assigned by NCloud",
			},
			"server_description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Server description to create",
			},
			"login_key_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The login key name to encrypt with the public key. Default : Uses the most recently created login key name",
			},
			"is_protect_server_termination": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateBoolValue,
				Description:  "You can set whether or not to protect return when creating. default : false",
			},
			"server_create_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Number of servers that can be created at a time, and not more than 20 servers can be created at a time. default: 1",
			},
			"server_create_start_no": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "If you create multiple servers at once, the server name will be serialized. You can set the starting number of the serial numbers. The total number of servers created and server starting number cannot exceed 1000. Default : If number of servers created(serverCreateCount) is greater than 1, and if there is no corresponding parameter value, the default will start from 001",
			},
			"internet_line_type_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInternetLineTypeCode,
				Description:  "Internet line identification code. PUBLC(Public), GLBL(Global). default : PUBLC(Public)",
			},
			"fee_system_type_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A rate system identification code. There are time plan(MTRAT) and flat rate (FXSUM). Default : Time plan(MTRAT)",
			},
			"zone_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "You can determine the ZONE where the server will be created. It can be obtained through the getZoneList action. Default : Assigned by NAVER Cloud Platform.",
			},
			"access_control_group_configuration_no_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				MinItems:    1,
				Description: "You can set the ACG created when creating the server. ACG setting number can be obtained through the getAccessControlGroupList action. Default : Default ACG number",
			},
			"user_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The server will execute the user data script set by the user at first boot. To view the column, it is returned only when viewing the server instance. You must need base64 Encoding, URL Encoding before put in value of userData. If you don't URL Encoding again it occurs signature invalid error.",
			},

			"server_instance_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"base_block_storage_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"platform_type": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"is_fee_charging_monitoring": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_image_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_instance_status": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"server_instance_operation": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"server_instance_status_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uptime": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port_forwarding_public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port_forwarding_external_port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port_forwarding_internal_port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_no": {
							Type: schema.TypeString,
						},
						"zone_name": {
							Type: schema.TypeString,
						},
						"zone_description": {
							Type: schema.TypeString,
						},
					},
				},
			},
			"region": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_no": {
							Type: schema.TypeString,
						},
						"region_code": {
							Type: schema.TypeString,
						},
						"region_name": {
							Type: schema.TypeString,
						},
					},
				},
			},
			"base_block_storage_disk_type": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"base_block_storage_disk_detail_type": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"internet_line_type": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
		},
	}
}

func resourceNcloudInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceNcloudInstanceCreate")
	conn := meta.(*NcloudSdk).conn

	reqParams := buildCreateServerInstanceReqParams(d)
	resp, err := conn.CreateServerInstances(reqParams)
	if err != nil {
		logErrorResponse("CreateServerInstances", err, reqParams)
		return err
	}
	logCommonResponse("CreateServerInstances", reqParams, resp.CommonResponse)

	serverInstance := &resp.ServerInstanceList[0]
	d.SetId(serverInstance.ServerInstanceNo)

	if err := waitForInstance(conn, serverInstance.ServerInstanceNo, "RUN", DefaultCreateTimeout); err != nil {
		return err
	}
	return resourceNcloudInstanceRead(d, meta)
}

func resourceNcloudInstanceRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceNcloudInstanceRead")
	conn := meta.(*NcloudSdk).conn

	instance, err := getServerInstance(conn, d.Id())
	if err != nil {
		return err
	}

	if instance != nil {
		d.Set("server_instance_no", instance.ServerInstanceNo)
		d.Set("server_name", instance.ServerName)
		d.Set("server_image_product_code", instance.ServerImageProductCode)
		d.Set("server_instance_status", map[string]interface{}{
			"code":      instance.PlatformType.Code,
			"code_name": instance.PlatformType.CodeName,
		})
		d.Set("server_instance_status_name", instance.ServerInstanceStatusName)
		d.Set("uptime", instance.Uptime)
		d.Set("server_image_name", instance.ServerImageName)
		d.Set("private_ip", instance.PrivateIP)
		d.Set("cpu_count", instance.CPUCount)
		d.Set("memory_size", instance.MemorySize)
		d.Set("base_block_storage_size", instance.BaseBlockStorageSize)
		d.Set("platform_type", map[string]interface{}{
			"code":      instance.PlatformType.Code,
			"code_name": instance.PlatformType.CodeName,
		})
		d.Set("is_fee_charging_monitoring", instance.IsFeeChargingMonitoring)
		d.Set("public_ip", instance.PublicIP)
		d.Set("private_ip", instance.PrivateIP)
		d.Set("server_instance_operation", map[string]interface{}{
			"code":      instance.ServerInstanceOperation.Code,
			"code_name": instance.ServerInstanceOperation.CodeName,
		})
		d.Set("create_date", instance.CreateDate)
		d.Set("uptime", instance.Uptime)
		d.Set("port_forwarding_public_ip", instance.PortForwardingPublicIP)
		d.Set("port_forwarding_external_port", instance.PortForwardingExternalPort)
		d.Set("port_forwarding_internal_port", instance.PortForwardingInternalPort)
		d.Set("zone", map[string]interface{}{
			"zone_no":          instance.Zone.ZoneNo,
			"zone_name":        instance.Zone.ZoneName,
			"zone_description": instance.Zone.ZoneDescription,
		})
		d.Set("region", map[string]interface{}{
			"region_no":   instance.Region.RegionNo,
			"region_code": instance.Region.RegionCode,
			"region_name": instance.Region.RegionName,
		})
		d.Set("base_block_storage_disk_type", map[string]interface{}{
			"code":      instance.BaseBlockStorageDiskType.Code,
			"code_name": instance.BaseBlockStorageDiskType.CodeName,
		})
		d.Set("base_block_storage_disk_detail_type", map[string]interface{}{
			"code":      instance.BaseBlockStroageDiskDetailType.Code,
			"code_name": instance.BaseBlockStroageDiskDetailType.CodeName,
		})
		d.Set("internet_line_type", map[string]interface{}{
			"code":      instance.InternetLineType.Code,
			"code_name": instance.InternetLineType.CodeName,
		})

		if userData, ok := d.GetOk("user_data"); ok {
			d.Set("user_data", Base64Decode(userData.(string)))
		}
	}

	return nil
}

func resourceNcloudInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceNcloudInstanceDelete")
	conn := meta.(*NcloudSdk).conn
	serverInstance, err := getServerInstance(conn, d.Id())
	if err != nil {
		return err
	}

	if serverInstance.ServerInstanceStatus.Code != "NSTOP" {
		if err := stopServerInstance(conn, d.Id()); err != nil {
			return err
		}
		if err := waitForInstance(conn, serverInstance.ServerInstanceNo, "NSTOP", DefaultStopTimeout); err != nil {
			return err
		}
	}

	err = deleteBlockStorageByServerInstanceNo(conn, d.Id())
	if err != nil {
		log.Printf("[ERROR] deleteBlockStorageByServerInstanceNo err: %s", err)
		return err
	}

	return terminateServerInstance(conn, d.Id())
}

func resourceNcloudInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceNcloudInstanceUpdate")
	conn := meta.(*NcloudSdk).conn

	if d.HasChange("serTestAccDataSourceServerImages_basicver_product_code") {
		reqParams := &sdk.RequestChangeServerInstanceSpec{
			ServerInstanceNo:  d.Get("server_instance_no").(string),
			ServerProductCode: d.Get("server_product_code").(string),
		}

		resp, err := conn.ChangeServerInstanceSpec(reqParams)
		if err != nil {
			logErrorResponse("ChangeServerInstanceSpec", err, reqParams)
			return err
		}
		logCommonResponse("ChangeServerInstanceSpec", reqParams, resp.CommonResponse)
	}

	return resourceNcloudInstanceRead(d, meta)
}

func buildCreateServerInstanceReqParams(d *schema.ResourceData) *sdk.RequestCreateServerInstance {

	var paramAccessControlGroupConfigurationNoList []string
	if param, ok := d.GetOk("access_control_group_configuration_no_list"); ok {
		paramAccessControlGroupConfigurationNoList = StringList(param.([]interface{}))
	}

	reqParams := &sdk.RequestCreateServerInstance{
		ServerImageProductCode:     d.Get("server_image_product_code").(string),
		ServerProductCode:          d.Get("server_product_code").(string),
		MemberServerImageNo:        d.Get("member_server_image_no").(string),
		ServerName:                 d.Get("server_name").(string),
		ServerDescription:          d.Get("server_description").(string),
		LoginKeyName:               d.Get("login_key_name").(string),
		IsProtectServerTermination: d.Get("is_protect_server_termination").(string),
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

func getServerInstance(conn *sdk.Conn, serverInstanceNo string) (*sdk.ServerInstance, error) {
	fmt.Printf("[DEBUG] getServerInstance")
	reqParams := new(sdk.RequestGetServerInstanceList)
	reqParams.ServerInstanceNoList = []string{serverInstanceNo}
	resp, err := conn.GetServerInstanceList(reqParams)

	if err != nil {
		logErrorResponse("GetServerInstanceList", err, reqParams)
		return nil, err
	}
	logCommonResponse("GetServerInstanceList", reqParams, resp.CommonResponse)
	if len(resp.ServerInstanceList) > 0 {
		inst := &resp.ServerInstanceList[0]
		log.Printf("[DEBUG] %s ServerName: %s, Status: %s", "GetServerInstanceList", inst.ServerName, inst.ServerInstanceStatusName)
		return inst, nil
	}
	return nil, nil
}

func stopServerInstance(conn *sdk.Conn, serverInstanceNo string) error {
	reqParams := &sdk.RequestStopServerInstances{
		ServerInstanceNoList: []string{serverInstanceNo},
	}
	resp, err := conn.StopServerInstances(reqParams)
	if err != nil {
		logErrorResponse("StopServerInstances", err, reqParams)
		return err
	}
	logCommonResponse("StopServerInstances", reqParams, resp.CommonResponse)

	return nil
}

func terminateServerInstance(conn *sdk.Conn, serverInstanceNo string) error {
	reqParams := &sdk.RequestTerminateServerInstances{
		ServerInstanceNoList: []string{serverInstanceNo},
	}
	resp, err := conn.TerminateServerInstances(reqParams)
	if err != nil {
		logErrorResponse("TerminateServerInstances", err, reqParams)
		// TODO: check 502 Bad Gateway error
		// return err
		return nil
	}
	logCommonResponse("TerminateServerInstances", reqParams, resp.CommonResponse)
	return nil
}

func waitForInstance(conn *sdk.Conn, instanceId string, status string, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultWaitForInterval
	}
	for {
		instance, err := getServerInstance(conn, instanceId)
		if err != nil {
			return err
		}
		if instance == nil || instance.ServerInstanceStatus.Code == status {
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
