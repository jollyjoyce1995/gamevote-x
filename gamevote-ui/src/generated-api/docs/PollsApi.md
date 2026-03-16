# PollsApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**partiesCodeVotesAttendeePost**](PollsApi.md#partiescodevotesattendeepost) | **POST** /parties/{code}/votes/{attendee} | Submit a vote |



## partiesCodeVotesAttendeePost

> partiesCodeVotesAttendeePost(code, attendee, choices)

Submit a vote

Submit a vote for an attendee

### Example

```ts
import {
  Configuration,
  PollsApi,
} from '';
import type { PartiesCodeVotesAttendeePostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PollsApi();

  const body = {
    // string | Poll ID
    code: code_example,
    // string | Attendee Name
    attendee: attendee_example,
    // { [key: string]: number; } | Choices (-1, 0, or 1)
    choices: ...,
  } satisfies PartiesCodeVotesAttendeePostRequest;

  try {
    const data = await api.partiesCodeVotesAttendeePost(body);
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
| **code** | `string` | Poll ID | [Defaults to `undefined`] |
| **attendee** | `string` | Attendee Name | [Defaults to `undefined`] |
| **choices** | `{ [key: string]: number; }` | Choices (-1, 0, or 1) | |

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
| **204** | No Content |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

