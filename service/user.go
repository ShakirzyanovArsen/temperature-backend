package service

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"temperature-backend/model"
	"temperature-backend/repository"
	"temperature-backend/service/err"
)

type UserService interface {
	Register(email string) (*model.User, *Error)
}

type userServiceImpl struct {
	repo *repository.UserRepository
	tokenGenerator tokenGenerator
}

type tokenGenerator interface {
	getToken() (string, error)
}

type userTokenGenerator struct {}

func (userTokenGenerator) getToken() (string, error) {
	b := make([]byte, 256)
	_, e := rand.Read(b)
	if e != nil {
		return "", e
	}
	hash := sha1.New()
	hash.Write(b)
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (s userServiceImpl) Register(email string) (*model.User, *Error) {
	userByEmail := (*s.repo).FindByEmail(email)
	if userByEmail != nil {
		msg := fmt.Sprintf("User with email %s already exists", email)
		return nil, &Error{Code: err.UserAlreadyExistsCode, Msg: msg}
	}
	token := ""
	for {
		generatedToken, e := s.tokenGenerator.getToken()
		if e != nil {
			return nil, &Error{Code: UnexpectedError, Msg: e.Error()}
		}
		if (*s.repo).FindByToken(token) == nil {
			token = generatedToken
			break
		}
	}
	newUser := &model.User{Email: email, Token: token}
	e := (*s.repo).Save(newUser)
	if e != nil {
		return nil, &Error{Code: UnexpectedError, Msg: e.Error()}
	}
	return newUser, nil
}

func NewUserService(userRepository *repository.UserRepository) UserService {
	return userServiceImpl{repo: userRepository, tokenGenerator: userTokenGenerator{}}
}

