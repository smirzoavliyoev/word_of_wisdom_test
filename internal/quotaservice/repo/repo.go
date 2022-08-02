package repo

import (
	"errors"
	"sync"
)

var ErrNoQuotas = errors.New("err no quotas")

type Repository struct {
	Rwlock *sync.RWMutex
	Quotas map[string]string
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
	}
}

func (r *Repository) Get() (string, error) {
	if len(r.Quotas) == 0 {
		return "", ErrNoQuotas
	}
	r.Rwlock.Lock()
	for k, v := range r.Quotas {
		delete(r.Quotas, k)
		return v, nil
	}
	r.Rwlock.Unlock()

	return "", nil
}
