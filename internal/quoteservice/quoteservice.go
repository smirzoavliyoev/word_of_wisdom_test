package quoteservice

import (
	"errors"

	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/quoteservice/repo"
)

var StorageError = errors.New("unknows error")

type QuoteService struct {
	Repo *repo.Repository
}

func NewQuotaService() *QuoteService {
	return &QuoteService{
		Repo: repo.NewRepo(),
	}
}

func (q *QuoteService) GetQuota() (string, error) {
	quota, err := q.Repo.Get()
	if err != nil {
		return "", err
	}

	if quota == "" {
		return "", StorageError
	}

	return quota, nil
}
