package hadoop

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vhadoop"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"time"
)

var (
	_ datasource.DataSource              = &hadoopsDataSource{}
	_ datasource.DataSourceWithConfigure = &hadoopsDataSource{}
)

func NewHadoopsDataSource() datasource.DataSource {
	return &hadoopsDataSource{}
}

type hadoopsDataSource struct {
	config *conn.ProviderConfig
}

func (h *hadoopsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hadoops"
}

func (h *hadoopsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (h *hadoopsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			},
			"server_name": schema.StringAttribute{
				Optional: true,
			},
			"server_instance_no": schema.StringAttribute{
				Optional: true,
			},
		},
		Blocks: map[string]schema.Block{
			"filter": common.DataSourceFiltersBlock(),
			"hadoops": schema.ListNestedBlock{
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"cluster_name": schema.StringAttribute{
							Computed: true,
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
				},
			},
		},
	}
}

func (h *hadoopsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data hadoopsDataSourceModel
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
	diag := data.refreshFromOutput(ctx, filteredList)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := data
	state.ID = types.StringValue(time.Now().UTC().String())
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type hadoopsDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	ZoneCode         types.String `tfsdk:"zone_code"`
	VpcNo            types.String `tfsdk:"vpc_no"`
	SubnetNo         types.String `tfsdk:"subnet_no"`
	ClusterName      types.String `tfsdk:"cluster_name"`
	ServerName       types.String `tfsdk:"server_name"`
	ServerInstanceNo types.String `tfsdk:"server_instance_no"`
	Hadoops          types.List   `tfsdk:"hadoops"`
	Filters          types.Set    `tfsdk:"filter"`
}

type hadoop struct {
	ID                       types.String `tfsdk:"id"`
	ClusterName              types.String `tfsdk:"cluster_name"`
	ClusterTypeCode          types.String `tfsdk:"cluster_type_code"`
	Version                  types.String `tfsdk:"version"`
	ImageProductCode         types.String `tfsdk:"image_product_code"`
	HadoopServerInstanceList types.List   `tfsdk:"hadoop_server_instance_list"`
}

func (h hadoop) attrTypes() map[string]attr.Type {
	serverListElementType := types.ListType{
		ElemType: types.ObjectType{
			AttrTypes: hadoopServer{}.attrTypes(),
		},
	}

	return map[string]attr.Type{
		"id":                          types.StringType,
		"cluster_name":                types.StringType,
		"cluster_type_code":           types.StringType,
		"version":                     types.StringType,
		"image_product_code":          types.StringType,
		"hadoop_server_instance_list": serverListElementType,
	}
}

func (h *hadoopsDataSourceModel) refreshFromOutput(ctx context.Context, output []*hadoopDataSourceModel) diag.Diagnostics {
	var hadoops []hadoop
	var diag diag.Diagnostics

	for _, instance := range output {
		hadoops = append(hadoops, hadoop{
			ID:                       instance.ID,
			ClusterName:              instance.ClusterName,
			ClusterTypeCode:          instance.ClusterTypeCode,
			Version:                  instance.Version,
			ImageProductCode:         instance.ImageProductCode,
			HadoopServerInstanceList: instance.HadoopServerInstanceList,
		})
	}
	h.Hadoops, diag = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: hadoop{}.attrTypes()}, hadoops)
	if diag.HasError() {
		return diag
	}
	return nil
}
