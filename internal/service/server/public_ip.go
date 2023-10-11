package server

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

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/zone"
)

func ResourceNcloudPublicIpInstance() *schema.Resource {
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
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 10000)),
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
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"PUBLC", "GLBL"}, false)),
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
	config := meta.(*conn.ProviderConfig)
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
	config := meta.(*conn.ProviderConfig)

	resource, err := GetPublicIp(config, d.Id())

	if err != nil {
		return err
	}

	if resource == nil {
		d.SetId("")
		return nil
	}

	instance := ConvertToMap(resource)
	SetSingularResourceDataFromMapSchema(ResourceNcloudPublicIpInstance(), d, instance)
	if err := d.Set("public_ip_no", resource.PublicIpInstanceNo); err != nil {
		return err
	}

	if err := d.Set("server_instance_no", resource.ServerInstanceNo); err != nil {
		return err
	}

	return nil
}

func resourceNcloudPublicIpDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)
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
	config := meta.(*conn.ProviderConfig)

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

func createClassicPublicIp(d *schema.ResourceData, config *conn.ProviderConfig) (*string, error) {
	client := config.Client

	zoneNo, err := zone.ParseZoneNoParameter(config, d)
	if err != nil {
		return nil, err
	}

	reqParams := &server.CreatePublicIpInstanceRequest{
		RegionNo:            &config.RegionNo,
		ZoneNo:              zoneNo,
		ServerInstanceNo:    StringPtrOrNil(d.GetOk("server_instance_no")),
		PublicIpDescription: StringPtrOrNil(d.GetOk("description")),
	}

	LogCommonRequest("createClassicPublicIp", reqParams)

	resp, err := client.Server.V2Api.CreatePublicIpInstance(reqParams)
	if err != nil {
		LogErrorResponse("createClassicPublicIp", err, reqParams)
		return nil, err
	}
	LogResponse("createClassicPublicIp", resp)

	publicIPInstance := resp.PublicIpInstanceList[0]

	return publicIPInstance.PublicIpInstanceNo, nil
}

func createVpcPublicIp(d *schema.ResourceData, config *conn.ProviderConfig) (*string, error) {
	client := config.Client

	reqParams := &vserver.CreatePublicIpInstanceRequest{
		RegionCode:          &config.RegionCode,
		ServerInstanceNo:    StringPtrOrNil(d.GetOk("server_instance_no")),
		PublicIpDescription: StringPtrOrNil(d.GetOk("description")),
	}

	LogCommonRequest("createVpcPublicIp", reqParams)

	resp, err := client.Vserver.V2Api.CreatePublicIpInstance(reqParams)
	if err != nil {
		LogErrorResponse("createVpcPublicIp", err, reqParams)
		return nil, err
	}
	LogResponse("createVpcPublicIp", resp)

	publicIPInstance := resp.PublicIpInstanceList[0]

	return publicIPInstance.PublicIpInstanceNo, nil
}

func deleteClassicPublicIp(d *schema.ResourceData, config *conn.ProviderConfig) error {
	client := config.Client

	reqParams := &server.DeletePublicIpInstancesRequest{
		PublicIpInstanceNoList: []*string{ncloud.String(d.Id())},
	}

	LogCommonRequest("deleteClassicPublicIp", reqParams)

	resp, err := client.Server.V2Api.DeletePublicIpInstances(reqParams)
	if err != nil {
		LogErrorResponse("deleteClassicPublicIp", err, reqParams)
		return err
	}
	LogResponse("deleteClassicPublicIp", resp)

	return nil
}

func deleteVpcPublicIp(d *schema.ResourceData, config *conn.ProviderConfig) error {
	client := config.Client

	reqParams := &vserver.DeletePublicIpInstanceRequest{
		RegionCode:         &config.RegionCode,
		PublicIpInstanceNo: ncloud.String(d.Id()),
	}

	LogCommonRequest("deleteVpcPublicIp", reqParams)

	resp, err := client.Vserver.V2Api.DeletePublicIpInstance(reqParams)
	if err != nil {
		LogErrorResponse("deleteVpcPublicIp", err, reqParams)
		return err
	}
	LogResponse("deleteVpcPublicIp", resp)

	return nil
}

func GetPublicIp(config *conn.ProviderConfig, id string) (*PublicIpInstance, error) {
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

func getClassicPublicIp(config *conn.ProviderConfig, id string) (*PublicIpInstance, error) {
	client := config.Client
	regionNo := config.RegionNo

	reqParams := &server.GetPublicIpInstanceListRequest{
		RegionNo:               &regionNo,
		PublicIpInstanceNoList: []*string{ncloud.String(id)},
	}

	LogCommonRequest("getClassicPublicIp", reqParams)
	resp, err := client.Server.V2Api.GetPublicIpInstanceList(reqParams)

	if err != nil {
		LogErrorResponse("getClassicPublicIp", err, reqParams)
		return nil, err
	}
	LogResponse("getClassicPublicIp", resp)

	if len(resp.PublicIpInstanceList) == 0 {
		return nil, nil
	}

	if err := verify.ValidateOneResult(len(resp.PublicIpInstanceList)); err != nil {
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

func getVpcPublicIp(config *conn.ProviderConfig, id string) (*PublicIpInstance, error) {
	client := config.Client
	regionCode := config.RegionCode

	reqParams := &vserver.GetPublicIpInstanceListRequest{
		RegionCode:             &regionCode,
		PublicIpInstanceNoList: []*string{ncloud.String(id)},
	}

	LogCommonRequest("getVpcPublicIp", reqParams)
	resp, err := client.Vserver.V2Api.GetPublicIpInstanceList(reqParams)

	if err != nil {
		LogErrorResponse("getVpcPublicIp", err, reqParams)
		return nil, err
	}
	LogResponse("getVpcPublicIp", resp)

	if len(resp.PublicIpInstanceList) == 0 {
		return nil, nil
	}

	if err := verify.ValidateOneResult(len(resp.PublicIpInstanceList)); err != nil {
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

func checkAssociatedPublicIP(config *conn.ProviderConfig, id string) (bool, error) {
	instance, err := GetPublicIp(config, id)

	if err != nil {
		return false, err
	}

	if instance == nil {
		return false, nil
	}

	return instance.ServerInstanceNo != nil && *instance.ServerInstanceNo != "", nil
}

func disassociatedPublicIp(config *conn.ProviderConfig, id string) error {
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

func disassociatedClassicPublicIp(config *conn.ProviderConfig, id string) error {
	reqParams := &server.DisassociatePublicIpFromServerInstanceRequest{PublicIpInstanceNo: ncloud.String(id)}

	LogCommonRequest("disassociatedClassicPublicIP", reqParams)

	resp, err := config.Client.Server.V2Api.DisassociatePublicIpFromServerInstance(reqParams)
	if err != nil {
		LogErrorResponse("disassociatedClassicPublicIP", err, id)
		return err
	}
	LogCommonResponse("disassociatedClassicPublicIP", GetCommonResponse(resp))

	return nil
}

func disassociatedVpcPublicIp(config *conn.ProviderConfig, id string) error {
	reqParams := &vserver.DisassociatePublicIpFromServerInstanceRequest{
		RegionCode:         &config.RegionCode,
		PublicIpInstanceNo: ncloud.String(id),
	}

	LogCommonRequest("disassociatedVpcPublicIp", reqParams)

	resp, err := config.Client.Vserver.V2Api.DisassociatePublicIpFromServerInstance(reqParams)
	if err != nil {
		LogErrorResponse("disassociatedVpcPublicIp", err, id)
		return err
	}
	LogCommonResponse("disassociatedVpcPublicIp", GetCommonResponse(resp))

	return nil
}

func waitForPublicIpDisassociation(config *conn.ProviderConfig, id string) error {
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
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Public IP (%s) to become disassociation: %s", id, err)
	}

	return nil
}

func getPublicIpInstanceOperationCode(config *conn.ProviderConfig, id string) (string, error) {
	instance, err := GetPublicIp(config, id)
	if err != nil {
		return "", err
	}
	return *instance.PublicIpInstanceOperationCode, nil
}

func waitForPublicIpAssociation(config *conn.ProviderConfig, id string) error {
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
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Public IP (%s) to become association: %s", id, err)
	}

	return nil
}

func associatedPublicIp(d *schema.ResourceData, config *conn.ProviderConfig) error {
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

func associatedClassicPublicIp(d *schema.ResourceData, config *conn.ProviderConfig) error {
	reqParams := &server.AssociatePublicIpWithServerInstanceRequest{
		PublicIpInstanceNo: ncloud.String(d.Id()),
		ServerInstanceNo:   ncloud.String(d.Get("server_instance_no").(string)),
	}

	LogCommonRequest("associatedClassicPublicIp", reqParams)

	resp, err := config.Client.Server.V2Api.AssociatePublicIpWithServerInstance(reqParams)
	if err != nil {
		LogErrorResponse("associatedClassicPublicIp", err, d.Id())
		return err
	}
	LogCommonResponse("associatedClassicPublicIp", GetCommonResponse(resp))

	return nil
}

func associatedVpcPublicIp(d *schema.ResourceData, config *conn.ProviderConfig) error {
	reqParams := &vserver.AssociatePublicIpWithServerInstanceRequest{
		RegionCode:         &config.RegionCode,
		PublicIpInstanceNo: ncloud.String(d.Id()),
		ServerInstanceNo:   ncloud.String(d.Get("server_instance_no").(string)),
	}

	LogCommonRequest("associatedVpcPublicIp", reqParams)

	resp, err := config.Client.Vserver.V2Api.AssociatePublicIpWithServerInstance(reqParams)
	if err != nil {
		LogErrorResponse("associatedVpcPublicIp", err, d.Id())
		return err
	}
	LogCommonResponse("associatedVpcPublicIp", GetCommonResponse(resp))

	return nil
}

func resourceNcloudPublicIpCustomizeDiff(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	if config.SupportVPC {
		if v, ok := diff.GetOk("zone"); ok {
			_ = diff.Clear("zone")
			return fmt.Errorf("you don't use 'zone' if SupportVPC is true. Please remove this value [%s]", v)
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
