package repo

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrNoChallengeUsageFoundForIP = errors.New("err no challenges found for current Ip")
var once *sync.Once

type Repository struct {
	Rwlock        *sync.RWMutex
	ChalangeUsage map[string]Item

	Expired []string
}

type Item struct {
	CreatedAt time.Time
	Value     string
}

func NewRepo() *Repository {
	repo := &Repository{
		Rwlock:        &sync.RWMutex{},
		ChalangeUsage: map[string]Item{},
	}
	// once.Do(repo.clean)
	return repo
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

// func (r *Repository) clean() {
// 	go func() {
// 		for {
// 			time.Sleep(5 * time.Minute)
// 			r.Rwlock.Lock()

// 			for k, v := range r.ChalangeUsage {
// 				if time.Since(v.CreatedAt) > time.Minute*20 {
// 					r.Expired = append(r.Expired, v.Value)
// 					delete(r.ChalangeUsage, k)
// 				}
// 			}
// 			r.Rwlock.Unlock()
// 		}
// 	}()
// }

func (r *Repository) Save(ip string, challenge string) {
	r.Rwlock.Lock()
	defer r.Rwlock.Unlock()
	r.ChalangeUsage[ip] = Item{
		CreatedAt: time.Now(),
		Value:     challenge,
	}

	fmt.Println(r.ChalangeUsage)
}

func (r *Repository) EmptyExpiredData() {
	r.Rwlock.Lock()
	r.Expired = nil
	r.Rwlock.Unlock()
}
