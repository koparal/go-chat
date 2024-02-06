package user

import (
	"chat/internal/utils"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type service struct {
	Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository,
	}
}

var timeOut = time.Duration(3) * time.Second

func (s *service) Register(c context.Context, req *CreateUserData) (*CreateUserResponse, error) {
	ctx, cancel := context.WithTimeout(c, timeOut)
	defer cancel()

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	if req.Username == "" {
		return nil, fmt.Errorf("username required field")
	}

	u := &User{
		Username: req.Username,
		Password: hashedPassword,
	}

	r, err := s.Repository.Register(ctx, u)
	if err != nil {
		return nil, err
	}

	res := &CreateUserResponse{
		ID:       strconv.Itoa(int(r.ID)),
		Username: r.Username,
	}

	return res, nil
}

type MyJWTClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

func (s *service) Login(c context.Context, req *LoginUserData) (*LoginUserResponse, error) {
	ctx, cancel := context.WithTimeout(c, timeOut)
	defer cancel()

	if req.Username == "" {
		return &LoginUserResponse{}, fmt.Errorf("username required field")
	}

	u, err := s.Repository.CheckExistsUser(ctx, req.Username)
	if err != nil {
		return &LoginUserResponse{}, err
	}

	err = utils.CheckPassword(req.Password, u.Password)
	if err != nil {
		fmt.Println(err)
		return &LoginUserResponse{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		ID:       strconv.Itoa(int(u.ID)),
		Username: u.Username,
		IsAdmin:  u.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(u.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	accessToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		return &LoginUserResponse{}, err
	}

	return &LoginUserResponse{
		Username:    u.Username,
		ID:          strconv.Itoa(int(u.ID)),
		AccessToken: accessToken,
		IsAdmin:     u.IsAdmin,
	}, nil
}
