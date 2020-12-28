# Crypto-ticks-downloader

**Crypto-ticks-downloader** is a tool for receiving ticks:
ETH-BTC, BTC-USD, BTC-EUR with WebSocket. Writes it to the Postgres database.

<img alt="" src="https://i.imgur.com/S28xHVo.gif"/>

## Installation

```
go get github.com/daniilsolovey/crypto-ticks-downloader
```

## Usage

##### -c --config \<path>
Read specified config file. [default: config.yaml].

##### --debug
Enable debug messages.

##### -v --version
Print version.

#####  -h --help
Show this help.

## Build
For build program use command:

```
make build
```


## Configuration

Before running a program configure config.yaml file and create a database

Base config settings:

```yaml
websocket_url: wss://ws-feed.pro.coinbase.com
database:
    name: "your_database_name"
    host: "localhost"
    port: 5432
    user: "postgres"
    password: "admin"
```

## RUN

For running program use command:

```
./crypto-ticks-downloader
```