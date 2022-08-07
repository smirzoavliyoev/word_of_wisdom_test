package challengeservice

import (
	"errors"

	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/challengeservice/repo"
)

var StorageError = errors.New("unknows error")

type ChallengeService struct {
	Repo *repo.Repository
}

func NewQuotaService() *ChallengeService {
	return &ChallengeService{
		Repo: repo.NewRepo(),
	}
}

func (q *ChallengeService) GetChalange() (string, error) {
	quota, err := q.Repo.Get()
	if err != nil {
		return "", err
	}

	if quota == "" {
		return "", StorageError
	}

	return quota, nil
}
