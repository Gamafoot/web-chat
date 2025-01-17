package domain

import (
	"errors"
)

var (
	ErrLoginIsExists          = errors.New("Login is exists")
	ErrInvalidLoginOrPassword = errors.New("Invalid login or password")
	ErrBadRequest             = errors.New("Invalid request data")
	ErrRecordNotFound         = errors.New("Record no found")
	ErrRoomFull               = errors.New("The room is full")
	ErrFailToJoinRoom         = errors.New("Fail to join room")
)
