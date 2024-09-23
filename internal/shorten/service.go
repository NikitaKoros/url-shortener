package shorten

import (
	"context"
	"log"
	"github.com/google/uuid"
	"url-shortener/internal/model"
)

type Storage interface {
	Put(ctx context.Context, shortening model.Shortening) (*model.Shortening, error)
	Get(ctx context.Context, identifier string) (*model.Shortening, error)
	IncrementVisits(ctx context.Context, identifier string) error
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) Shorten(ctx context.Context, input model.ShortenInput) (*model.Shortening, error) {
	var (
		id = uuid.New().ID()
		identifier string
	)

	if input.Identifier == "" {
		identifier = Shorten(id)
	} else {
		identifier = input.Identifier
	}

	dbShortening := model.Shortening {
		Identifier: identifier,
		OriginalURL: input.RawURL,
	}

	shortening, err := s.storage.Put(ctx, dbShortening)
	if err != nil {
		return nil, err
	}

	return shortening, nil
}

func (s *Service) Get(ctx context.Context, identifier string) (*model.Shortening, error) {
	shortening, err := s.storage.Get(ctx, identifier)
	if err != nil {
		return nil, err
	}

	return shortening, nil
}

func (s *Service) Redirect(ctx context.Context, identifier string) (string, error) {
	shortening, err := s.storage.Get(ctx, identifier)
	if err != nil {
		return "", err
	}

	if err := s.storage.IncrementVisits(ctx, identifier); err != nil {
		log.Printf("failed to increment visits for identifier %q: %v", identifier, err)
	}

	return shortening.OriginalURL, nil
}