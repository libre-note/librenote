package pgsql

import (
	"context"
	"librenote/app/model"
	repo "librenote/app/user/repository/pgsql"
	"time"
)

func (s *PgsqlRepositoryTestSuite) TestPgsqlUserRepository_CreateUser() {
	nowTime := time.Now().UTC().Format("2006-01-02 15:04:05")

	newUser := &model.User{
		FullName:  "Mr. Test",
		Email:     "mrtest@example.com",
		Hash:      "sfj34ksfdsfj$24247834skfjskdf",
		IsActive:  1,
		IsTrashed: 0,
		CreatedAt: nowTime,
		UpdatedAt: nowTime,
	}

	r := repo.NewPgsqlUserRepository(s.db)
	s.Assert().NoError(r.CreateUser(context.Background(), newUser))

	query := "SELECT id, full_name, created_at::text FROM users LIMIT 1"
	row := s.db.QueryRowContext(context.Background(), query)

	var res model.User
	err := row.Scan(
		&res.ID,
		&res.FullName,
		&res.CreatedAt,
	)

	s.Assert().NoError(err)
	s.Assert().Equal(newUser.FullName, res.FullName)
	s.Assert().Equal(newUser.CreatedAt, res.CreatedAt)
}

func (s *PgsqlRepositoryTestSuite) TestPgsqlUserRepository_GetUserByID() {
	nowTime := time.Now().UTC().Format("2006-01-02 15:04:05")

	newUser := &model.User{
		FullName:  "Mr. Test",
		Email:     "mrtest@example.com",
		Hash:      "abc123",
		IsActive:  1,
		IsTrashed: 0,
		CreatedAt: nowTime,
		UpdatedAt: nowTime,
	}

	r := repo.NewPgsqlUserRepository(s.db)
	s.Assert().NoError(r.CreateUser(context.Background(), newUser))

	id := 1 // new id is created

	result, err := r.GetUser(context.Background(), int32(id))
	s.Assert().NoError(err)
	s.Assert().Equal(newUser.FullName, result.FullName)
	s.Assert().Equal(newUser.CreatedAt, newUser.CreatedAt)
}

func (s *PgsqlRepositoryTestSuite) TestPgsqlUserRepository_GetUserByEmail() {
	nowTime := time.Now().UTC().Format("2006-01-02 15:04:05")

	firstUser := &model.User{
		FullName:  "Mr. Test 1",
		Email:     "mrtest1@example.com",
		Hash:      "abc123",
		IsActive:  1,
		IsTrashed: 0,
		CreatedAt: nowTime,
		UpdatedAt: nowTime,
	}

	secondUser := &model.User{
		FullName:  "Mr. Test 2",
		Email:     "mrtest2@example.com",
		Hash:      "abc123",
		IsActive:  1,
		IsTrashed: 0,
		CreatedAt: nowTime,
		UpdatedAt: nowTime,
	}

	r := repo.NewPgsqlUserRepository(s.db)
	s.Assert().NoError(r.CreateUser(context.Background(), firstUser))
	s.Assert().NoError(r.CreateUser(context.Background(), secondUser))

	result, err := r.GetUserByEmail(context.Background(), secondUser.Email)
	s.Assert().NoError(err)
	s.Assert().Equal(secondUser.FullName, result.FullName)
	s.Assert().Equal(secondUser.IsActive, result.IsActive)
}

func (s *PgsqlRepositoryTestSuite) TestPgsqlUserRepository_UpdateUser() {
	nowTime := time.Now().UTC().Format("2006-01-02 15:04:05")

	newUser := &model.User{
		FullName:  "Mr. Test 1",
		Email:     "mrtest1@example.com",
		Hash:      "abc123",
		IsActive:  1,
		IsTrashed: 0,
		CreatedAt: nowTime,
		UpdatedAt: nowTime,
	}

	r := repo.NewPgsqlUserRepository(s.db)
	s.Assert().NoError(r.CreateUser(context.Background(), newUser))

	updateUser := newUser
	updateUser.ID = 1
	updateUser.Hash = "changed_hash"
	updateUser.DarkModeEnabled = 1
	updateUser.UpdatedAt = time.Now().UTC().Format("2006-01-02 15:04:05")

	s.Assert().NoError(r.UpdateUser(context.Background(), updateUser))

	result, err := r.GetUserByEmail(context.Background(), newUser.Email)
	s.Assert().NoError(err)
	s.Assert().Equal(updateUser.Hash, result.Hash)
	s.Assert().Equal(updateUser.DarkModeEnabled, result.DarkModeEnabled)
}
