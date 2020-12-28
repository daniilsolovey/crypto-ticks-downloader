package main

import (
	"github.com/daniilsolovey/crypto-ticks-downloader/internal/config"
	"github.com/daniilsolovey/crypto-ticks-downloader/internal/database"
	"github.com/daniilsolovey/crypto-ticks-downloader/internal/operator"
	"github.com/daniilsolovey/crypto-ticks-downloader/internal/websocket"
	"github.com/docopt/docopt-go"
	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"
)

var version = "[manual build]"

var usage = `crypto-ticks-downloader

Download ticks and save it to the database.

Usage:
  crypto-ticks-downloader [options]

Options:
  -c --config <path>                Read specified config file. [default: config.yaml]
  --debug                           Enable debug messages.
  -v --version                      Print version.
  -h --help                         Show this help.
`

func main() {
	args, err := docopt.ParseArgs(
		usage,
		nil,
		"crypto-ticks-downloader "+version,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof(
		karma.Describe("version", version),
		"crypto-ticks-downloader started",
	)

	if args["--debug"].(bool) {
		log.SetLevel(log.LevelDebug)
	}

	log.Infof(nil, "loading configuration file: %q", args["--config"].(string))

	config, err := config.Load(args["--config"].(string))
	if err != nil {
		log.Fatal(err)
	}

	log.Infof(
		karma.Describe("database", config.Database.Name),
		"connecting to the database",
	)

	newDatabase := database.NewDatabase(
		config.Database.Name, config.Database.User, config.Database.Password,
	)
	defer newDatabase.Close()
	websocket, err := websocket.NewWebSocketConnection(config.WebsocketURL)
	if err != nil {
		log.Fatal(err)
	}

	channel := make(chan *database.Ticker)

	operator := operator.NewOperator(config, newDatabase, websocket, channel)

	// err = operator.WritePrices()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	go operator.DistributeTickers()

	err = operator.GetPrices()
	if err != nil {
		log.Fatal(err)
	}

}
