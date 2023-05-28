package server

import "ralts/internal/chat"

const (
	InternalServerError = "INTERNAL_SERVER_ERROR"
)

type Response struct {
	Payload *chat.Message
	Error   *Error
}

type Request struct {
	UserId  string
	Message string
}
type Error struct {
	Code    string
	Message string
}
