package challengeusagefixer

import (
	"sync"
	"time"

	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/challengeservice"
	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/challengeusageservice"
)

// may be it will be cool to join challenge services and create one repo....
// but it will be bad if in future we need to divide responsibilities
// so ill keep functionalities atomic
var once *sync.Once

type ChallengeUsageFixer struct {
	ChallengeUsageService *challengeusageservice.ChallengeUsageService
	Challengeservice      *challengeservice.ChallengeService
}

func NewChallengeUsageFixer(
	cs *challengeservice.ChallengeService,
	cus *challengeusageservice.ChallengeUsageService,
) *ChallengeUsageFixer {

	cuf := &ChallengeUsageFixer{
		ChallengeUsageService: cus,
		Challengeservice:      cs,
	}
	once.Do(cuf.fixData)
	return cuf
}

func (c *ChallengeUsageFixer) fixData() {
	go func() {
		for {
			time.Sleep(20 * time.Minute)

			expired := c.ChallengeUsageService.GetExpired()
			// good 1 systemcall to lock
			// prev - used for....too many systemcalls
			c.Challengeservice.SaveChunk(expired)

			c.ChallengeUsageService.EmptyExpiredData()
		}
	}()
}
