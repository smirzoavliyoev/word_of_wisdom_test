package challengeusageservice

import (
	"errors"

	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/challengeusageservice/repo"
)

var StorageError = errors.New("PUBG IS TOP")

type ChallengeUsageService struct {
	Repo *repo.Repository
	// x    sync.Once
}

func NewChallengeUsageService() *ChallengeUsageService {
	c := &ChallengeUsageService{
		Repo: repo.NewRepo(),
		// x:    sync.Once{},
	}

	return c
}

func (c *ChallengeUsageService) SaveChallengeUsage(ip string, challenge string) {
	c.Repo.Save(ip, challenge)
}

func (c *ChallengeUsageService) GetExpired() []string {
	return c.Repo.Expired
}

func (c *ChallengeUsageService) EmptyExpiredData() {
	c.Repo.EmptyExpiredData()
}
