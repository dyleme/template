package example

import (
	"context"
	"fmt"

	"github.com/dyleme/template/internal/domain"
	"github.com/dyleme/template/pkg/txmanager"
)

type Repository struct {
	getter txmanager.TxGetter
}

func NewRepository(getter txmanager.TxGetter) *Repository {
	return &Repository{getter: getter}
}

func (r *Repository) Get(ctx context.Context, id int) (domain.Example, error) {
	tx := r.getter.GetTx(ctx)

	row := tx.QueryRow(ctx, "SELECT * FROM example WHERE $1", id)
	var example domain.Example
	if err := row.Scan(&example); err != nil {
		return domain.Example{}, fmt.Errorf("query row: %w", err)
	}

	return example, nil
}

func (r *Repository) List(ctx context.Context) ([]domain.Example, error) {
	tx := r.getter.GetTx(ctx)

	rows, err := tx.Query(ctx, "SELECT * FROM example")
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var examples []domain.Example
	for rows.Next() {
		var example domain.Example
		if err := rows.Scan(&example); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		examples = append(examples, example)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows err: %w", rows.Err())
	}

	return examples, nil
}

func (r *Repository) Create(ctx context.Context, example domain.Example) error {
	tx := r.getter.GetTx(ctx)

	_, err := tx.Exec(ctx, "INSERT INTO example (id, name) VALUES ($1, $2)", example.ID, example.Name)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, example domain.Example) error {
	tx := r.getter.GetTx(ctx)

	_, err := tx.Exec(ctx, "UPDATE example SET name = $1 WHERE id = $2", example.Name, example.ID)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}
