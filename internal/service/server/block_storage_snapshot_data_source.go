package server

import (
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/provider"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

func init() {
	RegisterDataSource("ncloud_block_storage_snapshot", dataSourceNcloudBlockStorageSnapshot())
}

func dataSourceNcloudBlockStorageSnapshot() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudBlockStorageSnapshotRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"block_storage_no": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"snapshot_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filter": DataSourceFiltersSchema(),
		},
	}
}

func dataSourceNcloudBlockStorageSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	instances, err := GetBlockStorageSnapshot(d, config)
	if err != nil {
		return err
	}

	if len(instances) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	resources := ConvertToArrayMap(instances)
	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudBlockStorageSnapshot().Schema)
	}

	if err := ValidateOneResult(len(resources)); err != nil {
		return err
	}

	d.SetId(resources[0]["snapshot_no"].(string))
	SetSingularResourceDataFromMap(d, resources[0])

	return nil
}

func GetBlockStorageSnapshot(d *schema.ResourceData, config *ProviderConfig) ([]*BlockStorageSnapshot, error) {
	if config.SupportVPC {
		return getVpcBlockStorageSnapshot(d, config)
	}

	return getClassicBlockStorageSnapshot(d, config)
}

func getClassicBlockStorageSnapshot(d *schema.ResourceData, config *ProviderConfig) ([]*BlockStorageSnapshot, error) {
	regionNo, err := ParseRegionNoParameter(d)
	if err != nil {
		return nil, err
	}

	reqParams := &server.GetBlockStorageSnapshotInstanceListRequest{
		RegionNo: regionNo,
	}

	if v, ok := d.GetOk("block_storage_no"); ok {
		reqParams.OriginalBlockStorageInstanceNoList = []*string{ncloud.String(v.(string))}
	}

	if v, ok := d.GetOk("id"); ok {
		reqParams.BlockStorageSnapshotInstanceNoList = []*string{ncloud.String(v.(string))}
	}

	LogCommonRequest("getClassicBlockStorageSnapshot", reqParams)
	resp, err := config.Client.Server.V2Api.GetBlockStorageSnapshotInstanceList(reqParams)
	if err != nil {
		LogErrorResponse("getClassicBlockStorageSnapshot", err, reqParams)
		return nil, err
	}
	LogResponse("getClassicBlockStorageSnapshot", resp)

	var list []*BlockStorageSnapshot
	for _, r := range resp.BlockStorageSnapshotInstanceList {
		list = append(list, convertClassicSnapshotInstance(r))
	}

	return list, nil
}

func convertClassicSnapshotInstance(r *server.BlockStorageSnapshotInstance) *BlockStorageSnapshot {
	if r == nil {
		return nil
	}

	return &BlockStorageSnapshot{
		SnapshotNo:             r.BlockStorageSnapshotInstanceNo,
		Name:                   r.BlockStorageSnapshotName,
		VolumeSize:             r.BlockStorageSnapshotVolumeSize,
		BlockStorageInstanceNo: r.OriginalBlockStorageInstanceNo,
		Description:            r.BlockStorageSnapshotInstanceDescription,
	}
}

func getVpcBlockStorageSnapshot(d *schema.ResourceData, config *ProviderConfig) ([]*BlockStorageSnapshot, error) {
	reqParams := &vserver.GetBlockStorageSnapshotInstanceListRequest{
		RegionCode: &config.RegionCode,
	}

	if v, ok := d.GetOk("block_storage_no"); ok {
		reqParams.OriginalBlockStorageInstanceNoList = []*string{ncloud.String(v.(string))}
	}

	if v, ok := d.GetOk("id"); ok {
		reqParams.BlockStorageSnapshotInstanceNoList = []*string{ncloud.String(v.(string))}
	}

	LogCommonRequest("getVpcBlockStorageSnapshot", reqParams)
	resp, err := config.Client.Vserver.V2Api.GetBlockStorageSnapshotInstanceList(reqParams)
	if err != nil {
		LogErrorResponse("getVpcBlockStorageSnapshot", err, reqParams)
		return nil, err
	}
	LogResponse("getVpcBlockStorageSnapshot", resp)

	var list []*BlockStorageSnapshot
	for _, r := range resp.BlockStorageSnapshotInstanceList {
		list = append(list, convertVpcSnapshotInstance(r))
	}

	return list, nil
}

func convertVpcSnapshotInstance(r *vserver.BlockStorageSnapshotInstance) *BlockStorageSnapshot {
	if r == nil {
		return nil
	}

	return &BlockStorageSnapshot{
		SnapshotNo:             r.BlockStorageSnapshotInstanceNo,
		Name:                   r.BlockStorageSnapshotName,
		VolumeSize:             r.BlockStorageSnapshotVolumeSize,
		BlockStorageInstanceNo: r.OriginalBlockStorageInstanceNo,
		Description:            r.BlockStorageSnapshotDescription,
	}
}

// BlockStorageSnapshot Dto for block storage snapshot
type BlockStorageSnapshot struct {
	SnapshotNo             *string `json:"snapshot_no,omitempty"`
	Name                   *string `json:"name,omitempty"`
	VolumeSize             *int64  `json:"volume_size,omitempty"`
	BlockStorageInstanceNo *string `json:"block_storage_no,omitempty"`
	Description            *string `json:"description,omitempty"`
}
