package server

import (
	"log"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

func ResourceNcloudInitScript() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudInitScriptCreate,
		Read:   resourceNcloudInitScriptRead,
		Delete: resourceNcloudInitScriptDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: ToDiagFunc(ValidateInstanceName),
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(0, 1000)),
			},
			"os_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"LNX", "WND"}, false)),
			},

			"init_script_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNcloudInitScriptCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	instance, err := createInitScript(d, config)

	if err != nil {
		return err
	}

	d.SetId(*instance.InitScriptNo)
	log.Printf("[INFO] Init script ID: %s", d.Id())

	return resourceNcloudInitScriptRead(d, meta)
}

func resourceNcloudInitScriptRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	instance, err := GetInitScript(config, d.Id())
	if err != nil {
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*instance.InitScriptNo)
	d.Set("init_script_no", instance.InitScriptNo)
	d.Set("name", instance.InitScriptName)
	d.Set("description", instance.InitScriptDescription)
	d.Set("content", instance.InitScriptContent)
	d.Set("os_type", instance.OsType.Code)

	return nil
}

func resourceNcloudInitScriptDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	if err := DeleteInitScript(config, d.Id()); err != nil {
		return err
	}

	return nil
}

func GetInitScript(config *conn.ProviderConfig, id string) (*vserver.InitScript, error) {
	if config.SupportVPC {
		return getVpcInitScript(config, id)
	}

	return nil, NotSupportClassic("resource `ncloud_init_script`")
}

func getVpcInitScript(config *conn.ProviderConfig, id string) (*vserver.InitScript, error) {
	reqParams := &vserver.GetInitScriptDetailRequest{
		RegionCode:   &config.RegionCode,
		InitScriptNo: ncloud.String(id),
	}

	LogCommonRequest("GetInitScriptDetail", reqParams)
	resp, err := config.Client.Vserver.V2Api.GetInitScriptDetail(reqParams)
	if err != nil {
		LogErrorResponse("GetInitScriptDetail", err, reqParams)
		return nil, err
	}
	LogResponse("GetInitScriptDetail", resp)

	if len(resp.InitScriptList) > 0 {
		return resp.InitScriptList[0], nil
	}

	return nil, nil
}

func createInitScript(d *schema.ResourceData, config *conn.ProviderConfig) (*vserver.InitScript, error) {
	if config.SupportVPC {
		return createVpcInitScript(d, config)
	}

	return nil, NotSupportClassic("resource `ncloud_init_script`")
}

func createVpcInitScript(d *schema.ResourceData, config *conn.ProviderConfig) (*vserver.InitScript, error) {
	reqParams := &vserver.CreateInitScriptRequest{
		RegionCode:            &config.RegionCode,
		InitScriptContent:     ncloud.String(d.Get("content").(string)),
		InitScriptName:        StringPtrOrNil(d.GetOk("name")),
		InitScriptDescription: StringPtrOrNil(d.GetOk("description")),
		OsTypeCode:            StringPtrOrNil(d.GetOk("os_type")),
	}

	LogCommonRequest("createVpcInitScript", reqParams)
	resp, err := config.Client.Vserver.V2Api.CreateInitScript(reqParams)
	if err != nil {
		LogErrorResponse("createVpcInitScript", err, reqParams)
		return nil, err
	}
	LogResponse("createVpcInitScript", resp)

	return resp.InitScriptList[0], nil
}

func DeleteInitScript(config *conn.ProviderConfig, id string) error {
	if config.SupportVPC {
		return deleteVpcInitScript(config, id)
	}

	return NotSupportClassic("resource `ncloud_init_script`")
}

func deleteVpcInitScript(config *conn.ProviderConfig, id string) error {
	reqParams := &vserver.DeleteInitScriptsRequest{
		RegionCode:       &config.RegionCode,
		InitScriptNoList: []*string{ncloud.String(id)},
	}

	LogCommonRequest("deleteVpcInitScript", reqParams)
	resp, err := config.Client.Vserver.V2Api.DeleteInitScripts(reqParams)
	if err != nil {
		LogErrorResponse("deleteVpcInitScript", err, reqParams)
		return err
	}
	LogResponse("deleteVpcInitScript", resp)

	return nil
}
