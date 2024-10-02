package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

var (
	_ datasource.DataSource              = &loginKeyDataSource{}
	_ datasource.DataSourceWithConfigure = &loginKeyDataSource{}
)

func NewLoginKeyDataSource() datasource.DataSource {
	return &loginKeyDataSource{}
}

type loginKeyDataSource struct {
	config *conn.ProviderConfig
}

func (d *loginKeyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_login_key"
}

func (d *loginKeyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.config = config
}

func (d *loginKeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"output_file": schema.StringAttribute{
				Optional: true,
			},
			"login_key_list": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key_name": schema.StringAttribute{
							Computed: true,
						},
						"fingerprint": schema.StringAttribute{
							Computed: true,
						},
						"create_date": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": common.DataSourceFiltersBlock(),
		},
	}
}

func (d *loginKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data loginKeyDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := GetLoginKeyList(d.config)
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}

	if output == nil {
		resp.Diagnostics.AddError("READING ERROR", "no result. please change search criteria and try again.")
		return
	}

	loginKeyList := flattenLoginKey(output)
	fillteredList := common.FilterModels(ctx, data.Filters, loginKeyList)
	data.refreshFromOutput(ctx, fillteredList)

	if !data.OutputFile.IsNull() && data.OutputFile.String() != "" {
		outputPath := data.OutputFile.ValueString()

		if convertedList, err := convertToJsonStruct(data.KeyList.Elements()); err != nil {
			resp.Diagnostics.AddError("OUTPUT FILE ERROR", err.Error())
			return
		} else if err := common.WriteToFile(outputPath, convertedList); err != nil {
			resp.Diagnostics.AddError("OUTPUT FILE ERROR", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

type loginKeyDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	KeyList    types.List   `tfsdk:"login_key_list"`
	OutputFile types.String `tfsdk:"output_file"`
	Filters    types.Set    `tfsdk:"filter"`
}

type loginKeyModel struct {
	KeyName     types.String `tfsdk:"key_name"`
	Fingerprint types.String `tfsdk:"fingerprint"`
	CreateDate  types.String `tfsdk:"create_date"`
}

type loginKeyToJsonConvert struct {
	KeyName     string `json:"key_name"`
	Fingerprint string `json:"fingerprint"`
	CreateDate  string `json:"create_date"`
}

func (d loginKeyModel) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"key_name":    types.StringType,
		"fingerprint": types.StringType,
		"create_date": types.StringType,
	}
}

type loginKeyStruct struct {
	KeyName     *string `json:"key_name,omitempty"`
	Fingerprint *string `json:"fingerprint,omitempty"`
	CreateDate  *string `json:"create_date,omitempty"`
}

func convertToJsonStruct(keys []attr.Value) ([]loginKeyToJsonConvert, error) {
	var loginKeyToConvert = []loginKeyToJsonConvert{}

	for _, key := range keys {
		keyJson := loginKeyToJsonConvert{}
		if err := json.Unmarshal([]byte(key.String()), &keyJson); err != nil {
			return nil, err
		}
		loginKeyToConvert = append(loginKeyToConvert, keyJson)
	}

	return loginKeyToConvert, nil
}

func flattenLoginKey(list []*loginKeyStruct) []*loginKeyModel {
	var outputs []*loginKeyModel

	for _, v := range list {
		var output loginKeyModel
		output.refreshFromOutput(v)

		outputs = append(outputs, &output)
	}
	return outputs
}

func (d *loginKeyDataSourceModel) refreshFromOutput(ctx context.Context, output []*loginKeyModel) {
	keyListValue, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: loginKeyModel{}.attrTypes()}, output)
	d.KeyList = keyListValue
	d.ID = types.StringValue("")
}

func (d *loginKeyModel) refreshFromOutput(output *loginKeyStruct) {
	d.KeyName = types.StringPointerValue(output.KeyName)
	d.Fingerprint = types.StringPointerValue(output.Fingerprint)
	d.CreateDate = types.StringPointerValue(output.CreateDate)
}

func GetLoginKeyList(config *conn.ProviderConfig) ([]*loginKeyStruct, error) {
	if config.SupportVPC {
		return getVpcLoginKeyList(config)
	} else {
		return getClassicLoginKeyList(config)
	}
}

func getVpcLoginKeyList(config *conn.ProviderConfig) ([]*loginKeyStruct, error) {
	resp, err := config.Client.Vserver.V2Api.GetLoginKeyList(&vserver.GetLoginKeyListRequest{})

	if err != nil {
		return nil, err
	}

	if len(resp.LoginKeyList) < 1 {
		return nil, nil
	}

	var loginKeys []*loginKeyStruct
	for _, l := range resp.LoginKeyList {
		loginKeys = append(loginKeys, &loginKeyStruct{
			KeyName:     l.KeyName,
			Fingerprint: l.Fingerprint,
			CreateDate:  l.CreateDate,
		})
	}

	return loginKeys, nil
}

func getClassicLoginKeyList(config *conn.ProviderConfig) ([]*loginKeyStruct, error) {
	resp, err := config.Client.Server.V2Api.GetLoginKeyList(&server.GetLoginKeyListRequest{})

	if err != nil {
		return nil, err
	}

	if len(resp.LoginKeyList) < 1 {
		return nil, nil
	}

	var loginKeys []*loginKeyStruct
	for _, l := range resp.LoginKeyList {
		loginKeys = append(loginKeys, &loginKeyStruct{
			KeyName:     l.KeyName,
			Fingerprint: l.Fingerprint,
			CreateDate:  l.CreateDate,
		})
	}

	return loginKeys, nil
}
