package websocket

import (
	ws "github.com/gorilla/websocket"
	"github.com/reconquest/karma-go"
)

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
