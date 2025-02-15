package callbacks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
)

// GeoResponse는 OpenWeatherMap Geocoding API 응답 구조체
type GeoResponse struct {
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
	Name string  `json:"name"`
}

// WeatherSummaryResponse는 OpenWeatherMap One Call API 3.0 `/onecall/day_summary` 응답 구조체
type WeatherSummaryResponse struct {
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
	Timezone   string  `json:"tz"`
	Date       string  `json:"date"`
	Units      string  `json:"units"`
	CloudCover struct {
		Afternoon float64 `json:"afternoon"`
	} `json:"cloud_cover"`
	Humidity struct {
		Afternoon float64 `json:"afternoon"`
	} `json:"humidity"`
	Precipitation struct {
		Total float64 `json:"total"`
	} `json:"precipitation"`
	Pressure struct {
		Afternoon float64 `json:"afternoon"`
	} `json:"pressure"`
	Temperature struct {
		Min       float64 `json:"min"`
		Max       float64 `json:"max"`
		Afternoon float64 `json:"afternoon"`
		Night     float64 `json:"night"`
		Evening   float64 `json:"evening"`
		Morning   float64 `json:"morning"`
	} `json:"temperature"`
	Wind struct {
		Max struct {
			Speed     float64 `json:"speed"`
			Direction float64 `json:"direction"`
		} `json:"max"`
	} `json:"wind"`
}

// APIErrorResponse는 API 호출 실패 시 반환되는 JSON 구조체
type APIErrorResponse struct {
	Code       int      `json:"cod"`
	Message    string   `json:"message"`
	Parameters []string `json:"parameters"`
}

// getCoordinates: 도시명을 위도/경도로 변환
func getCoordinates(apiKey string, city string) (float64, float64, error) {
	baseURL := "http://api.openweathermap.org/geo/1.0/direct"
	params := url.Values{}
	params.Set("q", city)
	params.Set("limit", "1")
	params.Set("appid", apiKey)

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(reqURL)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("위도/경도 변환 API 호출 실패: %s", resp.Status)
	}

	var geoData []GeoResponse
	if err := json.NewDecoder(resp.Body).Decode(&geoData); err != nil {
		return 0, 0, err
	}

	if len(geoData) == 0 {
		return 0, 0, fmt.Errorf("도시명을 찾을 수 없습니다: %s", city)
	}

	return geoData[0].Lat, geoData[0].Lon, nil
}

// getWeatherSummary: `/onecall/day_summary` API 호출하여 특정 날짜의 날씨 요약 가져오기
func getWeatherSummary(apiKey string, date string, latitude, longitude float64, unit, lang string) (*WeatherSummaryResponse, error) {
	baseURL := "https://api.openweathermap.org/data/3.0/onecall/day_summary"
	params := url.Values{}
	params.Set("lat", fmt.Sprintf("%f", latitude))
	params.Set("lon", fmt.Sprintf("%f", longitude))
	params.Set("date", date)
	params.Set("appid", apiKey)
	params.Set("units", unit)
	params.Set("lang", lang)

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// API 에러 응답 처리
	if resp.StatusCode != http.StatusOK {
		var apiErr APIErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return nil, fmt.Errorf("API 호출 실패: HTTP %d (응답 해석 실패)", resp.StatusCode)
		}
		return nil, fmt.Errorf("API 호출 실패: HTTP %d, 메시지: %s, 매개변수: %v", apiErr.Code, apiErr.Message, apiErr.Parameters)
	}

	var weatherResp WeatherSummaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		return nil, err
	}

	return &weatherResp, nil
}

func GetWeather(s *service, ctx context.Context, args []byte, metadata Metadata) (any, error) {
	var req struct {
		Location string `json:"location"`
		Date     string `json:"date"`
		Unit     string `json:"unit"`
	}

	if err := json.Unmarshal(args, &req); err != nil {
		return nil, errors.WithStack(err)
	}

	latitude, longitude, err := getCoordinates(s.config.OpenWeatherApiKey, req.Location)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert coordinates")
	}

	weatherSummary, err := getWeatherSummary(s.config.OpenWeatherApiKey, req.Date, latitude, longitude, req.Unit, "en")
	if err != nil {
		return nil, errors.Wrapf(err, "error occurred while fetching weather information")
	}

	return &weatherSummary, nil
}

func init() {
	dispatchFunctions["get_weather"] = GetWeather
}
