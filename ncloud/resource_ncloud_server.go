package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	RegisterResource("ncloud_server", resourceNcloudServer())
}

func resourceNcloudServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudServerCreate,
		Read:   resourceNcloudServerRead,
		Update: resourceNcloudServerUpdate,
		Delete: resourceNcloudServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DefaultCreateTimeout),
			Delete: schema.DefaultTimeout(DefaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			"server_image_product_code": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"member_server_image_no"},
			},
			"server_product_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"member_server_image_no": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"server_image_product_code"},
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateDiagFunc: ToDiagFunc(validation.All(
					validation.StringLenBetween(3, 30),
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9-*]+$`), "Composed of alphabets, numbers, hyphen (-) and wild card (*)."),
					validation.StringMatch(regexp.MustCompile(`.*[^\\-]$`), "Hyphen (-) cannot be used for the last character and if wild card (*) is used, other characters cannot be input."),
				)),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"login_key_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"is_protect_server_termination": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			// Deprecated
			"internet_line_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"PUBLC", "GLBL"}, false)),
				Deprecated:       "This parameter is no longer used.",
			},
			"fee_system_type_code": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"raid_type_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tag_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"tag_value": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"subnet_no": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"init_script_no": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"placement_group_no": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network_interface": {
				Type:          schema.TypeList,
				ConflictsWith: []string{"access_control_group_configuration_no_list"},
				Optional:      true,
				Computed:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_interface_no": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"order": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"subnet_no": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"is_encrypted_base_block_storage_volume": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"instance_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_no": {
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
			"port_forwarding_public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port_forwarding_external_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"port_forwarding_internal_port": {
				Type:     schema.TypeInt,
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
			"is_fee_charging_monitoring": {
				Type:       schema.TypeBool,
				Computed:   true,
				Deprecated: "This field no longer support",
			},
			"region": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "This field no longer support",
			},
		},
	}
}

func resourceNcloudServerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	id, err := createServerInstance(d, config)

	if err != nil {
		return err
	}

	d.SetId(ncloud.StringValue(id))
	log.Printf("[INFO] Server instance ID: %s", d.Id())

	return resourceNcloudServerRead(d, meta)
}

func resourceNcloudServerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	r, err := getServerInstance(config, d.Id())
	if err != nil {
		return err
	}

	if r == nil {
		d.SetId("")
		return nil
	}

	if config.SupportVPC {
		buildNetworkInterfaceList(config, r)
	}

	instance := ConvertToMap(r)

	SetSingularResourceDataFromMapSchema(resourceNcloudServer(), d, instance)

	return nil
}

func resourceNcloudServerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	serverInstance, err := getServerInstance(config, d.Id())
	if err != nil {
		return err
	}

	if ncloud.StringValue(serverInstance.ServerInstanceStatus) != "NSTOP" {
		log.Printf("[INFO] Stopping Instance %q for terminate", d.Id())
		if err := stopThenWaitServerInstance(config, d.Id()); err != nil {
			return err
		}
	}

	blockStorageList, err := getAdditionalBlockStorageList(config, d.Id())
	if err != nil {
		return err
	}

	if len(blockStorageList) > 0 {
		for _, blockStorage := range blockStorageList {
			if err := disconnectBlockStorage(config, blockStorage); err != nil {
				return err
			}

			if err := waitForDisconnectBlockStorage(config, d, blockStorage); err != nil {
				return err
			}
		}
	}

	if err := terminateThenWaitServerInstance(config, d.Id()); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceNcloudServerUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if d.HasChange("server_product_code") {
		if err := updateServerInstanceSpec(d, config); err != nil {
			return err
		}
	}

	if d.HasChange("is_protect_server_termination") {
		if err := updateServerProtectionTermination(d, config); err != nil {
			return err
		}
	}

	return resourceNcloudServerRead(d, meta)
}

func createServerInstance(d *schema.ResourceData, config *ProviderConfig) (*string, error) {
	if config.SupportVPC {
		return createVpcServerInstance(d, config)
	}

	return createClassicServerInstance(d, config)
}

func createClassicServerInstance(d *schema.ResourceData, config *ProviderConfig) (*string, error) {
	zoneNo, err := parseZoneNoParameter(config, d)
	if err != nil {
		return nil, err
	}

	reqParams := &server.CreateServerInstancesRequest{
		ZoneNo:                     zoneNo,
		ServerImageProductCode:     StringPtrOrNil(d.GetOk("server_image_product_code")),
		ServerProductCode:          StringPtrOrNil(d.GetOk("server_product_code")),
		MemberServerImageNo:        StringPtrOrNil(d.GetOk("member_server_image_no")),
		ServerName:                 StringPtrOrNil(d.GetOk("name")),
		ServerDescription:          StringPtrOrNil(d.GetOk("description")),
		LoginKeyName:               StringPtrOrNil(d.GetOk("login_key_name")),
		IsProtectServerTermination: BoolPtrOrNil(d.GetOk("is_protect_server_termination")),
		FeeSystemTypeCode:          StringPtrOrNil(d.GetOk("fee_system_type_code")),
		UserData:                   StringPtrOrNil(d.GetOk("user_data")),
		RaidTypeName:               StringPtrOrNil(d.GetOk("raid_type_name")),
	}

	if instanceTagList, err := expandTagListParams(d.Get("tag_list").([]interface{})); err == nil {
		reqParams.InstanceTagList = instanceTagList
	}

	if param, ok := d.GetOk("access_control_group_configuration_no_list"); ok {
		reqParams.AccessControlGroupConfigurationNoList = expandStringInterfaceList(param.([]interface{}))
	}

	var resp *server.CreateServerInstancesResponse
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		var err error
		logCommonRequest("createClassicServerInstance", reqParams)
		resp, err = config.Client.server.V2Api.CreateServerInstances(reqParams)
		if err != nil {
			errBody, _ := GetCommonErrorBody(err)
			if containsInStringList(errBody.ReturnCode, []string{ApiErrorUnknown, ApiErrorAuthorityParameter, ApiErrorServerObjectInOperation, ApiErrorPreviousServersHaveNotBeenEntirelyTerminated}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		logResponse("createClassicServerInstance", resp)
		return nil
	})

	if err != nil {
		logErrorResponse("createClassicServerInstance", err, reqParams)
		return nil, err
	}
	logResponse("createClassicServerInstance", resp)

	serverInstance := resp.ServerInstanceList[0]

	if err := waitStateNcloudServerForCreation(config, *serverInstance.ServerInstanceNo); err != nil {
		return nil, err
	}

	return serverInstance.ServerInstanceNo, nil
}

func createVpcServerInstance(d *schema.ResourceData, config *ProviderConfig) (*string, error) {
	if _, ok := d.GetOk("subnet_no"); !ok {
		return nil, ErrorRequiredArgOnVpc("subnet_no")
	}

	if _, ok := d.GetOk("access_control_group_configuration_no_list"); ok {
		return nil, NotSupportVpc("`access_control_group_configuration_no_list` of ncloud_server")
	}

	if _, ok := d.GetOk("user_data"); ok {
		return nil, NotSupportVpc("`user_data` of ncloud_server")
	}

	subnet, err := getSubnetInstance(config, d.Get("subnet_no").(string))
	if err != nil {
		return nil, err
	}

	if subnet == nil {
		return nil, fmt.Errorf("no matching subnet(%s) found", d.Get("subnet_no"))
	}

	reqParams := &vserver.CreateServerInstancesRequest{
		RegionCode:                        &config.RegionCode,
		ServerProductCode:                 StringPtrOrNil(d.GetOk("server_product_code")),
		ServerImageProductCode:            StringPtrOrNil(d.GetOk("server_image_product_code")),
		MemberServerImageInstanceNo:       StringPtrOrNil(d.GetOk("member_server_image_no")),
		ServerName:                        StringPtrOrNil(d.GetOk("name")),
		ServerDescription:                 StringPtrOrNil(d.GetOk("description")),
		LoginKeyName:                      StringPtrOrNil(d.GetOk("login_key_name")),
		IsProtectServerTermination:        BoolPtrOrNil(d.GetOk("is_protect_server_termination")),
		FeeSystemTypeCode:                 StringPtrOrNil(d.GetOk("fee_system_type_code")),
		InitScriptNo:                      StringPtrOrNil(d.GetOk("init_script_no")),
		VpcNo:                             subnet.VpcNo,
		SubnetNo:                          subnet.SubnetNo,
		PlacementGroupNo:                  StringPtrOrNil(d.GetOk("placement_group_no")),
		IsEncryptedBaseBlockStorageVolume: BoolPtrOrNil(d.GetOk("is_encrypted_base_block_storage_volume")),
	}

	if networkInterfaceList, ok := d.GetOk("network_interface"); !ok {
		defaultAcgNo, err := getDefaultAccessControlGroup(config, *subnet.VpcNo)
		if err != nil {
			return nil, err
		}

		niParam := &vserver.NetworkInterfaceParameter{
			NetworkInterfaceOrder:    ncloud.Int32(0),
			AccessControlGroupNoList: []*string{ncloud.String(defaultAcgNo)},
		}

		reqParams.NetworkInterfaceList = []*vserver.NetworkInterfaceParameter{niParam}
	} else {
		for _, vi := range networkInterfaceList.([]interface{}) {
			m := vi.(map[string]interface{})
			order := m["order"].(int)
			networkInterfaceNo := m["network_interface_no"].(string)

			networkInterface, err := getNetworkInterface(config, networkInterfaceNo)
			if err != nil {
				return nil, err
			}

			if networkInterface == nil {
				return nil, fmt.Errorf("no matching network interface [%s] found", networkInterfaceNo)
			}

			niParam := &vserver.NetworkInterfaceParameter{
				NetworkInterfaceOrder: ncloud.Int32(int32(order)),
				NetworkInterfaceNo:    networkInterface.NetworkInterfaceNo,
				SubnetNo:              networkInterface.SubnetNo,
			}

			reqParams.NetworkInterfaceList = append(reqParams.NetworkInterfaceList, niParam)
		}
	}

	logCommonRequest("createVpcServerInstance", reqParams)
	resp, err := config.Client.vserver.V2Api.CreateServerInstances(reqParams)
	if err != nil {
		logErrorResponse("createVpcServerInstance", err, reqParams)
		return nil, err
	}
	logResponse("createVpcServerInstance", resp)
	serverInstance := resp.ServerInstanceList[0]

	if err := waitStateNcloudServerForCreation(config, *serverInstance.ServerInstanceNo); err != nil {
		return nil, err
	}

	return serverInstance.ServerInstanceNo, nil
}

func waitStateNcloudServerForCreation(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"INIT", "CREAT"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getServerInstance(config, id)
			if err != nil {
				return 0, "", err
			}
			return instance, ncloud.StringValue(instance.ServerInstanceStatus), nil
		},
		Timeout:    DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for ServerInstance state to be \"RUN\": %s", err)
	}

	return nil
}

func updateServerInstanceSpec(d *schema.ResourceData, config *ProviderConfig) error {
	serverInstance, err := getServerInstance(config, d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Stopping Instance %q for server_product_code change", d.Id())
	if ncloud.StringValue(serverInstance.ServerInstanceStatus) != "NSTOP" {
		if err := stopThenWaitServerInstance(config, d.Id()); err != nil {
			return err
		}
	}

	if err := changeServerInstanceSpec(d, config); err != nil {
		return err
	}

	log.Printf("[INFO] Start Instance %q for server_product_code change", d.Id())
	if err := startThenWaitServerInstance(config, d.Id()); err != nil {
		return err
	}

	return nil
}

func changeServerInstanceSpec(d *schema.ResourceData, config *ProviderConfig) error {
	var err error
	if config.SupportVPC {
		err = changeVpcServerInstanceSpec(d, config)
	} else {
		err = changeClassicServerInstanceSpec(d, config)
	}

	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"CHNG"},
		Target:  []string{"NULL"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getServerInstance(config, d.Id())

			if err != nil {
				return 0, "", err
			}

			return instance, ncloud.StringValue(instance.ServerInstanceOperation), nil
		},
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for ServerInstance operation to be \"NULL\": %s", err)
	}

	return nil
}

func changeClassicServerInstanceSpec(d *schema.ResourceData, config *ProviderConfig) error {
	reqParams := &server.ChangeServerInstanceSpecRequest{
		ServerInstanceNo:  ncloud.String(d.Get("instance_no").(string)),
		ServerProductCode: ncloud.String(d.Get("server_product_code").(string)),
	}

	logCommonRequest("changeClassicServerInstanceSpec", reqParams)
	resp, err := config.Client.server.V2Api.ChangeServerInstanceSpec(reqParams)
	if err != nil {
		logErrorResponse("changeClassicServerInstanceSpec", err, reqParams)
		return err
	}
	logCommonResponse("changeClassicServerInstanceSpec", GetCommonResponse(resp))

	return nil
}

func changeVpcServerInstanceSpec(d *schema.ResourceData, config *ProviderConfig) error {
	reqParams := &vserver.ChangeServerInstanceSpecRequest{
		RegionCode:        &config.RegionCode,
		ServerInstanceNo:  ncloud.String(d.Get("instance_no").(string)),
		ServerProductCode: ncloud.String(d.Get("server_product_code").(string)),
	}

	logCommonRequest("changeVpcServerInstanceSpec", reqParams)
	resp, err := config.Client.vserver.V2Api.ChangeServerInstanceSpec(reqParams)
	if err != nil {
		logErrorResponse("ChangeServerInstanceSpec", err, reqParams)
		return err
	}
	logResponse("changeVpcServerInstanceSpec", resp)

	return nil
}

func updateServerProtectionTermination(d *schema.ResourceData, config *ProviderConfig) error {
	if config.SupportVPC {
		return updateVpcServerProtectionTermination(d, config)
	}

	return updateClassicServerProtectionTermination(d, config)
}

func updateVpcServerProtectionTermination(d *schema.ResourceData, config *ProviderConfig) error {
	reqParams := &vserver.SetProtectServerTerminationRequest{
		RegionCode:                 &config.RegionCode,
		ServerInstanceNo:           ncloud.String(d.Id()),
		IsProtectServerTermination: ncloud.Bool(d.Get("is_protect_server_termination").(bool)),
	}

	logCommonRequest("SetProtectServerTermination", reqParams)
	resp, err := config.Client.vserver.V2Api.SetProtectServerTermination(reqParams)
	if err != nil {
		logErrorResponse("SetProtectServerTermination", err, reqParams)
		return err
	}
	logResponse("SetProtectServerTermination", resp)

	return nil
}

func updateClassicServerProtectionTermination(d *schema.ResourceData, config *ProviderConfig) error {
	reqParams := &server.SetProtectServerTerminationRequest{
		ServerInstanceNo:           ncloud.String(d.Id()),
		IsProtectServerTermination: ncloud.Bool(d.Get("is_protect_server_termination").(bool)),
	}

	logCommonRequest("SetProtectServerTermination", reqParams)
	resp, err := config.Client.server.V2Api.SetProtectServerTermination(reqParams)
	if err != nil {
		logErrorResponse("SetProtectServerTermination", err, reqParams)
		return err
	}
	logResponse("SetProtectServerTermination", resp)

	return nil
}

func startThenWaitServerInstance(config *ProviderConfig, id string) error {
	var err error
	if config.SupportVPC {
		err = startVpcServerInstance(config, id)
	} else {
		err = startClassicServerInstance(config, id)
	}

	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"NSTOP"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getServerInstance(config, id)
			if err != nil {
				return 0, "", err
			}

			return instance, ncloud.StringValue(instance.ServerInstanceStatus), nil
		},
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for ServerInstance state to be \"RUN\": %s", err)
	}

	return nil
}

func startClassicServerInstance(config *ProviderConfig, id string) error {
	reqParams := &server.StartServerInstancesRequest{
		ServerInstanceNoList: []*string{ncloud.String(id)},
	}
	logCommonRequest("startClassicServerInstance", reqParams)
	resp, err := config.Client.server.V2Api.StartServerInstances(reqParams)
	if err != nil {
		logErrorResponse("startClassicServerInstance", err, reqParams)
		return err
	}
	logResponse("startClassicServerInstance", resp)

	return nil
}

func startVpcServerInstance(config *ProviderConfig, id string) error {
	reqParams := &vserver.StartServerInstancesRequest{
		RegionCode:           &config.RegionCode,
		ServerInstanceNoList: []*string{ncloud.String(id)},
	}
	logCommonRequest("startVpcServerInstance", reqParams)
	resp, err := config.Client.vserver.V2Api.StartServerInstances(reqParams)
	if err != nil {
		logErrorResponse("startVpcServerInstance", err, reqParams)
		return err
	}
	logResponse("startVpcServerInstance", resp)

	return nil
}

func getServerInstance(config *ProviderConfig, id string) (*ServerInstance, error) {
	if config.SupportVPC {
		return getVpcServerInstance(config, id)
	}

	return getClassicServerInstance(config, id)
}

func getClassicServerInstance(config *ProviderConfig, id string) (*ServerInstance, error) {
	reqParams := &server.GetServerInstanceListRequest{
		ServerInstanceNoList: []*string{ncloud.String(id)},
	}

	logCommonRequest("getClassicServerInstance", reqParams)
	resp, err := config.Client.server.V2Api.GetServerInstanceList(reqParams)

	if err != nil {
		logErrorResponse("getClassicServerInstance", err, reqParams)
		return nil, err
	}

	logResponse("getClassicServerInstance", resp)

	if len(resp.ServerInstanceList) == 0 {
		return nil, nil
	}

	if err := validateOneResult(len(resp.ServerInstanceList)); err != nil {
		return nil, err
	}

	return convertClassicServerInstance(resp.ServerInstanceList[0]), nil
}

func convertClassicServerInstance(r *server.ServerInstance) *ServerInstance {
	if r == nil {
		return nil
	}

	return &ServerInstance{
		ZoneNo:                         r.Zone.ZoneNo,
		ServerImageProductCode:         r.ServerImageProductCode,
		ServerProductCode:              r.ServerProductCode,
		ServerName:                     r.ServerName,
		ServerDescription:              r.ServerDescription,
		LoginKeyName:                   r.LoginKeyName,
		IsProtectServerTermination:     r.IsProtectServerTermination,
		ServerInstanceNo:               r.ServerInstanceNo,
		ServerImageName:                r.ServerImageName,
		CpuCount:                       r.CpuCount,
		MemorySize:                     r.MemorySize,
		BaseBlockStorageSize:           r.BaseBlockStorageSize,
		IsFeeChargingMonitoring:        r.IsFeeChargingMonitoring,
		PublicIp:                       r.PublicIp,
		PrivateIp:                      r.PrivateIp,
		PortForwardingPublicIp:         r.PortForwardingPublicIp,
		PortForwardingExternalPort:     r.PortForwardingExternalPort,
		PortForwardingInternalPort:     r.PortForwardingInternalPort,
		ServerInstanceStatus:           r.ServerInstanceStatus.Code,
		PlatformType:                   r.PlatformType.Code,
		ServerInstanceOperation:        r.ServerInstanceOperation.Code,
		Zone:                           r.Zone.ZoneCode,
		BaseBlockStorageDiskType:       r.BaseBlockStorageDiskType.Code,
		BaseBlockStorageDiskDetailType: flattenMapByKey(r.BaseBlockStorageDiskDetailType, "code"),
		InstanceTagList:                r.InstanceTagList,
	}
}

func getVpcServerInstance(config *ProviderConfig, id string) (*ServerInstance, error) {
	reqParams := &vserver.GetServerInstanceDetailRequest{
		RegionCode:       &config.RegionCode,
		ServerInstanceNo: ncloud.String(id),
	}

	logCommonRequest("getVpcServerInstance", reqParams)
	resp, err := config.Client.vserver.V2Api.GetServerInstanceDetail(reqParams)

	if err != nil {
		logErrorResponse("getVpcServerInstance", err, reqParams)
		return nil, err
	}

	logResponse("getVpcServerInstance", resp)

	if len(resp.ServerInstanceList) == 0 {
		return nil, nil
	}

	if err := validateOneResult(len(resp.ServerInstanceList)); err != nil {
		return nil, err
	}

	return convertVcpServerInstance(resp.ServerInstanceList[0]), nil
}

func convertVcpServerInstance(r *vserver.ServerInstance) *ServerInstance {
	if r == nil {
		return nil
	}

	instance := &ServerInstance{
		ServerImageProductCode:         r.ServerImageProductCode,
		ServerProductCode:              r.ServerProductCode,
		ServerName:                     r.ServerName,
		ServerDescription:              r.ServerDescription,
		LoginKeyName:                   r.LoginKeyName,
		IsProtectServerTermination:     r.IsProtectServerTermination,
		ServerInstanceNo:               r.ServerInstanceNo,
		CpuCount:                       r.CpuCount,
		MemorySize:                     r.MemorySize,
		PublicIp:                       r.PublicIp,
		ServerInstanceStatus:           r.ServerInstanceStatus.Code,
		PlatformType:                   r.PlatformType.Code,
		ServerInstanceOperation:        r.ServerInstanceOperation.Code,
		Zone:                           r.ZoneCode,
		BaseBlockStorageDiskType:       r.BaseBlockStorageDiskType.Code,
		BaseBlockStorageDiskDetailType: flattenMapByKey(r.BaseBlockStorageDiskDetailType, "code"),
		VpcNo:                          r.VpcNo,
		SubnetNo:                       r.SubnetNo,
		InitScriptNo:                   r.InitScriptNo,
		PlacementGroupNo:               r.PlacementGroupNo,
	}

	for _, networkInterfaceNo := range r.NetworkInterfaceNoList {
		ni := &ServerInstanceNetworkInterface{
			NetworkInterfaceNo: networkInterfaceNo,
		}

		instance.NetworkInterfaceList = append(instance.NetworkInterfaceList, ni)
	}

	return instance
}

func buildNetworkInterfaceList(config *ProviderConfig, r *ServerInstance) error {
	for _, ni := range r.NetworkInterfaceList {
		networkInterface, err := getNetworkInterface(config, *ni.NetworkInterfaceNo)

		if err != nil {
			return err
		}

		if networkInterface == nil {
			continue
		}

		re := regexp.MustCompile("[0-9]+")
		order, err := strconv.Atoi(re.FindString(*networkInterface.DeviceName))

		if err != nil {
			return fmt.Errorf("error parsing network interface device name: %s", *networkInterface.DeviceName)
		}

		ni.PrivateIp = networkInterface.Ip
		ni.SubnetNo = networkInterface.SubnetNo
		ni.NetworkInterfaceNo = networkInterface.NetworkInterfaceNo
		ni.Order = ncloud.Int32(int32(order))
	}

	return nil
}

func stopThenWaitServerInstance(config *ProviderConfig, id string) error {
	var err error

	stateConf := &resource.StateChangeConf{
		Pending: []string{"SETUP"},
		Target:  []string{"NULL"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getServerInstance(config, id)
			if err != nil {
				return 0, "", err
			}

			return instance, ncloud.StringValue(instance.ServerInstanceOperation), nil
		},
		Timeout:    DefaultStopTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for ServerInstance operation to be \"NULL\": %s", err)
	}

	if config.SupportVPC {
		err = stopVpcServerInstance(config, id)
	} else {
		err = stopClassicServerInstance(config, id)
	}

	if err != nil {
		return err
	}

	stateConf = &resource.StateChangeConf{
		Pending: []string{"RUN"},
		Target:  []string{"NSTOP"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getServerInstance(config, id)
			if err != nil {
				return 0, "", err
			}

			return instance, ncloud.StringValue(instance.ServerInstanceStatus), nil
		},
		Timeout:    DefaultStopTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for ServerInstance state to be \"NSTOP\": %s", err)
	}

	return nil
}

func stopClassicServerInstance(config *ProviderConfig, id string) error {
	reqParams := &server.StopServerInstancesRequest{
		ServerInstanceNoList: []*string{ncloud.String(id)},
	}
	logCommonRequest("stopClassicServerInstance", reqParams)
	resp, err := config.Client.server.V2Api.StopServerInstances(reqParams)
	if err != nil {
		logErrorResponse("stopClassicServerInstance", err, reqParams)
		return err
	}
	logResponse("stopClassicServerInstance", resp)

	return nil
}

func stopVpcServerInstance(config *ProviderConfig, id string) error {
	reqParams := &vserver.StopServerInstancesRequest{
		RegionCode:           &config.RegionCode,
		ServerInstanceNoList: []*string{ncloud.String(id)},
	}
	logCommonRequest("stopVpcServerInstance", reqParams)
	resp, err := config.Client.vserver.V2Api.StopServerInstances(reqParams)
	if err != nil {
		logErrorResponse("stopClassicServerInstance", err, reqParams)
		return err
	}
	logResponse("stopVpcServerInstance", resp)

	return nil
}

func terminateThenWaitServerInstance(config *ProviderConfig, id string) error {
	var err error
	if config.SupportVPC {
		err = terminateVpcServerInstance(config, id)
	} else {
		err = terminateClassicServerInstance(config, id)
	}

	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"NSTOP"},
		Target:  []string{"TERMINATED"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getServerInstance(config, id)

			if err != nil {
				return 0, "", err
			}
			if instance == nil { // Instance is terminated.
				return instance, "TERMINATED", nil
			}
			return instance, ncloud.StringValue(instance.ServerInstanceStatus), nil
		},
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for ServerInstance state to be \"TERMINATED\": %s", err)
	}

	return nil
}

func terminateClassicServerInstance(config *ProviderConfig, id string) error {
	reqParams := &server.TerminateServerInstancesRequest{
		ServerInstanceNoList: []*string{ncloud.String(id)},
	}

	var resp *server.TerminateServerInstancesResponse
	err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		var err error
		logCommonRequest("terminateClassicServerInstance", reqParams)
		resp, err = config.Client.server.V2Api.TerminateServerInstances(reqParams)
		if err != nil {
			errBody, _ := GetCommonErrorBody(err)
			if containsInStringList(errBody.ReturnCode, []string{ApiErrorUnknown, ApiErrorServerObjectInOperation2}) {
				logErrorResponse("retry terminateClassicServerInstance", err, reqParams)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		logResponse("terminateClassicServerInstance", resp)
		return nil
	})

	if err != nil {
		logErrorResponse("terminateClassicServerInstance", err, reqParams)
		return err
	}

	return nil
}

func terminateVpcServerInstance(config *ProviderConfig, id string) error {
	reqParams := &vserver.TerminateServerInstancesRequest{
		RegionCode:           &config.RegionCode,
		ServerInstanceNoList: []*string{ncloud.String(id)},
	}

	logCommonRequest("terminateVpcServerInstance", reqParams)
	resp, err := config.Client.vserver.V2Api.TerminateServerInstances(reqParams)
	logResponse("terminateVpcServerInstance", resp)

	if err != nil {
		logErrorResponse("terminateVpcServerInstance", err, reqParams)
		return err
	}

	return nil
}

func getAdditionalBlockStorageList(config *ProviderConfig, id string) ([]*BlockStorage, error) {
	if config.SupportVPC {
		return getVpcAdditionalBlockStorageList(config, id)
	} else {
		return getClassicAdditionalBlockStorageList(config, id)
	}
}

func getVpcAdditionalBlockStorageList(config *ProviderConfig, id string) ([]*BlockStorage, error) {
	resp, err := config.Client.vserver.V2Api.GetBlockStorageInstanceList(&vserver.GetBlockStorageInstanceListRequest{
		RegionCode:               &config.RegionCode,
		ServerInstanceNo:         ncloud.String(id),
		BlockStorageTypeCodeList: []*string{ncloud.String("SVRBS")},
	})

	if err != nil {
		return nil, err
	}

	if len(resp.BlockStorageInstanceList) < 1 {
		return nil, nil
	}

	blockStorageList := make([]*BlockStorage, 0)
	for _, blockStorage := range resp.BlockStorageInstanceList {
		blockStorageList = append(blockStorageList, convertVpcBlockStorage(blockStorage))
	}

	return blockStorageList, nil
}

func getClassicAdditionalBlockStorageList(config *ProviderConfig, id string) ([]*BlockStorage, error) {
	resp, err := config.Client.server.V2Api.GetBlockStorageInstanceList(&server.GetBlockStorageInstanceListRequest{
		RegionNo:                 &config.RegionCode,
		ServerInstanceNo:         ncloud.String(id),
		BlockStorageTypeCodeList: []*string{ncloud.String("SVRBS")},
	})

	if err != nil {
		return nil, err
	}

	if len(resp.BlockStorageInstanceList) < 1 {
		return nil, nil
	}

	blockStorageList := make([]*BlockStorage, 0)
	for _, blockStorage := range resp.BlockStorageInstanceList {
		blockStorageList = append(blockStorageList, convertClassicBlockStorage(blockStorage))
	}

	return blockStorageList, nil
}

func convertVpcBlockStorage(storage *vserver.BlockStorageInstance) *BlockStorage {
	return &BlockStorage{
		BlockStorageInstanceNo:  storage.BlockStorageInstanceNo,
		ServerInstanceNo:        storage.ServerInstanceNo,
		BlockStorageType:        storage.BlockStorageType.Code,
		BlockStorageName:        storage.BlockStorageName,
		BlockStorageSize:        storage.BlockStorageSize,
		DeviceName:              storage.DeviceName,
		BlockStorageProductCode: storage.BlockStorageProductCode,
		Status:                  storage.BlockStorageInstanceStatus.Code,
		Operation:               storage.BlockStorageInstanceOperation.Code,
		Description:             storage.BlockStorageDescription,
		DiskType:                storage.BlockStorageDiskType.Code,
		DiskDetailType:          storage.BlockStorageDiskDetailType.Code,
		ZoneCode:                storage.ZoneCode,
	}
}

func convertClassicBlockStorage(storage *server.BlockStorageInstance) *BlockStorage {
	return &BlockStorage{
		BlockStorageInstanceNo:  storage.BlockStorageInstanceNo,
		ServerInstanceNo:        storage.ServerInstanceNo,
		ServerName:              storage.ServerName,
		BlockStorageType:        storage.BlockStorageType.Code,
		BlockStorageName:        storage.BlockStorageName,
		BlockStorageSize:        storage.BlockStorageSize,
		DeviceName:              storage.DeviceName,
		BlockStorageProductCode: storage.BlockStorageProductCode,
		Status:                  storage.BlockStorageInstanceStatus.Code,
		Operation:               storage.BlockStorageInstanceOperation.Code,
		Description:             storage.BlockStorageInstanceDescription,
		DiskType:                storage.DiskType.Code,
		DiskDetailType:          storage.DiskDetailType.Code,
		ZoneCode:                storage.Zone.ZoneCode,
	}
}

func disconnectBlockStorage(config *ProviderConfig, storage *BlockStorage) error {
	if config.SupportVPC {
		return disconnectVpcBlockStorage(config, storage)
	} else {
		return disconnectClassicBlockStorage(config, storage)
	}
}

func disconnectVpcBlockStorage(config *ProviderConfig, storage *BlockStorage) error {
	_, err := config.Client.vserver.V2Api.DetachBlockStorageInstances(&vserver.DetachBlockStorageInstancesRequest{
		RegionCode:                 &config.RegionCode,
		BlockStorageInstanceNoList: []*string{storage.BlockStorageInstanceNo},
	})

	if err != nil {
		return err
	}

	return nil
}

func disconnectClassicBlockStorage(config *ProviderConfig, storage *BlockStorage) error {
	_, err := config.Client.server.V2Api.DetachBlockStorageInstances(&server.DetachBlockStorageInstancesRequest{
		BlockStorageInstanceNoList: []*string{storage.BlockStorageInstanceNo},
	})

	if err != nil {
		return err
	}

	return nil
}

func waitForDisconnectBlockStorage(config *ProviderConfig, d *schema.ResourceData, storage *BlockStorage) error {
	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		blockStorage, err := getBlockStorage(config, *storage.BlockStorageInstanceNo)
		if err != nil {
			return resource.RetryableError(err)
		}
		if *blockStorage.Status != BlockStorageStatusCodeCreate {
			return resource.RetryableError(fmt.Errorf("sill connected block storage(%s)", *blockStorage.BlockStorageInstanceNo))
		}
		return nil
	})
}

func getServerZoneNo(config *ProviderConfig, serverInstanceNo string) (string, error) {
	instance, err := getServerInstance(config, serverInstanceNo)
	if err != nil || instance == nil || instance.ZoneNo == nil {
		return "", err
	}
	return *instance.ZoneNo, nil
}

// ServerInstance server instance model
type ServerInstance struct {
	// Request
	ZoneNo                         *string               `json:"zone_no,omitempty"`
	ServerImageProductCode         *string               `json:"server_image_product_code,omitempty"`
	ServerProductCode              *string               `json:"server_product_code,omitempty"`
	MemberServerImageNo            *string               `json:"member_server_image_no,omitempty"`
	ServerName                     *string               `json:"name,omitempty"`
	ServerDescription              *string               `json:"description,omitempty"`
	LoginKeyName                   *string               `json:"login_key_name,omitempty"`
	IsProtectServerTermination     *bool                 `json:"is_protect_server_termination,omitempty"`
	FeeSystemTypeCode              *string               `json:"fee_system_type_code,omitempty"`
	UserData                       *string               `json:"user_data,omitempty"`
	RaidTypeName                   *string               `json:"raid_type_name,omitempty"`
	ServerInstanceNo               *string               `json:"instance_no,omitempty"`
	ServerImageName                *string               `json:"server_image_name,omitempty"`
	CpuCount                       *int32                `json:"cpu_count,omitempty"`
	MemorySize                     *int64                `json:"memory_size,omitempty"`
	BaseBlockStorageSize           *int64                `json:"base_block_storage_size,omitempty"`
	IsFeeChargingMonitoring        *bool                 `json:"is_fee_charging_monitoring,omitempty"`
	PublicIp                       *string               `json:"public_ip,omitempty"`
	PrivateIp                      *string               `json:"private_ip,omitempty"`
	PortForwardingPublicIp         *string               `json:"port_forwarding_public_ip,omitempty"`
	PortForwardingExternalPort     *int32                `json:"port_forwarding_external_port,omitempty"`
	PortForwardingInternalPort     *int32                `json:"port_forwarding_internal_port,omitempty"`
	ServerInstanceStatus           *string               `json:"status,omitempty"`
	PlatformType                   *string               `json:"platform_type,omitempty"`
	ServerInstanceOperation        *string               `json:"operation,omitempty"`
	Zone                           *string               `json:"zone,omitempty"`
	BaseBlockStorageDiskType       *string               `json:"base_block_storage_disk_type,omitempty"`
	BaseBlockStorageDiskDetailType *string               `json:"base_block_storage_disk_detail_type,omitempty"`
	InstanceTagList                []*server.InstanceTag `json:"tag_list,omitempty"`
	// VPC
	VpcNo                *string                           `json:"vpc_no,omitempty"`
	SubnetNo             *string                           `json:"subnet_no,omitempty"`
	InitScriptNo         *string                           `json:"init_script_no,omitempty"`
	PlacementGroupNo     *string                           `json:"placement_group_no,omitempty"`
	NetworkInterfaceList []*ServerInstanceNetworkInterface `json:"network_interface"`
}

// ServerInstanceNetworkInterface network interface model in server instance
type ServerInstanceNetworkInterface struct {
	Order              *int32  `json:"order,omitempty"`
	NetworkInterfaceNo *string `json:"network_interface_no,omitempty"`
	PrivateIp          *string `json:"private_ip,omitempty"`
	SubnetNo           *string `json:"subnet_no,omitempty"`
}
