package env

// This package is a helper package:
// (1) it is very simple
// (2) it only depends on stdlib
//
// It's so simple in fact that I don't care about decoupling from it,
// so I wont write any interfaces here.

import (
	"fmt"
	"os"
	"strconv"
)

func GetString(key string, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}

func MustGetString(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(
			fmt.Sprintf("can't start program: required env variable '%s' is unset or empty", key),
		)
	}
	return v
}

func GetInt(key string, defaultValue int) int {
	v, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return defaultValue
	}
	return v
}

func MustGetInt(key string) int {
	v, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		panic(
			fmt.Sprintf("can't start program: required env variable '%s' is unset or empty", key),
		)
	}
	return v
}

func GetFloat(key string, defaultValue float64) float64 {
	v, err := strconv.ParseFloat(os.Getenv(key), 64)
	if err != nil {
		return defaultValue
	}
	return v
}

func MustGetFloat(key string) float64 {
	v, err := strconv.ParseFloat(os.Getenv(key), 64)
	if err != nil {
		panic(
			fmt.Sprintf("can't start program: required env variable '%s' is unset or empty", key),
		)
	}
	return v
}
