# DrinksApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**drinksPresetsGet**](DrinksApi.md#drinkspresetsget) | **GET** /drinks/presets | Get preset drinks |
| [**drinksPresetsPost**](DrinksApi.md#drinkspresetspost) | **POST** /drinks/presets | Add custom drink preset |



## drinksPresetsGet

> Array&lt;ModelsDrinkType&gt; drinksPresetsGet()

Get preset drinks

Get a list of all preset drinks

### Example

```ts
import {
  Configuration,
  DrinksApi,
} from '';
import type { DrinksPresetsGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new DrinksApi();

  try {
    const data = await api.drinksPresetsGet();
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**Array&lt;ModelsDrinkType&gt;**](ModelsDrinkType.md)

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


## drinksPresetsPost

> ModelsDrinkType drinksPresetsPost(drinkType)

Add custom drink preset

Create a new custom drink preset saving it to the database

### Example

```ts
import {
  Configuration,
  DrinksApi,
} from '';
import type { DrinksPresetsPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new DrinksApi();

  const body = {
    // ModelsDrinkType | Drink Type Details
    drinkType: ...,
  } satisfies DrinksPresetsPostRequest;

  try {
    const data = await api.drinksPresetsPost(body);
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
| **drinkType** | [ModelsDrinkType](ModelsDrinkType.md) | Drink Type Details | |

### Return type

[**ModelsDrinkType**](ModelsDrinkType.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

