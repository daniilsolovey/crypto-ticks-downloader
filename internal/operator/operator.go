package operator

import (
	"strconv"

	"github.com/daniilsolovey/crypto-ticks-downloader/internal/config"
	"github.com/daniilsolovey/crypto-ticks-downloader/internal/database"
	"github.com/gorilla/websocket"
	"github.com/preichenberger/go-coinbasepro/v2"
	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"
)

type Operator struct {
	config           *config.Config
	websocket        *websocket.Conn
	database         *database.Database
	channels         map[string](chan *database.Ticks)
	distributionChan chan *database.Ticks
}

func NewOperator(
	config *config.Config,
	database *database.Database,
	websocket *websocket.Conn,
	tickerChannels map[string](chan *database.Ticks),
	distributionChan chan *database.Ticks,
) *Operator {
	return &Operator{
		config:           config,
		websocket:        websocket,
		database:         database,
		channels:         tickerChannels,
		distributionChan: distributionChan,
	}
}

func CreateTickersWithChannels(config *config.Config) map[string](chan *database.Ticks) {
	tickerChannels := make(map[string](chan *database.Ticks))
	for _, value := range config.Tickers {
		tickerChannels[value] = make(chan *database.Ticks)
	}

	return tickerChannels
}

func (operator *Operator) CreateTicksTable() error {
	err := operator.database.CreateTicksTable()
	if err != nil {
		return karma.Format(
			err,
			"unable to create 'ticks' table in the database",
		)
	}

	return nil
}

func (operator *Operator) ReceiveTicks() error {
	subscription := operator.createSubscribtion()
	err := operator.websocket.WriteJSON(subscription)
	if err != nil {
		return karma.Format(
			err,
			"unable to write json-encoded message to websocket connection",
		)
	}

	for {
		message := coinbasepro.Message{}
		err = operator.websocket.ReadJSON(&message)
		if err != nil {
			return karma.Format(
				err,
				"unable to read json-encoded message from websocket connection",
			)
		}

		ticker, err := prepareTicker(
			message.ProductID, message.BestAsk, message.BestBid,
		)
		if err != nil {
			return karma.Format(
				err,
				"unable to prepare ticker for further sending to channel",
			)
		}

		operator.distributionChan <- ticker
	}
}

func prepareTicker(productID, ask, bid string) (*database.Ticks, error) {
	if ask == "" || bid == "" || productID == "" {
		return &database.Ticks{}, nil
	}

	resultAsk, err := strconv.ParseFloat(ask, 64)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to convert ask price to float64, ask: %s",
			ask,
		)
	}

	resultBid, err := strconv.ParseFloat(bid, 64)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to convert bid to float64, bid: %s",
			bid,
		)
	}

	return &database.Ticks{
		Symbol:   productID,
		AskPrice: resultAsk,
		BidPrice: resultBid,
	}, nil
}

func (operator *Operator) createSubscribtion() coinbasepro.Message {
	var channels []coinbasepro.MessageChannel
	for key := range operator.channels {
		subscription := coinbasepro.MessageChannel{
			Name: "ticker",
			ProductIds: []string{
				key,
			},
		}

		log.Debugf(nil, "subscription created: %s", key)
		channels = append(channels, subscription)
	}

	result := coinbasepro.Message{
		Type:     "subscribe",
		Channels: channels,
	}
	return result
}
