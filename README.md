## API

## Contents

### Requests

Requests should be formatted in the following way.

```TypeScript
{
	"name": string,
	"minMemory": number,
	"maxVcpu": number,
	"minInstances": number,
	"totalInstances": number
}
```

- `name` (string) : The name of the service.
  - Names of services within the same request must be unique.
- `minMemory` (float) : The minimum amount of memory required by the service.
  - The advising service assumes that more memory than the given value will provide significantly diminishing returns to the performance of the specified service.
- `maxVcpu` (int) : The maximum number of CPU cores that are useful to the service.
  - The advising service assumes that more CPU cores than the given value will provide significantly dimishing returns to the performance of the specified service.
- `minInstances` (int) : The minimum number of instances of the service that are required to be available at all times.
  - The advising service will use non-transient resources for this number of instances.
- `totalInstances` (int) : The total number of instances of the service that are desired to be available.
  - The advising service will use transient or non-transient resources for the difference between `minInstances` and `totalInstances`.

---

## Responses

TODO

TODO: Update example config
