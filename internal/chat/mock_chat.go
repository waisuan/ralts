package chat

import (
	"errors"
	"math/rand"
	"time"
)

type MockChatHandler struct {
	Config *MockChatHandlerConfig
}

type MockChatHandlerConfig struct {
	UnableToSave        bool
	UnableToGetMsgCount bool
}

func NewMockChatHandler(cfg *MockChatHandlerConfig) *MockChatHandler {
	return &MockChatHandler{
		Config: cfg,
	}
}

func (c *MockChatHandler) LoadAllMessages() (Messages, error) {
	return nil, nil
}

func (c *MockChatHandler) SaveMessage(username string, text string, now func() time.Time) (*Message, error) {
	if c.Config.UnableToSave {
		return nil, errors.New("unable to save message")
	}

	return &Message{
		ChatId:    int64(rand.Int()),
		Username:  username,
		Message:   text,
		CreatedAt: now(),
	}, nil
}

func (c *MockChatHandler) GetMessageCount(username string, today func() time.Time) (int, error) {
	if c.Config.UnableToGetMsgCount {
		return 0, errors.New("unable to get message count")
	}

	return 0, nil
}
