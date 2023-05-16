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

func TestChat_LoadAllMessages_MessagesAvailable(t *testing.T) {
	assert := assert.New(t)

	deps := dependencies.NewDependencies(cfg)
	defer deps.Disconnect()

	th := testHelper.TestHelper(cfg)
	defer th()

	chat := NewChat(deps)

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
}

func TestChat_LoadAllMessages_NoMessagesAvailable(t *testing.T) {
	assert := assert.New(t)

	deps := dependencies.NewDependencies(cfg)
	defer deps.Disconnect()

	th := testHelper.TestHelper(cfg)
	defer th()

	chat := NewChat(deps)

	msgs, err := chat.LoadAllMessages()
	assert.Nil(err)
	assert.Empty(msgs)
}

func TestChat_SaveMessage_Successful(t *testing.T) {
	assert := assert.New(t)

	deps := dependencies.NewDependencies(cfg)
	defer deps.Disconnect()

	th := testHelper.TestHelper(cfg)
	defer th()

	chat := NewChat(deps)

	msg, err := chat.SaveMessage("e.sia", "Lorem Ipsum", fakeNow)
	assert.Nil(err)
	assert.Equal(int64(1), msg.ChatId)
}

func TestChat_GetMessageCount_NoMessages(t *testing.T) {
	assert := assert.New(t)

	deps := dependencies.NewDependencies(cfg)
	defer deps.Disconnect()

	th := testHelper.TestHelper(cfg)
	defer th()

	chat := NewChat(deps)

	count, err := chat.GetMessageCount("e.sia", time.Now)
	assert.Nil(err)
	assert.Equal(0, count)
}

func TestChat_GetMessageCount_HasMessages(t *testing.T) {
	assert := assert.New(t)

	deps := dependencies.NewDependencies(cfg)
	defer deps.Disconnect()

	th := testHelper.TestHelper(cfg)
	defer th()

	chat := NewChat(deps)

	_, _ = chat.SaveMessage("e.sia", "testing", time.Now)

	count, err := chat.GetMessageCount("e.sia", time.Now)
	assert.Nil(err)
	assert.Equal(1, count)
}

func TestChat_GetMessageCount_DiffDay(t *testing.T) {
	assert := assert.New(t)

	deps := dependencies.NewDependencies(cfg)
	defer deps.Disconnect()

	th := testHelper.TestHelper(cfg)
	defer th()

	chat := NewChat(deps)

	_, _ = chat.SaveMessage("e.sia", "testing", fakeNow)

	count, err := chat.GetMessageCount("e.sia", time.Now)
	assert.Nil(err)
	assert.Equal(0, count)
}
