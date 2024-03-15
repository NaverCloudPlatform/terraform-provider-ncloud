package mysql

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmysql"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/vpc"
)

var (
	_ resource.Resource                = &mysqlResource{}
	_ resource.ResourceWithConfigure   = &mysqlResource{}
	_ resource.ResourceWithImportState = &mysqlResource{}
)

func NewMysqlResource() resource.Resource {
	return &mysqlResource{}
}

type mysqlResource struct {
	config *conn.ProviderConfig
}

func (m *mysqlResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mysql"
}

func (m *mysqlResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 20),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[ㄱ-ㅣ가-힣A-Za-z0-9-]+$`),
						"Composed of alphabets, numbers, hyphen (-).",
					),
				},
			},
			"id": framework.IDAttribute(),
			"server_name_prefix": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 30),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-z]+[a-z0-9-]+[a-z0-9]$`),
						"Composed of lowercase alphabets, numbers, hyphen (-). Must start with an alphabetic character, and the last character can only be an English letter or number.",
					),
				},
			},
			"user_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(4, 16),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z]+[a-zA-Z0-9-\\_,]+$`),
						"Composed of alphabets, numbers, hyphen (-), (\\), (_), (,). Must start with an alphabetic character.",
					),
				},
			},
			"user_password": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.All(
						stringvalidator.LengthBetween(8, 20),
						stringvalidator.RegexMatches(regexp.MustCompile(`[a-zA-Z]+`), "Must have at least one alphabet"),
						stringvalidator.RegexMatches(regexp.MustCompile(`\d+`), "Must have at least one number"),
						stringvalidator.RegexMatches(regexp.MustCompile(`[~!@#$%^*()\-_=\[\]\{\};:,.<>?]+`), "Must have at least one special character"),
						stringvalidator.RegexMatches(regexp.MustCompile(`^[^&+\\"'/\s`+"`"+`]*$`), "Must not have ` & + \\ \" ' / and white space."),
					),
				},
				Sensitive: true,
			},
			"host_ip": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"database_name": schema.StringAttribute{
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
			"subnet_no": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"image_product_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"product_code": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"data_storage_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"SSD", "HDD", "CB1"}...),
				},
				Description: "default: SSD",
			},
			"is_ha": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Description: "default: true",
			},
			"is_multi_zone": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Description: "default: false",
			},
			"is_storage_encryption": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Description: "default: false",
			},
			"is_backup": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Description: "default: true",
			},
			"backup_file_retention_period": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.Between(1, 30),
				},
			},
			"backup_time": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Description: "ex) 01:15",
			},
			"is_automatic_backup": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"port": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.Any(
						int64validator.Between(10000, 20000),
						int64validator.OneOf(3306),
					),
				},
				Description: "default: 3306",
			},
			"standby_master_subnet_no": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"engine_version_code": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region_code": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vpc_no": schema.StringAttribute{
				Computed: true,
			},
			"access_control_group_no_list": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"mysql_config_list": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
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

func (r *mysqlResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*conn.ProviderConfig)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ProviderConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.config = config
}

func (r *mysqlResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan mysqlResourceModel

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

	subnet, err := vpc.GetSubnetInstance(r.config, plan.SubnetNo.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"CREATING ERROR",
			err.Error(),
		)
	}

	reqParams := &vmysql.CreateCloudMysqlInstanceRequest{
		RegionCode:                 &r.config.RegionCode,
		CloudMysqlServiceName:      plan.ServiceName.ValueStringPointer(),
		CloudMysqlServerNamePrefix: plan.ServerNamePrefix.ValueStringPointer(),
		CloudMysqlUserName:         plan.UserName.ValueStringPointer(),
		CloudMysqlUserPassword:     plan.UserPassword.ValueStringPointer(),
		HostIp:                     plan.HostIp.ValueStringPointer(),
		CloudMysqlDatabaseName:     plan.DatabaseName.ValueStringPointer(),
		VpcNo:                      subnet.VpcNo,
		SubnetNo:                   subnet.SubnetNo,
	}
	plan.VpcNo = types.StringPointerValue(subnet.VpcNo)

	if !plan.Port.IsNull() && !plan.Port.IsUnknown() {
		reqParams.CloudMysqlPort = ncloud.Int32(int32(plan.Port.ValueInt64()))
	}

	if !plan.DataStorageTypeCode.IsNull() {
		reqParams.DataStorageTypeCode = plan.DataStorageTypeCode.ValueStringPointer()
	}

	if !plan.ProductCode.IsNull() {
		reqParams.CloudMysqlProductCode = plan.ProductCode.ValueStringPointer()
	}

	if !plan.ImageProductCode.IsNull() {
		reqParams.CloudMysqlImageProductCode = plan.ImageProductCode.ValueStringPointer()
	}
	if !plan.IsHa.IsNull() && !plan.IsHa.IsUnknown() {
		reqParams.IsHa = plan.IsHa.ValueBoolPointer()
		if plan.IsHa.ValueBool() {
			if !plan.IsMultiZone.IsNull() && !plan.IsMultiZone.IsUnknown() {
				reqParams.IsMultiZone = plan.IsMultiZone.ValueBoolPointer()
			}
			if !plan.IsStorageEncryption.IsNull() && !plan.IsStorageEncryption.IsUnknown() {
				reqParams.IsStorageEncryption = plan.IsStorageEncryption.ValueBoolPointer()
			}
			if !plan.IsBackup.IsNull() && !plan.IsBackup.IsUnknown() && !plan.IsBackup.ValueBool() {
				resp.Diagnostics.AddError(
					"CREATING ERROR",
					"when `is_ha` is true, `is_backup` must be true or not be inputted",
				)
				return
			}

		} else {
			if !plan.IsMultiZone.IsNull() && !plan.IsMultiZone.IsUnknown() {
				resp.Diagnostics.AddError(
					"CREATING ERROR",
					"when `is_ha` is false, `is_multi_zone` parameter is not used",
				)
				return
			}
			if !plan.StandbyMasterSubnetNo.IsNull() && !plan.StandbyMasterSubnetNo.IsUnknown() {
				resp.Diagnostics.AddError(
					"CREATING ERROR",
					"when `is_ha` is false, `standby_master_subnet_no` is not used",
				)
				return
			}
			if !plan.IsStorageEncryption.IsNull() && !plan.IsStorageEncryption.IsUnknown() && plan.IsStorageEncryption.ValueBool() {
				resp.Diagnostics.AddError(
					"CREATING ERROR",
					"when `is_ha` is false, can't set true for `is_storage_encryption`",
				)
				return
			}
			if !plan.IsBackup.IsNull() && !plan.IsBackup.IsUnknown() {
				reqParams.IsBackup = plan.IsBackup.ValueBoolPointer()
			}
		}
	}

	if plan.IsMultiZone.ValueBool() {
		if plan.StandbyMasterSubnetNo.IsNull() || plan.StandbyMasterSubnetNo.IsUnknown() {
			resp.Diagnostics.AddError(
				"CREATING ERROR",
				"when `is_multi_zone` is true, `standby_master_subnet_no` must be entered",
			)
			return
		}
		reqParams.StandbyMasterSubnetNo = plan.StandbyMasterSubnetNo.ValueStringPointer()
	} else if !plan.StandbyMasterSubnetNo.IsNull() && !plan.StandbyMasterSubnetNo.IsUnknown() {
		resp.Diagnostics.AddError(
			"CREATING ERROR",
			"when `is_multi_zone` is false, `standby_master_subnet_no` is not used",
		)
		return
	}

	if !plan.BackupFileRetentionPeriod.IsNull() && !plan.BackupFileRetentionPeriod.IsUnknown() {
		reqParams.BackupFileRetentionPeriod = ncloud.Int32(int32(plan.BackupFileRetentionPeriod.ValueInt64()))
	}

	if !plan.IsAutomaticBackup.IsNull() && !plan.IsAutomaticBackup.IsUnknown() {
		reqParams.IsAutomaticBackup = plan.IsAutomaticBackup.ValueBoolPointer()
	}

	if reqParams.IsBackup == nil || *reqParams.IsBackup {
		if reqParams.IsAutomaticBackup == nil || *reqParams.IsAutomaticBackup {
			if !plan.BackupTime.IsNull() && !plan.BackupTime.IsUnknown() {
				resp.Diagnostics.AddError(
					"CREATING ERROR",
					"when `is_backup` is true and `is_automactic_backup` is true, `backup_time` is not used",
				)
				return
			}
		} else {
			if plan.BackupTime.IsNull() || plan.BackupTime.IsUnknown() {
				resp.Diagnostics.AddError(
					"CREATING ERROR",
					"when `is_backup` is true and `is_automactic_backup` is false, `backup_time` must be entered",
				)
				return
			}
			reqParams.BackupTime = plan.BackupTime.ValueStringPointer()
		}
	}

	tflog.Info(ctx, "CreateMysql reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vmysql.V2Api.CreateCloudMysqlInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "CreateMysql response="+common.MarshalUncheckedString(response))

	if response == nil || len(response.CloudMysqlInstanceList) < 1 {
		resp.Diagnostics.AddError("CREATING ERROR", "response invalid")
		return
	}

	mysqlIns := response.CloudMysqlInstanceList[0]
	plan.ID = types.StringPointerValue(mysqlIns.CloudMysqlInstanceNo)

	output, err := waitMysqlCreation(ctx, r.config, *mysqlIns.CloudMysqlInstanceNo)
	if err != nil {
		resp.Diagnostics.AddError("WAITING FOR CREATION ERROR", err.Error())
		return
	}

	plan.refreshFromOutput(ctx, output)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *mysqlResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state mysqlResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetMysqlInstance(ctx, r.config, state.ID.ValueString())
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

func (m *mysqlResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *mysqlResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state mysqlResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vmysql.DeleteCloudMysqlInstanceRequest{
		RegionCode:           &r.config.RegionCode,
		CloudMysqlInstanceNo: state.ID.ValueStringPointer(),
	}
	tflog.Info(ctx, "DeleteMysql reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := r.config.Client.Vmysql.V2Api.DeleteCloudMysqlInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("DELETING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "DeleteMysql response="+common.MarshalUncheckedString(response))

	if err := waitMysqlDeletion(ctx, r.config, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("WAITING FOR DELETE ERROR", err.Error())
	}
}

func (r *mysqlResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func GetMysqlInstance(ctx context.Context, config *conn.ProviderConfig, no string) (*vmysql.CloudMysqlInstance, error) {
	reqParams := &vmysql.GetCloudMysqlInstanceDetailRequest{
		RegionCode:           &config.RegionCode,
		CloudMysqlInstanceNo: ncloud.String(no),
	}
	tflog.Info(ctx, "GetMysqlDetail reqParams="+common.MarshalUncheckedString(reqParams))

	resp, err := config.Client.Vmysql.V2Api.GetCloudMysqlInstanceDetail(reqParams)
	// If the lookup result is 0 or already deleted, it will respond with a 400 error with a 5001017 return code.
	if err != nil && !(strings.Contains(err.Error(), `"returnCode": "5001017"`)) {
		return nil, err
	}
	tflog.Info(ctx, "GetMysqlDetail response="+common.MarshalUncheckedString(resp))

	if resp == nil || len(resp.CloudMysqlInstanceList) < 1 || len(resp.CloudMysqlInstanceList[0].CloudMysqlServerInstanceList) < 1 {
		return nil, nil
	}

	return resp.CloudMysqlInstanceList[0], nil
}

func waitMysqlCreation(ctx context.Context, config *conn.ProviderConfig, id string) (*vmysql.CloudMysqlInstance, error) {
	var mysqlInstance *vmysql.CloudMysqlInstance
	stateConf := &retry.StateChangeConf{
		Pending: []string{"creating", "settingUp"},
		Target:  []string{"running"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetMysqlInstance(ctx, config, id)
			mysqlInstance = instance
			if err != nil {
				return 0, "", err
			}

			status := instance.CloudMysqlInstanceStatus.Code
			op := instance.CloudMysqlInstanceOperation.Code

			if *status == "INIT" && *op == "CREAT" {
				return instance, "creating", nil
			}

			if *status == "CREAT" && *op == "SETUP" {
				return instance, "settingUp", nil
			}

			if *status == "CREAT" && *op == "NULL" {
				return instance, "running", nil
			}

			return 0, "", fmt.Errorf("error occurred while waiting to create mysql")
		},
		Timeout:    6 * conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err := stateConf.WaitForState()
	if err != nil {
		return nil, fmt.Errorf("error waiting for MysqlInstance state to be \"CREAT\": %s", err)
	}

	return mysqlInstance, nil
}

func waitMysqlDeletion(ctx context.Context, config *conn.ProviderConfig, id string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"deleting"},
		Target:  []string{"deleted"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetMysqlInstance(ctx, config, id)
			if err != nil {
				return 0, "", err
			}

			if instance == nil {
				return instance, "deleted", nil
			}

			status := instance.CloudMysqlInstanceStatus.Code
			op := instance.CloudMysqlInstanceOperation.Code

			if *status == "DEL" && *op == "DEL" {
				return instance, "deleting", nil
			}

			return 0, "", fmt.Errorf("error occurred while waiting to delete mysql")
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for mysql (%s) to become termintaing: %s", id, err)
	}

	return nil
}

type mysqlResourceModel struct {
	ID                        types.String `tfsdk:"id"`
	ServiceName               types.String `tfsdk:"service_name"`
	ServerNamePrefix          types.String `tfsdk:"server_name_prefix"`
	UserName                  types.String `tfsdk:"user_name"`
	UserPassword              types.String `tfsdk:"user_password"`
	HostIp                    types.String `tfsdk:"host_ip"`
	DatabaseName              types.String `tfsdk:"database_name"`
	SubnetNo                  types.String `tfsdk:"subnet_no"`
	ImageProductCode          types.String `tfsdk:"image_product_code"`
	ProductCode               types.String `tfsdk:"product_code"`
	DataStorageTypeCode       types.String `tfsdk:"data_storage_type"`
	IsHa                      types.Bool   `tfsdk:"is_ha"`
	IsMultiZone               types.Bool   `tfsdk:"is_multi_zone"`
	IsStorageEncryption       types.Bool   `tfsdk:"is_storage_encryption"`
	IsBackup                  types.Bool   `tfsdk:"is_backup"`
	BackupFileRetentionPeriod types.Int64  `tfsdk:"backup_file_retention_period"`
	BackupTime                types.String `tfsdk:"backup_time"`
	IsAutomaticBackup         types.Bool   `tfsdk:"is_automatic_backup"`
	Port                      types.Int64  `tfsdk:"port"`
	StandbyMasterSubnetNo     types.String `tfsdk:"standby_master_subnet_no"`
	EngineVersionCode         types.String `tfsdk:"engine_version_code"`
	RegionCode                types.String `tfsdk:"region_code"`
	VpcNo                     types.String `tfsdk:"vpc_no"`
	AccessControlGroupNoList  types.List   `tfsdk:"access_control_group_no_list"`
	MysqlConfigList           types.List   `tfsdk:"mysql_config_list"`
	MysqlServerList           types.List   `tfsdk:"mysql_server_list"`
}

type mysqlServer struct {
	ServerInstanceNo    types.String `tfsdk:"server_instance_no"`
	ServerName          types.String `tfsdk:"server_name"`
	ServerRole          types.String `tfsdk:"server_role"`
	ZoneCode            types.String `tfsdk:"zone_code"`
	SubnetNo            types.String `tfsdk:"subnet_no"`
	ProductCode         types.String `tfsdk:"product_code"`
	IsPublicSubnet      types.Bool   `tfsdk:"is_public_subnet"`
	PublicDomain        types.String `tfsdk:"public_domain"`
	PrivateDomain       types.String `tfsdk:"private_domain"`
	DataStorageSize     types.Int64  `tfsdk:"data_storage_size"`
	UsedDataStorageSize types.Int64  `tfsdk:"used_data_storage_size"`
	CpuCount            types.Int64  `tfsdk:"cpu_count"`
	MemorySize          types.Int64  `tfsdk:"memory_size"`
	Uptime              types.String `tfsdk:"uptime"`
	CreateDate          types.String `tfsdk:"create_date"`
}

func (m mysqlServer) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"server_instance_no":     types.StringType,
		"server_name":            types.StringType,
		"server_role":            types.StringType,
		"zone_code":              types.StringType,
		"subnet_no":              types.StringType,
		"product_code":           types.StringType,
		"is_public_subnet":       types.BoolType,
		"public_domain":          types.StringType,
		"private_domain":         types.StringType,
		"data_storage_size":      types.Int64Type,
		"used_data_storage_size": types.Int64Type,
		"cpu_count":              types.Int64Type,
		"memory_size":            types.Int64Type,
		"uptime":                 types.StringType,
		"create_date":            types.StringType,
	}
}

func (m *mysqlResourceModel) refreshFromOutput(ctx context.Context, output *vmysql.CloudMysqlInstance) {
	m.ID = types.StringPointerValue(output.CloudMysqlInstanceNo)
	m.ServiceName = types.StringPointerValue(output.CloudMysqlServiceName)
	m.ImageProductCode = types.StringPointerValue(output.CloudMysqlImageProductCode)
	m.DataStorageTypeCode = types.StringPointerValue(output.CloudMysqlServerInstanceList[0].DataStorageType.Code)
	m.IsHa = types.BoolPointerValue(output.IsHa)
	m.IsMultiZone = types.BoolPointerValue(output.IsMultiZone)
	m.IsStorageEncryption = types.BoolPointerValue(output.CloudMysqlServerInstanceList[0].IsStorageEncryption)
	m.IsBackup = types.BoolPointerValue(output.IsBackup)
	m.BackupFileRetentionPeriod = common.Int64ValueFromInt32(output.BackupFileRetentionPeriod)
	m.BackupTime = types.StringPointerValue(output.BackupTime)
	m.Port = common.Int64ValueFromInt32(output.CloudMysqlPort)
	m.EngineVersionCode = types.StringPointerValue(output.EngineVersion)
	m.RegionCode = types.StringPointerValue(output.CloudMysqlServerInstanceList[0].RegionCode)
	m.VpcNo = types.StringPointerValue(output.CloudMysqlServerInstanceList[0].VpcNo)

	acgList, _ := types.ListValueFrom(ctx, types.StringType, output.AccessControlGroupNoList)
	m.AccessControlGroupNoList = acgList
	configList, _ := types.ListValueFrom(ctx, types.StringType, output.CloudMysqlConfigList)
	m.MysqlConfigList = configList

	var serverList []mysqlServer
	for _, server := range output.CloudMysqlServerInstanceList {
		mysqlServerInstance := mysqlServer{
			ServerInstanceNo: types.StringPointerValue(server.CloudMysqlServerInstanceNo),
			ServerName:       types.StringPointerValue(server.CloudMysqlServerName),
			ServerRole:       types.StringPointerValue(server.CloudMysqlServerRole.Code),
			ZoneCode:         types.StringPointerValue(server.ZoneCode),
			SubnetNo:         types.StringPointerValue(server.SubnetNo),
			ProductCode:      types.StringPointerValue(server.CloudMysqlProductCode),
			IsPublicSubnet:   types.BoolPointerValue(server.IsPublicSubnet),
			PrivateDomain:    types.StringPointerValue(server.PrivateDomain),
			DataStorageSize:  types.Int64PointerValue(server.DataStorageSize),
			CpuCount:         common.Int64ValueFromInt32(server.CpuCount),
			MemorySize:       types.Int64PointerValue(server.MemorySize),
			Uptime:           types.StringPointerValue(server.Uptime),
			CreateDate:       types.StringPointerValue(server.CreateDate),
		}

		if server.PublicDomain != nil {
			mysqlServerInstance.PublicDomain = types.StringPointerValue(server.PublicDomain)
		}

		if server.UsedDataStorageSize != nil {
			mysqlServerInstance.UsedDataStorageSize = types.Int64PointerValue(server.UsedDataStorageSize)
		}
		serverList = append(serverList, mysqlServerInstance)
	}

	mysqlServers, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: mysqlServer{}.attrTypes()}, serverList)

	m.MysqlServerList = mysqlServers
}
