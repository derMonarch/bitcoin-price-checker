**Last traded bitcoin price fetcher**

Fetch up-to-date Bitcoin prices by currency pair


## Status

Current supported currency pairs are BTCEUR / BTCCHF / BTCUSD


## Install with Docker

If you're using Docker, you can build the image with and run the service locally.

```bash
docker build -t btcprices .
```

The service is configured to run on port 8080 - this port is exposed in the docker image

## Run the built image

```
docker run -it -p 8080:8080 btcprices
```

## Build from Source

If you want to build the application from source you can use the following command:

```bash
go build -o btcprices cmd/btcprice.go
```

## Get Started
**How to use the service**

The service exposes a REST API with the route /api/v1/ltp

**Current Handlers**

Currently, there is one handler handling a GET request to /api/v1/ltp
Supported currency pairs are BTC/EUR BTC/CHF BTC/USD

You can either:
- Call the endpoint as is with /api/v1/ltp and get all supported currency pairs back
- Define query parameters to select specific supported currency pairs for e.g. /api/v1/ltp?pairs=BTCEUR
- You can also fetch more pairs selectively with /api/v1/ltp?pairs=BTCEUR&pairs=BTCCHF

**Environment variables**

API_LAST_TRADED_PRICE - Kraken REST API to fetch currency pairs

The environment variable is set in the .env file
