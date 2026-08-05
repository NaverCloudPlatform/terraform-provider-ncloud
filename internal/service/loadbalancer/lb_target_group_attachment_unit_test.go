package loadbalancer

import (
	"reflect"
	"testing"
)

func TestParseLbTargetGroupAttachmentID(t *testing.T) {
	tests := []struct {
		name              string
		id                string
		wantTargetGroupNo string
		wantTargetNoList  []string
		wantErr           bool
	}{
		{
			name:              "single target",
			id:                "12345:23456",
			wantTargetGroupNo: "12345",
			wantTargetNoList:  []string{"23456"},
		},
		{
			name:              "multiple targets",
			id:                "12345:23456,34567",
			wantTargetGroupNo: "12345",
			wantTargetNoList:  []string{"23456", "34567"},
		},
		{
			name:              "trims whitespace",
			id:                " 12345 : 23456, 34567 ",
			wantTargetGroupNo: "12345",
			wantTargetNoList:  []string{"23456", "34567"},
		},
		{
			name:    "missing separator",
			id:      "12345",
			wantErr: true,
		},
		{
			name:    "missing target group",
			id:      ":23456",
			wantErr: true,
		},
		{
			name:    "missing target",
			id:      "12345:",
			wantErr: true,
		},
		{
			name:    "empty target in list",
			id:      "12345:23456,",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetGroupNo, targetNoList, err := parseLbTargetGroupAttachmentID(tt.id)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if targetGroupNo != tt.wantTargetGroupNo {
				t.Fatalf("target group no = %q, want %q", targetGroupNo, tt.wantTargetGroupNo)
			}
			if !reflect.DeepEqual(targetNoList, tt.wantTargetNoList) {
				t.Fatalf("target no list = %#v, want %#v", targetNoList, tt.wantTargetNoList)
			}
		})
	}
}

func TestLbTargetGroupAttachmentID(t *testing.T) {
	got := lbTargetGroupAttachmentID("12345", []string{"23456", "34567"})
	want := "12345:23456,34567"
	if got != want {
		t.Fatalf("id = %q, want %q", got, want)
	}
}
