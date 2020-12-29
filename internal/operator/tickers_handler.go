package operator

import (
	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"
)

func (operator *Operator) HandleTickers() {
	operator.startRoutines()
	for ticker := range operator.channels.DistributionChan {
		log.Debugf(
			karma.Describe("ticker", ticker),
			"received ticker",
		)
		switch ticker.Symbol {
		case "BTC-USD":
			operator.channels.BTCUSDChan <- ticker
			break
		case "ETH-BTC":
			operator.channels.ETHBTCChan <- ticker
		case "BTC-EUR":
			operator.channels.BTCEURChan <- ticker
		}
	}
}

func (operator *Operator) startRoutines() {
	go operator.HandleBTCUSD()
	go operator.HandleETHBTC()
	go operator.HandleBTCEUR()
}

func (operator *Operator) HandleBTCUSD() {
	for ticker := range operator.channels.BTCUSDChan {
		err := operator.database.Write(*ticker)
		if err != nil {
			log.Errorf(
				err,
				"unable to write ticker to the database, ticker: %s",
				ticker.Symbol,
			)
		}
	}
}

func (operator *Operator) HandleETHBTC() {
	for ticker := range operator.channels.ETHBTCChan {
		err := operator.database.Write(*ticker)
		if err != nil {
			log.Errorf(
				err,
				"unable to write ticker to the database, ticker: %s",
				ticker.Symbol,
			)
		}
	}
}

func (operator *Operator) HandleBTCEUR() {
	for ticker := range operator.channels.BTCEURChan {
		err := operator.database.Write(*ticker)
		if err != nil {
			log.Errorf(
				err,
				"unable to write ticker to the database, ticker: %s",
				ticker.Symbol,
			)
		}
	}
}
