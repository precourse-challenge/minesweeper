package util

import (
	"log"
)

func FatalIfError[T any](value T, err error) T {
	if err != nil {
		log.Fatalf("[FATAL] %s\n", err)
	}
	return value
}
