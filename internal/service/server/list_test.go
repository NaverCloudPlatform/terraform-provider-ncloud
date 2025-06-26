package server

import (
	"testing"
)

func TestExpandBlockDevicePartitionListParams(t *testing.T) {
	blockDevicePartitions := []interface{}{
		map[string]interface{}{
			"mount_point":    "/data",
			"partition_size": "100",
		},
		map[string]interface{}{
			"mount_point":    "/backup",
			"partition_size": "200",
		},
	}

	result, err := expandBlockDevicePartitionListParams(blockDevicePartitions)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("result was nil")
	}

	if len(result) != 2 {
		t.Fatalf("expected result to have %d elements, but got %d", 2, len(result))
	}

	partition := result[0]
	if *partition.MountPoint != "/data" {
		t.Fatalf("expected MountPoint to be /data, but got %s", *partition.MountPoint)
	}

	if *partition.PartitionSize != "100" {
		t.Fatalf("expected PartitionSize to be 100, but got %s", *partition.PartitionSize)
	}

	partition = result[1]
	if *partition.MountPoint != "/backup" {
		t.Fatalf("expected MountPoint to be /backup, but got %s", *partition.MountPoint)
	}

	if *partition.PartitionSize != "200" {
		t.Fatalf("expected PartitionSize to be 200, but got %s", *partition.PartitionSize)
	}
}

func TestExpandBlockDevicePartitionListParams_EmptyInput(t *testing.T) {
	blockDevicePartitions := []interface{}{}

	result, err := expandBlockDevicePartitionListParams(blockDevicePartitions)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("result was nil")
	}

	if len(result) != 0 {
		t.Fatalf("expected result to have %d elements, but got %d", 0, len(result))
	}
}

func TestExpandBlockDevicePartitionListParams_InvalidInput(t *testing.T) {
	blockDevicePartitions := []interface{}{
		map[string]interface{}{
			"invalid_key": "value",
		},
	}

	result, err := expandBlockDevicePartitionListParams(blockDevicePartitions)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("result was nil")
	}

	if len(result) != 1 {
		t.Fatalf("expected result to have %d elements, but got %d", 1, len(result))
	}

	partition := result[0]
	if partition.MountPoint != nil {
		t.Fatalf("expected MountPoint to be nil, but got %s", *partition.MountPoint)
	}

	if partition.PartitionSize != nil {
		t.Fatalf("expected PartitionSize to be nil, but got %s", *partition.PartitionSize)
	}
}
