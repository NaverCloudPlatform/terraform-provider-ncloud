package server

import (
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
)

func TestNeedsServerStopBeforeDetaching(t *testing.T) {
	cases := []struct {
		name     string
		instance *ServerInstance
		want     bool
	}{
		{"server does not exist", nil, false},
		{"server already stopped", &ServerInstance{ServerInstanceStatus: ncloud.String("NSTOP")}, false},
		{"server running", &ServerInstance{ServerInstanceStatus: ncloud.String("RUN")}, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := needsServerStopBeforeDetaching(tc.instance); got != tc.want {
				t.Errorf("needsServerStopBeforeDetaching() = %v, want %v", got, tc.want)
			}
		})
	}
}
