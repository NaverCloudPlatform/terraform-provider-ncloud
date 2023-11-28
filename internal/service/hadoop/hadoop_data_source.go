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
	_ datasource.DataSource              = &hadoopDataSource{}
	_ datasource.DataSourceWithConfigure = &hadoopDataSource{}
)

func NewHadoopDataSource() datasource.DataSource {
	return &hadoopDataSource{}
}

type hadoopDataSource struct {
	config *conn.ProviderConfig
}

func (h *hadoopDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hadoop"
}

func (h *hadoopDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (h *hadoopDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"zone_code": schema.StringAttribute{
				Optional: true,
			},
			"vpc_no": schema.StringAttribute{
				Optional: true,
			},
			"subnet_no": schema.StringAttribute{
				Optional: true,
			},
			"cluster_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"server_name": schema.StringAttribute{
				Optional: true,
			},
			"server_instance_no": schema.StringAttribute{
				Optional: true,
			},
			"cluster_type_code": schema.StringAttribute{
				Computed: true,
			},
			"version": schema.StringAttribute{
				Computed: true,
			},
			"image_product_code": schema.StringAttribute{
				Computed: true,
			},
			"hadoop_server_instance_list": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"hadoop_server_name": schema.StringAttribute{
							Computed: true,
						},
						"hadoop_server_role": schema.StringAttribute{
							Computed: true,
						},
						"hadoop_product_code": schema.StringAttribute{
							Computed: true,
						},
						"region_code": schema.StringAttribute{
							Computed: true,
						},
						"zone_code": schema.StringAttribute{
							Computed: true,
						},
						"vpc_no": schema.StringAttribute{
							Computed: true,
						},
						"subnet_no": schema.StringAttribute{
							Computed: true,
						},
						"is_public_subnet": schema.BoolAttribute{
							Computed: true,
						},
						"data_storage_size": schema.Int64Attribute{
							Computed: true,
						},
						"cpu_count": schema.Int64Attribute{
							Computed: true,
						},
						"memory_size": schema.Int64Attribute{
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

func (h *hadoopDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data hadoopDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vhadoop.GetCloudHadoopInstanceListRequest{
		RegionCode: &h.config.RegionCode,
	}

	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		reqParams.CloudHadoopInstanceNoList = []*string{data.ID.ValueStringPointer()}
	}
	if !data.ZoneCode.IsNull() && !data.ZoneCode.IsUnknown() {
		reqParams.ZoneCode = data.ZoneCode.ValueStringPointer()
	}
	if !data.VpcNo.IsNull() && !data.VpcNo.IsUnknown() {
		reqParams.VpcNo = data.VpcNo.ValueStringPointer()
	}
	if !data.SubnetNo.IsNull() && !data.SubnetNo.IsUnknown() {
		reqParams.SubnetNo = data.SubnetNo.ValueStringPointer()
	}
	if !data.ClusterName.IsNull() && !data.ClusterName.IsUnknown() {
		reqParams.CloudHadoopClusterName = data.ClusterName.ValueStringPointer()
	}
	if !data.ServerName.IsNull() && !data.ServerName.IsUnknown() {
		reqParams.CloudHadoopServerName = data.ServerName.ValueStringPointer()
	}
	if !data.ServerInstanceNo.IsNull() && !data.ServerInstanceNo.IsUnknown() {
		reqParams.CloudHadoopServerInstanceNoList = []*string{data.ServerInstanceNo.ValueStringPointer()}
	}

	tflog.Info(ctx, "GetHadoopList", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})

	hadoopResp, err := h.config.Client.Vhadoop.V2Api.GetCloudHadoopInstanceList(reqParams)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetCloudHadoopInstanceList",
			fmt.Sprintf("error: %s, reqParams: %s", err.Error(), common.MarshalUncheckedString(reqParams)),
		)
		resp.Diagnostics.Append(diags...)
		return
	}

	tflog.Info(ctx, "GetHadoopList response", map[string]any{
		"hadoopResponse": common.MarshalUncheckedString(hadoopResp),
	})

	hadoopList := flattenHadoopList(ctx, hadoopResp.CloudHadoopInstanceList)

	filteredList := common.FilterModels(ctx, data.Filters, hadoopList)

	if err := verify.ValidateOneResult(len(filteredList)); err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"GetHadoopList result vaildation",
			err.Error(),
		)
		resp.Diagnostics.Append(diags...)
		return
	}

	state := filteredList[0]
	state.Filters = data.Filters
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenHadoopList(ctx context.Context, hadoops []*vhadoop.CloudHadoopInstance) []*hadoopDataSourceModel {
	var outputs []*hadoopDataSourceModel

	for _, v := range hadoops {
		var output hadoopDataSourceModel
		output.refreshFromOutput(ctx, v)
		outputs = append(outputs, &output)
	}

	return outputs
}

type hadoopDataSourceModel struct {
	ID                       types.String `tfsdk:"id"`
	ZoneCode                 types.String `tfsdk:"zone_code"`
	VpcNo                    types.String `tfsdk:"vpc_no"`
	SubnetNo                 types.String `tfsdk:"subnet_no"`
	ClusterName              types.String `tfsdk:"cluster_name"`
	ServerName               types.String `tfsdk:"server_name"`
	ServerInstanceNo         types.String `tfsdk:"server_instance_no"`
	ClusterTypeCode          types.String `tfsdk:"cluster_type_code"`
	Version                  types.String `tfsdk:"version"`
	ImageProductCode         types.String `tfsdk:"image_product_code"`
	HadoopServerInstanceList types.List   `tfsdk:"hadoop_server_instance_list"`
	Filters                  types.Set    `tfsdk:"filter"`
}

func (m *hadoopDataSourceModel) refreshFromOutput(ctx context.Context, output *vhadoop.CloudHadoopInstance) {
	m.ID = types.StringPointerValue(output.CloudHadoopInstanceNo)
	m.ClusterTypeCode = types.StringPointerValue(output.CloudHadoopClusterType.Code)
	m.Version = types.StringPointerValue(output.CloudHadoopVersion.Code)
	m.ImageProductCode = types.StringPointerValue(output.CloudHadoopImageProductCode)
	m.ClusterName = types.StringPointerValue(output.CloudHadoopClusterName)
	m.Version = types.StringPointerValue(output.CloudHadoopVersion.Code)
	m.ImageProductCode = types.StringPointerValue(output.CloudHadoopImageProductCode)
	m.HadoopServerInstanceList, _ = listValueFromHadoopServerInatanceList(ctx, output.CloudHadoopServerInstanceList)
}
