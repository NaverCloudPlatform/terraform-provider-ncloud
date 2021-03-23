package ncloud

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vloadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"time"
)

const (
	TargetGroupAttachmentBusyStateErrorCode            = "1200004"
	TargetGroupAttachmentPleaseTryAgainErrorCode       = "1250000"
	TargetGroupAttachmentInvalidTargetGroupNoErrorCode = "1205009"
)

func init() {
	RegisterResource("ncloud_lb_target_group_attachment", resourceNcloudLbTargetGroupAttachment())
}

func resourceNcloudLbTargetGroupAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNcloudLbTargetGroupAttachmentCreate,
		ReadContext:   resourceNcloudLbTargetGroupAttachmentRead,
		DeleteContext: resourceNcloudLbTargetGroupAttachmentDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DefaultCreateTimeout),
			Delete: schema.DefaultTimeout(DefaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			"target_group_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

	err := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		logCommonRequest("resourceNcloudLbTargetGroupAttachmentCreate", reqParams)
		resp, err := config.Client.vloadbalancer.V2Api.AddTarget(reqParams)
		if err != nil {
			errBody, _ := GetCommonErrorBody(err)
			if errBody.ReturnCode == TargetGroupAttachmentBusyStateErrorCode || errBody.ReturnCode == TargetGroupAttachmentPleaseTryAgainErrorCode {
				return resource.RetryableError(err)
			}
			logErrorResponse("resourceNcloudLbTargetGroupAttachmentCreate", err, reqParams)
			return resource.NonRetryableError(err)
		}

		logResponse("resourceNcloudLbTargetGroupAttachmentCreate", resp)
		target := getTargetFromList(resp.TargetList, ncloud.StringValue(reqParams.TargetNoList[0]))
		if target == nil {
			return resource.RetryableError(fmt.Errorf("target has not been created yet"))
		}
		return nil
	})

	if err != nil {
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
		if errorBody.ReturnCode == TargetGroupAttachmentInvalidTargetGroupNoErrorCode {
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

	err := resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		logCommonRequest("resourceNcloudLbTargetGroupAttachmentDelete", reqParams)
		resp, err := config.Client.vloadbalancer.V2Api.RemoveTarget(reqParams)
		if err != nil {
			errBody, _ := GetCommonErrorBody(err)
			if errBody.ReturnCode == TargetGroupAttachmentBusyStateErrorCode || errBody.ReturnCode == TargetGroupAttachmentPleaseTryAgainErrorCode {
				return resource.RetryableError(err)
			}
			logErrorResponse("resourceNcloudLbTargetGroupAttachmentDelete", err, reqParams)
			return resource.NonRetryableError(err)
		}
		logResponse("resourceNcloudLbTargetGroupAttachmentDelete", resp)

		target := getTargetFromList(resp.TargetList, ncloud.StringValue(reqParams.TargetNoList[0]))
		if target != nil {
			return resource.RetryableError(fmt.Errorf("target has not been removed yet"))
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func getTargetFromList(list []*vloadbalancer.Target, targetNo string) *vloadbalancer.Target {
	for _, target := range list {
		if ncloud.StringValue(target.TargetNo) == targetNo {
			return target
		}
	}
	return nil
}
