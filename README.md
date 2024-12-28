# Currency loader for Econumo

The Currency Loader for Econumo acts as a bridge between [Open Exchange Rates](https://openexchangerates.org) and
the [Econumo System API](https://econumo.com/docs/api/).

---
> [!NOTE]
> In the free tier of Open Exchange Rates, only **USD** is supported as a base currency.
> This shouldn't be an issue for most users, as the base currency is only used for conversions.
>
> For example: If you have most of your accounts in CAD (Canadian Dollar) and one savings account in USD, the currency
> conversion will occur in two scenarios:
> 1. When you transfer money from your CAD account to your USD account.
> 2. In your budget, if you choose to convert spending to other currencies.

## Configuration

Please, configure the following environment variables:

- `ECONUMO_CURRENCY_BASE` - the base currency symbol (e.g. USD) (required). This variable must match the one used for your Econumo instance.
- `ECONUMO_SYSTEM_API_KEY` - the **Econumo System API** key (required). This variable must also match the one used for your Econumo instance.
- `ECONUMO_BASE_URL` - Econumo Base URL (e.g. https://demo.econumo.com) (required).
- `OPEN_EXCHANGE_RATES_TOKEN` - your Open Exchange Rates **App ID** (required).
- `OPEN_EXCHANGE_RATES_SYMBOLS` - the list of currency symbols to load (e.g. USD,EUR,GBP). This variable can be left empty to load all available currencies.

## Usage

Loading the exchange rates once a day is sufficient. 

### Build the solution

```bash
go build -o currency-loader
```

```bash
GOOS=linux GOARCH=amd64 go build -o currency-loader
```

### Loading currencies and their exchange rates

To load the currency exchange rates for today, use the following command:

```bash
./currency-loader 
```

Alternatively, you can load currencies for a specific date in the past:

```bash
./currency-loader --date=2023-09-14
```
