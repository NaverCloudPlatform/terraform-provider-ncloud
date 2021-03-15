package ncloud

import (
	"context"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vloadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"time"
)

const (
	ErrorCodeInvalidTargetGroupNo = "1205009"
)

func init() {
	RegisterResource("ncloud_lb_target_group_attachment", resourceNcloudLbTargetGroupAttachment())
}

func resourceNcloudLbTargetGroupAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNcloudLbTargetGroupAttachmentCreate,
		ReadContext:   resourceNcloudLbTargetGroupAttachmentRead,
		DeleteContext: resourceNcloudLbTargetGroupAttachmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"target_group_no": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"target_no": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

func resourceNcloudLbTargetGroupAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_lb_target_group_attachment`"))
	}
	reqParams := &vloadbalancer.AddTargetRequest{
		RegionCode:    &config.RegionCode,
		TargetGroupNo: ncloud.String(d.Get("target_group_no").(string)),
		TargetNoList:  ncloud.StringList([]string{d.Get("target_no").(string)}),
	}
	if _, err := config.Client.vloadbalancer.V2Api.AddTarget(reqParams); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(time.Now().UTC().String())
	return nil
}

func resourceNcloudLbTargetGroupAttachmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_lb_target_group`"))
	}
	reqParams := &vloadbalancer.GetTargetListRequest{
		RegionCode:    &config.RegionCode,
		TargetGroupNo: ncloud.String(d.Get("target_group_no").(string)),
	}
	resp, err := config.Client.vloadbalancer.V2Api.GetTargetList(reqParams)
	if err != nil {
		errorBody, _ := GetCommonErrorBody(err)
		if errorBody.ReturnCode == ErrorCodeInvalidTargetGroupNo {
			log.Printf("[WARN] Target group does not exist, removing target attachment %s", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	var exist bool
	targetNo := d.Get("target_no").(string)
	for _, target := range resp.TargetList {
		if ncloud.StringValue(target.TargetNo) == targetNo {
			exist = true
			break
		}
	}

	if !exist {
		log.Printf("[WARN] Target dose not exist, removing target attachment %s", d.Id())
		d.SetId("")
	}
	return nil
}

func resourceNcloudLbTargetGroupAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_lb_target_group_attachment`"))
	}
	reqParams := &vloadbalancer.RemoveTargetRequest{
		RegionCode:    &config.RegionCode,
		TargetGroupNo: ncloud.String(d.Get("target_group_no").(string)),
		TargetNoList:  ncloud.StringList([]string{d.Get("target_no").(string)}),
	}
	if _, err := config.Client.vloadbalancer.V2Api.RemoveTarget(reqParams); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
