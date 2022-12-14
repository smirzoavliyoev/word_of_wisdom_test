package challengeservice

import (
	"errors"

	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/challengeservice/repo"
)

var StorageError = errors.New("unknows error")

type ChallengeService struct {
	Repo *repo.Repository
}

func NewChallengeService() *ChallengeService {
	return &ChallengeService{
		Repo: repo.NewRepo(),
	}
}

func (q *ChallengeService) GetChallenge() (string, error) {
	quota, err := q.Repo.Get()
	if err != nil {
		return "", err
	}

	if quota == "" {
		return "", StorageError
	}

	return quota, nil
}

func (q *ChallengeService) SaveChallenge(challenge string) {
	q.Repo.Save(challenge)
}

func (c *ChallengeService) SaveChunk(ch []string) {
	c.Repo.SaveChunk(ch)
}
