package quotaservice

import (
	"errors"

	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/quotaservice/repo"
)

var StorageError = errors.New("unknows error")

type QuotaService struct {
	Repo *repo.Repository
}

func NewQuotaService() *QuotaService {
	return &QuotaService{
		Repo: repo.NewRepo(),
	}
}

func (q *QuotaService) GetQuota() (string, error) {
	quota, err := q.Repo.Get()
	if err != nil {
		return "", err
	}

	if quota == "" {
		return "", StorageError
	}

	return quota, nil
}
