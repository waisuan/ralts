package chat

import (
	"github.com/stretchr/testify/assert"
	"ralts/internal/config"
	"ralts/internal/db"
	testHelper "ralts/internal/testing"
	"testing"
	"time"
)

var cfg = config.NewConfig(true)

func initTestHelper() func() {
	testHelper.ClearDB(cfg)
	testHelper.InitDB(cfg)

	return func() {
		testHelper.ClearDB(cfg)
	}
}

func now() time.Time {
	return time.Date(2022, 11, 17, 20, 34, 58, 651387000, time.UTC)
}

func TestChat_LoadAllMessages_MessagesAvailable(t *testing.T) {
	assert := assert.New(t)

	dbClient := db.NewRaltsDatabase(cfg)
	defer dbClient.Close()

	tH := initTestHelper()
	defer tH()

	chat := NewChat(dbClient)

	_, err := chat.SaveMessage("e.sia", "Lorem Ipsum", now)
	assert.Nil(err)

	msgs, err := chat.LoadAllMessages()
	assert.Nil(err)
	assert.Len(msgs, 1)

	m := msgs[0]
	assert.Equal(m.ChatId, int64(1))
	assert.Equal(m.Username, "e.sia")
	assert.Equal(m.Message, "Lorem Ipsum")
	assert.Equal(m.CreatedAt, now())
}

func TestChat_LoadAllMessages_NoMessagesAvailable(t *testing.T) {
	assert := assert.New(t)

	dbClient := db.NewRaltsDatabase(cfg)
	defer dbClient.Close()

	tH := initTestHelper()
	defer tH()

	chat := NewChat(dbClient)

	msgs, err := chat.LoadAllMessages()
	assert.Nil(err)
	assert.Empty(msgs)
}

func TestChat_SaveMessage_Successful(t *testing.T) {
	assert := assert.New(t)

	dbClient := db.NewRaltsDatabase(cfg)
	defer dbClient.Close()

	tH := initTestHelper()
	defer tH()

	chat := NewChat(dbClient)

	msg, err := chat.SaveMessage("e.sia", "Lorem Ipsum", now)
	assert.Nil(err)
	assert.Equal(msg.ChatId, int64(1))
}
