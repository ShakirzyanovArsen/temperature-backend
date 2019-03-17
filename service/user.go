package service

import (
	"fmt"
	"temperature-backend/model"
	"temperature-backend/repository"
)

type UserService interface {
	Register(email string) (*model.User, *Error)
}

type userServiceImpl struct {
	repo           *repository.UserRepository
	tokenGenerator tokenGenerator
}

func (s userServiceImpl) Register(email string) (*model.User, *Error) {
	userByEmail := (*s.repo).FindByEmail(email)
	if userByEmail != nil {
		msg := fmt.Sprintf("User with email %s already exists", email)
		return nil, &Error{Code: EntityAlreadyExists, Msg: msg}
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
	return userServiceImpl{repo: userRepository, tokenGenerator: tokenGeneratorImpl{}}
}
