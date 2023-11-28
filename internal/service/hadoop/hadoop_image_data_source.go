package hadoop

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vhadoop"
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
	_ datasource.DataSource              = &hadoopImageDataSource{}
	_ datasource.DataSourceWithConfigure = &hadoopImageDataSource{}
)

func NewHadoopImageDataSource() datasource.DataSource {
	return &hadoopImageDataSource{}
}

type hadoopImageDataSource struct {
	config *conn.ProviderConfig
}

func (h *hadoopImageDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (h *hadoopImageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hadoop_image"
}

func (h *hadoopImageDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"product_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"exclusion_product_code": schema.StringAttribute{
				Optional: true,
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
			"base_block_storage_size": schema.Int64Attribute{
				Computed: true,
			},
			"platform_type": schema.StringAttribute{
				Computed: true,
			},
			"os_information": schema.StringAttribute{
				Computed: true,
			},
			"generation_code": schema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"filter": common.DataSourceFiltersBlock(),
		},
	}
}

func (h *hadoopImageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data hadoopImageDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vhadoop.GetCloudHadoopImageProductListRequest{
		RegionCode: &h.config.RegionCode,
	}

	if !data.ProductCode.IsNull() && !data.ProductCode.IsUnknown() {
		reqParams.ProductCode = data.ProductCode.ValueStringPointer()
	}

	if !data.ExclusionProductCode.IsNull() && !data.ExclusionProductCode.IsUnknown() {
		reqParams.ExclusionProductCode = data.ExclusionProductCode.ValueStringPointer()
	}

	tflog.Info(ctx, "GetHadoopImageProductList", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})
	imageProductResp, err := h.config.Client.Vhadoop.V2Api.GetCloudHadoopImageProductList(reqParams)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetHadoopImageProductList",
			fmt.Sprintf("error: %s, reqParams: %s", err.Error(), common.MarshalUncheckedString(reqParams)),
		)
		resp.Diagnostics.Append(diags...)
		return
	}
	tflog.Info(ctx, "GetHadoopImageProductList response", map[string]any{
		"imageProductResponse": common.MarshalUncheckedString(imageProductResp),
	})

	imageProductList := flattenHadoopImageList(imageProductResp.ProductList)
	fillteredList := common.FilterModels(ctx, data.Filters, imageProductList)

	if err := verify.ValidateOneResult(len(fillteredList)); err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetHadoopImageProductList result validation",
			err.Error(),
		)
		resp.Diagnostics.Append(diags...)
		return
	}

	state := fillteredList[0]
	state.Filters = data.Filters

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenHadoopImageList(imageProducts []*vhadoop.Product) []*hadoopImageDataSourceModel {
	var outputs []*hadoopImageDataSourceModel

	for _, v := range imageProducts {
		var output hadoopImageDataSourceModel

		output.refreshFromOutput(v)

		outputs = append(outputs, &output)
	}
	return outputs
}

type hadoopImageDataSourceModel struct {
	ID                   types.String `tfsdk:"id"`
	ProductCode          types.String `tfsdk:"product_code"`
	ExclusionProductCode types.String `tfsdk:"exclusion_product_code"`
	ProductName          types.String `tfsdk:"product_name"`
	ProductType          types.String `tfsdk:"product_type"`
	ProductDescription   types.String `tfsdk:"product_description"`
	InfraResourceType    types.String `tfsdk:"infra_resource_type"`
	BaseBlockStorageSize types.Int64  `tfsdk:"base_block_storage_size"`
	PlatformType         types.String `tfsdk:"platform_type"`
	OsInformation        types.String `tfsdk:"os_information"`
	GenerationCode       types.String `tfsdk:"generation_code"`
	Filters              types.Set    `tfsdk:"filter"`
}

func (h *hadoopImageDataSourceModel) refreshFromOutput(output *vhadoop.Product) {
	h.ID = types.StringPointerValue(output.ProductCode)
	h.ProductCode = types.StringPointerValue(output.ProductCode)
	h.ProductName = types.StringPointerValue(output.ProductName)
	h.ProductType = types.StringPointerValue(output.ProductType.Code)
	h.ProductDescription = types.StringPointerValue(output.ProductDescription)
	h.InfraResourceType = types.StringPointerValue(output.InfraResourceType.Code)
	h.BaseBlockStorageSize = types.Int64PointerValue(output.BaseBlockStorageSize)
	h.PlatformType = types.StringPointerValue(output.PlatformType.Code)
	h.OsInformation = types.StringPointerValue(output.OsInformation)
	h.GenerationCode = types.StringPointerValue(output.GenerationCode)
}
