package spot

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var url = "https://api.spot-hinta.fi/TodayAndDayForward"

type SpotPrice struct {
	Rank         int
	DateTime     time.Time
	PriceNoTax   float64
	PriceWithTax float64
}

type SpotHintaApiState struct {
	SpotPrices []SpotPrice
	LastCheck  time.Time
}

type SpotHintaApiClient struct {
	State SpotHintaApiState
}

func CreateSpotHintaApiClient() SpotHintaApiClient {
	apiClient := SpotHintaApiClient{}
	data, err := os.ReadFile("spot_hinta_api_state.json")

	if err != nil {
		return apiClient
	}

	err = json.Unmarshal(data, &apiClient)
	if err != nil {
		log.Printf("Failed to load Spot Hinta Api Client: %v", err)
	}
	return apiClient
}

func (api SpotHintaApiClient) GetPrices() []SpotPrice {
	return api.State.SpotPrices
}

func (api *SpotHintaApiClient) SaveState() {
	data, err := json.Marshal(api)
	if err != nil {
		log.Printf("Failed to save Spot Hinta Api Client: %v", err)
		return
	}

	err = os.WriteFile("spot_hinta_api_state.json", data, 0644)
	if err != nil {
		log.Printf("Failed to save Spot Hinta Api Client: %v", err)
	}
}

func (api *SpotHintaApiClient) UpdatePrices() error {
	log.Println("Updating Spot Hinta API prices...")
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch data: %s", resp.Status)
	}
	var prices []SpotPrice
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read data: %v", err)
	}
	err = json.Unmarshal(body, &prices)
	if err != nil {
		return fmt.Errorf("failed to parse data: %v", err)
	}
	api.State.SpotPrices = prices
	api.State.LastCheck = time.Now()
	api.SaveState()
	return nil
}

func (api *SpotHintaApiClient) PollPrices() {
	log.Println("Polling Spot Hinta API...")
	for {
		minutesUntilNextCheck := time.Until(api.State.LastCheck.Add(1 * time.Hour))
		if minutesUntilNextCheck > 0 {
			time.Sleep(minutesUntilNextCheck)
		}
		err := api.UpdatePrices()
		if err != nil {
			log.Printf("Error updating prices: %v", err)
		}
	}
}

func (api *SpotHintaApiClient) GetCurrentPrice() *SpotPrice {
	if len(api.State.SpotPrices) == 0 {
		return nil
	}
	lastIndex := len(api.State.SpotPrices) - 1
	for i, price := range api.State.SpotPrices {
		if price.DateTime.After(time.Now()) {
			return &api.State.SpotPrices[i-1]
		}
		if price.DateTime.Equal(time.Now()) {
			return &price
		}
		if i == lastIndex {
			return &price
		}
	}
	return nil
}
