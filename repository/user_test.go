package repository

import (
	"temperature-backend/model"
	"temperature-backend/test_utils"
	"testing"
)

func Test_inMemoryRepository_Save_new(t *testing.T) {
	repo := NewUserRepository().(inMemoryUserRepository)
	user := model.User{Email: "test@test.ru", Token: "qwerty"}
	err := repo.Save(&user)
	test_utils.UnexpectedError(t, err)
	test_utils.AssertInt(t, 1, user.Id)
}

func Test_inMemoryRepository_Save_update(t *testing.T) {
	repo := NewUserRepository().(inMemoryUserRepository)
	user := model.User{Email: "test@test.ru", Token: "qwerty"}
	err := repo.Save(&user)
	test_utils.UnexpectedError(t, err)
	updatedEmail := "test2@test.ru"
	user.Email = updatedEmail
	err = repo.Save(&user)
	test_utils.UnexpectedError(t, err)
	test_utils.AssertString(t, updatedEmail, (*(repo.users))[0].Email)
}

// Test checks external immutability of repository. *model.User update shouldn't cause repository changes
func Test_inMemoryRepository_Save_immutability(t *testing.T) {
	repo := NewUserRepository().(inMemoryUserRepository)
	expectedToken := "qwerty"
	user := model.User{Email: "test@test.ru", Token: expectedToken}
	err := repo.Save(&user)
	test_utils.UnexpectedError(t, err)
	user.Token = "123"
	actualToken := (*(repo.users))[0].Token
	test_utils.AssertString(t, expectedToken, actualToken)
}

func Test_inMemoryUserRepository_FindByToken(t *testing.T) {
	repo := NewUserRepository().(inMemoryUserRepository)
	expectedToken := "qwerty"
	user := model.User{Email: "test@test.ru", Token: expectedToken}
	*repo.users = append(*repo.users, &user)
	findedUser := repo.FindByToken(expectedToken)
	test_utils.AssertNotNil(t, findedUser)
	test_utils.AssertString(t, expectedToken, findedUser.Token)
	findedUser.Token = "abc"
	actualToken := (*(repo.users))[0].Token
	test_utils.AssertString(t, expectedToken, actualToken)
}

func Test_inMemoryUserRepository_FIndByEmail(t *testing.T) {
	repo := NewUserRepository().(inMemoryUserRepository)
	expectedEmail := "test@test.ru"
	user := model.User{Email: expectedEmail, Token: "qwerty"}
	*repo.users = append(*repo.users, &user)
	findedUser := repo.FindByEmail(expectedEmail)
	test_utils.AssertNotNil(t, findedUser)
	test_utils.AssertString(t, expectedEmail, findedUser.Email)
	findedUser.Email = "abc"
	actualEmail := (*(repo.users))[0].Email
	test_utils.AssertString(t, expectedEmail, actualEmail)
}

func TestNewUserRepository(t *testing.T) {
	repo := NewUserRepository().(inMemoryUserRepository)
	test_utils.AssertInt(t, 0, repo.sequence.val)
	test_utils.AssertNotNil(t, repo.users)
}
