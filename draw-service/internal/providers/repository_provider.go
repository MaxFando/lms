package providers

import "github.com/jmoiron/sqlx"

type RepositoryProvider struct {
	db *sqlx.DB
}

func NewRepositoryProvider(db *sqlx.DB) *RepositoryProvider {
	return &RepositoryProvider{
		db: db,
	}
}

func (r *RepositoryProvider) RegisterDependencies() {

}
