package apigw

import (
	"context"
	"os"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/ncloudsdk"
)

func (plan *PostproductresponseModel) refreshFromOutput_createOp(ctx context.Context, diagnostics *diag.Diagnostics, createRes map[string]interface{}) {

	// Allocate resource id from create response
	id := createRes[""].(string)

	// Indicate where to get resource id from create response
	err := plan.waitResourceCreated(ctx, id)

	if err != nil {
		diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}

	var postPlan PostproductresponseModel

	c := ncloudsdk.NewClient("https://apigateway.apigw.ntruss.com/api/v1", os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"))
	response, err := c.GETProductsProductid_TF(ctx, &ncloudsdk.PrimitiveGETProductsProductidRequest{
			Productid: plan.Productid.ValueString(),

	})

	if err != nil {
		diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}

	// Fill required attributes
	
			if !response.Product.Attributes()["description"].IsNull() || !response.Product.Attributes()["description"].IsUnknown() {
				postPlan.Description = types.StringValue(response.Product.Attributes()["description"].String())
			}

			if !response.Product.Attributes()["product_name"].IsNull() || !response.Product.Attributes()["product_name"].IsUnknown() {
				postPlan.ProductName = types.StringValue(response.Product.Attributes()["product_name"].String())
			}

			if !response.Product.Attributes()["subscription_code"].IsNull() || !response.Product.Attributes()["subscription_code"].IsUnknown() {
				postPlan.SubscriptionCode = types.StringValue(response.Product.Attributes()["subscription_code"].String())
			}

			if !response.Product.Attributes()["product"].IsNull() || !response.Product.Attributes()["product"].IsUnknown() {
				objectRes, diag := types.ObjectValueFrom(ctx, postPlan.Product.AttributeTypes(ctx), response.Product)
				if diag.HasError() {
					diagnostics.AddError("CONVERSION ERROR", "Error occured while getting object value: product")
					return
				}
				postPlan.Product = objectRes
			}

			if !response.Product.Attributes()["productid"].IsNull() || !response.Product.Attributes()["productid"].IsUnknown() {
				postPlan.Productid = types.StringValue(response.Product.Attributes()["productid"].String())
			}


	*plan = postPlan
}

func (plan *PostproductresponseModel) refreshFromOutput(ctx context.Context, diagnostics *diag.Diagnostics, id string) {

	c := ncloudsdk.NewClient("https://apigateway.apigw.ntruss.com/api/v1", os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"))
	response, err := c.GETProductsProductid_TF(ctx, &ncloudsdk.PrimitiveGETProductsProductidRequest{
			Productid: plan.Productid.ValueString(),

	})

	if err != nil {
		 diagnostics.AddError("CREATING ERROR", err.Error())
		 return
	}

	var postPlan PostproductresponseModel

	// Fill required attributes
	
			if !response.Product.Attributes()["description"].IsNull() || !response.Product.Attributes()["description"].IsUnknown() {
				postPlan.Description = types.StringValue(response.Product.Attributes()["description"].String())
			}

			if !response.Product.Attributes()["product_name"].IsNull() || !response.Product.Attributes()["product_name"].IsUnknown() {
				postPlan.ProductName = types.StringValue(response.Product.Attributes()["product_name"].String())
			}

			if !response.Product.Attributes()["subscription_code"].IsNull() || !response.Product.Attributes()["subscription_code"].IsUnknown() {
				postPlan.SubscriptionCode = types.StringValue(response.Product.Attributes()["subscription_code"].String())
			}

			if !response.Product.Attributes()["product"].IsNull() || !response.Product.Attributes()["product"].IsUnknown() {
				objectRes, diag := types.ObjectValueFrom(ctx, postPlan.Product.AttributeTypes(ctx), response.Product)
				if diag.HasError() {
					diagnostics.AddError("CONVERSION ERROR", "Error occured while getting object value: product")
					return
				}
				postPlan.Product = objectRes
			}

			if !response.Product.Attributes()["productid"].IsNull() || !response.Product.Attributes()["productid"].IsUnknown() {
				postPlan.Productid = types.StringValue(response.Product.Attributes()["productid"].String())
			}


	*plan = postPlan
}

func (plan *PostproductresponseModel) waitResourceCreated(ctx context.Context, id string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREATING"},
		Target:  []string{"CREATED"},
		Refresh: func() (interface{}, string, error) {
			c := ncloudsdk.NewClient("https://apigateway.apigw.ntruss.com/api/v1", os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"))
			response, err := c.GETProductsProductid_TF(ctx, &ncloudsdk.PrimitiveGETProductsProductidRequest{
					// need to use id
					Productid: plan.Productid.ValueString(),

			})
			if err != nil {
				return response, "CREATING", nil
			}
			if response != nil {
				return response, "CREATED", nil
			}

			return response, "CREATING", nil
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error occured while waiting for resource to be created: %s", err)
	}
	return nil
}

func (plan *PostproductresponseModel) waitResourceDeleted(ctx context.Context, id string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			c := ncloudsdk.NewClient("https://apigateway.apigw.ntruss.com/api/v1", os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"))
			response, err := c.GETProductsProductid_TF(ctx, &ncloudsdk.PrimitiveGETProductsProductidRequest{
					// need to use id
					Productid: plan.Productid.ValueString(),

			})
			if err == nil {
				return response, "DELETED", nil
			}

			return response, "DELETING", nil
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error occured while waiting for resource to be deleted: %s", err)
	}
	return nil
}
