package cloudmssql

import (
	"context"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmssql"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

func DataSourceNcloudMssql() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"service_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"is_ha": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "default: true",
		},
		"is_multi_zone": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "default: false",
		},
		"is_backup": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "default: true",
		},
		"backup_file_retention_period": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "default: 1(1 day)",
		},
		"backup_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "ex) 01:15",
		},
		"port": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "default: 1433",
		},
		"image_product_code": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"instance_no": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"server_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"zone_code": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"vpc_no": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"subnet_no": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"server_instance_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"server_instance_no": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"server_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"server_instance_status_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"region_code": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"zone_code": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"vpc_no": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"subnet_no": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"data_storage_size": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"cpu_count": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"memory_size": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
		"filter": DataSourceFiltersSchema(),
	}
	return GetSingularDataSourceItemSchemaContext(ResourceNcloudMssql(), fieldMap, dataSourceNcloudMssqlRead)
}

func dataSourceNcloudMssqlRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("data source `ncloud_mssql`"))
	}

	msList, err := getCloudMssqlList(d, config)
	if err != nil {
		return diag.FromErr(err)
	}

	msListMap := ConvertToArrayMap(msList)
	if f, ok := d.GetOk("filter"); ok {
		msListMap = ApplyFilters(f.(*schema.Set), msListMap, DataSourceNcloudMssql().Schema)
	}

	if err := ValidateOneResult(len(msListMap)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(msListMap[0]["cloudMssqlInstanceNo"].(string))
	SetSingularResourceDataFromMapSchema(DataSourceNcloudMssql(), d, msListMap[0])
	return nil
}

func getCloudMssqlList(d *schema.ResourceData, config *conn.ProviderConfig) ([]*CloudMssqlInstance, error) {
	reqParams := &vmssql.GetCloudMssqlInstanceListRequest{
		CloudMssqlServiceName: StringPtrOrNil(d.GetOk("service_name")),
		CloudMssqlServerName:  StringPtrOrNil(d.GetOk("server_name")),
		ZoneCode:              StringPtrOrNil(d.GetOk("zone_code")),
		VpcNo:                 StringPtrOrNil(d.GetOk("vpc_no")),
		SubnetNo:              StringPtrOrNil(d.GetOk("subnet_no")),
		RegionCode:            &config.RegionCode,
	}

	if v, ok := d.GetOk("id"); ok {
		reqParams.CloudMssqlInstanceNoList = []*string{ncloud.String(v.(string))}
	}

	LogCommonRequest("getCloudMssqlList", reqParams)

	resp, err := config.Client.Vmssql.V2Api.GetCloudMssqlInstanceList(reqParams)
	if err != nil {
		LogErrorResponse("getCloudMssqlList", err, reqParams)
		return nil, err
	}

	LogResponse("getCloudMssqlList", resp)

	var msList []*CloudMssqlInstance
	for _, ms := range resp.CloudMssqlInstanceList {
		msList = append(msList, convertCloudMssql(ms))
	}

	return msList, nil
}
