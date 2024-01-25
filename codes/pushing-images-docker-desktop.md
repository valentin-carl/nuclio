# Pushing Docker Images to a Local Insecure Registry Inside Minikube 

Apparently, Docker Desktop has a problem with pushing images to `localhost:5000/...` ([see more](https://github.com/docker/for-mac/issues/3611)). 

Solution: don't push directly to `localhost`. Instead, do something like this:

```shell
sudo echo "127.0.0.1 mynameisjeff" >> /etc/hosts
docker image tag local-message-bus:latest mynameisjeff:5000/local-message-bus:latest
docker image push mynameisjeff:5000/local-message-bus:latest
curl mynameisjeff:5000/v2/_catalog 
```
