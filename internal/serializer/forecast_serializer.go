package serializer

func FormatForecast(data map[string]interface{}) map[string]interface{} {
	headline := data["Headline"].(map[string]interface{})["Text"]
	dailyForecasts := data["DailyForecasts"]

	return map[string]interface{}{
		"headline": headline,
		"days":     dailyForecasts,
	}
}
