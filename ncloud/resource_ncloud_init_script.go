package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
)

func init() {
	RegisterResource("ncloud_init_script", resourceNcloudInitScript())
}

func resourceNcloudInitScript() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudInitScriptCreate,
		Read:   resourceNcloudInitScriptRead,
		Update: resourceNcloudInitScriptUpdate,
		Delete: resourceNcloudInitScriptDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: ncloudVpcCommonCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateInstanceName,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 1000),
			},
			"os_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"LNX", "WND"}, false),
			},

			"init_script_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNcloudInitScriptCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	instance, err := createInitScript(d, config)

	if err != nil {
		return err
	}

	d.SetId(*instance.InitScriptNo)
	log.Printf("[INFO] Init script ID: %s", d.Id())

	return resourceNcloudInitScriptRead(d, meta)
}

func resourceNcloudInitScriptRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	instance, err := getInitScript(config, d.Id())
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
	d.Set("os_type", instance.OsType)

	return nil
}

func resourceNcloudInitScriptUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceNcloudInitScriptRead(d, meta)
}

func resourceNcloudInitScriptDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if err := deleteInitScript(config, d.Id()); err != nil {
		return err
	}

	return nil
}

func getInitScript(config *ProviderConfig, id string) (*vserver.InitScript, error) {
	if config.SupportVPC {
		return getVpcInitScript(config, id)
	}

	return nil, NotSupportClassic("resource `ncloud_init_script`")
}

func getVpcInitScript(config *ProviderConfig, id string) (*vserver.InitScript, error) {
	reqParams := &vserver.GetInitScriptDetailRequest{
		RegionCode:   &config.RegionCode,
		InitScriptNo: ncloud.String(id),
	}

	logCommonRequest("GetInitScriptDetail", reqParams)
	resp, err := config.Client.vserver.V2Api.GetInitScriptDetail(reqParams)
	if err != nil {
		logErrorResponse("GetInitScriptDetail", err, reqParams)
		return nil, err
	}
	logResponse("GetInitScriptDetail", resp)

	if len(resp.InitScriptList) > 0 {
		return resp.InitScriptList[0], nil
	}

	return nil, nil
}

func createInitScript(d *schema.ResourceData, config *ProviderConfig) (*vserver.InitScript, error) {
	if config.SupportVPC {
		return createVpcInitScript(d, config)
	}

	return nil, NotSupportClassic("resource `ncloud_init_script`")
}

func createVpcInitScript(d *schema.ResourceData, config *ProviderConfig) (*vserver.InitScript, error) {
	reqParams := &vserver.CreateInitScriptRequest{
		RegionCode:            &config.RegionCode,
		InitScriptContent:     ncloud.String(d.Get("content").(string)),
		InitScriptName:        StringPtrOrNil(d.GetOk("name")),
		InitScriptDescription: StringPtrOrNil(d.GetOk("description")),
		OsTypeCode:            StringPtrOrNil(d.GetOk("os_type")),
	}

	logCommonRequest("createVpcInitScript", reqParams)
	resp, err := config.Client.vserver.V2Api.CreateInitScript(reqParams)
	if err != nil {
		logErrorResponse("createVpcInitScript", err, reqParams)
		return nil, err
	}
	logResponse("createVpcInitScript", resp)

	return resp.InitScriptList[0], nil
}

func deleteInitScript(config *ProviderConfig, id string) error {
	if config.SupportVPC {
		return deleteVpcInitScript(config, id)
	}

	return NotSupportClassic("resource `ncloud_init_script`")
}

func deleteVpcInitScript(config *ProviderConfig, id string) error {
	reqParams := &vserver.DeleteInitScriptsRequest{
		RegionCode:       &config.RegionCode,
		InitScriptNoList: []*string{ncloud.String(id)},
	}

	logCommonRequest("deleteVpcInitScript", reqParams)
	resp, err := config.Client.vserver.V2Api.DeleteInitScripts(reqParams)
	if err != nil {
		logErrorResponse("deleteVpcInitScript", err, reqParams)
		return err
	}
	logResponse("deleteVpcInitScript", resp)

	return nil
}
