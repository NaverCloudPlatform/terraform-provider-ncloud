package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"sort"
	"time"
)

var defaultDateFormat = "2006-01-02T15:04:00+0900"

type serverImageSort []sdk.ServerImage

func (a serverImageSort) Len() int {
	return len(a)
}
func (a serverImageSort) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a serverImageSort) Less(i, j int) bool {
	iTime, _ := time.Parse(defaultDateFormat, a[i].CreateDate)
	jTime, _ := time.Parse(defaultDateFormat, a[j].CreateDate)

	return iTime.Unix() < jTime.Unix()
}

func mostRecentServerImage(images []sdk.ServerImage) sdk.ServerImage {
	sortedImages := images
	sort.Sort(serverImageSort(sortedImages))
	return sortedImages[len(sortedImages)-1]
}

type acgSort []sdk.AccessControlGroup

func (a acgSort) Len() int {
	return len(a)
}
func (a acgSort) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a acgSort) Less(i, j int) bool {
	iTime, _ := time.Parse(defaultDateFormat, a[i].CreateDate)
	jTime, _ := time.Parse(defaultDateFormat, a[j].CreateDate)
	return iTime.Unix() < jTime.Unix()
}

func mostRecentAccessControlGroup(acgs []sdk.AccessControlGroup) sdk.AccessControlGroup {
	sortedAcgs := acgs
	sort.Sort(acgSort(sortedAcgs))
	return sortedAcgs[len(sortedAcgs)-1]
}

type publicIPSort []sdk.PublicIPInstance

func (a publicIPSort) Len() int {
	return len(a)
}
func (a publicIPSort) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a publicIPSort) Less(i, j int) bool {
	iTime, _ := time.Parse(defaultDateFormat, a[i].CreateDate)
	jTime, _ := time.Parse(defaultDateFormat, a[j].CreateDate)
	return iTime.Unix() < jTime.Unix()
}

func mostRecentPublicIP(publicIPs []sdk.PublicIPInstance) sdk.PublicIPInstance {
	sortedPublicIPs := publicIPs
	sort.Sort(publicIPSort(sortedPublicIPs))
	return sortedPublicIPs[len(sortedPublicIPs)-1]
}
