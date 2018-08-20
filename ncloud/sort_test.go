package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"testing"
)

func TestMostRecentMemberServerImage(t *testing.T) {
	recentDate := "2018-06-22T15:21:00+0900"
	images := []*server.MemberServerImage{
		{MemberServerImageNo: ncloud.String("1755"), CreateDate: ncloud.String("2014-02-06T15:21:41+0900")},
		{MemberServerImageNo: ncloud.String("1756"), CreateDate: ncloud.String(recentDate)},
		{MemberServerImageNo: ncloud.String("1753"), CreateDate: ncloud.String("2012-06-22T15:21:00+0900")},
	}

	if mostRecent := mostRecentMemberServerImage(images); recentDate != *mostRecent.CreateDate {
		t.Fatalf("Expected: %s, Actual: %s", recentDate, *mostRecent.CreateDate)
	}
}

func TestMostRecentAccessControlGroup(t *testing.T) {
	recentDate := "2018-06-22T15:21:00+0900"
	images := []*server.AccessControlGroup{
		{AccessControlGroupConfigurationNo: ncloud.String("1"), CreateDate: ncloud.String("2014-02-06T15:21:41+0900")},
		{AccessControlGroupConfigurationNo: ncloud.String("2"), CreateDate: ncloud.String(recentDate)},
		{AccessControlGroupConfigurationNo: ncloud.String("3"), CreateDate: ncloud.String("2012-06-22T15:21:00+0900")},
	}

	if mostRecent := mostRecentAccessControlGroup(images); recentDate != *mostRecent.CreateDate {
		t.Fatalf("Expected: %s, Actual: %s", recentDate, *mostRecent.CreateDate)
	}
}

func TestMostRecentPublicIp(t *testing.T) {
	recentDate := "2018-06-22T15:21:00+0900"
	images := []*server.PublicIpInstance{
		{PublicIpInstanceNo: ncloud.String("1"), CreateDate: ncloud.String("2014-02-06T15:21:41+0900")},
		{PublicIpInstanceNo: ncloud.String("2"), CreateDate: ncloud.String(recentDate)},
		{PublicIpInstanceNo: ncloud.String("3"), CreateDate: ncloud.String("2012-06-22T15:21:00+0900")},
	}

	if mostRecent := mostRecentPublicIp(images); recentDate != *mostRecent.CreateDate {
		t.Fatalf("Expected: %s, Actual: %s", recentDate, *mostRecent.CreateDate)
	}
}
