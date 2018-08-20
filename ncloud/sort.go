package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"sort"
	"time"
)

var defaultDateFormat = "2006-01-02T15:04:00+0900"

type memberServerImageSort []*server.MemberServerImage

func (a memberServerImageSort) Len() int {
	return len(a)
}
func (a memberServerImageSort) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a memberServerImageSort) Less(i, j int) bool {
	iTime, _ := time.Parse(defaultDateFormat, ncloud.StringValue(a[i].CreateDate))
	jTime, _ := time.Parse(defaultDateFormat, ncloud.StringValue(a[j].CreateDate))

	return iTime.Unix() < jTime.Unix()
}

func mostRecentMemberServerImage(images []*server.MemberServerImage) *server.MemberServerImage {
	sortedImages := images
	sort.Sort(memberServerImageSort(sortedImages))
	return sortedImages[len(sortedImages)-1]
}

type acgSort []*server.AccessControlGroup

func (a acgSort) Len() int {
	return len(a)
}
func (a acgSort) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a acgSort) Less(i, j int) bool {
	iTime, _ := time.Parse(defaultDateFormat, ncloud.StringValue(a[i].CreateDate))
	jTime, _ := time.Parse(defaultDateFormat, ncloud.StringValue(a[j].CreateDate))
	return iTime.Unix() < jTime.Unix()
}

func mostRecentAccessControlGroup(acgs []*server.AccessControlGroup) *server.AccessControlGroup {
	sortedAcgs := acgs
	sort.Sort(acgSort(sortedAcgs))
	return sortedAcgs[len(sortedAcgs)-1]
}

type publicIPSort []*server.PublicIpInstance

func (a publicIPSort) Len() int {
	return len(a)
}
func (a publicIPSort) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a publicIPSort) Less(i, j int) bool {
	iTime, _ := time.Parse(defaultDateFormat, ncloud.StringValue(a[i].CreateDate))
	jTime, _ := time.Parse(defaultDateFormat, ncloud.StringValue(a[j].CreateDate))
	return iTime.Unix() < jTime.Unix()
}

func mostRecentPublicIp(publicIPs []*server.PublicIpInstance) *server.PublicIpInstance {
	sortedPublicIps := publicIPs
	sort.Sort(publicIPSort(sortedPublicIps))
	return sortedPublicIps[len(sortedPublicIps)-1]
}
