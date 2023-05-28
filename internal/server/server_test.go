package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"ralts/internal/config"
	"ralts/internal/dependencies"
	testHelper "ralts/internal/testing"
	"strings"
	"testing"
	"time"
)

var cfg = config.NewConfig(true)

func FakeServer(deps *dependencies.Dependencies) *httptest.Server {
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
	return server
}

/**
TODOs:
- Group related tests together
- Unhappy path for /conn_count
- Unhappy path when trying to save most recent message
- Unhappy path when trying to get the message count when handling read request
*/

func TestServer_ServeChat(t *testing.T) {
	assert := assert.New(t)

	deps := dependencies.NewDependencies(cfg)
	defer deps.Disconnect()

	th := testHelper.TestHelper(cfg)
	defer th()

	server := FakeServer(deps)
	defer server.Close()

	wsUrl := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = time.Second

	ws, _, err := dialer.Dial(wsUrl, nil)
	assert.Nil(err)
	assert.NotNil(ws)

	t.Run("happy path", func(t *testing.T) {
		// Test_Write
		req := &Request{UserId: "e.sia", Message: "Lorem Ipsum"}
		err = ws.WriteJSON(req)
		assert.Nil(err)

		_, msg, err := ws.ReadMessage()
		assert.Nil(err)

		// Test_Read
		var resp Response
		_ = json.Unmarshal(msg, &resp)
		payload := resp.Payload
		procErr := resp.Error
		assert.Nil(procErr)
		assert.Equal(int64(1), payload.ChatId)
		assert.Equal(req.Message, payload.Message)
		assert.Equal(req.UserId, payload.Username)
		assert.NotNil(payload.CreatedAt)
	})

	t.Run("bad request", func(t *testing.T) {
		_ = ws.WriteJSON("bad request")
		_, msg, _ := ws.ReadMessage()
		var resp Response
		_ = json.Unmarshal(msg, &resp)
		assert.Nil(resp.Payload)
		assert.Equal(InternalServerError, resp.Error.Code)
		assert.Contains(resp.Error.Message, "cannot unmarshal")
	})
}

func TestServer_IncrDecrConnCount(t *testing.T) {
	assert := assert.New(t)

	deps := dependencies.NewDependencies(cfg)
	defer deps.Disconnect()

	th := testHelper.TestHelper(cfg)
	defer th()

	server := FakeServer(deps)
	defer server.Close()

	wsUrl := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = time.Second

	ws, _, _ := dialer.Dial(wsUrl, nil)
	// Leave a bit of buffer time for connection to init completely.
	time.Sleep(100 * time.Millisecond)

	r, err := deps.Cache.Get(CONN_COUNT_KEY)
	assert.Nil(err)
	assert.Equal("1", r)

	err = ws.Close()
	assert.Nil(err)
	// Leave a bit of buffer time for connection to init completely.
	time.Sleep(100 * time.Millisecond)

	r, err = deps.Cache.Get(CONN_COUNT_KEY)
	assert.Nil(err)
	assert.Equal("0", r)
}

func TestServer_MaxConnCount(t *testing.T) {
	assert := assert.New(t)

	deps := dependencies.NewDependencies(cfg)
	defer deps.Disconnect()

	th := testHelper.TestHelper(cfg)
	defer th()

	server := FakeServer(deps)
	defer server.Close()

	wsUrl := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = time.Second

	_, _, _ = dialer.Dial(wsUrl, nil)
	// Leave a bit of buffer time for connection to init completely.
	time.Sleep(100 * time.Millisecond)
	_, _, _ = dialer.Dial(wsUrl, nil)
	// Leave a bit of buffer time for connection to init completely.
	time.Sleep(100 * time.Millisecond)
	ws, _, _ := dialer.Dial(wsUrl, nil)
	_, _, err := ws.ReadMessage()
	assert.Equal("websocket: close 1013: max no. of client connections reached", err.Error())
}

func TestServer_MaxSentMsgPerDay(t *testing.T) {
	assert := assert.New(t)

	deps := dependencies.NewDependencies(cfg)
	defer deps.Disconnect()

	th := testHelper.TestHelper(cfg)
	defer th()

	server := FakeServer(deps)
	defer server.Close()

	wsUrl := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = time.Second

	ws, _, err := dialer.Dial(wsUrl, nil)
	assert.Nil(err)
	assert.NotNil(ws)

	var sentPayload *Request
	for i := 0; i < 5; i++ {
		sentPayload = &Request{UserId: "e.sia.2", Message: "Lorem Ipsum"}
		err = ws.WriteJSON(sentPayload)
		assert.Nil(err)
		_, _, err = ws.ReadMessage()
		assert.Nil(err)
	}
	sentPayload = &Request{UserId: "e.sia.2", Message: "This won't go through..."}
	_ = ws.WriteJSON(sentPayload)
	_, _, err = ws.ReadMessage()
	assert.Equal("websocket: close 1013: reached max no. of messages sent today", err.Error())
}

func TestServer_GetConnCount_EmptyResponse(t *testing.T) {
	assert := assert.New(t)

	deps := dependencies.NewDependencies(cfg)
	defer deps.Disconnect()

	th := testHelper.TestHelper(cfg)
	defer th()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/conn_count")

	s := NewServer(deps)

	assert.NoError(s.GetConnCount(c))
	assert.Equal(http.StatusOK, rec.Code)
	assert.Equal(`{"count":0}`, strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestServer_GetConnCount_NonEmptyResponse(t *testing.T) {
	assert := assert.New(t)

	deps := dependencies.NewDependencies(cfg)
	defer deps.Disconnect()

	th := testHelper.TestHelper(cfg)
	defer th()

	_ = deps.Cache.Incr(CONN_COUNT_KEY)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/conn_count")
	s := NewServer(deps)

	assert.NoError(s.GetConnCount(c))
	assert.Equal(http.StatusOK, rec.Code)
	assert.Equal(`{"count":1}`, strings.TrimSuffix(rec.Body.String(), "\n"))
}
