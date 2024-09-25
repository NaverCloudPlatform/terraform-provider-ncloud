package mysql

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmysql"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

var (
	_ datasource.DataSource              = &mysqlProductsDataSource{}
	_ datasource.DataSourceWithConfigure = &mysqlProductsDataSource{}
)

func NewMysqlProductsDataSource() datasource.DataSource {
	return &mysqlProductsDataSource{}
}

type mysqlProductsDataSource struct {
	config *conn.ProviderConfig
}

func (m *mysqlProductsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mysql_products"
}

func (m *mysqlProductsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"image_product_code": schema.StringAttribute{
				Required: true,
			},
			"output_file": schema.StringAttribute{
				Optional: true,
			},
			"product_list": schema.ListNestedAttribute{
				Computed: true,
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
			},
		},
		Blocks: map[string]schema.Block{
			"filter": common.DataSourceFiltersBlock(),
		},
	}
}

func (m *mysqlProductsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (m *mysqlProductsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mysqlProductList

	if !m.config.SupportVPC {
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

	reqParams := &vmysql.GetCloudMysqlProductListRequest{
		RegionCode:                 &m.config.RegionCode,
		CloudMysqlImageProductCode: data.CloudMysqlImageProductCode.ValueStringPointer(),
	}
	tflog.Info(ctx, "GetMysqlProductsList reqParams="+common.MarshalUncheckedString(reqParams))

	mysqlProductResp, err := m.config.Client.Vmysql.V2Api.GetCloudMysqlProductList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "GetMysqlProductsList response="+common.MarshalUncheckedString(mysqlProductResp))

	if mysqlProductResp == nil || len(mysqlProductResp.ProductList) < 1 {
		resp.Diagnostics.AddError("READING ERROR", "no result.")
		return
	}

	mysqlProductList := flattenMysqlProduct(mysqlProductResp.ProductList)
	fillteredList := common.FilterModels(ctx, data.Filters, mysqlProductList)

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

func convertProductsToJsonStruct(images []attr.Value) ([]mysqlProductsToJsonConvert, error) {
	var mysqlProductsToConvert = []mysqlProductsToJsonConvert{}

	for _, image := range images {
		imageJasn := mysqlProductsToJsonConvert{}
		if err := json.Unmarshal([]byte(image.String()), &imageJasn); err != nil {
			return nil, err
		}
		mysqlProductsToConvert = append(mysqlProductsToConvert, imageJasn)
	}

	return mysqlProductsToConvert, nil
}

func flattenMysqlProduct(list []*vmysql.CloudDbProduct) []*mysqlProductModel {
	var outputs []*mysqlProductModel

	for _, v := range list {
		var output mysqlProductModel
		output.refreshFromOutput(v)

		outputs = append(outputs, &output)
	}
	return outputs
}

func (m *mysqlProductList) refreshFromOutput(ctx context.Context, list []*mysqlProductModel) {
	productListValue, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: mysqlProductModel{}.attrTypes()}, list)
	m.ProductList = productListValue
	m.ID = types.StringValue(time.Now().UTC().String())
}

type mysqlProductList struct {
	ID                         types.String `tfsdk:"id"`
	CloudMysqlImageProductCode types.String `tfsdk:"image_product_code"`
	ProductList                types.List   `tfsdk:"product_list"`
	OutputFile                 types.String `tfsdk:"output_file"`
	Filters                    types.Set    `tfsdk:"filter"`
}

type mysqlProductModel struct {
	ProductCode        types.String `tfsdk:"product_code"`
	ProductName        types.String `tfsdk:"product_name"`
	ProductType        types.String `tfsdk:"product_type"`
	ProductDescription types.String `tfsdk:"product_description"`
	InfraResourceType  types.String `tfsdk:"infra_resource_type"`
	CpuCount           types.Int64  `tfsdk:"cpu_count"`
	MemorySize         types.Int64  `tfsdk:"memory_size"`
	DiskType           types.String `tfsdk:"disk_type"`
}
type mysqlProductsToJsonConvert struct {
	ProductCode        string `json:"product_code"`
	ProductName        string `json:"product_name"`
	ProductType        string `json:"product_type"`
	ProductDescription string `json:"product_description"`
	InfraResourceType  string `json:"infra_resource_type"`
	CpuCount           int    `json:"cpu_count"`
	MemorySize         int    `json:"memory_size"`
	DiskType           string `json:"disk_type"`
}

func (m mysqlProductModel) attrTypes() map[string]attr.Type {
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
func (m *mysqlProductModel) refreshFromOutput(output *vmysql.CloudDbProduct) {
	m.ProductCode = types.StringPointerValue(output.ProductCode)
	m.ProductName = types.StringPointerValue(output.ProductName)
	m.ProductType = types.StringPointerValue(output.ProductType.Code)
	m.ProductDescription = types.StringPointerValue(output.ProductDescription)
	m.InfraResourceType = types.StringPointerValue(output.InfraResourceType.Code)
	m.CpuCount = common.Int64ValueFromInt32(output.CpuCount)
	m.MemorySize = types.Int64PointerValue(output.MemorySize)
	m.DiskType = types.StringPointerValue(output.DiskType.Code)
}
