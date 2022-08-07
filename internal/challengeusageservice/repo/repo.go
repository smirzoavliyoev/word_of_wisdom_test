package repo

import (
	"errors"
	"sync"
	"time"
)

var ErrNoChallengeUsageFoundForIP = errors.New("err no challenges found for current Ip")

type Repository struct {
	Rwlock        *sync.RWMutex
	ChalangeUsage map[string]Item
}

type Item struct {
	CreatedAt time.Time
	Value     string
}

func NewRepo() *Repository {
	return &Repository{
		Rwlock:        &sync.RWMutex{},
		ChalangeUsage: map[string]Item{},
	}
}

func (r *Repository) Get(key string) (string, error) {
	r.Rwlock.RLock()
	defer r.Rwlock.RUnlock()
	v, ok := r.ChalangeUsage[key]
	if !ok {
		return "", ErrNoChallengeUsageFoundForIP
	}
	return v.Value, nil
}

func (r *Repository) Clean() {
	r.Rwlock.Lock()
	defer r.Rwlock.Unlock()

	for k, v := range r.ChalangeUsage {
		if time.Since(v.CreatedAt) > time.Minute*20 {
			delete(r.ChalangeUsage, k)
		}
	}
}

func (r *Repository) Save(ip string, challenge string) {
	r.Rwlock.Lock()
	defer r.Rwlock.Unlock()
	r.ChalangeUsage[ip] = Item{
		CreatedAt: time.Now(),
		Value:     ip,
	}
}
