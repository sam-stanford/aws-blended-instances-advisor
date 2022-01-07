# API

### Requests

Requests should be formatted in the following way.

```TypeScript
{
  "services": {
      "name": string; // The service's name. Must be unique
      "minMemory": number; // The minimum amount on memory required
      "maxVcpu": number; // The maximum of CPU cores which can be utilised
      "minInstances": number; // The minimum number of instances of the service which should be running 
      "maxInstances": number; // The maximumum number of isntances of the service which should be running
    }[];
  "advisor": {
    "type": string; // Only accepts "weighted" for now
    "weights": {
      "availability": number; 
      "performance": number;
      "price": number;
    };
  };
  "options": {
    "avoidRepeatedInstanceTypes": boolean;
    "shareInstancesBetweenServices": boolean;
    "considerFreeInstances": boolean;
    "regions": string[];
  };
}
```

---

}

## Responses

Responses will be formatted in the following way.

```TypeScript

{
  [region: string]: {
    "score": number; 
    "instances": [id: string]: {
      "id": string; // UUID for the instance
      "name": string; // The AWS name/type of the instance
      "memory": number; // Memory in GB
      "vcpu": number; // Number of CPU cores
      "region": string; // AWS region
      "az": string; // AWS availability zone
      "os": string; // Operating system
      "price" number; // Price per hour in USD
      "revocProb": number; // The probability of revocation in the next month
    };
    "assignments": {
      "servicesToInstances": {[serviceName: string]: string};
      "instancesToServices": {[instanceId: string]: string};
    }
  };
}
```