package util

import "log"

func Must[T any](f func() (T, error)) T {
	result, err := f()
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func MustStr[T any](f func(string) (T, error), param string) T {
	result, err := f(param)
	if err != nil {
		log.Fatal(err)
	}
	return result
}
