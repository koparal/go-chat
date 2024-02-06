package user

import (
	"context"
	"database/sql"
	"fmt"
)

type PSqlCxt interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	db PSqlCxt
}

func NewRepository(db PSqlCxt) Repository {
	return &repository{db: db}
}

func (r *repository) Register(ctx context.Context, user *User) (*User, error) {
	var lastInsertId int
	query := "INSERT INTO users(username, password) VALUES ($1, $2) returning id"
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Password).Scan(&lastInsertId)
	if err != nil {
		return &User{}, err
	}

	user.ID = int64(lastInsertId)
	return user, nil
}

func (r *repository) CheckExistsUser(ctx context.Context, username string) (*User, error) {
	u := User{}
	query := "SELECT * FROM users WHERE username = $1"
	err := r.db.QueryRowContext(ctx, query, username).Scan(&u.ID, &u.Username, &u.Password, &u.IsAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found.")
		}
		return nil, err
	}

	return &u, nil
}
