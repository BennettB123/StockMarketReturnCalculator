package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	// Check for correct number of command-line arguments
	if len(os.Args) < 3 {
		fmt.Printf("Correct usage: %s <path_to_holdings_file> <api_key>\n"+
			"\nProgram expects <path_to_holdings_file> to be a JSON file with the following structure:\n"+
			"  [\n"+
			"    {\n"+
			"      \"Ticker\": <ticker_symbol> (string),\n"+
			"      \"NumShares\": <number_of_shares> (int)>,\n"+
			"      \"AvgPricePerShare\": <average_price_pad_per_share> (float)>\n"+
			"    },\n"+
			"    ...\n"+
			"  ]\n"+
			"\nProgram expects <api_key> to be a single line file containing a valid Tiingo API key\n", os.Args[0])
		os.Exit(1)
	}
	userHoldings := getHoldings(os.Args[1])
	printTotalReturn(userHoldings)

}

// Structure to hold daily stock info from a daily Tiingo API call
type dailyInfo struct {
	AdjClose    float32   `json:"adjClose"`
	AdjHigh     float32   `json:"adjHigh"`
	AdjLow      float32   `json:"adjLow"`
	AdjOpen     float32   `json:"adjOpen"`
	AdjVolume   float32   `json:"adjVolume"`
	Close       float32   `json:"close"`
	Date        time.Time `json:"date"`
	DivCash     float32   `json:"divCash"`
	High        float32   `json:"high"`
	Low         float32   `json:"low"`
	Open        float32   `json:"open"`
	SplitFactor float32   `json:"splitFactor"`
	Volume      float32   `json:"volume"`
}

// Structure to hold a user's holding info about a stock
type holdingInfo struct {
	Ticker           string
	NumShares        int
	AvgPricePerShare float32
}

// Reads and returns the first line of the file, 'keyPath'
func getAPIKey(keyPath string) string {
	file, err := os.Open(keyPath)
	if err != nil {
		log.Fatalf("There was an error opening '%s': %s", keyPath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

// Parses the JSON file 'filePath' into an array of holdingInfo structures
func getHoldings(filePath string) []holdingInfo {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("There was an error opening '%s': %s", filePath, err)
	}
	defer file.Close()

	var holdings []holdingInfo
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&holdings); err != nil {
		log.Fatalf("There was an error while parsing holdings file: %s", err)
	}

	return holdings
}

// getCurrentStockPrice will send the current info of the stock with ticker symbol, 'ticker', through channel 'c'
func getStockPrice(ticker string, prices map[string]float32, wg *sync.WaitGroup) {
	defer wg.Done()

	// craft API url
	url := "https://api.tiingo.com/tiingo/daily/" + ticker + "/prices?token=" + getAPIKey(os.Args[2])

	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("There was an error contacting the API: %s", err)
	}
	defer response.Body.Close()

	// read the response and unmarshal it to a dailyInfo type
	data, _ := ioutil.ReadAll(response.Body)
	var dataJSON []dailyInfo
	err = json.Unmarshal(data, &dataJSON)
	if err != nil {
		log.Fatalf("There was an error parsing the API response to a 'dailyInfo' type")
	}

	prices[ticker] = dataJSON[0].Close
}

func printTotalReturn(holdings []holdingInfo) {
	// concurrently grab current value for each stock in 'holdings' from Tiingo API
	currentValues := make(map[string]float32)
	var wg sync.WaitGroup
	for _, h := range holdings {
		wg.Add(1)
		go getStockPrice(h.Ticker, currentValues, &wg)
	}
	wg.Wait()

	fmt.Printf("-----------------------------------------------------------\n")
	fmt.Printf("| Ticker | Shares | Avg Cost | Current Value |   Return   |\n")
	fmt.Printf("-----------------------------------------------------------\n")

	for _, h := range holdings {
		fmt.Printf("|  %-6s|  %-6d| %-9.2f|    %-11.2f|  %-10.2f|\n", h.Ticker, h.NumShares, h.AvgPricePerShare, currentValues[h.Ticker], float32(h.NumShares)*(currentValues[h.Ticker]-h.AvgPricePerShare))
		fmt.Printf("-----------------------------------------------------------\n")
	}

	// calculate total number of shares, total cost, total price, and total return.
	var totalShares int = 0
	var totalCost float32 = 0.0
	var totalValue float32 = 0.0
	var totalReturn float32 = 0.0
	for _, h := range holdings {
		totalShares += h.NumShares
		totalCost += float32(h.NumShares) * h.AvgPricePerShare
		totalValue += float32(h.NumShares) * currentValues[h.Ticker]
	}
	totalReturn = totalValue - totalCost

	// print total return of all holdings
	fmt.Printf("| TOTALS |  %-6d| %-9.2f|    %-11.2f|  %-10.2f|\n", totalShares, totalCost, totalValue, totalReturn)
	fmt.Printf("-----------------------------------------------------------\n")
}
