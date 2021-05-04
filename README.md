# Dynatrace Kubernetes Adapter for External Metrics

This project implements a Metric API Server offering additional External Metrics for Kubernetes deployments monitored with Dynatrace.

## Installation

A prebuilt version of this project already is available on [Docker Hub](https://hub.docker.com/repository/docker/dtcookie/k8s-dynatrace-metrics-adapter-amd64).

* Edit `deploy.yaml` and navigate to line 62. Enter the URL of your Dynatrace Environment here.
* Line 59 in `deploy.yaml` refers to a secret called `dynametric`. It holds the `API TOKEN` that grants access to the metrics offered by the Metric API Server.
  The minimum scope of the API Token should be `DataExport (Access problem and event feed, metrics, and topology)`
  Unless you want to choose a different name, this is how to create that secret.
    ```
    kubectl -n dynatrace-metrics create secret generic dynametric --from-literal="apiToken=####"
    ```    
* Finally execute `kubectl apply -f deploy.yaml`. This will automatically create a namespace `dynatrace-metrics` with the required pods.

## Configuration

This is an example for a Horizontal Pod Scaler scaling `frontend-http-server` based on the metric `com.dynatrace.builtin:service.responsetime`, expecting the average response time of the monitored service to not increase 1.5 seconds.

```
apiVersion: autoscaling/v2beta2
kind: HorizontalPodAutoscaler
metadata:
  name: frontend-http-server
  namespace: dt-metrics-sample
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: frontend-http-server    
  minReplicas: 1
  maxReplicas: 20
  metrics:
  - type: External
    external:
      metric:
        name: com.dynatrace.builtin:service.responsetime
        selector:
          matchLabels:
            hpa: frontend-http-server
      target:
        type: AverageValue
        value: 1500000000m
```

Any metrics available via the Dynatrace REST API (`/api/v1/timeseries`) are eligible.
Additional documentatation about this example can be found [here](https://github.com/dtcookie/k8s-dynatrace-metrics-adapter/tree/main/sample).
