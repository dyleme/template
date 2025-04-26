package example

import (
	"context"
	"fmt"

	"github.com/dyleme/template/internal/domain"
	"github.com/dyleme/template/pkg/txmanager"
)

type Repository interface {
	Get(ctx context.Context, id int) (domain.Example, error)
	List(ctx context.Context) ([]domain.Example, error)
	Create(ctx context.Context, example domain.Example) error
	Update(ctx context.Context, example domain.Example) error
}

type Service struct {
	repo      Repository
	txManager txmanager.TxManager
}

func NewService(repo Repository, txManager txmanager.TxManager) *Service {
	return &Service{repo: repo, txManager: txManager}
}

func (s *Service) Get(ctx context.Context, id int) (domain.Example, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) Update(ctx context.Context, params domain.Example) error {
	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		exmpl, err := s.repo.Get(ctx, params.ID)
		if err != nil {
			return fmt.Errorf("get: %w", err)
		}

		exmpl.Name = params.Name
		err = s.repo.Update(ctx, exmpl)
		if err != nil {
			return fmt.Errorf("update: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("do: %w", err)
	}

	return nil
}
