# Local Message Bus

For reference: [Akkus et al. (2018) SAND: Towards High-Performance Serverless Computing](https://www.usenix.org/system/files/conference/atc18/atc18-akkus.pdf).

To optimize serverless workflows, the local message bus serves as a shortcut between function invocations of the same application. 
To do so, we deploy a RabbitMQ instance in the same pod as the functions of that application.
For each function of that application, a new queue is created at the broker.
With each function, a function worker (so-called “grain worker”) is deployed.
It listens for all events of the respective queue and invokes the function accordingly.

## Deploying the local message bus

This is still very much work in progress. For now, this is how to build locally for development.
(Note: Before executing the following commands, remember to change into the correct directories.)

### Message broker

```shell
docker build -t local-message-bus .
docker run -d -it -p 5672:5672 -p 15672:15672 --rm --name lmb local-message-bus
docker exec -it lmb bash
```

### Grain worker

At the moment, the config file has to be adjusted for the respective function before building the docker image.
It's not the best way, and we'll change it in the future; but it works for now.

```shell
docker build -t grainworker .
docker run -d -it --name gw grainworker
docker exec -it gw bash
```