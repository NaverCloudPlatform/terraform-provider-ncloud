package server

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/vpc"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

const (
	NetworkInterfaceStateNotUsed    = "NOTUSED"
	NetworkInterfaceStateUsed       = "USED"
	NetworkInterfaceStateSet        = "SET"
	NetworkInterfaceStateUnSet      = "UNSET"
	NetworkInterfaceStateTerminated = "TERMINATED"
)

func ResourceNcloudNetworkInterface() *schema.Resource {
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
				ValidateDiagFunc: ToDiagFunc(ValidateInstanceName),
			},
			"private_ip": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: ToDiagFunc(validation.IsIPv4Address),
			},
			"access_control_groups": {
				Type:     schema.TypeSet,
				Required: true,
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
	config := meta.(*conn.ProviderConfig)

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
	config := meta.(*conn.ProviderConfig)

	instance, err := GetNetworkInterface(config, d.Id())
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
	config := meta.(*conn.ProviderConfig)

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

	if d.HasChange("access_control_groups") {
		o, n := d.GetChange("access_control_groups")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		add := ns.Difference(os).List()
		remove := os.Difference(ns).List()

		removeAcgList := ExpandStringInterfaceList(remove)
		addAcgList := ExpandStringInterfaceList(add)

		// First do add ACG prevent error '[1002035] At least one Acg must remain on the network interface.'
		if len(addAcgList) > 0 {
			if err := addNetworkInterfaceAccessControlGroup(d, config, addAcgList); err != nil {
				return err
			}
		}

		if len(removeAcgList) > 0 {
			if err := removeNetworkInterfaceAccessControlGroup(d, config, removeAcgList); err != nil {
				return err
			}
		}
	}

	return resourceNcloudNetworkInterfaceRead(d, meta)
}

func removeNetworkInterfaceAccessControlGroup(d *schema.ResourceData, config *conn.ProviderConfig, accessControlGroupNoList []*string) error {
	var resp *vserver.RemoveNetworkInterfaceAccessControlGroupResponse
	var reqParams *vserver.RemoveNetworkInterfaceAccessControlGroupRequest

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		var err error
		reqParams = &vserver.RemoveNetworkInterfaceAccessControlGroupRequest{
			RegionCode:               &config.RegionCode,
			AccessControlGroupNoList: accessControlGroupNoList,
			NetworkInterfaceNo:       ncloud.String(d.Id()),
		}

		LogCommonRequest("RemoveNetworkInterfaceAccessControlGroup", reqParams)
		resp, err = config.Client.Vserver.V2Api.RemoveNetworkInterfaceAccessControlGroup(reqParams)

		if err != nil {
			errBody, _ := GetCommonErrorBody(err)
			if errBody.ReturnCode == ApiErrorNetworkInterfaceAtLeastOneAcgMustRemain {
				LogErrorResponse("retry RemoveNetworkInterfaceAccessControlGroup", err, reqParams)
				time.Sleep(time.Second * 5)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		LogErrorResponse("RemoveNetworkInterfaceAccessControlGroup", err, reqParams)
		return err
	}

	LogResponse("RemoveNetworkInterfaceAccessControlGroup", resp)

	if err = waitForVpcNetworkInterfaceState(config, d.Id(), []string{NetworkInterfaceStateSet}, []string{NetworkInterfaceStateNotUsed, NetworkInterfaceStateUsed}); err != nil {
		return err
	}

	return nil
}

func addNetworkInterfaceAccessControlGroup(d *schema.ResourceData, config *conn.ProviderConfig, accessControlGroupNoList []*string) error {
	reqParams := &vserver.AddNetworkInterfaceAccessControlGroupRequest{
		RegionCode:               &config.RegionCode,
		AccessControlGroupNoList: accessControlGroupNoList,
		NetworkInterfaceNo:       ncloud.String(d.Id()),
	}

	LogCommonRequest("AddNetworkInterfaceAccessControlGroup", reqParams)
	resp, err := config.Client.Vserver.V2Api.AddNetworkInterfaceAccessControlGroup(reqParams)

	if err != nil {
		LogErrorResponse("AddNetworkInterfaceAccessControlGroup", err, reqParams)
		return err
	}

	LogResponse("AddNetworkInterfaceAccessControlGroup", resp)

	if err = waitForVpcNetworkInterfaceState(config, d.Id(), []string{NetworkInterfaceStateSet}, []string{NetworkInterfaceStateNotUsed, NetworkInterfaceStateUsed}); err != nil {
		return err
	}

	return nil
}

func resourceNcloudNetworkInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	if err := DeleteNetworkInterface(config, d.Id()); err != nil {
		return err
	}

	return nil
}

func GetNetworkInterface(config *conn.ProviderConfig, id string) (*vserver.NetworkInterface, error) {
	if config.SupportVPC {
		return getVpcNetworkInterface(config, id)
	}

	return nil, NotSupportClassic("resource `ncloud_network_interface`")
}

func getVpcNetworkInterface(config *conn.ProviderConfig, id string) (*vserver.NetworkInterface, error) {
	reqParams := &vserver.GetNetworkInterfaceDetailRequest{
		RegionCode:         &config.RegionCode,
		NetworkInterfaceNo: ncloud.String(id),
	}

	LogCommonRequest("getVpcNetworkInterface", reqParams)
	resp, err := config.Client.Vserver.V2Api.GetNetworkInterfaceDetail(reqParams)
	if err != nil {
		LogErrorResponse("getVpcNetworkInterface", err, reqParams)
		return nil, err
	}
	LogResponse("getVpcNetworkInterface", resp)

	if len(resp.NetworkInterfaceList) > 0 {
		return resp.NetworkInterfaceList[0], nil
	}

	return nil, nil
}

func createNetworkInterface(d *schema.ResourceData, config *conn.ProviderConfig) (*vserver.NetworkInterface, error) {
	if config.SupportVPC {
		return createVpcNetworkInterface(d, config)
	} else {
		return nil, NotSupportClassic("resource `ncloud_network_interface`")
	}
}

func createVpcNetworkInterface(d *schema.ResourceData, config *conn.ProviderConfig) (*vserver.NetworkInterface, error) {
	subnet, err := vpc.GetSubnetInstance(config, d.Get("subnet_no").(string))
	if err != nil {
		return nil, err
	}

	if subnet == nil {
		return nil, fmt.Errorf("subnet [%s] is not exist", d.Get("subnet_no"))
	}

	reqParams := &vserver.CreateNetworkInterfaceRequest{
		RegionCode:                  &config.RegionCode,
		AccessControlGroupNoList:    ExpandStringInterfaceList(d.Get("access_control_groups").(*schema.Set).List()),
		SubnetNo:                    ncloud.String(d.Get("subnet_no").(string)),
		VpcNo:                       subnet.VpcNo,
		NetworkInterfaceName:        StringPtrOrNil(d.GetOk("name")),
		NetworkInterfaceDescription: StringPtrOrNil(d.GetOk("description")),
		ServerInstanceNo:            StringPtrOrNil(d.GetOk("server_instance_no")),
		Ip:                          StringPtrOrNil(d.GetOk("private_ip")),
	}

	LogCommonRequest("createVpcNetworkInterface", reqParams)
	resp, err := config.Client.Vserver.V2Api.CreateNetworkInterface(reqParams)
	if err != nil {
		LogErrorResponse("createVpcNetworkInterface", err, reqParams)
		return nil, err
	}
	LogResponse("createVpcNetworkInterface", resp)

	return resp.NetworkInterfaceList[0], nil
}

func DeleteNetworkInterface(config *conn.ProviderConfig, id string) error {
	if config.SupportVPC {
		return deleteVpcNetworkInterface(config, id)
	}

	return NotSupportClassic("resource `ncloud_network_interface`")
}

func deleteVpcNetworkInterface(config *conn.ProviderConfig, id string) error {
	reqParams := &vserver.DeleteNetworkInterfaceRequest{
		RegionCode:         &config.RegionCode,
		NetworkInterfaceNo: ncloud.String(id),
	}

	LogCommonRequest("deleteVpcNetworkInterface", reqParams)
	resp, err := config.Client.Vserver.V2Api.DeleteNetworkInterface(reqParams)
	if err != nil {
		LogErrorResponse("deleteVpcNetworkInterface", err, reqParams)
		return err
	}
	LogResponse("deleteVpcNetworkInterface", resp)

	if err := waitForVpcNetworkInterfaceState(config, id, []string{NetworkInterfaceStateUsed, NetworkInterfaceStateNotUsed, NetworkInterfaceStateUnSet}, []string{NetworkInterfaceStateTerminated}); err != nil {
		return err
	}

	return nil
}

func attachNetworkInterface(d *schema.ResourceData, config *conn.ProviderConfig) error {
	var err error

	if config.SupportVPC {
		err = attachVpcNetworkInterface(d, config)
	} else {
		err = NotSupportClassic("resource `ncloud_network_interface`")
	}

	if err != nil {
		return err
	}

	waitForPublicIpDisassociate(d, config)

	return nil
}

func attachVpcNetworkInterface(d *schema.ResourceData, config *conn.ProviderConfig) error {
	reqParams := &vserver.AttachNetworkInterfaceRequest{
		RegionCode:         &config.RegionCode,
		NetworkInterfaceNo: ncloud.String(d.Id()),
		SubnetNo:           ncloud.String(d.Get("subnet_no").(string)),
		ServerInstanceNo:   ncloud.String(d.Get("server_instance_no").(string)),
	}

	LogCommonRequest("attachVpcNetworkInterface", reqParams)

	resp, err := config.Client.Vserver.V2Api.AttachNetworkInterface(reqParams)
	if err != nil {
		LogErrorResponse("attachVpcNetworkInterface", err, d.Id())
		return err
	}
	LogCommonResponse("attachVpcNetworkInterface", GetCommonResponse(resp))

	if err := waitForNetworkInterfaceAttachment(config, d.Id()); err != nil {
		return err
	}

	return nil
}

func waitForPublicIpDisassociate(d *schema.ResourceData, config *conn.ProviderConfig) error {
	reqParams := &vserver.GetServerInstanceDetailRequest{
		RegionCode:         &config.RegionCode,
		ServerInstanceNo:   ncloud.String(d.Get("server_instance_no").(string)),
	}

	resp, err := config.Client.Vserver.V2Api.GetServerInstanceDetail(reqParams)
	if err != nil {
		return err
	}

	if err := ValidateOneResult(len(resp.ServerInstanceList)); err != nil {
		return err
	}

	if publicIpNo := *resp.ServerInstanceList[0].PublicIpInstanceNo; publicIpNo != "" {
		if err := waitForPublicIpDisassociation(config, publicIpNo); err != nil {
			return err
		}
	}

	return nil
}

func detachNetworkInterface(d *schema.ResourceData, config *conn.ProviderConfig, serverInstanceNo string) error {
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

func detachVpcNetworkInterface(d *schema.ResourceData, config *conn.ProviderConfig, serverInstanceNo string) error {
	reqParams := &vserver.DetachNetworkInterfaceRequest{
		RegionCode:         &config.RegionCode,
		NetworkInterfaceNo: ncloud.String(d.Id()),
		SubnetNo:           ncloud.String(d.Get("subnet_no").(string)),
		ServerInstanceNo:   ncloud.String(serverInstanceNo),
	}

	LogCommonRequest("detachVpcNetworkInterface", reqParams)

	resp, err := config.Client.Vserver.V2Api.DetachNetworkInterface(reqParams)
	if err != nil {
		LogErrorResponse("detachVpcNetworkInterface", err, d.Id())
		return err
	}
	LogCommonResponse("detachVpcNetworkInterface", GetCommonResponse(resp))

	if err := waitForVpcNetworkInterfaceState(config, d.Id(), []string{NetworkInterfaceStateUnSet}, []string{NetworkInterfaceStateNotUsed}); err != nil {
		return err
	}

	return nil
}

func waitForNetworkInterfaceAttachment(config *conn.ProviderConfig, id string) error {
	var err error

	if config.SupportVPC {
		err = waitForVpcNetworkInterfaceState(config, id, []string{NetworkInterfaceStateSet}, []string{NetworkInterfaceStateUsed})
	} else {
		err = NotSupportClassic("resource `ncloud_network_interface`")
	}

	if err != nil {
		return err
	}

	return nil
}

func waitForVpcNetworkInterfaceState(config *conn.ProviderConfig, id string, pending []string, target []string) error {
	stateConf := &resource.StateChangeConf{
		Pending: pending,
		Target:  target,
		Refresh: func() (interface{}, string, error) {
			instance, err := GetNetworkInterface(config, id)
			return vpc.VpcCommonStateRefreshFunc(instance, err, "NetworkInterfaceStatus")
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for Network Interface (%s) to become (%v): %s", id, target, err)
	}

	return nil
}
