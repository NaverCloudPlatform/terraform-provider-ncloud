package nks

import (
	"reflect"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestFlattenInt32ListToStringList(t *testing.T) {
	initialList := []*int32{
		ncloud.Int32(int32(1111)),
		ncloud.Int32(int32(2222)),
		ncloud.Int32(int32(3333)),
	}

	stringList := flattenInt32ListToStringList(initialList)
	expected := []*string{
		ncloud.String("1111"),
		ncloud.String("2222"),
		ncloud.String("3333")}
	if !reflect.DeepEqual(stringList, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			stringList,
			expected)
	}
}

func TestFlattenNKSClusterLogInput(t *testing.T) {
	logInput := &vnks.ClusterLogInput{Audit: ncloud.Bool(true)}

	result := flattenNKSClusterLogInput(logInput)

	if result == nil {
		t.Fatal("result was nil")
	}

	r := result[0]
	if r["audit"].(bool) != true {
		t.Fatalf("expected result enabled to be true, but was %v", r["enabled"])
	}
}

func TestExpandNKSClusterLogInput(t *testing.T) {
	log := []interface{}{
		map[string]interface{}{
			"audit": false,
		},
	}

	result := expandNKSClusterLogInput(log, &vnks.AuditLogDto{})

	if result == nil {
		t.Fatal("result was nil")
	}

	if ncloud.BoolValue(result.Audit) != false {
		t.Fatalf("expected false , but got %v", ncloud.BoolValue(result.Audit))
	}
}

func TestFlattenNKSClusterOIDCSpec(t *testing.T) {
	oidcSpec := &vnks.OidcRes{
		Status:         ncloud.Bool(true),
		UsernameClaim:  ncloud.String("email"),
		UsernamePrefix: ncloud.String("username:"),
		IssuerURL:      ncloud.String("https://sso.ntruss.com/iss"),
		ClientId:       ncloud.String("testClient"),
		GroupsPrefix:   ncloud.String("groups:"),
		GroupsClaim:    ncloud.String("group"),
		RequiredClaim:  ncloud.String("iss=https://sso.ntruss.com/iss"),
	}

	result := flattenNKSClusterOIDCSpec(oidcSpec)

	if len(result) == 0 {
		t.Fatal("empty result")
	}

	r := result[0]

	if r["username_claim"].(string) != "email" {
		t.Fatalf("expected result username_claim to be 'email', but was %v", r["username_claim"])
	}

	if r["username_prefix"].(string) != "username:" {
		t.Fatalf("expected result username_prefix to be 'username:', but was %v", r["username_prefix"])
	}

	if r["issuer_url"].(string) != "https://sso.ntruss.com/iss" {
		t.Fatalf("expected result issuer_url to be 'https://sso.ntruss.com/iss', but was %v", r["issuer_url"])
	}

	if r["client_id"].(string) != "testClient" {
		t.Fatalf("expected result client_id to be 'testClient', but was %v", r["client_id"])
	}

	if r["groups_claim"].(string) != "group" {
		t.Fatalf("expected result groups_claim to be 'group', but was %v", r["groups_claim"])
	}

	if r["groups_prefix"].(string) != "groups:" {
		t.Fatalf("expected result groups_prefix to be 'groups:', but was %v", r["groups_prefix"])
	}

	if r["required_claim"].(string) != "iss=https://sso.ntruss.com/iss" {
		t.Fatalf("expected result groups_prefix to be 'iss=https://sso.ntruss.com/iss', but was %v", r["required_claim"])
	}
}

func TestExpandNKSClusterOIDCSpec(t *testing.T) {
	oidc := []interface{}{
		map[string]interface{}{
			"issuer_url":      "https://sso.ntruss.com/iss",
			"client_id":       "testClient",
			"username_claim":  "email",
			"username_prefix": "username:",
			"groups_claim":    "group",
			"groups_prefix":   "groups:",
			"required_claim":  "iss=https://sso.ntruss.com/iss",
		},
	}

	result := expandNKSClusterOIDCSpec(oidc)

	if result == nil {
		t.Fatal("result was nil")
	}

	expected := &vnks.UpdateOidcDto{
		Status:         ncloud.Bool(true),
		IssuerURL:      ncloud.String("https://sso.ntruss.com/iss"),
		ClientId:       ncloud.String("testClient"),
		UsernameClaim:  ncloud.String("email"),
		UsernamePrefix: ncloud.String("username:"),
		GroupsClaim:    ncloud.String("group"),
		GroupsPrefix:   ncloud.String("groups:"),
		RequiredClaim:  ncloud.String("iss=https://sso.ntruss.com/iss"),
	}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected %v , but got %v", expected, result)
	}
}

func TestFlattenNKSClusterIPAcl(t *testing.T) {
	ipAcl := &vnks.IpAclsRes{
		DefaultAction: ncloud.String("deny"),
		Entries: []*vnks.IpAclsEntriesRes{
			{Address: ncloud.String("10.0.1.0/24"),
				Action:  ncloud.String("allow"),
				Comment: ncloud.String("master ip"),
			},
		},
	}

	result := flattenNKSClusterIPAclEntries(ipAcl)

	if len(result.List()) == 0 {
		t.Fatal("empty result")
	}

	r := result.List()[0]
	rr := r.(map[string]interface{})
	if rr["address"].(string) != "10.0.1.0/24" {
		t.Fatalf("expected result address to be '10.0.1.0/24', but was %v", rr["address"])
	}

	if rr["action"].(string) != "allow" {
		t.Fatalf("expected result action to be 'allow', but was %v", rr["action"])
	}

	if rr["comment"].(string) != "master ip" {
		t.Fatalf("expected result comment to be 'master ip', but was %v", rr["comment"])
	}
}

func TestExpandNKSClusterIPAcl(t *testing.T) {
	ipAclList := schema.NewSet(schema.HashResource(ResourceNcloudNKSCluster().Schema["ip_acl"].Elem.(*schema.Resource)), []interface{}{})

	ipAclList.Add(map[string]interface{}{
		"action":  "allow",
		"address": "10.0.1.0/24",
		"comment": "master ip",
	})

	result := expandNKSClusterIPAcl(ipAclList)

	if result == nil {
		t.Fatal("result was nil")
	}

	expected := []*vnks.IpAclsEntriesDto{
		{
			Address: ncloud.String("10.0.1.0/24"),
			Action:  ncloud.String("allow"),
			Comment: ncloud.String("master ip"),
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected %v , but got %v", expected, result)
	}
}

func TestFlattenNKSNodePoolTaints(t *testing.T) {

	taints := []*vnks.NodePoolTaint{
		{
			Key:    ncloud.String("foo"),
			Value:  ncloud.String("bar"),
			Effect: ncloud.String("NoExecute"),
		},
		{
			Key:    ncloud.String("bar"),
			Value:  ncloud.String(""),
			Effect: ncloud.String("NoSchedule"),
		},
	}

	result := flattenNKSNodePoolTaints(taints)

	if len(result.List()) == 0 {
		t.Fatal("empty result")
	}

	r := result.List()[0]
	rr := r.(map[string]interface{})
	if rr["key"].(string) != "foo" {
		t.Fatalf("expected result key to be 'foo', but was %v", rr["key"])
	}

	if rr["value"].(string) != "bar" {
		t.Fatalf("expected result value to be 'bar', but was %v", rr["value"])
	}

	if rr["effect"].(string) != "NoExecute" {
		t.Fatalf("expected result effect to be 'NoExecute', but was %v", rr["effect"])
	}

	r = result.List()[1]
	rr = r.(map[string]interface{})
	if rr["key"].(string) != "bar" {
		t.Fatalf("expected result key to be 'bar', but was %v", rr["key"])
	}

	if rr["value"].(string) != "" {
		t.Fatalf("expected result value to be '', but was %v", rr["value"])
	}

	if rr["effect"].(string) != "NoSchedule" {
		t.Fatalf("expected result effect to be 'NoSchedule', but was %v", rr["effect"])
	}

}

func TestExpandNKSNodePoolTaints(t *testing.T) {
	taints := schema.NewSet(schema.HashResource(ResourceNcloudNKSNodePool().Schema["taint"].Elem.(*schema.Resource)), []interface{}{})

	taints.Add(map[string]interface{}{
		"key":    "foo",
		"value":  "bar",
		"effect": "NoExecute",
	})
	taints.Add(map[string]interface{}{
		"key":    "bar",
		"value":  "",
		"effect": "NoSchedule",
	})

	result := expandNKSNodePoolTaints(taints)

	if result == nil {
		t.Fatal("result was nil")
	}

	expected := []*vnks.NodePoolTaint{
		{
			Key:    ncloud.String("foo"),
			Value:  ncloud.String("bar"),
			Effect: ncloud.String("NoExecute"),
		},
		{
			Key:    ncloud.String("bar"),
			Value:  ncloud.String(""),
			Effect: ncloud.String("NoSchedule"),
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected %v , but got %v", expected, result)
	}
}

func TestFlattenNKSNodePoolLabels(t *testing.T) {

	labels := []*vnks.NodePoolLabel{
		{
			Key:   ncloud.String("foo"),
			Value: ncloud.String("bar"),
		},
		{
			Key:   ncloud.String("bar"),
			Value: ncloud.String("foo"),
		},
	}

	result := flattenNKSNodePoolLabels(labels)

	if len(result.List()) == 0 {
		t.Fatal("empty result")
	}

	r := result.List()[0]
	rr := r.(map[string]interface{})
	if rr["key"].(string) != "foo" {
		t.Fatalf("expected result key to be 'foo', but was %v", rr["key"])
	}

	if rr["value"].(string) != "bar" {
		t.Fatalf("expected result value to be 'bar', but was %v", rr["value"])
	}

	r = result.List()[1]
	rr = r.(map[string]interface{})
	if rr["key"].(string) != "bar" {
		t.Fatalf("expected result key to be 'bar', but was %v", rr["key"])
	}

	if rr["value"].(string) != "foo" {
		t.Fatalf("expected result value to be 'foo', but was %v", rr["value"])
	}

}

func TestExpandNKSNodePoolLabels(t *testing.T) {
	labels := schema.NewSet(schema.HashResource(ResourceNcloudNKSNodePool().Schema["label"].Elem.(*schema.Resource)), []interface{}{})

	labels.Add(map[string]interface{}{
		"key":   "foo",
		"value": "bar",
	})
	labels.Add(map[string]interface{}{
		"key":   "bar",
		"value": "foo",
	})

	result := expandNKSNodePoolLabels(labels)

	if result == nil {
		t.Fatal("result was nil")
	}

	expected := []*vnks.NodePoolLabel{
		{
			Key:   ncloud.String("foo"),
			Value: ncloud.String("bar"),
		},
		{
			Key:   ncloud.String("bar"),
			Value: ncloud.String("foo"),
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected %v , but got %v", expected, result)
	}
}

func TestFlattenNKSNodePoolAutoscale(t *testing.T) {
	expanded := &vnks.AutoscaleOption{
		Enabled: ncloud.Bool(true),
		Max:     ncloud.Int32(2),
		Min:     ncloud.Int32(2),
	}

	result := flattenNKSNodePoolAutoScale(expanded)

	if result == nil {
		t.Fatal("result was nil")
	}

	r := result[0]
	if r["enabled"].(bool) != true {
		t.Fatalf("expected result enabled to be true, but was %v", r["enabled"])
	}

	if r["min"].(int32) != 2 {
		t.Fatalf("expected result min to be 2, but was %d", r["min"])
	}

	if r["max"].(int32) != 2 {
		t.Fatalf("expected result max to be 2, but was %d", r["max"])
	}
}

func TestFlattenNKSWorkerNodes(t *testing.T) {
	expanded := []*vnks.WorkerNode{
		{
			Id:            ncloud.Int32(1),
			Name:          ncloud.String("node1"),
			ServerSpec:    ncloud.String("[Standard] vCPU 2EA, Memory 8GB"),
			PrivateIp:     ncloud.String("10.0.1.4"),
			PublicIp:      ncloud.String(""),
			K8sStatus:     ncloud.String("Ready"),
			DockerVersion: ncloud.String("containerd://1.3.7"),
			KernelVersion: ncloud.String("5.4.0-65-generic"),
		},
	}

	result := flattenNKSWorkerNodes(expanded)

	if result == nil {
		t.Fatal("result was nil")
	}

	r := result[0]
	if r["instance_no"].(int32) != 1 {
		t.Fatalf("expected result instance_no to be 1, but was %v", r["instance_no"])
	}

	if r["name"].(string) != "node1" {
		t.Fatalf("expected result name to be node1, but was %s", r["name"])
	}

	if r["spec"].(string) != "[Standard] vCPU 2EA, Memory 8GB" {
		t.Fatalf("expected result spec to be [Standard] vCPU 2EA, Memory 8GB, but was %s", r["spec"])
	}

	if r["private_ip"].(string) != "10.0.1.4" {
		t.Fatalf("expected result private_ip to be 10.0.1.4, but was %s", r["private_ip"])
	}

	if r["public_ip"].(string) != "" {
		t.Fatalf("expected result public_ip to be emtpy, but was %s", r["public_ip"])
	}

	if r["node_status"].(string) != "Ready" {
		t.Fatalf("expected result node_status to be Ready, but was %s", r["node_status"])
	}

	if r["container_version"].(string) != "containerd://1.3.7" {
		t.Fatalf("expected result container_version to be containerd://1.3.7, but was %s", r["container_version"])
	}

	if r["kernel_version"].(string) != "5.4.0-65-generic" {
		t.Fatalf("expected result kernel_version to be 5.4.0-65-generic, but was %s", r["kernel_version"])
	}
}

func TestExpandNKSNodePoolAutoScale(t *testing.T) {
	autoscaleList := []interface{}{
		map[string]interface{}{
			"enabled": true,
			"min":     2,
			"max":     2,
		},
	}

	result := expandNKSNodePoolAutoScale(autoscaleList)

	if result == nil {
		t.Fatal("result was nil")
	}

	if ncloud.BoolValue(result.Enabled) != true {
		t.Fatalf("expected result true, but got %v", ncloud.BoolValue(result.Enabled))
	}

	if ncloud.Int32Value(result.Min) != int32(2) {
		t.Fatalf("expected result 2, but got %d", ncloud.Int32Value(result.Min))
	}

	if ncloud.Int32Value(result.Max) != int32(2) {
		t.Fatalf("expected result 2, but got %d", ncloud.Int32Value(result.Max))
	}
}
