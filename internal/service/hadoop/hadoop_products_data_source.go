package hadoop

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vhadoop"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"time"
)

var (
	_ datasource.DataSource              = &hadoopProductsDataSource{}
	_ datasource.DataSourceWithConfigure = &hadoopProductsDataSource{}
)

func NewHadoopProductsDataSource() datasource.DataSource {
	return &hadoopProductsDataSource{}
}

type hadoopProductsDataSource struct {
	config *conn.ProviderConfig
}

func (h *hadoopProductsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*conn.ProviderConfig)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ProviderConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	h.config = config
}

func (h *hadoopProductsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hadoop_products"
}

func (h *hadoopProductsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"image_product_code": schema.StringAttribute{
				Required: true,
			},
			"product_code": schema.StringAttribute{
				Optional: true,
			},
			"infra_resource_detail_type_code": schema.StringAttribute{
				Optional: true,
			},
			"exclusion_product_code": schema.StringAttribute{
				Optional: true,
			},
			"output_file": schema.StringAttribute{
				Optional: true,
			},
		},
		Blocks: map[string]schema.Block{
			"filter": common.DataSourceFiltersBlock(),
			"product_list": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"product_name": schema.StringAttribute{
							Computed: true,
						},
						"product_code": schema.StringAttribute{
							Computed: true,
						},
						"product_type": schema.StringAttribute{
							Computed: true,
						},
						"product_description": schema.StringAttribute{
							Computed: true,
						},
						"infra_resource_type": schema.StringAttribute{
							Computed: true,
						},
						"infra_resource_detail_type": schema.StringAttribute{
							Computed: true,
						},
						"cpu_count": schema.Int64Attribute{
							Computed: true,
						},
						"memory_size": schema.Int64Attribute{
							Computed: true,
						},
						"disk_type": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (h *hadoopProductsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data hadoopProductsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vhadoop.GetCloudHadoopProductListRequest{
		RegionCode: &h.config.RegionCode,
	}

	if !data.ImageProductCode.IsNull() && !data.ImageProductCode.IsUnknown() {
		reqParams.CloudHadoopImageProductCode = data.ImageProductCode.ValueStringPointer()
	}

	if !data.ProductCode.IsNull() && !data.ProductCode.IsUnknown() {
		reqParams.ProductCode = data.ProductCode.ValueStringPointer()
	}

	if !data.InfraResourceDetailTypeCode.IsNull() && !data.InfraResourceDetailTypeCode.IsUnknown() {
		reqParams.InfraResourceDetailTypeCode = data.InfraResourceDetailTypeCode.ValueStringPointer()
	}

	if !data.ExclusionProductCode.IsNull() && !data.ExclusionProductCode.IsUnknown() {
		reqParams.ExclusionProductCode = data.ExclusionProductCode.ValueStringPointer()
	}

	tflog.Info(ctx, "GetHadoopProductList", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})

	productsResp, err := h.config.Client.Vhadoop.V2Api.GetCloudHadoopProductList(reqParams)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetHadoopProductList",
			fmt.Sprintf("error: %s, reqParams: %s", err.Error(), common.MarshalUncheckedString(reqParams)),
		)
		resp.Diagnostics.Append(diags...)
		return
	}
	tflog.Info(ctx, "GetHadoopProductList response", map[string]any{
		"productResponse": common.MarshalUncheckedString(productsResp),
	})

	hadoopProductsList := flattenHadoopProductList(productsResp.ProductList)
	fillteredList := common.FilterModels(ctx, data.Filters, hadoopProductsList)
	data.refreshFromOutput(ctx, fillteredList)
	if !data.OutputFile.IsNull() && data.OutputFile.String() != "" {
		outputPath := data.OutputFile.ValueString()

		if diags := writeHadoopProductsToFile(outputPath, data.ProductList); diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}
	data.ID = types.StringValue(time.Now().UTC().String())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func writeHadoopProductsToFile(path string, products types.List) diag.Diagnostics {
	var hadoopProducts []hadoopProductsToJsonConvert
	var diags diag.Diagnostics

	for _, product := range products.Elements() {
		hadoopProduct := hadoopProductsToJsonConvert{}
		if err := json.Unmarshal([]byte(product.String()), &hadoopProduct); err != nil {
			diags.AddError(
				"Unmarshal",
				fmt.Sprintf("error: %s", err.Error()),
			)
			return diags
		}
		hadoopProducts = append(hadoopProducts, hadoopProduct)
	}

	if err := common.WriteToFile(path, hadoopProducts); err != nil {
		diags.AddError(
			"WriteToFile",
			fmt.Sprintf("error: %s", err.Error()),
		)
		return diags
	}
	return nil
}

func flattenHadoopProductList(products []*vhadoop.Product) []*hadoopProductModel {
	var outputs []*hadoopProductModel

	for _, v := range products {
		var output hadoopProductModel
		output.refreshFromOutput(v)

		outputs = append(outputs, &output)
	}
	return outputs
}

func (h *hadoopProductsDataSourceModel) refreshFromOutput(ctx context.Context, output []*hadoopProductModel) {
	h.ProductList, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: hadoopProductModel{}.attrTypes()}, output)
}

type hadoopProductsDataSourceModel struct {
	ID                          types.String `tfsdk:"id"`
	ImageProductCode            types.String `tfsdk:"image_product_code"`
	ProductCode                 types.String `tfsdk:"product_code"`
	InfraResourceDetailTypeCode types.String `tfsdk:"infra_resource_detail_type_code"`
	ExclusionProductCode        types.String `tfsdk:"exclusion_product_code"`
	OutputFile                  types.String `tfsdk:"output_file"`
	ProductList                 types.List   `tfsdk:"product_list"`
	Filters                     types.Set    `tfsdk:"filter"`
}

type hadoopProductModel struct {
	ProductName             types.String `tfsdk:"product_name"`
	ProductCode             types.String `tfsdk:"product_code"`
	ProductType             types.String `tfsdk:"product_type"`
	ProductDescription      types.String `tfsdk:"product_description"`
	InfraResourceType       types.String `tfsdk:"infra_resource_type"`
	InfraResourceDetailType types.String `tfsdk:"infra_resource_detail_type"`
	CpuCount                types.Int64  `tfsdk:"cpu_count"`
	MemorySize              types.Int64  `tfsdk:"memory_size"`
	DiskType                types.String `tfsdk:"disk_type"`
}

type hadoopProductsToJsonConvert struct {
	ProductName             string `json:"product_name"`
	ProductCode             string `json:"product_code"`
	ProductType             string `json:"product_type"`
	ProductDescription      string `json:"product_description"`
	InfraResourceType       string `json:"infra_resource_type"`
	InfraResourceDetailType string `json:"infra_resource_detail_type"`
	CpuCount                int64  `json:"cpu_count"`
	MemorySize              int64  `json:"memory_size"`
	DiskType                string `json:"disk_type"`
}

func (_ hadoopProductModel) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"product_name":               types.StringType,
		"product_code":               types.StringType,
		"product_type":               types.StringType,
		"product_description":        types.StringType,
		"infra_resource_type":        types.StringType,
		"infra_resource_detail_type": types.StringType,
		"cpu_count":                  types.Int64Type,
		"memory_size":                types.Int64Type,
		"disk_type":                  types.StringType,
	}
}

func (h *hadoopProductModel) refreshFromOutput(output *vhadoop.Product) {
	h.ProductName = types.StringPointerValue(output.ProductName)
	h.ProductCode = types.StringPointerValue(output.ProductCode)
	h.ProductType = types.StringPointerValue(output.ProductType.Code)
	h.ProductDescription = types.StringPointerValue(output.ProductDescription)
	h.InfraResourceType = types.StringPointerValue(output.InfraResourceType.Code)
	h.InfraResourceDetailType = types.StringPointerValue(output.InfraResourceDetailType.Code)
	h.CpuCount = types.Int64Value(int64(*output.CpuCount))
	h.MemorySize = types.Int64PointerValue(output.MemorySize)
	h.DiskType = types.StringPointerValue(output.DiskType.Code)
}
