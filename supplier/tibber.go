package supplier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type tibber struct {
	Supplier
	config Config
}

type TibberResponse struct {
	Data struct {
		Viewer struct {
			Homes []struct {
				TodayPrices    []HourlyPrice `json:"todayPrices"`
				TomorrowPrices []HourlyPrice `json:"tomorrowPrices"`
			} `json:"homes"`
		} `json:"viewer"`
	} `json:"data"`
}

func (t *tibber) GetHourlyPrices() (*Prices, error) {

}

func (t *tibber) Refresh() error {
	var jsonData = []byte(`{"query": "YourGraphQLQuery"}`) // replace "YourGraphQLQuery"
	req, err := http.NewRequest("POST", "https://api.tibber.com/v1-beta/gql", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer YourAccessToken") // replace "YourAccessToken"

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var tibberResponse TibberResponse
	err = json.Unmarshal(body, &tibberResponse)
	if err != nil {
		return nil, err
	}

	return &tibberResponse, nil
}

func main() {
	data, err := getTibberData()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Accessing the hourly price data
	for _, home := range data.Data.Viewer.Homes {
		fmt.Println("Today's hourly prices:")
		for _, price := range home.TodayPrices {
			fmt.Printf("Hour: %d, Price: %.2f\n", price.Hour, price.Price)
		}

		fmt.Println("Tomorrow's hourly prices:")
		for _, price := range home.TomorrowPrices {
			fmt.Printf("Hour: %d, Price: %.2f\n", price.Hour, price.Price)
		}
	}
}
