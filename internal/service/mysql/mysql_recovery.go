package mysql

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmysql"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"
)

const (
	CREATING = "creating"
	SETTING  = "settingUp"
	RUNNING  = "running"
	DELETING = "deleting"
	DELETED  = "deleted"
)

var (
	_ resource.Resource                = &mysqlRecoveryResource{}
	_ resource.ResourceWithConfigure   = &mysqlRecoveryResource{}
	_ resource.ResourceWithImportState = &mysqlRecoveryResource{}
)

func NewMysqlRecoveryResource() resource.Resource {
	return &mysqlRecoveryResource{}
}

type mysqlRecoveryResource struct {
	config *conn.ProviderConfig
}

func (r *mysqlRecoveryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *mysqlRecoveryResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *mysqlRecoveryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mysql_recovery"
}

func (r *mysqlRecoveryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttribute(),
			"mysql_instance_no": schema.StringAttribute{
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
			"recovery_server_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 25),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-z]+[a-z0-9-]+[a-z0-9]$`),
						"Composed of lowercase alphabets, numbers, hyphen (-). Must start with an alphabetic character, and the last character can only be an English letter or number.",
					),
				},
			},
			"file_name": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(path.Expressions{
						path.MatchRoot("recovery_time"),
					}...),
				},
			},
			"recovery_time": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *mysqlRecoveryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan mysqlRecoveryResourceModel

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

	reqParams := &vmysql.CreateCloudMysqlRecoveryInstanceRequest{
		RegionCode:                   &r.config.RegionCode,
		CloudMysqlInstanceNo:         plan.MysqlInstanceNo.ValueStringPointer(),
		CloudMysqlRecoveryServerName: plan.MysqlRecoveryServerName.ValueStringPointer(),
	}

	if !plan.FileName.IsNull() {
		reqParams.FileName = plan.FileName.ValueStringPointer()
	}

	if !plan.RecoveryTime.IsNull() {
		reqParams.RecoveryTime = plan.RecoveryTime.ValueStringPointer()
	}

	tflog.Info(ctx, "CreateCloudMysqlRecovery reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vmysql.V2Api.CreateCloudMysqlRecoveryInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}

	if response == nil || len(response.CloudMysqlInstanceList) < 1 {
		resp.Diagnostics.AddError("CREATING ERROR", "response valid")
		return
	}

	mysqlIns := response.CloudMysqlInstanceList[0]
	serverList := mysqlIns.CloudMysqlServerInstanceList

	if *mysqlIns.IsMultiZone {
		if !plan.SubnetNo.IsNull() {
			reqParams.SubnetNo = plan.SubnetNo.ValueStringPointer()
		} else {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"when `is_multi_zone` is true, SubnetNo should be set",
			)
		}
	}

	if len(serverList) < 2 {
		resp.Diagnostics.AddError("CREATING ERROR", "response invalid")
		return
	}

	recoveryServer, err := GetMysqlRecovery(ctx, r.config, plan.MysqlInstanceNo.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("READ ERROR", err.Error())
		return
	}

	plan.ID = types.StringPointerValue(recoveryServer.CloudMysqlServerInstanceNo)

	output, err := waitMysqlRecoveryCreation(ctx, r.config, *mysqlIns.CloudMysqlInstanceNo)
	if err != nil {
		resp.Diagnostics.AddError("WAITING FOR CREATION ERROR", err.Error())
		return
	}

	plan.refreshFromOutput(output, plan.MysqlInstanceNo.ValueStringPointer())

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *mysqlRecoveryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state mysqlRecoveryResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetMysqlRecovery(ctx, r.config, state.MysqlInstanceNo.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.refreshFromOutput(output, state.MysqlInstanceNo.ValueStringPointer())

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *mysqlRecoveryResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *mysqlRecoveryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state mysqlRecoveryResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vmysql.DeleteCloudMysqlServerInstanceRequest{
		RegionCode:                 &r.config.RegionCode,
		CloudMysqlServerInstanceNo: state.ID.ValueStringPointer(),
	}
	tflog.Info(ctx, "DeleteMysqlRecovery reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vmysql.V2Api.DeleteCloudMysqlServerInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("DELETING ERROR", err.Error())
		return
	}

	tflog.Info(ctx, "DeleteMysqlRecovery response="+common.MarshalUncheckedString(response))

	if err := waitMysqlRecoveryDeletion(ctx, r.config, state.MysqlInstanceNo.ValueString()); err != nil {
		resp.Diagnostics.AddError("WAITING FOR DELETE ERROR", err.Error())
	}
}

func waitMysqlRecoveryCreation(ctx context.Context, config *conn.ProviderConfig, id string) (*vmysql.CloudMysqlServerInstance, error) {
	var mysqlInstance *vmysql.CloudMysqlServerInstance
	stateConf := &retry.StateChangeConf{
		Pending: []string{CREATING, SETTING},
		Target:  []string{RUNNING},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetMysqlRecovery(ctx, config, id)
			mysqlInstance = instance
			if err != nil {
				return 0, "", err
			}

			status := instance.CloudMysqlServerInstanceStatusName
			if *status == CREATING {
				return instance, CREATING, nil
			}

			if *status == SETTING {
				return instance, SETTING, nil
			}

			if *status == RUNNING {
				return instance, RUNNING, nil
			}

			return 0, "", fmt.Errorf("error occurred while waiting to create mysql recovery")
		},
		Timeout:    6 * conn.DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error waiting for mysql recovery state to be \"running\": %s", err)
	}

	return mysqlInstance, nil
}

func waitMysqlRecoveryDeletion(ctx context.Context, config *conn.ProviderConfig, id string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{DELETING},
		Target:  []string{DELETED},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetMysqlRecovery(ctx, config, id)
			if err != nil {
				return 0, "", err
			}

			if instance == nil {
				return instance, DELETED, nil
			}

			status := instance.CloudMysqlServerInstanceStatusName

			if *status == DELETING {
				return instance, DELETING, nil
			}

			if *status == DELETED {
				return instance, DELETED, nil
			}

			return 0, "", fmt.Errorf("error occurred while waiting to delete mysql recovery")
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for mysql recovery (%s) to become terminating: %s", id, err)
	}

	return nil
}

func GetMysqlRecovery(ctx context.Context, config *conn.ProviderConfig, no string) (*vmysql.CloudMysqlServerInstance, error) {
	reqParams := &vmysql.GetCloudMysqlInstanceDetailRequest{
		RegionCode:           &config.RegionCode,
		CloudMysqlInstanceNo: ncloud.String(no),
	}
	tflog.Info(ctx, "GetMysqlDetail reqParams="+common.MarshalUncheckedString(reqParams))

	resp, err := config.Client.Vmysql.V2Api.GetCloudMysqlInstanceDetail(reqParams)
	if err != nil && !CheckIfAlreadyDeleted(err) {
		return nil, err
	}
	tflog.Info(ctx, "GetMysqlDetail response="+common.MarshalUncheckedString(resp))

	if resp == nil || len(resp.CloudMysqlInstanceList) < 1 {
		return nil, nil
	}

	serverList := resp.CloudMysqlInstanceList[0].CloudMysqlServerInstanceList

	for _, server := range serverList {
		if *server.CloudMysqlServerRole.CodeName == "Recovery" {
			return server, nil
		}
	}
	return nil, nil
}

type mysqlRecoveryResourceModel struct {
	ID                      types.String `tfsdk:"id"`
	MysqlInstanceNo         types.String `tfsdk:"mysql_instance_no"`
	SubnetNo                types.String `tfsdk:"subnet_no"`
	MysqlRecoveryServerName types.String `tfsdk:"recovery_server_name"`
	FileName                types.String `tfsdk:"file_name"`
	RecoveryTime            types.String `tfsdk:"recovery_time"`
}

func (r *mysqlRecoveryResourceModel) refreshFromOutput(output *vmysql.CloudMysqlServerInstance, id *string) {
	r.ID = types.StringPointerValue(output.CloudMysqlServerInstanceNo)
	r.MysqlInstanceNo = types.StringPointerValue(id)
	r.SubnetNo = types.StringPointerValue(output.SubnetNo)
}
