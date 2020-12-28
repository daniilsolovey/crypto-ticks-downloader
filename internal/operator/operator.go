package operator

import (
	"fmt"
	"strconv"
	"time"

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
	distributionChan chan *database.Ticker
}

func NewOperator(
	config *config.Config,
	database *database.Database,
	websocket *websocket.Conn,
	channel chan *database.Ticker,
) *Operator {
	return &Operator{
		config:           config,
		websocket:        websocket,
		database:         database,
		distributionChan: channel,
	}
}

func (operator *Operator) DistributeTickers() {
	for {
		select {
		case data, ok := <-operator.distributionChan:
			if ok {
				log.Warning("x ", data)
				fmt.Printf("Value %v was read.\n", data)
				switch data.Symbol {
				case "BTC-USD":
					go operator.WriteBTCUSD(data)
				case "ETH-BTC":
					go operator.WriteETHBTC(data)
				case "BTC-EUR":
					go operator.WriteBTCEUR(data)
				}

				continue
			}
		}
	}
}

func (operator *Operator) WriteBTCUSD(ticker *database.Ticker) {
	fmt.Println("ticker: ", ticker)
	err := operator.database.Write(*ticker)
	if err != nil {
		log.Error(err)
	}
}

func (operator *Operator) WriteETHBTC(ticker *database.Ticker) {
	fmt.Println("ticker: ", ticker)
	err := operator.database.Write(*ticker)
	if err != nil {
		log.Error(err)
	}

}

func (operator *Operator) WriteBTCEUR(ticker *database.Ticker) {
	fmt.Println("ticker: ", ticker)
	err := operator.database.Write(*ticker)
	if err != nil {
		log.Error(err)
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
				"unable to read json as message",
			)
		}

		ask, bid, err := preparePrices(message.BestAsk, message.BestBid)
		if err != nil {
			return karma.Format(
				err,
				"unable to handle prices",
			)
		}

		operator.distributionChan <- &database.Ticker{
			Symbol:   message.ProductID,
			AskPrice: ask,
			BidPrice: bid,
		}
	}

	return nil
}

func preparePrices(ask, bid string) (float64, float64, error) {
	if ask == "" || bid == "" {
		return 0, 0, nil
	}

	resultAsk, err := strconv.ParseFloat(ask, 64)
	if err != nil {
		return 0, 0, karma.Format(
			err,
			"unable to convert ask price to float64, ask: %s",
			ask,
		)
	}

	resultBid, err := strconv.ParseFloat(bid, 64)
	if err != nil {
		return 0, 0, karma.Format(
			err,
			"unable to convert bid to float64, bid: %s",
			bid,
		)
	}

	return resultAsk, resultBid, nil
}

func (operator *Operator) WritePrices() error {
	err := operator.database.CreateTicksTable()
	if err != nil {
		return karma.Format(
			err,
			"unable to create schema in the database",
		)
	}

	currency_1 := database.Ticker{
		Timestamp: time.Now().Unix(),
		Symbol:    "BTC",
		AskPrice:  23.343,
		BidPrice:  55.32,
	}

	err = operator.database.Write(currency_1)
	if err != nil {
		return karma.Format(
			err,
			"unable to write currency to the database",
		)
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
