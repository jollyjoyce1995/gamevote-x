# PartiesApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**partiesCodeAttendeesAttendeeIdDelete**](PartiesApi.md#partiescodeattendeesattendeeiddelete) | **DELETE** /parties/{code}/attendees/{attendeeId} | Delete an attendee |
| [**partiesCodeAttendeesGet**](PartiesApi.md#partiescodeattendeesget) | **GET** /parties/{code}/attendees | Get party attendees |
| [**partiesCodeAttendeesPost**](PartiesApi.md#partiescodeattendeespost) | **POST** /parties/{code}/attendees | Add an attendee |
| [**partiesCodeBeersPost**](PartiesApi.md#partiescodebeerspost) | **POST** /parties/{code}/beers | Add a beer |
| [**partiesCodeGet**](PartiesApi.md#partiescodeget) | **GET** /parties/{code} | Get a party |
| [**partiesCodeOptionsGameNameDelete**](PartiesApi.md#partiescodeoptionsgamenamedelete) | **DELETE** /parties/{code}/options/{gameName} | Delete an option |
| [**partiesCodeOptionsGet**](PartiesApi.md#partiescodeoptionsget) | **GET** /parties/{code}/options | Get party options |
| [**partiesCodeOptionsPost**](PartiesApi.md#partiescodeoptionspost) | **POST** /parties/{code}/options | Add an option |
| [**partiesCodePatch**](PartiesApi.md#partiescodepatch) | **PATCH** /parties/{code} | Patch a party |
| [**partiesCodeStreamGet**](PartiesApi.md#partiescodestreamget) | **GET** /parties/{code}/stream | SSE stream for a party |
| [**partiesGet**](PartiesApi.md#partiesget) | **GET** /parties | Get all parties |
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


## partiesCodeOptionsGameNameDelete

> partiesCodeOptionsGameNameDelete(code, gameName)

Delete an option

Delete an option from a party by its name

### Example

```ts
import {
  Configuration,
  PartiesApi,
} from '';
import type { PartiesCodeOptionsGameNameDeleteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PartiesApi();

  const body = {
    // string | Party Code
    code: code_example,
    // string | Game Name
    gameName: gameName_example,
  } satisfies PartiesCodeOptionsGameNameDeleteRequest;

  try {
    const data = await api.partiesCodeOptionsGameNameDelete(body);
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
| **gameName** | `string` | Game Name | [Defaults to `undefined`] |

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


## partiesCodeOptionsPost

> ModelsPartyOption partiesCodeOptionsPost(code, option)

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
    // ModelsPartyOption | Option Details
    option: ...,
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
| **option** | [ModelsPartyOption](ModelsPartyOption.md) | Option Details | |

### Return type

[**ModelsPartyOption**](ModelsPartyOption.md)

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


## partiesCodeStreamGet

> string partiesCodeStreamGet(code, username)

SSE stream for a party

Opens a Server-Sent Events stream for real-time party updates

### Example

```ts
import {
  Configuration,
  PartiesApi,
} from '';
import type { PartiesCodeStreamGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PartiesApi();

  const body = {
    // string | Party Code
    code: code_example,
    // string | Username of the connected client
    username: username_example,
  } satisfies PartiesCodeStreamGetRequest;

  try {
    const data = await api.partiesCodeStreamGet(body);
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
| **username** | `string` | Username of the connected client | [Defaults to `undefined`] |

### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `text/event-stream`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## partiesGet

> Array&lt;ServicePartyDTO&gt; partiesGet()

Get all parties

Get all parties ordered by ID

### Example

```ts
import {
  Configuration,
  PartiesApi,
} from '';
import type { PartiesGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PartiesApi();

  try {
    const data = await api.partiesGet();
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

[**Array&lt;ServicePartyDTO&gt;**](ServicePartyDTO.md)

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
    // ServicePartyDTO | Party details
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
| **party** | [ServicePartyDTO](ServicePartyDTO.md) | Party details | |

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

