package repo

import (
	"errors"
	"sync"
)

type Repository struct {
	Rwlock     *sync.RWMutex
	Quotas     map[string]string
	UsedHashes map[string]struct{}
}

func NewRepo() *Repository {
	return &Repository{
		Rwlock: &sync.RWMutex{},
		Quotas: map[string]string{
			"quota1": "abc1",
			"quota2": "abc2",
			"quota3": "abc3",
			"quota4": "abc4",
			"quota5": "abc5",
		},
		UsedHashes: make(map[string]struct{}),
	}
}

func (r *Repository) Get() (string, error) {
	if len(r.Quotas) == 0 {
		return "", errors.New("no quotas")
	}
	r.Rwlock.Lock()
	for k, v := range r.Quotas {
		delete(r.Quotas, k)
		return v, nil
	}
	r.Rwlock.Unlock()

	return "", nil
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
