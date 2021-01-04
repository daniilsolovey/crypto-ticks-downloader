package websocket

import (
	"github.com/gorilla/websocket"
	ws "github.com/gorilla/websocket"
	"github.com/preichenberger/go-coinbasepro/v2"
	"github.com/reconquest/karma-go"
)

type Websocket interface {
	WriteJSON(coinbasepro.Message) error
	ReadJSON() (*coinbasepro.Message, error)
}

type Client struct {
	websocket *websocket.Conn
}

func NewClient(websocket *websocket.Conn) *Client {
	return &Client{
		websocket: websocket,
	}
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

func (client *Client) WriteJSON(subscription coinbasepro.Message) error {
	err := client.websocket.WriteJSON(subscription)
	if err != nil {
		return karma.Format(
			err,
			"unable to write json-encoded message to websocket connection",
		)
	}

	return nil
}

func (client *Client) ReadJSON() (*coinbasepro.Message, error) {
	message := coinbasepro.Message{}
	err := client.websocket.ReadJSON(&message)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to read json-encoded message from websocket connection",
		)
	}

	return &message, nil
}
