package main

import (
	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/tcp"
	"github.com/smirzoavliyoev/word_of_wisdom_test/pkg/config"
)

type HashcashSolutionVerifierHandler struct {
	server *tcp.Server
}

func NewRequestChallengeMessage(s *tcp.Server) *HashcashSolutionVerifierHandler {

	cfg, err := config.ReadConfig(config.WithSpecificConfigPathOption)

	if err != nil {
		panic(err)
	}

	return &HashcashSolutionVerifierHandler{
		server: tcp.NewServer(cfg),
	}
}
