FROM BASEIMAGE
COPY k8s-dynatrace-metrics-adapter /
ENTRYPOINT ["/k8s-dynatrace-metrics-adapter", "--logtostderr=true"]
