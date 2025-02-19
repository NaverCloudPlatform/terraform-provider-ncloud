
/* =================================================================================
 * NCLOUD SDK LAYER FOR TERRAFORM CODEGEN - DO NOT EDIT
 * =================================================================================
 * Refresh Template
 * Required data are as follows
 *
 *		MethodName         string
 *		PrimitiveRequest   string
 *		StringifiedRequest string
 *		Query              string
 *		Body               string
 *		Path               string
 *		Method             string
 * ================================================================================= */

package ncloudsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PrimitivePATCHProductsProductidApisApiidStagesStageidDeploymentsDeploymentnoResourcesResourceidMethodsMethodnameRequest struct {
    Productid string `json:"product-id"`
Apiid string `json:"api-id"`
Stageid string `json:"stage-id"`
Deploymentno string `json:"deployment-no"`
Resourceid string `json:"resource-id"`
Methodname string `json:"method-name"`
IsInherit bool `json:"isInherit"`
ThrottleRps int32 `json:"throttleRps"`
CacheTtlSec int32 `json:"cacheTtlSec"`
EndpointDomain string `json:"endpointDomain"`

}

type StringifiedPATCHProductsProductidApisApiidStagesStageidDeploymentsDeploymentnoResourcesResourceidMethodsMethodnameRequest struct {
	Productid string `json:"product-id"`
Apiid string `json:"api-id"`
Stageid string `json:"stage-id"`
Deploymentno string `json:"deployment-no"`
Resourceid string `json:"resource-id"`
Methodname string `json:"method-name"`
IsInherit string `json:"isInherit"`
ThrottleRps string `json:"throttleRps"`
CacheTtlSec string `json:"cacheTtlSec"`
EndpointDomain string `json:"endpointDomain"`

}

func (n *NClient) PATCHProductsProductidApisApiidStagesStageidDeploymentsDeploymentnoResourcesResourceidMethodsMethodname(ctx context.Context, primitiveReq *PrimitivePATCHProductsProductidApisApiidStagesStageidDeploymentsDeploymentnoResourcesResourceidMethodsMethodnameRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

 	

	initBody["isInherit"] = r.IsInherit

			if r.ThrottleRps != "" {
				initBody["throttleRps"] = r.ThrottleRps
			}

			if r.CacheTtlSec != "" {
				initBody["cacheTtlSec"] = r.CacheTtlSec
			}

			if r.EndpointDomain != "" {
				initBody["endpointDomain"] = r.EndpointDomain
			}


	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL +"/"+"products"+"/"+ClearDoubleQuote(r.Productid)+"/"+"apis"+"/"+ClearDoubleQuote(r.Apiid)+"/"+"stages"+"/"+ClearDoubleQuote(r.Stageid)+"/"+"deployments"+"/"+ClearDoubleQuote(r.Deploymentno)+"/"+"resources"+"/"+ClearDoubleQuote(r.Resourceid)+"/"+"methods"+"/"+ClearDoubleQuote(r.Methodname)

	response, err := n.MakeRequestWithContext(ctx, "PATCH", url, body, query)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, fmt.Errorf("output is nil")
	}

	snake_case_response := convertKeys(response).(map[string]interface{})

	return snake_case_response, nil
}

func (n *NClient) PATCHProductsProductidApisApiidStagesStageidDeploymentsDeploymentnoResourcesResourceidMethodsMethodname_TF(ctx context.Context, r *PrimitivePATCHProductsProductidApisApiidStagesStageidDeploymentsDeploymentnoResourcesResourceidMethodsMethodnameRequest) (*PATCHProductsProductidApisApiidStagesStageidDeploymentsDeploymentnoResourcesResourceidMethodsMethodnameResponse, error) {
	t, err := n.PATCHProductsProductidApisApiidStagesStageidDeploymentsDeploymentnoResourcesResourceidMethodsMethodname(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_PATCHProductsProductidApisApiidStagesStageidDeploymentsDeploymentnoResourcesResourceidMethodsMethodname(context.TODO(), t)
	if err != nil {
		return nil, err
	}

	return res, nil
}


/* =================================================================================
 * NCLOUD SDK LAYER FOR TERRAFORM CODEGEN - DO NOT EDIT
 * =================================================================================
 * Refresh Template
 * Required data are as follows
 *
 *		Model             string
 *		MethodName        string
 *		RefreshLogic      string
 *		PossibleTypes     string
 *		ConditionalObjectFieldsWithNull string
 * ================================================================================= */

type PATCHProductsProductidApisApiidStagesStageidDeploymentsDeploymentnoResourcesResourceidMethodsMethodnameResponse struct {
    StageMethod         types.Object `tfsdk:"stage_method"`

}

func ConvertToFrameworkTypes_PATCHProductsProductidApisApiidStagesStageidDeploymentsDeploymentnoResourcesResourceidMethodsMethodname(ctx context.Context, data map[string]interface{}) (*PATCHProductsProductidApisApiidStagesStageidDeploymentsDeploymentnoResourcesResourceidMethodsMethodnameResponse, error) {
	var dto PATCHProductsProductidApisApiidStagesStageidDeploymentsDeploymentnoResourcesResourceidMethodsMethodnameResponse

    
			if data["stage_method"] != nil {
				tempStageMethod := data["stage_method"].(map[string]interface{})

				allFields := []string{
					"throttle_rps",
"stage_deployment_no",
"rest_url",
"resource_path",
"resource_id",
"method_name",
"method_code",
"is_inherit",
"invoke_url",
"endpoint_domain",
"endpoint_config_json",
"endpoint_code",
"endpoint_action_id",
"cache_ttl_sec",

				}

				convertedMap := make(map[string]interface{})
				for _, field := range allFields {
					if val, ok := tempStageMethod[field]; ok {
						convertedMap[field] = val
					}
				}

				convertedTempStageMethod, err := convertToObject_PATCHProductsProductidApisApiidStagesStageidDeploymentsDeploymentnoResourcesResourceidMethodsMethodname(ctx, convertedMap)
				if err != nil {
					return nil, err
				}

				dto.StageMethod = diagOff(types.ObjectValueFrom, ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
					"throttle_rps": types.Int32Type,
"stage_deployment_no": types.Int64Type,
"rest_url": types.StringType,
"resource_path": types.StringType,
"resource_id": types.StringType,
"method_name": types.StringType,
"method_code": types.StringType,
"is_inherit": types.BoolType,
"invoke_url": types.StringType,
"endpoint_domain": types.StringType,
"endpoint_config_json": types.StringType,
"endpoint_code": types.StringType,
"endpoint_action_id": types.StringType,
"cache_ttl_sec": types.Int32Type,

				}}.AttributeTypes(), convertedTempStageMethod)
			}


	return &dto, nil
}

func convertToObject_PATCHProductsProductidApisApiidStagesStageidDeploymentsDeploymentnoResourcesResourceidMethodsMethodname(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

    possibleTypes := map[string]attr.Type{
        "throttle_rps": types.Int32Type,
"stage_deployment_no": types.Int64Type,
"rest_url": types.StringType,
"resource_path": types.StringType,
"resource_id": types.StringType,
"method_name": types.StringType,
"method_code": types.StringType,
"is_inherit": types.BoolType,
"invoke_url": types.StringType,
"endpoint_domain": types.StringType,
"endpoint_config_json": types.StringType,
"endpoint_code": types.StringType,
"endpoint_action_id": types.StringType,
"cache_ttl_sec": types.Int32Type,


	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			

			attrValue, err := convertValueToAttr_PATCHProductsProductidApisApiidStagesStageidDeploymentsDeploymentnoResourcesResourceidMethodsMethodname(value)
			if err != nil {
				return types.Object{}, fmt.Errorf("error converting field %s: %v", field, err)
			}
			attrValues[field] = attrValue
		} else {
            

			switch fieldType {
			case types.StringType:
				attrValues[field] = types.StringNull()
			case types.Int64Type:
				attrValues[field] = types.Int64Null()
			case types.BoolType:
				attrValues[field] = types.BoolNull()
			}
		}
	}

	r, diag := types.ObjectValue(attrTypes, attrValues)
	if diag.HasError() {
		return types.Object{}, fmt.Errorf("error from converting object: %v", diag)
	}

	// OK
	return r, nil
}

func convertValueToAttr_PATCHProductsProductidApisApiidStagesStageidDeploymentsDeploymentnoResourcesResourceidMethodsMethodname(value interface{}) (attr.Value, error) {
     switch v := value.(type) {
     case string:
         return types.StringValue(v), nil
     case int32:
         return types.Int32Value(v), nil
     case int64:
         return types.Int64Value(v), nil
     case float64:
         return types.Float64Value(v), nil
     case bool:
         return types.BoolValue(v), nil
     case nil:
         return types.StringNull(), nil
     default:
         return nil, fmt.Errorf("unsupported type: %T", value)
     }
}

