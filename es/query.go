package es

import (
	"context"
	"fmt"
	"math"

	"github.com/contextcloud/eventstore/es/filters"
	"github.com/contextcloud/eventstore/es/utils"
	"github.com/google/uuid"
)

type Pagination[T any] struct {
	Limit      int `json:"limit"`
	Page       int `json:"page"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
	Items      []T `json:"items"`
}

type Query[T Entity] interface {
	Load(ctx context.Context, id uuid.UUID) (T, error)
	Find(ctx context.Context, filter filters.Filter) ([]T, error)
	Count(ctx context.Context, filter filters.Filter) (int, error)
	Pagination(ctx context.Context, filter filters.Filter) (*Pagination[T], error)
}

type query[T Entity] struct {
	name string
}

func (q *query[T]) Load(ctx context.Context, id uuid.UUID) (T, error) {
	var item T

	unit, err := GetUnit(ctx)
	if err != nil {
		return item, err
	}

	out, err := unit.Load(ctx, q.name, id)
	if err != nil {
		return item, err
	}

	result, ok := out.(T)
	if !ok {
		return item, fmt.Errorf("unexpected type: %T", out)
	}
	return result, nil
}

func (q *query[T]) Save(ctx context.Context, entities ...T) error {
	unit, err := GetUnit(ctx)
	if err != nil {
		return err
	}

	for _, entity := range entities {
		if err := unit.Save(ctx, q.name, entity); err != nil {
			return err
		}
	}
	return nil
}

func (q *query[T]) Find(ctx context.Context, filter filters.Filter) ([]T, error) {
	unit, err := GetUnit(ctx)
	if err != nil {
		return nil, err
	}

	var items []T
	if err := unit.Find(ctx, q.name, filter, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (q *query[T]) Count(ctx context.Context, filter filters.Filter) (int, error) {
	unit, err := GetUnit(ctx)
	if err != nil {
		return 0, err
	}

	return unit.Count(ctx, q.name, filter)
}

func (q *query[T]) Pagination(ctx context.Context, filter filters.Filter) (*Pagination[T], error) {
	if filter.Limit == nil {
		return nil, fmt.Errorf("Limit required for pagination")
	}
	if filter.Offset == nil {
		return nil, fmt.Errorf("Offset required for pagination")
	}

	unit, err := GetUnit(ctx)
	if err != nil {
		return nil, err
	}

	totalItems, err := unit.Count(ctx, q.name, filter)
	if err != nil {
		return nil, err
	}

	var items []T
	if err := unit.Find(ctx, q.name, filter, &items); err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(*filter.Limit)))
	return &Pagination[T]{
		Limit:      *filter.Limit,
		Page:       *filter.Offset + 1,
		TotalItems: totalItems,
		TotalPages: totalPages,
		Items:      items,
	}, nil
}

func NewQuery[T Entity]() Query[T] {
	var item T
	typeOf := utils.GetElemType(item)

	return &query[T]{
		name: typeOf.String(),
	}
}
