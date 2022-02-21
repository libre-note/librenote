package response

import (
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
)

var (
	ErrNotFound            = errors.New("resource not found")
	ErrInvalidPage         = errors.New("invalid page request")
	ErrConflict            = errors.New("data conflict or already exist")
	ErrBadRequest          = errors.New("bad request, check param or body")
	ErrInternalServerError = errors.New("internal server error")
)

func getStatusCode(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrInvalidPage:
		return http.StatusNotFound
	case ErrConflict:
		return http.StatusConflict
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrInternalServerError:
		return http.StatusInternalServerError
	default:
		wrapErr := &wrapErr{}
		if errors.As(err, wrapErr) {
			return wrapErr.StatusCode
		}
		return http.StatusInternalServerError
	}
}

// RespondError takes an `error` and a `customErr message` args
// to log the error to system and return to client
func RespondError(err error, customErr ...error) (int, Response) {
	logrus.Errorln(err, customErr)
	if len(customErr) > 0 {
		return getStatusCode(err), Response{Success: false, Message: customErr[0].Error()}
	}
	return getStatusCode(err), Response{Success: false, Message: err.Error()}
}

type wrapErr struct {
	StatusCode int
	Err        error
}

// implements error interface
func (e wrapErr) Error() string {
	return e.Err.Error()
}

// Unwrap Implements the errors.Unwrap interface
func (e wrapErr) Unwrap() error {
	return e.Err // Returns inner error
}

func WrapError(err error, statusCode int, errCode string) error {
	return wrapErr{
		Err:        err,
		StatusCode: statusCode,
	}
}
