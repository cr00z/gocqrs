package util

import "log"

func Must[T any](f func() (T, error)) T {
	result, err := f()
	if err != nil {
		log.Fatal(err)
	}
	return result
}
