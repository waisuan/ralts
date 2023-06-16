package chat

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"ralts/internal/dependencies"
	"time"
)

type ChatHandler interface {
	LoadAllMessages() (Messages, error)
	SaveMessage(username string, text string, now func() time.Time) (*Message, error)
	GetMessageCount(username string, today func() time.Time) (int, error)
}

type Message struct {
	ChatId    int64     `json:"chatId"`
	Username  string    `json:"username"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}

type Messages []Message

func (m *Message) ToString() string {
	return fmt.Sprintf("[%s] %s: %s",
		m.CreatedAt.Format("2006-01-02 15:01:05"),
		m.Username,
		m.Message,
	)
}

type Chat struct {
	deps *dependencies.Dependencies
}

func NewChat(deps *dependencies.Dependencies) *Chat {
	return &Chat{
		deps: deps,
	}
}

func (c *Chat) LoadAllMessages() (Messages, error) {
	rows, err := c.deps.Storage.Query(context.Background(), "select * from chat order by created_at limit 50;")
	if err != nil {
		log.Error(fmt.Sprintf("Unable to execute query -> %s", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var messages Messages
	for rows.Next() {
		var m Message
		err := rows.Scan(&m.ChatId, &m.Username, &m.Message, &m.CreatedAt)
		if err != nil {
			log.Error(fmt.Sprintf("Unable load query results -> %s", err.Error()))
			return nil, err
		}
		messages = append(messages, m)
	}
	if err := rows.Err(); err != nil {
		log.Error(fmt.Sprintf("Something went wrong -> %s", err.Error()))
		return nil, err
	}

	return messages, nil
}

func (c *Chat) SaveMessage(username string, text string, now func() time.Time) (*Message, error) {
	var m Message
	err := c.deps.Storage.QueryRow(context.Background(), `
    	insert into chat (username, message, created_at)
    	VALUES ($1, $2, $3)
    	RETURNING chat_id, username, message, created_at 
    `, username, text, now()).Scan(&m.ChatId, &m.Username, &m.Message, &m.CreatedAt)
	if err != nil {
		log.Error(fmt.Sprintf("Unable to save chat message -> %s", err.Error()))
		return nil, err
	}

	return &m, nil
}

func (c *Chat) GetMessageCount(username string, today func() time.Time) (int, error) {
	rows, err := c.deps.Storage.Query(context.Background(), `
		select count(*) from chat where username = $1 and created_at::date = $2;
	`, username, today().Format("2006-01-02"))
	if err != nil {
		log.Error(fmt.Sprintf("Unable to execute query -> %s", err.Error()))
		return 0, err
	}
	defer rows.Close()

	rows.Next()

	var count int
	err = rows.Scan(&count)
	if err != nil {
		log.Error(fmt.Sprintf("Unable load query results -> %s", err.Error()))
		return 0, err
	}

	if err := rows.Err(); err != nil {
		log.Error(fmt.Sprintf("Something went wrong -> %s", err.Error()))
		return 0, err
	}

	return count, nil
}
