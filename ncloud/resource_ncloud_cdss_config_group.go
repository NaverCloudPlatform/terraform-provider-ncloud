package ncloud

import (
	"context"
	"strconv"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vcdss"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	RegisterResource("ncloud_cdss_config_group", resourceNcloudCDSSConfigGroup())
}

func resourceNcloudCDSSConfigGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNcloudCDSSConfigGroupCreate,
		ReadContext:   resourceNcloudCDSSConfigGroupRead,
		UpdateContext: resourceNcloudCDSSConfigGroupUpdate,
		DeleteContext: resourceNcloudCDSSConfigGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(3, 15)),
			},
			"description": {
				Type:             schema.TypeString,
				ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(0, 255)),
				Optional:         true,
			},
			"kafka_version_code": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

func resourceNcloudCDSSConfigGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_cdss_config_group`"))
	}

	reqParams := vcdss.CreateConfigGroup{
		ConfigGroupName:  *StringPtrOrNil(d.GetOk("name")),
		KafkaVersionCode: *StringPtrOrNil(d.GetOk("kafka_version_code")),
	}

	description := StringPtrOrNil(d.GetOk("description"))
	if description != nil {
		reqParams.Description = *description
	}

	logCommonRequest("resourceNcloudCDSSConfigGroupCreate", reqParams)
	resp, _, err := config.Client.vcdss.V1Api.ConfigGroupCreateConfigGroupPost(ctx, reqParams)
	if err != nil {
		logErrorResponse("resourceNcloudCDSSConfigGroupCreate", err, reqParams)
		return diag.FromErr(err)
	}
	logResponse("resourceNcloudCDSSConfigGroupCreate", resp)

	id := strconv.Itoa(int(ncloud.Int32Value(&resp.Result.ConfigGroupNo)))
	d.SetId(id)
	return resourceNcloudCDSSConfigGroupRead(ctx, d, meta)
}

func resourceNcloudCDSSConfigGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_cdss_config_group`"))
	}

	configGroup, err := getCDSSConfigGroup(ctx, config, *StringPtrOrNil(d.GetOk("kafka_version_code")), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if configGroup == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", configGroup.ConfigGroupName)
	d.Set("kafka_version_code", configGroup.KafkaVersionCode)
	d.Set("description", configGroup.Description)

	return nil
}

func resourceNcloudCDSSConfigGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_cdss_cluster`"))
	}

	if d.HasChanges("description") {
		_, n := d.GetChange("description")

		newDescription := n.(string)
		logCommonRequest("resourceNcloudCDSSConfigGroupUpdate", d.Id())

		reqParams := vcdss.SetKafkaConfigGroupMemoRequest{
			KafkaVersionCode: *StringPtrOrNil(d.GetOk("kafka_version_code")),
			Description:      newDescription,
		}

		if _, _, err := config.Client.vcdss.V1Api.ConfigGroupSetKafkaConfigGroupMemoConfigGroupNoPost(ctx, reqParams, d.Id()); err != nil {
			logErrorResponse("resourceNcloudCDSSConfigGroupUpdate", err, d.Id())
			return diag.FromErr(err)
		}
	}
	return nil
}

func resourceNcloudCDSSConfigGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_cdss_config_group`"))
	}

	logCommonRequest("resourceNcloudCDSSConfigGroupDelete", d.Id())
	if _, _, err := config.Client.vcdss.V1Api.ConfigGroupDeleteConfigGroupConfigGroupNoDelete(ctx, d.Id()); err != nil {
		logErrorResponse("resourceNcloudCDSSConfigGroupDelete", err, d.Id())
		return diag.FromErr(err)
	}

	return nil
}

func getCDSSConfigGroup(ctx context.Context, config *ProviderConfig, kafkaVersionCode string, id string) (*vcdss.GetKafkaConfigGroupResponseVo, error) {
	reqParams := vcdss.GetKafkaConfigGroupRequest{
		KafkaVersionCode: kafkaVersionCode,
	}
	resp, _, err := config.Client.vcdss.V1Api.ConfigGroupGetKafkaConfigGroupConfigGroupNoPost(ctx, reqParams, id)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}
