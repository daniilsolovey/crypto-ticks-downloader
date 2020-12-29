package operator

import (
	"strconv"

	"github.com/daniilsolovey/crypto-ticks-downloader/internal/config"
	"github.com/daniilsolovey/crypto-ticks-downloader/internal/database"
	"github.com/gorilla/websocket"
	"github.com/preichenberger/go-coinbasepro/v2"
	"github.com/reconquest/karma-go"
)

type Operator struct {
	config    *config.Config
	websocket *websocket.Conn
	database  *database.Database
	channels  *Channels
}

type Channels struct {
	DistributionChan chan *database.Ticks
	BTCUSDChan       chan *database.Ticks
	BTCEURChan       chan *database.Ticks
	ETHBTCChan       chan *database.Ticks
}

func NewOperator(
	config *config.Config,
	database *database.Database,
	websocket *websocket.Conn,
	channels *Channels,
) *Operator {
	return &Operator{
		config:    config,
		websocket: websocket,
		database:  database,
		channels:  channels,
	}
}

func CreateChannelsForTickers() *Channels {
	return &Channels{
		DistributionChan: make(chan *database.Ticks),
		BTCUSDChan:       make(chan *database.Ticks),
		BTCEURChan:       make(chan *database.Ticks),
		ETHBTCChan:       make(chan *database.Ticks),
	}
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
	subscription := createSubscribtion()
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

		operator.channels.DistributionChan <- ticker
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

func createSubscribtion() coinbasepro.Message {
	subscription := coinbasepro.Message{
		Type: "subscribe",
		Channels: []coinbasepro.MessageChannel{
			{
				Name: "ticker",
				ProductIds: []string{
					"BTC-USD",
				},
			},

			{
				Name: "ticker",
				ProductIds: []string{
					"BTC-EUR",
				},
			},

			{
				Name: "ticker",
				ProductIds: []string{
					"ETH-BTC",
				},
			},
		},
	}

	return subscription
}
