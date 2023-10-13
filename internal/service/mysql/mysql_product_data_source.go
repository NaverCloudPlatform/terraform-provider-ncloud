package mysql

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmysql"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

var (
	_ datasource.DataSource              = &mysqlProductDataSource{}
	_ datasource.DataSourceWithConfigure = &mysqlProductDataSource{}
)

func NewMysqlProductDataSource() datasource.DataSource {
	return &mysqlProductDataSource{}
}

type mysqlProductDataSource struct {
	config *conn.ProviderConfig
}

func (m *mysqlProductDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	m.config = config
}

func (m *mysqlProductDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mysql_product"
}

func (m *mysqlProductDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"cloud_mysql_image_product_code": schema.StringAttribute{
				Required: true,
			},
			"product_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"exclusion_product_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"product_name": schema.StringAttribute{
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
		Blocks: map[string]schema.Block{
			"filter": common.DataSourceFiltersBlock(),
		},
	}
}

func (m *mysqlProductDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mysqlProductDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vmysql.GetCloudMysqlProductListRequest{
		RegionCode:                 &m.config.RegionCode,
		CloudMysqlImageProductCode: data.CloudMysqlImageProductCode.ValueStringPointer(),
	}

	if !data.ProductCode.IsNull() && !data.ProductCode.IsUnknown() {
		reqParams.ProductCode = data.ProductCode.ValueStringPointer()
	}

	if !data.ExclusionProductCode.IsNull() && !data.ExclusionProductCode.IsUnknown() {
		reqParams.ExclusionProductCode = data.ExclusionProductCode.ValueStringPointer()
	}

	tflog.Info(ctx, "GetMysqlProductList", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})

	mysqlProductResp, err := m.config.Client.Vmysql.V2Api.GetCloudMysqlProductList(reqParams)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetMysqlProductList",
			fmt.Sprintf("error: %s, reqParams: %s", err.Error(), common.MarshalUncheckedString(reqParams)),
		)
		resp.Diagnostics.Append(diags...)
		return
	}

	tflog.Info(ctx, "GetMysqlProductList response", map[string]any{
		"mysqlProductResponse": common.MarshalUncheckedString(mysqlProductResp),
	})

	mysqlProductList := flattenMysqlProducts(ctx, mysqlProductResp.ProductList)

	fillteredList := common.FilterModels(ctx, data.Filters, mysqlProductList)

	if err := verify.ValidateOneResult(len(fillteredList)); err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetMysqlProductList result more than one",
			err.Error(),
		)
		resp.Diagnostics.Append(diags...)
		return
	}

	state := fillteredList[0]
	state.Filters = data.Filters

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenMysqlProducts(ctx context.Context, list []*vmysql.Product) []*mysqlProductDataSourceModel {
	var outputs []*mysqlProductDataSourceModel

	for _, v := range list {
		var output mysqlProductDataSourceModel
		output.refreshFromOutput(v)

		outputs = append(outputs, &output)
	}
	return outputs
}

type mysqlProductDataSourceModel struct {
	ID                         types.String `tfsdk:"id"`
	CloudMysqlImageProductCode types.String `tfsdk:"cloud_mysql_image_product_code"`
	ProductCode                types.String `tfsdk:"product_code"`
	ProductName                types.String `tfsdk:"product_name"`
	ProductType                types.String `tfsdk:"product_type"`
	ProductDescription         types.String `tfsdk:"product_description"`
	InfraResourceType          types.String `tfsdk:"infra_resource_type"`
	CpuCount                   types.Int64  `tfsdk:"cpu_count"`
	MemorySize                 types.Int64  `tfsdk:"memory_size"`
	DiskType                   types.String `tfsdk:"disk_type"`
	ExclusionProductCode       types.String `tfsdk:"exclusion_product_code"`
	Filters                    types.Set    `tfsdk:"filter"`
}

func (m *mysqlProductDataSourceModel) refreshFromOutput(output *vmysql.Product) {
	m.ProductCode = types.StringPointerValue(output.ProductCode)
	m.ProductName = types.StringPointerValue(output.ProductName)
	m.ProductType = types.StringPointerValue(output.ProductType.Code)
	m.ProductDescription = types.StringPointerValue(output.ProductDescription)
	m.InfraResourceType = types.StringPointerValue(output.InfraResourceType.Code)
	m.CpuCount = types.Int64Value(int64(*output.CpuCount))
	m.MemorySize = types.Int64PointerValue(output.MemorySize)
	m.DiskType = types.StringPointerValue(output.DiskType.Code)
	m.ID = types.StringPointerValue(output.ProductCode)
}
