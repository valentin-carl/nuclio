# Benchmarking SAND in Nuclio

## Starting the workflow

### Invoke function `entrytask` from within pod

```shell
TARGET=10000
kubectl exec -it <sand-pod-name> -- curl localhost:8080 -d "$TARGET"
```
