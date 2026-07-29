package subaccount

import (
	"context"
	"fmt"
	"time"

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
	subaccountsdk "github.com/terraform-providers/terraform-provider-ncloud/internal/sdk/subaccount"
)

var (
	_ resource.Resource              = &subAccountAccessKeyResource{}
	_ resource.ResourceWithConfigure = &subAccountAccessKeyResource{}
)

func NewSubAccountAccessKeyResource() resource.Resource {
	return &subAccountAccessKeyResource{}
}

type subAccountAccessKeyResource struct {
	config *conn.ProviderConfig
}

type subAccountAccessKeyResourceModel struct {
	ID           types.String `tfsdk:"id"`
	SubAccountId types.String `tfsdk:"sub_account_id"`
	AccessKey    types.String `tfsdk:"access_key"`
	SecretKey    types.String `tfsdk:"secret_key"`
	CreateTime   types.String `tfsdk:"create_time"`
}

func (r *subAccountAccessKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subaccount_access_key"
}

func (r *subAccountAccessKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttribute(),
			"sub_account_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "ID of the sub account to issue the access key for. " +
					"The sub account must have `can_api_gateway_access` enabled.",
			},
			"access_key": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Access key ID.",
			},
			"secret_key": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Secret key. Returned by the API only at creation time and stored in state.",
			},
			"create_time": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Creation time of the access key.",
			},
		},
	}
}

func (r *subAccountAccessKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *subAccountAccessKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan subAccountAccessKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subAccountId := plan.SubAccountId.ValueString()

	tflog.Info(ctx, "CreateSubAccountAccessKey", map[string]any{
		"subAccountId": subAccountId,
	})

	key, err := r.config.Client.SubAccount.CreateAccessKey(ctx, subAccountId)
	if err != nil {
		common.LogErrorResponse("CreateSubAccountAccessKey", err, subAccountId)
		resp.Diagnostics.AddError("Error Creating SubAccount Access Key", err.Error())
		return
	}

	plan.ID = types.StringValue(key.AccessKey)
	plan.AccessKey = types.StringValue(key.AccessKey)
	plan.SecretKey = types.StringValue(key.KeySecret)
	plan.CreateTime = types.StringValue(key.CreateTime)

	// Persist state before the visibility wait: the key already exists and
	// its secret can never be retrieved again, so a transient failure below
	// must not lose it.
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := waitForAccessKey(ctx, r.config, subAccountId, key.AccessKey); err != nil {
		resp.Diagnostics.AddError("Error Creating SubAccount Access Key", err.Error())
	}
}

func (r *subAccountAccessKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state subAccountAccessKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	keys, err := r.config.Client.SubAccount.ListAccessKeys(ctx, state.SubAccountId.ValueString())
	if err != nil {
		if subaccountsdk.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		common.LogErrorResponse("ListSubAccountAccessKeys", err, state.SubAccountId.ValueString())
		resp.Diagnostics.AddError("Error Reading SubAccount Access Key", err.Error())
		return
	}

	// The list response never contains the secret; keep the value captured at creation.
	for _, key := range keys {
		if key.AccessKey == state.ID.ValueString() {
			state.AccessKey = types.StringValue(key.AccessKey)
			state.CreateTime = types.StringValue(key.CreateTime)

			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	resp.State.RemoveResource(ctx)
}

func (r *subAccountAccessKeyResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *subAccountAccessKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state subAccountAccessKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "DeleteSubAccountAccessKey", map[string]any{
		"subAccountId": state.SubAccountId.ValueString(),
		"accessKey":    state.ID.ValueString(),
	})

	err := r.config.Client.SubAccount.DeleteAccessKey(ctx, state.SubAccountId.ValueString(), state.ID.ValueString())
	if err != nil && !subaccountsdk.IsNotFound(err) {
		common.LogErrorResponse("DeleteSubAccountAccessKey", err, state.ID.ValueString())
		resp.Diagnostics.AddError("Error Deleting SubAccount Access Key", err.Error())
	}
}

func waitForAccessKey(ctx context.Context, config *conn.ProviderConfig, subAccountId, accessKey string) error {
	listContains := func() (bool, error) {
		keys, err := config.Client.SubAccount.ListAccessKeys(ctx, subAccountId)
		if err != nil {
			return false, err
		}
		for _, key := range keys {
			if key.AccessKey == accessKey {
				return true, nil
			}
		}
		return false, nil
	}

	// Issuance is normally synchronous; poll only when the first read
	// misses the key (eventual consistency).
	if found, err := listContains(); err != nil || found {
		return err
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{""},
		Target:  []string{"OK"},
		Refresh: func() (any, string, error) {
			found, err := listContains()
			if err != nil {
				return 0, "", err
			}
			if found {
				return 0, "OK", nil
			}
			return 0, "", nil
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for SubAccount access key (%s) to become available: %w", accessKey, err)
	}

	return nil
}
