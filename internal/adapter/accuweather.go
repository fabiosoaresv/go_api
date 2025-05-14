package adapter

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

const (
	baseURL = "https://dataservice.accuweather.com"
	apiKey  = "Z1F1GUzpMaHfSKq7Qz3e7lqygFhPVliP"
)

func FetchForecast(locationKey string) (map[string]interface{}, error) {
	client := resty.New()
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"apikey":  apiKey,
			"details": "true",
		}).
		SetResult(map[string]interface{}{}).
		Get(fmt.Sprintf("%s/forecasts/v1/daily/5day/%s", baseURL, locationKey))

	if err != nil {
		return nil, err
	}

	result := resp.Result().(*map[string]interface{})
	return *result, nil
}
