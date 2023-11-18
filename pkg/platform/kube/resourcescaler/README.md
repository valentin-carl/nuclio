The component from Nuclio's Kubernetes platform located at `nuclio/pkg/platform/kube/resourcescaler` is `resourcescaler`. The job of this component is to handle the scaling of functions to and from zero automatically based on configurable metrics and settings. It interacts with the Kubernetes API to adjust the number of replicas for a function based on its current demand.

The component `resourcescaler` contains two .go files:

1. `resourcescaler.go`
This file contains the `NuclioResourceScaler` struct and the logic for scaling functions. The `NuclioResourceScaler` implements methods to:
   - Change the scale of resources by setting them to zero or scaling them back up.
   - Parse scaling configurations and resources.
   - Update the status of Nuclio functions in relation to scaling actions.
   - Manage HTTP requests for health checks to ensure a function is ready after scaling from zero.
   - Retrieve the list of resources that can be scaled and the resource scaling configuration.

2. `resourcescaler_test.go`
This file contains integration tests for the `resourcescaler` component, ensuring that the scaling functionality behaves as expected in different scenarios. The tests include:
   - Sanity check to ensure the scaler can scale a function down to zero and back to ready state.
   - Multi-target test to confirm that the scaler correctly handles multiple functions that share the same API gateway and that it respects scaling configurations.

The `NuclioResourceScaler` struct is also responsible for integrating with the V3IO scaler package to extend the functionality of scaling Nuclio functions asynchronously based on metrics.
