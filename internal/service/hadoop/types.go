package hadoop

import (
	"context"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vhadoop"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type hadoopServer struct {
	HadoopServerName  types.String `tfsdk:"hadoop_server_name"`
	HadoopServerRole  types.String `tfsdk:"hadoop_server_role"`
	HadoopProductCode types.String `tfsdk:"hadoop_product_code"`
	RegionCode        types.String `tfsdk:"region_code"`
	ZoneCode          types.String `tfsdk:"zone_code"`
	VpcNo             types.String `tfsdk:"vpc_no"`
	SubnetNo          types.String `tfsdk:"subnet_no"`
	IsPublicSubnet    types.Bool   `tfsdk:"is_public_subnet"`
	DataStorageSize   types.Int64  `tfsdk:"data_storage_size"`
	CpuCount          types.Int64  `tfsdk:"cpu_count"`
	MemorySize        types.Int64  `tfsdk:"memory_size"`
}

func (h hadoopServer) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"hadoop_server_name":  types.StringType,
		"hadoop_server_role":  types.StringType,
		"hadoop_product_code": types.StringType,
		"region_code":         types.StringType,
		"zone_code":           types.StringType,
		"vpc_no":              types.StringType,
		"subnet_no":           types.StringType,
		"is_public_subnet":    types.BoolType,
		"data_storage_size":   types.Int64Type,
		"cpu_count":           types.Int64Type,
		"memory_size":         types.Int64Type,
	}
}

func listValueFromHadoopServerInatanceList(ctx context.Context, serverInatances []*vhadoop.CloudHadoopServerInstance) (basetypes.ListValue, diag.Diagnostics) {
	var hadoopServerList []hadoopServer
	for _, serverInstance := range serverInatances {
		hadoopServerList = append(hadoopServerList, hadoopServer{
			HadoopServerName:  types.StringPointerValue(serverInstance.CloudHadoopServerName),
			HadoopServerRole:  types.StringPointerValue(serverInstance.CloudHadoopServerRole.CodeName),
			HadoopProductCode: types.StringPointerValue(serverInstance.CloudHadoopProductCode),
			RegionCode:        types.StringPointerValue(serverInstance.RegionCode),
			ZoneCode:          types.StringPointerValue(serverInstance.ZoneCode),
			VpcNo:             types.StringPointerValue(serverInstance.VpcNo),
			SubnetNo:          types.StringPointerValue(serverInstance.SubnetNo),
			IsPublicSubnet:    types.BoolPointerValue(serverInstance.IsPublicSubnet),
			DataStorageSize:   types.Int64Value(*serverInstance.DataStorageSize),
			CpuCount:          types.Int64Value(int64(*serverInstance.CpuCount)),
			MemorySize:        types.Int64Value(*serverInstance.MemorySize),
		})
	}

	return types.ListValueFrom(ctx, types.ObjectType{AttrTypes: hadoopServer{}.attrTypes()}, hadoopServerList)
}
