package ncloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func Test_validateInstanceName(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "tEsting123",
			ErrCount: 1,
		},
		{
			Value:    "testing123!",
			ErrCount: 1,
		},
		{
			Value:    "1testing123",
			ErrCount: 1,
		},
		{
			Value:    "한글-123",
			ErrCount: 1,
		},
		{
			Value:    "te",
			ErrCount: 1,
		},
		{
			Value:    "testing",
			ErrCount: 0,
		},
		{
			Value:    "testing-123",
			ErrCount: 0,
		},
		{
			Value:    "testing--123",
			ErrCount: 0,
		},
		{
			Value:    "testing_123",
			ErrCount: 1,
		},
		{
			Value:    "testing123-",
			ErrCount: 1,
		},
		{
			Value:    "testing123*",
			ErrCount: 1,
		},
		{
			Value:    acctest.RandStringFromCharSet(256, acctest.CharSetAlpha),
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateInstanceName(tc.Value, "name")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Instance Name to trigger a validation error for %q", tc.Value)
		}
	}
}
