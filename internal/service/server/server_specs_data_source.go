package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

var (
	_ datasource.DataSource              = &serverSpecsDataSource{}
	_ datasource.DataSourceWithConfigure = &serverSpecsDataSource{}
)

func NewServerSpecsDataSource() datasource.DataSource {
	return &serverSpecsDataSource{}
}

type serverSpecsDataSource struct {
	config *conn.ProviderConfig
}

func (d *serverSpecsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_specs"
}

func (d *serverSpecsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"output_file": schema.StringAttribute{
				Optional: true,
			},
			"server_spec_list": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"server_spec_code": schema.StringAttribute{
							Computed: true,
						},
						"hypervisor_type": schema.StringAttribute{
							Computed: true,
						},
						"generation_code": schema.StringAttribute{
							Computed: true,
						},
						"cpu_architecture_type": schema.StringAttribute{
							Computed: true,
						},
						"cpu_count": schema.Int32Attribute{
							Computed: true,
						},
						"memory_size": schema.Int64Attribute{
							Computed: true,
						},
						"block_storage_max_count": schema.Int32Attribute{
							Computed: true,
						},
						"block_storage_max_iops": schema.Int32Attribute{
							Computed: true,
						},
						"block_storage_max_throughput": schema.Int32Attribute{
							Computed: true,
						},
						"network_performance": schema.Int64Attribute{
							Computed: true,
						},
						"network_interface_max_count": schema.Int32Attribute{
							Computed: true,
						},
						"gpu_count": schema.Int32Attribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"product_code": schema.StringAttribute{
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

func (d *serverSpecsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.config = config
}

func (d *serverSpecsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data serverSpecsDataSourceModel

	if !d.config.SupportVPC {
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

	reqParams := &vserver.GetServerSpecListRequest{
		RegionCode: &d.config.RegionCode,
	}
	tflog.Info(ctx, "GetServerSpecListRequest reqParams="+common.MarshalUncheckedString(reqParams))

	specResp, err := d.config.Client.Vserver.V2Api.GetServerSpecList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "GetServerSpecListRequest response="+common.MarshalUncheckedString(specResp))

	if specResp == nil || len(specResp.ServerSpecList) < 1 {
		resp.Diagnostics.AddError("READING ERROR", "no result.")
		return
	}

	specList := flattenServerSpecList(specResp.ServerSpecList)
	fillteredList := common.FilterModels(ctx, data.Filters, specList)
	diags := data.refreshFromOutput(ctx, fillteredList)
	if diags.HasError() {
		resp.Diagnostics.AddError("READING ERROR", "refreshFromOutput error")
		return
	}

	if !data.OutputFile.IsNull() && data.OutputFile.String() != "" {
		outputPath := data.OutputFile.ValueString()

		if convertedList, err := convertSpecToJsonStruct(data.ServerSpecList.Elements()); err != nil {
			resp.Diagnostics.AddError("OUTPUT FILE ERROR", err.Error())
			return
		} else if err := common.WriteToFile(outputPath, convertedList); err != nil {
			resp.Diagnostics.AddError("OUTPUT FILE ERROR", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func convertSpecToJsonStruct(specs []attr.Value) ([]serverSpecToJsonConvert, error) {
	var serverSpecToConvert = []serverSpecToJsonConvert{}

	for _, spec := range specs {
		specJson := serverSpecToJsonConvert{}
		if err := json.Unmarshal([]byte(common.ReplaceNull(spec.String())), &specJson); err != nil {
			return nil, err
		}
		serverSpecToConvert = append(serverSpecToConvert, specJson)
	}

	return serverSpecToConvert, nil
}

func flattenServerSpecList(list []*vserver.ServerSpec) []*serverSpec {
	var outputs []*serverSpec

	for _, v := range list {
		var output serverSpec
		output.refreshFromOutput(v)

		outputs = append(outputs, &output)
	}
	return outputs
}

type serverSpecsDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	ServerSpecList types.List   `tfsdk:"server_spec_list"`
	OutputFile     types.String `tfsdk:"output_file"`
	Filters        types.Set    `tfsdk:"filter"`
}

type serverSpec struct {
	ServerSpecCode            types.String `tfsdk:"server_spec_code"`
	HypervisorType            types.String `tfsdk:"hypervisor_type"`
	GenerationCode            types.String `tfsdk:"generation_code"`
	CpuArchitectureType       types.String `tfsdk:"cpu_architecture_type"`
	CpuCount                  types.Int32  `tfsdk:"cpu_count"`
	MemorySize                types.Int64  `tfsdk:"memory_size"`
	BlockStorageMaxCount      types.Int32  `tfsdk:"block_storage_max_count"`
	BlockStorageMaxIops       types.Int32  `tfsdk:"block_storage_max_iops"`
	BlockStorageMaxThroughput types.Int32  `tfsdk:"block_storage_max_throughput"`
	NetworkPerformance        types.Int64  `tfsdk:"network_performance"`
	NetworkInterfaceMaxCount  types.Int32  `tfsdk:"network_interface_max_count"`
	GpuCount                  types.Int32  `tfsdk:"gpu_count"`
	Description               types.String `tfsdk:"description"`
	ProductCode               types.String `tfsdk:"product_code"`
}

type serverSpecToJsonConvert struct {
	ServerSpecCode            string `json:"server_spec_code"`
	HypervisorType            string `json:"hypervisor_type"`
	GenerationCode            string `json:"generation_code"`
	CpuArchitectureType       string `json:"cpu_architecture_type"`
	CpuCount                  int    `json:"cpu_count,omitempty"`
	MemorySize                int64  `json:"memory_size,omitempty"`
	BlockStorageMaxCount      int    `json:"block_storage_max_count,omitempty"`
	BlockStorageMaxIops       int    `json:"block_storage_max_iops,omitempty"`
	BlockStorageMaxThroughput int    `json:"block_storage_max_throughput,omitempty"`
	NetworkPerformance        int64  `json:"network_performance,omitempty"`
	NetworkInterfaceMaxCount  int    `json:"network_interface_max_count,omitempty"`
	GpuCount                  int    `json:"gpu_count,omitempty"`
	Description               string `json:"description"`
	ProductCode               string `json:"product_code"`
}

func (d serverSpec) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"server_spec_code":             types.StringType,
		"hypervisor_type":              types.StringType,
		"generation_code":              types.StringType,
		"cpu_architecture_type":        types.StringType,
		"cpu_count":                    types.Int32Type,
		"memory_size":                  types.Int64Type,
		"block_storage_max_count":      types.Int32Type,
		"block_storage_max_iops":       types.Int32Type,
		"block_storage_max_throughput": types.Int32Type,
		"network_performance":          types.Int64Type,
		"network_interface_max_count":  types.Int32Type,
		"gpu_count":                    types.Int32Type,
		"description":                  types.StringType,
		"product_code":                 types.StringType,
	}
}

func (d *serverSpecsDataSourceModel) refreshFromOutput(ctx context.Context, list []*serverSpec) diag.Diagnostics {
	specListValue, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: serverSpec{}.attrTypes()}, list)
	if diags.HasError() {
		return diags
	}

	d.ServerSpecList = specListValue
	d.ID = types.StringValue("")

	return diags
}

func (d *serverSpec) refreshFromOutput(output *vserver.ServerSpec) {
	d.ServerSpecCode = types.StringPointerValue(output.ServerSpecCode)
	d.GenerationCode = types.StringPointerValue(output.GenerationCode)
	d.CpuArchitectureType = types.StringPointerValue(output.CpuArchitectureType.Code)
	d.CpuCount = types.Int32PointerValue(output.CpuCount)
	d.MemorySize = types.Int64PointerValue(output.MemorySize)
	d.BlockStorageMaxCount = types.Int32PointerValue(output.BlockStorageMaxCount)
	d.BlockStorageMaxIops = types.Int32PointerValue(output.BlockStorageMaxIops)
	d.BlockStorageMaxThroughput = types.Int32PointerValue(output.BlockStorageMaxThroughput)
	d.NetworkPerformance = types.Int64PointerValue(output.NetworkPerformance)
	d.NetworkInterfaceMaxCount = types.Int32PointerValue(output.NetworkInterfaceMaxCount)
	d.GpuCount = types.Int32PointerValue(output.GpuCount)
	d.Description = types.StringPointerValue(output.ServerSpecDescription)
	d.ProductCode = types.StringPointerValue(output.ServerProductCode)

	if output.HypervisorType != nil {
		d.HypervisorType = types.StringPointerValue(output.HypervisorType.Code)
	}
}
