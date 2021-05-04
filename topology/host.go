package topology

import (
	"strings"
)

type Host struct {
	EntityBase
	Tags              []*TagInfo        `json:"tags"`
	KubernetesLabels  map[string]string `json:"kubernetesLabels"`
	KubernetesCluster string            `json:"kubernetesCluster"`
	KubernetesNode    string            `json:"kubernetesNode"`
}

func (host *Host) Labels() map[string]string {
	m := map[string]string{}
	if host.Tags != nil {
		for _, tag := range host.Tags {
			key := strings.TrimSpace(tag.Key)
			if key != "" {
				m[tag.Key] = strings.TrimSpace(tag.Value)
			}
		}
	}
	cluster := strings.TrimSpace(host.KubernetesCluster)
	if cluster != "" {
		m["kubernetesCluster"] = cluster
	}
	node := strings.TrimSpace(host.KubernetesNode)
	if node != "" {
		m["kubernetesNode"] = node
	}
	return m
}
