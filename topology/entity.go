package topology

import "github.com/dtcookie/k8s-dynatrace-metrics-adapter/cache"

type Entity interface {
	Labels() map[string]string
}

type EntityBase struct {
	cache.ItemBase
	ID          string `json:"entityId"`
	DisplayName string `json:"displayName"`
}

func (entity *EntityBase) Name() string {
	return entity.DisplayName
}

func (entity *EntityBase) GetID() string {
	return entity.ID
}
