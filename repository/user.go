package repository

import (
	"context"
	"expense-tracker/entity"
	"time"
)

type UserRepoItf interface {
	InsertUser(ctx context.Context, user entity.InsertUserReq) (*int, error)
	IsEmailExists(ctx context.Context, email string) bool
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

func (ur UserRepoStruct) IsEmailExists(ctx context.Context, email string) bool {
	var emailExists bool
	sql := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	db := ChooseDbOrTx(ctx, ur.db)
	err := db.QueryRowContext(ctx, sql, email).Scan(&emailExists)
	if err != nil {
		return true
	}

	return emailExists
}
