package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"ralts/internal/config"
	"ralts/internal/db"
	testHelper "ralts/internal/testing"
	"strings"
	"testing"
)

type WsHandler struct {
	handler echo.HandlerFunc
}

var cfg = config.NewConfig(true)

func ServeHttp(w http.ResponseWriter, r *http.Request) {
	e := echo.New()
	c := e.NewContext(r, w)

	pool := NewPool()
	go pool.Start()

	s := &Server{}

	forever := make(chan struct{})
	_ = s.ServeChat(c, pool)
	<-forever
}

func TestServer_ServeChat(t *testing.T) {
	assert := assert.New(t)

	dbClient := db.NewRaltsDatabase(cfg)
	defer dbClient.Close()

	th := testHelper.TestHelper(cfg)
	defer th()

	server := httptest.NewServer(http.HandlerFunc(ServeHttp))
	defer server.Close()

	wsUrl := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	ws, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	assert.Nil(err)
	assert.NotNil(ws)

	sentPayload := &Payload{UserId: "e.sia", Message: "Lorem Ipsum"}
	err = ws.WriteJSON(sentPayload)
	assert.Nil(err)

	_, msg, err := ws.ReadMessage()
	assert.Nil(err)

	var receivedPayload Payload
	_ = json.Unmarshal(msg, &receivedPayload)
	assert.Equal(receivedPayload, *sentPayload)
}
