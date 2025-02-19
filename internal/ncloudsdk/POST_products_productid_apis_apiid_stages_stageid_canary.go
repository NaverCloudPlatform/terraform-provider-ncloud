
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

type PrimitivePOSTProductsProductidApisApiidStagesStageidCanaryRequest struct {
    Productid string `json:"product-id"`
Apiid string `json:"api-id"`
Stageid string `json:"stage-id"`

}

type StringifiedPOSTProductsProductidApisApiidStagesStageidCanaryRequest struct {
	Productid string `json:"product-id"`
Apiid string `json:"api-id"`
Stageid string `json:"stage-id"`

}

func (n *NClient) POSTProductsProductidApisApiidStagesStageidCanary(ctx context.Context, primitiveReq *PrimitivePOSTProductsProductidApisApiidStagesStageidCanaryRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

 	

	

	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL +"/"+"products"+"/"+ClearDoubleQuote(r.Productid)+"/"+"apis"+"/"+ClearDoubleQuote(r.Apiid)+"/"+"stages"+"/"+ClearDoubleQuote(r.Stageid)+"/"+"canary"

	response, err := n.MakeRequestWithContext(ctx, "POST", url, body, query)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, fmt.Errorf("output is nil")
	}

	snake_case_response := convertKeys(response).(map[string]interface{})

	return snake_case_response, nil
}

func (n *NClient) POSTProductsProductidApisApiidStagesStageidCanary_TF(ctx context.Context, r *PrimitivePOSTProductsProductidApisApiidStagesStageidCanaryRequest) (*POSTProductsProductidApisApiidStagesStageidCanaryResponse, error) {
	t, err := n.POSTProductsProductidApisApiidStagesStageidCanary(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_POSTProductsProductidApisApiidStagesStageidCanary(context.TODO(), t)
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

type POSTProductsProductidApisApiidStagesStageidCanaryResponse struct {
    Usedistributionrate         types.Bool `tfsdk:"use_distribution_rate"`
Stageid         types.String `tfsdk:"stage_id"`
Deployedstagedeploymentno         types.Int64`tfsdk:"deployed_stage_deployment_no"`
Canarythrottlerps         types.Int32`tfsdk:"canary_throttle_rps"`
Canaryendpointdomain         types.String `tfsdk:"canary_endpoint_domain"`
Canarydistributionrate         types.Float64 `tfsdk:"canary_distribution_rate"`
Canarydeploymentno         types.Int64`tfsdk:"canary_deployment_no"`
Canarydeploymentdescription         types.String `tfsdk:"canary_deployment_description"`
Canarydeployedtime         types.String `tfsdk:"canary_deployed_time"`
CanaryConditions         types.List `tfsdk:"canary_conditions"`
Canarycachettlsec         types.Int32`tfsdk:"canary_cache_ttl_sec"`

}

func ConvertToFrameworkTypes_POSTProductsProductidApisApiidStagesStageidCanary(ctx context.Context, data map[string]interface{}) (*POSTProductsProductidApisApiidStagesStageidCanaryResponse, error) {
	var dto POSTProductsProductidApisApiidStagesStageidCanaryResponse

    
			if data["use_distribution_rate"] != nil {
				dto.Usedistributionrate = types.BoolValue(data["use_distribution_rate"].(bool))
			}

			if data["stage_id"] != nil {
				dto.Stageid = types.StringValue(data["stage_id"].(string))
			}

				if data["deployed_stage_deployment_no"] != nil {
					dto.Deployedstagedeploymentno = types.Int64Value(data["deployed_stage_deployment_no"].(int64))
				}

				if data["canary_throttle_rps"] != nil {
					dto.Canarythrottlerps = types.Int32Value(data["canary_throttle_rps"].(int32))
				}

			if data["canary_endpoint_domain"] != nil {
				dto.Canaryendpointdomain = types.StringValue(data["canary_endpoint_domain"].(string))
			}

			if data["canary_distribution_rate"] != nil {
				dto.Canarydistributionrate = types.Float64Value(data["canary_distribution_rate"].(float64))
			}

				if data["canary_deployment_no"] != nil {
					dto.Canarydeploymentno = types.Int64Value(data["canary_deployment_no"].(int64))
				}

			if data["canary_deployment_description"] != nil {
				dto.Canarydeploymentdescription = types.StringValue(data["canary_deployment_description"].(string))
			}

			if data["canary_deployed_time"] != nil {
				dto.Canarydeployedtime = types.StringValue(data["canary_deployed_time"].(string))
			}

				if data["canary_conditions"] != nil {
					tempCanaryConditions := data["canary_conditions"].([]interface{})
					dto.CanaryConditions = diagOff(types.ListValueFrom, ctx, types.ListType{ElemType:
						
	types.ObjectType{AttrTypes: map[string]attr.Type{
		
		"parameter_value": types.StringType,
"parameter_name": types.StringType,
"parameter_code": types.StringType,

	},

					}}.ElementType(), tempCanaryConditions)
				}
				if data["canary_cache_ttl_sec"] != nil {
					dto.Canarycachettlsec = types.Int32Value(data["canary_cache_ttl_sec"].(int32))
				}


	return &dto, nil
}

func convertToObject_POSTProductsProductidApisApiidStagesStageidCanary(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

    possibleTypes := map[string]attr.Type{
        
	}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			

			attrValue, err := convertValueToAttr_POSTProductsProductidApisApiidStagesStageidCanary(value)
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

func convertValueToAttr_POSTProductsProductidApisApiidStagesStageidCanary(value interface{}) (attr.Value, error) {
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

