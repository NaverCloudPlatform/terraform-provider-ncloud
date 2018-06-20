package ncloud

import "testing"

func TestValidateBoolValue(t *testing.T) {
	if _, errs := validateBoolValue("true", "boolValue"); len(errs) > 0 {
		t.Fatalf("Error: %s", errs)
	}
	if _, errs := validateBoolValue("false", "boolValue"); len(errs) > 0 {
		t.Fatalf("Error: %s", errs)
	}
}

func TestValidateBoolValue_shouldReturnError(t *testing.T) {
	if _, errs := validateBoolValue("1", "boolValue"); len(errs) == 0 {
		t.Fatalf("Expected: boolValue should be true or false")
	}
	if _, errs := validateBoolValue("a", "boolValue"); len(errs) == 0 {
		t.Fatalf("Expected: boolValue should be true or false")
	}
}

func TestValidateInternetLineTypeCode(t *testing.T) {
	if _, errs := validateInternetLineTypeCode("PUBLC", "InternetLineTypeCode"); len(errs) > 0 {
		t.Fatalf("Error: %s", errs)
	}
	if _, errs := validateInternetLineTypeCode("GLBL", "InternetLineTypeCode"); len(errs) > 0 {
		t.Fatalf("Error: %s", errs)
	}
}

func TestValidateInternetLineTypeCode_shouldReturnError(t *testing.T) {
	if _, errs := validateInternetLineTypeCode("a", "InternetLineTypeCode"); len(errs) == 0 {
		t.Fatalf("Expected: InternetLineTypeCode must be one of PUBLC GLBL")
	}
}

func TestValidateServerName(t *testing.T) {
	if _, errs := validateServerName("test-server", "ServerName"); len(errs) > 0 {
		t.Fatalf("Error: %s", errs)
	}
}

func TestValidateServerName_shouldReturnError(t *testing.T) {
	if _, errs := validateServerName("1", "ServerName"); len(errs) == 0 {
		t.Fatalf("Expected: must be a valid \"ServerName\" characters between 1 and 30")
	}
	if _, errs := validateServerName("1234567890123456789012345678901", "ServerName"); len(errs) == 0 {
		t.Fatalf("Expected: must be a valid \"ServerName\" characters between 1 and 30")
	}
	if _, errs := validateServerName("!@#$", "ServerName"); len(errs) == 0 {
		t.Fatalf("Expected: server name is composed of alphabets, numbers, hyphen (-) and wild card (*). Hyphen (-) cannot be used for the last character and if wild card (*) is used, other characters cannot be input. Maximum length is 63Bytes, and the minimum is 1Byte")
	}
}

func TestValidateStringLengthInRange(t *testing.T) {
	f := validateStringLengthInRange(1, 5)

	if _, errs := f("1", "test"); len(errs) > 0 {
		t.Fatalf("Error: %s", errs)
	}
	if _, errs := f("12345", "test"); len(errs) > 0 {
		t.Fatalf("Error: %s", errs)
	}
}

func TestValidateStringLengthInRange_shouldReturnError(t *testing.T) {
	f := validateStringLengthInRange(2, 5)

	if _, errs := f("1", "test"); len(errs) == 0 {
		t.Fatalf("Expected: must be a valid \"test\" characters between 2 and 5")
	}
}

func TestValidateIntegerInRange(t *testing.T) {
	f := validateIntegerInRange(1, 5)

	if _, errs := f("1", "test"); len(errs) > 0 {
		t.Fatalf("Error: %s", errs)
	}
	if _, errs := f("4", "test"); len(errs) > 0 {
		t.Fatalf("Error: %s", errs)
	}
}

func TestValidateIntegerInRange_shouldReturnError(t *testing.T) {
	f := validateIntegerInRange(1, 5)

	if _, errs := f("0", "test"); len(errs) == 0 {
		t.Fatalf("Expected: \"test\" cannot be lower than 1: 0")
	}
	if _, errs := f("7", "test"); len(errs) == 0 {
		t.Fatalf("Expected: \"test\" cannot be higher than 5: 7")
	}
}

func TestValidateRegexp(t *testing.T) {
	if _, errs := validateRegexp("[A-Z|a-z|0-9]", "regex"); len(errs) > 0 {
		t.Fatalf("Error: %s", errs)
	}
}

func TestValidateRegexp_shouldReturnError(t *testing.T) {
	if _, errs := validateRegexp("[", "regex"); len(errs) == 0 {
		t.Fatalf("Expected: \"regex\": error parsing regexp: missing closing ]: `[`]")
	}
}

func TestValidateIncludeValues(t *testing.T) {
	f := validateIncludeValues([]string{"a", "b", "c"})
	if _, errs := f("a", "test"); len(errs) > 0 {
		t.Fatalf("Error: %s", errs)
	}

	if _, errs := f([]string{"b"}, "test"); len(errs) > 0 {
		t.Fatalf("Error: %s", errs)
	}
}

func TestValidateIncludeValues_shouldReturnError(t *testing.T) {
	f := validateIncludeValues([]string{"a", "b", "c"})
	if _, errs := f("d", "test"); len(errs) == 0 {
		t.Fatalf("test should be a or b or c")
	}

	if _, errs := f([]string{"d"}, "test"); len(errs) == 0 {
		t.Fatalf("test should be a or b or c")
	}
}
