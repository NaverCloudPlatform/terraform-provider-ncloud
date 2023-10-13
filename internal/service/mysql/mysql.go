package mysql

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sdkresource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/vpc"
	"regexp"
	"strings"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmysql"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/framework"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
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

func (r *mysqlResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
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
					stringvalidator.LengthAtLeast(3),
					stringvalidator.LengthAtMost(20),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[ㄱ-ㅣ가-힣A-Za-z0-9-]+$`),
						"Composed of alphabets, numbers, hyphen (-).",
					),
				},
			},
			"id": framework.IDAttribute(),
			"name_prefix": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(3),
					stringvalidator.LengthAtMost(30),
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
					stringvalidator.LengthAtLeast(4),
					stringvalidator.LengthAtMost(16),
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
					stringvalidator.LengthAtLeast(8),
					stringvalidator.LengthAtMost(20),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9~!@#$%^*()\-_=\[\]\{\};:,.<>?]{8,20}$`),
						"Must Combine at least one each of alphabets, numbers, special characters except ` & + \\ \" ' / and white space.",
					),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`.+[a-zA-Z]{1,}.+|.+[a-zA-Z]{1,}|[a-zA-Z]{1,}.+`),
						"Must have at least 1 alphabet.",
					),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`.+[0-9]{1,}.+|.+[0-9]{1,}|[0-9]{1,}.+`),
						"Must have at least 1 Number.",
					),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`.+[~!@#$%^*()\-_=\[\]\{\};:,.<>?].+|.+[~!@#$%^*()\-_=\[\]\{\};:,.<>?]|[~!@#$%^*()\-_=\[\]\{\};:,.<>?].+`),
						"Must have at least 1 special characters except ` & + \\ \" ' / and white space.",
					),
				},
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
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(30),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z]+[a-zA-Z0-9-\\_,]+$`),
						"Composed of alphabets, numbers, hyphen (-), (\\), (_), (,). Must start with an alphabetic character.",
					),
				},
			},
			"subnet_no": schema.StringAttribute{
				Required: true,
			},
			"engine_version_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"data_storage_type_code": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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
					int64planmodifier.RequiresReplace(),
				},
				Description: "default: false",
			},
			"backup_time": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
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
					int64planmodifier.RequiresReplace(),
				},
				Description: "default: 3306",
			},
			"standby_master_subnet_no": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vpc_no": schema.StringAttribute{
				Computed: true,
			},
			"image_product_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"product_code": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"instance_no": schema.StringAttribute{
				Computed: true,
			},

			"access_control_group_no_list": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},

			"mysql_config_list": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"create_date": schema.StringAttribute{
				Computed: true,
			},
			"mysql_server_list": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"mysql_server_instance_no": schema.StringAttribute{
							Computed: true,
						},
						"mysql_server_name": schema.StringAttribute{
							Computed: true,
						},
						"mysql_server_role": schema.StringAttribute{
							Computed: true,
						},
						"mysql_server_product_code": schema.StringAttribute{
							Computed: true,
						},
						"region_code": schema.StringAttribute{
							Computed: true,
						},
						"zone_code": schema.StringAttribute{
							Computed: true,
						},
						"vpc_no": schema.StringAttribute{
							Computed: true,
						},
						"subnet_no": schema.StringAttribute{
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
						"data_storage_type": schema.StringAttribute{
							Computed: true,
						},
						"is_storage_encryption": schema.BoolAttribute{
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
				Computed: true,
			},
		},
	}
}

func (r *mysqlResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan mysqlResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !r.config.SupportVPC {
		resp.Diagnostics.AddError(
			"Not support classic",
			fmt.Sprintf("resource %s does not support classic", req.Config.Schema.Type().String()),
		)
		return
	}

	subnet, err := vpc.GetSubnetInstance(r.config, plan.SubnetNo.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Get Subnet Instance, SubnetNo=%v", plan.SubnetNo.ValueString()),
			err.Error(),
		)
	}

	reqParams := &vmysql.CreateCloudMysqlInstanceRequest{
		RegionCode: &r.config.RegionCode,
		VpcNo:      subnet.VpcNo,
		SubnetNo:   subnet.SubnetNo,
	}

	if !plan.DatabaseName.IsNull() {
		reqParams.CloudMysqlDatabaseName = plan.DatabaseName.ValueStringPointer()
	}

	if !plan.Port.IsNull() && !plan.Port.IsUnknown() {
		reqParams.CloudMysqlPort = ncloud.Int32(int32(plan.Port.ValueInt64()))
	}

	if !plan.HostIp.IsNull() {
		reqParams.HostIp = plan.HostIp.ValueStringPointer()
	}

	if !plan.UserPassword.IsNull() {
		reqParams.CloudMysqlUserPassword = plan.UserPassword.ValueStringPointer()
	}

	if !plan.UserName.IsNull() {
		reqParams.CloudMysqlUserName = plan.UserName.ValueStringPointer()
	}

	if !plan.NamePrefix.IsNull() {
		reqParams.CloudMysqlServerNamePrefix = plan.NamePrefix.ValueStringPointer()
	}

	if !plan.ServiceName.IsNull() {
		reqParams.CloudMysqlServiceName = plan.ServiceName.ValueStringPointer()
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
			if !plan.IsStorageEncryption.IsNull() {
				reqParams.IsStorageEncryption = plan.IsStorageEncryption.ValueBoolPointer()
			}
			if !plan.IsBackup.IsNull() && !plan.IsBackup.ValueBool() {
				resp.Diagnostics.AddError(
					fmt.Sprintf("when `is_ha` is true, `is_backup` must be true or not be inputted for default value"),
					err.Error(),
				)
				return
			}

		} else {
			if !plan.IsMultiZone.IsNull() && plan.IsMultiZone.ValueBool() {
				resp.Diagnostics.AddError(
					fmt.Sprintf("when `is_ha` is false, `is_backup` must be false"),
					err.Error(),
				)
				return
			}
			if !plan.StandbyMasterSubnetNo.IsNull() {
				resp.Diagnostics.AddError(
					fmt.Sprintf("when `is_ha` is false, `standby_master_subnet_no` must not be inputed"),
					err.Error(),
				)
				return
			}
			if !plan.IsBackup.IsNull() && !plan.IsBackup.IsUnknown() {
				reqParams.IsBackup = plan.IsBackup.ValueBoolPointer()
			}
		}
	}

	if !plan.IsMultiZone.IsNull() && plan.IsMultiZone.ValueBool() {
		if plan.StandbyMasterSubnetNo.IsNull() {
			resp.Diagnostics.AddError(
				fmt.Sprintf("when `is_multi_zone` is true, `standby_master_subnet_no` must be entered"),
				err.Error(),
			)
			return
		}
		reqParams.StandbyMasterSubnetNo = plan.StandbyMasterSubnetNo.ValueStringPointer()
	}

	if !plan.BackupFileRetentionPeriod.IsNull() && !plan.BackupFileRetentionPeriod.IsUnknown() {
		reqParams.BackupFileRetentionPeriod = ncloud.Int32(int32(plan.BackupFileRetentionPeriod.ValueInt64()))
	}

	if !plan.IsAutomaticBackup.IsNull() {
		reqParams.IsAutomaticBackup = plan.IsAutomaticBackup.ValueBoolPointer()
	}

	if (reqParams.IsBackup == nil || *reqParams.IsBackup) && !plan.IsAutomaticBackup.IsNull() && !plan.IsAutomaticBackup.ValueBool() {
		if plan.BackupTime.IsNull() {
			resp.Diagnostics.AddError(
				fmt.Sprintf("when `is_backup` is true and `is_automactic_backup` is false, `backup_time` must be entered"),
				err.Error(),
			)
			return
		}
		reqParams.BackupTime = plan.BackupTime.ValueStringPointer()
	}

	tflog.Info(ctx, "CreateMysql", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})
	response, err := r.config.Client.Vmysql.V2Api.CreateCloudMysqlInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Create Mysql Instance, err params=%v", *reqParams),
			err.Error(),
		)
		return
	}
	tflog.Info(ctx, "CreateMysql response", map[string]any{
		"createMysqlResponse": common.MarshalUncheckedString(response),
	})

	mysqlIns := response.CloudMysqlInstanceList[0]
	plan.ID = types.StringPointerValue(mysqlIns.CloudMysqlInstanceNo)
	tflog.Info(ctx, "Mysql ID", map[string]any{
		"MysqlNo": *mysqlIns.CloudMysqlInstanceNo,
	})
	output, err := waitMysqlForCreation(ctx, r.config, *mysqlIns.CloudMysqlInstanceNo)
	if err != nil {
		resp.Diagnostics.AddError("waiting for Mysql creation", err.Error())
		return
	}

	if err := plan.refreshFromOutput(ctx, output); err != nil {
		resp.Diagnostics.AddError("refreshing mysql details", err.Error())
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

}

func (r *mysqlResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state mysqlResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetMysqlInstance(ctx, r.config, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("GetMysql", err.Error())
		return
	}

	if output == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if err := state.refreshFromOutput(ctx, output); err != nil {
		resp.Diagnostics.AddError("refreshing mysql details", err.Error())
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (m *mysqlResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *mysqlResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state mysqlResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vmysql.DeleteCloudMysqlInstanceRequest{
		RegionCode:           &r.config.RegionCode,
		CloudMysqlInstanceNo: state.InstanceNo.ValueStringPointer(),
	}

	tflog.Info(ctx, "DeleteMysql", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})

	response, err := r.config.Client.Vmysql.V2Api.DeleteCloudMysqlInstance(reqParams)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Delete Mysql Instance params=%v", *reqParams),
			err.Error(),
		)
		return
	}

	tflog.Info(ctx, "DeleteCloudMysql response", map[string]any{
		"deleteMysqlResponse": common.MarshalUncheckedString(response),
	})

	if err := WaitForMysqlDeletion(ctx, r.config, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"fail to wait for mysql deletion",
			err.Error(),
		)
	}
}

func GetMysqlInstance(ctx context.Context, config *conn.ProviderConfig, id string) (*vmysql.CloudMysqlInstance, error) {
	reqParams := &vmysql.GetCloudMysqlInstanceDetailRequest{
		RegionCode:           &config.RegionCode,
		CloudMysqlInstanceNo: ncloud.String(id),
	}

	tflog.Info(ctx, "GetMysql", map[string]any{
		"reqParams": common.MarshalUncheckedString(reqParams),
	})

	resp, err := config.Client.Vmysql.V2Api.GetCloudMysqlInstanceDetail(reqParams)
	if err != nil {
		return nil, err
	}

	tflog.Info(ctx, "GetMysql response", map[string]any{
		"getMysqlResponse": common.MarshalUncheckedString(resp),
	})

	if len(resp.CloudMysqlInstanceList) > 0 {
		mysql := resp.CloudMysqlInstanceList[0]
		return mysql, nil
	}
	return nil, nil
}
func WaitForMysqlDeletion(ctx context.Context, config *conn.ProviderConfig, id string) error {
	stateConf := &sdkresource.StateChangeConf{
		Pending: []string{"deleting"},
		Target:  []string{"deleted"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetMysqlInstance(ctx, config, id)

			if err != nil && !strings.Contains(err.Error(), `"returnCode": "5001017"`) {
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
		return fmt.Errorf("Error waiting for mysql (%s) to become termintaing: %s", id, err)
	}

	return nil
}

func waitMysqlForCreation(ctx context.Context, config *conn.ProviderConfig, id string) (*vmysql.CloudMysqlInstance, error) {
	var mysqlInstance *vmysql.CloudMysqlInstance
	stateConf := &sdkresource.StateChangeConf{
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

type mysqlResourceModel struct {
	ServiceName               types.String `tfsdk:"service_name"`
	NamePrefix                types.String `tfsdk:"name_prefix"`
	UserName                  types.String `tfsdk:"user_name"`
	UserPassword              types.String `tfsdk:"user_password"`
	HostIp                    types.String `tfsdk:"host_ip"`
	DatabaseName              types.String `tfsdk:"database_name"`
	SubnetNo                  types.String `tfsdk:"subnet_no"`
	EngineVersionCode         types.String `tfsdk:"engine_version_code"`
	DataStorageTypeCode       types.String `tfsdk:"data_storage_type_code"`
	IsHa                      types.Bool   `tfsdk:"is_ha"`
	IsMultiZone               types.Bool   `tfsdk:"is_multi_zone"`
	IsStorageEncryption       types.Bool   `tfsdk:"is_storage_encryption"`
	IsBackup                  types.Bool   `tfsdk:"is_backup"`
	BackupFileRetentionPeriod types.Int64  `tfsdk:"backup_file_retention_period"`
	BackupTime                types.String `tfsdk:"backup_time"`
	IsAutomaticBackup         types.Bool   `tfsdk:"is_automatic_backup"`
	Port                      types.Int64  `tfsdk:"port"`
	StandbyMasterSubnetNo     types.String `tfsdk:"standby_master_subnet_no"`
	VpcNo                     types.String `tfsdk:"vpc_no"`
	ImageProductCode          types.String `tfsdk:"image_product_code"`
	ProductCode               types.String `tfsdk:"product_code"`
	InstanceNo                types.String `tfsdk:"instance_no"`
	ID                        types.String `tfsdk:"id"`
	CreateDate                types.String `tfsdk:"create_date"`
	AccessControlGroupNoList  types.List   `tfsdk:"access_control_group_no_list"`
	MysqlConfigList           types.List   `tfsdk:"mysql_config_list"`
	MysqlServerList           types.List   `tfsdk:"mysql_server_list"`
}

type mysqlServer struct {
	MysqlServerInstanceNo  types.String `tfsdk:"mysql_server_instance_no"`
	MysqlServerName        types.String `tfsdk:"mysql_server_name"`
	MysqlServerRole        types.String `tfsdk:"mysql_server_role"`
	MysqlServerProductCode types.String `tfsdk:"mysql_server_product_code"`
	RegionCode             types.String `tfsdk:"region_code"`
	ZoneCode               types.String `tfsdk:"zone_code"`
	VpcNo                  types.String `tfsdk:"vpc_no"`
	SubnetNo               types.String `tfsdk:"subnet_no"`
	IsPublicSubnet         types.Bool   `tfsdk:"is_public_subnet"`
	PublicDomain           types.String `tfsdk:"public_domain"`
	PrivateDomain          types.String `tfsdk:"private_domain"`
	DataStorageType        types.String `tfsdk:"data_storage_type"`
	IsStorageEncryption    types.Bool   `tfsdk:"is_storage_encryption"`
	DataStorageSize        types.Int64  `tfsdk:"data_storage_size"`
	UsedDataStorageSize    types.Int64  `tfsdk:"used_data_storage_size"`
	CpuCount               types.Int64  `tfsdk:"cpu_count"`
	MemorySize             types.Int64  `tfsdk:"memory_size"`
	Uptime                 types.String `tfsdk:"uptime"`
	CreateDate             types.String `tfsdk:"create_date"`
}

func (m mysqlServer) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"mysql_server_instance_no":  types.StringType,
		"mysql_server_name":         types.StringType,
		"mysql_server_role":         types.StringType,
		"mysql_server_product_code": types.StringType,
		"region_code":               types.StringType,
		"zone_code":                 types.StringType,
		"vpc_no":                    types.StringType,
		"subnet_no":                 types.StringType,
		"is_public_subnet":          types.BoolType,
		"public_domain":             types.StringType,
		"private_domain":            types.StringType,
		"data_storage_type":         types.StringType,
		"is_storage_encryption":     types.BoolType,
		"data_storage_size":         types.Int64Type,
		"used_data_storage_size":    types.Int64Type,
		"cpu_count":                 types.Int64Type,
		"memory_size":               types.Int64Type,
		"uptime":                    types.StringType,
		"create_date":               types.StringType,
	}
}

func (m *mysqlResourceModel) refreshFromOutput(ctx context.Context, output *vmysql.CloudMysqlInstance) error {
	m.ID = types.StringPointerValue(output.CloudMysqlInstanceNo)
	m.ServiceName = types.StringPointerValue(output.CloudMysqlServiceName)
	m.EngineVersionCode = types.StringPointerValue(output.EngineVersion)
	m.IsHa = types.BoolPointerValue(output.IsHa)
	m.IsMultiZone = types.BoolPointerValue(output.IsMultiZone)
	m.IsBackup = types.BoolPointerValue(output.IsBackup)
	m.BackupFileRetentionPeriod = types.Int64Value(int64(*output.BackupFileRetentionPeriod))
	m.BackupTime = types.StringPointerValue(output.BackupTime)
	m.Port = types.Int64Value(int64(*output.CloudMysqlPort))
	m.ImageProductCode = types.StringPointerValue(output.CloudMysqlImageProductCode)
	m.CreateDate = types.StringPointerValue(output.CreateDate)
	m.InstanceNo = types.StringPointerValue(output.CloudMysqlInstanceNo)
	m.VpcNo = types.StringPointerValue(output.CloudMysqlServerInstanceList[0].VpcNo)
	m.IsStorageEncryption = types.BoolPointerValue(output.CloudMysqlServerInstanceList[0].IsStorageEncryption)

	acgList, _ := types.ListValueFrom(ctx, types.StringType, output.AccessControlGroupNoList)
	m.AccessControlGroupNoList = acgList

	configList, _ := types.ListValueFrom(ctx, types.StringType, output.CloudMysqlConfigList)
	m.MysqlConfigList = configList

	var serverList []mysqlServer
	for _, server := range output.CloudMysqlServerInstanceList {
		mysqlServerInstance := mysqlServer{
			MysqlServerInstanceNo:  types.StringPointerValue(server.CloudMysqlServerInstanceNo),
			MysqlServerName:        types.StringPointerValue(server.CloudMysqlServerName),
			MysqlServerRole:        types.StringPointerValue(server.CloudMysqlServerRole.Code),
			MysqlServerProductCode: types.StringPointerValue(server.CloudMysqlProductCode),
			RegionCode:             types.StringPointerValue(server.RegionCode),
			ZoneCode:               types.StringPointerValue(server.ZoneCode),
			VpcNo:                  types.StringPointerValue(server.VpcNo),
			SubnetNo:               types.StringPointerValue(server.SubnetNo),
			IsPublicSubnet:         types.BoolPointerValue(server.IsPublicSubnet),
			PrivateDomain:          types.StringPointerValue(server.PrivateDomain),
			DataStorageType:        types.StringPointerValue(server.DataStorageType.Code),
			IsStorageEncryption:    types.BoolPointerValue(server.IsStorageEncryption),
			DataStorageSize:        types.Int64Value(*server.DataStorageSize),
			CpuCount:               types.Int64Value(int64(*server.CpuCount)),
			MemorySize:             types.Int64Value(*server.MemorySize),
			Uptime:                 types.StringPointerValue(server.Uptime),
			CreateDate:             types.StringPointerValue(server.CreateDate),
		}
		if server.PublicDomain != nil {
			mysqlServerInstance.PublicDomain = types.StringPointerValue(server.PublicDomain)
		}

		if server.UsedDataStorageSize != nil {
			mysqlServerInstance.UsedDataStorageSize = types.Int64Value(*server.UsedDataStorageSize)
		}
		serverList = append(serverList, mysqlServerInstance)
	}

	mysqlServers, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: mysqlServer{}.attrTypes()}, serverList)

	m.MysqlServerList = mysqlServers
	return nil
}
