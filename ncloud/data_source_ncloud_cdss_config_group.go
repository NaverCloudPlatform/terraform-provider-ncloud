package ncloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_cdss_config_group", dataSourceNcloudCDSSConfigGroup())
}

func dataSourceNcloudCDSSConfigGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNcloudCDSSConfigGroupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"kafka_version_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNcloudCDSSConfigGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("dataSource `ncloud_cdss_config_group`"))
	}

	configGroup, err := getCDSSConfigGroup(ctx, config, *StringPtrOrNil(d.GetOk("kafka_version_code")), *StringPtrOrNil(d.GetOk("id")))
	if err != nil {
		return diag.FromErr(err)
	}

	if configGroup == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*StringPtrOrNil(d.GetOk("id")))
	d.Set("name", configGroup.ConfigGroupName)
	d.Set("description", configGroup.Description)

	return nil
}
