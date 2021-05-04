package topology

import (
	"strings"
)

type Service struct {
	EntityBase
	Tags []*TagInfo `json:"tags"`
}

func (service *Service) Labels() map[string]string {
	m := map[string]string{}
	if service.Tags != nil {
		for _, tag := range service.Tags {
			key := strings.TrimSpace(tag.Key)
			if key != "" {
				m[tag.Key] = strings.TrimSpace(tag.Value)
			}
		}
	}
	return m
}
