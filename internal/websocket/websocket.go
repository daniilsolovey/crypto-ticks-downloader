package websocket

import (
	"github.com/daniilsolovey/crypto-ticks-downloader/internal/config"
	ws "github.com/gorilla/websocket"
	"github.com/reconquest/karma-go"
)

type Client struct {
	config   *config.Config
	wsDialer ws.Dialer
}

func NewWebSocketConnection(
	url string,
) (*ws.Conn, error) {
	var wsDialer ws.Dialer
	websocket, _, err := wsDialer.Dial(url, nil)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to connect to the websocket url: %s",
			url,
		)
	}

	return websocket, nil
}
