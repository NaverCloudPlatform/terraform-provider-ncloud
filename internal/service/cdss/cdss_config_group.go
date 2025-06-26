package cdss

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vcdss"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func ResourceNcloudCDSSConfigGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNcloudCDSSConfigGroupCreate,
		ReadContext:   resourceNcloudCDSSConfigGroupRead,
		UpdateContext: resourceNcloudCDSSConfigGroupUpdate,
		DeleteContext: resourceNcloudCDSSConfigGroupDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), ":")
				if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected id:kafka_version_code", d.Id())
				}
				id := idParts[0]
				KafkaVersionCode := idParts[1]
				d.SetId(id)
				d.Set("kafka_version_code", KafkaVersionCode)
				return []*schema.ResourceData{d}, nil
			},
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
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(3, 15)),
			},
			"description": {
				Type:             schema.TypeString,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 255)),
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
	config := meta.(*conn.ProviderConfig)

	reqParams := vcdss.CreateConfigGroup{
		ConfigGroupName:  *StringPtrOrNil(d.GetOk("name")),
		KafkaVersionCode: *StringPtrOrNil(d.GetOk("kafka_version_code")),
	}

	description := StringPtrOrNil(d.GetOk("description"))
	if description != nil {
		reqParams.Description = *description
	}

	LogCommonRequest("resourceNcloudCDSSConfigGroupCreate", reqParams)
	resp, _, err := config.Client.Vcdss.V1Api.ConfigGroupCreateConfigGroupPost(ctx, reqParams)
	if err != nil {
		LogErrorResponse("resourceNcloudCDSSConfigGroupCreate", err, reqParams)
		return diag.FromErr(err)
	}
	LogResponse("resourceNcloudCDSSConfigGroupCreate", resp)

	id := strconv.Itoa(int(ncloud.Int32Value(&resp.Result.ConfigGroupNo)))
	d.SetId(id)
	return resourceNcloudCDSSConfigGroupRead(ctx, d, meta)
}

func resourceNcloudCDSSConfigGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)

	configGroup, err := getCDSSConfigGroup(ctx, config, d.Get("kafka_version_code").(string), d.Id())
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
	config := meta.(*conn.ProviderConfig)

	if d.HasChanges("description") {
		_, n := d.GetChange("description")

		newDescription := n.(string)
		LogCommonRequest("resourceNcloudCDSSConfigGroupUpdate", d.Id())

		reqParams := vcdss.SetKafkaConfigGroupMemoRequest{
			KafkaVersionCode: *StringPtrOrNil(d.GetOk("kafka_version_code")),
			Description:      newDescription,
		}

		if _, _, err := config.Client.Vcdss.V1Api.ConfigGroupSetKafkaConfigGroupMemoConfigGroupNoPost(ctx, reqParams, d.Id()); err != nil {
			LogErrorResponse("resourceNcloudCDSSConfigGroupUpdate", err, d.Id())
			return diag.FromErr(err)
		}
	}
	return nil
}

func resourceNcloudCDSSConfigGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)

	LogCommonRequest("resourceNcloudCDSSConfigGroupDelete", d.Id())
	if _, _, err := config.Client.Vcdss.V1Api.ConfigGroupDeleteConfigGroupConfigGroupNoDelete(ctx, d.Id()); err != nil {
		LogErrorResponse("resourceNcloudCDSSConfigGroupDelete", err, d.Id())
		return diag.FromErr(err)
	}

	return nil
}

func getCDSSConfigGroup(ctx context.Context, config *conn.ProviderConfig, kafkaVersionCode string, id string) (*vcdss.GetKafkaConfigGroupResponseVo, error) {
	reqParams := vcdss.GetKafkaConfigGroupRequest{
		KafkaVersionCode: kafkaVersionCode,
	}
	LogCommonRequest("getCDSSConfigGroup", reqParams)

	resp, _, err := config.Client.Vcdss.V1Api.ConfigGroupGetKafkaConfigGroupConfigGroupNoPost(ctx, reqParams, id)
	if err != nil {
		return nil, err
	}
	LogResponse("getCDSSConfigGroup", resp)

	return resp.Result, nil
}
