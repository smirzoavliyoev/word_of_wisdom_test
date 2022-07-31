package service

import (
	"fmt"

	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/repo"
)

type Service struct {
	Repo *repo.Repository
}

func NewService(repo *repo.Repository) *Service {
	return &Service{
		Repo: repo,
	}
}

func (s *Service) GetQuota() (string, error) {
	v, err := s.Repo.Get()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return v, nil
}
