package topology

import (
	"strings"
)

type ProcessGroup struct {
	EntityBase
	Tags     []*TagInfo `json:"tags"`
	Metadata *Metadata  `json:"metadata"`
}

func (pg *ProcessGroup) Labels() map[string]string {
	m := map[string]string{}
	if pg.Tags != nil {
		for _, tag := range pg.Tags {
			key := strings.TrimSpace(tag.Key)
			if key != "" {
				m[tag.Key] = strings.TrimSpace(tag.Value)
			}
		}
	}
	if pg.Metadata != nil && pg.Metadata.Properties != nil {
		for k, v := range pg.Metadata.Properties {
			key := strings.TrimSpace(k)
			if key != "" {
				m[key] = strings.TrimSpace(v)
			}
		}
	}
	return m
}
