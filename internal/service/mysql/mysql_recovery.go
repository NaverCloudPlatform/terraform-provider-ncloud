package mysql

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmysql"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	idParts := strings.Split(req.ID, ":")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: mysql_instance_no:id Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("mysql_instance_no"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
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
			"subnet_no": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
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
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("recovery_time"),
					),
				},
			},
			"recovery_time": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("file_name"),
					),
				},
			},
			"mysql_server_list": schema.ListNestedAttribute{
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
						"public_domain": schema.StringAttribute{
							Computed: true,
						},
						"private_domain": schema.StringAttribute{
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

	if !plan.SubnetNo.IsNull() && !plan.SubnetNo.IsUnknown() {
		// In `gov`, multi_zone is always false, so subnet is auto-generated with default value
		if r.config.Site == "gov" {
			resp.Diagnostics.AddError(
				"NOT SUPPORT GOV SITE",
				"`subnet_no` does not support gov site",
			)
			return
		}
		reqParams.SubnetNo = plan.SubnetNo.ValueStringPointer()
	}

	tflog.Info(ctx, "CreateMysqlRecovery reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vmysql.V2Api.CreateCloudMysqlRecoveryInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "CreateMysqlRecovery response="+common.MarshalUncheckedString(response))

	if response == nil || len(response.CloudMysqlInstanceList) < 1 {
		resp.Diagnostics.AddError("CREATING ERROR", "response valid")
		return
	}

	mysqlIns := response.CloudMysqlInstanceList[0]
	serverList := mysqlIns.CloudMysqlServerInstanceList
	var index int

	for i, server := range serverList {
		if (server.CloudMysqlServerRole != nil && *server.CloudMysqlServerRole.Code == "R") && (*server.CloudMysqlServerInstanceStatusName == CREATING) {
			index = i
			break
		}
	}

	output, err := waitMysqlServerCreation(ctx, r.config, *mysqlIns.CloudMysqlInstanceNo, index)
	if err != nil {
		resp.Diagnostics.AddError("WAITING FOR CREATION ERROR", err.Error())
		return
	}

	if diags := plan.refreshFromOutput(ctx, output, mysqlIns.CloudMysqlInstanceNo); diags.HasError() {
		resp.Diagnostics.AddError("CREATING ERROR", "refreshFromOutput error")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *mysqlRecoveryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state mysqlRecoveryResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetMysqlRecovery(ctx, r.config, state.MysqlInstanceNo.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if diags := state.refreshFromOutput(ctx, output, state.MysqlInstanceNo.ValueStringPointer()); diags.HasError() {
		resp.Diagnostics.AddError("READING ERROR", "refreshFromOutput error")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
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

	if err := waitMysqlRecoveryDeletion(ctx, r.config, state.MysqlInstanceNo.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("WAITING FOR DELETE ERROR", err.Error())
		return
	}
}

func waitMysqlRecoveryDeletion(ctx context.Context, config *conn.ProviderConfig, instanceNo string, serverInstanceNo string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{DELETING},
		Target:  []string{DELETED},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetMysqlRecovery(ctx, config, instanceNo, serverInstanceNo)
			if err != nil {
				return 0, "", err
			}

			if instance == nil {
				return instance, DELETED, nil
			}

			status := instance[0].CloudMysqlServerInstanceStatusName

			if *status == DELETING || *status == DELETED {
				return instance, *status, nil
			}

			return 0, "", fmt.Errorf("error occurred while waiting to delete mysql recovery")
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      1 * time.Minute,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for mysql recovery (%s) to become terminating: %s", serverInstanceNo, err)
	}

	return nil
}

func GetMysqlRecovery(ctx context.Context, config *conn.ProviderConfig, instanceNo string, serverInstanceNo string) ([]*vmysql.CloudMysqlServerInstance, error) {
	reqParams := &vmysql.GetCloudMysqlInstanceDetailRequest{
		RegionCode:           &config.RegionCode,
		CloudMysqlInstanceNo: ncloud.String(instanceNo),
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
		if (server.CloudMysqlServerRole != nil && *server.CloudMysqlServerRole.Code == "R") && (*server.CloudMysqlServerInstanceNo == serverInstanceNo) {
			return []*vmysql.CloudMysqlServerInstance{server}, nil
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
	MysqlServerList         types.List   `tfsdk:"mysql_server_list"`
}

func (r *mysqlRecoveryResourceModel) refreshFromOutput(ctx context.Context, output []*vmysql.CloudMysqlServerInstance, instanceNo *string) diag.Diagnostics {
	r.ID = types.StringPointerValue(output[0].CloudMysqlServerInstanceNo)
	r.MysqlInstanceNo = types.StringPointerValue(instanceNo)
	r.SubnetNo = types.StringPointerValue(output[0].SubnetNo)

	serverList, diags := listValueFromMysqlServerList(ctx, output)
	r.MysqlServerList = serverList

	return diags
}
