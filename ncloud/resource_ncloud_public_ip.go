package ncloud

import (
	"context"
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
		Update: resourceNcloudPublicIpUpdate,
		Delete: resourceNcloudPublicIpDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: resourceNcloudPublicIpCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"server_instance_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(1, 10000)),
			},

			"public_ip_no": {
				Type:     schema.TypeString,
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

	instance := ConvertToMap(resource)
	SetSingularResourceDataFromMapSchema(resourceNcloudPublicIpInstance(), d, instance)
	if err := d.Set("public_ip_no", resource.PublicIpInstanceNo); err != nil {
		return err
	}

	if err := d.Set("server_instance_no", resource.ServerInstanceNo); err != nil {
		return err
	}

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
			if err := resource.Retry(time.Minute, func() *resource.RetryError {
				if err := associatedPublicIp(d, config); err != nil {
					errBody, _ := GetCommonErrorBody(err)
					if errBody.ReturnCode == "1003016" {
						time.Sleep(time.Second * 1)
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			}); err != nil {
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
		RegionNo:            &config.RegionNo,
		ZoneNo:              zoneNo,
		ServerInstanceNo:    StringPtrOrNil(d.GetOk("server_instance_no")),
		PublicIpDescription: StringPtrOrNil(d.GetOk("description")),
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

func getPublicIp(config *ProviderConfig, id string) (*PublicIpInstance, error) {
	var r *PublicIpInstance
	var err error
	if config.SupportVPC {
		r, err = getVpcPublicIp(config, id)
	} else {
		r, err = getClassicPublicIp(config, id)
	}

	if err != nil {
		return nil, err
	}

	return r, nil
}

func getClassicPublicIp(config *ProviderConfig, id string) (*PublicIpInstance, error) {
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

	p := &PublicIpInstance{
		PublicIpInstanceNo:            r.PublicIpInstanceNo,
		PublicIp:                      r.PublicIp,
		PublicIpDescription:           r.PublicIpDescription,
		PublicIpKindTypeCode:          r.PublicIpKindType.Code,
		ZoneCode:                      r.Zone.ZoneCode,
		PublicIpInstanceStatusCode:    r.PublicIpInstanceStatus.Code,
		PublicIpInstanceOperationCode: r.PublicIpInstanceOperation.Code,
	}

	if r.ServerInstanceAssociatedWithPublicIp != nil {
		p.ServerInstanceNo = r.ServerInstanceAssociatedWithPublicIp.ServerInstanceNo
		p.PrivateIp = r.ServerInstanceAssociatedWithPublicIp.PrivateIp
	}

	return p, nil
}

func getVpcPublicIp(config *ProviderConfig, id string) (*PublicIpInstance, error) {
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

	p := &PublicIpInstance{
		PublicIpInstanceNo:            r.PublicIpInstanceNo,
		PublicIp:                      r.PublicIp,
		PublicIpDescription:           r.PublicIpDescription,
		PublicIpInstanceStatusCode:    r.PublicIpInstanceStatus.Code,
		ServerInstanceNo:              r.ServerInstanceNo,
		PrivateIp:                     r.PrivateIp,
		LastModifyDate:                r.LastModifyDate,
		PublicIpInstanceOperationCode: r.PublicIpInstanceOperation.Code,
		ZoneCode:                      nil,
	}

	return p, nil
}

func checkAssociatedPublicIP(config *ProviderConfig, id string) (bool, error) {
	instance, err := getPublicIp(config, id)

	if err != nil {
		return false, err
	}

	return instance.ServerInstanceNo != nil && *instance.ServerInstanceNo != "", nil
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
			opCode, opErr := getPublicIpInstanceOperationCode(config, id)

			if err != nil || opErr != nil {
				return 0, "", err
			}

			if !isAssociated && opCode == "NULL" {
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

func getPublicIpInstanceOperationCode(config *ProviderConfig, id string) (string, error) {
	instance, err := getPublicIp(config, id)
	if err != nil {
		return "", err
	}
	return *instance.PublicIpInstanceOperationCode, nil
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

func resourceNcloudPublicIpCustomizeDiff(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if config.SupportVPC {
		if v, ok := diff.GetOk("zone"); ok {
			diff.Clear("zone")
			return fmt.Errorf("You don't use 'zone' if SupportVPC is true. Please remove this value [%s]", v)
		}
	}
	return nil
}

type PublicIpInstance struct {
	PublicIpInstanceNo   *string `json:"instance_no,omitempty"`
	PublicIp             *string `json:"public_ip,omitempty"`
	PublicIpDescription  *string `json:"description,omitempty"`
	ServerInstanceNo     *string `json:"server_instance_no,omitempty"`
	PublicIpKindTypeCode *string `json:"kind_type,omitempty"`
	ZoneCode             *string `json:"zone,omitempty"`

	PublicIpInstanceStatusCode    *string
	PublicIpInstanceOperationCode *string
	// Server Instance
	PrivateIp      *string
	LastModifyDate *string
}
