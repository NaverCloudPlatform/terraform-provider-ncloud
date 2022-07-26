package ncloud

import (
	"context"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vcdss"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strconv"
)

func init() {
	RegisterResource("ncloud_vcdss_config_group", resourceNcloudVCDSSConfigGroup())
}

func resourceNcloudVCDSSConfigGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNcloudVCDSSConfigGroupCreate,
		ReadContext:   resourceNcloudVCDSSConfigGroupRead,
		DeleteContext: resourceNcloudVCDSSConfigGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config_group_name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(3, 15)),
			},
			"description": {
				Type:             schema.TypeString,
				ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(0, 255)),
				ForceNew:         true,
				Required:         true,
			},
			"kafka_version_code": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

func resourceNcloudVCDSSConfigGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_cdss_config_group`"))
	}

	reqParams := vcdss.CreateConfigGroup{
		ConfigGroupName:  *StringPtrOrNil(d.GetOk("config_group_name")),
		Description:      *StringPtrOrNil(d.GetOk("description")),
		KafkaVersionCode: *StringPtrOrNil(d.GetOk("kafka_version_code")),
	}

	logCommonRequest("resourceNcloudVCDSSClusterCreate", reqParams)
	resp, _, err := config.Client.vcdss.V1Api.ConfigGroupCreateConfigGroupPost(ctx, reqParams)
	if err != nil {
		logErrorResponse("resourceNcloudVCDSSConfigGroupCreate", err, reqParams)
		return diag.FromErr(err)
	}
	logResponse("resourceNcloudVCDSSConfigGroupCreate", resp)

	uuid := strconv.Itoa(int(ncloud.Int32Value(&resp.Result.ConfigGroupNo)))
	d.SetId(uuid)
	return resourceNcloudVCDSSConfigGroupRead(ctx, d, meta)
}

func resourceNcloudVCDSSConfigGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_cdss_config_group`"))
	}

	configGroup, err := getVCDSSConfigGroup(ctx, config, *StringPtrOrNil(d.GetOk("kafka_version_code")), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if configGroup == nil {
		d.SetId("")
		return nil
	}

	d.Set("config_group_name", configGroup.ConfigGroupName)
	d.Set("kafkaVersionCode", configGroup.KafkaVersionCode)
	d.Set("description", configGroup.Description)

	return nil
}

func resourceNcloudVCDSSConfigGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_cdss_config_group`"))
	}

	logCommonRequest("resourceNcloudVCDSConfigGroupDelete", d.Id())
	if _, _, err := config.Client.vcdss.V1Api.ConfigGroupDeleteConfigGroupConfigGroupNoDelete(ctx, d.Id()); err != nil {
		logErrorResponse("resourceNcloudVCDSSConfigGroupDelete", err, d.Id())
		return diag.FromErr(err)
	}

	return nil
}

func getVCDSSConfigGroup(ctx context.Context, config *ProviderConfig, kafkaVersionCode string, uuid string) (*vcdss.GetKafkaConfigGroupResponseVo, error) {
	reqParams := vcdss.GetKafkaConfigGroupRequest{
		KafkaVersionCode: kafkaVersionCode,
	}
	resp, _, err := config.Client.vcdss.V1Api.ConfigGroupGetKafkaConfigGroupConfigGroupNoPost(ctx, reqParams, uuid)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}
