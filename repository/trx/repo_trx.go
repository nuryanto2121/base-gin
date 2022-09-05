package repotrx

import (
	itrx "app/interface/trx"
	"app/pkg/db"
	"context"

	"gorm.io/gorm"
)

type repoTrx struct {
	db db.DBGormDelegate
}

func NewRepoTrx(Conn db.DBGormDelegate) itrx.Repository {
	return &repoTrx{Conn}
}

// Run implements itrx.Repository
func (r *repoTrx) Run(ctx context.Context, fn func(ctx context.Context) error) error {
	tx := r.db.BeginTx()
	// defer todo: klo panic rollback

	// set tx to context
	newCtx := setContextWithTx(ctx, tx)

	// execute function
	err := fn(newCtx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
func setContextWithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, "tx", tx)
}
