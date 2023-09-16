package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"
)

type CurrenciesResponse struct {
	Currencies map[string]string `json:"currencies"`
}

type ExchangeRatesResponse struct {
	Timestamp int64              `json:"timestamp"`
	Base      string             `json:"base"`
	Rates     map[string]float64 `json:"rates"`
}

func main() {
	// Parse command line arguments
	datePtr := flag.String("date", "", "Specify the date in the format Y-M-d (e.g., 2023-09-14)")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		print("Error loading .env file")
	}
	if os.Getenv("OPEN_EXCHANGE_RATES_TOKEN") == "" {
		panic("OPEN_EXCHANGE_RATES_TOKEN is not specified")
	}
	if os.Getenv("BASE_SYMBOL") == "" {
		panic("BASE_SYMBOL is not specified")
	}
	if os.Getenv("ECONUMO_API_URL") == "" {
		panic("ECONUMO_API_URL is not specified")
	}
	if os.Getenv("ECONUMO_API_KEY") == "" {
		panic("ECONUMO_API_KEY is not specified")
	}

	openExchangeRatesToken := os.Getenv("OPEN_EXCHANGE_RATES_TOKEN")
	symbols := os.Getenv("SYMBOLS")
	baseSymbol := os.Getenv("BASE_SYMBOL")
	econumoAPIURL := os.Getenv("ECONUMO_API_URL")
	econumoAPIKey := os.Getenv("ECONUMO_API_KEY")

	currenciesURL := fmt.Sprintf("https://openexchangerates.org/api/currencies.json?app_id=%s", openExchangeRatesToken)
	respCurrencies, err := http.Get(currenciesURL)

	if err != nil {
		log.Fatalf("Failed to fetch currencies: %v", err)
	}
	defer respCurrencies.Body.Close()

	if respCurrencies.StatusCode != 200 {
		log.Fatalf("Failed to fetch currencies. Status code: %d", respCurrencies.StatusCode)
	}

	// Parse the currencies JSON response
	var currenciesResp CurrenciesResponse
	err = json.NewDecoder(respCurrencies.Body).Decode(&currenciesResp)
	if err != nil {
		log.Fatalf("Failed to parse currencies response: %v", err)
	}

	// Create a filtered map of currencies based on SYMBOLS
	filteredCurrencies := make(map[string]string)
	symbolList := strings.Split(symbols, ",")
	for _, symbol := range symbolList {
		if currencyName, ok := currenciesResp.Currencies[symbol]; ok {
			filteredCurrencies[symbol] = currencyName
		}
	}

	// Convert filteredCurrencies to JSON
	currenciesJSON, err := json.Marshal(filteredCurrencies)
	if err != nil {
		log.Fatalf("Failed to marshal filteredCurrencies to JSON: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/system/import-currency-list", econumoAPIURL), strings.NewReader(string(currenciesJSON)))
	if err != nil {
		log.Fatalf("Failed to create request to send currencies data to econumo API: %v", err)
	}
	req.Header.Set("Authorization", econumoAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Failed to send currencies data to econumo API: %v", err)
	}
	defer resp.Body.Close()

	// Fetch currency rates from openexchangerates.org
	var openExchangeRatesURL string

	if *datePtr == "" {
		openExchangeRatesURL = fmt.Sprintf("https://openexchangerates.org/api/latest.json?app_id=%s", openExchangeRatesToken)
	} else {
		openExchangeRatesURL = fmt.Sprintf("https://openexchangerates.org/api/historical/%s.json?app_id=%s", *datePtr, openExchangeRatesToken)
	}
	openExchangeRatesURL += fmt.Sprintf("&base=%s", baseSymbol)
	if symbols != "" {
		openExchangeRatesURL += fmt.Sprintf("&symbols=%s", symbols)
	}

	resp, err = http.Get(openExchangeRatesURL)

	if err != nil {
		log.Fatalf("Failed to fetch currency rates: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Failed to fetch currency rates. Status code: %d", resp.StatusCode)
	}

	// Parse the currency rates JSON response
	var exchangeRates ExchangeRatesResponse
	err = json.NewDecoder(resp.Body).Decode(&exchangeRates)
	if err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}

	// Send currency rates data to econumo API
	exchangeRatesJSON, err := json.Marshal(exchangeRates)
	if err != nil {
		log.Fatalf("Failed to marshal combinedData to JSON: %v", err)
	}

	req, err = http.NewRequest("POST", fmt.Sprintf("%s/api/v1/system/import-currency-rates", econumoAPIURL), strings.NewReader(string(exchangeRatesJSON)))
	if err != nil {
		log.Fatalf("Failed to create request to send data to econumo API: %v", err)
	}
	req.Header.Set("Authorization", econumoAPIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)

	if err != nil {
		log.Fatalf("Failed to send data to econumo API: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("Data sent to econumo API successfully")
}
