package repository

import (
	"database/sql"
	"time"
)

type SystemRepository interface {
	DBCheck() (bool, error)
	CurrentTime() int64
}

type systemRepository struct {
	db *sql.DB
}

func NewSystemRepository(db *sql.DB) SystemRepository {
	return &systemRepository{
		db: db,
	}
}

func (r *systemRepository) DBCheck() (bool, error) {
	if err := r.db.Ping(); err == nil {
		return true, nil
	}

	return false, nil
}

func (r *systemRepository) CurrentTime() int64 {
	return time.Now().Unix()
}
