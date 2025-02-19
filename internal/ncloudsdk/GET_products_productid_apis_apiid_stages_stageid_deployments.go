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

type PrimitiveGETProductsProductidApisApiidStagesStageidDeploymentsRequest struct {
	Productid string `json:"product-id"`
	Apiid     string `json:"api-id"`
	Stageid   string `json:"stage-id"`
	Limit     int64  `json:"limit"`
	Offset    int64  `json:"offset"`
}

type StringifiedGETProductsProductidApisApiidStagesStageidDeploymentsRequest struct {
	Productid string `json:"product-id"`
	Apiid     string `json:"api-id"`
	Stageid   string `json:"stage-id"`
	Limit     string `json:"limit"`
	Offset    string `json:"offset"`
}

func (n *NClient) GETProductsProductidApisApiidStagesStageidDeployments(ctx context.Context, primitiveReq *PrimitiveGETProductsProductidApisApiidStagesStageidDeploymentsRequest) (map[string]interface{}, error) {
	query := map[string]string{}
	initBody := map[string]string{}

	convertedReq, err := ConvertStructToStringMap(*primitiveReq)
	if err != nil {
		return nil, err
	}

	if r.Limit != "" {
		query["limit"] = r.Limit
	}

	if r.Offset != "" {
		query["offset"] = r.Offset
	}

	rawBody, err := json.Marshal(initBody)
	if err != nil {
		return nil, err
	}

	body := strings.Replace(string(rawBody), `\"`, "", -1)

	url := n.BaseURL + "/" + "products" + "/" + ClearDoubleQuote(r.Productid) + "/" + "apis" + "/" + ClearDoubleQuote(r.Apiid) + "/" + "stages" + "/" + ClearDoubleQuote(r.Stageid) + "/" + "deployments"

	response, err := n.MakeRequestWithContext(ctx, "GET", url, body, query)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, fmt.Errorf("output is nil")
	}

	snake_case_response := convertKeys(response).(map[string]interface{})

	return snake_case_response, nil
}

func (n *NClient) GETProductsProductidApisApiidStagesStageidDeployments_TF(ctx context.Context, r *PrimitiveGETProductsProductidApisApiidStagesStageidDeploymentsRequest) (*GETProductsProductidApisApiidStagesStageidDeploymentsResponse, error) {
	t, err := n.GETProductsProductidApisApiidStagesStageidDeployments(ctx, r)
	if err != nil {
		return nil, err
	}

	res, err := ConvertToFrameworkTypes_GETProductsProductidApisApiidStagesStageidDeployments(context.TODO(), t)
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

type GETProductsProductidApisApiidStagesStageidDeploymentsResponse struct {
	Total               types.Int64 `tfsdk:"total"`
	Defaultdeploymentno types.Int64 `tfsdk:"default_deployment_no"`
	Content             types.List  `tfsdk:"content"`
}

func ConvertToFrameworkTypes_GETProductsProductidApisApiidStagesStageidDeployments(ctx context.Context, data map[string]interface{}) (*GETProductsProductidApisApiidStagesStageidDeploymentsResponse, error) {
	var dto GETProductsProductidApisApiidStagesStageidDeploymentsResponse

	if data["total"] != nil {
		dto.Total = types.Int64Value(data["total"].(int64))
	}

	if data["default_deployment_no"] != nil {
		dto.Defaultdeploymentno = types.Int64Value(data["default_deployment_no"].(int64))
	}

	if data["content"] != nil {
		tempContent := data["content"].([]interface{})
		dto.Content = diagOff(types.ListValueFrom, ctx, types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{

			"stage_id":               types.StringType,
			"stage_deployment_no":    types.Int64Type,
			"document_json":          types.StringType,
			"deployment_description": types.StringType,
			"deployed_time":          types.StringType,
		},
		}}.ElementType(), tempContent)
	}

	return &dto, nil
}

func convertToObject_GETProductsProductidApisApiidStagesStageidDeployments(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

	possibleTypes := map[string]attr.Type{}

	for field, fieldType := range possibleTypes {
		attrTypes[field] = fieldType

		if value, exists := data[field]; exists {

			attrValue, err := convertValueToAttr_GETProductsProductidApisApiidStagesStageidDeployments(value)
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

func convertValueToAttr_GETProductsProductidApisApiidStagesStageidDeployments(value interface{}) (attr.Value, error) {
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
