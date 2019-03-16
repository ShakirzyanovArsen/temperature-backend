package service

import "fmt"

const UnexpectedError = 500

type Error struct {
	Code int
	Msg  string
}

func (e Error) Error() string {
	return fmt.Sprintf("service error %d: %s", e.Code, e.Msg)
}
