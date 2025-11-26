package database

import "errors"

var (
	ErrUserDoesNotExist            = errors.New("user does not exist")
	ErrConversationtDoesNotExist   = errors.New("conversation does not exist")
	ErrMessageDoesNotExist         = errors.New("message does not exist")
	ErrCommentDoesNotExist         = errors.New("comment does not exist")
	ErrUnauthorizedToDeleteMessage = errors.New("unauthorized To Delete Message")
	ErrGroupDoesNotExist           = errors.New("group does not exist")
)
