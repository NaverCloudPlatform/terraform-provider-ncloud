package mysql

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmysql"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
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
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
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

	plan.ID = plan.MysqlInstanceNo

	tflog.Info(ctx, "CreateMysqlDatabaseList reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vmysql.V2Api.AddCloudMysqlDatabaseList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR-dbs", err.Error())
		return
	}
	tflog.Info(ctx, "CreateMysqlDatabaseList response="+common.MarshalUncheckedString(response))

	if response == nil || *response.ReturnCode != "0" {
		resp.Diagnostics.AddError("CREATING ERROR-dbs", "response invalid")
		return
	}

	_, err = waitMysqlCreation(ctx, r.config, plan.MysqlInstanceNo.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("WAITING FOR MYSQL CREATION ERROR", err.Error())
		return
	}

	output, err := GetMysqlDatabaseList(ctx, r.config, plan.MysqlInstanceNo.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	plan.refreshFromOutput(ctx, output, plan.MysqlInstanceNo.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *mysqlDatabasesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state mysqlDatabasesResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetMysqlDatabaseList(ctx, r.config, state.MysqlInstanceNo.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.refreshFromOutput(ctx, output, state.MysqlInstanceNo.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *mysqlDatabasesResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *mysqlDatabasesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state mysqlDatabasesResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fmt.Println(len(state.MysqlDatabaseList.Elements()))
	_, err := waitMysqlCreation(ctx, r.config, state.MysqlInstanceNo.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("WAITING FOR MYSQAL CREATION ERROR", err.Error())
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

	if err := waitMysqlDatabasesDeletion(ctx, r.config, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("WAITING FOR DELETE ERROR", err.Error())
	}
}

func waitMysqlDatabasesDeletion(ctx context.Context, config *conn.ProviderConfig, id string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{DELETING},
		Target:  []string{DELETED},
		Refresh: func() (interface{}, string, error) {
			dbList, err := GetMysqlDatabaseList(ctx, config, id)
			if err != nil {
				return 0, "", err
			}

			if len(dbList) > 0 {
				return dbList, DELETING, nil
			}

			if len(dbList) == 0 || dbList == nil {
				return dbList, DELETED, nil
			}

			return 0, "", fmt.Errorf("error occurred while waiting to delete postgresql database")
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      1 * time.Minute,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for postgresql database (%s) to become terminating: %s", id, err)
	}

	return nil
}

func GetMysqlDatabaseList(ctx context.Context, config *conn.ProviderConfig, id string) ([]*vmysql.CloudMysqlDatabase, error) {
	reqParams := &vmysql.GetCloudMysqlDatabaseListRequest{
		RegionCode:           &config.RegionCode,
		CloudMysqlInstanceNo: ncloud.String(id),
		PageNo:               ncloud.Int32(0),
		PageSize:             ncloud.Int32(2147483647),
	}
	tflog.Info(ctx, "GetMysqlDatabaseList reqParams="+common.MarshalUncheckedString(reqParams))

	resp, err := config.Client.Vmysql.V2Api.GetCloudMysqlDatabaseList(reqParams)
	if err != nil {
		return nil, err
	}
	tflog.Info(ctx, "GetMysqlDatabaseList response="+common.MarshalUncheckedString(resp))

	fmt.Println("database list:")
	for _, i := range resp.CloudMysqlDatabaseList {
		fmt.Println(*i.DatabaseName)
	}

	if resp == nil || len(resp.CloudMysqlDatabaseList) == 1 {
		return nil, nil
	}

	tflog.Info(ctx, "GetMysqlUserList response="+common.MarshalUncheckedString(resp))

	return resp.CloudMysqlDatabaseList[1:], nil
}

type mysqlDatabasesResourceModel struct {
	ID                types.String `tfsdk:"id"`
	MysqlInstanceNo   types.String `tfsdk:"mysql_instance_no"`
	MysqlDatabaseList types.List   `tfsdk:"mysql_database_list"`
}

type mysqlDatabase struct {
	DatabaseName types.String `tfsdk:"name"`
}

func (r mysqlDatabase) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
	}
}

func (r *mysqlDatabasesResourceModel) refreshFromOutput(ctx context.Context, output []*vmysql.CloudMysqlDatabase, instance string) {
	r.ID = types.StringValue(instance)
	r.MysqlInstanceNo = types.StringValue(instance)

	var databaseList []mysqlDatabase
	for _, db := range output {
		mysqlDb := mysqlDatabase{
			DatabaseName: types.StringPointerValue(db.DatabaseName),
		}
		fmt.Printf("refresh: %s", *db.DatabaseName)
		databaseList = append(databaseList, mysqlDb)
	}

	mysqlDatabases, err := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: mysqlDatabase{}.AttrTypes()}, databaseList)
	if err != nil {
		log.Printf("Error converting database list: %v", err)
		return
	}

	r.MysqlDatabaseList = mysqlDatabases
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
