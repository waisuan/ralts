package chat

import (
	"github.com/stretchr/testify/assert"
	"ralts/internal/config"
	"ralts/internal/dependencies"
	testHelper "ralts/internal/testing"
	"testing"
	"time"
)

var cfg = config.NewConfig(true)

func fakeNow() time.Time {
	return time.Date(2022, 11, 17, 20, 34, 58, 651387000, time.UTC)
}

func TestChat_LoadAllMessages(t *testing.T) {
	assert := assert.New(t)

	deps := dependencies.NewDependencies(cfg)
	defer deps.Disconnect()

	th := testHelper.TestHelper(cfg)
	defer th()

	chat := NewChat(deps)

	t.Run("no messages available", func(t *testing.T) {
		msgs, err := chat.LoadAllMessages()
		assert.Nil(err)
		assert.Empty(msgs)
	})

	t.Run("messages available", func(t *testing.T) {
		_, err := chat.SaveMessage("e.sia", "Lorem Ipsum", fakeNow)
		assert.Nil(err)

		msgs, err := chat.LoadAllMessages()
		assert.Nil(err)
		assert.Len(msgs, 1)

		m := msgs[0]
		assert.Equal(int64(1), m.ChatId)
		assert.Equal("e.sia", m.Username)
		assert.Equal("Lorem Ipsum", m.Message)
		assert.Equal(fakeNow(), m.CreatedAt)
	})
}

func TestChat_SaveMessage(t *testing.T) {
	assert := assert.New(t)

	deps := dependencies.NewDependencies(cfg)
	defer deps.Disconnect()

	th := testHelper.TestHelper(cfg)
	defer th()

	chat := NewChat(deps)

	t.Run("happy path", func(t *testing.T) {
		msg, err := chat.SaveMessage("e.sia", "Lorem Ipsum", fakeNow)
		assert.Nil(err)
		assert.Equal(int64(1), msg.ChatId)
	})
}

func TestChat_GetMessageCount(t *testing.T) {
	assert := assert.New(t)

	deps := dependencies.NewDependencies(cfg)
	defer deps.Disconnect()

	th := testHelper.TestHelper(cfg)
	defer th()

	chat := NewChat(deps)

	t.Run("no messages", func(t *testing.T) {
		count, err := chat.GetMessageCount("e.sia", time.Now)
		assert.Nil(err)
		assert.Equal(0, count)
	})

	t.Run("has messages", func(t *testing.T) {
		_, _ = chat.SaveMessage("e.sia", "testing", time.Now)

		count, err := chat.GetMessageCount("e.sia", time.Now)
		assert.Nil(err)
		assert.Equal(1, count)
	})

	t.Run("diff day", func(t *testing.T) {
		_, _ = chat.SaveMessage("e.sia.2", "testing", fakeNow)

		count, err := chat.GetMessageCount("e.sia.2", time.Now)
		assert.Nil(err)
		assert.Equal(0, count)
	})
}
