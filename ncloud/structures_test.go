package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnks"
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

	if result[0] != "25363" {
		t.Fatalf("expected configuration_no to be 25363, but was %s", result[0])
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

	if result[0] != "4964" {
		t.Fatalf("expected configuration_no to be 4964, but was %s", result[0])
	}

	if result[1] != "30067" {
		t.Fatalf("expected configuration_no to be 30067, but was %s", result[0])
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

	if result[0] != "4653" {
		t.Fatalf("expected result no to be 4653, but was %s", result[0])
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

	result := flattenServerProducts(expanded)

	if result == nil {
		t.Fatal("result was nil")
	}

	if len(result) != 1 {
		t.Fatalf("expected result had %d elements, but got %d", 1, len(result))
	}

	if result[0] != "SPSVRSTAND000056" {
		t.Fatalf("expected result product_code to be SPSVRSTAND000056, but was %s", result[0])
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
			SnapshotVolumeConfigurationRatio: ncloud.Float32(0.0),
			SnapshotVolumeSize:               ncloud.Int64(0),
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

	if result[0] != "856180" {
		t.Fatalf("expected result instance_no to be 856180, but was %s", result[0])
	}
}

func TestExpandLoadBalancerRuleParams(t *testing.T) {
	lbrulelist := []interface{}{
		map[string]interface{}{
			"protocol_type":        "HTTP",
			"load_balancer_port":   80,
			"server_port":          80,
			"l7_health_check_path": "/monitor/l7check",
		},
		map[string]interface{}{
			"protocol_type":        "HTTPS",
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

func TestFlattenMapByKey(t *testing.T) {
	expanded := &server.CommonCode{
		Code: ncloud.String("test"),
	}

	result := flattenMapByKey(expanded, "code")

	if result == nil {
		t.Fatal("result was nil")
	}

	if *result != "test" {
		t.Fatalf("result expected 'test' but was %s", *result)
	}
}

func TestFlattenArrayStructByKey(t *testing.T) {
	list := []*server.InstanceTag{
		{
			InstanceNo: ncloud.String("foo"),
		},
		{
			InstanceNo: ncloud.String("bar"),
		},
	}

	result := flattenArrayStructByKey(list, "instanceNo")

	if result == nil {
		t.Fatal("result was nil")
	}

	if *result[0] != "foo" {
		t.Fatalf("result expected 'foo' but was %s", *result[0])
	}

	if *result[1] != "bar" {
		t.Fatalf("result expected 'bar' but was %s", *result[1])
	}

	list = []*server.InstanceTag{}
	result = flattenArrayStructByKey(list, "instanceNo")

	if result == nil {
		t.Fatal("result was nil")
	}

	if len(result) != 0 {
		t.Fatalf("result len(result) expected 0 but was %d", len(result))
	}

}

func TestGetInt32FromString(t *testing.T) {

	var tests = []struct {
		input    string
		expected int32
	}{
		{"1", 1},
		{"-1", -1},
		{"0", 0},
		{"", 0},
		{"foo", 0},
	}

	for k, v := range tests {
		result := getInt32FromString(v.input, true)
		if ncloud.Int32Value(result) != v.expected {
			t.Fatalf("case %d: expected %d, but got %d", k, v.expected, result)
		}
	}

	result := getInt32FromString("1", false)

	if ncloud.Int32Value(result) != 0 {
		t.Fatalf("result expected '0' but was %d", *result)
	}
}

func TestExpandStringInterfaceListToInt32List(t *testing.T) {
	initialList := []string{"1111", "2222", "3333"}
	l := make([]interface{}, len(initialList))
	for i, v := range initialList {
		l[i] = v
	}
	int32List := expandStringInterfaceListToInt32List(l)
	expected := []*int32{
		ncloud.Int32(int32(1111)),
		ncloud.Int32(int32(2222)),
		ncloud.Int32(int32(3333)),
	}
	if !reflect.DeepEqual(int32List, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			int32List,
			expected)
	}
}

func TestFlattenInt32ListToStringList(t *testing.T) {
	initialList := []*int32{
		ncloud.Int32(int32(1111)),
		ncloud.Int32(int32(2222)),
		ncloud.Int32(int32(3333)),
	}

	stringList := flattenInt32ListToStringList(initialList)
	expected := []*string{
		ncloud.String("1111"),
		ncloud.String("2222"),
		ncloud.String("3333")}
	if !reflect.DeepEqual(stringList, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			stringList,
			expected)
	}
}

func TestFlattenNKSClusterLogInput(t *testing.T) {
	logInput := &vnks.ClusterLogInput{Audit: ncloud.Bool(true)}

	result := flattenNKSClusterLogInput(logInput)

	if result == nil {
		t.Fatal("result was nil")
	}

	r := result[0]
	if r["audit"].(bool) != true {
		t.Fatalf("expected result enabled to be true, but was %v", r["enabled"])
	}
}

func TestExpandNKSClusterLogInput(t *testing.T) {
	log := []interface{}{
		map[string]interface{}{
			"audit": false,
		},
	}

	result := expandNKSClusterLogInput(log)

	if result == nil {
		t.Fatal("result was nil")
	}

	if ncloud.BoolValue(result.Audit) != false {
		t.Fatalf("expected false , but got %v", ncloud.BoolValue(result.Audit))
	}
}

func TestFlattenNKSNodePoolAutoscale(t *testing.T) {
	expanded := &vnks.AutoscaleOption{
		Enabled: ncloud.Bool(true),
		Max:     ncloud.Int32(2),
		Min:     ncloud.Int32(2),
	}

	result := flattenNKSNodePoolAutoScale(expanded)

	if result == nil {
		t.Fatal("result was nil")
	}

	r := result[0]
	if r["enabled"].(bool) != true {
		t.Fatalf("expected result enabled to be true, but was %v", r["enabled"])
	}

	if r["min"].(int32) != 2 {
		t.Fatalf("expected result min to be 2, but was %d", r["min"])
	}

	if r["max"].(int32) != 2 {
		t.Fatalf("expected result max to be 2, but was %d", r["max"])
	}
}

func TestFlattenNKSWorkerNodes(t *testing.T) {
	expanded := []*vnks.WorkerNode{
		{
			Id:            ncloud.Int32(1),
			Name:          ncloud.String("node1"),
			ServerSpec:    ncloud.String("[Standard] vCPU 2EA, Memory 8GB"),
			PrivateIp:     ncloud.String("10.0.1.4"),
			PublicIp:      ncloud.String(""),
			K8sStatus:     ncloud.String("Ready"),
			DockerVersion: ncloud.String("containerd://1.3.7"),
			KernelVersion: ncloud.String("5.4.0-65-generic"),
		},
	}

	result := flattenNKSWorkerNodes(expanded)

	if result == nil {
		t.Fatal("result was nil")
	}

	r := result[0]
	if r["instance_no"].(int32) != 1 {
		t.Fatalf("expected result instance_no to be 1, but was %v", r["instance_no"])
	}

	if r["name"].(string) != "node1" {
		t.Fatalf("expected result name to be node1, but was %s", r["name"])
	}

	if r["spec"].(string) != "[Standard] vCPU 2EA, Memory 8GB" {
		t.Fatalf("expected result spec to be [Standard] vCPU 2EA, Memory 8GB, but was %s", r["spec"])
	}

	if r["private_ip"].(string) != "10.0.1.4" {
		t.Fatalf("expected result private_ip to be 10.0.1.4, but was %s", r["private_ip"])
	}

	if r["public_ip"].(string) != "" {
		t.Fatalf("expected result public_ip to be emtpy, but was %s", r["public_ip"])
	}

	if r["node_status"].(string) != "Ready" {
		t.Fatalf("expected result node_status to be Ready, but was %s", r["node_status"])
	}

	if r["container_version"].(string) != "containerd://1.3.7" {
		t.Fatalf("expected result container_version to be containerd://1.3.7, but was %s", r["container_version"])
	}

	if r["kernel_version"].(string) != "5.4.0-65-generic" {
		t.Fatalf("expected result kernel_version to be 5.4.0-65-generic, but was %s", r["kernel_version"])
	}
}

func TestExpandNKSNodePoolAutoScale(t *testing.T) {
	autoscaleList := []interface{}{
		map[string]interface{}{
			"enabled": true,
			"min":     2,
			"max":     2,
		},
	}

	result := expandNKSNodePoolAutoScale(autoscaleList)

	if result == nil {
		t.Fatal("result was nil")
	}

	if ncloud.BoolValue(result.Enabled) != true {
		t.Fatalf("expected result true, but got %v", ncloud.BoolValue(result.Enabled))
	}

	if ncloud.Int32Value(result.Min) != int32(2) {
		t.Fatalf("expected result 2, but got %d", ncloud.Int32Value(result.Min))
	}

	if ncloud.Int32Value(result.Max) != int32(2) {
		t.Fatalf("expected result 2, but got %d", ncloud.Int32Value(result.Max))
	}
}

func TestExpandSourceBuildEnvVarsParams(t *testing.T) {
	envVars := []interface{}{
		map[string]interface{}{
			"key":   "key1",
			"value": "value1",
		},
		map[string]interface{}{
			"key":   "key2",
			"value": "value2",
		},
	}

	result, _ := expandSourceBuildEnvVarsParams(envVars)

	if result == nil {
		t.Fatal("result was nil")
	}

	if len(result) != 2 {
		t.Fatalf("expected result had %d elements, but got %d", 2, len(result))
	}

	env := result[0]
	if *env.Key != "key1" {
		t.Fatalf("expected result key1, but got %s", *env.Key)
	}

	if *env.Value != "value1" {
		t.Fatalf("expected result value1, but got %s", *env.Value)
	}

	env2 := result[1]
	if *env2.Key != "key2" {
		t.Fatalf("expected result key2, but got %s", *env2.Key)
	}

	if *env2.Value != "value2" {
		t.Fatalf("expected result value2, but got %s", *env2.Value)
	}
}
