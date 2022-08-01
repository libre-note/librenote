package usecase

import (
	"errors"
	"librenote/app/response"
	"librenote/app/system/repository"
)

type SystemUsecase interface {
	GetHealth() error
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

func (u *systemUsecase) GetHealth() error {
	// check db
	dbOnline, err := u.repo.DBCheck()
	if err != nil {
		return err
	}

	if !dbOnline {
		return response.WrapError(errors.New("DB is offline"), 503)
	}

	return nil
}

func (u *systemUsecase) GetTime() *TimeResp {
	resp := TimeResp{}
	resp.CurrentTimeUnix = u.repo.CurrentTime()

	return &resp
}

type TimeResp struct {
	CurrentTimeUnix int64 `json:"current_time_unix"`
}
