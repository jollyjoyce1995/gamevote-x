
# ModelsDrinkType


## Properties

Name | Type
------------ | -------------
`alcoholPercent` | number
`beerEquivalent` | number
`id` | string
`name` | string
`unitName` | string
`volumeMl` | number

## Example

```typescript
import type { ModelsDrinkType } from ''

// TODO: Update the object below with actual values
const example = {
  "alcoholPercent": null,
  "beerEquivalent": null,
  "id": null,
  "name": null,
  "unitName": null,
  "volumeMl": null,
} satisfies ModelsDrinkType

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ModelsDrinkType
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


