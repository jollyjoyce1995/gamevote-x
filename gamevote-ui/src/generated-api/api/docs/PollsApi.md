# PollsApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**pollsGet**](PollsApi.md#pollsget) | **GET** /polls | Get all polls |
| [**pollsIdGet**](PollsApi.md#pollsidget) | **GET** /polls/{id} | Get a poll |
| [**pollsIdOutstandingGet**](PollsApi.md#pollsidoutstandingget) | **GET** /polls/{id}/outstanding | Get outstanding voters |
| [**pollsIdPut**](PollsApi.md#pollsidput) | **PUT** /polls/{id} | Update a poll |
| [**pollsIdResultsGet**](PollsApi.md#pollsidresultsget) | **GET** /polls/{id}/results | Get poll results |
| [**pollsIdVotesAttendeePut**](PollsApi.md#pollsidvotesattendeeput) | **PUT** /polls/{id}/votes/{attendee} | Submit a vote |
| [**pollsIdVotesGet**](PollsApi.md#pollsidvotesget) | **GET** /polls/{id}/votes | Get all votes |
| [**pollsPost**](PollsApi.md#pollspost) | **POST** /polls | Create a poll |



## pollsGet

> Array&lt;ModelsPoll&gt; pollsGet()

Get all polls

Retrieve a list of all polls

### Example

```ts
import {
  Configuration,
  PollsApi,
} from '';
import type { PollsGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PollsApi();

  try {
    const data = await api.pollsGet();
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

[**Array&lt;ModelsPoll&gt;**](ModelsPoll.md)

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


## pollsIdGet

> ModelsPoll pollsIdGet(id)

Get a poll

Get a poll by its ID

### Example

```ts
import {
  Configuration,
  PollsApi,
} from '';
import type { PollsIdGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PollsApi();

  const body = {
    // string | Poll ID
    id: id_example,
  } satisfies PollsIdGetRequest;

  try {
    const data = await api.pollsIdGet(body);
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
| **id** | `string` | Poll ID | [Defaults to `undefined`] |

### Return type

[**ModelsPoll**](ModelsPoll.md)

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


## pollsIdOutstandingGet

> Array&lt;string&gt; pollsIdOutstandingGet(id)

Get outstanding voters

Get attendees who have not yet voted

### Example

```ts
import {
  Configuration,
  PollsApi,
} from '';
import type { PollsIdOutstandingGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PollsApi();

  const body = {
    // string | Poll ID
    id: id_example,
  } satisfies PollsIdOutstandingGetRequest;

  try {
    const data = await api.pollsIdOutstandingGet(body);
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
| **id** | `string` | Poll ID | [Defaults to `undefined`] |

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


## pollsIdPut

> ModelsPoll pollsIdPut(id, poll)

Update a poll

Update a poll details (used to resume or complete)

### Example

```ts
import {
  Configuration,
  PollsApi,
} from '';
import type { PollsIdPutRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PollsApi();

  const body = {
    // string | Poll ID
    id: id_example,
    // ModelsPoll | Poll Details
    poll: ...,
  } satisfies PollsIdPutRequest;

  try {
    const data = await api.pollsIdPut(body);
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
| **id** | `string` | Poll ID | [Defaults to `undefined`] |
| **poll** | [ModelsPoll](ModelsPoll.md) | Poll Details | |

### Return type

[**ModelsPoll**](ModelsPoll.md)

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


## pollsIdResultsGet

> { [key: string]: number; } pollsIdResultsGet(id)

Get poll results

Get aggregated poll results

### Example

```ts
import {
  Configuration,
  PollsApi,
} from '';
import type { PollsIdResultsGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PollsApi();

  const body = {
    // string | Poll ID
    id: id_example,
  } satisfies PollsIdResultsGetRequest;

  try {
    const data = await api.pollsIdResultsGet(body);
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
| **id** | `string` | Poll ID | [Defaults to `undefined`] |

### Return type

**{ [key: string]: number; }**

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


## pollsIdVotesAttendeePut

> { [key: string]: number; } pollsIdVotesAttendeePut(id, attendee, choices)

Submit a vote

Submit a vote for an attendee

### Example

```ts
import {
  Configuration,
  PollsApi,
} from '';
import type { PollsIdVotesAttendeePutRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PollsApi();

  const body = {
    // string | Poll ID
    id: id_example,
    // string | Attendee Name
    attendee: attendee_example,
    // { [key: string]: number; } | Choices (-1, 0, or 1)
    choices: ...,
  } satisfies PollsIdVotesAttendeePutRequest;

  try {
    const data = await api.pollsIdVotesAttendeePut(body);
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
| **id** | `string` | Poll ID | [Defaults to `undefined`] |
| **attendee** | `string` | Attendee Name | [Defaults to `undefined`] |
| **choices** | `{ [key: string]: number; }` | Choices (-1, 0, or 1) | |

### Return type

**{ [key: string]: number; }**

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


## pollsIdVotesGet

> { [key: string]: { [key: string]: number; }; } pollsIdVotesGet(id)

Get all votes

Get votes mapping the attendee to their choices

### Example

```ts
import {
  Configuration,
  PollsApi,
} from '';
import type { PollsIdVotesGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PollsApi();

  const body = {
    // string | Poll ID
    id: id_example,
  } satisfies PollsIdVotesGetRequest;

  try {
    const data = await api.pollsIdVotesGet(body);
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
| **id** | `string` | Poll ID | [Defaults to `undefined`] |

### Return type

**{ [key: string]: { [key: string]: number; }; }**

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


## pollsPost

> ModelsPoll pollsPost(poll)

Create a poll

Create a new poll manually

### Example

```ts
import {
  Configuration,
  PollsApi,
} from '';
import type { PollsPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PollsApi();

  const body = {
    // ModelsPoll | Poll Details
    poll: ...,
  } satisfies PollsPostRequest;

  try {
    const data = await api.pollsPost(body);
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
| **poll** | [ModelsPoll](ModelsPoll.md) | Poll Details | |

### Return type

[**ModelsPoll**](ModelsPoll.md)

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

