# StockMarketReturnCalculator
Go application to calculate your current returns from the stock market

## How to Use
The application can be built by using the following command in the directory where main.go is located:
```
go build -o <executable_name>
```

To use the application, run the command:
```
<executable_name> <path_to_holdings_file> <path_to_api_key>
```
where <path_to_holdings_file> is a JSON file of the form:
```JSON
[
  {
    "Ticker": ticker_symbol,
    "NumShares": number_of_shares>,
    "AvgPricePerShare": average_price_paid_per_share>
  }
]
```
and <path_to_api_key> is a single-line file containing a valid Tiingo API key. You can sign up for a Tiingo API key [here](https://api.tiingo.com/ "Tiingo API Homepage")

### Example Input and Output
Holdings.json contains
```JSON
[
  {
    "Ticker": "MSFT",
    "NumShares": 1,
    "AvgPricePerShare": 184.12
  },
  {
    "Ticker": "GOOG",
    "NumShares": 3,
    "AvgPricePerShare": 1250.90
  }
]
```
running the command: ``` go run main.go Holdings.json <path_to_api_key> ``` will output:
```
-----------------------------------------------------------
| Ticker | Shares | Avg Cost | Current Value |   Return   |
-----------------------------------------------------------
|  MSFT  |  1     | 184.12   |    161.57     |  -22.55    |
-----------------------------------------------------------
|  GOOG  |  3     | 1250.90  |    1298.41    |  142.53    |
-----------------------------------------------------------
| TOTALS |  4     | 3936.82  |    4056.80    |  119.98    |
-----------------------------------------------------------
```
