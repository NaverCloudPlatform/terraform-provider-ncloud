package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"testing"
)

func TestMostRecentServerImage(t *testing.T) {
	recentDate := "2018-06-22T15:21:00+0900"
	images := []sdk.ServerImage{
		{MemberServerImageNo: "1755", CreateDate: "2014-02-06T15:21:41+0900"},
		{MemberServerImageNo: "1756", CreateDate: recentDate},
		{MemberServerImageNo: "1753", CreateDate: "2012-06-22T15:21:00+0900"},
	}

	if mostRecent := mostRecentServerImage(images); recentDate != mostRecent.CreateDate {
		t.Fatalf("Expected: %s, Actual: %s", recentDate, mostRecent.CreateDate)
	}
}

func TestMostRecentAccessControlGroup(t *testing.T) {
	recentDate := "2018-06-22T15:21:00+0900"
	images := []sdk.AccessControlGroup{
		{AccessControlGroupConfigurationNo: "1", CreateDate: "2014-02-06T15:21:41+0900"},
		{AccessControlGroupConfigurationNo: "2", CreateDate: recentDate},
		{AccessControlGroupConfigurationNo: "3", CreateDate: "2012-06-22T15:21:00+0900"},
	}

	if mostRecent := mostRecentAccessControlGroup(images); recentDate != mostRecent.CreateDate {
		t.Fatalf("Expected: %s, Actual: %s", recentDate, mostRecent.CreateDate)
	}
}

func TestMostRecentPublicIP(t *testing.T) {
	recentDate := "2018-06-22T15:21:00+0900"
	images := []sdk.PublicIPInstance{
		{PublicIPInstanceNo: "1", CreateDate: "2014-02-06T15:21:41+0900"},
		{PublicIPInstanceNo: "2", CreateDate: recentDate},
		{PublicIPInstanceNo: "3", CreateDate: "2012-06-22T15:21:00+0900"},
	}

	if mostRecent := mostRecentPublicIP(images); recentDate != mostRecent.CreateDate {
		t.Fatalf("Expected: %s, Actual: %s", recentDate, mostRecent.CreateDate)
	}
}
