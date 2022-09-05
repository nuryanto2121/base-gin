package itrx

import "context"

type Repository interface {
	Run(ctx context.Context, fn func(ctx context.Context) error) error
}
