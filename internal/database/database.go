package database

import (
	"strings"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"
)

const (
	ERR_CODE_TABLE_ALREADY_EXISTS = "#42P07"
)

type Database struct {
	name     string
	user     string
	password string
	client   *pg.DB
}

type Ticks struct {
	Timestamp int64
	Symbol    string
	AskPrice  float64
	BidPrice  float64
}

func NewDatabase(
	name, user, password string,
) *Database {
	database := &Database{
		name:     name,
		user:     user,
		password: password,
	}

	database.connect()
	return database
}

func (database *Database) connect() {
	database.client = pg.Connect(
		&pg.Options{
			Database: database.name,
			User:     database.user,
			Password: database.password,
		})
}

func (database *Database) Close() error {
	err := database.client.Close()
	if err != nil {
		return karma.Format(
			err,
			"unable to close connection to the database",
		)
	}

	return nil
}

func (database *Database) Write(ticker Ticks) error {
	ticker.Timestamp = makeTimestamp()
	_, err := database.client.Model(&ticker).Insert()
	if err != nil {
		return karma.Describe("ticker", ticker).Format(
			err,
			"unable to write ticker to the database",
		)
	}

	log.Debugf(
		karma.Describe("symbol", ticker.Symbol).
			Describe("timestamp", ticker.Timestamp).
			Describe("ask_price", ticker.AskPrice).
			Describe("bid_price", ticker.BidPrice),
		"ticker was successfully written to the database",
	)

	return nil
}

func (database *Database) CreateTicksTable() error {
	var model *Ticks
	err := database.client.Model(model).CreateTable(&orm.CreateTableOptions{})
	if strings.Contains(err.Error(), ERR_CODE_TABLE_ALREADY_EXISTS) {
		log.Info("table 'ticks' already exists in the database")
		return nil
	}

	if err != nil {
		return karma.Format(
			err,
			"unable to create 'ticks' table in the database",
		)
	}

	log.Info("table 'ticks' successfully created in the database")
	return nil
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
