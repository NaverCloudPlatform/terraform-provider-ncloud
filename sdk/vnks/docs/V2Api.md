# \V2Api

All URIs are relative to *https://nks.apigw.ntruss.com/vnks/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ClustersGet**](V2Api.md#ClustersGet) | **Get** /clusters | 
[**ClustersPost**](V2Api.md#ClustersPost) | **Post** /clusters | 
[**ClustersUuidDelete**](V2Api.md#ClustersUuidDelete) | **Delete** /clusters/{uuid} | 
[**ClustersUuidGet**](V2Api.md#ClustersUuidGet) | **Get** /clusters/{uuid} | 
[**ClustersUuidKubeconfigGet**](V2Api.md#ClustersUuidKubeconfigGet) | **Get** /clusters/{uuid}/kubeconfig | 
[**ClustersUuidKubeconfigResetPatch**](V2Api.md#ClustersUuidKubeconfigResetPatch) | **Patch** /clusters/{uuid}/kubeconfig/reset | 
[**ClustersUuidNodePoolGet**](V2Api.md#ClustersUuidNodePoolGet) | **Get** /clusters/{uuid}/node-pool | 
[**ClustersUuidNodePoolInstanceNoDelete**](V2Api.md#ClustersUuidNodePoolInstanceNoDelete) | **Delete** /clusters/{uuid}/node-pool/{instanceNo} | 
[**ClustersUuidNodePoolInstanceNoPatch**](V2Api.md#ClustersUuidNodePoolInstanceNoPatch) | **Patch** /clusters/{uuid}/node-pool/{instanceNo} | 
[**ClustersUuidNodePoolPost**](V2Api.md#ClustersUuidNodePoolPost) | **Post** /clusters/{uuid}/node-pool | 
[**ClustersUuidNodesGet**](V2Api.md#ClustersUuidNodesGet) | **Get** /clusters/{uuid}/nodes | 
[**OptionVersionGet**](V2Api.md#OptionVersionGet) | **Get** /option/version | 
[**RootGet**](V2Api.md#RootGet) | **Get** / | 


# **ClustersGet**
> ClustersRes ClustersGet()


### Required Parameters
This endpoint does not need any parameter.

### Return type

*[**ClustersRes**](ClustersRes.md)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ClustersPost**
> CreateClusterRes ClustersPost(clusterInputBody)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**clusterInputBody** | **[\*ClusterInputBody](ClusterInputBody.md)** |  | 

### Return type

*[**CreateClusterRes**](CreateClusterRes.md)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ClustersUuidDelete**
> ClustersUuidDelete(uuid)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**uuid** | **string** | uuid | 

### Return type

 (empty response body)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ClustersUuidGet**
> ClusterRes ClustersUuidGet(uuid)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**uuid** | **string** | uuid | 

### Return type

*[**ClusterRes**](ClusterRes.md)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ClustersUuidKubeconfigGet**
> KubeconfigRes ClustersUuidKubeconfigGet(uuid)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**uuid** | **string** | uuid | 

### Return type

*[**KubeconfigRes**](KubeconfigRes.md)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ClustersUuidKubeconfigResetPatch**
> ClustersUuidKubeconfigResetPatch(uuid)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**uuid** | **string** | uuid | 

### Return type

 (empty response body)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ClustersUuidNodePoolGet**
> NodePoolsRes ClustersUuidNodePoolGet(uuid)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**uuid** | **string** | uuid | 

### Return type

*[**NodePoolsRes**](NodePoolsRes.md)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ClustersUuidNodePoolInstanceNoDelete**
> ClustersUuidNodePoolInstanceNoDelete(uuid, instanceNo)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**uuid** | **string** | uuid | **instanceNo** | **string** | instanceNo | 

### Return type

 (empty response body)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ClustersUuidNodePoolInstanceNoPatch**
> ClustersUuidNodePoolInstanceNoPatch(nodePoolUpdateBody, uuid, instanceNo)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**nodePoolUpdateBody** | **[\*NodePoolUpdateBody](NodePoolUpdateBody.md)** |  | **uuid** | **string** | uuid | **instanceNo** | **string** | instanceNo | 

### Return type

 (empty response body)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ClustersUuidNodePoolPost**
> ClustersUuidNodePoolPost(nodePoolCreationBody, uuid)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**nodePoolCreationBody** | **[\*NodePoolCreationBody](NodePoolCreationBody.md)** |  | **uuid** | **string** | uuid | 

### Return type

 (empty response body)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ClustersUuidNodesGet**
> WorkerNodeRes ClustersUuidNodesGet(uuid)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**uuid** | **string** | uuid | 

### Return type

*[**WorkerNodeRes**](WorkerNodeRes.md)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **OptionVersionGet**
> OptionsRes OptionVersionGet()


### Required Parameters
This endpoint does not need any parameter.

### Return type

*[**OptionsRes**](OptionsRes.md)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **RootGet**
> RootGet()


### Required Parameters
This endpoint does not need any parameter.

### Return type

 (empty response body)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

