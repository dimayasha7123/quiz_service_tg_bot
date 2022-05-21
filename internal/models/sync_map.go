package models

import (
	"sync"
)

type SyncMap struct {
	sync.RWMutex
	M map[int64]*User
}

func NewSyncMap() *SyncMap {
	return &SyncMap{M: make(map[int64]*User)}
}
