package mongodb

import (
	"context"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmongodb"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

var (
	_ datasource.DataSource              = &mongodbProductDataSource{}
	_ datasource.DataSourceWithConfigure = &mongodbProductDataSource{}
)

func NewMongoDbProductDataSource() datasource.DataSource {
	return &mongodbProductDataSource{}
}

type mongodbProductDataSource struct {
	config *conn.ProviderConfig
}

func (m *mongodbProductDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mongodb_product"
}

func (m *mongodbProductDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (m *mongodbProductDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
				Computed: true,
			},
			"product_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"infra_resource_detail_type_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"exclusion_product_code": schema.StringAttribute{
				Optional: true,
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

func (m *mongodbProductDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mongodbProductDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vmongodb.GetCloudMongoDbProductListRequest{
		RegionCode:                   &m.config.RegionCode,
		CloudMongoDbImageProductCode: data.CloudMongoDbImageProductCode.ValueStringPointer(),
	}

	if !data.ProductCode.IsNull() && !data.ProductCode.IsUnknown() {
		reqParams.ProductCode = data.ProductCode.ValueStringPointer()
	}

	if !data.ExclusionProductCode.IsNull() && !data.ExclusionProductCode.IsUnknown() {
		reqParams.ExclusionProductCode = data.ExclusionProductCode.ValueStringPointer()
	}

	tflog.Info(ctx, "GetMongoDbProductList", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})

	mongodbProductResp, err := m.config.Client.Vmongodb.V2Api.GetCloudMongoDbProductList(reqParams)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetMongoDbProductList",
			fmt.Sprintf("error: %s, reqParams: %s", err.Error(), common.MarshalUncheckedString(reqParams)),
		)
		resp.Diagnostics.Append(diags...)
		return
	}

	tflog.Info(ctx, "GetMongoDbProductList response", map[string]any{
		"mongodbProductResponse": common.MarshalUncheckedString(mongodbProductResp),
	})

	mongodbProductList := flattenMongoDbProducts(ctx, mongodbProductResp.ProductList)

	fillteredList := common.FilterModels(ctx, data.Filters, mongodbProductList)

	if err := verify.ValidateOneResult(len(fillteredList)); err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetVpcList result validation",
			err.Error(),
		)
		resp.Diagnostics.Append(diags...)
		return
	}

	state := fillteredList[0]
	state.Filters = data.Filters

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenMongoDbProducts(ctx context.Context, list []*vmongodb.Product) []*mongodbProductDataSourceModel {
	var outputs []*mongodbProductDataSourceModel

	for _, v := range list {
		var output mongodbProductDataSourceModel
		output.refreshFromOutput(v)

		outputs = append(outputs, &output)
	}
	return outputs
}

type mongodbProductDataSourceModel struct {
	ID                           types.String `tfsdk:"id"`
	CloudMongoDbImageProductCode types.String `tfsdk:"image_product_code"`
	ProductCode                  types.String `tfsdk:"product_code"`
	ProductName                  types.String `tfsdk:"product_name"`
	ExclusionProductCode         types.String `tfsdk:"exclusion_product_code"`
	InfraResourceDetailType      types.String `tfsdk:"infra_resource_detail_type_code"`
	ProductType                  types.String `tfsdk:"product_type"`
	ProductDescription           types.String `tfsdk:"product_description"`
	InfraResourceType            types.String `tfsdk:"infra_resource_type"`
	CpuCount                     types.Int64  `tfsdk:"cpu_count"`
	MemorySize                   types.Int64  `tfsdk:"memory_size"`
	DiskType                     types.String `tfsdk:"disk_type"`
	Filters                      types.Set    `tfsdk:"filter"`
}

func (m *mongodbProductDataSourceModel) refreshFromOutput(output *vmongodb.Product) {
	m.ID = types.StringPointerValue(output.ProductCode)
	m.ProductCode = types.StringPointerValue(output.ProductCode)
	m.ProductName = types.StringPointerValue(output.ProductName)
	m.ProductType = types.StringPointerValue(output.ProductType.Code)
	m.ProductDescription = types.StringPointerValue(output.ProductDescription)
	m.InfraResourceType = types.StringPointerValue(output.InfraResourceType.Code)
	m.InfraResourceDetailType = types.StringPointerValue(output.InfraResourceDetailType.Code)
	m.CpuCount = types.Int64Value(int64(*output.CpuCount))
	m.MemorySize = types.Int64PointerValue(output.MemorySize)
	m.DiskType = types.StringPointerValue(output.DiskType.Code)
}
