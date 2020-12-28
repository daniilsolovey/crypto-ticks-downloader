package database

import (
	"strings"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"
)

type Database struct {
	name     string
	user     string
	password string
	client   *pg.DB
}

type Ticker struct {
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

func (database *Database) Write(ticker Ticker) error {
	ticker.Timestamp = makeTimestamp()
	_, err := database.client.Model(&ticker).Insert()
	if err != nil {
		return karma.Format(
			err,
			"unable to write ticks to the database",
		)
	}

	log.Debugf(
		karma.Describe("symbol", ticker.Symbol).
			Describe("timestamp", ticker.Timestamp).
			Describe("timestamp", ticker.AskPrice).
			Describe("timestamp", ticker.BidPrice),
		"ticks were successfully written to the database",
	)

	return nil
}

func (database *Database) CreateTicksTable() error {
	var model *Ticker
	err := database.client.Model(model).CreateTable(&orm.CreateTableOptions{})
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return karma.Format(
			err,
			"unable to create table in the database",
		)
	}

	log.Info("table succesfully created")
	return nil
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
