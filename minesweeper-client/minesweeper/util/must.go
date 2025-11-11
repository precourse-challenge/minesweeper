package util

import (
	"fmt"
	"os"
)

func Must[T any](value T, err error) T {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return value
}
