package sdk

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"

	common "github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	request "github.com/NaverCloudPlatform/ncloud-sdk-go/request"
)

func processChangeServerInstanceSpecParams(reqParams *RequestChangeServerInstanceSpec) (map[string]string, error) {
	params := make(map[string]string)

	if reqParams.ServerInstanceNo == "" {
		return nil, errors.New("ServerInstanceNo field is required")
	}
	if reqParams.ServerProductCode == "" {
		return nil, errors.New("ServerProductCode field is required")
	}

	params["serverInstanceNo"] = reqParams.ServerInstanceNo
	params["serverProductCode"] = reqParams.ServerProductCode

	return params, nil
}

func (s *Conn) ChangeServerInstanceSpec(reqParams *RequestChangeServerInstanceSpec) (*ServerInstanceList, error) {
	params, err := processChangeServerInstanceSpecParams(reqParams)
	if err != nil {
		return nil, err
	}

	params["action"] = "changeServerInstanceSpec"

	bytes, resp, err := request.NewRequest(s.accessKey, s.secretKey, "POST", s.apiURL+"server/", params)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		responseError, err := common.ParseErrorResponse(bytes)
		if err != nil {
			return nil, err
		}

		respError := ServerInstanceList{}
		respError.ReturnCode = responseError.ReturnCode
		respError.ReturnMessage = responseError.ReturnMessage

		return &respError, fmt.Errorf("%s %s - error code: %d , error message: %s", resp.Status, string(bytes), responseError.ReturnCode, responseError.ReturnMessage)
	}

	var serverInstanceList = ServerInstanceList{}
	if err := xml.Unmarshal([]byte(bytes), &serverInstanceList); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &serverInstanceList, nil
}
