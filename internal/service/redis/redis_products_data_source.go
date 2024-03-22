package redis

import (
	"context"
	"encoding/json"
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
	_ datasource.DataSource              = &redisProductsDataSource{}
	_ datasource.DataSourceWithConfigure = &redisProductsDataSource{}
)

func NewRedisProductsDataSource() datasource.DataSource {
	return &redisProductsDataSource{}
}

type redisProductsDataSource struct {
	config *conn.ProviderConfig
}

func (r *redisProductsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_redis_products"
}

func (r *redisProductsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"redis_image_product_code": schema.StringAttribute{
				Required: true,
			},
			"output_file": schema.StringAttribute{
				Optional: true,
			},
			"product_list": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"product_code": schema.StringAttribute{
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
				},
				Computed: true,
			},
		},

		Blocks: map[string]schema.Block{
			"filter": common.DataSourceFiltersBlock(),
		},
	}
}

func (r *redisProductsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *redisProductsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data redisProductList

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

	reqParams := &vredis.GetCloudRedisProductListRequest{
		RegionCode:                 &r.config.RegionCode,
		CloudRedisImageProductCode: data.CloudRedisImageProductCode.ValueStringPointer(),
	}
	tflog.Info(ctx, "GetRedisProductList reqParams="+common.MarshalUncheckedString(reqParams))

	redisProductResp, err := r.config.Client.Vredis.V2Api.GetCloudRedisProductList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "GetRedisProductList response="+common.MarshalUncheckedString(redisProductResp))

	if redisProductResp == nil || len(redisProductResp.ProductList) < 1 {
		resp.Diagnostics.AddError("READING ERROR", "no result.")
		return
	}

	redisProductList := flattenRedisProduct(redisProductResp.ProductList)
	fillteredList := common.FilterModels(ctx, data.Filters, redisProductList)

	data.refreshFromOutput(ctx, fillteredList)

	if !data.OutputFile.IsNull() && data.OutputFile.String() != "" {
		outputPath := data.OutputFile.ValueString()

		if convertedList, err := convertProductsToJsonStruct(data.ProductList.Elements()); err != nil {
			resp.Diagnostics.AddError("OUTPUT FILE ERROR", err.Error())
		} else if err := common.WriteToFile(outputPath, convertedList); err != nil {
			resp.Diagnostics.AddError("OUTPUT FILE ERROR", err.Error())
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func convertProductsToJsonStruct(images []attr.Value) ([]redisProductsToJsonConvert, error) {
	var redisProductsToConvert = []redisProductsToJsonConvert{}

	for _, image := range images {
		imageJasn := redisProductsToJsonConvert{}
		if err := json.Unmarshal([]byte(image.String()), &imageJasn); err != nil {
			return nil, err
		}
		redisProductsToConvert = append(redisProductsToConvert, imageJasn)
	}

	return redisProductsToConvert, nil
}

func flattenRedisProduct(list []*vredis.Product) []*redisProductModel {
	var outputs []*redisProductModel

	for _, v := range list {
		var output redisProductModel
		output.refreshFromOutput(v)

		outputs = append(outputs, &output)
	}
	return outputs
}

func (r *redisProductList) refreshFromOutput(ctx context.Context, list []*redisProductModel) {
	productListValue, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: redisProductModel{}.attrTypes()}, list)
	r.ProductList = productListValue
	r.ID = types.StringValue("")
}

type redisProductList struct {
	ID                         types.String `tfsdk:"id"`
	CloudRedisImageProductCode types.String `tfsdk:"redis_image_product_code"`
	ProductList                types.List   `tfsdk:"product_list"`
	OutputFile                 types.String `tfsdk:"output_file"`
	Filters                    types.Set    `tfsdk:"filter"`
}

type redisProductModel struct {
	ProductCode        types.String `tfsdk:"product_code"`
	ProductName        types.String `tfsdk:"product_name"`
	ProductType        types.String `tfsdk:"product_type"`
	ProductDescription types.String `tfsdk:"product_description"`
	InfraResourceType  types.String `tfsdk:"infra_resource_type"`
	CpuCount           types.Int64  `tfsdk:"cpu_count"`
	MemorySize         types.Int64  `tfsdk:"memory_size"`
	DiskType           types.String `tfsdk:"disk_type"`
}

type redisProductsToJsonConvert struct {
	ProductCode        string `json:"product_code"`
	ProductName        string `json:"product_name"`
	ProductType        string `json:"product_type"`
	ProductDescription string `json:"product_description"`
	InfraResourceType  string `json:"infra_resource_type"`
	CpuCount           int    `json:"cpu_count"`
	MemorySize         int    `json:"memory_size"`
	DiskType           string `json:"disk_type"`
}

func (r redisProductModel) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"product_code":        types.StringType,
		"product_name":        types.StringType,
		"product_type":        types.StringType,
		"infra_resource_type": types.StringType,
		"cpu_count":           types.Int64Type,
		"memory_size":         types.Int64Type,
		"disk_type":           types.StringType,
		"product_description": types.StringType,
	}
}

func (r *redisProductModel) refreshFromOutput(output *vredis.Product) {
	r.ProductCode = types.StringPointerValue(output.ProductCode)
	r.ProductName = types.StringPointerValue(output.ProductName)
	r.ProductType = types.StringPointerValue(output.ProductType.Code)
	r.ProductDescription = types.StringPointerValue(output.ProductDescription)
	r.InfraResourceType = types.StringPointerValue(output.InfraResourceType.Code)
	r.CpuCount = common.Int64ValueFromInt32(output.CpuCount)
	r.MemorySize = types.Int64PointerValue(output.MemorySize)
	r.DiskType = types.StringPointerValue(output.DiskType.Code)
}
