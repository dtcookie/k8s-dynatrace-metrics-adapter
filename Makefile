ndef = $(if $(value $(1)),,$(error $(1) not set))
IMAGE?=k8s-dynatrace-metrics-adapter
EXECUTABLE=$(IMAGE)
TEMP_DIR:=$(shell mktemp -d)
ARCH?=amd64
OUT_DIR?=./_output
VERSION?=latest

.PHONY: all build test container push clean

all: build

clean:
	rm -rf $(TEMP_DIR)
	rm -rf $(OUT_DIR)

build: tidy	
	CGO_ENABLED=0 GOOS=linux GOARCH=$(ARCH) go build -o $(OUT_DIR)/$(ARCH)/$(EXECUTABLE)

tidy:
	go mod tidy

test:
	CGO_ENABLED=0 go test ./...

container: build
	$(call ndef,REGISTRY)
	cp Dockerfile $(TEMP_DIR)
	cp $(OUT_DIR)/$(ARCH)/$(EXECUTABLE) $(TEMP_DIR)/$(EXECUTABLE)
	cd $(TEMP_DIR) && sed -i.bak "s|BASEIMAGE|scratch|g" Dockerfile
	sed -i.bak 's|REGISTRY|'${REGISTRY}'|g' deploy.yaml
	docker build -t $(REGISTRY)/$(IMAGE)-$(ARCH):$(VERSION) $(TEMP_DIR)
	rm -rf $(TEMP_DIR) deploy.yaml.bak

push: container
	docker push $(REGISTRY)/$(IMAGE)-$(ARCH):$(VERSION)