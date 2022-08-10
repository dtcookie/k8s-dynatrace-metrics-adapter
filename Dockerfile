FROM debian
COPY ./build/k8s-dynatrace-metrics-adapter /k8s-dynatrace-metrics-adapter
ENTRYPOINT ["/k8s-dynatrace-metrics-adapter", "--logtostderr=true"]
