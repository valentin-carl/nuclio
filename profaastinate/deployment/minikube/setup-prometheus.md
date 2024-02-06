### Starting Minikube

```bash
minikube start --memory 8192 --cpus 6 
```

### Enabling Metrics Server

```bash
minikube addons enable metrics-server
```

###  start function registry

```bash
# Set docker env to minikube, so that we can push images to the minikube registry
eval $(minikube docker-env)

# Run a Docker container to expose Minikube's registry
docker run -d -p 5000:5000 --name function-registry registry:latest
```

### Build Nuclio
    
```bash
make build
```

### Installing Prometheus

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prometheus prometheus-community/prometheus --namespace monitoring --create-namespace
```

### Installing Prometheus Adapter

```bash
helm install prometheus-adapter prometheus-community/prometheus-adapter
helm upgrade prometheus-adapter prometheus-community/prometheus-adapter --namespace monitoring --values ./prometheus-adapter-config.yaml
```