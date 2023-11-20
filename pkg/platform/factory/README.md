The factory component does the job of creating new platforms based on the requested type and configuration it receives. It has two main functions `CreatePlatform` and `GetPlatformByType`.

`CreatePlatform` function receives context, logger, platformType, platformConfiguration, defaultNamespace as inputs. Then it creates a new platform of type either local or Kubernetes based on the received platformType. If the platformType is neither local nor Kubernetes, it throws an error.

The `GetPlatformByType` function returns the platform type based on the given platformType and platformConfiguration. If the platformType is auto, the function first checks if there is a configured kubeconfig path or if the platform is running in a Kubernetes cluster, if either is true, it returns the Kubernetes platform type, otherwise, it returns the local platform type. If the platformType is neither local, Kubernetes, nor auto, it throws an error.
