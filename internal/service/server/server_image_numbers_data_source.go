package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

var (
	_ datasource.DataSource              = &serverImageNumbersDataSource{}
	_ datasource.DataSourceWithConfigure = &serverImageNumbersDataSource{}
)

func NewServerImageNumbersDataSource() datasource.DataSource {
	return &serverImageNumbersDataSource{}
}

type serverImageNumbersDataSource struct {
	config *conn.ProviderConfig
}

func (d *serverImageNumbersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_image_numbers"
}

func (d *serverImageNumbersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"server_image_name": schema.StringAttribute{
				Optional: true,
			},
			"hypervisor_type": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"XEN", "KVM"}...),
				},
			},
			"output_file": schema.StringAttribute{
				Optional: true,
			},
			"image_number_list": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"server_image_number": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"hypervisor_type": schema.StringAttribute{
							Computed: true,
						},
						"cpu_architecture_type": schema.StringAttribute{
							Computed: true,
						},
						"os_category_type": schema.StringAttribute{
							Computed: true,
						},
						"os_type": schema.StringAttribute{
							Computed: true,
						},
						"product_code": schema.StringAttribute{
							Computed: true,
						},
						"block_storage_mapping_list": schema.ListNestedAttribute{
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"order": schema.Int32Attribute{
										Computed: true,
									},
									"block_storage_snapshot_instance_no": schema.Int32Attribute{
										Computed: true,
									},
									"block_storage_snapshot_name": schema.StringAttribute{
										Computed: true,
									},
									"block_storage_size": schema.Int64Attribute{
										Computed: true,
									},
									"block_storage_name": schema.StringAttribute{
										Computed: true,
									},
									"block_storage_volume_type": schema.StringAttribute{
										Computed: true,
									},
									"iops": schema.Int32Attribute{
										Computed: true,
									},
									"throughput": schema.Int64Attribute{
										Computed: true,
									},
									"is_encrypted_volume": schema.BoolAttribute{
										Computed: true,
									},
								},
							},
							Computed: true,
						},
					},
				},
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"filter": common.DataSourceFiltersBlock(),
		},
	}
}

func (d *serverImageNumbersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *serverImageNumbersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data serverImageNumbersDataSourceModel

	if !d.config.SupportVPC {
		resp.Diagnostics.AddError(
			"NOT SUPPORT CLASSIC",
			"does not support CLASSIC. only VPC.",
		)
		return
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqParams := &vserver.GetServerImageListRequest{
		RegionCode:         &d.config.RegionCode,
		ServerImageName:    data.ServerImageName.ValueStringPointer(),
		HypervisorCodeList: []*string{data.HypervisorType.ValueStringPointer()},
	}
	tflog.Info(ctx, "GetServerImageListRequest reqParams="+common.MarshalUncheckedString(reqParams))

	imageNoResp, err := d.config.Client.Vserver.V2Api.GetServerImageList(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("READING ERROR", err.Error())
		return
	}
	tflog.Info(ctx, "GetServerImageListRequest response="+common.MarshalUncheckedString(imageNoResp))

	if imageNoResp == nil || len(imageNoResp.ServerImageList) < 1 {
		resp.Diagnostics.AddError("READING ERROR", "no result.")
		return
	}

	imagesNoList, diags := flattenServerImageList(ctx, imageNoResp.ServerImageList)
	if diags.HasError() {
		resp.Diagnostics.AddError("READING ERROR", "flattenServerImageList error")
		return
	}
	fillteredList := common.FilterModels(ctx, data.Filters, imagesNoList)
	diags = data.refreshFromOutput(ctx, fillteredList)
	if diags.HasError() {
		resp.Diagnostics.AddError("READING ERROR", "refreshFromOutput error")
		return
	}

	if !data.OutputFile.IsNull() && data.OutputFile.String() != "" {
		outputPath := data.OutputFile.ValueString()

		if convertedList, err := convertImagesToJsonStruct(data.ImageNumberList.Elements()); err != nil {
			resp.Diagnostics.AddError("OUTPUT FILE ERROR", err.Error())
			return
		} else if err := common.WriteToFile(outputPath, convertedList); err != nil {
			resp.Diagnostics.AddError("OUTPUT FILE ERROR", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func convertImagesToJsonStruct(images []attr.Value) ([]serverImageNoToJsonConvert, error) {
	var serverImagesToConvert = []serverImageNoToJsonConvert{}

	for _, image := range images {
		imageJasn := serverImageNoToJsonConvert{}
		if err := json.Unmarshal([]byte(common.ReplaceNull(image.String())), &imageJasn); err != nil {
			return nil, err
		}
		serverImagesToConvert = append(serverImagesToConvert, imageJasn)
	}

	return serverImagesToConvert, nil
}

func flattenServerImageList(ctx context.Context, list []*vserver.ServerImage) ([]*serverImageNo, diag.Diagnostics) {
	var outputs []*serverImageNo
	var diags diag.Diagnostics

	for _, v := range list {
		var output serverImageNo
		diags = output.refreshFromOutput(ctx, v)
		if diags.HasError() {
			return nil, diags
		}

		outputs = append(outputs, &output)
	}
	return outputs, diags
}

type serverImageNumbersDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	ServerImageName types.String `tfsdk:"server_image_name"`
	HypervisorType  types.String `tfsdk:"hypervisor_type"`
	ImageNumberList types.List   `tfsdk:"image_number_list"`
	OutputFile      types.String `tfsdk:"output_file"`
	Filters         types.Set    `tfsdk:"filter"`
}

type serverImageNo struct {
	Number              types.String `tfsdk:"server_image_number"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	Type                types.String `tfsdk:"type"`
	HypervisorType      types.String `tfsdk:"hypervisor_type"`
	CpuArchitectureType types.String `tfsdk:"cpu_architecture_type"`
	OsCategoryType      types.String `tfsdk:"os_category_type"`
	OsType              types.String `tfsdk:"os_type"`
	ProductCode         types.String `tfsdk:"product_code"`
	BlockStorageMapList types.List   `tfsdk:"block_storage_mapping_list"`
}

type blockStorageMap struct {
	Order                          types.Int32  `tfsdk:"order"`
	BlockStorageSnapshotInstanceNo types.Int32  `tfsdk:"block_storage_snapshot_instance_no"`
	BlockStorageSnapshotName       types.String `tfsdk:"block_storage_snapshot_name"`
	BlockStorageSize               types.Int64  `tfsdk:"block_storage_size"`
	BlockStorageName               types.String `tfsdk:"block_storage_name"`
	BlockStorageVolumeType         types.String `tfsdk:"block_storage_volume_type"`
	Iops                           types.Int32  `tfsdk:"iops"`
	Throughput                     types.Int64  `tfsdk:"throughput"`
	IsEncryptedVolume              types.Bool   `tfsdk:"is_encrypted_volume"`
}

type serverImageNoToJsonConvert struct {
	Number              string                         `json:"server_image_number"`
	Name                string                         `json:"name"`
	Description         string                         `json:"description"`
	Type                string                         `json:"type"`
	HypervisorType      string                         `json:"hypervisor_type"`
	CpuArchitectureType string                         `json:"cpu_architecture_type"`
	OsCategoryType      string                         `json:"os_category_type"`
	OsType              string                         `json:"os_type"`
	ProductCode         string                         `json:"product_code"`
	BlockStorageMapList []blockStorageMapToJsonConvert `json:"block_storage_mapping_list"`
}

type blockStorageMapToJsonConvert struct {
	Order                          int32  `json:"order"`
	BlockStorageSnapshotInstanceNo int32  `json:"block_storage_snapshot_instance_no,omitempty"`
	BlockStorageSnapshotName       string `json:"block_storage_snapshot_name,omitempty"`
	BlockStorageSize               int64  `json:"block_storage_size"`
	BlockStorageName               string `json:"block_storage_name,omitempty"`
	BlockStorageVolumeType         string `json:"block_storage_volume_type"`
	Iops                           int32  `json:"iops,omitempty"`
	Throughput                     int64  `json:"throughput,omitempty"`
	IsEncryptedVolume              bool   `json:"is_encrypted_volume"`
}

func (d serverImageNo) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"server_image_number":        types.StringType,
		"name":                       types.StringType,
		"description":                types.StringType,
		"type":                       types.StringType,
		"hypervisor_type":            types.StringType,
		"cpu_architecture_type":      types.StringType,
		"os_category_type":           types.StringType,
		"os_type":                    types.StringType,
		"product_code":               types.StringType,
		"block_storage_mapping_list": types.ListType{ElemType: types.ObjectType{AttrTypes: blockStorageMap{}.attrTypes()}},
	}
}

func (d blockStorageMap) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"order":                              types.Int32Type,
		"block_storage_snapshot_instance_no": types.Int32Type,
		"block_storage_snapshot_name":        types.StringType,
		"block_storage_size":                 types.Int64Type,
		"block_storage_name":                 types.StringType,
		"block_storage_volume_type":          types.StringType,
		"iops":                               types.Int32Type,
		"throughput":                         types.Int64Type,
		"is_encrypted_volume":                types.BoolType,
	}
}

func (d *serverImageNumbersDataSourceModel) refreshFromOutput(ctx context.Context, list []*serverImageNo) diag.Diagnostics {
	imageNoListValue, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: serverImageNo{}.attrTypes()}, list)
	if diags.HasError() {
		return diags
	}

	d.ImageNumberList = imageNoListValue
	d.ID = types.StringValue("")

	return diags
}

func (d *serverImageNo) refreshFromOutput(ctx context.Context, output *vserver.ServerImage) diag.Diagnostics {
	d.Number = types.StringPointerValue(output.ServerImageNo)
	d.Name = types.StringPointerValue(output.ServerImageName)
	d.Description = types.StringPointerValue(output.ServerImageDescription)
	d.Type = types.StringPointerValue(output.ServerImageType.Code)
	d.HypervisorType = types.StringPointerValue(output.HypervisorType.Code)
	d.CpuArchitectureType = types.StringPointerValue(output.CpuArchitectureType.Code)
	d.OsCategoryType = types.StringPointerValue(output.OsCategoryType.Code)
	d.OsType = types.StringPointerValue(output.OsType.Code)
	d.ProductCode = types.StringPointerValue(output.ServerImageProductCode)

	var blockStorageList []blockStorageMap
	for _, block := range output.BlockStorageMappingList {
		blockStorage := blockStorageMap{
			Order:                          types.Int32PointerValue(block.Order),
			BlockStorageSnapshotInstanceNo: types.Int32PointerValue(block.BlockStorageSnapshotInstanceNo),
			BlockStorageSnapshotName:       types.StringPointerValue(block.BlockStorageSnapshotName),
			BlockStorageSize:               types.Int64PointerValue(block.BlockStorageSize),
			BlockStorageName:               types.StringPointerValue(block.BlockStorageName),
			BlockStorageVolumeType:         types.StringPointerValue(block.BlockStorageVolumeType.Code),
			Iops:                           types.Int32PointerValue(block.Iops),
			Throughput:                     types.Int64PointerValue(block.Throughput),
			IsEncryptedVolume:              types.BoolPointerValue(block.IsEncryptedVolume),
		}
		blockStorageList = append(blockStorageList, blockStorage)
	}
	blockStorageMaps, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: blockStorageMap{}.attrTypes()}, blockStorageList)
	if diags.HasError() {
		return diags
	}
	d.BlockStorageMapList = blockStorageMaps

	return diags
}
