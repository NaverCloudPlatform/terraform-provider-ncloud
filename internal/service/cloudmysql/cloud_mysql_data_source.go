package cloudmysql

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmysql"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudMysql() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"service_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"engine_version": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"is_ha": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "default: true",
		},
		"is_multi_zone": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "default: false",
		},
		"is_backup": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "default: true",
		},
		"backup_file_retention_period": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "default: 1(1 day)",
		},
		"backup_time": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "ex) 01:15",
		},
		"port": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "default: 3306",
		},
		"image_product_code": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"instance_no": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"cloud_mysql_server_instance_list": {
			Type: schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"cloud_mysql_server_instance_no": {
						Type: schema.TypeString,
						Computed: true,
						ForceNew: true,

					},
					"cloud_mysql_server_name": {
						Type: schema.TypeString,
						Computed: true,
						ForceNew: true,

					},
					"region_code": {
						Type: schema.TypeString,
						Computed: true,
						ForceNew: true,

					},
					"zone_code": {
						Type: schema.TypeString,
						Computed: true,
						ForceNew: true,

					},
					"vpc_no": {
						Type: schema.TypeString,
						Computed: true,
						ForceNew: true,
					},
					"subnet_no": {
						Type: schema.TypeString,
						Computed: true,
						ForceNew: true,
					},
					"data_storage_size": {
						Type: schema.TypeInt,
						Computed: true,
						ForceNew: true,
					},
					"used_data_storage_size": {
						Type: schema.TypeInt,
						Computed: true,
						ForceNew: true,
					},
					"cpu_count": {
						Type: schema.TypeInt,
						Computed: true,
						ForceNew: true,
					},
				},
			},
		},
		"filter": DataSourceFiltersSchema(),
	}

	return GetSingularDataSourceItemSchema(ResourceNcloudMySql(), fieldMap, dataSourceNcloudMysqlRead)
}

func dataSourceNcloudMysqlRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	if !config.SupportVPC {
		return NotSupportClassic("data source `ncloud_mysql`")
	}

	instance, err := getMysqlList(d, config)
	if err != nil {
		return err
	}
	d.SetId(instance[0]["id"].(string))

	detailInstance, err := getMysqlDetail(config, instance)
	if err != nil {
		return err
	}


	resources := ConvertToArrayMap(detailInstance)
	SetSingularResourceDataFromMap(d, resources[0])
	return nil
}

func getMysqlDetail(config *conn.ProviderConfig, r []map[string]interface{}) ([]map[string]interface{}, error){
	reqParams := &vmysql.GetCloudMysqlInstanceDetailRequest{
		RegionCode:           &config.RegionCode,
		CloudMysqlInstanceNo: ncloud.String(r[0]["id"].(string)),
	}

	resp, err := config.Client.Vmysql.V2Api.GetCloudMysqlInstanceDetail(reqParams)
	if err != nil {
		LogErrorResponse("getCloudMysqlDetail", err, reqParams)
		return nil, err
	}
	LogResponse("GetMysqlDetail", resp)

	var resources []map[string]interface{}
	for _, r := range resp.CloudMysqlInstanceList {
		instance := map[string]interface{}{
			"id":                               *r.CloudMysqlInstanceNo,
			"instance_no":                         *r.CloudMysqlInstanceNo,
			"service_name":                     *r.CloudMysqlServiceName,
			"engine_version":                   *r.EngineVersion,
			"is_ha":                            *r.IsHa,
			"is_multi_zone":                    *r.IsMultiZone,
			"is_backup":                        *r.IsBackup,
			"backup_file_retention_period":     *r.BackupFileRetentionPeriod,
			"backup_time":                      *r.BackupTime,
			"port":                             *r.CloudMysqlPort,
			"image_product_code":               *r.CloudMysqlImageProductCode,
			"cloud_mysql_server_instance_list": map[string]interface{}{},
		}
		var serverInstanceList []map[string]interface{}
		for _, v := range r.CloudMysqlServerInstanceList {
			ins := map[string]interface{}{
				"cloud_mysql_server_instance_no": v.CloudMysqlServerInstanceNo,
				"cloud_mysql_server_name":        v.CloudMysqlServerName,
				"region_code":                    v.RegionCode,
				"zone_code":                      v.ZoneCode,
				"vpc_no":                         v.VpcNo,
				"subnet_no":                      v.SubnetNo,
				"data_storage_size":              v.DataStorageSize,
				"used_data_storage_size":         v.UsedDataStorageSize,
				"cpu_count":                      v.CpuCount,
			}

			serverInstanceList = append(serverInstanceList, ins)
		}
		instance["cloud_mysql_server_instance_list"] = serverInstanceList
		resources = append(resources, instance)
	}

	return resources, nil
}

func getMysqlList(d *schema.ResourceData, config *conn.ProviderConfig) ([]map[string]interface{}, error) {
	client := config.Client

	reqParams := &vmysql.GetCloudMysqlInstanceListRequest{
		RegionCode: &config.RegionCode,
	}

	if v, ok := d.GetOk("id"); ok {
		reqParams.CloudMysqlInstanceNoList = []*string{ncloud.String(v.(string))}
	}
	if v, ok := d.GetOk("zone_code"); ok {
		reqParams.ZoneCode = ncloud.String(v.(string))
	}
	if v, ok := d.GetOk("vpc_no"); ok {
		reqParams.VpcNo = ncloud.String(v.(string))
	}
	if v, ok := d.GetOk("subnet_no"); ok {
		reqParams.SubnetNo = ncloud.String(v.(string))
	}
	if v, ok := d.GetOk("service_name"); ok {
		reqParams.CloudMysqlServiceName = ncloud.String(v.(string))
	}

	LogCommonRequest("getCloudMysqlList", reqParams)

	resp, err := client.Vmysql.V2Api.GetCloudMysqlInstanceList(reqParams)
	if err != nil {
		LogErrorResponse("getCloudMysqlList", err, reqParams)
		return nil, err
	}
	LogResponse("getCloudMysqlList", resp)

	var resources []map[string]interface{}
	for _, r := range resp.CloudMysqlInstanceList {
		instance := map[string]interface{}{
			"id":                           *r.CloudMysqlInstanceNo,
			"instance_no":                     *r.CloudMysqlInstanceNo,
			"service_name":                 *r.CloudMysqlServiceName,
			"is_ha":                        *r.IsHa,
			"is_multi_zone":                *r.IsMultiZone,
		}

		resources = append(resources, instance)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudMysql().Schema)
	}

	return resources, nil
}
