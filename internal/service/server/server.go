package server

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/vpc"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/zone"
)

func ResourceNcloudServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudServerCreate,
		Read:   resourceNcloudServerRead,
		Update: resourceNcloudServerUpdate,
		Delete: resourceNcloudServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(conn.DefaultCreateTimeout),
			Delete: schema.DefaultTimeout(conn.DefaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			"server_image_product_code": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"member_server_image_no"},
			},
			"member_server_image_no": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"server_image_product_code"},
			},
			"server_image_number": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"server_image_product_code"},
			},
			"server_product_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"server_spec_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(
					validation.StringLenBetween(3, 30),
					validation.StringMatch(regexp.MustCompile(`^[a-z]+[a-z0-9-]+[a-z0-9]$`), "Allows only lowercase letters(a-z), numbers, hyphen (-). Must start with an alphabetic character, must end with an English letter or number"),
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
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"PUBLC", "GLBL"}, false)),
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
			"block_device_partition_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_point": {
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile(`^/(?:[a-z][a-z0-9]*)?$`), "Must start with an / character. Only lowercase English letters and numbers are allowed for names under /, and must start with a lowercase English letter.")),
						},
						"partition_size": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
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
			"hypervisor_type": {
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
	config := meta.(*conn.ProviderConfig)

	id, err := createServerInstance(d, config)

	if err != nil {
		return err
	}

	d.SetId(ncloud.StringValue(id))
	log.Printf("[INFO] Server instance ID: %s", d.Id())

	return resourceNcloudServerRead(d, meta)
}

func resourceNcloudServerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	r, err := GetServerInstance(config, d.Id())
	if err != nil {
		return err
	}

	if r == nil {
		d.SetId("")
		return nil
	}

	if config.SupportVPC {
		_ = buildNetworkInterfaceList(config, r)
	}

	instance := ConvertToMap(r)

	SetSingularResourceDataFromMapSchema(ResourceNcloudServer(), d, instance)

	return nil
}

func resourceNcloudServerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)
	serverInstance, err := GetServerInstance(config, d.Id())
	if err != nil {
		return err
	}

	if serverInstance == nil {
		d.SetId("")
		return nil
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

			if err := waitForDisconnectBlockStorage(config, *blockStorage.BlockStorageInstanceNo); err != nil {
				return err
			}
		}

		if err := detachThenWaitServerInstance(config, d.Id()); err != nil {
			return err
		}
	}

	if err := terminateThenWaitServerInstance(config, d.Id()); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceNcloudServerUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	if d.HasChange("server_product_code") || d.HasChange("server_spec_code") {
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

func createServerInstance(d *schema.ResourceData, config *conn.ProviderConfig) (*string, error) {
	if config.SupportVPC {
		return createVpcServerInstance(d, config)
	}

	return createClassicServerInstance(d, config)
}

func createClassicServerInstance(d *schema.ResourceData, config *conn.ProviderConfig) (*string, error) {
	zoneNo, err := zone.ParseZoneNoParameter(config, d)
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
		reqParams.AccessControlGroupConfigurationNoList = ExpandStringInterfaceList(param.([]interface{}))
	}

	var resp *server.CreateServerInstancesResponse
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		var err error
		LogCommonRequest("createClassicServerInstance", reqParams)
		resp, err = config.Client.Server.V2Api.CreateServerInstances(reqParams)
		if err != nil {
			errBody, _ := GetCommonErrorBody(err)
			if ContainsInStringList(errBody.ReturnCode, []string{ApiErrorUnknown, ApiErrorAuthorityParameter, ApiErrorServerObjectInOperation, ApiErrorPreviousServersHaveNotBeenEntirelyTerminated}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		LogResponse("createClassicServerInstance", resp)
		return nil
	})

	if err != nil {
		LogErrorResponse("createClassicServerInstance", err, reqParams)
		return nil, err
	}
	LogResponse("createClassicServerInstance", resp)

	serverInstance := resp.ServerInstanceList[0]

	if err := waitStateNcloudServerForCreation(config, *serverInstance.ServerInstanceNo); err != nil {
		return nil, err
	}

	return serverInstance.ServerInstanceNo, nil
}

func createVpcServerInstance(d *schema.ResourceData, config *conn.ProviderConfig) (*string, error) {
	if _, ok := d.GetOk("subnet_no"); !ok {
		return nil, ErrorRequiredArgOnVpc("subnet_no")
	}

	if _, ok := d.GetOk("access_control_group_configuration_no_list"); ok {
		return nil, NotSupportVpc("`access_control_group_configuration_no_list` of ncloud_server")
	}

	if _, ok := d.GetOk("user_data"); ok {
		return nil, NotSupportVpc("`user_data` of ncloud_server")
	}

	subnet, err := vpc.GetSubnetInstance(config, d.Get("subnet_no").(string))
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
		ServerImageNo:                     StringPtrOrNil(d.GetOk("server_image_number")),
		ServerSpecCode:                    StringPtrOrNil(d.GetOk("server_spec_code")),
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
		RaidTypeName:                      StringPtrOrNil(d.GetOk("raid_type_name")),
	}

	if networkInterfaceList, ok := d.GetOk("network_interface"); !ok {
		defaultAcgNo, err := vpc.GetDefaultAccessControlGroup(config, *subnet.VpcNo)
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

			networkInterface, err := GetNetworkInterface(config, networkInterfaceNo)
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

	if blockDevicePartitionList, err := expandBlockDevicePartitionListParams(d.Get("block_device_partition_list").([]interface{})); err == nil {
		reqParams.BlockDevicePartitionList = blockDevicePartitionList
	}

	LogCommonRequest("createVpcServerInstance", reqParams)
	resp, err := config.Client.Vserver.V2Api.CreateServerInstances(reqParams)
	if err != nil {
		LogErrorResponse("createVpcServerInstance", err, reqParams)
		return nil, err
	}
	LogResponse("createVpcServerInstance", resp)
	serverInstance := resp.ServerInstanceList[0]

	if err := waitStateNcloudServerForCreation(config, *serverInstance.ServerInstanceNo); err != nil {
		return nil, err
	}

	blockStorageList, err := getVpcBasicBlockStorageList(config, *serverInstance.ServerInstanceNo)
	if err != nil {
		return nil, err
	}

	if len(blockStorageList) > 0 {
		for _, blockStorage := range blockStorageList {
			if err := waitForAttachedBlockStorage(config, *blockStorage.BlockStorageInstanceNo); err != nil {
				return nil, err
			}
		}
	}

	return serverInstance.ServerInstanceNo, nil
}

func waitStateNcloudServerForCreation(config *conn.ProviderConfig, id string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"INIT", "CREAT"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetServerInstance(config, id)
			if err != nil {
				return 0, "", err
			}

			if instance == nil {
				return 0, "", fmt.Errorf("fail to get Server instance, %s doesn't exist", id)
			}

			return instance, ncloud.StringValue(instance.ServerInstanceStatus), nil
		},
		Timeout:    conn.DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for ServerInstance state to be \"RUN\": %s", err)
	}

	return nil
}

func updateServerInstanceSpec(d *schema.ResourceData, config *conn.ProviderConfig) error {
	serverInstance, err := GetServerInstance(config, d.Id())
	if err != nil {
		return err
	}

	if serverInstance == nil {
		return fmt.Errorf("fail to get Server instance, %s doesn't exist", d.Id())
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

func changeServerInstanceSpec(d *schema.ResourceData, config *conn.ProviderConfig) error {
	var err error
	if config.SupportVPC {
		err = changeVpcServerInstanceSpec(d, config)
	} else {
		err = changeClassicServerInstanceSpec(d, config)
	}

	if err != nil {
		return err
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"CHNG"},
		Target:  []string{"NULL"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetServerInstance(config, d.Id())

			if err != nil {
				return 0, "", err
			}

			if instance == nil {
				return 0, "", fmt.Errorf("fail to get Server instance, %s doesn't exist", d.Id())
			}

			return instance, ncloud.StringValue(instance.ServerInstanceOperation), nil
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for ServerInstance operation to be \"NULL\": %s", err)
	}

	return nil
}

func changeClassicServerInstanceSpec(d *schema.ResourceData, config *conn.ProviderConfig) error {
	reqParams := &server.ChangeServerInstanceSpecRequest{
		ServerInstanceNo:  ncloud.String(d.Get("instance_no").(string)),
		ServerProductCode: ncloud.String(d.Get("server_product_code").(string)),
	}

	LogCommonRequest("changeClassicServerInstanceSpec", reqParams)
	resp, err := config.Client.Server.V2Api.ChangeServerInstanceSpec(reqParams)
	if err != nil {
		LogErrorResponse("changeClassicServerInstanceSpec", err, reqParams)
		return err
	}
	LogCommonResponse("changeClassicServerInstanceSpec", GetCommonResponse(resp))

	return nil
}

func changeVpcServerInstanceSpec(d *schema.ResourceData, config *conn.ProviderConfig) error {
	reqParams := &vserver.ChangeServerInstanceSpecRequest{
		RegionCode:       &config.RegionCode,
		ServerInstanceNo: ncloud.String(d.Get("instance_no").(string)),
	}

	if d.HasChange("server_product_code") {
		reqParams.ServerProductCode = ncloud.String(d.Get("server_product_code").(string))
	}
	if d.HasChange("server_spec_code") {
		reqParams.ServerSpecCode = ncloud.String(d.Get("server_spec_code").(string))
	}

	LogCommonRequest("changeVpcServerInstanceSpec", reqParams)
	resp, err := config.Client.Vserver.V2Api.ChangeServerInstanceSpec(reqParams)
	if err != nil {
		LogErrorResponse("ChangeServerInstanceSpec", err, reqParams)
		return err
	}
	LogResponse("changeVpcServerInstanceSpec", resp)

	return nil
}

func updateServerProtectionTermination(d *schema.ResourceData, config *conn.ProviderConfig) error {
	if config.SupportVPC {
		return updateVpcServerProtectionTermination(d, config)
	}

	return updateClassicServerProtectionTermination(d, config)
}

func updateVpcServerProtectionTermination(d *schema.ResourceData, config *conn.ProviderConfig) error {
	reqParams := &vserver.SetProtectServerTerminationRequest{
		RegionCode:                 &config.RegionCode,
		ServerInstanceNo:           ncloud.String(d.Id()),
		IsProtectServerTermination: ncloud.Bool(d.Get("is_protect_server_termination").(bool)),
	}

	LogCommonRequest("SetProtectServerTermination", reqParams)
	resp, err := config.Client.Vserver.V2Api.SetProtectServerTermination(reqParams)
	if err != nil {
		LogErrorResponse("SetProtectServerTermination", err, reqParams)
		return err
	}
	LogResponse("SetProtectServerTermination", resp)

	return nil
}

func updateClassicServerProtectionTermination(d *schema.ResourceData, config *conn.ProviderConfig) error {
	reqParams := &server.SetProtectServerTerminationRequest{
		ServerInstanceNo:           ncloud.String(d.Id()),
		IsProtectServerTermination: ncloud.Bool(d.Get("is_protect_server_termination").(bool)),
	}

	LogCommonRequest("SetProtectServerTermination", reqParams)
	resp, err := config.Client.Server.V2Api.SetProtectServerTermination(reqParams)
	if err != nil {
		LogErrorResponse("SetProtectServerTermination", err, reqParams)
		return err
	}
	LogResponse("SetProtectServerTermination", resp)

	return nil
}

func startThenWaitServerInstance(config *conn.ProviderConfig, id string) error {
	var err error
	if config.SupportVPC {
		err = startVpcServerInstance(config, id)
	} else {
		err = startClassicServerInstance(config, id)
	}

	if err != nil {
		return err
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"NSTOP"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetServerInstance(config, id)
			if err != nil {
				return 0, "", err
			}

			if instance == nil {
				return 0, "", fmt.Errorf("fail to get Server instance, %s doesn't exist", id)
			}

			return instance, ncloud.StringValue(instance.ServerInstanceStatus), nil
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for ServerInstance state to be \"RUN\": %s", err)
	}

	return nil
}

func startClassicServerInstance(config *conn.ProviderConfig, id string) error {
	reqParams := &server.StartServerInstancesRequest{
		ServerInstanceNoList: []*string{ncloud.String(id)},
	}
	LogCommonRequest("startClassicServerInstance", reqParams)
	resp, err := config.Client.Server.V2Api.StartServerInstances(reqParams)
	if err != nil {
		LogErrorResponse("startClassicServerInstance", err, reqParams)
		return err
	}
	LogResponse("startClassicServerInstance", resp)

	return nil
}

func startVpcServerInstance(config *conn.ProviderConfig, id string) error {
	reqParams := &vserver.StartServerInstancesRequest{
		RegionCode:           &config.RegionCode,
		ServerInstanceNoList: []*string{ncloud.String(id)},
	}
	LogCommonRequest("startVpcServerInstance", reqParams)
	resp, err := config.Client.Vserver.V2Api.StartServerInstances(reqParams)
	if err != nil {
		LogErrorResponse("startVpcServerInstance", err, reqParams)
		return err
	}
	LogResponse("startVpcServerInstance", resp)

	return nil
}

func GetServerInstance(config *conn.ProviderConfig, id string) (*ServerInstance, error) {
	if config.SupportVPC {
		return getVpcServerInstance(config, id)
	}

	return getClassicServerInstance(config, id)
}

func getClassicServerInstance(config *conn.ProviderConfig, id string) (*ServerInstance, error) {
	reqParams := &server.GetServerInstanceListRequest{
		ServerInstanceNoList: []*string{ncloud.String(id)},
	}

	LogCommonRequest("getClassicServerInstance", reqParams)
	resp, err := config.Client.Server.V2Api.GetServerInstanceList(reqParams)

	if err != nil {
		LogErrorResponse("getClassicServerInstance", err, reqParams)
		return nil, err
	}

	LogResponse("getClassicServerInstance", resp)

	if len(resp.ServerInstanceList) == 0 {
		return nil, nil
	}

	if err := ValidateOneResult(len(resp.ServerInstanceList)); err != nil {
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
		ServerInstanceStatus:           common.GetCodePtrByCommonCode(r.ServerInstanceStatus),
		PlatformType:                   common.GetCodePtrByCommonCode(r.PlatformType),
		ServerInstanceOperation:        common.GetCodePtrByCommonCode(r.ServerInstanceOperation),
		Zone:                           r.Zone.ZoneCode,
		BaseBlockStorageDiskType:       common.GetCodePtrByCommonCode(r.BaseBlockStorageDiskType),
		BaseBlockStorageDiskDetailType: flattenMapByKey(r.BaseBlockStorageDiskDetailType, "code"),
		InstanceTagList:                r.InstanceTagList,
	}
}

func getVpcServerInstance(config *conn.ProviderConfig, id string) (*ServerInstance, error) {
	reqParams := &vserver.GetServerInstanceDetailRequest{
		RegionCode:       &config.RegionCode,
		ServerInstanceNo: ncloud.String(id),
	}

	LogCommonRequest("getVpcServerInstance", reqParams)
	resp, err := config.Client.Vserver.V2Api.GetServerInstanceDetail(reqParams)

	if err != nil {
		LogErrorResponse("getVpcServerInstance", err, reqParams)
		return nil, err
	}

	LogResponse("getVpcServerInstance", resp)

	if resp == nil || len(resp.ServerInstanceList) == 0 {
		return nil, nil
	}

	if err := ValidateOneResult(len(resp.ServerInstanceList)); err != nil {
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
		ServerImageNo:                  r.ServerImageNo,
		ServerProductCode:              r.ServerProductCode,
		ServerSpecCode:                 r.ServerSpecCode,
		ServerName:                     r.ServerName,
		ServerDescription:              r.ServerDescription,
		LoginKeyName:                   r.LoginKeyName,
		IsProtectServerTermination:     r.IsProtectServerTermination,
		ServerInstanceNo:               r.ServerInstanceNo,
		CpuCount:                       r.CpuCount,
		MemorySize:                     r.MemorySize,
		PublicIp:                       r.PublicIp,
		ServerInstanceStatus:           common.GetCodePtrByCommonCode(r.ServerInstanceStatus),
		PlatformType:                   common.GetCodePtrByCommonCode(r.PlatformType),
		ServerInstanceOperation:        common.GetCodePtrByCommonCode(r.ServerInstanceOperation),
		Zone:                           r.ZoneCode,
		BaseBlockStorageDiskType:       common.GetCodePtrByCommonCode(r.BaseBlockStorageDiskType),
		BaseBlockStorageDiskDetailType: flattenMapByKey(r.BaseBlockStorageDiskDetailType, "code"),
		VpcNo:                          r.VpcNo,
		SubnetNo:                       r.SubnetNo,
		InitScriptNo:                   r.InitScriptNo,
		PlacementGroupNo:               r.PlacementGroupNo,
		HypervisorType:                 common.GetCodePtrByCommonCode(r.HypervisorType),
		BlockDevicePartitionList:       r.BlockDevicePartitionList,
	}

	for _, networkInterfaceNo := range r.NetworkInterfaceNoList {
		ni := &ServerInstanceNetworkInterface{
			NetworkInterfaceNo: networkInterfaceNo,
		}

		instance.NetworkInterfaceList = append(instance.NetworkInterfaceList, ni)
	}

	return instance
}

func buildNetworkInterfaceList(config *conn.ProviderConfig, r *ServerInstance) error {
	for _, ni := range r.NetworkInterfaceList {
		networkInterface, err := GetNetworkInterface(config, *ni.NetworkInterfaceNo)

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

func stopThenWaitServerInstance(config *conn.ProviderConfig, id string) error {
	var err error

	stateConf := &retry.StateChangeConf{
		Pending: []string{"SETUP"},
		Target:  []string{"NULL"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetServerInstance(config, id)
			if err != nil {
				return 0, "", err
			}

			if instance == nil {
				return &server.ServerInstance{}, "NULL", nil
			}

			return instance, ncloud.StringValue(instance.ServerInstanceOperation), nil
		},
		Timeout:    conn.DefaultStopTimeout,
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

	stateConf = &retry.StateChangeConf{
		Pending: []string{"RUN"},
		Target:  []string{"NSTOP"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetServerInstance(config, id)
			if err != nil {
				return 0, "", err
			}

			if instance == nil {
				return &server.ServerInstance{}, "NULL", nil
			}

			return instance, ncloud.StringValue(instance.ServerInstanceStatus), nil
		},
		Timeout:    conn.DefaultStopTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for ServerInstance state to be \"NSTOP\": %s", err)
	}

	return nil
}

func detachThenWaitServerInstance(config *conn.ProviderConfig, id string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"SETUP"},
		Target:  []string{"NULL"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetServerInstance(config, id)
			if err != nil {
				return 0, "", err
			}

			// FIXME: When deleting a server if user detach block storage what they attached by themself
			//        and keep that block storage alive
			// 1. during the server deletion process, block storage detached
			// 2. attempt to detach from the server during the block storage in-place update process
			// 2-1. but the server is already destroyed and detachThenWaitServerInstance() is called
			//      by block storage to get serverInstance info, and the instance inquiry result is nil,
			//      causing panic when access ServerInstance(nil) field.
			if instance == nil {
				return &server.ServerInstance{}, "NULL", nil
			}

			return instance, ncloud.StringValue(instance.ServerInstanceOperation), nil
		},
		Timeout:    conn.DefaultStopTimeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for ServerInstance operation to be \"NULL\": %s", err)
	}

	return nil
}

func stopClassicServerInstance(config *conn.ProviderConfig, id string) error {
	reqParams := &server.StopServerInstancesRequest{
		ServerInstanceNoList: []*string{ncloud.String(id)},
	}
	LogCommonRequest("stopClassicServerInstance", reqParams)
	resp, err := config.Client.Server.V2Api.StopServerInstances(reqParams)
	if err != nil {
		LogErrorResponse("stopClassicServerInstance", err, reqParams)
		return err
	}
	LogResponse("stopClassicServerInstance", resp)

	return nil
}

func stopVpcServerInstance(config *conn.ProviderConfig, id string) error {
	reqParams := &vserver.StopServerInstancesRequest{
		RegionCode:           &config.RegionCode,
		ServerInstanceNoList: []*string{ncloud.String(id)},
	}
	LogCommonRequest("stopVpcServerInstance", reqParams)
	resp, err := config.Client.Vserver.V2Api.StopServerInstances(reqParams)
	if err != nil {
		LogErrorResponse("stopVpcServerInstance", err, reqParams)
		return err
	}
	LogResponse("stopVpcServerInstance", resp)

	return nil
}

func terminateThenWaitServerInstance(config *conn.ProviderConfig, id string) error {
	var err error
	if config.SupportVPC {
		err = terminateVpcServerInstance(config, id)
	} else {
		err = terminateClassicServerInstance(config, id)
	}

	if err != nil {
		return err
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"NSTOP"},
		Target:  []string{"TERMINATED"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetServerInstance(config, id)

			if err != nil {
				return 0, "", err
			}
			if instance == nil { // Instance is terminated.
				return instance, "TERMINATED", nil
			}
			return instance, ncloud.StringValue(instance.ServerInstanceStatus), nil
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for ServerInstance state to be \"TERMINATED\": %s", err)
	}

	return nil
}

func terminateClassicServerInstance(config *conn.ProviderConfig, id string) error {
	reqParams := &server.TerminateServerInstancesRequest{
		ServerInstanceNoList: []*string{ncloud.String(id)},
	}

	var resp *server.TerminateServerInstancesResponse
	err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		var err error
		LogCommonRequest("terminateClassicServerInstance", reqParams)
		resp, err = config.Client.Server.V2Api.TerminateServerInstances(reqParams)
		if err != nil {
			errBody, _ := GetCommonErrorBody(err)
			if ContainsInStringList(errBody.ReturnCode, []string{ApiErrorUnknown, ApiErrorServerObjectInOperation2}) {
				LogErrorResponse("retry terminateClassicServerInstance", err, reqParams)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		LogResponse("terminateClassicServerInstance", resp)
		return nil
	})

	if err != nil {
		LogErrorResponse("terminateClassicServerInstance", err, reqParams)
		return err
	}

	return nil
}

func terminateVpcServerInstance(config *conn.ProviderConfig, id string) error {
	reqParams := &vserver.TerminateServerInstancesRequest{
		RegionCode:           &config.RegionCode,
		ServerInstanceNoList: []*string{ncloud.String(id)},
	}

	LogCommonRequest("terminateVpcServerInstance", reqParams)
	resp, err := config.Client.Vserver.V2Api.TerminateServerInstances(reqParams)
	LogResponse("terminateVpcServerInstance", resp)

	if err != nil {
		LogErrorResponse("terminateVpcServerInstance", err, reqParams)
		return err
	}

	return nil
}

func getAdditionalBlockStorageList(config *conn.ProviderConfig, id string) ([]*BlockStorage, error) {
	if config.SupportVPC {
		return getVpcAdditionalBlockStorageList(config, id)
	} else {
		return getClassicAdditionalBlockStorageList(config, id)
	}
}

func getVpcAdditionalBlockStorageList(config *conn.ProviderConfig, id string) ([]*BlockStorage, error) {
	resp, err := config.Client.Vserver.V2Api.GetBlockStorageInstanceList(&vserver.GetBlockStorageInstanceListRequest{
		RegionCode:               &config.RegionCode,
		ServerInstanceNo:         ncloud.String(id),
		BlockStorageTypeCodeList: []*string{ncloud.String("SVRBS")},
	})

	if err != nil {
		return nil, err
	}

	LogResponse("getVpcAdditionalBlockStorageList", resp)

	if len(resp.BlockStorageInstanceList) < 1 {
		return nil, nil
	}

	blockStorageList := make([]*BlockStorage, 0)
	for _, blockStorage := range resp.BlockStorageInstanceList {
		blockStorageList = append(blockStorageList, convertVpcBlockStorage(blockStorage))
	}

	return blockStorageList, nil
}

func getVpcBasicBlockStorageList(config *conn.ProviderConfig, id string) ([]*BlockStorage, error) {
	resp, err := config.Client.Vserver.V2Api.GetBlockStorageInstanceList(&vserver.GetBlockStorageInstanceListRequest{
		RegionCode:       &config.RegionCode,
		ServerInstanceNo: ncloud.String(id),
	})

	if err != nil {
		return nil, err
	}

	LogResponse("getVpcBasicBlockStorageList", resp)

	if len(resp.BlockStorageInstanceList) < 1 {
		return nil, nil
	}

	blockStorageList := make([]*BlockStorage, 0)
	for _, blockStorage := range resp.BlockStorageInstanceList {
		blockStorageList = append(blockStorageList, convertVpcBlockStorage(blockStorage))
	}

	return blockStorageList, nil
}

func getClassicAdditionalBlockStorageList(config *conn.ProviderConfig, id string) ([]*BlockStorage, error) {
	resp, err := config.Client.Server.V2Api.GetBlockStorageInstanceList(&server.GetBlockStorageInstanceListRequest{
		RegionNo:                 &config.RegionCode,
		ServerInstanceNo:         ncloud.String(id),
		BlockStorageTypeCodeList: []*string{ncloud.String("SVRBS")},
	})

	if err != nil {
		return nil, err
	}

	LogResponse("getClassicAdditionalBlockStorageList", resp)

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
		BlockStorageType:        common.GetCodePtrByCommonCode(storage.BlockStorageType),
		BlockStorageName:        storage.BlockStorageName,
		BlockStorageSize:        storage.BlockStorageSize,
		DeviceName:              storage.DeviceName,
		BlockStorageProductCode: storage.BlockStorageProductCode,
		Status:                  common.GetCodePtrByCommonCode(storage.BlockStorageInstanceStatus),
		StatusName:              storage.BlockStorageInstanceStatusName,
		Description:             storage.BlockStorageDescription,
		DiskType:                common.GetCodePtrByCommonCode(storage.BlockStorageDiskType),
		DiskDetailType:          common.GetCodePtrByCommonCode(storage.BlockStorageDiskDetailType),
		ZoneCode:                storage.ZoneCode,
	}
}

func convertClassicBlockStorage(storage *server.BlockStorageInstance) *BlockStorage {
	return &BlockStorage{
		BlockStorageInstanceNo:  storage.BlockStorageInstanceNo,
		ServerInstanceNo:        storage.ServerInstanceNo,
		ServerName:              storage.ServerName,
		BlockStorageType:        common.GetCodePtrByCommonCode(storage.BlockStorageType),
		BlockStorageName:        storage.BlockStorageName,
		BlockStorageSize:        storage.BlockStorageSize,
		DeviceName:              storage.DeviceName,
		BlockStorageProductCode: storage.BlockStorageProductCode,
		Status:                  common.GetCodePtrByCommonCode(storage.BlockStorageInstanceStatus),
		Operation:               common.GetCodePtrByCommonCode(storage.BlockStorageInstanceOperation),
		Description:             storage.BlockStorageInstanceDescription,
		DiskType:                common.GetCodePtrByCommonCode(storage.DiskType),
		DiskDetailType:          common.GetCodePtrByCommonCode(storage.DiskDetailType),
		ZoneCode:                storage.Zone.ZoneCode,
	}
}

func disconnectBlockStorage(config *conn.ProviderConfig, storage *BlockStorage) error {
	if config.SupportVPC {
		return disconnectVpcBlockStorage(config, storage)
	} else {
		return disconnectClassicBlockStorage(config, storage)
	}
}

func disconnectVpcBlockStorage(config *conn.ProviderConfig, storage *BlockStorage) error {
	_, err := config.Client.Vserver.V2Api.DetachBlockStorageInstances(&vserver.DetachBlockStorageInstancesRequest{
		RegionCode:                 &config.RegionCode,
		BlockStorageInstanceNoList: []*string{storage.BlockStorageInstanceNo},
	})

	if err != nil {
		return err
	}

	return nil
}

func disconnectClassicBlockStorage(config *conn.ProviderConfig, storage *BlockStorage) error {
	_, err := config.Client.Server.V2Api.DetachBlockStorageInstances(&server.DetachBlockStorageInstancesRequest{
		BlockStorageInstanceNoList: []*string{storage.BlockStorageInstanceNo},
	})

	if err != nil {
		return err
	}

	return nil
}

func waitForDisconnectBlockStorage(config *conn.ProviderConfig, no string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{BlockStorageStatusNameAttach},
		Target:  []string{BlockStorageStatusNameDetach},
		Refresh: func() (interface{}, string, error) {
			resp, err := GetBlockStorage(config, no)
			if err != nil {
				return 0, "", err
			}

			if resp == nil {
				return 0, "", fmt.Errorf("fail to get BlockStorage instance, %s doesn't exist", no)
			}

			if *resp.StatusName == BlockStorageStatusNameAttach {
				return resp, BlockStorageStatusNameAttach, nil
			} else if *resp.StatusName == BlockStorageStatusNameDetach {
				return resp, BlockStorageStatusNameDetach, nil
			}

			return 0, "", fmt.Errorf("error occurred while waiting to detached")
		},
		Timeout:    6 * conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for BlockStorage (%s) to become available: %s", no, err)
	}

	return nil
}

func waitForAttachedBlockStorage(config *conn.ProviderConfig, no string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{BlockStorageStatusNameInit, BlockStorageStatusNameOptimizing},
		Target:  []string{BlockStorageStatusNameAttach},
		Refresh: func() (interface{}, string, error) {
			resp, err := GetBlockStorage(config, no)
			if err != nil {
				return 0, "", err
			}

			if resp == nil {
				return 0, "", fmt.Errorf("fail to get BlockStorage instance, %s doesn't exist", no)
			}

			if *resp.StatusName == BlockStorageStatusNameInit {
				return resp, BlockStorageStatusNameInit, nil
			} else if *resp.StatusName == BlockStorageStatusNameOptimizing {
				return resp, BlockStorageStatusNameOptimizing, nil
			} else if *resp.StatusName == BlockStorageStatusNameAttach {
				return resp, BlockStorageStatusNameAttach, nil
			}

			return 0, "", fmt.Errorf("error occurred while waiting to attached")
		},
		Timeout:    6 * conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for BlockStorage (%s) to become available: %s", no, err)
	}

	return nil
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
	ServerImageNo            *string                           `json:"server_image_number,omitempty"`
	ServerSpecCode           *string                           `json:"server_spec_code,omitempty"`
	HypervisorType           *string                           `json:"hypervisor_type,omitempty"`
	VpcNo                    *string                           `json:"vpc_no,omitempty"`
	SubnetNo                 *string                           `json:"subnet_no,omitempty"`
	InitScriptNo             *string                           `json:"init_script_no,omitempty"`
	PlacementGroupNo         *string                           `json:"placement_group_no,omitempty"`
	NetworkInterfaceList     []*ServerInstanceNetworkInterface `json:"network_interface"`
	BlockDevicePartitionList []*vserver.BlockDevicePartition   `json:"block_device_partition_list,omitempty"`
}

// ServerInstanceNetworkInterface network interface model in server instance
type ServerInstanceNetworkInterface struct {
	Order              *int32  `json:"order,omitempty"`
	NetworkInterfaceNo *string `json:"network_interface_no,omitempty"`
	PrivateIp          *string `json:"private_ip,omitempty"`
	SubnetNo           *string `json:"subnet_no,omitempty"`
}
