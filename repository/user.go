package repository

import (
	"context"
	"expense-tracker/entity"
	"time"
)

type UserRepoItf interface {
	InsertUser(ctx context.Context, user entity.InsertUserReq) (*int, error)
	IsEmailExists(ctx context.Context, email string) (bool, error)
	IsEmailNotExists(ctx context.Context, email string) (bool, error)
	SelectHashPasswordByEmail(ctx context.Context, email string) (*string, error)
	SelectIdByEmail(ctx context.Context, email string) (*int, error)
}

type UserRepoStruct struct {
	db DBTX
}

func NewUserRepo(db DBTX) UserRepoStruct {
	return UserRepoStruct{
		db: db,
	}
}

func (ur UserRepoStruct) InsertUser(ctx context.Context, user entity.InsertUserReq) (*int, error) {
	sql := `INSERT INTO users(name, email, password, created_at, updated_at)
		VALUES
		($1, $2, $3, $4, $5)
		RETURNING (id)
	`

	var userId int
	db := ChooseDbOrTx(ctx, ur.db)
	err := db.QueryRowContext(ctx, sql, &user.Name, &user.Email, &user.Password, time.Now(), time.Now()).Scan(&userId)
	if err != nil {
		return nil, err
	}

	return &userId, nil
}

func (ur UserRepoStruct) IsEmailExists(ctx context.Context, email string) (bool, error) {
	var emailExists bool
	sql := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	db := ChooseDbOrTx(ctx, ur.db)
	err := db.QueryRowContext(ctx, sql, email).Scan(&emailExists)
	if err != nil {
		return true, err
	}
	return emailExists, nil
}

func (ur UserRepoStruct) IsEmailNotExists(ctx context.Context, email string) (bool, error) {
	var emailExists bool
	sql := `SELECT NOT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	db := ChooseDbOrTx(ctx, ur.db)
	err := db.QueryRowContext(ctx, sql, email).Scan(&emailExists)
	if err != nil {
		return true, err
	}
	return emailExists, nil
}

func (ur UserRepoStruct) SelectHashPasswordByEmail(ctx context.Context, email string) (*string, error) {
	var hashPassword string
	sql := `SELECT password FROM users WHERE email = $1`
	db := ChooseDbOrTx(ctx, ur.db)
	err := db.QueryRowContext(ctx, sql, email).Scan(&hashPassword)
	if err != nil {
		return nil, err
	}
	return &hashPassword, nil
}

func (ur UserRepoStruct) SelectIdByEmail(ctx context.Context, email string) (*int, error) {
	var id int
	sql := `SELECT id FROM users WHERE email = $1`
	db := ChooseDbOrTx(ctx, ur.db)
	err := db.QueryRowContext(ctx, sql, email).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
