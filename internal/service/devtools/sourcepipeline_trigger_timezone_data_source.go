package devtools

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudSourcePipelineTimeZone() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudSourcePipelineTimeZoneRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"timezone": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceNcloudSourcePipelineTimeZoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)

	timeZone, err := getSourcePipelineTimeZone(ctx, config)
	if err != nil {
		LogErrorResponse("getSourcePipelineTimeZone", err, timeZone)
		return diag.FromErr(err)
	}
	LogResponse("getSourcePipelineTimeZone", timeZone)

	if timeZone == nil {
		d.SetId("")
		return nil
	}

	d.SetId(time.Now().UTC().String())
	d.Set("timezone", timeZone)

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		return diag.FromErr(WriteToFile(output.(string), timeZone))
	}
	return nil
}

func getSourcePipelineTimeZone(ctx context.Context, config *conn.ProviderConfig) ([]*string, error) {
	resp, err := config.Client.Vsourcepipeline.V1Api.GetTimeZone(ctx)
	if err != nil {
		return nil, err
	}
	return resp.TimeZone, nil
}
