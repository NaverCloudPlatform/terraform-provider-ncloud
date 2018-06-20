package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"sort"
	"time"
)

type serverImageSort []sdk.ServerImage

func (a serverImageSort) Len() int {
	return len(a)
}
func (a serverImageSort) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a serverImageSort) Less(i, j int) bool {
	iTime, _ := time.Parse(time.RFC3339, a[i].CreateDate)
	jTime, _ := time.Parse(time.RFC3339, a[j].CreateDate)
	return iTime.Unix() < jTime.Unix()
}

type acgSort []sdk.AccessControlGroup

func (a acgSort) Len() int {
	return len(a)
}
func (a acgSort) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a acgSort) Less(i, j int) bool {
	iTime, _ := time.Parse(time.RFC3339, a[i].CreateDate)
	jTime, _ := time.Parse(time.RFC3339, a[j].CreateDate)
	return iTime.Unix() < jTime.Unix()
}

func mostRecentAccessControlGroup(acgs []sdk.AccessControlGroup) sdk.AccessControlGroup {
	sortedAcgs := acgs
	sort.Sort(acgSort(sortedAcgs))
	return sortedAcgs[len(sortedAcgs)-1]
}
