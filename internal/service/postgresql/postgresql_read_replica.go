package postgresql

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpostgresql"
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
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
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
			"subnet_no": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *postgresqlReadReplicaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan postgresqlReadReplicaResourceModel

	if !r.config.SupportVPC {
		resp.Diagnostics.AddError(
			"NOT SUPPORT CLASSIC",
			"resource does not support CLASSIC. only VPC.",
		)
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vpostgresql.CreateCloudPostgresqlReadReplicaInstanceRequest{
		CloudPostgresqlInstanceNo: plan.PostgresqlInstanceNo.ValueStringPointer(),
	}

	if !plan.SubnetNo.IsNull() {
		reqParams.SubnetNo = plan.SubnetNo.ValueStringPointer()
	}

	tflog.Info(ctx, "CreateCloudPostgresqlReadReplica reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vpostgresql.V2Api.CreateCloudPostgresqlReadReplicaInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}

	if response == nil || len(response.CloudPostgresqlInstanceList) < 1 {
		resp.Diagnostics.AddError("CREATING ERROR", "response invalid")
		return
	}

	postgresqlIns := response.CloudPostgresqlInstanceList[0]
	serverList := postgresqlIns.CloudPostgresqlServerInstanceList

	if !*postgresqlIns.IsHa {
		resp.Diagnostics.AddError(
			"CREATIG ERROR",
			"when `is_ha` is false, Read Replica can't be created",
		)
	}

	if len(serverList) < 2 {
		resp.Diagnostics.AddError("CREATING ERROR", "response invalid")
		return
	}

	readReplicaServer, err := GetPostgresqlReadReplica(ctx, r.config, plan.PostgresqlInstanceNo.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("READ ERROR", err.Error())
		return
	}

	plan.ID = types.StringPointerValue(readReplicaServer.CloudPostgresqlServerInstanceNo)

	output, err := waitPostgresqlReadReplicaCreation(ctx, r.config, *postgresqlIns.CloudPostgresqlInstanceNo)
	if err != nil {
		resp.Diagnostics.AddError("WAITING FOR CREATION ERROR", err.Error())
		return
	}

	plan.refreshFromOutput(output, plan.PostgresqlInstanceNo.ValueStringPointer())

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *postgresqlReadReplicaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state postgresqlReadReplicaResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetPostgresqlReadReplica(ctx, r.config, state.PostgresqlInstanceNo.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.refreshFromOutput(output, state.PostgresqlInstanceNo.ValueStringPointer())

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
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

	if err := waitPostgresqlReadReplicaDeletion(ctx, r.config, state.PostgresqlInstanceNo.ValueString()); err != nil {
		resp.Diagnostics.AddError("WAITING FOR DELETE ERROR", err.Error())
	}
}

func waitPostgresqlReadReplicaCreation(ctx context.Context, config *conn.ProviderConfig, id string) (*vpostgresql.CloudPostgresqlServerInstance, error) {
	var postgresqlInstance *vpostgresql.CloudPostgresqlServerInstance
	stateConf := &retry.StateChangeConf{
		Pending: []string{"creating", "settingUp"},
		Target:  []string{"running"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetPostgresqlReadReplica(ctx, config, id)
			postgresqlInstance = instance
			if err != nil {
				return 0, "", err
			}

			status := instance.CloudPostgresqlServerInstanceStatus.Code
			op := instance.CloudPostgresqlServerInstanceOperation.Code

			if *status == "PEND" && *op == "CREAT" {
				return instance, "creating", nil
			}

			if *status == "RUN" && *op == "SETUP" {
				return instance, "settingUp", nil
			}

			if *status == "RUN" && *op == "NOOP" {
				return instance, "running", nil
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

	return postgresqlInstance, nil
}

func waitPostgresqlReadReplicaDeletion(ctx context.Context, config *conn.ProviderConfig, id string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{DELETING},
		Target:  []string{DELETED},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetPostgresqlReadReplica(ctx, config, id)
			if err != nil {
				return 0, "", err
			}

			if instance == nil {
				return instance, DELETED, nil
			}

			statusName := instance.CloudPostgresqlServerInstanceStatusName

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
		return fmt.Errorf("error waiting for postgresql read replica (%s) to become termintaing: %s", id, err)
	}

	return nil
}

func GetPostgresqlReadReplica(ctx context.Context, config *conn.ProviderConfig, no string) (*vpostgresql.CloudPostgresqlServerInstance, error) {
	reqParams := &vpostgresql.GetCloudPostgresqlInstanceDetailRequest{
		RegionCode:                &config.RegionCode,
		CloudPostgresqlInstanceNo: ncloud.String(no),
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
		if *server.CloudPostgresqlServerRole.Code == "S" {
			return server, nil
		}
	}
	return nil, nil
}

type postgresqlReadReplicaResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	PostgresqlInstanceNo types.String `tfsdk:"postgresql_instance_no"`
	SubnetNo             types.String `tfsdk:"subnet_no"`
}

func (r *postgresqlReadReplicaResourceModel) refreshFromOutput(output *vpostgresql.CloudPostgresqlServerInstance, id *string) {
	r.ID = types.StringPointerValue(output.CloudPostgresqlServerInstanceNo)
	r.PostgresqlInstanceNo = types.StringPointerValue(id)
	r.SubnetNo = types.StringPointerValue(output.SubnetNo)
}
