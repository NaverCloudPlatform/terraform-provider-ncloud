# Cluster

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Uuid** | ***string** | 클러스터 uuid | [default to null]
**AcgName** | ***string** | 클러스터 acg 이름 | [default to null]
**Name** | ***string** | 클러스터 이름 | [default to null]
**Capacity** | ***string** | 클러스터 용량 | [default to null]
**ClusterType** | ***string** | 클러스터 타입 | [default to null]
**NodeCount** | ***int32** | 등록된 노드 총 개수 | [default to null]
**NodeMaxCount** | ***int32** | 사용할 수 있는 노드의 최대 개수 | [default to null]
**CpuCount** | ***int32** | cpu 개수 | [default to null]
**MemorySize** | ***int32** | 메모리 용량 | [default to null]
**CreatedAt** | ***string** | 생성 일자 | [default to null]
**Endpoint** | ***string** | Control Plane API 주소 | [default to null]
**K8sVersion** | ***string** | 쿠버네티스 버전 | [default to null]
**RegionCode** | ***string** | region의 코드 | [default to null]
**Status** | ***string** | 클러스터의 상태 | [default to null]
**SubnetLbName** | ***string** | 로드밸런서 전용 서브넷 이름 | [default to null]
**SubnetLbNo** | ***int32** | 로드밸런서 전용 서브넷 No | [default to null]
**SubnetName** | ***string** | 서브넷 이름 | [default to null]
**SubnetNoList** | **[]\*string** | 서브넷 No 목록 | [default to null]
**UpdatedAt** | ***string** | 최근 업데이트 일자 | [default to null]
**VpcName** | ***string** | vpc 이름 | [default to null]
**VpcNo** | ***int32** | vpc 번호 | [default to null]
**ZoneNo** | ***int32** | zone 번호 | [default to null]
**LoginKeyName** | ***string** | 로그인 키 이름 | [default to null]
**Log** | **[*ClusterLogInput](ClusterLogInput.md)** | log | [optional] [default to null]
**NodePool** | **[[]\*NodePoolRes](NodePoolRes.md)** | 노드풀 | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


