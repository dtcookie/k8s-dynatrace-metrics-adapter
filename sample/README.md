# Sample deployment

This sample deployment showcases how to use the Dynatrace Metrics Adapter in combination with a Horizontal Pod Scaler.

Prerequisite for using the sample is that you already deployed OneAgent Operator or Dynatrace Operator as explained [here](https://www.dynatrace.com/support/help/technology-support/cloud-platforms/kubernetes/) on your Kubernetes cluster.

# Deploy the application

The sample consists of a backend and a frontend HTTP Server. The frontend server comes with its own load generator, therefore there is no need to expose it as a service.

Executing `kubectl apply -f deploy.yaml` will deploy both, frontend and backend within the namespace `dt-metrics-sample`. By default only one instance of should will be created.

# Validate Dynatrace Monitoring

If all the prerequisites for this sample are met, Two new Services will show up within `Transactions and Services` in the Dynatrace WebUI.

* `frontend-http-server-*`
* `backend-http-server-*`

With just one instance of `frontend-http-server` running, the average response time of this service will be around 5 seconds. Bu the more instances of the same application are getting launched the smaller the response time will be.

# Applying Tags

Later on you need to be able to narrow down the entities Dynatrace should provide the metrics for.
In this example, the easiest way to do that is to apply a Tag to the monitored service. We're choosing the tag `hpa` and give it the value `frontend-http-server`.
![Applying Tags](https://github.com/dtcookie/k8s-dynatrace-metrics-adapter/blob/main/sample/img/tags.png?raw=true "Applying Tags")

# Apply the Horizontal Pod Scaler

After executing `kubectl apply -f hpa.yaml` Kubernetes will decide to launch additional 3 pods for `fronted=http-server`.
Reason for that are the requirement within `hpa.yaml`.

```apiVersion: autoscaling/v2beta2
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

* Only the metric `com.dynatrace.builtin:service.responsetime` contributes to the decision whether to scale the number of replicas.
  - The Endpoint `Timeseries` of the Dynatrace REST API (/api/v1/timeseries) offers a complete list of available metrics - not just built in metrics but even custom metrics configured by any user
* The `selector` states that only metrics for entities with the tag `hpa` having the value `frontend-http-server` are of interest.
* The `target` states that scaling up is necessary until a response time of 1500 milliseconds is reached
