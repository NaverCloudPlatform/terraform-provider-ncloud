package redis

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sdkresource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vredis"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

var (
	_ resource.Resource                = &redisConfigGroupResource{}
	_ resource.ResourceWithConfigure   = &redisConfigGroupResource{}
	_ resource.ResourceWithImportState = &redisConfigGroupResource{}
)

func NewRedisConfigGroupResource() resource.Resource {
	return &redisConfigGroupResource{}
}

type redisConfigGroupResource struct {
	config *conn.ProviderConfig
}

func (r *redisConfigGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_redis_config_group"
}

func (r *redisConfigGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 15),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-z]+[a-z0-9-]+[a-z0-9]$`),
						"Composed of lowercase alphabets, numbers, hyphen (-). Must start with an alphabetic character, and the last character can only be an English letter or number.",
					),
				},
			},
			"redis_version": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"id": framework.IDAttribute(),
			"config_group_no": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *redisConfigGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *redisConfigGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan redisConfigGroupResourceModel

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

	reqParams := &vredis.CreateCloudRedisConfigGroupRequest{
		RegionCode:             &r.config.RegionCode,
		CloudRedisVersion:      plan.RedisVersion.ValueStringPointer(),
		ConfigGroupName:        plan.Name.ValueStringPointer(),
		ConfigGroupDescription: plan.Description.ValueStringPointer(),
	}

	tflog.Info(ctx, "CreateCloudRedisConfigGroup reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vredis.V2Api.CreateCloudRedisConfigGroup(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "CreateCloudRedisConfigGroup response="+common.MarshalUncheckedString(response))

	if response == nil || len(response.CloudRedisConfigGroupList) < 1 {
		resp.Diagnostics.AddError("CREATING ERROR", "response invalid")
		return
	}

	// The response contains all the CloudRedisConfigGroupList.
	var confGroup *vredis.CloudRedisConfigGroup
	for _, v := range response.CloudRedisConfigGroupList {
		if *v.ConfigGroupName == plan.Name.ValueString() {
			confGroup = v
			break
		}
	}

	if confGroup == nil {
		resp.Diagnostics.AddError("CREATING ERROR", "response invalid")
		return
	}

	// Couldn't check Config Group number in console, so set ID to name
	plan.ID = types.StringPointerValue(confGroup.ConfigGroupName)

	output, err := waitRedisConfigGroupCreated(ctx, r.config, *confGroup.ConfigGroupName)
	if err != nil {
		resp.Diagnostics.AddError("WAITING FOR CREATION ERROR", err.Error())
		return
	}

	plan.refreshFromOutput(ctx, output)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *redisConfigGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state redisConfigGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetRedisConfigGroup(ctx, r.config, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.refreshFromOutput(ctx, output)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *redisConfigGroupResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *redisConfigGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state redisConfigGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vredis.DeleteCloudRedisConfigGroupRequest{
		RegionCode:    &r.config.RegionCode,
		ConfigGroupNo: state.ConfigGroupNo.ValueStringPointer(),
	}

	tflog.Info(ctx, "DeleteCloudRedisConfigGroup reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vredis.V2Api.DeleteCloudRedisConfigGroup(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("DELETING ERROR", err.Error())
		return
	}

	tflog.Info(ctx, "DeleteCloudRedisConfigGroup response="+common.MarshalUncheckedString(response))

	if err := waitRedisConfigGroupDeleted(ctx, r.config, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("WAITING FOR DELETE ERROR", err.Error())
	}
}

func (r *redisConfigGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func GetRedisConfigGroup(ctx context.Context, config *conn.ProviderConfig, name string) (*vredis.CloudRedisConfigGroup, error) {
	reqParams := &vredis.GetCloudRedisConfigGroupListRequest{
		RegionCode:      &config.RegionCode,
		ConfigGroupName: &name,
	}

	tflog.Info(ctx, "GetRedisConfigGroup reqParams="+common.MarshalUncheckedString(reqParams))

	resp, err := config.Client.Vredis.V2Api.GetCloudRedisConfigGroupList(reqParams)
	if err != nil {
		return nil, err
	}
	tflog.Info(ctx, "GetRedisConfigGroup response="+common.MarshalUncheckedString(resp))

	if resp == nil || len(resp.CloudRedisConfigGroupList) < 1 {
		return nil, nil
	}

	// The response contains all CloudRedisConfigGroupList with the same prefix name.
	var confGroup *vredis.CloudRedisConfigGroup
	for _, v := range resp.CloudRedisConfigGroupList {
		if *v.ConfigGroupName == name {
			confGroup = v
			break
		}
	}

	if confGroup == nil {
		return nil, nil
	}

	return confGroup, nil
}

func waitRedisConfigGroupDeleted(ctx context.Context, config *conn.ProviderConfig, name string) error {
	stateConf := &sdkresource.StateChangeConf{
		Pending: []string{"deleting"},
		Target:  []string{"deleted"},
		Refresh: func() (interface{}, string, error) {
			resp, err := GetRedisConfigGroup(ctx, config, name)
			if err != nil {
				return 0, "", err
			}

			if resp == nil || *resp.ConfigGroupStatusName == "deleted" {
				return resp, "deleted", nil
			}

			if *resp.ConfigGroupStatusName == "deleting" {
				return resp, "deleting", nil
			}

			return 0, "", fmt.Errorf("error occurred while waiting to delete")
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for Redis Config Group (%s) to become termintaing: %s", name, err)
	}

	return nil
}

func waitRedisConfigGroupCreated(ctx context.Context, config *conn.ProviderConfig, name string) (*vredis.CloudRedisConfigGroup, error) {
	var redisConfigGroup *vredis.CloudRedisConfigGroup
	stateConf := &sdkresource.StateChangeConf{
		Pending: []string{"creating", "settingUp"},
		Target:  []string{"running"},
		Refresh: func() (interface{}, string, error) {
			resp, err := GetRedisConfigGroup(ctx, config, name)
			redisConfigGroup = resp
			if err != nil {
				return 0, "", err
			}

			if resp == nil {
				return 0, "", fmt.Errorf("GetRedisConfigGroup is nil")
			}

			if *resp.ConfigGroupStatusName == "running" {
				return resp, "running", nil
			} else if *resp.ConfigGroupStatusName == "creating" {
				return resp, "creating", nil
			} else if *resp.ConfigGroupStatusName == "settingUp" {
				return resp, "settingUp", nil
			}

			return 0, "", fmt.Errorf("error occurred while waiting to create")
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return nil, fmt.Errorf("Error waiting for Redis Config Group (%s) to become available: %s", name, err)
	}

	return redisConfigGroup, nil
}

type redisConfigGroupResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	RedisVersion  types.String `tfsdk:"redis_version"`
	Description   types.String `tfsdk:"description"`
	ConfigGroupNo types.String `tfsdk:"config_group_no"`
}

func (r *redisConfigGroupResourceModel) refreshFromOutput(ctx context.Context, output *vredis.CloudRedisConfigGroup) {
	r.ID = types.StringPointerValue(output.ConfigGroupName)
	r.Name = types.StringPointerValue(output.ConfigGroupName)
	r.RedisVersion = types.StringPointerValue(output.CloudRedisVersion)
	r.Description = types.StringPointerValue(output.ConfigGroupDescription)
	r.ConfigGroupNo = types.StringPointerValue(output.ConfigGroupNo)

	return
}
