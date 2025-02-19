package apigw

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
			"invoke_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Invoke Id",
				MarkdownDescription: "Invoke Id",
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

	c := ncloudsdk.NewClient("https://apigateway.apigw.ntruss.com/api/v1", os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"))

	reqParams := &ncloudsdk.PrimitivePOSTProductsRequest{
		ProductName:      plan.ProductName.ValueString(),
		SubscriptionCode: plan.SubscriptionCode.ValueString(),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		reqParams.Description = plan.Description.ValueString()
	}

	tflog.Info(ctx, "CreateProduct reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := c.POSTProducts(ctx, reqParams)
	if err != nil {
		resp.Diagnostics.AddError("Error with POSTProducts_TF", err.Error())
		return
	}

	tflog.Info(ctx, "CreateProduct response="+common.MarshalUncheckedString(response))

	plan.refreshFromOutput_createOp(ctx, &resp.Diagnostics, response)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (a *productResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan PostproductresponseModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.refreshFromOutput(ctx, &resp.Diagnostics, plan.ID.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (a *productResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan, state PostproductresponseModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &ncloudsdk.PrimitivePATCHProductsProductidRequest{
		Productid:        plan.Productid.ValueString(),
		ProductName:      plan.ProductName.ValueString(),
		SubscriptionCode: plan.SubscriptionCode.ValueString(),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		reqParams.Description = plan.Description.ValueString()
	}

	tflog.Info(ctx, "UpdatePATCHProductsProductid reqParams="+common.MarshalUncheckedString(reqParams))

	c := ncloudsdk.NewClient("https://apigateway.apigw.ntruss.com/api/v1", os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"))

	response, err := c.PATCHProductsProductid_TF(ctx, reqParams)
	if err != nil {
		resp.Diagnostics.AddError("UPDATING ERROR", err.Error())
		return
	}
	if response == nil {
		resp.Diagnostics.AddError("UPDATING ERROR", "response invalid")
		return
	}

	tflog.Info(ctx, "UpdatePATCHProductsProductid response="+common.MarshalUncheckedString(response))

	plan.refreshFromOutput(ctx, &resp.Diagnostics, state.ID.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

}

func (a *productResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state PostproductresponseModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &ncloudsdk.PrimitiveDELETEProductsProductidRequest{
		Productid: plan.Productid.ValueString(),
	}

	tflog.Info(ctx, "UpdateDELETEProductsProductid reqParams="+common.MarshalUncheckedString(reqParams))

	c := ncloudsdk.NewClient("https://apigateway.apigw.ntruss.com/api/v1", os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"))

	_, err := c.DELETEProductsProductid_TF(ctx, reqParams)
	if err != nil {
		resp.Diagnostics.AddError("DELETING ERROR", err.Error())
		return
	}

	err = state.waitResourceDeleted(ctx, state.ID.ValueString())
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
	InvodkeId        types.String `tfsdk:"invoke_id"`
	TenantId         types.String `tfsdk:"tenant_id"`
	Published        types.Bool   `tfsdk:"published"`
	Modifier         types.String `tfsdk:"modifier"`
	DomainCode       types.String `tfsdk:"domain_code"`
	Deleted          types.Bool   `tfsdk:"deleted"`
	ModTime          types.String `tfsdk:"mod_time"`
	ZoneCode         types.String `tfsdk:"zone_code"`
}
