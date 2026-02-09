package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	User UserRepository
	// Add other repositories here
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		User: NewUserRepository(db),
	}
}
