package sdk

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"

	common "github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	request "github.com/NaverCloudPlatform/ncloud-sdk-go/request"
)

func processAddNasVolumeAccessControlParams(reqParams *RequestNasVolumeAccessControl) (map[string]string, error) {
	params := make(map[string]string)

	if reqParams == nil {
		return nil, errors.New("NasVolumeInstanceNo is required")
	}

	if err := validateRequiredField("NasVolumeInstanceNo", reqParams.NasVolumeInstanceNo); err != nil {
		return nil, err
	}
	params["nasVolumeInstanceNo"] = reqParams.NasVolumeInstanceNo

	if len(reqParams.ServerInstanceNoList) > 0 {
		for k, v := range reqParams.ServerInstanceNoList {
			params[fmt.Sprintf("serverInstanceNoList.%d", k+1)] = v
		}
	}

	if len(reqParams.CustomIPList) > 0 {
		for k, v := range reqParams.CustomIPList {
			params[fmt.Sprintf("customIpList.%d", k+1)] = v
		}
	}

	return params, nil
}

// AddNasVolumeAccessControl add Nas Volume Access Control
func (s *Conn) AddNasVolumeAccessControl(reqParams *RequestNasVolumeAccessControl) (*NasVolumeInstanceList, error) {
	params, err := processAddNasVolumeAccessControlParams(reqParams)

	if err != nil {
		return nil, err
	}

	params["action"] = "addNasVolumeAccessControl"

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
