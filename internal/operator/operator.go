package operator

import (
	"time"

	"github.com/daniilsolovey/crypto-ticks-downloader/internal/config"
	"github.com/gorilla/websocket"
	"github.com/preichenberger/go-coinbasepro/v2"
	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"
)

type Operator struct {
	config    *config.Config
	websocket *websocket.Conn
	// database    *database.Database

}

func NewOperator(
	config *config.Config,
	// database *database.Database,
	websocket *websocket.Conn,
) *Operator {
	return &Operator{
		config:    config,
		websocket: websocket,
		// database:         database,
	}
}

func (operator *Operator) GetPrices() error {
	subscription := createSubscribtion()
	err := operator.websocket.WriteJSON(subscription)
	if err != nil {
		return karma.Format(
			err,
			"unable to write json as message",
		)
	}

	for true {
		message := coinbasepro.Message{}
		err = operator.websocket.ReadJSON(&message)
		if err != nil {
			return karma.Format(
				err,
				"unable to write json as message",
			)
		}

		log.Infof(
			karma.
				Describe("best ask: ", message.BestAsk).
				Describe("best bid: ", message.BestBid), "currency: %s", message.ProductID)

		time.Sleep(1 * time.Second)
	}

	return nil
}

func createSubscribtion() coinbasepro.Message {
	subscription := coinbasepro.Message{
		Type: "subscribe",
		Channels: []coinbasepro.MessageChannel{
			coinbasepro.MessageChannel{
				Name: "ticker",
				ProductIds: []string{
					"BTC-USD",
				},
			},

			coinbasepro.MessageChannel{
				Name: "ticker",
				ProductIds: []string{
					"BTC-EUR",
				},
			},

			coinbasepro.MessageChannel{
				Name: "ticker",
				ProductIds: []string{
					"ETH-BTC",
				},
			},
		},
	}

	return subscription
}