package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenInt32ListToStringList(list []*int32) []*string {
	res := make([]*string, 0)
	for _, v := range list {
		res = append(res, ncloud.IntString(int(ncloud.Int32Value(v))))
	}
	return res
}

func flattenNKSClusterLogInput[T *vnks.ClusterLogInput | *vnks.AuditLogDto](logInput T) []map[string]interface{} {
	if logInput == nil {
		return nil
	}

	var audit bool
	switch v := any(logInput).(type) {
	case *vnks.ClusterLogInput:
		audit = ncloud.BoolValue(v.Audit)
	case *vnks.AuditLogDto:
		audit = ncloud.BoolValue(v.Audit)
	default:
		return nil
	}

	return []map[string]interface{}{
		{
			"audit": audit,
		},
	}
}
func expandNKSClusterLogInput[T *vnks.ClusterLogInput | *vnks.AuditLogDto](logList []interface{}, returnType T) T {
	if len(logList) == 0 {
		return nil
	}
	log := logList[0].(map[string]interface{})
	switch any(returnType).(type) {
	case *vnks.ClusterLogInput:
		return T(&vnks.ClusterLogInput{
			Audit: ncloud.Bool(log["audit"].(bool)),
		})
	case *vnks.AuditLogDto:
		return T(&vnks.AuditLogDto{
			Audit: ncloud.Bool(log["audit"].(bool)),
		})
	default:
		return nil
	}

}

func flattenNKSClusterOIDCSpec(oidcSpec *vnks.OidcRes) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if oidcSpec == nil || !*oidcSpec.Status {
		return res
	}

	res = []map[string]interface{}{
		{
			"issuer_url":      ncloud.StringValue(oidcSpec.IssuerURL),
			"client_id":       ncloud.StringValue(oidcSpec.ClientId),
			"username_claim":  ncloud.StringValue(oidcSpec.UsernameClaim),
			"username_prefix": ncloud.StringValue(oidcSpec.UsernamePrefix),
			"groups_claim":    ncloud.StringValue(oidcSpec.GroupsClaim),
			"groups_prefix":   ncloud.StringValue(oidcSpec.GroupsPrefix),
			"required_claim":  ncloud.StringValue(oidcSpec.RequiredClaim),
		},
	}
	return res
}

func expandNKSClusterOIDCSpec(oidc []interface{}) *vnks.UpdateOidcDto {
	res := &vnks.UpdateOidcDto{Status: ncloud.Bool(false)}
	if len(oidc) == 0 {
		return res
	}

	oidcSpec := oidc[0].(map[string]interface{})
	if oidcSpec["issuer_url"].(string) != "" && oidcSpec["client_id"].(string) != "" {
		res.Status = ncloud.Bool(true)
		res.IssuerURL = ncloud.String(oidcSpec["issuer_url"].(string))
		res.ClientId = ncloud.String(oidcSpec["client_id"].(string))

		usernameClaim, ok := oidcSpec["username_claim"]
		if ok {
			res.UsernameClaim = ncloud.String(usernameClaim.(string))
		}
		usernamePrefix, ok := oidcSpec["username_prefix"]
		if ok {
			res.UsernamePrefix = ncloud.String(usernamePrefix.(string))
		}
		groupsClaim, ok := oidcSpec["groups_claim"]
		if ok {
			res.GroupsClaim = ncloud.String(groupsClaim.(string))
		}
		groupsPrefix, ok := oidcSpec["groups_prefix"]
		if ok {
			res.GroupsPrefix = ncloud.String(groupsPrefix.(string))
		}
		requiredClaims, ok := oidcSpec["required_claim"]
		if ok {
			res.RequiredClaim = ncloud.String(requiredClaims.(string))
		}
	}

	return res
}

func flattenNKSClusterIPAclEntries(ipAcl *vnks.IpAclsRes) *schema.Set {

	ipAclList := schema.NewSet(schema.HashResource(resourceNcloudNKSCluster().Schema["ip_acl"].Elem.(*schema.Resource)), []interface{}{})

	for _, entry := range ipAcl.Entries {
		m := map[string]interface{}{
			"action":  *entry.Action,
			"address": *entry.Address,
		}
		if entry.Comment != nil {
			m["comment"] = *entry.Comment
		}
		ipAclList.Add(m)
	}

	return ipAclList

}

func expandNKSClusterIPAcl(acl interface{}) []*vnks.IpAclsEntriesDto {
	if acl == nil {
		return []*vnks.IpAclsEntriesDto{}
	}

	set := acl.(*schema.Set)
	res := make([]*vnks.IpAclsEntriesDto, 0)
	for _, raw := range set.List() {
		entry := raw.(map[string]interface{})

		add := &vnks.IpAclsEntriesDto{
			Address: ncloud.String(entry["address"].(string)),
			Action:  ncloud.String(entry["action"].(string)),
		}
		if comment, exist := entry["comment"].(string); exist {
			add.Comment = ncloud.String(comment)
		}
		res = append(res, add)
	}

	return res
}

func flattenNKSNodePoolTaints(taints []*vnks.NodePoolTaint) *schema.Set {

	res := schema.NewSet(schema.HashResource(resourceNcloudNKSNodePool().Schema["taint"].Elem.(*schema.Resource)), []interface{}{})

	for _, taint := range taints {
		m := map[string]interface{}{
			"key":    *taint.Key,
			"effect": *taint.Effect,
			"value":  *taint.Value,
		}
		res.Add(m)
	}

	return res

}

func expandNKSNodePoolTaints(taints interface{}) []*vnks.NodePoolTaint {
	if taints == nil {
		return nil
	}

	set := taints.(*schema.Set)
	res := make([]*vnks.NodePoolTaint, 0)
	for _, raw := range set.List() {
		taint := raw.(map[string]interface{})

		add := &vnks.NodePoolTaint{
			Key:    ncloud.String(taint["key"].(string)),
			Effect: ncloud.String(taint["effect"].(string)),
			Value:  ncloud.String(taint["value"].(string)),
		}

		res = append(res, add)
	}

	return res
}

func flattenNKSNodePoolLabels(labels []*vnks.NodePoolLabel) *schema.Set {

	res := schema.NewSet(schema.HashResource(resourceNcloudNKSNodePool().Schema["label"].Elem.(*schema.Resource)), []interface{}{})

	for _, label := range labels {
		m := map[string]interface{}{
			"key":   *label.Key,
			"value": *label.Value,
		}
		res.Add(m)
	}

	return res

}

func expandNKSNodePoolLabels(labels interface{}) []*vnks.NodePoolLabel {
	if labels == nil {
		return nil
	}

	set := labels.(*schema.Set)
	res := make([]*vnks.NodePoolLabel, 0)
	for _, raw := range set.List() {
		labels := raw.(map[string]interface{})

		add := &vnks.NodePoolLabel{
			Key:   ncloud.String(labels["key"].(string)),
			Value: ncloud.String(labels["value"].(string)),
		}

		res = append(res, add)
	}

	return res
}

func flattenNKSNodePoolAutoScale(ao *vnks.AutoscaleOption) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if ao == nil {
		return res
	}
	m := map[string]interface{}{
		"enabled": ncloud.BoolValue(ao.Enabled),
		"min":     ncloud.Int32Value(ao.Min),
		"max":     ncloud.Int32Value(ao.Max),
	}
	res = append(res, m)
	return res
}

func expandNKSNodePoolAutoScale(as []interface{}) *vnks.AutoscalerUpdate {
	if len(as) == 0 {
		return nil
	}
	autoScale := as[0].(map[string]interface{})
	return &vnks.AutoscalerUpdate{
		Enabled: ncloud.Bool(autoScale["enabled"].(bool)),
		Min:     ncloud.Int32(int32(autoScale["min"].(int))),
		Max:     ncloud.Int32(int32(autoScale["max"].(int))),
	}
}

func flattenNKSWorkerNodes(wns []*vnks.WorkerNode) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if wns == nil {
		return res
	}
	for _, wn := range wns {
		m := map[string]interface{}{
			"name":              ncloud.StringValue(wn.Name),
			"instance_no":       ncloud.Int32Value(wn.Id),
			"spec":              ncloud.StringValue(wn.ServerSpec),
			"private_ip":        ncloud.StringValue(wn.PrivateIp),
			"public_ip":         ncloud.StringValue(wn.PublicIp),
			"node_status":       ncloud.StringValue(wn.K8sStatus),
			"container_version": ncloud.StringValue(wn.DockerVersion),
			"kernel_version":    ncloud.StringValue(wn.KernelVersion),
		}
		res = append(res, m)
	}

	return res
}
