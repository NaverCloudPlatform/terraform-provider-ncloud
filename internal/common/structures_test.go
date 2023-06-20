package common

import (
	"reflect"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
)

func TestExpandStringInterfaceList(t *testing.T) {
	initialList := []string{"1111", "2222", "3333"}
	l := make([]interface{}, len(initialList))
	for i, v := range initialList {
		l[i] = v
	}
	stringList := ExpandStringInterfaceList(l)
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

	result := FlattenCommonCode(expanded)

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

func TestFlattenArrayStructByKey(t *testing.T) {
	list := []*server.InstanceTag{
		{
			InstanceNo: ncloud.String("foo"),
		},
		{
			InstanceNo: ncloud.String("bar"),
		},
	}

	result := FlattenArrayStructByKey(list, "instanceNo")

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
	result = FlattenArrayStructByKey(list, "instanceNo")

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
		result := GetInt32FromString(v.input, true)
		if ncloud.Int32Value(result) != v.expected {
			t.Fatalf("case %d: expected %d, but got %d", k, v.expected, result)
		}
	}

	result := GetInt32FromString("1", false)

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
	int32List := ExpandStringInterfaceListToInt32List(l)
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
