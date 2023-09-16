# Currency loader for Econumo

Currency loader for Econumo is the bridge between the https://openexchangerates.org and the Econumo API.

## Configuration

Please, configure the following environment variables:
- `OPEN_EXCHANGE_RATES_TOKEN` - your Open Exchange Rates API key
- `BASE_SYMBOL` - the base currency symbol (e.g. USD)
- `SYMBOLS` - the list of currency symbols to load (e.g. USD,EUR,GBP)
- `ECONUMO_API_URL` - the Econumo API URL (e.g. https://api.econumo.com)
- `ECONUMO_API_TOKEN` - the Econumo API token

## Usage

1. Build the solution
```bash
go build -o currency-loader
```

1. Run the command
```bash
./currency-loader --date=2023-09-14
```
