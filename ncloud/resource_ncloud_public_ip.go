package ncloud

import (
	"fmt"
	"log"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	RegisterResource("ncloud_public_ip", resourceNcloudPublicIpInstance())
}

func resourceNcloudPublicIpInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudPublicIpCreate,
		Read:   resourceNcloudPublicIpRead,
		Delete: resourceNcloudPublicIpDelete,
		Update: resourceNcloudPublicIpUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: resourceNcloudPublicIpCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"server_instance_no": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 10000),
			},

			"public_ip_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internet_line_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PUBLC", "GLBL"}, false),
			},
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kind_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_no": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Use 'public_ip_no' instead",
			},
		},
	}
}

func resourceNcloudPublicIpCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	var publicIpInstanceNo *string
	var err error

	if config.SupportVPC {
		publicIpInstanceNo, err = createVpcPublicIp(d, config)
	} else {
		publicIpInstanceNo, err = createClassicPublicIp(d, config)
	}

	if err != nil {
		return err
	}

	d.SetId(ncloud.StringValue(publicIpInstanceNo))
	log.Printf("[INFO] Public IP ID: %s", d.Id())

	if v, ok := d.GetOk("server_instance_no"); ok && v != "" {
		if err := waitForPublicIpAssociation(config, d.Id()); err != nil {
			return err
		}
	}

	return resourceNcloudPublicIpRead(d, meta)
}

func resourceNcloudPublicIpRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	resource, err := getPublicIp(config, d.Id())

	if err != nil {
		return err
	}

	if resource == nil {
		d.SetId("")
		return nil
	}

	SetSingularResourceDataFromMap(d, resource)

	return nil
}

func resourceNcloudPublicIpDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	var err error

	// Check associated public ip
	if associated, err := checkAssociatedPublicIP(config, d.Id()); associated {
		// if associated public ip, disassociated the public ip
		if err := disassociatedPublicIp(config, d.Id()); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	if config.SupportVPC {
		err = deleteVpcPublicIp(d, config)
	} else {
		err = deleteClassicPublicIp(d, config)
	}

	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceNcloudPublicIpUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if d.HasChange("server_instance_no") {
		o, n := d.GetChange("server_instance_no")
		if len(o.(string)) > 0 {
			if err := disassociatedPublicIp(config, d.Id()); err != nil {
				return err
			}
		}

		if len(n.(string)) > 0 {
			if err := associatedPublicIp(d, config); err != nil {
				return err
			}
		}
	}

	return resourceNcloudPublicIpRead(d, meta)
}

func createClassicPublicIp(d *schema.ResourceData, config *ProviderConfig) (*string, error) {
	client := config.Client

	zoneNo, err := parseZoneNoParameter(config, d)
	if err != nil {
		return nil, err
	}

	reqParams := &server.CreatePublicIpInstanceRequest{
		RegionNo:             &config.RegionNo,
		ZoneNo:               zoneNo,
		InternetLineTypeCode: StringPtrOrNil(d.GetOk("internet_line_type")),
		ServerInstanceNo:     StringPtrOrNil(d.GetOk("server_instance_no")),
		PublicIpDescription:  StringPtrOrNil(d.GetOk("description")),
	}

	logCommonRequest("createClassicPublicIp", reqParams)

	resp, err := client.server.V2Api.CreatePublicIpInstance(reqParams)
	if err != nil {
		logErrorResponse("createClassicPublicIp", err, reqParams)
		return nil, err
	}
	logResponse("createClassicPublicIp", resp)

	publicIPInstance := resp.PublicIpInstanceList[0]

	return publicIPInstance.PublicIpInstanceNo, nil
}

func createVpcPublicIp(d *schema.ResourceData, config *ProviderConfig) (*string, error) {
	client := config.Client

	reqParams := &vserver.CreatePublicIpInstanceRequest{
		RegionCode:          &config.RegionCode,
		ServerInstanceNo:    StringPtrOrNil(d.GetOk("server_instance_no")),
		PublicIpDescription: StringPtrOrNil(d.GetOk("description")),
	}

	logCommonRequest("createVpcPublicIp", reqParams)

	resp, err := client.vserver.V2Api.CreatePublicIpInstance(reqParams)
	if err != nil {
		logErrorResponse("createVpcPublicIp", err, reqParams)
		return nil, err
	}
	logResponse("createVpcPublicIp", resp)

	publicIPInstance := resp.PublicIpInstanceList[0]

	return publicIPInstance.PublicIpInstanceNo, nil
}

func deleteClassicPublicIp(d *schema.ResourceData, config *ProviderConfig) error {
	client := config.Client

	reqParams := &server.DeletePublicIpInstancesRequest{
		PublicIpInstanceNoList: []*string{ncloud.String(d.Id())},
	}

	logCommonRequest("deleteClassicPublicIp", reqParams)

	resp, err := client.server.V2Api.DeletePublicIpInstances(reqParams)
	if err != nil {
		logErrorResponse("deleteClassicPublicIp", err, reqParams)
		return err
	}
	logResponse("deleteClassicPublicIp", resp)

	return nil
}

func deleteVpcPublicIp(d *schema.ResourceData, config *ProviderConfig) error {
	client := config.Client

	reqParams := &vserver.DeletePublicIpInstanceRequest{
		RegionCode:         &config.RegionCode,
		PublicIpInstanceNo: ncloud.String(d.Id()),
	}

	logCommonRequest("deleteVpcPublicIp", reqParams)

	resp, err := client.vserver.V2Api.DeletePublicIpInstance(reqParams)
	if err != nil {
		logErrorResponse("deleteVpcPublicIp", err, reqParams)
		return err
	}
	logResponse("deleteVpcPublicIp", resp)

	return nil
}

func getPublicIp(config *ProviderConfig, id string) (map[string]interface{}, error) {
	var resource map[string]interface{}
	var err error

	if config.SupportVPC {
		resource, err = getVpcPublicIp(config, id)
	} else {
		resource, err = getClassicPublicIp(config, id)
	}

	if err != nil {
		return nil, err
	}

	return resource, nil
}

func getClassicPublicIp(config *ProviderConfig, id string) (map[string]interface{}, error) {
	client := config.Client
	regionNo := config.RegionNo

	reqParams := &server.GetPublicIpInstanceListRequest{
		RegionNo:               &regionNo,
		PublicIpInstanceNoList: []*string{ncloud.String(id)},
	}

	logCommonRequest("getClassicPublicIp", reqParams)
	resp, err := client.server.V2Api.GetPublicIpInstanceList(reqParams)

	if err != nil {
		logErrorResponse("getClassicPublicIp", err, reqParams)
		return nil, err
	}
	logResponse("getClassicPublicIp", resp)

	if len(resp.PublicIpInstanceList) == 0 {
		return nil, nil
	}

	if err := validateOneResult(len(resp.PublicIpInstanceList)); err != nil {
		return nil, err
	}

	r := resp.PublicIpInstanceList[0]

	instance := map[string]interface{}{
		"id":                 *r.PublicIpInstanceNo,
		"public_ip_no":       *r.PublicIpInstanceNo,
		"public_ip":          *r.PublicIp,
		"description":        *r.PublicIpDescription,
		"zone":               *r.Zone.ZoneCode,
		"instance_no":        *r.PublicIpInstanceNo, // Deprecated
		"server_instance_no": nil,
	}

	if m := flattenCommonCode(r.InternetLineType); m["code"] != nil {
		instance["internet_line_type"] = m["code"]
	}

	if m := flattenCommonCode(r.PublicIpInstanceStatus); m["code"] != nil {
		instance["status"] = m["code"]
	}

	if m := flattenCommonCode(r.PublicIpKindType); m["code"] != nil {
		instance["kind_type"] = m["code"]
	}

	if r.ServerInstanceAssociatedWithPublicIp != nil {
		SetStringIfNotNilAndEmpty(instance, "server_instance_no", r.ServerInstanceAssociatedWithPublicIp.ServerInstanceNo)
	}

	return instance, nil
}

func getVpcPublicIp(config *ProviderConfig, id string) (map[string]interface{}, error) {
	client := config.Client
	regionCode := config.RegionCode

	reqParams := &vserver.GetPublicIpInstanceListRequest{
		RegionCode:             &regionCode,
		PublicIpInstanceNoList: []*string{ncloud.String(id)},
	}

	logCommonRequest("getVpcPublicIp", reqParams)
	resp, err := client.vserver.V2Api.GetPublicIpInstanceList(reqParams)

	if err != nil {
		logErrorResponse("getVpcPublicIp", err, reqParams)
		return nil, err
	}
	logResponse("getVpcPublicIp", resp)

	if len(resp.PublicIpInstanceList) == 0 {
		return nil, nil
	}

	if err := validateOneResult(len(resp.PublicIpInstanceList)); err != nil {
		return nil, err
	}

	r := resp.PublicIpInstanceList[0]

	instance := map[string]interface{}{
		"id":                 *r.PublicIpInstanceNo,
		"public_ip_no":       *r.PublicIpInstanceNo,
		"public_ip":          *r.PublicIp,
		"description":        *r.PublicIpDescription,
		"server_instance_no": nil,
	}

	SetStringIfNotNilAndEmpty(instance, "server_instance_no", r.ServerInstanceNo)

	if m := flattenCommonCode(r.PublicIpInstanceStatus); m["code"] != nil {
		instance["status"] = m["code"]
	}

	return instance, nil
}

func checkAssociatedPublicIP(config *ProviderConfig, id string) (bool, error) {
	instance, err := getPublicIp(config, id)

	if err != nil {
		return false, err
	}

	return instance["server_instance_no"] != nil && instance["server_instance_no"] != "", nil
}

func disassociatedPublicIp(config *ProviderConfig, id string) error {
	var err error

	if config.SupportVPC {
		err = disassociatedVpcPublicIp(config, id)
	} else {
		err = disassociatedClassicPublicIp(config, id)
	}

	if err != nil {
		return err
	}

	if err := waitForPublicIpDisassociation(config, id); err != nil {
		return err
	}

	return nil
}

func disassociatedClassicPublicIp(config *ProviderConfig, id string) error {
	reqParams := &server.DisassociatePublicIpFromServerInstanceRequest{PublicIpInstanceNo: ncloud.String(id)}

	logCommonRequest("disassociatedClassicPublicIP", reqParams)

	resp, err := config.Client.server.V2Api.DisassociatePublicIpFromServerInstance(reqParams)
	if err != nil {
		logErrorResponse("disassociatedClassicPublicIP", err, id)
		return err
	}
	logCommonResponse("disassociatedClassicPublicIP", GetCommonResponse(resp))

	return nil
}

func disassociatedVpcPublicIp(config *ProviderConfig, id string) error {
	reqParams := &vserver.DisassociatePublicIpFromServerInstanceRequest{PublicIpInstanceNo: ncloud.String(id)}

	logCommonRequest("disassociatedVpcPublicIp", reqParams)

	resp, err := config.Client.vserver.V2Api.DisassociatePublicIpFromServerInstance(reqParams)
	if err != nil {
		logErrorResponse("disassociatedVpcPublicIp", err, id)
		return err
	}
	logCommonResponse("disassociatedVpcPublicIp", GetCommonResponse(resp))

	return nil
}

func waitForPublicIpDisassociation(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"NOT OK"},
		Target:  []string{"OK"},
		Refresh: func() (interface{}, string, error) {
			isAssociated, err := checkAssociatedPublicIP(config, id)

			if err != nil {
				return 0, "", err
			}

			if !isAssociated {
				return 0, "OK", nil
			}

			return nil, "NOT OK", nil
		},
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Public IP (%s) to become disassociation: %s", id, err)
	}

	return nil
}

func waitForPublicIpAssociation(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"NOT OK"},
		Target:  []string{"OK"},
		Refresh: func() (interface{}, string, error) {
			isAssociated, err := checkAssociatedPublicIP(config, id)

			if err != nil {
				return 0, "", err
			}

			if isAssociated {
				return 0, "OK", nil
			}

			return nil, "NOT OK", nil
		},
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Public IP (%s) to become association: %s", id, err)
	}

	return nil
}

func associatedPublicIp(d *schema.ResourceData, config *ProviderConfig) error {
	var err error

	if config.SupportVPC {
		err = associatedVpcPublicIp(d, config)
	} else {
		err = associatedClassicPublicIp(d, config)
	}

	if err != nil {
		return err
	}

	if err := waitForPublicIpAssociation(config, d.Id()); err != nil {
		return err
	}

	return nil
}

func associatedClassicPublicIp(d *schema.ResourceData, config *ProviderConfig) error {
	reqParams := &server.AssociatePublicIpWithServerInstanceRequest{
		PublicIpInstanceNo: ncloud.String(d.Id()),
		ServerInstanceNo:   ncloud.String(d.Get("server_instance_no").(string)),
	}

	logCommonRequest("associatedClassicPublicIp", reqParams)

	resp, err := config.Client.server.V2Api.AssociatePublicIpWithServerInstance(reqParams)
	if err != nil {
		logErrorResponse("associatedClassicPublicIp", err, d.Id())
		return err
	}
	logCommonResponse("associatedClassicPublicIp", GetCommonResponse(resp))

	return nil
}

func associatedVpcPublicIp(d *schema.ResourceData, config *ProviderConfig) error {
	reqParams := &vserver.AssociatePublicIpWithServerInstanceRequest{
		RegionCode:         &config.RegionCode,
		PublicIpInstanceNo: ncloud.String(d.Id()),
		ServerInstanceNo:   ncloud.String(d.Get("server_instance_no").(string)),
	}

	logCommonRequest("associatedVpcPublicIp", reqParams)

	resp, err := config.Client.vserver.V2Api.AssociatePublicIpWithServerInstance(reqParams)
	if err != nil {
		logErrorResponse("associatedVpcPublicIp", err, d.Id())
		return err
	}
	logCommonResponse("associatedVpcPublicIp", GetCommonResponse(resp))

	return nil
}

func resourceNcloudPublicIpCustomizeDiff(diff *schema.ResourceDiff, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if config.SupportVPC {
		if v, ok := diff.GetOk("zone"); ok {
			diff.Clear("zone")
			return fmt.Errorf("You don't use 'zone' if SupportVPC is true. Please remove this value [%s]", v)
		}

		if v, ok := diff.GetOk("internet_line_type"); ok {
			return fmt.Errorf("You don't use 'internet_line_type' if SupportVPC is true. Please remove this value [%s]", v)
		}
	}

	if diff.HasChange("description") {
		old, new := diff.GetChange("description")
		if len(old.(string)) > 0 {
			return fmt.Errorf("Change 'description' is not support, Please set `description` as a old value = [%s -> %s]", new, old)
		}
	}

	return nil
}
