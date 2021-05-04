package cache

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type Repo interface {
	Get(id string, options map[string]interface{}) (Item, error)
}

type repo struct {
	name      string
	maxAge    int
	fnGet     func(id string, options map[string]interface{}) (Item, error)
	fnOnError func(e error) Item
	items     sync.Map
	channel   chan *request
}

type request struct {
	ID      string
	Options map[string]interface{}
}

func NewRepo(name string, maxAge int, fnGet func(id string, options map[string]interface{}) (Item, error), fnOnError func(e error) Item) Repo {
	r := &repo{name: name, maxAge: maxAge, items: sync.Map{}, fnGet: fnGet, fnOnError: fnOnError, channel: make(chan *request, 100)}
	go r.fetcher()
	return r
}

func (rep *repo) fetcher() {
	for v := range rep.channel {
		item := rep.find(v.ID, v.Options)
		if item == nil || item.GetTSRequested()-item.GetTSQueried() > int64(rep.maxAge) {
			rep.fetch(v.ID, v.Options)
		}
	}
}

func (rep *repo) fetch(id string, options map[string]interface{}) {
	var err error
	var item Item
	if item, err = rep.fnGet(id, options); err != nil {
		if rep.fnOnError == nil {
			return
		}
		if item = rep.fnOnError(err); item == nil {
			return
		}
	}
	if item == nil {
		return
	}
	item.SetTSRequested(time.Now().Unix())
	item.SetTSQueried(item.GetTSRequested())
	// data, _ := json.Marshal(item)
	// klog.Info(fmt.Sprintf("[%v] STORE %v -> %v", rep.name, keyFor(id, options), string(data)))
	rep.items.Store(keyFor(id, options), item)
}

func keyFor(id string, options map[string]interface{}) string {
	key := id

	if len(options) > 0 {
		for k, v := range options {
			data, _ := json.Marshal(v)
			key = fmt.Sprintf("%s#%v:%v", key, k, string(data))
		}
	}
	return key
}

func (rep *repo) find(id string, options map[string]interface{}) Item {
	if uitem, found := rep.items.Load(keyFor(id, options)); found {
		item := uitem.(Item)
		item.SetTSRequested(time.Now().Unix())
		return item
	}
	return nil
}

func (rep *repo) Get(id string, options map[string]interface{}) (Item, error) {
	// klog.Info(fmt.Sprintf("[%v] Get %v", rep.name, keyFor(id, options)))
	uitem := rep.find(id, options)
	if uitem != nil {
		item := uitem.(Item)
		if item.GetTSRequested()-item.GetTSQueried() > int64(rep.maxAge) {
			rep.channel <- &request{ID: id, Options: options}
		}
		// data, _ := json.Marshal(item)
		// klog.Info(fmt.Sprintf("[%v]   => %v", rep.name, string(data)))
		return item, nil
	}
	rep.channel <- &request{ID: id, Options: options}
	return nil, nil
}
