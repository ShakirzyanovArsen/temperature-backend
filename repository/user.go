package repository

import (
	"errors"
	"sync"
	"temperature-backend/model"
)

type UserRepository interface {
	Save(user *model.User) error
	FindByToken(token string) *model.User
	FindByEmail(email string) *model.User
}

type inMemoryUserRepository struct {
	mux sync.Mutex
	users      *[]*model.User
	sequence   int
}

func (repo inMemoryUserRepository) findById(id int) *model.User{
	for _, u := range *repo.users {
		if u.Id == id {
			return u
		}
	}
	return nil
}

func (repo inMemoryUserRepository) Save(user *model.User) error {
	repo.mux.Lock()
	defer repo.mux.Unlock()
	if user.Id == 0 {
		repo.sequence++
		user.Id = repo.sequence
		newUser := *user
		*repo.users = append(*repo.users, &newUser)
		return nil
	} else {
		userToUpdate := repo.findById(user.Id)
		if userToUpdate == nil {
			return errors.New(`can't find user to update`)
		}
		*userToUpdate = *user
		return nil
	}
}

func (repo inMemoryUserRepository) FindByToken(token string) *model.User {
	for _, u := range *repo.users {
		if u.Token == token {
			userToReturn := *u
			return &userToReturn
		}
	}
	return nil
}

func (repo inMemoryUserRepository) FindByEmail(email string) *model.User {
	for _, u := range *repo.users {
		if u.Email == email {
			userToReturn := *u
			return &userToReturn
		}
	}
	return nil
}

func NewUserRepository() UserRepository {
	users := new([]*model.User)
	repo := inMemoryUserRepository{users: users, sequence: 0}
	return repo
}
