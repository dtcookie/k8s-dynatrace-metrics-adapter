package cache

type Item interface {
	// GetID() string
	SetTSQueried(int64)
	GetTSQueried() int64
	SetTSRequested(int64)
	GetTSRequested() int64
}
