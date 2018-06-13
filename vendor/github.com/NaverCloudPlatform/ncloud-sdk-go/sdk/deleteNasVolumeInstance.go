package sdk

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/request"
)

func processDeleteNasVolumeInstanceParams(nasVolumeInstanceNo string) (map[string]string, error) {
	params := make(map[string]string)

	if err := validateRequiredField("nasVolumeInstanceNo", nasVolumeInstanceNo); err != nil {
		return params, err
	}
	params["nasVolumeInstanceNo"] = nasVolumeInstanceNo

	return params, nil
}

func (s *Conn) DeleteNasVolumeInstance(nasVolumeInstanceNo string) (*NasVolumeInstanceList, error) {
	params, err := processDeleteNasVolumeInstanceParams(nasVolumeInstanceNo)
	if err != nil {
		return nil, err
	}

	params["action"] = "deleteNasVolumeInstance"

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
		fmt.Println(err)
		return nil, err
	}

	return &nasVolumeInstanceList, nil
}
