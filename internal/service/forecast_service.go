package service

import (
	"go_api/internal/adapter"
	"go_api/internal/serializer"
)

func GetForecast(locationKey string) (interface{}, error) {
	data, err := adapter.FetchForecast(locationKey)
	if err != nil {
		return nil, err
	}
	return serializer.FormatForecast(data), nil
}
