## My try to enable Scale-To-Zero

#### Edit the values.yaml

```yaml
...
dlx:
  enabled: true
...
autoscaler:
  enabled: true
...

platform: 
  scaleToZero:
      httpTriggerIngressAnnotations: null
      inactivityWindowPresets:
      - 1m
      - 2m
      - 5m
      - 10m
      - 30m
      mode: enabled
      resourceReadinessTimeout: 5m
      scaleResources:
      - metricName: nuclio_processor_handled_events_total
        threshold: 0
        windowSize: 10m
      scalerInterval: 1m 
...
```

### Starting Minikube

```bash
minikube start --memory 8192 --cpus 6 
minikube addons enable metrics-server
```

### Setting Up the Function Registry

```bash
eval $(minikube docker-env)
docker run -d -p 5000:5000 --name function-registry registry:latest
```

### Building Nuclio

```bash
make build
```

### Installing Prometheus

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prometheus prometheus-community/prometheus
```

### Installing and Configuring Prometheus Adapter

```bash
helm install prometheus-adapter prometheus-community/prometheus-adapter
helm upgrade prometheus-adapter prometheus-community/prometheus-adapter -f ./prometheus-adapter-config.yaml
```

#### Prometheus Adapter Configuration

```yaml
prometheusAdapter:
  resourceMappings:
    rules:
      - metricsQuery: sum(rate(<<.Series>>{<<.LabelMatchers>>}[10m])) by (<<.GroupBy>>)
        resources:
          overrides:
            kubernetes_pod_name: pod
            function: nucliofunction
            namespace: default
        name:
          matches: "^(.*)_total"
          as: "${1}_per_10m"
        seriesQuery: nuclio_processor_handled_events_total{namespace!="",exported_instance!="",function!="",trigger_kind="http"}
```


The metric `nuclio_processor_handled_events_total_per_10m` should be available at Prometheus ... it isnt ....
