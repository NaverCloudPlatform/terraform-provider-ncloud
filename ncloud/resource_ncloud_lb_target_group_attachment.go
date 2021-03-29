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
		UpdateContext: resourceNcloudLbTargetGroupAttachmentUpdate,
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
			"target_no_list": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
		TargetNoList:  ncloud.StringInterfaceList(d.Get("target_no_list").([]interface{})),
	}

	err := waitForAddTarget(ctx, d, config, reqParams)

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

	targetNoList, err := getVpcLoadBalancerTargetGroupAttachment(config, d.Get("target_group_no").(string), ncloud.StringListValue(ncloud.StringInterfaceList(d.Get("target_no_list").([]interface{}))))
	if err != nil {
		errorBody, _ := GetCommonErrorBody(err)
		if errorBody.ReturnCode == TargetGroupAttachmentInvalidTargetGroupNoErrorCode {
			log.Printf("[WARN] Target group does not exist, removing target attachment %s", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if targetNoList == nil {
		log.Printf("[WARN] Target dose not exist, removing target attachment %s", d.Id())
		d.SetId("")
	}

	d.Set("target_no_list", targetNoList)
	return nil
}

func resourceNcloudLbTargetGroupAttachmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_lb_target_group`"))
	}
	if d.HasChange("target_no_list") {
		o, n := d.GetChange("target_no_list")
		oldTargetNoList := ncloud.StringInterfaceList(o.([]interface{}))
		newTargetNoList := ncloud.StringInterfaceList(n.([]interface{}))

		oldTargetNoMap := make(map[string]bool)
		newTargetNoMap := make(map[string]bool)

		for _, oldTargetNo := range oldTargetNoList {
			oldTargetNoMap[*oldTargetNo] = true
		}

		for _, newTargetNo := range newTargetNoList {
			newTargetNoMap[*newTargetNo] = true
		}

		removeTargetNoList := make([]string, 0)
		addTargetNoList := make([]string, 0)

		for key, _ := range newTargetNoMap {
			if oldTargetNoMap[key] {
				delete(oldTargetNoMap, key)
			} else {
				addTargetNoList = append(addTargetNoList, key)
			}
		}

		for key, _ := range oldTargetNoMap {
			removeTargetNoList = append(removeTargetNoList, key)
		}

		if len(addTargetNoList) >= 1 {
			addReqParams := &vloadbalancer.AddTargetRequest{
				RegionCode:    &config.RegionCode,
				TargetGroupNo: ncloud.String(d.Get("target_group_no").(string)),
				TargetNoList:  ncloud.StringList(addTargetNoList),
			}

			addErr := waitForAddTarget(ctx, d, config, addReqParams)

			if addErr != nil {
				return diag.FromErr(addErr)
			}
		}

		if len(removeTargetNoList) >= 1 {
			removeReqParams := &vloadbalancer.RemoveTargetRequest{
				RegionCode:    &config.RegionCode,
				TargetGroupNo: ncloud.String(d.Get("target_group_no").(string)),
				TargetNoList:  ncloud.StringList(removeTargetNoList),
			}

			removeErr := waitForRemoveTarget(ctx, d, config, removeReqParams)

			if removeErr != nil {
				return diag.FromErr(removeErr)
			}
		}
	}
	return resourceNcloudLbTargetGroupAttachmentRead(ctx, d, config)
}

func resourceNcloudLbTargetGroupAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_lb_target_group_attachment`"))
	}
	reqParams := &vloadbalancer.RemoveTargetRequest{
		RegionCode:    &config.RegionCode,
		TargetGroupNo: ncloud.String(d.Get("target_group_no").(string)),
		TargetNoList:  ncloud.StringInterfaceList(d.Get("target_no_list").([]interface{})),
	}

	err := waitForRemoveTarget(ctx, d, config, reqParams)

	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func getVpcLoadBalancerTargetGroupAttachment(config *ProviderConfig, targetGroupNo string, targetNoList []string) ([]string, error) {
	reqParams := &vloadbalancer.GetTargetListRequest{
		RegionCode:    &config.RegionCode,
		TargetGroupNo: ncloud.String(targetGroupNo),
	}
	resp, err := config.Client.vloadbalancer.V2Api.GetTargetList(reqParams)
	if err != nil {
		return nil, err
	}
	matchTargetNoList := getMatchTargetNoListFromResponse(resp.TargetList, targetNoList)

	if len(matchTargetNoList) < 1 {
		return nil, nil
	}

	return matchTargetNoList, nil
}

func getMatchTargetNoListFromResponse(respTargetList []*vloadbalancer.Target, targetNoList []string) []string {
	matchTargetNoList := make([]string, 0)
	respTargetNoList := make([]string, 0)

	for _, respTarget := range respTargetList {
		respTargetNoList = append(respTargetNoList, ncloud.StringValue(respTarget.TargetNo))
	}

	for _, targetNo := range targetNoList {
		if containsInStringList(respTargetNoList, targetNo) {
			matchTargetNoList = append(matchTargetNoList, targetNo)
		}
	}

	return matchTargetNoList
}

func waitForAddTarget(ctx context.Context, d *schema.ResourceData, config *ProviderConfig, reqParams *vloadbalancer.AddTargetRequest) error {
	return resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
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
		return nil
	})
}

func waitForRemoveTarget(ctx context.Context, d *schema.ResourceData, config *ProviderConfig, reqParams *vloadbalancer.RemoveTargetRequest) error {
	return resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
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

		matchTargetNoList := getMatchTargetNoListFromResponse(resp.TargetList, ncloud.StringListValue(reqParams.TargetNoList))
		if len(matchTargetNoList) > 0 {
			return resource.RetryableError(fmt.Errorf("target has not been removed yet"))
		}
		return nil
	})
}
