package repository

import (
	"fmt"
	"temperature-backend/model"
	"testing"
)

func Test_inMemoryRepository_Save_new(t *testing.T) {
	repo := NewUserRepository().(inMemoryUserRepository)
	user := model.User{Email: "test@test.ru", Token: "qwerty"}
	err := repo.Save(&user)
	if err != nil {
		t.Error(fmt.Sprintf("unexpected err: %s", err))
	}
	if user.Id != 1 {
		t.Error(fmt.Sprintf("expected user.id:1, actual: %d", user.Id))
	}
}

func Test_inMemoryRepository_Save_update(t *testing.T) {
	repo := NewUserRepository().(inMemoryUserRepository)
	user := model.User{Email: "test@test.ru", Token: "qwerty"}
	err := repo.Save(&user)
	if err != nil {
		t.Error(fmt.Sprintf("unexpected err: %s", err))
	}
	updatedEmail := "test2@test.ru"
	user.Email = updatedEmail
	err = repo.Save(&user)
	if err != nil {
		t.Error(fmt.Sprintf("unexpected err: %s", err))
	}
	if (*(repo.users))[0].Email != updatedEmail {
		t.Error(fmt.Sprintf("expected email:%s, actual email: %s", updatedEmail, (*(repo.users))[0].Email))
	}
}

// Test checks external immutability of repository. *model.User update shouldn't cause repository changes
func Test_inMemoryRepository_Save_immutability(t *testing.T) {
	repo := NewUserRepository().(inMemoryUserRepository)
	expectedToken := "qwerty"
	user := model.User{Email: "test@test.ru", Token: expectedToken}
	err := repo.Save(&user)
	if err != nil {
		t.Error(fmt.Sprintf("unexpected err: %s", err))
	}
	user.Token = "123"
	actualToken := (*(repo.users))[0].Token
	if actualToken != expectedToken {
		t.Error(fmt.Sprintf("expected token: %s, actualToken: %s. Save() method causes mutability of repository", expectedToken, actualToken))
	}
}

func Test_inMemoryUserRepository_FindByToken(t *testing.T) {
	repo := NewUserRepository().(inMemoryUserRepository)
	expectedToken := "qwerty"
	user := model.User{Email: "test@test.ru", Token: expectedToken}
	*repo.users = append(*repo.users, &user)
	findedUser := repo.FindByToken(expectedToken)
	if findedUser == nil {
		t.Error("expected user model, actual: nil")
	}
	if expectedToken != findedUser.Token {
		t.Error(fmt.Sprintf("expected user token: %s, actual user token: %s", expectedToken, findedUser.Token))
	}
	findedUser.Token = "abc"
	actualToken := (*(repo.users))[0].Token
	if expectedToken != actualToken {
		t.Error(fmt.Sprintf("expected token: %s, actualToken: %s. FindByToken() causes mutability of repository", expectedToken, actualToken))
	}
}

func Test_inMemoryUserRepository_FIndByEmail(t *testing.T) {
	repo := NewUserRepository().(inMemoryUserRepository)
	expectedEmail := "test@test.ru"
	user := model.User{Email: expectedEmail, Token: "qwerty"}
	*repo.users = append(*repo.users, &user)
	findedUser := repo.FindByEmail(expectedEmail)
	if findedUser == nil {
		t.Error("expected user model, actual: nil")
	}
	if expectedEmail != findedUser.Email {
		t.Error(fmt.Sprintf("expected user email: %s, actual user email: %s", expectedEmail, findedUser.Email))
	}
	findedUser.Email = "abc"
	actualEmail := (*(repo.users))[0].Email
	if expectedEmail != actualEmail {
		t.Error(fmt.Sprintf("expected email: %s, actual email: %s. FindByEmail() causes mutability of repository", expectedEmail, actualEmail))
	}
}

func TestNewUserRepository(t *testing.T) {
	repo := NewUserRepository().(inMemoryUserRepository)
	if repo.sequence != 0 {
		t.Error(fmt.Sprintf("sequence should be zero, actual: %d", repo.sequence))
	}
	if repo.users == nil {
		t.Error("users slice should be initialized")
	}
}
