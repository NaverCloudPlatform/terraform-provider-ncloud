package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"log"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	RegisterResource("ncloud_network_interface", resourceNcloudNetworkInterface())
}

func resourceNcloudNetworkInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudNetworkInterfaceCreate,
		Read:   resourceNcloudNetworkInterfaceRead,
		Update: resourceNcloudNetworkInterfaceUpdate,
		Delete: resourceNcloudNetworkInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"subnet_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: ToDiagFunc(validateInstanceName),
			},
			"private_ip": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: ToDiagFunc(validation.IsIPv4Address),
			},
			"access_control_groups": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"server_instance_no": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(0, 1000)),
			},
			"network_interface_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceNcloudNetworkInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	instance, err := createNetworkInterface(d, config)

	if err != nil {
		return err
	}

	d.SetId(*instance.NetworkInterfaceNo)
	log.Printf("[INFO] Network Interface ID: %s", d.Id())

	if v, ok := d.GetOk("server_instance_no"); ok && v != "" {
		if err := waitForNetworkInterfaceAttachment(config, d.Id()); err != nil {
			return err
		}
	}

	return resourceNcloudNetworkInterfaceRead(d, meta)
}

func resourceNcloudNetworkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	instance, err := getNetworkInterface(config, d.Id())
	if err != nil {
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*instance.NetworkInterfaceNo)
	d.Set("network_interface_no", instance.NetworkInterfaceNo)
	d.Set("name", instance.NetworkInterfaceName)
	d.Set("description", instance.NetworkInterfaceDescription)
	d.Set("subnet_no", instance.SubnetNo)
	d.Set("private_ip", instance.Ip)
	d.Set("server_instance_no", instance.InstanceNo)
	d.Set("status", instance.NetworkInterfaceStatus.Code)
	d.Set("access_control_groups", instance.AccessControlGroupNoList)
	d.Set("is_default", instance.IsDefault)

	if instance.InstanceType != nil {
		d.Set("instance_type", instance.InstanceType.Code)
	}

	return nil
}

func resourceNcloudNetworkInterfaceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if d.HasChange("server_instance_no") {
		o, n := d.GetChange("server_instance_no")
		if len(o.(string)) > 0 {
			if err := detachNetworkInterface(d, config, o.(string)); err != nil {
				return err
			}
		}

		if len(n.(string)) > 0 {
			if err := attachNetworkInterface(d, config); err != nil {
				return err
			}
		}
	}
	return resourceNcloudNetworkInterfaceRead(d, meta)
}

func resourceNcloudNetworkInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if err := deleteNetworkInterface(config, d.Id()); err != nil {
		return err
	}

	return nil
}

func getNetworkInterface(config *ProviderConfig, id string) (*vserver.NetworkInterface, error) {
	if config.SupportVPC {
		return getVpcNetworkInterface(config, id)
	}

	return nil, NotSupportClassic("resource `ncloud_network_interface`")
}

func getVpcNetworkInterface(config *ProviderConfig, id string) (*vserver.NetworkInterface, error) {
	reqParams := &vserver.GetNetworkInterfaceDetailRequest{
		RegionCode:         &config.RegionCode,
		NetworkInterfaceNo: ncloud.String(id),
	}

	logCommonRequest("getVpcNetworkInterface", reqParams)
	resp, err := config.Client.vserver.V2Api.GetNetworkInterfaceDetail(reqParams)
	if err != nil {
		logErrorResponse("getVpcNetworkInterface", err, reqParams)
		return nil, err
	}
	logResponse("getVpcNetworkInterface", resp)

	if len(resp.NetworkInterfaceList) > 0 {
		return resp.NetworkInterfaceList[0], nil
	}

	return nil, nil
}

func createNetworkInterface(d *schema.ResourceData, config *ProviderConfig) (*vserver.NetworkInterface, error) {
	if config.SupportVPC {
		return createVpcNetworkInterface(d, config)
	} else {
		return nil, NotSupportClassic("resource `ncloud_network_interface`")
	}

	if v, ok := d.GetOk("server_instance_no"); ok && v != "" {
		if err := waitForVpcNetworkInterfaceAttachment(config, d.Id()); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func createVpcNetworkInterface(d *schema.ResourceData, config *ProviderConfig) (*vserver.NetworkInterface, error) {
	subnet, err := getSubnetInstance(config, d.Get("subnet_no").(string))
	if err != nil {
		return nil, err
	}

	if subnet == nil {
		return nil, fmt.Errorf("subnet [%s] is not exist", d.Get("subnet_no"))
	}

	reqParams := &vserver.CreateNetworkInterfaceRequest{
		RegionCode:                  &config.RegionCode,
		AccessControlGroupNoList:    expandStringInterfaceList(d.Get("access_control_groups").([]interface{})),
		SubnetNo:                    ncloud.String(d.Get("subnet_no").(string)),
		VpcNo:                       subnet.VpcNo,
		NetworkInterfaceName:        StringPtrOrNil(d.GetOk("name")),
		NetworkInterfaceDescription: StringPtrOrNil(d.GetOk("description")),
		ServerInstanceNo:            StringPtrOrNil(d.GetOk("server_instance_no")),
		Ip:                          StringPtrOrNil(d.GetOk("private_ip")),
	}

	logCommonRequest("createVpcNetworkInterface", reqParams)
	resp, err := config.Client.vserver.V2Api.CreateNetworkInterface(reqParams)
	if err != nil {
		logErrorResponse("createVpcNetworkInterface", err, reqParams)
		return nil, err
	}
	logResponse("createVpcNetworkInterface", resp)

	return resp.NetworkInterfaceList[0], nil
}

func deleteNetworkInterface(config *ProviderConfig, id string) error {
	if config.SupportVPC {
		return deleteVpcNetworkInterface(config, id)
	}

	return NotSupportClassic("resource `ncloud_network_interface`")
}

func deleteVpcNetworkInterface(config *ProviderConfig, id string) error {
	reqParams := &vserver.DeleteNetworkInterfaceRequest{
		RegionCode:         &config.RegionCode,
		NetworkInterfaceNo: ncloud.String(id),
	}

	logCommonRequest("deleteVpcNetworkInterface", reqParams)
	resp, err := config.Client.vserver.V2Api.DeleteNetworkInterface(reqParams)
	if err != nil {
		logErrorResponse("deleteVpcNetworkInterface", err, reqParams)
		return err
	}
	logResponse("deleteVpcNetworkInterface", resp)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"USED", "NOTUSED", "UNSET"},
		Target:  []string{"TERMINATED"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getNetworkInterface(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "NetworkInterfaceStatus")
		},
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for Network Interface (%s) to become terminated: %s", id, err)
	}

	return nil
}

func attachNetworkInterface(d *schema.ResourceData, config *ProviderConfig) error {
	var err error

	if config.SupportVPC {
		err = attachVpcNetworkInterface(d, config)
	} else {
		err = NotSupportClassic("resource `ncloud_network_interface`")
	}

	if err != nil {
		return err
	}

	if err := waitForPublicIpDisassociation(config, d.Id()); err != nil {
		return err
	}

	return nil
}

func attachVpcNetworkInterface(d *schema.ResourceData, config *ProviderConfig) error {
	reqParams := &vserver.AttachNetworkInterfaceRequest{
		RegionCode:         &config.RegionCode,
		NetworkInterfaceNo: ncloud.String(d.Id()),
		SubnetNo:           ncloud.String(d.Get("subnet_no").(string)),
		ServerInstanceNo:   ncloud.String(d.Get("server_instance_no").(string)),
	}

	logCommonRequest("attachVpcNetworkInterface", reqParams)

	resp, err := config.Client.vserver.V2Api.AttachNetworkInterface(reqParams)
	if err != nil {
		logErrorResponse("attachVpcNetworkInterface", err, d.Id())
		return err
	}
	logCommonResponse("attachVpcNetworkInterface", GetCommonResponse(resp))

	stateConf := &resource.StateChangeConf{
		Pending: []string{"SET"},
		Target:  []string{"USED"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getNetworkInterface(config, d.Id())
			return VpcCommonStateRefreshFunc(instance, err, "NetworkInterfaceStatus")
		},
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for Network Interface (%s) to become attachmented: %s", d.Id(), err)
	}

	return nil
}

func detachNetworkInterface(d *schema.ResourceData, config *ProviderConfig, serverInstanceNo string) error {
	var err error

	if config.SupportVPC {
		err = detachVpcNetworkInterface(d, config, serverInstanceNo)
	} else {
		err = NotSupportClassic("resource `ncloud_network_interface`")
	}

	if err != nil {
		return err
	}

	return nil
}

func detachVpcNetworkInterface(d *schema.ResourceData, config *ProviderConfig, serverInstanceNo string) error {
	reqParams := &vserver.DetachNetworkInterfaceRequest{
		RegionCode:         &config.RegionCode,
		NetworkInterfaceNo: ncloud.String(d.Id()),
		SubnetNo:           ncloud.String(d.Get("subnet_no").(string)),
		ServerInstanceNo:   ncloud.String(serverInstanceNo),
	}

	logCommonRequest("detachVpcNetworkInterface", reqParams)

	resp, err := config.Client.vserver.V2Api.DetachNetworkInterface(reqParams)
	if err != nil {
		logErrorResponse("detachVpcNetworkInterface", err, d.Id())
		return err
	}
	logCommonResponse("detachVpcNetworkInterface", GetCommonResponse(resp))

	stateConf := &resource.StateChangeConf{
		Pending: []string{"UNSET"},
		Target:  []string{"NOTUSED"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getNetworkInterface(config, d.Id())
			return VpcCommonStateRefreshFunc(instance, err, "NetworkInterfaceStatus")
		},
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for Network Interface (%s) to become detachmented: %s", d.Id(), err)
	}

	return nil
}

func waitForNetworkInterfaceAttachment(config *ProviderConfig, id string) error {
	var err error

	if config.SupportVPC {
		err = waitForVpcNetworkInterfaceAttachment(config, id)
	} else {
		err = NotSupportClassic("resource `ncloud_network_interface`")
	}

	if err != nil {
		return err
	}

	return nil
}

func waitForVpcNetworkInterfaceAttachment(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"SET"},
		Target:  []string{"USED"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getNetworkInterface(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "NetworkInterfaceStatus")
		},
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for Network Interface (%s) to become attachmented: %s", id, err)
	}

	return nil
}
