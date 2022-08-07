package quoteservice

import (
	"errors"
	"time"

	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/challengeusageservice/repo"
)

var StorageError = errors.New("PUBG IS TOP")

type ChallengeUsageService struct {
	Repo *repo.Repository
}

func NewChallengeUsageService() *ChallengeUsageService {

	return &ChallengeUsageService{
		Repo: repo.NewRepo(),
	}
}

func (c *ChallengeUsageService) clean() {

	for {
		time.Sleep(20 * time.Minute)
		c.Repo.Clean()
	}

}
