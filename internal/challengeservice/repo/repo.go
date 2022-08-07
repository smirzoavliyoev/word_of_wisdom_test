package repo

import (
	"errors"
	"strconv"
	"sync"
)

// HERE can use generics
var ErrNoQuotas = errors.New("err no chalanges")

type Repository struct {
	Rwlock    *sync.RWMutex
	Chalanges map[string]string
}

func NewRepo() *Repository {
	return &Repository{
		Rwlock: &sync.RWMutex{},
		Chalanges: map[string]string{
			"chalange1": "abc1",
			"chalange2": "abc2",
			"chalange3": "abc3",
			"chalange4": "abc4",
			"chalange5": "abc5",
		},
	}
}

func (r *Repository) Get() (string, error) {
	r.Rwlock.Lock()
	defer r.Rwlock.Unlock()

	if len(r.Chalanges) == 0 {
		return "", ErrNoQuotas
	}
	for k, v := range r.Chalanges {
		delete(r.Chalanges, k)
		return v, nil
	}

	return "", nil
}

func (r *Repository) Save(challenge string) {
	r.Rwlock.Lock()
	defer r.Rwlock.Unlock()
	r.Chalanges["challenge"+strconv.Itoa(len(r.Chalanges))] = challenge
}
