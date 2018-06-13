package sdk

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"

	"strconv"

	common "github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	request "github.com/NaverCloudPlatform/ncloud-sdk-go/request"
)

func processChangeNasVolumeSizeParams(reqParams *RequestChangeNasVolumeSize) (map[string]string, error) {
	params := make(map[string]string)

	if reqParams == nil || reqParams.NasVolumeInstanceNo == "" {
		return nil, errors.New("NasVolumeInstanceNo field is required")
	}
	params["nasVolumeInstanceNo"] = reqParams.NasVolumeInstanceNo

	if err := validateIntegerInRange("VolumeSize", reqParams.VolumeSize, 500, 10000); err != nil {
		return nil, err
	}
	if err := validateMultipleValue("VolumeSize", reqParams.VolumeSize, 100); err != nil {
		return nil, err
	}
	params["volumeSize"] = strconv.Itoa(reqParams.VolumeSize)

	return params, nil
}

// ChangeNasVolumeSize changes nas volume size
func (s *Conn) ChangeNasVolumeSize(reqParams *RequestChangeNasVolumeSize) (*NasVolumeInstanceList, error) {
	params, err := processChangeNasVolumeSizeParams(reqParams)
	if err != nil {
		return nil, err
	}

	params["action"] = "changeNasVolumeSize"

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
