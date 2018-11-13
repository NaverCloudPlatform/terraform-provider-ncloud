package ncloud

import (
	"reflect"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/loadbalancer"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
)

func TestExpandStringInterfaceList(t *testing.T) {
	initialList := []string{"1111", "2222", "3333"}
	l := make([]interface{}, len(initialList))
	for i, v := range initialList {
		l[i] = v
	}
	stringList := expandStringInterfaceList(l)
	expected := []*string{
		ncloud.String("1111"),
		ncloud.String("2222"),
		ncloud.String("3333"),
	}
	if !reflect.DeepEqual(stringList, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			stringList,
			expected)
	}
}

func TestFlattenCommonCode(t *testing.T) {
	expanded := &server.CommonCode{
		Code:     ncloud.String("code"),
		CodeName: ncloud.String("codename"),
	}

	result := flattenCommonCode(expanded)

	if result == nil {
		t.Fatal("result was nil")
	}

	if result["code"] != "code" {
		t.Fatalf("expected result code to be code, but was %s", result["code"])
	}

	if result["code_name"] != "codename" {
		t.Fatalf("expected result code_name to be codename, but was %s", result["code_name"])
	}
}

func TestFlattenAccessControlRules(t *testing.T) {
	expected := []*server.AccessControlRule{
		{
			AccessControlRuleConfigurationNo: ncloud.String("25363"),
			ProtocolType: &server.CommonCode{
				Code:     ncloud.String("TCP"),
				CodeName: ncloud.String("tcp"),
			},
			SourceIp:                               ncloud.String("0.0.0.0/0"),
			SourceAccessControlRuleConfigurationNo: ncloud.String("4964"),
			SourceAccessControlRuleName:            ncloud.String("ncloud-default-acg"),
			DestinationPort:                        ncloud.String("1-65535"),
			AccessControlRuleDescription:           ncloud.String("for test"),
		},
	}

	result := flattenAccessControlRules(expected)

	if result == nil {
		t.Fatal("result was nil")
	}

	if len(result) != 1 {
		t.Fatalf("expected result had %d elements, but got %d", 1, len(result))
	}

	if result[0]["access_control_rule_configuration_no"] != "25363" {
		t.Fatalf("expected access_control_rule_configuration_no to be 25363, but was %s", result[0]["access_control_rule_configuration_no"])
	}

	if result[0]["source_ip"] != "0.0.0.0/0" {
		t.Fatalf("expected source_ip to be 0.0.0.0/0, but was %s", result[0]["source_ip"])
	}

	if result[0]["destination_port"] != "1-65535" {
		t.Fatalf("expected destination_port to be 1-65535, but was %s", result[0]["destination_port"])
	}

	if result[0]["source_access_control_rule_configuration_no"] != "4964" {
		t.Fatalf("expected source_access_control_rule_configuration_no to be 4964, but was %s", result[0]["source_access_control_rule_configuration_no"])
	}

	if result[0]["source_access_control_rule_name"] != "ncloud-default-acg" {
		t.Fatalf("expected source_access_control_rule_name to be ncloud-default-acg, but was %s", result[0]["source_access_control_rule_name"])
	}

	if result[0]["access_control_rule_description"] != "for test" {
		t.Fatalf("expected access_control_rule_description to be 'for test', but was %s", result[0]["access_control_rule_description"])
	}
}

func TestFlattenAccessControlGroups(t *testing.T) {
	expected := []*server.AccessControlGroup{
		{
			AccessControlGroupConfigurationNo: ncloud.String("4964"),
			AccessControlGroupName:            ncloud.String("ncloud-default-acg"),
			AccessControlGroupDescription:     ncloud.String("for test"),
			IsDefaultGroup:                    ncloud.Bool(true),
			CreateDate:                        ncloud.String("2017-02-23T10:25:39+0900"),
		},
		{
			AccessControlGroupConfigurationNo: ncloud.String("30067"),
			AccessControlGroupName:            ncloud.String("httpdport"),
			AccessControlGroupDescription:     ncloud.String("for httpd test"),
			IsDefaultGroup:                    ncloud.Bool(false),
			CreateDate:                        ncloud.String("2018-01-07T10:17:14+0900"),
		},
	}

	result := flattenAccessControlGroups(expected)

	if result == nil {
		t.Fatal("result was nil")
	}

	if len(result) != 2 {
		t.Fatalf("expected result had %d elements, but got %d", 2, len(result))
	}

	if result[0]["access_control_group_configuration_no"] != "4964" {
		t.Fatalf("expected access_control_group_configuration_no to be 4964, but was %s", result[0]["access_control_group_configuration_no"])
	}

	if result[0]["access_control_group_name"] != "ncloud-default-acg" {
		t.Fatalf("expected access_control_group_name to be ncloud-default-acg, but was %s", result[0]["access_control_group_name"])
	}

	if result[0]["access_control_group_description"] != "for test" {
		t.Fatalf("expected access_control_group_description to be 'for test', but was %s", result[0]["access_control_group_description"])
	}

	if result[0]["is_default_group"] != true {
		t.Fatalf("expected is_default_group to be true, but was %b", result[0]["is_default_group"])
	}

	if result[0]["create_date"] != "2017-02-23T10:25:39+0900" {
		t.Fatalf("expected create_date to be 2017-02-23T10:25:39+0900, but was %s", result[0]["create_date"])
	}
}

func TestFlattenRegion(t *testing.T) {
	expanded := &server.Region{
		RegionNo:   ncloud.String("1"),
		RegionCode: ncloud.String("KR"),
		RegionName: ncloud.String("Korea"),
	}

	result := flattenRegion(expanded)

	if result == nil {
		t.Fatal("result was nil")
	}

	if result["region_no"] != "1" {
		t.Fatalf("expected result region_no to be 1, but was %s", result["region_no"])
	}

	if result["region_code"] != "KR" {
		t.Fatalf("expected result region_code to be KR, but was %s", result["region_code"])
	}

	if result["region_name"] != "Korea" {
		t.Fatalf("expected result region_name to be Korea, but was %s", result["region_name"])
	}
}

func TestFlattenZone(t *testing.T) {
	expanded := &server.Zone{
		ZoneNo:          ncloud.String("3"),
		ZoneName:        ncloud.String("KR-2"),
		ZoneCode:        ncloud.String("KR-2"),
		ZoneDescription: ncloud.String("평촌 zone"),
		RegionNo:        ncloud.String("1"),
	}

	result := flattenZone(expanded)

	if result == nil {
		t.Fatal("result was nil")
	}

	if result["zone_no"] != "3" {
		t.Fatalf("expected result zone_no to be 3, but was %s", result["zone_no"])
	}

	if result["zone_code"] != "KR-2" {
		t.Fatalf("expected result zone_code to be KR-2, but was %s", result["zone_code"])
	}

	if result["zone_name"] != "KR-2" {
		t.Fatalf("expected result zone_name to be KR-2, but was %s", result["zone_name"])
	}

	if result["zone_description"] != "평촌 zone" {
		t.Fatalf("expected result zone_description to be 평촌 zone, but was %s", result["zone_description"])
	}

	if result["region_no"] != "1" {
		t.Fatalf("expected result region_no to be 1, but was %s", result["region_no"])
	}
}

func TestFlattenMemberServerImages(t *testing.T) {
	expanded := []*server.MemberServerImage{
		{
			MemberServerImageNo:          ncloud.String("4653"),
			MemberServerImageName:        ncloud.String("test-1514385790"),
			MemberServerImageDescription: ncloud.String("server description"),
			OriginalServerInstanceNo:     ncloud.String("572053"),
			OriginalServerProductCode:    ncloud.String("SPSVRSTAND000004"),
			OriginalServerName:           ncloud.String("svr-9bbaf27a2902b5c"),
			OriginalBaseBlockStorageDiskType: &server.CommonCode{
				Code:     ncloud.String("NET"),
				CodeName: ncloud.String("Network Storage"),
			},
			OriginalServerImageProductCode: ncloud.String("SPSW0LINUX000043"),
			OriginalOsInformation:          ncloud.String("CentOS 5.11 (64-bit)"),
			OriginalServerImageName:        ncloud.String("centos-5.11-64"),
			MemberServerImageStatusName:    ncloud.String("creating"),
			MemberServerImageStatus: &server.CommonCode{
				Code:     ncloud.String("INIT"),
				CodeName: ncloud.String("NSI INIT state"),
			},
			MemberServerImageOperation: &server.CommonCode{
				Code:     ncloud.String("CREAT"),
				CodeName: ncloud.String("NSI CREAT OP"),
			},
			MemberServerImagePlatformType: &server.CommonCode{
				Code:     ncloud.String("LNX64"),
				CodeName: ncloud.String("Linux 64 Bit"),
			},
			CreateDate: ncloud.String("2018-01-07T10:17:14+0900"),
			Region: &server.Region{
				RegionNo:   ncloud.String("1"),
				RegionCode: ncloud.String("KR"),
				RegionName: ncloud.String("Korea"),
			},
			MemberServerImageBlockStorageTotalRows: ncloud.Int32(2),
			MemberServerImageBlockStorageTotalSize: ncloud.Int64(1127428915200),
		},
	}

	result := flattenMemberServerImages(expanded)

	if result == nil {
		t.Fatal("result was nil")
	}

	if len(result) != 1 {
		t.Fatalf("expected result had %d elements, but got %d", 1, len(result))
	}

	r := result[0]

	if r["member_server_image_no"] != "4653" {
		t.Fatalf("expected result member_server_image_no to be 4653, but was %s", r["member_server_image_no"])
	}

	if r["member_server_image_name"] != "test-1514385790" {
		t.Fatalf("expected result member_server_image_name to be test-1514385790, but was %s", r["member_server_image_name"])
	}

	if r["member_server_image_description"] != "server description" {
		t.Fatalf("expected result member_server_image_description to be 'server description', but was %s", r["member_server_image_description"])
	}

	if r["original_server_instance_no"] != "572053" {
		t.Fatalf("expected result original_server_instance_no to be 572053, but was %s", r["original_server_instance_no"])
	}

	if r["original_server_product_code"] != "SPSVRSTAND000004" {
		t.Fatalf("expected result original_server_product_code to be SPSVRSTAND000004, but was %s", r["original_server_product_code"])
	}

	if r["original_server_name"] != "svr-9bbaf27a2902b5c" {
		t.Fatalf("expected result original_server_name to be svr-9bbaf27a2902b5c, but was %s", r["original_server_name"])
	}

	if r["original_server_image_product_code"] != "SPSW0LINUX000043" {
		t.Fatalf("expected result original_server_image_product_code to be SPSW0LINUX000043, but was %s", r["original_server_image_product_code"])
	}

	if r["original_os_information"] != "CentOS 5.11 (64-bit)" {
		t.Fatalf("expected result original_os_information to be CentOS 5.11 (64-bit), but was %s", r["original_os_information"])
	}

	if r["original_server_image_name"] != "centos-5.11-64" {
		t.Fatalf("expected result original_server_image_name to be centos-5.11-64, but was %s", r["original_server_image_name"])
	}

	if r["member_server_image_status_name"] != "creating" {
		t.Fatalf("expected result member_server_image_status_name to be creating, but was %s", r["member_server_image_status_name"])
	}

	if r["create_date"] != "2018-01-07T10:17:14+0900" {
		t.Fatalf("expected result create_date to be 2018-01-07T10:17:14+0900, but was %s", r["create_date"])
	}

	if r["member_server_image_block_storage_total_rows"] != 2 {
		t.Fatalf("expected result member_server_image_block_storage_total_rows to be 2, but was %s", r["member_server_image_block_storage_total_rows"])
	}

	if r["member_server_image_block_storage_total_size"] != 1127428915200 {
		t.Fatalf("expected result member_server_image_block_storage_total_size to be 1127428915200, but was %s", r["member_server_image_block_storage_total_size"])
	}
}

func TestFlattenCustomIPList(t *testing.T) {
	expanded := []*server.NasVolumeInstanceCustomIp{
		{
			CustomIp: ncloud.String("1.1.1.1"),
		},
		{
			CustomIp: ncloud.String("2.2.2.2"),
		},
	}

	result := flattenCustomIPList(expanded)

	if result == nil {
		t.Fatal("result was nil")
	}

	if len(result) != 2 {
		t.Fatalf("expected result had %d elements, but got %d", 2, len(result))
	}

	if result[0] != "1.1.1.1" {
		t.Fatalf("expected result first element to be 1.1.1.1, but was %s", result[0])
	}

	if result[1] != "2.2.2.2" {
		t.Fatalf("expected result first element to be 2.2.2.2, but was %s", result[1])
	}
}

func TestFlattenServerImages(t *testing.T) {
	expanded := []*server.Product{
		{
			ProductCode: ncloud.String("SPSVRSTAND000056"),
			ProductName: ncloud.String("vCPU 1EA, Memory 1GB, Disk 50GB"),
			ProductType: &server.CommonCode{
				Code:     ncloud.String("MICRO"),
				CodeName: ncloud.String("Micro Server"),
			},
			ProductDescription: ncloud.String("vCPU 1EA, Memory 1GB, Disk 50GB"),
			InfraResourceType: &server.CommonCode{
				Code:     ncloud.String("SVR"),
				CodeName: ncloud.String("Server"),
			},
			CpuCount:             ncloud.Int32(1),
			MemorySize:           ncloud.Int64(1073741824),
			BaseBlockStorageSize: ncloud.Int64(53687091200),
			DiskType: &server.CommonCode{
				Code:     ncloud.String("NET"),
				CodeName: ncloud.String("Network Storage"),
			},
			AddBlockStorageSize: ncloud.Int64(0),
		},
	}

	result := flattenServerImages(expanded)

	if result == nil {
		t.Fatal("result was nil")
	}

	if len(result) != 1 {
		t.Fatalf("expected result had %d elements, but got %d", 1, len(result))
	}

	r := result[0]

	if r["product_code"] != "SPSVRSTAND000056" {
		t.Fatalf("expected result product_code to be SPSVRSTAND000056, but was %s", r["product_code"])
	}

	if r["product_name"] != "vCPU 1EA, Memory 1GB, Disk 50GB" {
		t.Fatalf("expected result product_name to be 'vCPU 1EA, Memory 1GB, Disk 50GB', but was %s", r["product_name"])
	}

	if r["product_description"] != "vCPU 1EA, Memory 1GB, Disk 50GB" {
		t.Fatalf("expected result product_description to be 'vCPU 1EA, Memory 1GB, Disk 50GB', but was %s", r["product_description"])
	}

	if r["cpu_count"] != 1 {
		t.Fatalf("expected result cpu_count to be 1, but was %d", r["cpu_count"])
	}

	if r["memory_size"] != 1073741824 {
		t.Fatalf("expected result memory_size to be 1073741824, but was %d", r["memory_size"])
	}

	if r["base_block_storage_size"] != 53687091200 {
		t.Fatalf("expected result base_block_storage_size to be 53687091200, but was %d", r["base_block_storage_size"])
	}

	if r["add_block_storage_size"] != 0 {
		t.Fatalf("expected result add_block_storage_size to be 0, but was %d", r["add_block_storage_size"])
	}
}

func TestFlattenNasVolumeInstances(t *testing.T) {
	expanded := []*server.NasVolumeInstance{
		{
			NasVolumeInstanceNo: ncloud.String("856180"),
			NasVolumeInstanceStatus: &server.CommonCode{
				Code:     ncloud.String("CREAT"),
				CodeName: ncloud.String("NAS create"),
			},
			NasVolumeInstanceOperation: &server.CommonCode{
				Code:     ncloud.String("NULL"),
				CodeName: ncloud.String("NAS NULL OP"),
			},
			NasVolumeInstanceStatusName: ncloud.String("created"),
			CreateDate:                  ncloud.String("2018-07-12T20:32:45+0900"),
			MountInformation:            ncloud.String("10.250.53.74:/n003666_aaa"),
			VolumeAllotmentProtocolType: &server.CommonCode{
				Code:     ncloud.String("NFS"),
				CodeName: ncloud.String("NFS"),
			},
			VolumeName:                       ncloud.String("n003666_aaa"),
			VolumeTotalSize:                  ncloud.Int64(536870912000),
			VolumeSize:                       ncloud.Int64(536870912000),
			VolumeUseSize:                    ncloud.Int64(1314816),
			VolumeUseRatio:                   ncloud.Float32(0.0),
			SnapshotVolumeConfigurationRatio: ncloud.Float32(0.0),
			SnapshotVolumeSize:               ncloud.Int64(0),
			SnapshotVolumeUseSize:            ncloud.Int64(0),
			SnapshotVolumeUseRatio:           ncloud.Float32(0.0),
			IsSnapshotConfiguration:          ncloud.Bool(false),
			IsEventConfiguration:             ncloud.Bool(false),
			Region: &server.Region{
				RegionNo:   ncloud.String("1"),
				RegionCode: ncloud.String("KR"),
				RegionName: ncloud.String("Korea"),
			},
			Zone: &server.Zone{
				ZoneNo:          ncloud.String("3"),
				ZoneName:        ncloud.String("KR-2"),
				ZoneCode:        ncloud.String("KR-2"),
				ZoneDescription: ncloud.String("평촌 zone"),
				RegionNo:        ncloud.String("1"),
			},
		},
	}

	result := flattenNasVolumeInstances(expanded)

	if result == nil {
		t.Fatal("result was nil")
	}

	if len(result) != 1 {
		t.Fatalf("expected result had %d elements, but got %d", 1, len(result))
	}

	r := result[0]

	if r["nas_volume_instance_no"] != "856180" {
		t.Fatalf("expected result nas_volume_instance_no to be 856180, but was %s", r["nas_volume_instance_no"])
	}

	if r["create_date"] != "2018-07-12T20:32:45+0900" {
		t.Fatalf("expected result create_date to be 2018-07-12T20:32:45+0900, but was %s", r["create_date"])
	}

	if r["volume_name"] != "n003666_aaa" {
		t.Fatalf("expected result volume_name to be n003666_aaa, but was %s", r["volume_name"])
	}

	if r["volume_total_size"] != 536870912000 {
		t.Fatalf("expected result volume_total_size to be 536870912000, but was %d", r["volume_total_size"])
	}

	if r["volume_size"] != 536870912000 {
		t.Fatalf("expected result volume_size to be 536870912000, but was %d", r["volume_size"])
	}

	if r["volume_use_size"] != 1314816 {
		t.Fatalf("expected result volume_use_size to be 1314816, but was %d", r["volume_use_size"])
	}

	if r["volume_use_ratio"] != float32(0.0) {
		t.Fatalf("expected result volume_use_ratio to be 0.0, but was %f", r["volume_use_ratio"])
	}

	if r["snapshot_volume_size"] != int64(0) {
		t.Fatalf("expected result snapshot_volume_size to be 0, but was %d", r["snapshot_volume_size"])
	}

	if r["snapshot_volume_use_size"] != int64(0) {
		t.Fatalf("expected result snapshot_volume_use_size to be 0, but was %d", r["snapshot_volume_use_size"])
	}

	if r["snapshot_volume_use_ratio"] != float32(0.0) {
		t.Fatalf("expected result snapshot_volume_use_ratio to be 0.0 , but was %f", r["snapshot_volume_use_ratio"])
	}

	if r["is_snapshot_configuration"] != false {
		t.Fatalf("expected result is_snapshot_configuration to be false, but was %b", r["is_snapshot_configuration"])
	}

	if r["is_event_configuration"] != false {
		t.Fatalf("expected result is_event_configuration to be false, but was %b", r["is_event_configuration"])
	}
}

func TestExpandLoadBalancerRuleParams(t *testing.T) {
	lbrulelist := []interface{}{
		map[string]interface{}{
			"protocol_type_code":   "HTTP",
			"load_balancer_port":   80,
			"server_port":          80,
			"l7_health_check_path": "/monitor/l7check",
		},
		map[string]interface{}{
			"protocol_type_code":   "HTTPS",
			"load_balancer_port":   443,
			"server_port":          443,
			"l7_health_check_path": "/monitor/l7check",
			"certificate_name":     "aaa",
		},
	}

	result, _ := expandLoadBalancerRuleParams(lbrulelist)

	if result == nil {
		t.Fatal("result was nil")
	}

	if len(result) != 2 {
		t.Fatalf("expected result had %d elements, but got %d", 2, len(result))
	}

	r := result[0]
	if *r.ProtocolTypeCode != "HTTP" {
		t.Fatalf("expected result ProtocolTypeCode to be HTTP, but was %s", *r.ProtocolTypeCode)
	}

	if *r.LoadBalancerPort != 80 {
		t.Fatalf("expected result LoadBalancerPort to be 80, but was %d", *r.LoadBalancerPort)
	}

	if *r.ServerPort != 80 {
		t.Fatalf("expected result ServerPort to be 80, but was %d", *r.ServerPort)
	}

	if *r.L7HealthCheckPath != "/monitor/l7check" {
		t.Fatalf("expected result L7HealthCheckPath to be '/monitor/l7check', but was %s", *r.L7HealthCheckPath)
	}

	if r.CertificateName != nil {
		t.Fatalf("expected result CertificateName to be nil, but was %s", *r.CertificateName)
	}

	r = result[1]
	if *r.ProtocolTypeCode != "HTTPS" {
		t.Fatalf("expected result ProtocolTypeCode to be HTTPS, but was %s", *r.ProtocolTypeCode)
	}

	if *r.LoadBalancerPort != 443 {
		t.Fatalf("expected result LoadBalancerPort to be 443, but was %d", *r.LoadBalancerPort)
	}

	if *r.ServerPort != 443 {
		t.Fatalf("expected result ServerPort to be 443, but was %d", *r.ServerPort)
	}

	if *r.L7HealthCheckPath != "/monitor/l7check" {
		t.Fatalf("expected result L7HealthCheckPath to be '/monitor/l7check', but was %s", *r.L7HealthCheckPath)
	}

	if *r.CertificateName != "aaa" {
		t.Fatalf("expected result CertificateName to be aaa, but was %s", *r.CertificateName)
	}
}

func TestFlattenLoadBalancerRuleList(t *testing.T) {
	expanded := []*loadbalancer.LoadBalancerRule{
		{
			ProtocolType: &loadbalancer.CommonCode{
				Code:     ncloud.String("HTTP"),
				CodeName: ncloud.String("http"),
			},
			LoadBalancerPort:   ncloud.Int32(80),
			ServerPort:         ncloud.Int32(80),
			L7HealthCheckPath:  ncloud.String("/monitor/l7check"),
			CertificateName:    ncloud.String("aaa"),
			ProxyProtocolUseYn: ncloud.String("Y"),
		},
	}

	result := flattenLoadBalancerRuleList(expanded)

	if result == nil {
		t.Fatal("result was nil")
	}

	if len(result) != 1 {
		t.Fatalf("expected result had %d elements, but got %d", 1, len(result))
	}

	r := result[0]

	if r["load_balancer_port"] != int32(80) {
		t.Fatalf("expected result load_balancer_port to be 80, but was %d", r["load_balancer_port"])
	}

	if r["server_port"] != int32(80) {
		t.Fatalf("expected result server_port to be 80, but was %d", r["server_port"])
	}

	if r["l7_health_check_path"] != "/monitor/l7check" {
		t.Fatalf("expected result l7_health_check_path to be /monitor/l7check, but was %s", r["l7_health_check_path"])
	}

	if r["certificate_name"] != "aaa" {
		t.Fatalf("expected result certificate_name to be aaa, but was %s", r["certificate_name"])
	}

	if r["proxy_protocol_use_yn"] != "Y" {
		t.Fatalf("expected result proxy_protocol_use_yn to be Y, but was %s", r["proxy_protocol_use_yn"])
	}

}

func TestFlattenLoadBalancedServerInstanceList(t *testing.T) {
	expanded := []*loadbalancer.LoadBalancedServerInstance{
		{
			ServerInstance: &loadbalancer.ServerInstance{
				ServerInstanceNo: ncloud.String("123456"),
			},
		},
		{
			ServerInstance: &loadbalancer.ServerInstance{
				ServerInstanceNo: ncloud.String("234567"),
			},
		},
	}

	result := flattenLoadBalancedServerInstanceList(expanded)

	if result == nil {
		t.Fatal("result was nil")
	}

	if len(result) != 2 {
		t.Fatalf("expected result had %d elements, but got %d", 2, len(result))
	}

	if result[0] != "123456" {
		t.Fatalf("expected result load_balancer_port to be '123456', but was %s", result[0])
	}

	if result[1] != "234567" {
		t.Fatalf("expected result load_balancer_port to be '234567', but was %s", result[1])
	}
}

func TestExpandTagListParams(t *testing.T) {
	lbrulelist := []interface{}{
		map[string]interface{}{
			"tag_key":   "dev",
			"tag_value": "web",
		},
		map[string]interface{}{
			"tag_key":   "prod",
			"tag_value": "auth",
		},
	}

	result, _ := expandTagListParams(lbrulelist)

	if result == nil {
		t.Fatal("result was nil")
	}

	if len(result) != 2 {
		t.Fatalf("expected result had %d elements, but got %d", 2, len(result))
	}

	tag := result[0]
	if *tag.TagKey != "dev" {
		t.Fatalf("expected result ProtocolTypeCode to be dev, but was %s", *tag.TagKey)
	}

	if *tag.TagValue != "web" {
		t.Fatalf("expected result ProtocolTypeCode to be web, but was %s", *tag.TagValue)
	}

	tag = result[1]
	if *tag.TagKey != "prod" {
		t.Fatalf("expected result ProtocolTypeCode to be prod, but was %s", *tag.TagKey)
	}

	if *tag.TagValue != "auth" {
		t.Fatalf("expected result ProtocolTypeCode to be auth, but was %s", *tag.TagValue)
	}
}

func TestFlattenInstanceTagList(t *testing.T) {
	expanded := []*server.InstanceTag{
		{
			TagKey:   ncloud.String("dev"),
			TagValue: ncloud.String("web"),
		},
		{
			TagKey:   ncloud.String("prod"),
			TagValue: ncloud.String("auth"),
		},
	}

	result := flattenInstanceTagList(expanded)

	if result == nil {
		t.Fatal("result was nil")
	}

	if len(result) != 2 {
		t.Fatalf("expected result had %d elements, but got %d", 2, len(result))
	}

	r := result[0]
	if r["tag_key"] != "dev" {
		t.Fatalf("expected result tag_key to be 'dev', but was %s", r["tag_key"])
	}

	if r["tag_value"] != "web" {
		t.Fatalf("expected result tag_value to be 'web', but was %s", r["tag_value"])
	}

	r = result[1]
	if r["tag_key"] != "prod" {
		t.Fatalf("expected result tag_key to be 'prod', but was %s", r["tag_key"])
	}

	if r["tag_value"] != "auth" {
		t.Fatalf("expected result tag_value to be 'auth', but was %s", r["tag_value"])
	}
}
