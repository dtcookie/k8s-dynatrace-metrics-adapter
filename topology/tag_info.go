package topology

type TagInfo struct {
	Context Context `json:"context"`
	Key     string  `json:"key"`
	Value   string  `json:"value"`
}

type Context string

var Contexts = struct {
	AWS          Context
	AWSGeneric   Context
	Azure        Context
	CloudFoundry Context
	Contextless  Context
	Environment  Context
	GoogleCloud  Context
	Kubernetes   Context
}{
	AWS:          Context("AWS"),
	AWSGeneric:   Context("AWS_GENERIC"),
	Azure:        Context("AZURE"),
	CloudFoundry: Context("CLOUD_FOUNDRY"),
	Contextless:  Context("CONTEXTLESS"),
	Environment:  Context("ENVIRONMENT"),
	GoogleCloud:  Context("GOOGLE_CLOUD"),
	Kubernetes:   Context("KUBERNETES"),
}
