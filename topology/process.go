package topology

import (
	"encoding/json"
	"strings"
)

type Process struct {
	EntityBase
	Tags     []*TagInfo `json:"tags"`
	Metadata *Metadata  `json:"metadata"`
}

func (process *Process) Labels() map[string]string {
	m := map[string]string{}
	if process.Tags != nil {
		for _, tag := range process.Tags {
			key := strings.TrimSpace(tag.Key)
			if key != "" {
				m[tag.Key] = strings.TrimSpace(tag.Value)
			}
		}
	}
	if process.Metadata != nil && process.Metadata.Properties != nil {
		for k, v := range process.Metadata.Properties {
			key := strings.TrimSpace(k)
			if key != "" {
				m[key] = strings.TrimSpace(v)
			}
		}
	}
	return m
}

type Metadata struct {
	Properties map[string]string
}

func (md *Metadata) UnmarshalJSON(data []byte) error {
	md.Properties = map[string]string{}
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	for k, v := range m {
		arr := []string{}
		if err := json.Unmarshal(v, &arr); err != nil {
			return err
		}
		if len(arr) == 1 {
			md.Properties[k] = arr[0]
		}

	}
	return nil
}
