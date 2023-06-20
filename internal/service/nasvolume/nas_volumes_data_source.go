package nasvolume

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/provider"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

func init() {
	RegisterDataSource("ncloud_nas_volumes", dataSourceNcloudNasVolumes())
}

func dataSourceNcloudNasVolumes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudNasVolumesRead,

		Schema: map[string]*schema.Schema{
			"volume_allotment_protocol_type_code": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"NFS", "CIFS"}, false)),
			},
			"is_event_configuration": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_snapshot_configuration": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"no_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region code. Get available values using the `data ncloud_regions`.",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Zone code. Get available values using the `data ncloud_zones`.",
			},
			"filter": DataSourceFiltersSchema(),
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nas_volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     GetDataSourceItemSchema(resourceNcloudNasVolume()),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceNcloudNasVolumesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	instances, err := getNasVolumeList(d, config)
	if err != nil {
		return err
	}

	resources := ConvertToArrayMap(instances)

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudNasVolumes().Schema)
	}

	if len(resources) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	var ids []string
	for _, r := range resources {
		ids = append(ids, r["nas_volume_no"].(string))
	}

	d.SetId(DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return err
	}

	if err := d.Set("nas_volumes", resources); err != nil {
		return err
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		return WriteToFile(output.(string), d.Get("nas_volumes"))
	}

	return nil
}
