package user

import (
	"context"
	"server/util"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repository Repository) Service {
	return &service{
		Repository: repository,
		timeout:    time.Duration(2) * time.Second,
	}
}

func (s *service) CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	newUser, err := s.Repository.CreateUser(ctx, &User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	})
	if err != nil {
		return nil, err
	}

	return &CreateUserRes{
		ID:       strconv.Itoa(int(newUser.ID)),
		Username: newUser.Username,
		Email:    newUser.Email,
	}, nil
}

type JWTClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

const secretKey = "HEHE-Super-Secret"

func (s *service) LoginUser(c context.Context, req *LoginUserReq) (*LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	existingUser, err := s.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	err = util.ComparePassword(req.Password, existingUser.Password)
	if err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		ID:       strconv.Itoa(int(existingUser.ID)),
		Username: existingUser.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(existingUser.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	return &LoginUserRes{
		accessToken: ss,
		ID:          strconv.Itoa(int(existingUser.ID)),
		Username:    existingUser.Username,
	}, nil
}
