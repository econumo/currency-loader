package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type CurrenciesResponse struct {
	Currencies map[string]string `json:"-"`
}

type CurrenciesRequest struct {
	Items []string `json:"items"`
}

type ExchangeRatesResponse struct {
	Timestamp int64              `json:"timestamp"`
	Base      string             `json:"base"`
	Rates     map[string]float64 `json:"rates"`
}

type ExchangeRatesRequest struct {
	Timestamp int64            `json:"timestamp"`
	Base      string           `json:"base"`
	Items     []CurrenciesRate `json:"items"`
}

type CurrenciesRate struct {
	Code string `json:"code"`
	Rate string `json:"rate"`
}

func main() {
	// Parse command line arguments
	datePtr := flag.String("date", "", "Specify the date in the format Y-M-d (e.g., 2023-09-14)")
	envPathPtr := flag.String("env", ".env", "Specify the path to the .env file")
	flag.Parse()

	err := godotenv.Load(*envPathPtr)
	if err != nil {
		fmt.Printf(".env file not found at %s", *envPathPtr)
	}
	if os.Getenv("OPEN_EXCHANGE_RATES_TOKEN") == "" {
		panic("OPEN_EXCHANGE_RATES_TOKEN is not specified")
	}
	if os.Getenv("ECONUMO_CURRENCY_BASE") == "" {
		panic("ECONUMO_CURRENCY_BASE is not specified")
	}
	if os.Getenv("ECONUMO_BASE_URL") == "" {
		panic("ECONUMO_BASE_URL is not specified")
	}
	if os.Getenv("ECONUMO_SYSTEM_API_KEY") == "" {
		panic("ECONUMO_SYSTEM_API_KEY is not specified")
	}

	openExchangeRatesToken := os.Getenv("OPEN_EXCHANGE_RATES_TOKEN")
	symbols := os.Getenv("OPEN_EXCHANGE_RATES_SYMBOLS")
	baseSymbol := os.Getenv("ECONUMO_CURRENCY_BASE")
	econumoBaseURL := os.Getenv("ECONUMO_BASE_URL")
	econumoAPIKey := os.Getenv("ECONUMO_SYSTEM_API_KEY")

	currenciesURL := fmt.Sprintf("https://openexchangerates.org/api/currencies.json?app_id=%s", openExchangeRatesToken)
	respCurrencies, err := http.Get(currenciesURL)

	if err != nil {
		log.Fatalf("Failed to fetch currencies: %v", err)
	}

	if respCurrencies.StatusCode != 200 {
		log.Fatalf("Failed to fetch currencies. Status code: %d", respCurrencies.StatusCode)
	}

	// Parse the currencies JSON response
	var currenciesResp CurrenciesResponse
	// Unmarshal the JSON data into the struct
	err = json.NewDecoder(respCurrencies.Body).Decode(&currenciesResp.Currencies)
	if err != nil {
		log.Fatalf("Failed to parse currencies response: %v", err)
	}
	defer respCurrencies.Body.Close()

	// Create a filtered map of currencies based on OPEN_EXCHANGE_RATES_SYMBOLS
	var filteredCurrencies CurrenciesRequest
	symbolList := strings.Split(symbols, ",")
	for code, _ := range currenciesResp.Currencies {
		if symbols == "" {
			filteredCurrencies.Items = append(filteredCurrencies.Items, code)
		} else {
			for _, symbol := range symbolList {
				if symbol == code {
					filteredCurrencies.Items = append(filteredCurrencies.Items, code)
				}
			}
		}
	}

	// Convert filteredCurrencies to JSON
	currenciesJSON, err := json.Marshal(filteredCurrencies)
	if err != nil {
		log.Fatalf("Failed to marshal filteredCurrencies to JSON: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/system/import-currency-list", econumoBaseURL), strings.NewReader(string(currenciesJSON)))
	if err != nil {
		log.Fatalf("Failed to create request to send currencies to Econumo API: %v", err)
	}
	req.Header.Set("Authorization", econumoAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Failed to send currencies to Econumo API: %v\n\n", err)
	}
	if resp.StatusCode != 200 {
		log.Fatalf("Failed to send currencies to econumo API. Status code: %d\n\n", resp.StatusCode)
	}
	defer resp.Body.Close()
	fmt.Printf("Currencies sent to econumo API successfully: %s\n\n", string(currenciesJSON))

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
	var exchangeRatesRequest ExchangeRatesRequest
	exchangeRatesRequest.Timestamp = exchangeRates.Timestamp
	exchangeRatesRequest.Base = exchangeRates.Base
	for code, rate := range exchangeRates.Rates {
		exchangeRatesRequest.Items = append(exchangeRatesRequest.Items, CurrenciesRate{Code: code, Rate: fmt.Sprintf("%f", rate)})
	}
	exchangeRatesJSON, err := json.Marshal(exchangeRatesRequest)
	if err != nil {
		log.Fatalf("Failed to marshal combinedData to JSON: %v", err)
	}

	req, err = http.NewRequest("POST", fmt.Sprintf("%s/api/v1/system/import-currency-rates", econumoBaseURL), strings.NewReader(string(exchangeRatesJSON)))
	if err != nil {
		log.Fatalf("Failed to create request to send currency rates to Econumo API: %v", err)
	}
	req.Header.Set("Authorization", econumoAPIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send currency rates to Econumo API: %v\n\n", err)
	}
	if resp.StatusCode != 200 {
		log.Fatalf("Failed to send currency rates to Econumo API. Status code: %d\n\n", resp.StatusCode)
	}
	defer resp.Body.Close()

	// Add better error handling and response logging
	if resp.StatusCode != http.StatusOK {
		// Read the response body for error details
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed to read error response body: %v", err)
		}
		log.Fatalf("Failed to send data to econumo API. Status code: %d, Response: %s",
			resp.StatusCode, string(bodyBytes))
	}

	fmt.Printf("Currency rates sent to econumo API successfully: %s\n", string(exchangeRatesJSON))
}
