package nks

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/provider"
)

func init() {
	RegisterDataSource("ncloud_nks_server_products", dataSourceNcloudNKSServerProducts())
}

func dataSourceNcloudNKSServerProducts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudNKSServerProductsRead,
		Schema: map[string]*schema.Schema{
			"filter": DataSourceFiltersSchema(),
			"products": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"detail": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cpu_count": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"memory_size": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"gpu_count": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"gpu_memory_size": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"product_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"product_code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"product_korean_desc": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"product_english_desc": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"software_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceNcloudNKSServerProductsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return NotSupportClassic("datasource `ncloud_nks_server_products`")
	}

	resources, err := getNKSServerProducts(config, d)
	if err != nil {
		LogErrorResponse("GetNKSServerProducts", err, "")
		return err
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudNKSServerProducts().Schema["products"].Elem.(*schema.Resource).Schema)
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("products", resources); err != nil {
		return fmt.Errorf("Error setting Codes: %s", err)
	}

	return nil
}

func getNKSServerProducts(config *ProviderConfig, d *schema.ResourceData) ([]map[string]interface{}, error) {
	LogCommonRequest("GetNKSServerProducts", "")

	softwareCode := StringPtrOrNil(d.GetOk("software_code"))
	zoneCode := StringPtrOrNil(d.GetOk("zone"))

	opt := make(map[string]interface{})
	opt["zoneCode"] = zoneCode
	resp, err := config.Client.Vnks.V2Api.OptionServerProductCodeGet(context.Background(), softwareCode, opt)

	if err != nil {
		LogErrorResponse("GetNKSServerProducts", err, "")
		return nil, err
	}

	LogResponse("GetNKSServerProducts", resp)

	resources := []map[string]interface{}{}

	for _, r := range *resp {
		instance := map[string]interface{}{
			"label": ncloud.StringValue(r.Detail.ProductName),
			"value": ncloud.StringValue(r.Detail.ProductCode),
			"detail": []map[string]interface{}{
				{
					"product_type":         ncloud.StringValue(r.Detail.ProductType2Code),
					"product_code":         ncloud.StringValue(r.Detail.ProductCode),
					"product_korean_desc":  ncloud.StringValue(r.Detail.ProductKoreanDesc),
					"product_english_desc": ncloud.StringValue(r.Detail.ProductEnglishDesc),
					"cpu_count":            strconv.Itoa(int(ncloud.Int32Value(r.Detail.CpuCount))),
					"memory_size":          strconv.Itoa(int(ncloud.Int32Value(r.Detail.MemorySizeGb))) + "GB",
					"gpu_count":            strconv.Itoa(int(ncloud.Int32Value(r.Detail.GpuCount))),
					"gpu_memory_size":      strconv.Itoa(int(ncloud.Int32Value(r.Detail.GpuMemorySizeGb))) + "GB",
				},
			},
		}

		resources = append(resources, instance)
	}

	return resources, nil
}
