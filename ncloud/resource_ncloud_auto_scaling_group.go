package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterResource("ncloud_auto_scaling_group", resourceNcloudAutoScalingGroup())
}

func resourceNcloudAutoScalingGroup() *schema.Resource {
	return &schema.Resource{
		//Create: resourceNcloudAutoScalingGroupCreate,
		//Read:   resourceNcloudAutoScalingGroupRead,
		//Update: resourceNcloudAutoScalingGroupUpdate,
		//Delete: resourceNcloudAutoScalingGroupDelete,
		//Schema: map[string]*schema.Schema{},
	}
}

func resourceNcloudAutoScalingGroupCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	id, err := createAutoScalingGroup(d, config)
	if err != nil {
		return err
	}

	d.SetId(ncloud.StringValue(id))

	return resourceNcloudAutoScalingGroupRead(d, meta)
}

func resourceNcloudAutoScalingGroupRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceNcloudAutoScalingGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceNcloudAutoScalingGroupDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func createAutoScalingGroup(d *schema.ResourceData, config *ProviderConfig) (*string, error) {
	if config.SupportVPC {
		return createVpcAutoScalingGroup(d, config)
	} else {
		return nil, createClassicAutoScalingGroup(d, config)
	}
}

func createClassicAutoScalingGroup(d *schema.ResourceData, config *ProviderConfig) error {
	return nil
}

func createVpcAutoScalingGroup(d *schema.ResourceData, config *ProviderConfig) (*string, error) {
	return nil, nil
}
