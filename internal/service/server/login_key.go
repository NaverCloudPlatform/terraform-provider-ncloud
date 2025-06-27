package server

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sdkresource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"
)

var (
	_ resource.Resource                = &loginKeyResource{}
	_ resource.ResourceWithConfigure   = &loginKeyResource{}
	_ resource.ResourceWithImportState = &loginKeyResource{}
)

type loginKeyResourceModel struct {
	KeyName     types.String `tfsdk:"key_name"`
	PrivateKey  types.String `tfsdk:"private_key"`
	Fingerprint types.String `tfsdk:"fingerprint"`
	ID          types.String `tfsdk:"id"`
}

type loginKeyResource struct {
	config *conn.ProviderConfig
}

func NewLoginKeyResource() resource.Resource {
	return &loginKeyResource{}
}

func (l *loginKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("key_name"), req, resp)
}

func (l *loginKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_login_key"
}

func (l *loginKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"key_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 30),
				},
				Description: "Key name to generate. If the generated key name exists, an error occurs.",
			},
			"private_key": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"fingerprint": schema.StringAttribute{
				Computed: true,
			},
			"id": framework.IDAttribute(),
		},
	}
}

func (l *loginKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	l.config = config
}

func (l *loginKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan loginKeyResourceModel
	var err error
	var privatekey *string

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	keyName := plan.KeyName.ValueStringPointer()

	privatekey, err = createVpcLoginKey(ctx, l.config, keyName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating LoginKey",
			err.Error(),
		)
		return
	}

	output, err := waitForNcloudLoginKeyCreation(l.config, *keyName)
	if err != nil {
		resp.Diagnostics.AddError("waiting for LoginKey creation", err.Error())
		return
	}

	plan.refreshFromOutput(output)
	plan.PrivateKey = types.StringValue(strings.TrimSpace(*privatekey))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (l *loginKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state loginKeyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetLoginKey(l.config, state.KeyName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("GetLoginKey", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.refreshFromOutput(output)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (l *loginKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (l *loginKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state loginKeyResourceModel
	var err error

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	keyName := state.KeyName.ValueString()

	tflog.Info(ctx, "DeleteLoginKey", map[string]any{
		"KeyName": common.MarshalUncheckedString(keyName),
	})

	err = deleteVpcLoginKey(ctx, l.config, keyName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting LoginKey",
			err.Error(),
		)
		return
	}
}

func createVpcLoginKey(ctx context.Context, config *conn.ProviderConfig, keyName *string) (*string, error) {
	reqParams := &vserver.CreateLoginKeyRequest{KeyName: keyName}
	tflog.Info(ctx, "DeleteVpcLoginKey", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})

	resp, err := config.Client.Vserver.V2Api.CreateLoginKey(reqParams)
	tflog.Info(ctx, "CreateVpcLoginKey response", map[string]any{
		"createVpcLoginKeyResponse": common.MarshalUncheckedString(resp),
	})

	return resp.PrivateKey, err
}

func waitForNcloudLoginKeyCreation(config *conn.ProviderConfig, keyName string) (*LoginKey, error) {
	var loginkey *LoginKey

	stateConf := &sdkresource.StateChangeConf{
		Pending: []string{""},
		Target:  []string{"OK"},
		Refresh: func() (interface{}, string, error) {
			resp, err := GetLoginKey(config, keyName)
			loginkey = resp
			if err != nil {
				return 0, "", err
			}

			if *resp.KeyName == keyName {
				return 0, "OK", err
			}

			return resp, "", nil
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return nil, fmt.Errorf("error waiting for Loginkey (%s) to become available: %s", keyName, err)
	}

	return loginkey, nil
}

type LoginKey struct {
	KeyName     *string `json:"key_name,omitempty"`
	Fingerprint *string `json:"fingerprint,omitempty"`
}

func GetLoginKey(config *conn.ProviderConfig, keyName string) (*LoginKey, error) {
	resp, err := config.Client.Vserver.V2Api.GetLoginKeyList(&vserver.GetLoginKeyListRequest{
		KeyName: ncloud.String(keyName),
	})

	if err != nil {
		return nil, err
	}

	if len(resp.LoginKeyList) < 1 {
		return nil, nil
	}

	l := resp.LoginKeyList[0]
	return &LoginKey{
		KeyName:     l.KeyName,
		Fingerprint: l.Fingerprint,
	}, nil
}

func deleteVpcLoginKey(ctx context.Context, config *conn.ProviderConfig, keyName string) error {
	reqParams := &vserver.DeleteLoginKeysRequest{KeyNameList: []*string{ncloud.String(keyName)}}
	tflog.Info(ctx, "DeletVpcLoginKey", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})

	resp, err := config.Client.Vserver.V2Api.DeleteLoginKeys(reqParams)
	if err != nil {
		common.LogErrorResponse("deleteVpcLoginKey", err, keyName)
		return err
	}
	tflog.Info(ctx, "DeleteVpcLoginKey response", map[string]any{
		"deleteVpcLoginKeyResponse": common.MarshalUncheckedString(resp),
	})

	stateConf := &sdkresource.StateChangeConf{
		Pending: []string{""},
		Target:  []string{"OK"},
		Refresh: func() (interface{}, string, error) {
			resp, err := GetLoginKey(config, keyName)
			if err != nil {
				return 0, "", err
			}

			if resp == nil {
				return 0, "OK", err
			}

			return resp, "", nil
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting to delete LoginKey: %v", err)
	}

	return nil
}

func (l *loginKeyResourceModel) refreshFromOutput(output *LoginKey) {
	l.ID = types.StringPointerValue(output.KeyName)
	l.KeyName = types.StringPointerValue(output.KeyName)
	l.Fingerprint = types.StringPointerValue(output.Fingerprint)
}
