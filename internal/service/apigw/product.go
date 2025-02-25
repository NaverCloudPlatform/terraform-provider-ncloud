package apigw

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/ncloudsdk"
)

func ProductResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"product_name": schema.StringAttribute{
				Required:            true,
				Description:         "Product Name<br>Length(Min/Max): 0/100",
				MarkdownDescription: "Product Name<br>Length(Min/Max): 0/100",
			},
			"subscription_code": schema.StringAttribute{
				Required:            true,
				Description:         "Subscription Code<br>Allowable values: PROTECTED, PUBLIC",
				MarkdownDescription: "Subscription Code<br>Allowable values: PROTECTED, PUBLIC",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"PROTECTED",
						"PUBLIC",
					),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Description<br>Length(Min/Max): 0/300",
				MarkdownDescription: "Description<br>Length(Min/Max): 0/300",
			},
			"tenant_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Tenant Id",
				MarkdownDescription: "Tenant Id",
			},
			"published": schema.BoolAttribute{
				Computed:            true,
				Description:         "Is Published",
				MarkdownDescription: "Is Published",
			},
			"modifier": schema.StringAttribute{
				Computed:            true,
				Description:         "Modifier",
				MarkdownDescription: "Modifier",
			},
			"domain_code": schema.StringAttribute{
				Computed:            true,
				Description:         "Domain Code",
				MarkdownDescription: "Domain Code",
			},
			"deleted": schema.BoolAttribute{
				Computed:            true,
				Description:         "Is Deleted",
				MarkdownDescription: "Is Deleted",
			},
			"mod_time": schema.StringAttribute{
				Computed:            true,
				Description:         "Mod Time",
				MarkdownDescription: "Mod Time",
			},
			"zone_code": schema.StringAttribute{
				Computed:            true,
				Description:         "Zone Code",
				MarkdownDescription: "Zone Code",
			},
		},
	}
}

func NewProductResource() resource.Resource {
	return &productResource{}
}

type productResource struct {
	config *conn.ProviderConfig
}

func (a *productResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	a.config = config
}

func (a *productResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_apigw_product"
}

func (a *productResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ProductResourceSchema(ctx)
}

func (a *productResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (a *productResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PostproductresponseModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &ncloudsdk.POSTProductsRequestBody{
		ProductName:      plan.ProductName.ValueStringPointer(),
		SubscriptionCode: plan.SubscriptionCode.ValueStringPointer(),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		reqParams.Description = plan.Description.ValueStringPointer()
	}

	tflog.Info(ctx, "CreateProduct reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := a.config.Client.Apigw.POSTProducts(ctx, reqParams)
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}

	if response == nil {
		resp.Diagnostics.AddError("CREATING ERROR", "response invalid")
		return
	}

	tflog.Info(ctx, "CreateProduct response="+common.MarshalUncheckedString(response))

	plan.refreshFromOutput_createOp(response)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (a *productResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan PostproductresponseModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.refreshFromOutput(ctx, &resp.Diagnostics, plan.ID.ValueString(), a)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (a *productResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state PostproductresponseModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqQueryParams := &ncloudsdk.PATCHProductsProductidRequestQuery{
		Productid: state.ID.ValueStringPointer(),
	}

	tflog.Info(ctx, "UpdateProducts reqQueryParams="+common.MarshalUncheckedString(reqQueryParams))

	reqBodyParams := &ncloudsdk.PATCHProductsProductidRequestBody{
		ProductName:      plan.ProductName.ValueStringPointer(),
		SubscriptionCode: plan.SubscriptionCode.ValueStringPointer(),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		reqBodyParams.Description = plan.Description.ValueStringPointer()
	}

	tflog.Info(ctx, "UpdateProducts reqBodyParams="+common.MarshalUncheckedString(reqBodyParams))

	response, err := a.config.Client.Apigw.PATCHProductsProductid(ctx, reqQueryParams, reqBodyParams)
	if err != nil {
		resp.Diagnostics.AddError("UPDATING ERROR", err.Error())
		return
	}
	if response == nil {
		resp.Diagnostics.AddError("UPDATING ERROR", "response invalid")
		return
	}

	tflog.Info(ctx, "UpdateProducts response="+common.MarshalUncheckedString(response))

	plan.refreshFromOutput(ctx, &resp.Diagnostics, state.ID.ValueString(), a)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (a *productResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state PostproductresponseModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &ncloudsdk.DELETEProductsProductidRequestQuery{
		Productid: state.ID.ValueStringPointer(),
	}

	tflog.Info(ctx, "DELETEProducts reqParams="+common.MarshalUncheckedString(reqParams))

	_, err := a.config.Client.Apigw.DELETEProductsProductid(ctx, reqParams)
	if err != nil {
		resp.Diagnostics.AddError("DELETING ERROR", err.Error())
		return
	}
}

type PostproductresponseModel struct {
	ID               types.String `tfsdk:"id"`
	ProductName      types.String `tfsdk:"product_name"`
	SubscriptionCode types.String `tfsdk:"subscription_code"`
	Description      types.String `tfsdk:"description"`
	TenantId         types.String `tfsdk:"tenant_id"`
	Published        types.Bool   `tfsdk:"published"`
	Modifier         types.String `tfsdk:"modifier"`
	DomainCode       types.String `tfsdk:"domain_code"`
	Deleted          types.Bool   `tfsdk:"deleted"`
	ModTime          types.String `tfsdk:"mod_time"`
	ZoneCode         types.String `tfsdk:"zone_code"`
}

func (plan *PostproductresponseModel) refreshFromOutput_createOp(resp map[string]interface{}) {
	// Allocate resource id from create response
	plan.ID = types.StringValue(resp["product"].(map[string]interface{})["product_id"].(string))
	plan.ProductName = types.StringValue(resp["product"].(map[string]interface{})["product_name"].(string))
	plan.SubscriptionCode = types.StringValue(resp["product"].(map[string]interface{})["subscription_code"].(string))
	plan.Description = types.StringValue(resp["product"].(map[string]interface{})["product_description"].(string))
	plan.TenantId = types.StringValue(resp["product"].(map[string]interface{})["tenant_id"].(string))
	plan.Published = types.BoolValue(resp["product"].(map[string]interface{})["is_published"].(bool))
	plan.Modifier = types.StringValue(resp["product"].(map[string]interface{})["modifier"].(string))
	plan.DomainCode = types.StringValue(resp["product"].(map[string]interface{})["domain_code"].(string))
	plan.Deleted = types.BoolValue(resp["product"].(map[string]interface{})["is_deleted"].(bool))
	plan.ModTime = types.StringValue(resp["product"].(map[string]interface{})["mod_time"].(string))
	plan.ZoneCode = types.StringValue(resp["product"].(map[string]interface{})["zone_code"].(string))
}

func (plan *PostproductresponseModel) refreshFromOutput(ctx context.Context, diagnostics *diag.Diagnostics, id string, a *productResource) {
	resp, err := a.config.Client.Apigw.GETProductsProductid(ctx, &ncloudsdk.GETProductsProductidRequestQuery{
		Productid: &id,
	})
	tflog.Info(ctx, "GetProduct response="+common.MarshalUncheckedString(resp))

	if err != nil {
		diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	plan.ID = types.StringValue(resp["product"].(map[string]interface{})["product_id"].(string))
	plan.ProductName = types.StringValue(resp["product"].(map[string]interface{})["product_name"].(string))
	plan.SubscriptionCode = types.StringValue(resp["product"].(map[string]interface{})["subscription_code"].(string))
	plan.Description = types.StringValue(resp["product"].(map[string]interface{})["product_description"].(string))
	plan.TenantId = types.StringValue(resp["product"].(map[string]interface{})["tenant_id"].(string))
	plan.Published = types.BoolValue(resp["product"].(map[string]interface{})["is_published"].(bool))
	plan.Modifier = types.StringValue(resp["product"].(map[string]interface{})["modifier"].(string))
	plan.DomainCode = types.StringValue(resp["product"].(map[string]interface{})["domain_code"].(string))
	plan.Deleted = types.BoolValue(resp["product"].(map[string]interface{})["is_deleted"].(bool))
	plan.ModTime = types.StringValue(resp["product"].(map[string]interface{})["mod_time"].(string))
	plan.ZoneCode = types.StringValue(resp["product"].(map[string]interface{})["zone_code"].(string))
}
