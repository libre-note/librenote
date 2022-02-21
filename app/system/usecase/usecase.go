package usecase

import "librenote/app/system/repository"

type SystemUsecase interface {
	GetHealth() (*HealthResp, error)
	GetTime() *TimeResp
}

type systemUsecase struct {
	repo repository.SystemRepository
}

func NewSystemUsecase(repo repository.SystemRepository) SystemUsecase {
	return &systemUsecase{
		repo: repo,
	}
}

func (u *systemUsecase) GetHealth() (*HealthResp, error) {
	resp := HealthResp{}

	// check db
	dbOnline, err := u.repo.DBCheck()
	resp.DBOnline = dbOnline
	if err != nil {
		return &resp, err
	}

	return &resp, nil
}

func (u *systemUsecase) GetTime() *TimeResp {
	resp := TimeResp{}
	resp.CurrentTimeUnix = u.repo.CurrentTime()
	return &resp
}

type TimeResp struct {
	CurrentTimeUnix int64 `json:"current_time_unix"`
}

type HealthResp struct {
	DBOnline    bool `json:"db_online"`
	CacheOnline bool `json:"cache_online"`
}
