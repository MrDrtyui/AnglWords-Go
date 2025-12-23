package word

import (
	"app/domain"
	"app/ent"
	"app/internal/neuronet"
	"context"
)

type Service struct {
	Repo         Repository
	NeuroService neuronet.Repository
}

func NewService(repo Repository, neruRepo neuronet.Repository) *Service {
	return &Service{Repo: repo, NeuroService: neruRepo}
}

func (s *Service) GetWord(ctx context.Context, word string) (*ent.Word, error) {
	return s.Repo.GetWord(ctx, word)
}

func (s *Service) GetMyWords(ctx context.Context, userId int) ([]*ent.Word, error) {
	return s.Repo.GetMyWords(ctx, userId)
}

func (s *Service) GetAllWords(ctx context.Context) ([]*ent.Word, error) {
	return s.Repo.GetAllWords(ctx)
}
func (s *Service) CreateWord(ctx context.Context, wordDto domain.WordCreateDto, userId int) (*ent.Word, error) {

	data, err := s.NeuroService.Translate(ctx, wordDto.Word)
	if err != nil {
		return nil, err
	}

	return s.Repo.CreateWord(ctx, userId, wordDto.Word, data.Level, data.WordRu)
}
