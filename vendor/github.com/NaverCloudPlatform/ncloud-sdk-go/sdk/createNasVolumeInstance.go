package sdk

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"strconv"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/request"
)

func processCreateNasVolumeInstance(reqParams *RequestCreateNasVolumeInstance) (map[string]string, error) {
	params := make(map[string]string)

	if err := validateRequiredField("VolumeName", reqParams.VolumeName); err != nil {
		return nil, err
	}
	params["volumeName"] = reqParams.VolumeName

	if err := validateIntegerInRange("VolumeSize", reqParams.VolumeSize, 500, 10000); err != nil {
		return nil, err
	}

	if err := validateMultipleValue("VolumeSize", reqParams.VolumeSize, 100); err != nil {
		return nil, err
	}
	params["volumeSize"] = strconv.Itoa(reqParams.VolumeSize)

	if err := validateIncludeValues("VolumeAllotmentProtocolTypeCode", reqParams.VolumeAllotmentProtocolTypeCode, []string{"NFS", "CIFS"}); err != nil {
		return nil, err
	}
	params["volumeAllotmentProtocolTypeCode"] = reqParams.VolumeAllotmentProtocolTypeCode

	if len(reqParams.ServerInstanceNoList) > 0 {
		for k, v := range reqParams.ServerInstanceNoList {
			params[fmt.Sprintf("serverInstanceNoList.%d", k+1)] = v
		}
	}

	if len(reqParams.CustomIpList) > 0 {
		for k, v := range reqParams.CustomIpList {
			params[fmt.Sprintf("customIpList.%d", k+1)] = v
		}
	}

	if reqParams.CifsUserName != "" {
		params["cifsUserName"] = reqParams.CifsUserName
	}

	if reqParams.CifsUserPassword != "" {
		params["cifsUserPassword"] = reqParams.CifsUserPassword
	}

	if reqParams.NasVolumeDescription != "" {
		if err := validateStringMaxLen("NasVolumeDescription", reqParams.NasVolumeDescription, 1000); err != nil {
			return nil, err
		}
		params["nasVolumeDescription"] = reqParams.NasVolumeDescription
	}

	if reqParams.RegionNo != "" {
		params["regionNo"] = reqParams.RegionNo
	}

	if reqParams.ZoneNo != "" {
		params["zoneNo"] = reqParams.ZoneNo
	}

	return params, nil
}

func (s *Conn) CreateNasVolumeInstance(reqParams *RequestCreateNasVolumeInstance) (*NasVolumeInstanceList, error) {

	params, err := processCreateNasVolumeInstance(reqParams)
	if err != nil {
		return nil, err
	}

	params["action"] = "createNasVolumeInstance"

	bytes, resp, err := request.NewRequest(s.accessKey, s.secretKey, "POST", s.apiURL+"server/", params)
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

	var nasVolumeInstanceList = NasVolumeInstanceList{}
	if err := xml.Unmarshal([]byte(bytes), &nasVolumeInstanceList); err != nil {
		return nil, err
	}

	return &nasVolumeInstanceList, nil
}
