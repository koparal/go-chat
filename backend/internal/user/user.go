package user

import "context"

type Repository interface {
	Register(ctx context.Context, user *User) (*User, error)
	CheckExistsUser(ctx context.Context, username string) (*User, error)
}

type Service interface {
	Register(c context.Context, req *CreateUserData) (*CreateUserResponse, error)
	Login(c context.Context, req *LoginUserData) (*LoginUserResponse, error)
}
