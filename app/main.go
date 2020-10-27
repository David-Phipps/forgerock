package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
)

// sorts the Time series date keys
func CustomSort(s map[string]interface{}) []string {
	keys := make([]string, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})
	return keys
}

// Parse json and return an average of closing stocks from x number of days from days env var
func ParseJSON(keys []string, prices map[string]interface{}, childKey string) float64 {
	days, _ := strconv.Atoi(os.Getenv("DAYS"))
	i := 1
	var sum float64 = 0

	for _, k := range keys {
		jsonElement := prices[k].(map[string]interface{})
		temp := jsonElement[childKey].(string)
		tempFlt, err := strconv.ParseFloat(temp, 64)
		if err != nil {
			fmt.Print(err.Error())
		}
		sum += tempFlt
		if i >= days {
			break
		}
		i++
	}
	avg := sum / float64(days)
	return avg

}

func main() {
	os.Setenv("STOCK", "MSFT")
	os.Setenv("DAYS", "3")
	os.Setenv("API_KEY", "C227WD9W3LUVKVV9")

	var baseURL string = "https://www.alphavantage.co/query?"
	var combinedURL string = baseURL + "apikey=" + (os.Getenv("API_KEY")) + "&function=TIME_SERIES_DAILY_ADJUSTED&symbol=" + os.Getenv("STOCK")

	response, err := http.Get(combinedURL)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	// Declared an empty interface (wasn't sure how to use a strut for our type since all the date
	// under Time Series Daily are unique)
	var result map[string]interface{}

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal([]byte(responseData), &result)

	// Grab our top level json key
	prices := result["Time Series (Daily)"].(map[string]interface{})

	// Sort the subkeys and get the avg price based on child key
	keys := CustomSort(prices)
	strPrice := fmt.Sprintf("%.2f", ParseJSON(keys, prices, "4. close"))

	//------------------
	//Web server portion
	//------------------

	// create a new `ServeMux`
	mux := http.NewServeMux()

	// handle `/` route
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprint(res, "The 3 day average closing price of "+os.Getenv("STOCK")+" is "+strPrice)
	})

	// listen and serve using `ServeMux`
	http.ListenAndServe(":9000", mux)

}
