
# ServicePartyDTO


## Properties

Name | Type
------------ | -------------
`links` | [{ [key: string]: ServiceLink; }](ServiceLink.md)
`attendees` | Array&lt;string&gt;
`beerCount` | number
`beerPerAttendee` | { [key: string]: number; }
`code` | string
`id` | string
`options` | [Array&lt;ModelsPartyOption&gt;](ModelsPartyOption.md)
`results` | { [key: string]: number; }
`status` | string

## Example

```typescript
import type { ServicePartyDTO } from ''

// TODO: Update the object below with actual values
const example = {
  "links": null,
  "attendees": null,
  "beerCount": null,
  "beerPerAttendee": null,
  "code": null,
  "id": null,
  "options": null,
  "results": null,
  "status": null,
} satisfies ServicePartyDTO

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ServicePartyDTO
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


