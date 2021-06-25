package xconfig

import (
	"os"
	"strconv"
)

func Get(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func GetInt(value string, fallback int) int {
	newValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return newValue
}
