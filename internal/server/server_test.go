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
	"ralts/internal/dependencies"
	testHelper "ralts/internal/testing"
	"strings"
	"testing"
	"time"
)

var cfg = config.NewConfig(true)

func TestServer_ServeChat(t *testing.T) {
	assert := assert.New(t)

	deps := dependencies.NewDependencies(cfg)
	defer deps.Disconnect()

	th := testHelper.TestHelper(cfg)
	defer th()

	callbacks := NewCallbacks(deps)
	go callbacks.Listen()

	pool := NewPool(callbacks)
	go pool.Start()

	s := &Server{
		Deps: deps,
	}

	serveHttp := func(w http.ResponseWriter, r *http.Request) {
		e := echo.New()
		c := e.NewContext(r, w)

		forever := make(chan struct{})
		_ = s.ServeChat(c, pool)
		<-forever
	}

	server := httptest.NewServer(http.HandlerFunc(serveHttp))
	defer server.Close()

	wsUrl := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = time.Second

	ws, _, err := dialer.Dial(wsUrl, nil)
	assert.Nil(err)
	assert.NotNil(ws)

	// Test_Write
	sentPayload := &Payload{UserId: "e.sia", Message: "Lorem Ipsum"}
	err = ws.WriteJSON(sentPayload)
	assert.Nil(err)

	_, msg, err := ws.ReadMessage()
	assert.Nil(err)

	// Test_Read
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

	// Test_CachedEntries
	r, err := deps.Cache.Get(CONN_COUNT_KEY)
	assert.Nil(err)
	assert.Equal(r, "1")
}
