# Currency loader for Econumo

Currency loader for Econumo is the bridge between the https://openexchangerates.org and Econumo API.

## Configuration

Please, configure the following environment variables:
- `OPEN_EXCHANGE_RATES_TOKEN` - your Open Exchange Rates API key
- `ECONUMO_CURRENCY_BASE` - the base currency symbol (e.g. USD)
- `OPEN_EXCHANGE_RATES_SYMBOLS` - the list of currency symbols to load (e.g. USD,EUR,GBP)
- `ECONUMO_API_URL` - Econumo URL (e.g. https://app.econumo.com)
- `ECONUMO_SYSTEM_API_KEY` - System Econumo API key

## Usage

#### Build the solution
```bash
go build -o currency-loader
```

```bash
GOOS=linux GOARCH=amd64 go build -o currency-loader
```

#### Run the command
```bash
./currency-loader --date=2023-09-14
```

OR

```bash
./currency-loader 
```

