# ClusterInputBody

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | ***string** | 클러스터 이름 | [default to null]
**ClusterType** | ***string** | 클러스터 타입 | [default to null]
**LoginKeyName** | ***string** | 로그인 키 이름 | [default to null]
**K8sVersion** | ***string** | 쿠버네티스 버전 | [optional] [default to null]
**RegionCode** | ***string** | Region의 코드 | [default to null]
**ZoneNo** | ***int32** | Zone 번호 | [default to null]
**VpcNo** | ***int32** | vpc의 No | [default to null]
**SubnetNoList** | **[]\*int32** | 서브넷 No 목록 | [default to null]
**SubnetLbNo** | ***int32** | 로드밸런서 전용 서브넷 No | [default to null]
**Log** | **[*ClusterLogInput](ClusterLogInput.md)** | log | [optional] [default to null]
**DefaultNodePool** | **[*DefaultNodePoolParam](DefaultNodePoolParam.md)** | 기본 노드풀 | [default to null]
**NodePool** | **[[]\*NodePool](NodePool.md)** | 추가 노드풀 | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


