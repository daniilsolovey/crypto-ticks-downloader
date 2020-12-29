package operator

import (
	"github.com/daniilsolovey/crypto-ticks-downloader/internal/database"
	"github.com/reconquest/pkg/log"
)

func (operator *Operator) DistributeTickers() {
	for key, value := range operator.channels {
		log.Debugf(nil, "go routine for handling %s ticker started ", key)
		go operator.HandleTicker(value)
	}

	for ticker := range operator.distributionChan {
		for key := range operator.channels {
			if key == ticker.Symbol {
				operator.channels[key] <- ticker
			}
		}
	}
}

func (operator *Operator) HandleTicker(channel chan *database.Ticks) {
	for ticker := range channel {
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
