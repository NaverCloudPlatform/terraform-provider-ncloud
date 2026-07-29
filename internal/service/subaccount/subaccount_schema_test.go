package subaccount_test

import (
	"context"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/subaccount"
)

func TestSubAccountResourceSchema(t *testing.T) {
	testResourceSchema(t, subaccount.NewSubAccountResource())
}

func TestSubAccountAccessKeyResourceSchema(t *testing.T) {
	testResourceSchema(t, subaccount.NewSubAccountAccessKeyResource())
}

func testResourceSchema(t *testing.T, r fwresource.Resource) {
	t.Helper()
	ctx := context.Background()

	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(ctx, fwresource.SchemaRequest{}, schemaResponse)
	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if diagnostics := schemaResponse.Schema.ValidateImplementation(ctx); diagnostics.HasError() {
		t.Fatalf("schema implementation diagnostics: %+v", diagnostics)
	}
}
