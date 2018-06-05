package common

type Zone struct {
	ZoneNo          string `xml:"zoneNo"`
	ZoneName        string `xml:"zoneName"`
	ZoneCode        string `xml:"zoneCode"`
	ZoneDescription string `xml:"zoneDescription"`
	RegionNo        string `xml:"regionNo"`
}
