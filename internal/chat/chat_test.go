package chat

import (
	"errors"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"ralts/internal/config"
	"ralts/internal/dependencies"
	testHelper "ralts/internal/testing"
	"regexp"
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

	t.Run("unable to execute query", func(t *testing.T) {
		mockPg, _ := pgxmock.NewPool()

		mockPg.ExpectQuery(regexp.QuoteMeta(`select * from chat`)).WillReturnError(errors.New("something went wrong"))

		deps = &dependencies.Dependencies{
			Cfg:     cfg,
			Storage: mockPg,
		}
		defer deps.Disconnect()

		chat = NewChat(deps)
		_, err := chat.LoadAllMessages()
		assert.Contains(err.Error(), "something went wrong")
	})

	t.Run("unable to scan query results", func(t *testing.T) {
		mockPg, _ := pgxmock.NewPool()

		rows := mockPg.NewRows([]string{"name", "dob"}).AddRow("evan", "12/11/1991")
		mockPg.ExpectQuery(regexp.QuoteMeta(`select * from chat`)).WillReturnRows(rows)

		deps = &dependencies.Dependencies{
			Cfg:     cfg,
			Storage: mockPg,
		}
		defer deps.Disconnect()

		chat = NewChat(deps)
		_, err := chat.LoadAllMessages()
		assert.Contains(err.Error(), "Incorrect argument number")
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

	t.Run("unsuccessful", func(t *testing.T) {
		mockPg, _ := pgxmock.NewPool()

		rows := mockPg.NewRows([]string{})
		mockPg.ExpectQuery("insert into chat").WithArgs("e.sia", "Lorem Ipsum", fakeNow()).WillReturnRows(rows)

		deps = &dependencies.Dependencies{
			Cfg:     cfg,
			Storage: mockPg,
		}
		defer deps.Disconnect()

		chat = NewChat(deps)
		_, err := chat.SaveMessage("e.sia", "Lorem Ipsum", fakeNow)
		assert.Contains(err.Error(), "Incorrect argument number")
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

	t.Run("unable to execute query", func(t *testing.T) {
		mockPg, _ := pgxmock.NewPool()

		mockPg.ExpectQuery(regexp.QuoteMeta(`select count(*) from chat`)).WithArgs("e.sia", time.Now().Format("2006-01-02")).WillReturnError(errors.New("something went wrong"))

		deps = &dependencies.Dependencies{
			Cfg:     cfg,
			Storage: mockPg,
		}
		defer deps.Disconnect()

		chat = NewChat(deps)
		_, err := chat.GetMessageCount("e.sia", time.Now)
		assert.Contains(err.Error(), "something went wrong")
	})

	t.Run("unable to scan query results", func(t *testing.T) {
		mockPg, _ := pgxmock.NewPool()

		rows := mockPg.NewRows([]string{"name", "dob"}).AddRow("evan", "12/11/1991")
		mockPg.ExpectQuery(regexp.QuoteMeta(`select count(*) from chat`)).WithArgs("e.sia", time.Now().Format("2006-01-02")).WillReturnRows(rows)

		deps = &dependencies.Dependencies{
			Cfg:     cfg,
			Storage: mockPg,
		}
		defer deps.Disconnect()

		chat = NewChat(deps)
		_, err := chat.GetMessageCount("e.sia", time.Now)
		assert.Contains(err.Error(), "Incorrect argument number")
	})
}
