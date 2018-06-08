package sdk

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/request"
)

func processGetPortForwardingRuleList(reqParams *RequestPortForwardingRuleList) (map[string]string, error) {
	params := make(map[string]string)

	if reqParams == nil {
		return params, nil
	}

	if reqParams.InternetLineTypeCode != "" {
		if reqParams.InternetLineTypeCode != "PUBLC" && reqParams.InternetLineTypeCode != "GLBL" {
			return nil, errors.New("InternetLineTypeCode should be PUBLC or GLBL")
		}
		params["internetLineTypeCode"] = reqParams.InternetLineTypeCode
	}

	if reqParams.RegionNo != "" {
		params["regionNo"] = reqParams.RegionNo
	}

	if reqParams.ZoneNo != "" {
		params["zoneNo"] = reqParams.ZoneNo
	}

	return params, nil
}

func (s *Conn) GetPortForwardingRuleList(reqParams *RequestPortForwardingRuleList) (*PortForwardingRuleList, error) {
	params, err := processGetPortForwardingRuleList(reqParams)
	if err != nil {
		return nil, err
	}

	params["action"] = "getPortForwardingRuleList"

	bytes, resp, err := request.NewRequest(s.accessKey, s.secretKey, "GET", s.apiURL+"server/", params)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		responseError, err := common.ParseErrorResponse(bytes)
		if err != nil {
			return nil, err
		}

		respError := PortForwardingRuleList{}
		respError.ReturnCode = responseError.ReturnCode
		respError.ReturnMessage = responseError.ReturnMessage

		return &respError, fmt.Errorf("%s %s - error code: %d , error message: %s", resp.Status, string(bytes), responseError.ReturnCode, responseError.ReturnMessage)
	}

	var portForwardingRuleList = PortForwardingRuleList{}
	if err := xml.Unmarshal([]byte(bytes), &portForwardingRuleList); err != nil {
		return nil, err
	}

	return &portForwardingRuleList, nil
}
