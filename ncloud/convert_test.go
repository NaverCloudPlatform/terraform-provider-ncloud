package ncloud

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBase64EncodeDecode(t *testing.T) {
	data := `CreateObject("WScript.Shell").run("cmd.exe /c powershell Set-ExecutionPolicy RemoteSigned & winrm set winrm/config/service/auth @{Basic="true"} & winrm set winrm/config/service @{AllowUnencrypted="true"} & winrm quickconfig -q & sc config WinRM start= auto & winrm get winrm/config/service")`
	encoded := Base64Encode(data)
	fmt.Printf("encoded: %s\n", encoded)

	decoded := Base64Decode(encoded)

	if !reflect.DeepEqual(data, decoded) {
		t.Fatalf("Expected: %s, Actual: %s", data, decoded)
	}
}
