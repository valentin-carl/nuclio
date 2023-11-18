The `apigatewayres` component is responsible for managing resources associated with an API gateway. This includes creating, updating, waiting for, listing, and deleting resources. An API gateway in this context acts as an entry point for an application, allowing clients to interact with a service such as invoking a function.

The following Go files are part of this component:

1. `lazy.go`: This file contains the `lazyClient` struct and methods, which manage API gateway resources in Kubernetes. These resources include Kubernetes services, ingresses, and secrets needed for the API gateway functionality.

2. `lazy_test.go`: This file contains the tests for the functionalities provided in the `lazy.go` file. Here, the lazyClient methods are tested to validate their functionality.

3. `types.go`: This file defines the `Client` and `Resources` interfaces that are used to manage the API gateway resources. Any struct implementing the `Client` interface will need to provide the methods for creating, updating, waiting for, listing, and deleting resources. The `Resources` interface defines a method to retrieve a map of ingress resources.
