package hadoop

import (
	"context"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vhadoop"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
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

func (h *hadoopDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("cluster_name"),
					),
				},
			},
			"cluster_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("id"),
					),
				},
			},
			"region_code": schema.StringAttribute{
				Computed: true,
			},
			"vpc_no": schema.StringAttribute{
				Optional: true,
			},
			"image_product_code": schema.StringAttribute{
				Computed: true,
			},
			"cluster_type_code": schema.StringAttribute{
				Computed: true,
			},
			"version": schema.StringAttribute{
				Computed: true,
			},
			"ambari_server_host": schema.StringAttribute{
				Computed: true,
			},
			"cluster_direct_access_account": schema.StringAttribute{
				Computed: true,
			},
			"login_key": schema.StringAttribute{
				Computed: true,
			},
			"object_storage_bucket": schema.StringAttribute{
				Computed: true,
			},
			"kdc_realm": schema.StringAttribute{
				Computed: true,
			},
			"domain": schema.StringAttribute{
				Computed: true,
			},
			"is_ha": schema.BoolAttribute{
				Computed: true,
			},
			"add_on_code_list": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"access_control_group_no_list": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"hadoop_server_list": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"server_instance_no": schema.StringAttribute{
							Computed: true,
						},
						"server_name": schema.StringAttribute{
							Computed: true,
						},
						"server_role": schema.StringAttribute{
							Computed: true,
						},
						"zone_code": schema.StringAttribute{
							Computed: true,
						},
						"subnet_no": schema.StringAttribute{
							Computed: true,
						},
						"product_code": schema.StringAttribute{
							Computed: true,
						},
						"is_public_subnet": schema.BoolAttribute{
							Computed: true,
						},
						"cpu_count": schema.Int64Attribute{
							Computed: true,
						},
						"memory_size": schema.Int64Attribute{
							Computed: true,
						},
						"data_storage_type": schema.StringAttribute{
							Computed: true,
						},
						"data_storage_size": schema.Int64Attribute{
							Computed: true,
						},
						"uptime": schema.StringAttribute{
							Computed: true,
						},
						"create_date": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
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

func (h *hadoopDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data hadoopDataSourceModel
	var hadoopId string

	if !h.config.SupportVPC {
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

	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		hadoopId = data.ID.ValueString()
	}

	if !data.ClusterName.IsNull() && !data.ClusterName.IsUnknown() {
		reqParams := &vhadoop.GetCloudHadoopInstanceListRequest{
			RegionCode:             &h.config.RegionCode,
			CloudHadoopClusterName: data.ClusterName.ValueStringPointer(),
		}
		tflog.Info(ctx, "GetHadoopList reqParams="+common.MarshalUncheckedString(reqParams))

		listResp, err := h.config.Client.Vhadoop.V2Api.GetCloudHadoopInstanceList(reqParams)
		if err != nil {
			resp.Diagnostics.AddError("READING ERROR", err.Error())
			return
		}
		tflog.Info(ctx, "GetHadoopList response="+common.MarshalUncheckedString(listResp))

		if listResp == nil || len(listResp.CloudHadoopInstanceList) < 1 {
			resp.Diagnostics.AddError("READING ERROR", "no result. please change search criteria and try again.")
			return
		}
		hadoopId = *listResp.CloudHadoopInstanceList[0].CloudHadoopInstanceNo
	}

	output, err := GetHadoopInstance(ctx, h.config, hadoopId)
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.Diagnostics.AddError("READING ERROR", "no result. please change search criteria and try again.")
		return
	}

	data.refreshFromOutput(ctx, output)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

type hadoopDataSourceModel struct {
	ID                         types.String `tfsdk:"id"`
	ClusterName                types.String `tfsdk:"cluster_name"`
	RegionCode                 types.String `tfsdk:"region_code"`
	VpcNo                      types.String `tfsdk:"vpc_no"`
	ImageProductCode           types.String `tfsdk:"image_product_code"`
	ClusterTypeCode            types.String `tfsdk:"cluster_type_code"`
	Version                    types.String `tfsdk:"version"`
	AmbariServerHost           types.String `tfsdk:"ambari_server_host"`
	ClusterDirectAccessAccount types.String `tfsdk:"cluster_direct_access_account"`
	LoginKey                   types.String `tfsdk:"login_key"`
	ObjectStorageBucket        types.String `tfsdk:"object_storage_bucket"`
	KdcRealm                   types.String `tfsdk:"kdc_realm"`
	Domain                     types.String `tfsdk:"domain"`
	IsHa                       types.Bool   `tfsdk:"is_ha"`
	AddOnCodeList              types.List   `tfsdk:"add_on_code_list"`
	AccessControlGroupNoList   types.List   `tfsdk:"access_control_group_no_list"`
	HadoopServerList           types.List   `tfsdk:"hadoop_server_list"`
}

func (d *hadoopDataSourceModel) refreshFromOutput(ctx context.Context, output *vhadoop.CloudHadoopInstance) {
	d.ID = types.StringPointerValue(output.CloudHadoopInstanceNo)
	d.ClusterName = types.StringPointerValue(output.CloudHadoopClusterName)
	d.RegionCode = types.StringPointerValue(output.CloudHadoopServerInstanceList[0].RegionCode)
	d.VpcNo = types.StringPointerValue(output.CloudHadoopServerInstanceList[0].VpcNo)
	d.ImageProductCode = types.StringPointerValue(output.CloudHadoopImageProductCode)
	d.ClusterTypeCode = types.StringPointerValue(output.CloudHadoopClusterType.Code)
	d.Version = types.StringPointerValue(output.CloudHadoopVersion.Code)
	d.AmbariServerHost = types.StringPointerValue(output.AmbariServerHost)
	d.ClusterDirectAccessAccount = types.StringPointerValue(output.ClusterDirectAccessAccount)
	d.LoginKey = types.StringPointerValue(output.LoginKey)
	d.ObjectStorageBucket = types.StringPointerValue(output.ObjectStorageBucket)
	d.KdcRealm = types.StringPointerValue(output.KdcRealm)
	d.Domain = types.StringPointerValue(output.Domain)
	d.IsHa = types.BoolPointerValue(output.IsHa)

	var addOnList []string
	for _, addOn := range output.CloudHadoopAddOnList {
		addOnList = append(addOnList, *addOn.Code)
	}
	d.AddOnCodeList, _ = types.ListValueFrom(ctx, types.StringType, addOnList)
	d.AccessControlGroupNoList, _ = types.ListValueFrom(ctx, types.StringType, output.AccessControlGroupNoList)
	d.HadoopServerList, _ = listValueFromHadoopServerInatanceList(ctx, output.CloudHadoopServerInstanceList)
}
