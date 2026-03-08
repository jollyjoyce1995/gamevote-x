
# ModelsParty


## Properties

Name | Type
------------ | -------------
`attendees` | Array&lt;string&gt;
`code` | string
`id` | string
`options` | Array&lt;string&gt;
`pollId` | string
`results` | { [key: string]: number; }
`status` | [ModelsPartyStatus](ModelsPartyStatus.md)

## Example

```typescript
import type { ModelsParty } from ''

// TODO: Update the object below with actual values
const example = {
  "attendees": null,
  "code": null,
  "id": null,
  "options": null,
  "pollId": null,
  "results": null,
  "status": null,
} satisfies ModelsParty

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ModelsParty
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


