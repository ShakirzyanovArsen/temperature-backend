package service

import (
	"errors"
	"fmt"
	"reflect"
	"temperature-backend/model"
	"temperature-backend/repository"
	"testing"
)

type userRepoMock struct {
	saveError       error
	findEmailResult *model.User
	findTokenResult *model.User
}

func (r userRepoMock) Save(user *model.User) error {
	if r.saveError == nil {
		user.Id = 1
	}
	return r.saveError
}

func (r userRepoMock) FindByToken(token string) *model.User {
	return r.findTokenResult
}

func (r userRepoMock) FindByEmail(email string) *model.User {
	return r.findEmailResult
}

func createUserRepoMock(saveError error, findEmailResult *model.User, findTokenResult *model.User) repository.UserRepository {
	result := &(userRepoMock{saveError: saveError, findEmailResult: findEmailResult, findTokenResult: findTokenResult})
	return result
}

type tokenGeneratorMock struct {
	tokens []string
	next   int
}

func (g tokenGeneratorMock) getToken() (string, error) {
	result := g.tokens[g.next]
	if len(g.tokens) <= g.next {
		return "", errors.New("error in token generation")
	}
	g.next++
	return result, nil
}

func creatMockedTokenGenerator(tokens []string) tokenGeneratorMock {
	return tokenGeneratorMock{tokens: tokens, next: 0}
}

func Test_userServiceImpl_Register(t *testing.T) {
	type fields struct {
		repo           *repository.UserRepository
		tokenGenerator tokenGenerator
	}
	type args struct {
		email string
	}
	successSaveRepo := createUserRepoMock(nil, nil, nil)
	existingEmail := "user_with@email.exists"
	userExistsRepoMock := createUserRepoMock(nil, &model.User{}, nil)
	saveErrorText := "Error on "
	saveWithErrorRepoMock := createUserRepoMock(errors.New(saveErrorText), nil, nil)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr *Error
	}{
		{
			name:    "Success save",
			fields:  fields{repo: &successSaveRepo, tokenGenerator: creatMockedTokenGenerator([]string{"123"})},
			args:    args{email: "test@test.ru"},
			want:    &(model.User{Id: 1, Email: "test@test.ru", Token: "123"}),
			wantErr: nil,
		},
		{
			name:    "User with email exists",
			fields:  fields{repo: &userExistsRepoMock, tokenGenerator: creatMockedTokenGenerator([]string{"123"})},
			args:    args{email: existingEmail},
			want:    nil,
			wantErr: &Error{Code: EntityAlreadyExists, Msg: fmt.Sprintf("User with email %s already exists", existingEmail)},
		},
		{
			name:    "Error on save",
			fields:  fields{repo: &saveWithErrorRepoMock, tokenGenerator: creatMockedTokenGenerator([]string{"123"})},
			args:    args{email: "test@test.ru"},
			want:    nil,
			wantErr: &Error{Code: UnexpectedError, Msg: saveErrorText},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := userServiceImpl{
				repo:           tt.fields.repo,
				tokenGenerator: tt.fields.tokenGenerator,
			}
			got, gotErr := s.Register(tt.args.email)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userServiceImpl.Register() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("userServiceImpl.Register() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
