
# ModelsPoll


## Properties

Name | Type
------------ | -------------
`attendees` | Array&lt;string&gt;
`id` | string
`options` | [Array&lt;ModelsPartyOption&gt;](ModelsPartyOption.md)
`status` | [ModelsPollStatus](ModelsPollStatus.md)

## Example

```typescript
import type { ModelsPoll } from ''

// TODO: Update the object below with actual values
const example = {
  "attendees": null,
  "id": null,
  "options": null,
  "status": null,
} satisfies ModelsPoll

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ModelsPoll
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


