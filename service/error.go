package service

import "fmt"

const UnexpectedError = 500
const EntityNotFound = 404
const EntityAlreadyExists = 400
const ParseError = 422
const AuthError = 401

type Error struct {
	Code int
	Msg  string
}

func (e Error) Error() string {
	return fmt.Sprintf("service error %d: %s", e.Code, e.Msg)
}
