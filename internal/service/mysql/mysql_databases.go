package mysql

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmysql"
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
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"
)

var (
	_ resource.Resource                = &mysqlDatabasesResource{}
	_ resource.ResourceWithConfigure   = &mysqlDatabasesResource{}
	_ resource.ResourceWithImportState = &mysqlDatabasesResource{}
)

func NewMysqlDatabasesResource() resource.Resource {
	return &mysqlDatabasesResource{}
}

type mysqlDatabasesResource struct {
	config *conn.ProviderConfig
}

func (r *mysqlDatabasesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var plan mysqlDatabasesResourceModel
	var dbList []mysqlDatabase
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
			plan.MysqlInstanceNo = types.StringValue(v)
		} else {
			db := mysqlDatabase{
				DatabaseName: types.StringValue(v),
			}
			dbList = append(dbList, db)
		}
	}

	mysqlDatabases, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: mysqlDatabase{}.attrTypes()}, dbList)
	plan.MysqlDatabaseList = mysqlDatabases

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *mysqlDatabasesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *mysqlDatabasesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mysql_databases"
}

func (r *mysqlDatabasesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttribute(),
			"mysql_instance_no": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"mysql_database_list": schema.ListNestedAttribute{
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
								stringvalidator.LengthBetween(1, 30),
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^[a-zA-Z]+[a-zA-Z0-9-\\_,]+$`),
									"Composed of alphabets, numbers, hyphen (-), (\\), (_), (,). Must start with an alphabetic character.",
								),
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

func (r *mysqlDatabasesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan mysqlDatabasesResourceModel

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

	reqParams := &vmysql.AddCloudMysqlDatabaseListRequest{
		RegionCode:                 &r.config.RegionCode,
		CloudMysqlInstanceNo:       plan.MysqlInstanceNo.ValueStringPointer(),
		CloudMysqlDatabaseNameList: convertToStringList(plan.MysqlDatabaseList),
	}

	tflog.Info(ctx, "CreateMysqlDatabaseList reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vmysql.V2Api.AddCloudMysqlDatabaseList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "CreateMysqlDatabaseList response="+common.MarshalUncheckedString(response))

	if response == nil || *response.ReturnCode != "0" {
		resp.Diagnostics.AddError("CREATING ERROR", "response invalid")
		return
	}

	_, err = waitMysqlCreation(ctx, r.config, plan.MysqlInstanceNo.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("WAITING FOR MYSQL CREATION ERROR", err.Error())
		return
	}

	output, err := GetMysqlDatabaseList(ctx, r.config, plan.MysqlInstanceNo.ValueString(), common.ConvertToStringList(plan.MysqlDatabaseList, "name"))
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if diags := plan.refreshFromOutput(ctx, output, plan.MysqlInstanceNo.ValueString()); diags.HasError() {
		resp.Diagnostics.AddError("READING ERROR", "refreshFromOutput error")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *mysqlDatabasesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state mysqlDatabasesResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetMysqlDatabaseList(ctx, r.config, state.MysqlInstanceNo.ValueString(), common.ConvertToStringList(state.MysqlDatabaseList, "name"))
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if diags := state.refreshFromOutput(ctx, output, state.MysqlInstanceNo.ValueString()); diags.HasError() {
		resp.Diagnostics.AddError("READING ERROR", "refreshFromOutput error")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *mysqlDatabasesResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *mysqlDatabasesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state mysqlDatabasesResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := waitMysqlCreation(ctx, r.config, state.MysqlInstanceNo.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("WAITING FOR MYSQL DELETE ERROR", err.Error())
		return
	}

	reqParams := &vmysql.DeleteCloudMysqlDatabaseListRequest{
		RegionCode:                 &r.config.RegionCode,
		CloudMysqlInstanceNo:       state.MysqlInstanceNo.ValueStringPointer(),
		CloudMysqlDatabaseNameList: convertToStringList(state.MysqlDatabaseList),
	}
	tflog.Info(ctx, "DeleteMysqlDatabaseList reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vmysql.V2Api.DeleteCloudMysqlDatabaseList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("DELETING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "DeleteMysqlDatabaseList response="+common.MarshalUncheckedString(response))
}

func GetMysqlDatabaseList(ctx context.Context, config *conn.ProviderConfig, id string, dbs []string) ([]*vmysql.CloudMysqlDatabase, error) {
	var allDbs []*vmysql.CloudMysqlDatabase
	pageNo := int32(0)
	pageSize := int32(100)
	hasMore := true

	for hasMore {
		reqParams := &vmysql.GetCloudMysqlDatabaseListRequest{
			RegionCode:           &config.RegionCode,
			CloudMysqlInstanceNo: ncloud.String(id),
			PageNo:               ncloud.Int32(pageNo),
			PageSize:             ncloud.Int32(pageSize),
		}
		tflog.Info(ctx, "GetMysqlDatabaseList reqParams="+common.MarshalUncheckedString(reqParams))

		resp, err := config.Client.Vmysql.V2Api.GetCloudMysqlDatabaseList(reqParams)
		if err != nil {
			return nil, err
		}

		if resp == nil {
			break
		}

		allDbs = append(allDbs, resp.CloudMysqlDatabaseList...)

		hasMore = len(resp.CloudMysqlDatabaseList) == int(pageSize)
		pageNo++
	}

	dbMap := make(map[string]*vmysql.CloudMysqlDatabase)
	for _, db := range allDbs {
		if db != nil && db.DatabaseName != nil {
			dbMap[*db.DatabaseName] = db
		}
	}

	var filteredDbs []*vmysql.CloudMysqlDatabase
	for _, dbname := range dbs {
		if db, exists := dbMap[dbname]; exists {
			filteredDbs = append(filteredDbs, db)
		}
	}

	if len(filteredDbs) == 0 {
		return nil, nil
	}

	tflog.Info(ctx, "GetMysqlDatabseList response="+common.MarshalUncheckedString(filteredDbs))

	return filteredDbs, nil
}

type mysqlDatabasesResourceModel struct {
	ID                types.String `tfsdk:"id"`
	MysqlInstanceNo   types.String `tfsdk:"mysql_instance_no"`
	MysqlDatabaseList types.List   `tfsdk:"mysql_database_list"`
}

type mysqlDatabase struct {
	DatabaseName types.String `tfsdk:"name"`
}

func (r mysqlDatabase) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
	}
}

func (r *mysqlDatabasesResourceModel) refreshFromOutput(ctx context.Context, output []*vmysql.CloudMysqlDatabase, instanceNo string) diag.Diagnostics {
	r.ID = types.StringValue(instanceNo)
	r.MysqlInstanceNo = types.StringValue(instanceNo)

	var databaseList []mysqlDatabase
	for _, db := range output {
		mysqlDb := mysqlDatabase{
			DatabaseName: types.StringPointerValue(db.DatabaseName),
		}
		databaseList = append(databaseList, mysqlDb)
	}

	mysqlDatabases, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: mysqlDatabase{}.attrTypes()}, databaseList)
	if diags.HasError() {
		return diags
	}

	r.MysqlDatabaseList = mysqlDatabases

	return diags
}

func convertToStringList(values basetypes.ListValue) []*string {
	result := make([]*string, 0, len(values.Elements()))

	for _, v := range values.Elements() {
		obj := v.(types.Object)

		attrs := obj.Attributes()

		name := attrs["name"].(types.String).ValueStringPointer()
		result = append(result, name)
	}

	return result
}
