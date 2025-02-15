package callbacks_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/habiliai/habiliai/api/pkg/callbacks"
	"github.com/habiliai/habiliai/api/pkg/config"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGetWeather(t *testing.T) {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		t.Skip("OPENWEATHER_API_KEY 환경 변수가 설정되지 않았습니다")
	}

	s := callbacks.NewService(&config.HabApiConfig{OpenWeatherApiKey: apiKey})
	contents, err := s.Dispatch(context.TODO(), "get_weather", []byte(`{
"location": "Seoul",
"date": "2025-02-10",
"unit": "metric"
}`), callbacks.Metadata{})

	require.NoError(t, err)

	var weatherSummary callbacks.WeatherSummaryResponse
	require.NoError(t, json.Unmarshal(contents, &weatherSummary))

	// 3. 출력
	fmt.Printf("🕒 시간대: %s\n", weatherSummary.Timezone)
	fmt.Printf("📅 날짜: %s\n", weatherSummary.Date)
	fmt.Printf("🌡️ 최고 기온: %.2f°C\n", weatherSummary.Temperature.Max)
	fmt.Printf("🌡️ 최저 기온: %.2f°C\n", weatherSummary.Temperature.Min)
	fmt.Printf("🌡️ 오후 기온(12:00): %.2f°C\n", weatherSummary.Temperature.Afternoon)
	fmt.Printf("🌡️ 아침 기온(06:00): %.2f°C\n", weatherSummary.Temperature.Morning)
	fmt.Printf("🌡️ 저녁 기온(18:00): %.2f°C\n", weatherSummary.Temperature.Evening)
	fmt.Printf("🌡️ 밤 기온(00:00): %.2f°C\n", weatherSummary.Temperature.Night)
	fmt.Printf("☁️ 오후 구름량: %.2f\n", weatherSummary.CloudCover.Afternoon)
	fmt.Printf("💧 오후 습도: %.2f\n", weatherSummary.Humidity.Afternoon)
	fmt.Printf("🌬️ 최대 풍속: %.2fm/s (방향: %.2f°)\n", weatherSummary.Wind.Max.Speed, weatherSummary.Wind.Max.Direction)
	fmt.Printf("🌧️ 강수량: %.2fmm\n", weatherSummary.Precipitation.Total)
	fmt.Printf("🗜️ 오후 기압: %.2fhPa\n", weatherSummary.Pressure.Afternoon)
}
