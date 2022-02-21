package repository

type SystemRepository interface {
	DBCheck() (bool, error)
	CurrentTime() int64
}
