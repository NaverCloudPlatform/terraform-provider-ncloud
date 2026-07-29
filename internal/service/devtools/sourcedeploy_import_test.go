package devtools_test

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/devtools"
)

func Test_SourceDeployStageImporter_rejects_non_numeric_project_id(t *testing.T) {
	// Given
	resource := devtools.ResourceNcloudSourceDeployStage()
	data := schema.TestResourceDataRaw(t, resource.Schema, nil)
	data.SetId("not-a-number:stage-id")

	// When
	_, err := resource.Importer.StateContext(context.Background(), data, nil)

	// Then
	if err == nil {
		t.Fatal("expected an error for a non-numeric project ID")
	}
	if !strings.Contains(err.Error(), "project ID") {
		t.Fatalf("expected a project ID error, got %q", err)
	}
}

func Test_SourceDeployStageImporter_preserves_valid_id(t *testing.T) {
	// Given
	resource := devtools.ResourceNcloudSourceDeployStage()
	data := schema.TestResourceDataRaw(t, resource.Schema, nil)
	data.SetId("123:stage-id")

	// When
	states, err := resource.Importer.StateContext(context.Background(), data, nil)

	// Then
	if err != nil {
		t.Fatalf("expected a valid ID to import, got %v", err)
	}
	if len(states) != 1 {
		t.Fatalf("expected one imported state, got %d", len(states))
	}
	if got := states[0].Get("project_id"); got != 123 {
		t.Fatalf("expected project_id 123, got %v", got)
	}
	if got := states[0].Id(); got != "stage-id" {
		t.Fatalf("expected resource ID stage-id, got %q", got)
	}
}

func Test_SourceDeployScenarioImporter_rejects_non_numeric_ids(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr string
	}{
		{
			name:    "non-numeric project ID",
			id:      "not-a-number:456:scenario-id",
			wantErr: "project ID",
		},
		{
			name:    "non-numeric stage ID",
			id:      "123:not-a-number:scenario-id",
			wantErr: "stage ID",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given
			resource := devtools.ResourceNcloudSourceDeployScenario()
			data := schema.TestResourceDataRaw(t, resource.Schema, nil)
			data.SetId(test.id)

			// When
			_, err := resource.Importer.StateContext(context.Background(), data, nil)

			// Then
			if err == nil {
				t.Fatalf("expected an error for %s", test.name)
			}
			if !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf("expected a %s error, got %q", test.wantErr, err)
			}
		})
	}
}

func Test_SourceDeployScenarioImporter_preserves_valid_id(t *testing.T) {
	// Given
	resource := devtools.ResourceNcloudSourceDeployScenario()
	data := schema.TestResourceDataRaw(t, resource.Schema, nil)
	data.SetId("123:456:scenario-id")

	// When
	states, err := resource.Importer.StateContext(context.Background(), data, nil)

	// Then
	if err != nil {
		t.Fatalf("expected a valid ID to import, got %v", err)
	}
	if len(states) != 1 {
		t.Fatalf("expected one imported state, got %d", len(states))
	}
	if got := states[0].Get("project_id"); got != 123 {
		t.Fatalf("expected project_id 123, got %v", got)
	}
	if got := states[0].Get("stage_id"); got != 456 {
		t.Fatalf("expected stage_id 456, got %v", got)
	}
	if got := states[0].Id(); got != "scenario-id" {
		t.Fatalf("expected resource ID scenario-id, got %q", got)
	}
}
