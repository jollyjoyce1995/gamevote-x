# GamesApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**gamesGet**](GamesApi.md#gamesget) | **GET** /games | Search games |



## gamesGet

> Array&lt;ModelsGame&gt; gamesGet(q)

Search games

Search the cached Steam game list by name

### Example

```ts
import {
  Configuration,
  GamesApi,
} from '';
import type { GamesGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new GamesApi();

  const body = {
    // string | Search Query
    q: q_example,
  } satisfies GamesGetRequest;

  try {
    const data = await api.gamesGet(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **q** | `string` | Search Query | [Defaults to `undefined`] |

### Return type

[**Array&lt;ModelsGame&gt;**](ModelsGame.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

