package sdk

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/request"
)

func processGetNasVolumeInstanceListParams(reqParams *RequestGetNasVolumeInstanceList) (map[string]string, error) {
	params := make(map[string]string)

	if reqParams == nil {
		return params, nil
	}

	if reqParams.VolumeAllotmentProtocolTypeCode != "" {
		if err := validateIncludeValues("VolumeAllotmentProtocolTypeCode", reqParams.VolumeAllotmentProtocolTypeCode, []string{"NFS", "CIFS"}); err != nil {
			return nil, err
		}
		params["volumeAllotmentProtocolTypeCode"] = reqParams.VolumeAllotmentProtocolTypeCode
	}

	if reqParams.IsEventConfiguration != "" {
		if err := validateBoolValue("IsEventConfiguration", reqParams.IsEventConfiguration); err != nil {
			return nil, err
		}
		params["isEventConfiguration"] = reqParams.IsEventConfiguration
	}

	if reqParams.IsSnapshotConfiguration != "" {
		if err := validateBoolValue("IsSnapshotConfiguration", reqParams.IsSnapshotConfiguration); err != nil {
			return nil, err
		}
		params["isSnapshotConfiguration"] = reqParams.IsSnapshotConfiguration
	}

	if len(reqParams.NasVolumeInstanceNoList) > 0 {
		for k, v := range reqParams.NasVolumeInstanceNoList {
			params[fmt.Sprintf("nasVolumeInstanceNoList.%d", k+1)] = v
		}
	}

	if reqParams.RegionNo != "" {
		params["regionNo"] = reqParams.RegionNo
	}

	if reqParams.ZoneNo != "" {
		params["zoneNo"] = reqParams.ZoneNo
	}

	return params, nil
}

func (s *Conn) GetNasVolumeInstanceList(reqParams *RequestGetNasVolumeInstanceList) (*NasVolumeInstanceList, error) {
	params, err := processGetNasVolumeInstanceListParams(reqParams)
	if err != nil {
		return nil, err
	}

	params["action"] = "getNasVolumeInstanceList"

	bytes, resp, err := request.NewRequest(s.accessKey, s.secretKey, "GET", s.apiURL+"server/", params)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		responseError, err := common.ParseErrorResponse(bytes)
		if err != nil {
			return nil, err
		}

		respError := NasVolumeInstanceList{}
		respError.ReturnCode = responseError.ReturnCode
		respError.ReturnMessage = responseError.ReturnMessage

		return &respError, fmt.Errorf("%s %s - error code: %d , error message: %s", resp.Status, string(bytes), responseError.ReturnCode, responseError.ReturnMessage)
	}

	NasVolumeInstanceList := NasVolumeInstanceList{}
	if err := xml.Unmarshal([]byte(bytes), &NasVolumeInstanceList); err != nil {
		return nil, err
	}

	return &NasVolumeInstanceList, nil
}
