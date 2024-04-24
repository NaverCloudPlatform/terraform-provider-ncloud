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

func (d *hadoopDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hadoop"
}

func (d *hadoopDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
				Computed: true,
			},
			"edge_node_subnet_no": schema.StringAttribute{
				Computed: true,
			},
			"master_node_subnet_no": schema.StringAttribute{
				Computed: true,
			},
			"worker_node_subnet_no": schema.StringAttribute{
				Computed: true,
			},
			"master_node_data_storage_type": schema.StringAttribute{
				Computed: true,
			},
			"worker_node_data_storage_type": schema.StringAttribute{
				Computed: true,
			},
			"master_node_data_storage_size": schema.Int64Attribute{
				Computed: true,
			},
			"worker_node_data_storage_size": schema.Int64Attribute{
				Computed: true,
			},
			"image_product_code": schema.StringAttribute{
				Computed: true,
			},
			"edge_node_product_code": schema.StringAttribute{
				Computed: true,
			},
			"master_node_product_code": schema.StringAttribute{
				Computed: true,
			},
			"worker_node_product_code": schema.StringAttribute{
				Computed: true,
			},
			"worker_node_count": schema.Int64Attribute{
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
			"login_key_name": schema.StringAttribute{
				Computed: true,
			},
			"bucket_name": schema.StringAttribute{
				Computed: true,
			},
			"use_kdc": schema.BoolAttribute{
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

func (d *hadoopDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *hadoopDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data hadoopDataSourceModel
	var hadoopId string

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

	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		hadoopId = data.ID.ValueString()
	}

	if !data.ClusterName.IsNull() && !data.ClusterName.IsUnknown() {
		reqParams := &vhadoop.GetCloudHadoopInstanceListRequest{
			RegionCode:             &d.config.RegionCode,
			CloudHadoopClusterName: data.ClusterName.ValueStringPointer(),
		}
		tflog.Info(ctx, "GetHadoopList reqParams="+common.MarshalUncheckedString(reqParams))

		listResp, err := d.config.Client.Vhadoop.V2Api.GetCloudHadoopInstanceList(reqParams)
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

	output, err := GetHadoopInstance(ctx, d.config, hadoopId)
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
	EdgeNodeSubnetNo           types.String `tfsdk:"edge_node_subnet_no"`
	MasterNodeSubnetNo         types.String `tfsdk:"master_node_subnet_no"`
	WorkerNodeSubnetNo         types.String `tfsdk:"worker_node_subnet_no"`
	MasterNodeDataStorageType  types.String `tfsdk:"master_node_data_storage_type"`
	WorkerNodeDataStorageType  types.String `tfsdk:"worker_node_data_storage_type"`
	MasterNodeDataStorageSize  types.Int64  `tfsdk:"master_node_data_storage_size"`
	WorkerNodeDataStorageSize  types.Int64  `tfsdk:"worker_node_data_storage_size"`
	ImageProductCode           types.String `tfsdk:"image_product_code"`
	EdgeNodeProductCode        types.String `tfsdk:"edge_node_product_code"`
	MasterNodeProductCode      types.String `tfsdk:"master_node_product_code"`
	WorkerNodeProductCode      types.String `tfsdk:"worker_node_product_code"`
	WorkerNodeCount            types.Int64  `tfsdk:"worker_node_count"`
	ClusterTypeCode            types.String `tfsdk:"cluster_type_code"`
	Version                    types.String `tfsdk:"version"`
	AmbariServerHost           types.String `tfsdk:"ambari_server_host"`
	ClusterDirectAccessAccount types.String `tfsdk:"cluster_direct_access_account"`
	LoginKey                   types.String `tfsdk:"login_key_name"`
	BucketName                 types.String `tfsdk:"bucket_name"`
	UseKdc                     types.Bool   `tfsdk:"use_kdc"`
	KdcRealm                   types.String `tfsdk:"kdc_realm"`
	Domain                     types.String `tfsdk:"domain"`
	IsHa                       types.Bool   `tfsdk:"is_ha"`
	AddOnCodeList              types.List   `tfsdk:"add_on_code_list"`
	AccessControlGroupNoList   types.List   `tfsdk:"access_control_group_no_list"`
	HadoopServerList           types.List   `tfsdk:"hadoop_server_list"`
}

func (m *hadoopDataSourceModel) refreshFromOutput(ctx context.Context, output *vhadoop.CloudHadoopInstance) {
	m.ID = types.StringPointerValue(output.CloudHadoopInstanceNo)
	m.ClusterName = types.StringPointerValue(output.CloudHadoopClusterName)
	m.RegionCode = types.StringPointerValue(output.CloudHadoopServerInstanceList[0].RegionCode)
	m.VpcNo = types.StringPointerValue(output.CloudHadoopServerInstanceList[0].VpcNo)
	m.ImageProductCode = types.StringPointerValue(output.CloudHadoopImageProductCode)
	m.ClusterTypeCode = types.StringPointerValue(output.CloudHadoopClusterType.Code)
	m.Version = types.StringPointerValue(output.CloudHadoopVersion.Code)
	m.AmbariServerHost = types.StringPointerValue(output.AmbariServerHost)
	m.ClusterDirectAccessAccount = types.StringPointerValue(output.ClusterDirectAccessAccount)
	m.LoginKey = types.StringPointerValue(output.LoginKey)
	m.BucketName = types.StringPointerValue(output.ObjectStorageBucket)
	m.KdcRealm = types.StringPointerValue(output.KdcRealm)
	m.Domain = types.StringPointerValue(output.Domain)
	m.IsHa = types.BoolPointerValue(output.IsHa)

	if output.KdcRealm != nil {
		m.UseKdc = types.BoolValue(true)
	} else {
		m.UseKdc = types.BoolValue(false)
	}

	var count int64
	var storageSize int64
	for _, server := range output.CloudHadoopServerInstanceList {
		if server.CloudHadoopServerRole != nil {
			if *server.CloudHadoopServerRole.Code == "E" {
				m.EdgeNodeProductCode = types.StringPointerValue(server.CloudHadoopProductCode)
				m.EdgeNodeSubnetNo = types.StringPointerValue(server.SubnetNo)
			}
			if *server.CloudHadoopServerRole.Code == "M" {
				m.MasterNodeProductCode = types.StringPointerValue(server.CloudHadoopProductCode)
				m.MasterNodeSubnetNo = types.StringPointerValue(server.SubnetNo)
				if server.DataStorageType != nil {
					m.MasterNodeDataStorageType = types.StringPointerValue(server.DataStorageType.Code)
				}
				// Byte to GBi
				storageSize = *server.DataStorageSize / 1024 / 1024 / 1024
				m.MasterNodeDataStorageSize = types.Int64Value(storageSize)
			}
			if *server.CloudHadoopServerRole.Code == "D" {
				m.WorkerNodeProductCode = types.StringPointerValue(server.CloudHadoopProductCode)
				m.WorkerNodeSubnetNo = types.StringPointerValue(server.SubnetNo)
				if server.DataStorageType != nil {
					m.WorkerNodeDataStorageType = types.StringPointerValue(server.DataStorageType.Code)
				}
				// Byte to GBi
				storageSize = *server.DataStorageSize / 1024 / 1024 / 1024
				m.WorkerNodeDataStorageSize = types.Int64Value(storageSize)
				count++
			}
		}
	}
	m.WorkerNodeCount = types.Int64Value(count)

	var addOnList []string
	for _, addOn := range output.CloudHadoopAddOnList {
		addOnList = append(addOnList, *addOn.Code)
	}
	m.AddOnCodeList, _ = types.ListValueFrom(ctx, types.StringType, addOnList)
	m.AccessControlGroupNoList, _ = types.ListValueFrom(ctx, types.StringType, output.AccessControlGroupNoList)
	m.HadoopServerList, _ = listValueFromHadoopServerInatanceList(ctx, output.CloudHadoopServerInstanceList)
}
