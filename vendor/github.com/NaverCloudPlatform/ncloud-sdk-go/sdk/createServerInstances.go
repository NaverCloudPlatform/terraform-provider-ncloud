package sdk

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	common "github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	request "github.com/NaverCloudPlatform/ncloud-sdk-go/request"
)

func processCreateServerInstancesParams(reqParams *RequestCreateServerInstance) (map[string]string, error) {
	params := make(map[string]string)

	if reqParams == nil {
		return params, nil
	}

	if reqParams.ServerImageProductCode != "" {
		if err := validateStringMaxLen("ServerImageProductCode", reqParams.ServerImageProductCode, 20); err != nil {
			return nil, err
		}
		params["serverImageProductCode"] = reqParams.ServerImageProductCode
	}

	if reqParams.ServerProductCode != "" {
		if err := validateStringMaxLen("ServerProductCode", reqParams.ServerProductCode, 20); err != nil {
			return nil, err
		}
		params["serverProductCode"] = reqParams.ServerProductCode
	}

	if reqParams.MemberServerImageNo != "" {
		params["memberServerImageNo"] = reqParams.MemberServerImageNo
	}

	if reqParams.ServerName != "" {
		if err := validateStringLenBetween("ServerName", reqParams.ServerName, 3, 30); err != nil {
			return nil, err
		}
		params["serverName"] = reqParams.ServerName
	}

	if reqParams.ServerDescription != "" {
		if err := validateStringMaxLen("ServerDescription", reqParams.ServerDescription, 1000); err != nil {
			return nil, err
		}
		params["serverDescription"] = reqParams.ServerDescription
	}

	if reqParams.LoginKeyName != "" {
		if err := validateStringLenBetween("LoginKeyName", reqParams.LoginKeyName, 3, 30); err != nil {
			return nil, err
		}
		params["loginKeyName"] = reqParams.LoginKeyName
	}

	if reqParams.IsProtectServerTermination != "" {
		if err := validateBoolValue("IsProtectServerTermination", reqParams.IsProtectServerTermination); err != nil {
			return nil, err
		}
		params["isProtectServerTermination"] = reqParams.IsProtectServerTermination
	}

	if reqParams.ServerCreateCount > 0 {
		if err := validateIntegerInRange("ServerCreateCount", reqParams.ServerCreateCount, 1, 20); err != nil {
			return nil, err
		}
		params["serverCreateCount"] = strconv.Itoa(reqParams.ServerCreateCount)
	}

	if reqParams.ServerCreateStartNo > 0 {
		if err := validateIntegerInRange("Sum of ServerCreateCount and ServerCreateStartNo", reqParams.ServerCreateCount+reqParams.ServerCreateStartNo, 0, 1000); err != nil {
			return nil, err
		}

		params["serverCreateStartNo"] = strconv.Itoa(reqParams.ServerCreateStartNo)
	}

	if reqParams.InternetLineTypeCode != "" {
		if err := validateIncludeValues("InternetLineTypeCode", reqParams.InternetLineTypeCode, []string{"PUBLC", "GLBL"}); err != nil {
			return nil, err
		}
		params["internetLineTypeCode"] = reqParams.InternetLineTypeCode
	}

	if reqParams.FeeSystemTypeCode != "" {
		if err := validateIncludeValues("FeeSystemTypeCode", reqParams.FeeSystemTypeCode, []string{"FXSUM", "MTRAT"}); err != nil {
			return nil, err
		}
		params["feeSystemTypeCode"] = reqParams.FeeSystemTypeCode
	}

	if reqParams.UserData != "" {
		if err := validateStringMaxLen("UserData", reqParams.UserData, 21847); err != nil {
			return nil, err
		}
		params["userData"] = base64.StdEncoding.EncodeToString([]byte(reqParams.UserData))
	}

	if reqParams.ZoneNo != "" {
		params["zoneNo"] = reqParams.ZoneNo
	}

	if len(reqParams.AccessControlGroupConfigurationNoList) > 0 {
		for k, v := range reqParams.AccessControlGroupConfigurationNoList {
			params[fmt.Sprintf("accessControlGroupConfigurationNoList.%d", k+1)] = v
		}
	}

	if reqParams.RaidTypeName != "" {
		params["raidTypeName"] = reqParams.RaidTypeName
	}

	return params, nil
}

// CreateServerInstances create server instances
func (s *Conn) CreateServerInstances(reqParams *RequestCreateServerInstance) (*ServerInstanceList, error) {
	params, err := processCreateServerInstancesParams(reqParams)
	if err != nil {
		return nil, err
	}

	params["action"] = "createServerInstances"

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

	var responseCreateServerInstances = ServerInstanceList{}
	if err := xml.Unmarshal([]byte(bytes), &responseCreateServerInstances); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &responseCreateServerInstances, nil
}
