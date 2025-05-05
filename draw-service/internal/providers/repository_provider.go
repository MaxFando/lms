package providers

import (
	"github.com/MaxFando/lms/draw-service/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
)

type RepositoryProvider struct {
	db *sqlx.DB

	drawRepository *postgres.DrawRepository
}

func NewRepositoryProvider(db *sqlx.DB) *RepositoryProvider {
	return &RepositoryProvider{
		db: db,
	}
}

func (r *RepositoryProvider) RegisterDependencies() {
	r.drawRepository = postgres.NewDrawRepository(r.db)
}
