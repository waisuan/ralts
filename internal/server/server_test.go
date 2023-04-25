package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"ralts/internal/chat"
	"ralts/internal/config"
	"ralts/internal/db"
	testHelper "ralts/internal/testing"
	"strings"
	"testing"
	"time"
)

type TestServer struct {
	serve func(e echo.Context)
}

var cfg = config.NewConfig(true)

func (s *TestServer) ServeHttp(w http.ResponseWriter, r *http.Request) {
	e := echo.New()
	c := e.NewContext(r, w)

	forever := make(chan struct{})
	s.serve(c)
	<-forever
}

func TestServer_ServeChat(t *testing.T) {
	assert := assert.New(t)

	th := testHelper.TestHelper(cfg)
	defer th()

	pool := NewPool()
	go pool.Start()

	s := &Server{
		Config: cfg,
	}

	ts := &TestServer{
		serve: func(e echo.Context) {
			dbClient := db.NewRaltsDatabase(cfg)
			defer dbClient.Close()

			s.Handlers = &Handlers{
				ChatHandler: chat.NewChat(dbClient),
			}
			_ = s.ServeChat(e, pool)
		},
	}

	server := httptest.NewServer(http.HandlerFunc(ts.ServeHttp))
	defer server.Close()

	wsUrl := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = time.Second

	ws, _, err := dialer.Dial(wsUrl, nil)
	assert.Nil(err)
	assert.NotNil(ws)

	sentPayload := &Payload{UserId: "e.sia", Message: "Lorem Ipsum"}
	err = ws.WriteJSON(sentPayload)
	assert.Nil(err)

	_, msg, err := ws.ReadMessage()
	assert.Nil(err)

	var receivedPayload chat.Message
	_ = json.Unmarshal(msg, &receivedPayload)
	assert.Equal(int64(1), receivedPayload.ChatId)
	assert.Equal(sentPayload.Message, receivedPayload.Message)
	assert.Equal(sentPayload.UserId, receivedPayload.Username)
	assert.NotNil(receivedPayload.CreatedAt)

	// Test_MaxConnCount
	_, _, _ = dialer.Dial(wsUrl, nil)
	ws, _, _ = dialer.Dial(wsUrl, nil)
	_, _, err = ws.ReadMessage()
	assert.Equal("websocket: close 1013: max no. of client connections reached", err.Error())
}
