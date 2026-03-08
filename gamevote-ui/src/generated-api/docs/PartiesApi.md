# PartiesApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**partiesCodeAttendeesAttendeeIdDelete**](PartiesApi.md#partiescodeattendeesattendeeiddelete) | **DELETE** /parties/{code}/attendees/{attendeeId} | Delete an attendee |
| [**partiesCodeAttendeesGet**](PartiesApi.md#partiescodeattendeesget) | **GET** /parties/{code}/attendees | Get party attendees |
| [**partiesCodeAttendeesPost**](PartiesApi.md#partiescodeattendeespost) | **POST** /parties/{code}/attendees | Add an attendee |
| [**partiesCodeBeersPost**](PartiesApi.md#partiescodebeerspost) | **POST** /parties/{code}/beers | Add a beer |
| [**partiesCodeGet**](PartiesApi.md#partiescodeget) | **GET** /parties/{code} | Get a party |
| [**partiesCodeOptionsGet**](PartiesApi.md#partiescodeoptionsget) | **GET** /parties/{code}/options | Get party options |
| [**partiesCodeOptionsOptionIdDelete**](PartiesApi.md#partiescodeoptionsoptioniddelete) | **DELETE** /parties/{code}/options/{optionId} | Delete an option |
| [**partiesCodeOptionsPost**](PartiesApi.md#partiescodeoptionspost) | **POST** /parties/{code}/options | Add an option |
| [**partiesCodePatch**](PartiesApi.md#partiescodepatch) | **PATCH** /parties/{code} | Patch a party |
| [**partiesPost**](PartiesApi.md#partiespost) | **POST** /parties | Create a new party |



## partiesCodeAttendeesAttendeeIdDelete

> partiesCodeAttendeesAttendeeIdDelete(code, attendeeId)

Delete an attendee

Delete an attendee from a party by index

### Example

```ts
import {
  Configuration,
  PartiesApi,
} from '';
import type { PartiesCodeAttendeesAttendeeIdDeleteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PartiesApi();

  const body = {
    // string | Party Code
    code: code_example,
    // number | Attendee Index
    attendeeId: 56,
  } satisfies PartiesCodeAttendeesAttendeeIdDeleteRequest;

  try {
    const data = await api.partiesCodeAttendeesAttendeeIdDelete(body);
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
| **code** | `string` | Party Code | [Defaults to `undefined`] |
| **attendeeId** | `number` | Attendee Index | [Defaults to `undefined`] |

### Return type

`void` (Empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## partiesCodeAttendeesGet

> Array&lt;string&gt; partiesCodeAttendeesGet(code)

Get party attendees

Get attendees for a specific party

### Example

```ts
import {
  Configuration,
  PartiesApi,
} from '';
import type { PartiesCodeAttendeesGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PartiesApi();

  const body = {
    // string | Party Code
    code: code_example,
  } satisfies PartiesCodeAttendeesGetRequest;

  try {
    const data = await api.partiesCodeAttendeesGet(body);
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
| **code** | `string` | Party Code | [Defaults to `undefined`] |

### Return type

**Array<string>**

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


## partiesCodeAttendeesPost

> HandlerStringValue partiesCodeAttendeesPost(code, value)

Add an attendee

Add an attendee to a party

### Example

```ts
import {
  Configuration,
  PartiesApi,
} from '';
import type { PartiesCodeAttendeesPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PartiesApi();

  const body = {
    // string | Party Code
    code: code_example,
    // HandlerStringValue | Attendee Name
    value: ...,
  } satisfies PartiesCodeAttendeesPostRequest;

  try {
    const data = await api.partiesCodeAttendeesPost(body);
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
| **code** | `string` | Party Code | [Defaults to `undefined`] |
| **value** | [HandlerStringValue](HandlerStringValue.md) | Attendee Name | |

### Return type

[**HandlerStringValue**](HandlerStringValue.md)

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


## partiesCodeBeersPost

> partiesCodeBeersPost(code, beer)

Add a beer

Add a beer for an attendee in a party

### Example

```ts
import {
  Configuration,
  PartiesApi,
} from '';
import type { PartiesCodeBeersPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PartiesApi();

  const body = {
    // string | Party Code
    code: code_example,
    // HandlerBeerDTO | Beer Details
    beer: ...,
  } satisfies PartiesCodeBeersPostRequest;

  try {
    const data = await api.partiesCodeBeersPost(body);
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
| **code** | `string` | Party Code | [Defaults to `undefined`] |
| **beer** | [HandlerBeerDTO](HandlerBeerDTO.md) | Beer Details | |

### Return type

`void` (Empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: Not defined


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## partiesCodeGet

> ServicePartyDTO partiesCodeGet(code)

Get a party

Get a party by its code

### Example

```ts
import {
  Configuration,
  PartiesApi,
} from '';
import type { PartiesCodeGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PartiesApi();

  const body = {
    // string | Party Code
    code: code_example,
  } satisfies PartiesCodeGetRequest;

  try {
    const data = await api.partiesCodeGet(body);
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
| **code** | `string` | Party Code | [Defaults to `undefined`] |

### Return type

[**ServicePartyDTO**](ServicePartyDTO.md)

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


## partiesCodeOptionsGet

> Array&lt;string&gt; partiesCodeOptionsGet(code)

Get party options

Get options for a specific party

### Example

```ts
import {
  Configuration,
  PartiesApi,
} from '';
import type { PartiesCodeOptionsGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PartiesApi();

  const body = {
    // string | Party Code
    code: code_example,
  } satisfies PartiesCodeOptionsGetRequest;

  try {
    const data = await api.partiesCodeOptionsGet(body);
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
| **code** | `string` | Party Code | [Defaults to `undefined`] |

### Return type

**Array<string>**

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


## partiesCodeOptionsOptionIdDelete

> partiesCodeOptionsOptionIdDelete(code, optionId)

Delete an option

Delete an option from a party by index

### Example

```ts
import {
  Configuration,
  PartiesApi,
} from '';
import type { PartiesCodeOptionsOptionIdDeleteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PartiesApi();

  const body = {
    // string | Party Code
    code: code_example,
    // number | Option Index
    optionId: 56,
  } satisfies PartiesCodeOptionsOptionIdDeleteRequest;

  try {
    const data = await api.partiesCodeOptionsOptionIdDelete(body);
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
| **code** | `string` | Party Code | [Defaults to `undefined`] |
| **optionId** | `number` | Option Index | [Defaults to `undefined`] |

### Return type

`void` (Empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## partiesCodeOptionsPost

> HandlerStringValue partiesCodeOptionsPost(code, value)

Add an option

Add an option to a party

### Example

```ts
import {
  Configuration,
  PartiesApi,
} from '';
import type { PartiesCodeOptionsPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PartiesApi();

  const body = {
    // string | Party Code
    code: code_example,
    // HandlerStringValue | Option Value
    value: ...,
  } satisfies PartiesCodeOptionsPostRequest;

  try {
    const data = await api.partiesCodeOptionsPost(body);
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
| **code** | `string` | Party Code | [Defaults to `undefined`] |
| **value** | [HandlerStringValue](HandlerStringValue.md) | Option Value | |

### Return type

[**HandlerStringValue**](HandlerStringValue.md)

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


## partiesCodePatch

> ServicePartyDTO partiesCodePatch(code, patchReq)

Patch a party

Update a party status

### Example

```ts
import {
  Configuration,
  PartiesApi,
} from '';
import type { PartiesCodePatchRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PartiesApi();

  const body = {
    // string | Party Code
    code: code_example,
    // HandlerPatchPartyRequest | Patch Request
    patchReq: ...,
  } satisfies PartiesCodePatchRequest;

  try {
    const data = await api.partiesCodePatch(body);
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
| **code** | `string` | Party Code | [Defaults to `undefined`] |
| **patchReq** | [HandlerPatchPartyRequest](HandlerPatchPartyRequest.md) | Patch Request | |

### Return type

[**ServicePartyDTO**](ServicePartyDTO.md)

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


## partiesPost

> ServicePartyDTO partiesPost(party)

Create a new party

Creates a new party with a generated 6-character code

### Example

```ts
import {
  Configuration,
  PartiesApi,
} from '';
import type { PartiesPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PartiesApi();

  const body = {
    // ModelsParty | Party details
    party: ...,
  } satisfies PartiesPostRequest;

  try {
    const data = await api.partiesPost(body);
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
| **party** | [ModelsParty](ModelsParty.md) | Party details | |

### Return type

[**ServicePartyDTO**](ServicePartyDTO.md)

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

