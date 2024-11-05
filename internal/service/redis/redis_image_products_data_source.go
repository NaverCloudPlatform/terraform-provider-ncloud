package redis

import (
	"context"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vredis"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

var (
	_ datasource.DataSource              = &redisImageProductsDataSource{}
	_ datasource.DataSourceWithConfigure = &redisImageProductsDataSource{}
)

func NewRedisImageProductsDataSource() datasource.DataSource {
	return &redisImageProductsDataSource{}
}

type redisImageProductsDataSource struct {
	config *conn.ProviderConfig
}

func (r *redisImageProductsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_redis_image_products"
}

func (r *redisImageProductsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"output_file": schema.StringAttribute{
				Optional: true,
			},
			"image_product_list": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"product_code": schema.StringAttribute{
							Computed: true,
						},
						"generation_code": schema.StringAttribute{
							Computed: true,
						},
						"product_name": schema.StringAttribute{
							Computed: true,
						},
						"product_type": schema.StringAttribute{
							Computed: true,
						},
						"platform_type": schema.StringAttribute{
							Computed: true,
						},
						"os_information": schema.StringAttribute{
							Computed: true,
						},
						"engine_version_code": schema.StringAttribute{
							Computed: true,
						},
					},
				},
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"filter": common.DataSourceFiltersBlock(),
		},
	}
}

func (r *redisImageProductsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.config = config
}

func (r *redisImageProductsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data redisImageProductsDataSourceModel

	if !r.config.SupportVPC {
		resp.Diagnostics.AddError(
			"NOT SUPPORT CLASSIC",
			"does not support CLASSIC. only VPC.",
		)
		return
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vredis.GetCloudRedisImageProductListRequest{
		RegionCode: &r.config.RegionCode,
	}
	tflog.Info(ctx, "GetRedisImageProductList reqParams="+common.MarshalUncheckedString(reqParams))

	redisImageProductResp, err := r.config.Client.Vredis.V2Api.GetCloudRedisImageProductList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "GetRedisImageProductList response="+common.MarshalUncheckedString(redisImageProductResp))

	if redisImageProductResp == nil || len(redisImageProductResp.ProductList) < 1 {
		resp.Diagnostics.AddError("READING ERROR", "no result.")
		return
	}

	redisImageProductList := flattenRedisImageProduct(redisImageProductResp.ProductList)
	fillteredList := common.FilterModels(ctx, data.Filters, redisImageProductList)
	data.refreshFromOutput(ctx, fillteredList)

	if !data.OutputFile.IsNull() && data.OutputFile.String() != "" {
		outputPath := data.OutputFile.ValueString()

		if err := common.WriteImageProductToFile(outputPath, data.ImageProductList); err != nil {
			resp.Diagnostics.AddError("OUTPUT FILE ERROR", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenRedisImageProduct(list []*vredis.Product) []*redisImageProduct {
	var outputs []*redisImageProduct

	for _, v := range list {
		var output redisImageProduct
		output.refreshFromOutput(v)

		outputs = append(outputs, &output)
	}
	return outputs
}

func (r *redisImageProductsDataSourceModel) refreshFromOutput(ctx context.Context, list []*redisImageProduct) {
	imageProductListValue, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: redisImageProduct{}.attrTypes()}, list)
	r.ImageProductList = imageProductListValue
	r.ID = types.StringValue("")
}

type redisImageProductsDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	ImageProductList types.List   `tfsdk:"image_product_list"`
	OutputFile       types.String `tfsdk:"output_file"`
	Filters          types.Set    `tfsdk:"filter"`
}

type redisImageProduct struct {
	ProductCode       types.String `tfsdk:"product_code"`
	GenerationCode    types.String `tfsdk:"generation_code"`
	ProductName       types.String `tfsdk:"product_name"`
	ProductType       types.String `tfsdk:"product_type"`
	PlatformType      types.String `tfsdk:"platform_type"`
	OsInformation     types.String `tfsdk:"os_information"`
	EngineVersionCode types.String `tfsdk:"engine_version_code"`
}

func (r redisImageProduct) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"product_code":        types.StringType,
		"generation_code":     types.StringType,
		"product_name":        types.StringType,
		"product_type":        types.StringType,
		"platform_type":       types.StringType,
		"os_information":      types.StringType,
		"engine_version_code": types.StringType,
	}
}

func (r *redisImageProduct) refreshFromOutput(output *vredis.Product) {
	r.ProductCode = types.StringPointerValue(output.ProductCode)
	r.GenerationCode = types.StringPointerValue(output.GenerationCode)
	r.ProductName = types.StringPointerValue(output.ProductName)
	r.ProductType = types.StringPointerValue(output.ProductType.Code)
	r.PlatformType = types.StringPointerValue(output.PlatformType.Code)
	r.OsInformation = types.StringPointerValue(output.OsInformation)
	r.EngineVersionCode = types.StringPointerValue(output.EngineVersionCode)
}
