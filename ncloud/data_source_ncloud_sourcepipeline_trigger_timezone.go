package ncloud

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_sourcepipeline_trigger_timezone", dataSourceNcloudSourcePipelineTimeZone())
}

func dataSourceNcloudSourcePipelineTimeZone() *schema.Resource {
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
	config := meta.(*ProviderConfig)

	timeZone, err := getSourcePipelineTimeZone(ctx, config)
	if err != nil {
		logErrorResponse("getSourcePipelineTimeZone", err, timeZone)
		return diag.FromErr(err)
	}
	logResponse("getSourcePipelineTimeZone", timeZone)

	if timeZone == nil {
		d.SetId("")
		return nil
	}

	d.SetId(time.Now().UTC().String())
	d.Set("timezone", timeZone)

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		return diag.FromErr(writeToFile(output.(string), timeZone))
	}
	return nil
}

func getSourcePipelineTimeZone(ctx context.Context, config *ProviderConfig) ([]*string, error) {
	if config.SupportVPC {
		return getVpcSourcePipelineTimeZone(ctx, config)
	}
	return getClassicSourcePipelineTimeZone(ctx, config)
}

func getClassicSourcePipelineTimeZone(ctx context.Context, config *ProviderConfig) ([]*string, error) {
	resp, err := config.Client.sourcepipeline.V1Api.GetTimeZone(ctx)
	if err != nil {
		return nil, err
	}
	return resp.TimeZone, nil
}

func getVpcSourcePipelineTimeZone(ctx context.Context, config *ProviderConfig) ([]*string, error) {
	resp, err := config.Client.vsourcepipeline.V1Api.GetTimeZone(ctx)
	if err != nil {
		return nil, err
	}
	return resp.TimeZone, nil
}
