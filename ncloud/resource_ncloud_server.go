package ncloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceNcloudServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudServerCreate,
		Read:   resourceNcloudServerRead,
		Delete: resourceNcloudServerDelete,
		Update: resourceNcloudServerUpdate,
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
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(3, 30),
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9-*]+$`), "Composed of alphabets, numbers, hyphen (-) and wild card (*)."),
					validation.StringMatch(regexp.MustCompile(`.*[^\\-]$`), "Hyphen (-) cannot be used for the last character and if wild card (*) is used, other characters cannot be input."),
				),
				Description: "Server name to create. default: Assigned by ncloud",
			},
			"description": {
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
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "You can set whether or not to protect return when creating. default : false",
			},
			"internet_line_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"PUBLC", "GLBL"}, false),
				Description:  "Internet line identification code. PUBLC(Public), GLBL(Global). default : PUBLC(Public)",
			},
			"fee_system_type_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A rate system identification code. There are time plan(MTRAT) and flat rate (FXSUM). Default : Time plan(MTRAT)",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Zone code. You can determine the ZONE where the server will be created. It can be obtained through the getZoneList action. Default : Assigned by NAVER Cloud Platform.",
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
			"raid_type_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Raid Type Name",
			},
			"tag_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        tagListSchemaResource,
				Description: "Instance tag list",
			},

			"instance_no": {
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
				Type:     schema.TypeString,
				Computed: true,
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
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"base_block_storage_disk_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"base_block_storage_disk_detail_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNcloudServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	reqParams, err := buildCreateServerInstanceReqParams(client, d)
	if err != nil {
		return err
	}

	var resp *server.CreateServerInstancesResponse
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		var err error
		logCommonRequest("CreateServerInstances", reqParams)
		resp, err = client.server.V2Api.CreateServerInstances(reqParams)

		log.Printf("[DEBUG] resourceNcloudServerCreate resp: %v", resp)
		if resp != nil && isRetryableErr(GetCommonResponse(resp), []string{ApiErrorUnknown, ApiErrorAuthorityParameter, ApiErrorServerObjectInOperation, ApiErrorPreviousServersHaveNotBeenEntirelyTerminated}) {
			return resource.RetryableError(err)
		}
		return resource.NonRetryableError(err)
	})

	if err != nil {
		logErrorResponse("CreateServerInstances", err, reqParams)
		return err
	}
	logCommonResponse("CreateServerInstances", GetCommonResponse(resp))

	serverInstance := resp.ServerInstanceList[0]
	d.SetId(ncloud.StringValue(serverInstance.ServerInstanceNo))

	stateConf := &resource.StateChangeConf{
		Pending: []string{"INIT", "CREAT"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getServerInstance(client, ncloud.StringValue(serverInstance.ServerInstanceNo))
			if err != nil {
				return 0, "", err
			}
			return instance, ncloud.StringValue(instance.ServerInstanceStatus.Code), nil
		},
		Timeout:    DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for ServerInstance state to be \"RUN\": %s", err)
	}

	return resourceNcloudServerRead(d, meta)
}

func resourceNcloudServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	instance, err := getServerInstance(client, d.Id())
	if err != nil {
		return err
	}

	if instance != nil {
		d.Set("instance_no", instance.ServerInstanceNo)
		d.Set("name", instance.ServerName)
		d.Set("server_image_product_code", instance.ServerImageProductCode)
		d.Set("instance_status_name", instance.ServerInstanceStatusName)
		d.Set("server_image_name", instance.ServerImageName)
		d.Set("private_ip", instance.PrivateIp)
		d.Set("cpu_count", instance.CpuCount)
		d.Set("memory_size", instance.MemorySize)
		d.Set("base_block_storage_size", instance.BaseBlockStorageSize)
		d.Set("is_fee_charging_monitoring", instance.IsFeeChargingMonitoring)
		d.Set("public_ip", instance.PublicIp)
		d.Set("private_ip", instance.PrivateIp)
		d.Set("port_forwarding_public_ip", instance.PortForwardingPublicIp)
		d.Set("port_forwarding_external_port", instance.PortForwardingExternalPort)
		d.Set("port_forwarding_internal_port", instance.PortForwardingInternalPort)
		d.Set("user_data", d.Get("user_data").(string))

		if instanceStatus := flattenCommonCode(instance.ServerInstanceStatus); instanceStatus["code"] != nil {
			d.Set("instance_status", instanceStatus["code"])
		}

		if platformType := flattenCommonCode(instance.PlatformType); platformType["code"] != nil {
			d.Set("platform_type", platformType["code"])
		}

		if instanceOperation := flattenCommonCode(instance.ServerInstanceOperation); instanceOperation["code"] != nil {
			d.Set("instance_operation", instanceOperation["code"])
		}

		if zone := flattenZone(instance.Zone); zone["zone_code"] != nil {
			d.Set("zone", zone["zone_code"])
		}

		if region := flattenRegion(instance.Region); region["region_code"] != nil {
			d.Set("region", region["region_code"])
		}

		if diskType := flattenCommonCode(instance.BaseBlockStorageDiskType); diskType["code"] != nil {
			d.Set("base_block_storage_disk_type", diskType["code"])
		}

		if diskDetailType := flattenCommonCode(instance.BaseBlockStroageDiskDetailType); diskDetailType["code"] != nil {
			d.Set("base_block_storage_disk_detail_type", diskDetailType["code"])
		}

		if LineType := flattenCommonCode(instance.InternetLineType); LineType["code"] != nil {
			d.Set("internet_line_type", LineType["code"])
		}

		if len(instance.InstanceTagList) != 0 {
			d.Set("tag_list", flattenInstanceTagList(instance.InstanceTagList))
		}
	} else {
		log.Printf("unable to find resource: %s", d.Id())
		d.SetId("") // resource not found
	}

	return nil
}

func resourceNcloudServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)
	serverInstance, err := getServerInstance(client, d.Id())
	if err != nil {
		return err
	}

	if serverInstance == nil || ncloud.StringValue(serverInstance.ServerInstanceStatus.Code) != "NSTOP" {
		if err := stopServerInstance(client, d.Id()); err != nil {
			return err
		}

		stateConf := &resource.StateChangeConf{
			Pending: []string{"RUN"},
			Target:  []string{"NSTOP"},
			Refresh: func() (interface{}, string, error) {
				instance, err := getServerInstance(client, ncloud.StringValue(serverInstance.ServerInstanceNo))
				if err != nil {
					return 0, "", err
				}
				return instance, ncloud.StringValue(instance.ServerInstanceStatus.Code), nil
			},
			Timeout:    DefaultTimeout,
			Delay:      2 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for ServerInstance state to be \"NSTOP\": %s", err)
		}
	}

	err = detachBlockStorageByServerInstanceNo(d, client, d.Id())
	if err != nil {
		log.Printf("[ERROR] detachBlockStorageByServerInstanceNo err: %s", err)
		return err
	}

	if err := terminateServerInstance(client, d.Id()); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceNcloudServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	if d.HasChange("server_product_code") {
		reqParams := &server.ChangeServerInstanceSpecRequest{
			ServerInstanceNo:  ncloud.String(d.Get("instance_no").(string)),
			ServerProductCode: ncloud.String(d.Get("server_product_code").(string)),
		}

		var resp *server.ChangeServerInstanceSpecResponse
		err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			var err error
			logCommonRequest("ChangeServerInstanceSpec", reqParams)
			resp, err = client.server.V2Api.ChangeServerInstanceSpec(reqParams)

			if resp != nil && isRetryableErr(GetCommonResponse(resp), []string{ApiErrorUnknown, ApiErrorObjectInOperation, ApiErrorObjectInOperation}) {
				logErrorResponse("retry ChangeServerInstanceSpec", err, reqParams)
				time.Sleep(time.Second * 5)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		})

		if err != nil {
			logErrorResponse("ChangeServerInstanceSpec", err, reqParams)
			return err
		}
		logCommonResponse("ChangeServerInstanceSpec", GetCommonResponse(resp))
	}

	return resourceNcloudServerRead(d, meta)
}

func buildCreateServerInstanceReqParams(client *NcloudAPIClient, d *schema.ResourceData) (*server.CreateServerInstancesRequest, error) {

	var paramAccessControlGroupConfigurationNoList []*string
	if param, ok := d.GetOk("access_control_group_configuration_no_list"); ok {
		paramAccessControlGroupConfigurationNoList = expandStringInterfaceList(param.([]interface{}))
	}
	zoneNo, err := parseZoneNoParameter(client, d)
	if err != nil {
		return nil, err
	}
	reqParams := &server.CreateServerInstancesRequest{
		InternetLineTypeCode:                  StringPtrOrNil(d.GetOk("internet_line_type")),
		ZoneNo:                                zoneNo,
		AccessControlGroupConfigurationNoList: paramAccessControlGroupConfigurationNoList,
	}

	if serverImageProductCode, ok := d.GetOk("server_image_product_code"); ok {
		reqParams.ServerImageProductCode = ncloud.String(serverImageProductCode.(string))
	}

	if serverProductCode, ok := d.GetOk("server_product_code"); ok {
		reqParams.ServerProductCode = ncloud.String(serverProductCode.(string))
	}

	if memberServerImageNo, ok := d.GetOk("member_server_image_no"); ok {
		reqParams.MemberServerImageNo = ncloud.String(memberServerImageNo.(string))
	}

	if serverName, ok := d.GetOk("name"); ok {
		reqParams.ServerName = ncloud.String(serverName.(string))
	}

	if serverDescription, ok := d.GetOk("description"); ok {
		reqParams.ServerDescription = ncloud.String(serverDescription.(string))
	}

	if loginKeyName, ok := d.GetOk("login_key_name"); ok {
		reqParams.LoginKeyName = ncloud.String(loginKeyName.(string))
	}

	if feeSystemTypeCode, ok := d.GetOk("fee_system_type_code"); ok {
		reqParams.FeeSystemTypeCode = ncloud.String(feeSystemTypeCode.(string))
	}

	if userData, ok := d.GetOk("user_data"); ok {
		reqParams.UserData = ncloud.String(userData.(string))
	}

	if raidTypeName, ok := d.GetOk("raid_type_name"); ok {
		reqParams.RaidTypeName = ncloud.String(raidTypeName.(string))
	}

	if instanceTagList, err := expandTagListParams(d.Get("tag_list").([]interface{})); err == nil {
		reqParams.InstanceTagList = instanceTagList
	}

	if IsProtectServerTermination, ok := d.GetOk("is_protect_server_termination"); ok {
		reqParams.IsProtectServerTermination = ncloud.Bool(IsProtectServerTermination.(bool))
	}

	return reqParams, nil
}

func getServerInstance(client *NcloudAPIClient, serverInstanceNo string) (*server.ServerInstance, error) {
	reqParams := new(server.GetServerInstanceListRequest)
	reqParams.ServerInstanceNoList = []*string{ncloud.String(serverInstanceNo)}
	logCommonRequest("GetServerInstanceList", reqParams)

	resp, err := client.server.V2Api.GetServerInstanceList(reqParams)

	if err != nil {
		logErrorResponse("GetServerInstanceList", err, reqParams)
		return nil, err
	}
	logCommonResponse("GetServerInstanceList", GetCommonResponse(resp))
	if len(resp.ServerInstanceList) > 0 {
		inst := resp.ServerInstanceList[0]
		return inst, nil
	}
	return nil, nil
}

func getServerZoneNo(client *NcloudAPIClient, serverInstanceNo string) (string, error) {
	serverInstance, err := getServerInstance(client, serverInstanceNo)
	if err != nil || serverInstance == nil || serverInstance.Zone == nil {
		return "", err
	}
	return *serverInstance.Zone.ZoneNo, nil
}

func stopServerInstance(client *NcloudAPIClient, serverInstanceNo string) error {
	reqParams := &server.StopServerInstancesRequest{
		ServerInstanceNoList: []*string{ncloud.String(serverInstanceNo)},
	}
	logCommonRequest("StopServerInstances", reqParams)
	resp, err := client.server.V2Api.StopServerInstances(reqParams)
	if err != nil {
		logErrorResponse("StopServerInstances", err, reqParams)
		return err
	}
	logCommonResponse("StopServerInstances", GetCommonResponse(resp))

	return nil
}

func terminateServerInstance(client *NcloudAPIClient, serverInstanceNo string) error {
	reqParams := &server.TerminateServerInstancesRequest{
		ServerInstanceNoList: []*string{ncloud.String(serverInstanceNo)},
	}

	var resp *server.TerminateServerInstancesResponse
	err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		var err error
		logCommonRequest("TerminateServerInstances", reqParams)
		resp, err = client.server.V2Api.TerminateServerInstances(reqParams)
		if err == nil && resp == nil {
			return resource.NonRetryableError(err)
		}
		if resp != nil && isRetryableErr(GetCommonResponse(resp), []string{ApiErrorUnknown, ApiErrorServerObjectInOperation2}) {
			logErrorResponse("retry TerminateServerInstances", err, reqParams)
			return resource.RetryableError(err)
		}
		logCommonResponse("TerminateServerInstances", GetCommonResponse(resp))
		return resource.NonRetryableError(err)
	})

	if err != nil {
		logErrorResponse("TerminateServerInstances", err, reqParams)
		return err
	}
	return nil
}

var tagListSchemaResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"tag_key": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Instance Tag Key",
		},
		"tag_value": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Instance Tag Value",
		},
	},
}
