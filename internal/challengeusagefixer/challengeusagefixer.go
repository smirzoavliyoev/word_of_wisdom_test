package challengeusagefixer

import "sync"

type ChallengeUsageFixer struct {
	once sync.Once
}
