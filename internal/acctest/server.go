package acctest

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func GetTestServerName() string {
	rInt := acctest.RandIntRange(1, 9999)
	testServerName := fmt.Sprintf("tf-%d-vm", rInt)
	return testServerName
}
