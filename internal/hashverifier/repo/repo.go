package repo

import (
	"sync"
)

type Repository struct {
	Rwlock     *sync.RWMutex
	UsedHashes map[string]struct{}
}

func NewRepo() *Repository {
	return &Repository{
		Rwlock:     &sync.RWMutex{},
		UsedHashes: make(map[string]struct{}),
	}
}

func (r *Repository) Spent(s string) bool {
	ans := false
	r.Rwlock.RLock()
	_, ans = r.UsedHashes[s]
	r.Rwlock.RUnlock()

	return ans
}

func (r *Repository) Add(s string) error {
	r.Rwlock.Lock()
	r.UsedHashes[s] = struct{}{}
	r.Rwlock.Unlock()
	return nil
}
