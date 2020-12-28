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

Download ticks and save them to the database.

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

	postgresDB := database.NewDatabase(
		config.Database.Name, config.Database.User, config.Database.Password,
	)
	defer postgresDB.Close()

	websocket, err := websocket.NewWebSocketConnection(config.WebsocketURL)
	if err != nil {
		log.Fatal(err)
	}

	operator := operator.NewOperator(
		config, postgresDB, websocket, make(chan *database.Ticks),
	)

	log.Info("creating 'ticks' table")
	err = operator.CreateTicksTable()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("start of ticks handling")
	go operator.HandleTicks()

	log.Info("receiving ticks")
	err = operator.ReceiveTicks()
	if err != nil {
		log.Fatal(err)
	}
}
