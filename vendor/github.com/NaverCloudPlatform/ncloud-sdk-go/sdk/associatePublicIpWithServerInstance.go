package sdk

import (
	"encoding/xml"
	"fmt"
	"net/http"

	common "github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	request "github.com/NaverCloudPlatform/ncloud-sdk-go/request"
)

func processAssociatePublicIPParams(reqParams *RequestAssociatePublicIP) (map[string]string, error) {
	params := make(map[string]string)

	if reqParams.PublicIPInstanceNo != "" {
		params["publicIpInstanceNo"] = reqParams.PublicIPInstanceNo
	}

	if reqParams.ServerInstanceNo != "" {
		params["serverInstanceNo"] = reqParams.ServerInstanceNo
	}

	return params, nil
}

// AssociatePublicIP associate ip instance to server instance
func (s *Conn) AssociatePublicIP(reqParams *RequestAssociatePublicIP) (*PublicIPInstanceList, error) {
	params, err := processAssociatePublicIPParams(reqParams)
	if err != nil {
		return nil, err
	}

	params["action"] = "associatePublicIpWithServerInstance"

	bytes, resp, err := request.NewRequest(s.accessKey, s.secretKey, "POST", s.apiURL+"server/", params)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		responseError, err := common.ParseErrorResponse(bytes)
		if err != nil {
			return nil, err
		}

		respError := PublicIPInstanceList{}
		respError.ReturnCode = responseError.ReturnCode
		respError.ReturnMessage = responseError.ReturnMessage

		return &respError, fmt.Errorf("%s %s - error code: %d , error message: %s", resp.Status, string(bytes), responseError.ReturnCode, responseError.ReturnMessage)
	}

	var responseCreatePublicIPInstances = PublicIPInstanceList{}
	if err := xml.Unmarshal([]byte(bytes), &responseCreatePublicIPInstances); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &responseCreatePublicIPInstances, nil
}
