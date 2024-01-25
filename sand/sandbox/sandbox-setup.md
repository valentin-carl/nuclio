# Sandbox Setup

## Background

In [SAND](https://www.usenix.org/system/files/conference/atc18/atc18-akkus.pdf), serverless workflows are optimized by removing isolation between functions of the same application.
To do so, functions are grouped into sandboxes.
To this end, Docker containers are used in the paper; we use Kubernetes pods as it fits our level of abstraction and Nuclio's architecture better.
Additionally, this sandbox contains a message bus, which allows functions invoking the next part of a workflow to take a shortcut.
Instead of sending a message to the Nuclio dashboard, which is in a different pod, function invocations of the same workflow can go directly from container to container.
This reduces the total workflow duration without having to optimize the duration of individual function calls.

## Target setup

Nuclio uses custom resource definitions for creating projects, api gateways, functions, and function events.
Kubernetes and Nuclio's controller create a pod for each function, containing one container each. 
This allows functions to be handled, scaled, restarted, etc. separately but is not optimal for optimizing the duration of workflows.
Our goal is to have one pod for an application, consisting of multiple containers for functions, a message bus, and a grain worker (see paper) for each function.

For now, we do this manually as the Nuclio controller doesn't know (yet) how to handle instances of the custom Kubernetes resource `nucliofunction` with multiple containers in their spec.
The relevant code for this is in the file `pkg/platform/kube/functionres/lazy.go`.

## Manually creating application sandboxes

### 1. Creating images for your Nuclio functions

We still want to use the Docker images Nuclio creates for the functions, we just want to put them in a different pod.
Hence, the first step is to create Nuclio function as usual and deploy them. 
This can be done using the UI Nuclio provides or with `nuctl`.
Note: We don't do this because we actually want the pods, deployments, etc. that are created by this; we only want the Docker images that are pushed into the insecure registry running inside the cluster in the process. 
Everything else will be deleted later on.

### 2. Pushing images of the message bus and grain worker to the registry

To create a deployment with our custom containers running inside the same pod, we also need these images to be available inside the registry running in the cluster.
First of all, create the necessary Docker images as usual.
Next, add the insecure registry to your Docker (if you haven't already).
On Mac, you can use Docker Desktop to add the registry's address to the Docker daemon's config. 
On Linux, you can (probably) find this file as `/etc/docker/daemon.json`.

Here's an example configuration.

```json
{
  "builder": {
    "gc": {
      "defaultKeepStorage": "20GB",
      "enabled": true
    }
  },
  "experimental": false,
  "insecure-registries": [
    "192.168.49.2:5000"
  ]
}
```

Make sure to replace the address with the correct address.
You can find out correct ip by running

```shell
minikube ip
```

and add the correct port.

Next, tag the images you want to add to the registry.
(Make sure to replace the port if your registry is listening on a different one.)

TODO ADD COMMANDS FOR GRAIN WORKER

```shell
docker image tag local-message-bus:latest $(minikube ip):5000/local-message-bus:latest
```

Now you can push the image to the insecure registry.

```shell
docker push $(minikube ip):5000/local-message-bus:latest
```

Repeat these steps for all images you want to use in the sandbox.

You can get the list of images inside the registry as follows.

```shell
minikube ssh -- curl localhost:5000/v2/_catalog
```

### 3. Creating a new deployment

Next, we want to use the images we pushed into our cluster to create a pod with containers using these images inside it.
Note: Kubernetes doesn't like it when we try to add containers to already running pods, and we can't create a new `nucliofunction` resource because the controller doesn't know how to handle these with multiple containers in their spec.
Hence, we create a new deployment.
To do so, we choose the function that starts the workflow and get its configuration yaml-file.
It is important that we choose the function that starts the workflow because the Nuclio dashboard can only associate one function with a deployment and, in turn, only invoke one of the functions that will be running inside our pod.
(But we can still invoke these functions in other ways, for example, through the message bus.)

#### Creating a file to specify the new deployment

Assuming our function is called “start”, and its deployment is called “nuclio-start”, we can get a copy of the config file we need by running this command.

```shell
kubectl get deployment nuclio-start -o yaml > sandbox.yaml
```

Next, remove the fields 

- metadata > creationTimestamp
- metadata > resourceVersion
- metadata > uid

from the file. 
These will be created automatically by kubernetes once we create the new resource, and we would get an error if we included them in a file specifying a new resource.

#### Adding new containers

Finally, we can add new containers to the deployment under `spec > template > spec > containers`.
The easiest way to do this is to copy the spec of the already existing container in that deployment (which is for the function starting the workflow) and to modify the necessary fields.
For each new container, we need to adjust the name, ports, and image.
The ports used by the containers cannot overlap, since only the pod gets an ip-address, and the containers running inside don't. 
For the image, we need to specify that it is located in the insecure registry running inside the cluster by calling it, e.g., `localhost:5000/local-message-bus:latest`.
(Here, localhost refers to the minikube VM.)

#### Creating the deployment

We can create the deployment with this command.

```shell
kubectl create -f sand.yaml
```

Now, check the Minikube dashboard to see if the deployment and all of our containers are running. 
Alternatively, you can also use the command line.

```shell
kubectl get deployments
```

### 4. Setting up the broker and grain workers

#### Connecting to the pods and containers

To connect to the pod, find out its name (it will have some random number as part of it) using 

```shell
kubectl get pods
```

and get a shell this way.

```shell
kubectl exec -it <podname> -- bash
```

To connect to the individual containers, use Docker instead. 
Important: Make sure to use the docker running inside the minikube VM.

```shell
# use docker inside minikube
eval $(minikube docker-env)

# connect to the container
docker exec -it <containername> bash
```

#### Networking within the sandbox

Now that we can connect to the containers directly, we can set up the local message bus and the grain workers. 
Because all of our containers are running within the same pod, networking between them becomes surprisingly easy: The other containers can all be reached via `localhost`.
Hence, we just have to adjust the configuration files for the grain workers to include the correct function names and ports.
(To ensure the grain workers use the correct configuration, restart them after reconfiguring.) 

> TODO change the grain workers to be configured via environment variables instead, those can be set in the deployment config instead and this step becomes easier.

#### Using the broker's management interface

> TODO

### Removing unnecessary deployments created in the process

At the end, use the Minikube dashboard or `kubectl` to remove the resources that were deployed while creating the functions' Docker images in the beginning.