package vpc

import (
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
)

func TestVpcCommonStateRefreshFunc(t *testing.T) {
	run := "RUN"
	instance := &vpc.RouteTable{
		RouteTableStatus: &vpc.CommonCode{
			Code:     &run,
			CodeName: &run,
		},
	}

	var empty *vpc.RouteTable

	testVpcCommonStateRefreshFunc(instance, "RUN", t)
	testVpcCommonStateRefreshFunc(empty, "TERMINATED", t)
}

func testVpcCommonStateRefreshFunc(instance interface{}, expected string, t *testing.T) {
	_, status, err := VpcCommonStateRefreshFunc(instance, nil, "routeTableStatus")

	if err != nil {
		t.Fatal("Got Error")
	}

	if status != expected {
		t.Fatalf("Expected: %s, Actual: %s", expected, status)
	}
}
