package autoscaling

import (
	"context"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/autoscaling"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vautoscaling"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

var (
	_ datasource.DataSource              = &launchConfigurationDataSource{}
	_ datasource.DataSourceWithConfigure = &launchConfigurationDataSource{}
)

func NewLaunchConfigDataSource() datasource.DataSource {
	return &launchConfigurationDataSource{}
}

type launchConfigurationDataSource struct {
	config *conn.ProviderConfig
}

func (l *launchConfigurationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	l.config = config
}

func (l *launchConfigurationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_launch_configuration"
}

func (l *launchConfigurationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"filter": common.DataSourceFiltersBlock(),
		},
	}
}

func (l *launchConfigurationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data launchConfigurationDataSourceModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	launchConfigResp, err := getLaunchConfigurationList(ctx, l.config, data.LaunchConfigurationNo.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Error Reading LaunchConfigurationList", err.Error())
		return
	}

	if launchConfigResp == nil {
		return
	}

	tflog.Info(ctx, "LaunchConfigurationList response="+common.MarshalUncheckedString(launchConfigResp))

	launchConfigList, diags := flattenLaunchConfigurationList(ctx, launchConfigResp)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	filteredList := common.FilterModels(ctx, data.Filters, launchConfigList)

	if err := verify.ValidateOneResult(len(filteredList)); err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"Get LaunchConfigurationList result validation",
			err.Error(),
		)
		resp.Diagnostics.Append(diags...)
		return
	}

	state := filteredList[0]
	state.Filters = data.Filters

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenLaunchConfigurationList(ctx context.Context, launchConfigList []*LaunchConfiguration) ([]*launchConfigurationDataSourceModel, diag.Diagnostics) {
	var outputs []*launchConfigurationDataSourceModel

	for _, v := range launchConfigList {
		var output launchConfigurationDataSourceModel
		output.refreshFromOutput(ctx, v)
		outputs = append(outputs, &output)
	}
	return outputs, nil
}

func getLaunchConfigurationList(ctx context.Context, config *conn.ProviderConfig, id string) ([]*LaunchConfiguration, error) {
	if config.SupportVPC {
		return getVpcLaunchConfigurationList(ctx, config, id)
	} else {
		return getClassicLaunchConfigurationList(ctx, config, id)
	}
}
func getClassicLaunchConfigurationList(ctx context.Context, config *conn.ProviderConfig, id string) ([]*LaunchConfiguration, error) {
	no := ncloud.String(id)
	reqParams := &autoscaling.GetLaunchConfigurationListRequest{
		RegionNo: &config.RegionNo,
	}

	tflog.Info(ctx, "getClassicLaunchConfigurationList reqParams", map[string]any{
		"getClassicLaunchConfigurationList reqParams": common.MarshalUncheckedString(reqParams),
	})
	resp, err := config.Client.Autoscaling.V2Api.GetLaunchConfigurationList(reqParams)

	if err != nil {
		return nil, err
	}
	tflog.Info(ctx, "getClassicLaunchConfiguratioList response", map[string]any{
		"getClassicLaunchConfigurationList response": common.MarshalUncheckedString(resp),
	})

	list := make([]*LaunchConfiguration, 0)
	for _, l := range resp.LaunchConfigurationList {
		launchConfiguration := &LaunchConfiguration{
			LaunchConfigurationNo:       l.LaunchConfigurationNo,
			LaunchConfigurationName:     l.LaunchConfigurationName,
			ServerImageProductCode:      l.ServerImageProductCode,
			MemberServerImageInstanceNo: l.MemberServerImageNo,
			ServerProductCode:           l.ServerProductCode,
			LoginKeyName:                l.LoginKeyName,
			UserData:                    l.UserData,
			AccessControlGroupNoList:    flattenAccessControlGroupList(l.AccessControlGroupList),
		}

		if *l.LaunchConfigurationNo == *no {
			return []*LaunchConfiguration{launchConfiguration}, nil
		}

		list = append(list, launchConfiguration)
	}

	if *no != "" {
		return nil, nil
	}
	return list, nil
}

func getVpcLaunchConfigurationList(ctx context.Context, config *conn.ProviderConfig, id string) ([]*LaunchConfiguration, error) {
	reqParams := &vautoscaling.GetLaunchConfigurationListRequest{
		RegionCode: &config.RegionCode,
	}

	if id != "" {
		reqParams.LaunchConfigurationNoList = []*string{ncloud.String(id)}
	}

	tflog.Info(ctx, "getVpcLaunchConfiguratioList reqParams", map[string]any{
		"getVpcLaunchConfigurationList": common.MarshalUncheckedString(reqParams),
	})

	resp, err := config.Client.Vautoscaling.V2Api.GetLaunchConfigurationList(reqParams)

	if err != nil {
		return nil, err
	}

	tflog.Info(ctx, "getVpcLaunchConfigurationList response", map[string]any{
		"getVpcLaunchConfigurationList": common.MarshalUncheckedString(resp),
	})

	if len(resp.LaunchConfigurationList) < 1 {
		return nil, nil
	}

	list := make([]*LaunchConfiguration, 0)
	for _, l := range resp.LaunchConfigurationList {
		list = append(list, &LaunchConfiguration{
			LaunchConfigurationName:     l.LaunchConfigurationName,
			ServerImageProductCode:      l.ServerImageProductCode,
			MemberServerImageInstanceNo: l.MemberServerImageInstanceNo,
			ServerProductCode:           l.ServerProductCode,
			LoginKeyName:                l.LoginKeyName,
			InitScriptNo:                l.InitScriptNo,
			IsEncryptedVolume:           l.IsEncryptedVolume,
			LaunchConfigurationNo:       l.LaunchConfigurationNo,
		})
	}
	return list, nil

}

func (l *launchConfigurationDataSourceModel) refreshFromOutput(ctx context.Context, output *LaunchConfiguration) {
	l.ID = types.StringPointerValue(output.LaunchConfigurationNo)
	l.LaunchConfigurationNo = types.StringPointerValue(output.LaunchConfigurationNo)
	l.LaunchConfigurationName = types.StringPointerValue(output.LaunchConfigurationName)
	l.ServerImageProductCode = types.StringPointerValue(output.ServerImageProductCode)
	l.ServerProductCode = types.StringPointerValue(output.ServerProductCode)
	l.MemberServerImageInstanceNo = types.StringPointerValue(output.MemberServerImageInstanceNo)
	l.LoginKeyName = types.StringPointerValue(output.LoginKeyName)
	l.InitScriptNo = types.StringPointerValue(output.InitScriptNo)
	l.UserData = types.StringPointerValue(output.UserData)
	accessControlGroupNoList, _ := types.ListValueFrom(ctx, types.StringType, output.AccessControlGroupNoList)
	l.AccessControlGroupNoList = accessControlGroupNoList
	l.IsEncryptedVolume = types.BoolPointerValue(output.IsEncryptedVolume)
}

type launchConfigurationDataSourceModel struct {
	ID                          types.String `tfsdk:"id"`
	LaunchConfigurationNo       types.String `tfsdk:"launch_configuration_no"`
	LaunchConfigurationName     types.String `tfsdk:"name"`
	ServerImageProductCode      types.String `tfsdk:"server_image_product_code"`
	ServerProductCode           types.String `tfsdk:"server_product_code"`
	MemberServerImageInstanceNo types.String `tfsdk:"member_server_image_no"`
	LoginKeyName                types.String `tfsdk:"login_key_name"`
	InitScriptNo                types.String `tfsdk:"init_script_no"`
	UserData                    types.String `tfsdk:"user_data"`
	AccessControlGroupNoList    types.List   `tfsdk:"access_control_group_no_list"`
	IsEncryptedVolume           types.Bool   `tfsdk:"is_encrypted_volume"`
	Filters                     types.Set    `tfsdk:"filter"`
}
