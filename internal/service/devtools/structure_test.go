package devtools

import "testing"

func TestExpandSourceBuildEnvVarsParams(t *testing.T) {
	envVars := []interface{}{
		map[string]interface{}{
			"key":   "key1",
			"value": "value1",
		},
		map[string]interface{}{
			"key":   "key2",
			"value": "value2",
		},
	}

	result, _ := expandSourceBuildEnvVarsParams(envVars)

	if result == nil {
		t.Fatal("result was nil")
	}

	if len(result) != 2 {
		t.Fatalf("expected result had %d elements, but got %d", 2, len(result))
	}

	env := result[0]
	if *env.Key != "key1" {
		t.Fatalf("expected result key1, but got %s", *env.Key)
	}

	if *env.Value != "value1" {
		t.Fatalf("expected result value1, but got %s", *env.Value)
	}

	env2 := result[1]
	if *env2.Key != "key2" {
		t.Fatalf("expected result key2, but got %s", *env2.Key)
	}

	if *env2.Value != "value2" {
		t.Fatalf("expected result value2, but got %s", *env2.Value)
	}
}
