package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"time"
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
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_product_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"member_server_image_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateServerName,
			},
			"server_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"login_key_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_protect_server_termination": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"server_create_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"server_create_start_no": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"internet_line_type_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInternetLineTypeCode,
			},
			"fee_system_type_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_control_group_configuration_no_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
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
	conn := meta.(*NcloudSdk).conn

	reqParams := buildCreateServerInstanceReqParams(d)
	resp, err := conn.CreateServerInstances(reqParams)
	logCommonResponse("CreateServerInstances", err, reqParams, resp.CommonResponse)

	if err != nil {
		return err
	}

	serverInstance := &resp.ServerInstanceList[0]
	d.SetId(serverInstance.ServerInstanceNo)

	if err := WaitForInstance(conn, serverInstance.ServerInstanceNo, "RUN", DefaultCreateTimeout); err != nil {
		return err
	}
	return resourceNcloudInstanceRead(d, meta)
}

func resourceNcloudInstanceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	instance, err := getServerInstance(conn, d.Id())
	if err != nil {
		return err
	}

	if instance != nil {
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

		if d.Get("user_data").(string) != "" {
			if err != nil {
				log.Printf("[ERROR] DescribeUserData for instance got error: %#v", err)
			}
			d.Set("user_data", Base64Decode(d.Get("user_data").(string)))
		}
	}

	return nil
}

func resourceNcloudInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn
	serverInstance, err := getServerInstance(conn, d.Id())
	if err != nil {
		return err
	}

	if serverInstance.ServerInstanceStatus.Code != "NSTOP" {
		if err := stopServerInstance(conn, d.Id()); err != nil {
			return err
		}
		if err := WaitForInstance(conn, serverInstance.ServerInstanceNo, "NSTOP", DefaultStopTimeout); err != nil {
			return err
		}
	}

	return terminateServerInstance(conn, d.Id())
}

func resourceNcloudInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	if d.HasChange("serTestAccDataSourceServerImages_basicver_product_code") {
		reqParams := &sdk.RequestChangeServerInstanceSpec{
			ServerInstanceNo:  d.Get("server_instance_no").(string),
			ServerProductCode: d.Get("server_product_code").(string),
		}

		resp, err := conn.ChangeServerInstanceSpec(reqParams)
		logCommonResponse("ChangeServerInstanceSpec", err, reqParams, resp.CommonResponse)

		if err != nil {
			return err
		}
	}

	return resourceNcloudInstanceRead(d, meta)
}

func buildCreateServerInstanceReqParams(d *schema.ResourceData) *sdk.RequestCreateServerInstance {

	var paramAccessControlGroupConfigurationNoList []string
	if param, ok := d.GetOk("access_control_group_configuration_no_list"); ok {
		paramAccessControlGroupConfigurationNoList = StringList(param.(*schema.Set).List())
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

func getServerInstance(conn *sdk.Conn, serverInstanceNo string) (*sdk.ServerInstance, error) {
	reqParams := new(sdk.RequestGetServerInstanceList)
	reqParams.ServerInstanceNoList = []string{serverInstanceNo}
	resp, err := conn.GetServerInstanceList(reqParams)

	if err != nil {
		return nil, err
	}
	logCommonResponse("GetServerInstanceList", err, reqParams, resp.CommonResponse)
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
	logCommonResponse("StopServerInstances", err, reqParams, resp.CommonResponse)

	if err != nil {
		return err
	}

	return nil
}

func terminateServerInstance(conn *sdk.Conn, serverInstanceNo string) error {
	reqParams := &sdk.RequestTerminateServerInstances{
		ServerInstanceNoList: []string{serverInstanceNo},
	}
	resp, err := conn.TerminateServerInstances(reqParams)
	logCommonResponse("TerminateServerInstances", err, reqParams, resp.CommonResponse)

	if err != nil {
		// TODO: check 502 Bad Gateway error
		// return err
		return nil
	}
	return nil
}

func WaitForInstance(conn *sdk.Conn, instanceId string, status string, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultWaitForInterval
	}
	for {
		instance, err := getServerInstance(conn, instanceId)
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
