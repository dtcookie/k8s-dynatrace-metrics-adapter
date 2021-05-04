package cache

type ItemBase struct {
	TSQueried   int64 `json:"-"`
	TSRequested int64 `json:"-"`
}

func (item *ItemBase) GetTSQueried() int64 {
	return item.TSQueried
}

func (item *ItemBase) SetTSQueried(v int64) {
	item.TSQueried = v
}

func (item *ItemBase) GetTSRequested() int64 {
	return item.TSRequested
}

func (item *ItemBase) SetTSRequested(v int64) {
	item.TSRequested = v
}
