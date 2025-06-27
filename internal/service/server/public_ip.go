package server

import (
	"fmt"
	"log"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
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

	publicIpInstanceNo, err = createVpcPublicIp(d, config)
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

	err = deleteVpcPublicIp(d, config)
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
	err := disassociatedVpcPublicIp(config, id)
	if err != nil {
		return err
	}

	if err := waitForPublicIpDisassociation(config, id); err != nil {
		return err
	}

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
	err := associatedVpcPublicIp(d, config)
	if err != nil {
		return err
	}

	if err := waitForPublicIpAssociation(config, d.Id()); err != nil {
		return err
	}

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

type PublicIpInstance struct {
	PublicIpInstanceNo   *string `json:"instance_no,omitempty"`
	PublicIp             *string `json:"public_ip,omitempty"`
	PublicIpDescription  *string `json:"description,omitempty"`
	ServerInstanceNo     *string `json:"server_instance_no,omitempty"`
	PublicIpKindTypeCode *string `json:"kind_type,omitempty"`

	PublicIpInstanceStatusCode    *string
	PublicIpInstanceOperationCode *string
	// Server Instance
	PrivateIp      *string
	LastModifyDate *string
}
