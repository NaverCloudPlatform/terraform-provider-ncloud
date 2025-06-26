package ses

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vses2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudSESNodeProduct() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudSESNodeProductRead,
		Schema: map[string]*schema.Schema{
			"filter": DataSourceFiltersSchema(),
			"codes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_count": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"os_image_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet_no": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func dataSourceNcloudSESNodeProductRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	resources, err := getSESNodeProduct(config, d)
	if err != nil {
		LogErrorResponse("GetSESNodeProduct", err, "")
		return err
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudSESNodeProduct().Schema["codes"].Elem.(*schema.Resource).Schema)
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("codes", resources); err != nil {
		return fmt.Errorf("Error setting Codes: %s", err)
	}

	return nil
}

func getSESNodeProduct(config *conn.ProviderConfig, d *schema.ResourceData) ([]map[string]interface{}, error) {
	LogCommonRequest("GetSESSoftwareProduct", "")

	reqParams := &vses2.V2ApiGetNodeProductListWithGetMethodUsingGETOpts{
		SoftwareProductCode: *StringPtrOrNil(d.GetOk("os_image_code")),
		SubnetNo:            *Int32PtrOrNil(d.GetOk("subnet_no")),
	}
	resp, _, err := config.Client.Vses.V2Api.GetNodeProductListWithGetMethodUsingGET(context.Background(), reqParams)

	if err != nil {
		LogErrorResponse("GetSESNodeProduct", err, "")
		return nil, err
	}

	LogResponse("GetSESNodeProduct", resp)

	resources := []map[string]interface{}{}

	for _, r := range resp.Result.ProductList {
		instance := map[string]interface{}{
			"id":          ncloud.StringValue(&r.ProductCode),
			"name":        ncloud.StringValue(&r.ProductEnglishDesc),
			"cpu_count":   ncloud.StringValue(&r.CpuCount),
			"memory_size": strconv.Itoa(int(ncloud.Int64Value(&r.MemorySize)/int64(1024*1024*1024))) + "GB",
		}

		resources = append(resources, instance)
	}

	return resources, nil
}
