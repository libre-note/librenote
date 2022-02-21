package pgsql

import (
	"database/sql"
	"librenote/app/system/repository"
	"time"
)

type systemRepository struct {
	db *sql.DB
}

func NewPgsqlSystemRepository(db *sql.DB) repository.SystemRepository {
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
