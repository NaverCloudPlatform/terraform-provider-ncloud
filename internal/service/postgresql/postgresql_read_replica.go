package postgresql

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpostgresql"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"
)

var (
	_ resource.Resource                = &postgresqlReadReplicaResource{}
	_ resource.ResourceWithConfigure   = &postgresqlReadReplicaResource{}
	_ resource.ResourceWithImportState = &postgresqlReadReplicaResource{}
)

func NewPostgresqlReadReplicaResource() resource.Resource {
	return &postgresqlReadReplicaResource{}
}

type postgresqlReadReplicaResource struct {
	config *conn.ProviderConfig
}

func (r *postgresqlReadReplicaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: postgresql_instance_no:id Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("postgresql_instance_no"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

func (r *postgresqlReadReplicaResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*conn.ProviderConfig)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ProviderConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.config = config
}

func (r *postgresqlReadReplicaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgresql_read_replica"
}

func (r *postgresqlReadReplicaResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttribute(),
			"postgresql_instance_no": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// Available only `pub` and `fin` site. But GOV response message have both values.
			"subnet_no": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"postgresql_server_list": schema.ListNestedAttribute{
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
						"product_code": schema.StringAttribute{
							Computed: true,
						},
						"zone_code": schema.StringAttribute{
							Computed: true,
						},
						"subnet_no": schema.StringAttribute{
							Computed: true,
						},
						"public_subnet": schema.BoolAttribute{
							Computed: true,
						},
						"public_domain": schema.StringAttribute{
							Computed: true,
						},
						"private_domain": schema.StringAttribute{
							Computed: true,
						},
						"private_ip": schema.StringAttribute{
							Computed: true,
						},
						"data_storage_size": schema.Int64Attribute{
							Computed: true,
						},
						"used_data_storage_size": schema.Int64Attribute{
							Computed: true,
						},
						"cpu_count": schema.Int64Attribute{
							Computed: true,
						},
						"memory_size": schema.Int64Attribute{
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

func (r *postgresqlReadReplicaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan postgresqlReadReplicaResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vpostgresql.CreateCloudPostgresqlReadReplicaInstanceRequest{
		CloudPostgresqlInstanceNo: plan.PostgresqlInstanceNo.ValueStringPointer(),
	}

	if !plan.SubnetNo.IsNull() && !plan.SubnetNo.IsUnknown() {
		if r.config.Site == "gov" {
			resp.Diagnostics.AddError(
				"NOT SUPPORT GOV SITE",
				"`subnet_no` does not support gov site",
			)
			return
		}
		reqParams.SubnetNo = plan.SubnetNo.ValueStringPointer()
	}

	tflog.Info(ctx, "CreateCloudPostgresqlReadReplica reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vpostgresql.V2Api.CreateCloudPostgresqlReadReplicaInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "CreateCloudPostgresqlReadReplica response="+common.MarshalUncheckedString(response))

	if response == nil || len(response.CloudPostgresqlInstanceList) < 1 {
		resp.Diagnostics.AddError("CREATING ERROR", "response invalid")
		return
	}

	postgresqlIns := response.CloudPostgresqlInstanceList[0]
	serverList := postgresqlIns.CloudPostgresqlServerInstanceList
	var index int

	for i, server := range serverList {
		if (*server.CloudPostgresqlServerRole.Code == "S") && (*server.CloudPostgresqlServerInstanceStatusName == CREATING) {
			index = i
			break
		}
	}

	output, err := waitPostgresqlReadReplicaCreation(ctx, r.config, *postgresqlIns.CloudPostgresqlInstanceNo, index)
	if err != nil {
		resp.Diagnostics.AddError("WAITING FOR CREATION ERROR", err.Error())
		return
	}

	if diags := plan.refreshFromOutput(ctx, output, postgresqlIns.CloudPostgresqlInstanceNo); diags.HasError() {
		resp.Diagnostics.AddError("CREATING ERROR", "refreshFromOutput error")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *postgresqlReadReplicaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state postgresqlReadReplicaResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetPostgresqlReadReplicaServer(ctx, r.config, state.PostgresqlInstanceNo.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if diags := state.refreshFromOutput(ctx, output, state.PostgresqlInstanceNo.ValueStringPointer()); diags.HasError() {
		resp.Diagnostics.AddError("READING ERROR", "refreshFromOutput error")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *postgresqlReadReplicaResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *postgresqlReadReplicaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state postgresqlReadReplicaResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vpostgresql.DeleteCloudPostgresqlReadReplicaInstanceRequest{
		RegionCode:                      &r.config.RegionCode,
		CloudPostgresqlServerInstanceNo: state.ID.ValueStringPointer(),
	}
	tflog.Info(ctx, "DeletePostgresqlReadReplica reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vpostgresql.V2Api.DeleteCloudPostgresqlReadReplicaInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("DELETING ERROR", err.Error())
		return
	}

	tflog.Info(ctx, "DeletePostgresqlReadReplica response="+common.MarshalUncheckedString(response))

	if err := waitPostgresqlReadReplicaDeletion(ctx, r.config, state.PostgresqlInstanceNo.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("WAITING FOR DELETE ERROR", err.Error())
	}
}

func waitPostgresqlReadReplicaCreation(ctx context.Context, config *conn.ProviderConfig, instanceNo string, index int) ([]*vpostgresql.CloudPostgresqlServerInstance, error) {
	var serverInstance []*vpostgresql.CloudPostgresqlServerInstance
	stateConf := &retry.StateChangeConf{
		Pending: []string{CREATING, SETTING},
		Target:  []string{RUNNING},
		Refresh: func() (interface{}, string, error) {
			instance, err := findPostgresqlReadReplicaServer(ctx, config, instanceNo, index)
			serverInstance = instance
			if err != nil {
				return 0, "", err
			}

			if instance == nil {
				return 0, "", fmt.Errorf("Instance is nil")
			}

			status := instance[0].CloudPostgresqlServerInstanceStatusName
			if *status == CREATING || *status == SETTING || *status == RUNNING {
				return instance, *status, nil
			}

			return 0, "", fmt.Errorf("error occurred while waiting to create postgresql read replica")
		},
		Timeout:    6 * conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error waiting for postgresql read replica state to be \"CREAT\": %s", err)
	}

	return serverInstance, nil
}

func waitPostgresqlReadReplicaDeletion(ctx context.Context, config *conn.ProviderConfig, instanceNo string, serverInstanceNo string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{DELETING},
		Target:  []string{DELETED},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetPostgresqlReadReplicaServer(ctx, config, instanceNo, serverInstanceNo)
			if err != nil {
				return 0, "", err
			}

			if instance == nil {
				return instance, DELETED, nil
			}

			statusName := instance[0].CloudPostgresqlServerInstanceStatusName

			if *statusName == DELETING {
				return instance, DELETING, nil
			}

			if *statusName == DELETED {
				return instance, DELETED, nil
			}

			return 0, "", fmt.Errorf("error occurred while waiting to delete postgresql read replica")
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      1 * time.Minute,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for postgresql read replica (%s) to become termintaing: %s", serverInstanceNo, err)
	}

	return nil
}

func findPostgresqlReadReplicaServer(ctx context.Context, config *conn.ProviderConfig, instanceNo string, index int) ([]*vpostgresql.CloudPostgresqlServerInstance, error) {
	reqParams := &vpostgresql.GetCloudPostgresqlInstanceDetailRequest{
		RegionCode:                &config.RegionCode,
		CloudPostgresqlInstanceNo: ncloud.String(instanceNo),
	}
	tflog.Info(ctx, "GetPostgresqlDetail reqParams="+common.MarshalUncheckedString(reqParams))

	resp, err := config.Client.Vpostgresql.V2Api.GetCloudPostgresqlInstanceDetail(reqParams)
	if err != nil && !(strings.Contains(err.Error(), `"returnCode": "5001017"`)) {
		return nil, err
	}
	tflog.Info(ctx, "GetPostgresqlDetail response="+common.MarshalUncheckedString(resp))

	if resp == nil || len(resp.CloudPostgresqlInstanceList) < 1 {
		return nil, fmt.Errorf("response is nil")
	}

	serverList := resp.CloudPostgresqlInstanceList[0].CloudPostgresqlServerInstanceList

	for i, server := range serverList {
		if i == index {
			return []*vpostgresql.CloudPostgresqlServerInstance{server}, nil
		}
	}
	return nil, nil
}

func GetPostgresqlReadReplicaServer(ctx context.Context, config *conn.ProviderConfig, instanceNo string, serverInstanceNo string) ([]*vpostgresql.CloudPostgresqlServerInstance, error) {
	reqParams := &vpostgresql.GetCloudPostgresqlInstanceDetailRequest{
		RegionCode:                &config.RegionCode,
		CloudPostgresqlInstanceNo: ncloud.String(instanceNo),
	}
	tflog.Info(ctx, "GetPostgresqlDetail reqParams="+common.MarshalUncheckedString(reqParams))

	resp, err := config.Client.Vpostgresql.V2Api.GetCloudPostgresqlInstanceDetail(reqParams)
	if err != nil && !(strings.Contains(err.Error(), `"returnCode": "5001017"`)) {
		return nil, err
	}
	tflog.Info(ctx, "GetPostgresqlDetail response="+common.MarshalUncheckedString(resp))

	if resp == nil || len(resp.CloudPostgresqlInstanceList) < 1 {
		return nil, nil
	}

	serverList := resp.CloudPostgresqlInstanceList[0].CloudPostgresqlServerInstanceList

	for _, server := range serverList {
		if (*server.CloudPostgresqlServerRole.Code == "S") && (*server.CloudPostgresqlServerInstanceNo == serverInstanceNo) {
			return []*vpostgresql.CloudPostgresqlServerInstance{server}, nil
		}
	}
	return nil, nil
}

type postgresqlReadReplicaResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	PostgresqlInstanceNo types.String `tfsdk:"postgresql_instance_no"`
	SubnetNo             types.String `tfsdk:"subnet_no"`
	PostgresqlServerList types.List   `tfsdk:"postgresql_server_list"`
}

func (r *postgresqlReadReplicaResourceModel) refreshFromOutput(ctx context.Context, output []*vpostgresql.CloudPostgresqlServerInstance, instanceNo *string) diag.Diagnostics {
	r.ID = types.StringPointerValue(output[0].CloudPostgresqlServerInstanceNo)
	r.PostgresqlInstanceNo = types.StringPointerValue(instanceNo)
	r.SubnetNo = types.StringPointerValue(output[0].SubnetNo)

	serverList, diags := listValueFromPostgresqlServerInatanceList(ctx, output)
	r.PostgresqlServerList = serverList

	return diags
}
