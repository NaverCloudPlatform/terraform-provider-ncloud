package postgresql

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpostgresql"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

var (
	_ resource.Resource                = &postgresqlDatabasesResource{}
	_ resource.ResourceWithConfigure   = &postgresqlDatabasesResource{}
	_ resource.ResourceWithImportState = &postgresqlDatabasesResource{}
)

func NewPostgresqlDatabasesResource() resource.Resource {
	return &postgresqlDatabasesResource{}
}

type postgresqlDatabasesResource struct {
	config *conn.ProviderConfig
}

func (r *postgresqlDatabasesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var plan postgresqlDatabasesResourceModel
	var databaseList []PostgresqlDatabase
	idParts := strings.Split(req.ID, ":")

	if len(idParts) < 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: id:name1:name2:... Got: %q", req.ID),
		)
		return
	}

	for idx, v := range idParts {
		if idx == 0 {
			plan.ID = types.StringValue(v)
		} else {
			postgresqlDatabase := PostgresqlDatabase{
				DatabaseName: types.StringValue(v),
			}
			databaseList = append(databaseList, postgresqlDatabase)
		}
	}

	postgresqlDatabases, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: PostgresqlDatabase{}.AttrTypes()}, databaseList)
	plan.PostgresqlDatabaseList = postgresqlDatabases

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *postgresqlDatabasesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *postgresqlDatabasesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgresql_databases"
}

func (r *postgresqlDatabasesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"postgresql_database_list": schema.ListNestedAttribute{
				Required: true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Validators: []validator.String{
								stringvalidator.LengthBetween(4, 16),
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^[a-z]+[a-z0-9_]+$`),
									"Composed of lowercase alphabets, numbers, underbar (_). Must start with an alphabetic character.",
								),
							},
						},
						"owner": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(10),
				},
			},
		},
	}
}

func (r *postgresqlDatabasesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan postgresqlDatabasesResourceModel

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

	reqParams := &vpostgresql.AddCloudPostgresqlDatabaseListRequest{
		RegionCode:                  &r.config.RegionCode,
		CloudPostgresqlInstanceNo:   plan.ID.ValueStringPointer(),
		CloudPostgresqlDatabaseList: convertToCloudPostgresqlDatabaseParameters(plan.PostgresqlDatabaseList),
	}

	tflog.Info(ctx, "CreatePostgresqlDatabaseList reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vpostgresql.V2Api.AddCloudPostgresqlDatabaseList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "CreatePostgresqlDatabaseList response="+common.MarshalUncheckedString(response))

	if response == nil || *response.ReturnCode != "0" {
		resp.Diagnostics.AddError("CREATING ERROR", "response invalid")
		return
	}

	_, err = WaitPostgresqlCreation(ctx, r.config, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("WATING FOR POSTGRESQL CREATION ERROR", err.Error())
		return
	}

	output, err := GetPostgresqlDatabaseList(ctx, r.config, plan.ID.ValueString(), common.ConvertToStringList(plan.PostgresqlDatabaseList, "name"))
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if diags := plan.refreshFromOutput(ctx, output, plan.ID.ValueString()); diags.HasError() {
		resp.Diagnostics.AddError("READING ERROR", "refreshFromOutput error")
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *postgresqlDatabasesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state postgresqlDatabasesResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetPostgresqlDatabaseList(ctx, r.config, state.ID.ValueString(), common.ConvertToStringList(state.PostgresqlDatabaseList, "name"))
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if diags := state.refreshFromOutput(ctx, output, state.ID.ValueString()); diags.HasError() {
		resp.Diagnostics.AddError("READING ERROR", "refreshFromOutput error")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *postgresqlDatabasesResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *postgresqlDatabasesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state postgresqlDatabasesResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vpostgresql.DeleteCloudPostgresqlDatabaseListRequest{
		RegionCode:                  &r.config.RegionCode,
		CloudPostgresqlInstanceNo:   state.ID.ValueStringPointer(),
		CloudPostgresqlDatabaseList: convertToCloudPostgresqlDatabaseKeyParameter(state.PostgresqlDatabaseList),
	}
	tflog.Info(ctx, "DeletePostgresqlDatabseList reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vpostgresql.V2Api.DeleteCloudPostgresqlDatabaseList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("DELETING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "DeletePostgresqlDatabseList response="+common.MarshalUncheckedString(response))

	_, err = WaitPostgresqlCreation(ctx, r.config, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("WAITING FOR DELETE ERROR", err.Error())
		return
	}
}

func GetPostgresqlDatabaseList(ctx context.Context, config *conn.ProviderConfig, id string, dbs []string) ([]*vpostgresql.CloudPostgresqlDatabase, error) {
	reqParams := &vpostgresql.GetCloudPostgresqlDatabaseListRequest{
		RegionCode:                &config.RegionCode,
		CloudPostgresqlInstanceNo: ncloud.String(id),
	}
	tflog.Info(ctx, "GetPostgresqlDatabaseList reqParams="+common.MarshalUncheckedString(reqParams))

	resp, err := config.Client.Vpostgresql.V2Api.GetCloudPostgresqlDatabaseList(reqParams)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, nil
	}

	dbMap := make(map[string]*vpostgresql.CloudPostgresqlDatabase)
	for _, db := range resp.CloudPostgresqlDatabaseList {
		if db != nil && db.DatabaseName != nil {
			dbMap[*db.DatabaseName] = db
		}
	}

	var filteredDbs []*vpostgresql.CloudPostgresqlDatabase
	for _, dbname := range dbs {
		if db, exists := dbMap[dbname]; exists {
			filteredDbs = append(filteredDbs, db)
		}
	}

	if len(filteredDbs) == 0 {
		return nil, nil
	}

	tflog.Info(ctx, "GetPostgresqlDatabaseList response="+common.MarshalUncheckedString(resp))

	return filteredDbs, nil
}

type postgresqlDatabasesResourceModel struct {
	ID                     types.String `tfsdk:"id"`
	PostgresqlDatabaseList types.List   `tfsdk:"postgresql_database_list"`
}

type PostgresqlDatabase struct {
	DatabaseName types.String `tfsdk:"name"`
	Owner        types.String `tfsdk:"owner"`
}

func (r PostgresqlDatabase) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":  types.StringType,
		"owner": types.StringType,
	}
}

func (r *postgresqlDatabasesResourceModel) refreshFromOutput(ctx context.Context, output []*vpostgresql.CloudPostgresqlDatabase, instance string) diag.Diagnostics {
	r.ID = types.StringValue(instance)

	var databaseList []PostgresqlDatabase

	for _, database := range output {
		postgresqlDatabase := PostgresqlDatabase{
			DatabaseName: types.StringPointerValue(database.DatabaseName),
			Owner:        types.StringPointerValue(database.Owner),
		}

		databaseList = append(databaseList, postgresqlDatabase)
	}

	postgresqlDatabases, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: PostgresqlDatabase{}.AttrTypes()}, databaseList)
	if diags.HasError() {
		return diags
	}

	r.PostgresqlDatabaseList = postgresqlDatabases

	return nil
}

func convertToCloudPostgresqlDatabaseParameters(values basetypes.ListValue) []*vpostgresql.CloudPostgresqlDatabaseParameter {
	result := make([]*vpostgresql.CloudPostgresqlDatabaseParameter, 0, len(values.Elements()))

	for _, v := range values.Elements() {
		obj := v.(types.Object)
		attrs := obj.Attributes()

		param := &vpostgresql.CloudPostgresqlDatabaseParameter{
			Name:  attrs["name"].(types.String).ValueStringPointer(),
			Owner: attrs["owner"].(types.String).ValueStringPointer(),
		}
		result = append(result, param)
	}

	return result
}

func convertToCloudPostgresqlDatabaseKeyParameter(values basetypes.ListValue) []*vpostgresql.CloudPostgresqlDatabaseKeyParameter {
	result := make([]*vpostgresql.CloudPostgresqlDatabaseKeyParameter, 0, len(values.Elements()))

	for _, v := range values.Elements() {
		obj := v.(types.Object)
		attrs := obj.Attributes()

		param := &vpostgresql.CloudPostgresqlDatabaseKeyParameter{
			Name: attrs["name"].(types.String).ValueStringPointer(),
		}
		result = append(result, param)
	}

	return result
}
