package habapi

import (
	"github.com/openai/openai-go"
	"strconv"
	"strings"
)

func getStringData(m openai.Metadata, key string) string {
	value, ok := m[key]
	if !ok {
		return ""
	}

	return value
}

func getIntData(m openai.Metadata, key string) int {
	value, ok := m[key]
	if !ok {
		return 0
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return intValue
}

func getIntSliceData(m openai.Metadata, key string) []int {
	value, ok := m[key]
	if !ok {
		return nil
	}

	values := strings.Split(value, ",")
	if len(values) == 0 {
		return nil
	}

	var err error
	intValues := make([]int, len(values))
	for i, v := range values {
		intValues[i], err = strconv.Atoi(v)
		if err != nil {
			return nil
		}
	}

	return intValues
}

func setIntSliceData(m openai.Metadata, key string, value []int) {
	var values []string
	for _, v := range value {
		values = append(values, strconv.Itoa(v))
	}

	m[key] = strings.Join(values, ",")
}

func setStringData(m openai.Metadata, key string, value string) {
	m[key] = value
}

func setIntData(m openai.Metadata, key string, value int) {
	m[key] = strconv.Itoa(value)
}
