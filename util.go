package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

func MapMerge[M ~map[K]V, K comparable, V any](maps ...M) M {
	result := make(map[K]V)

	for _, toMerge := range maps {
		for k, v := range toMerge {
			if _, found := result[k]; found {
				log.Fatalln(fmt.Errorf("EY NOT DUPLICATE KEYS"))
			}

			result[k] = v
		}
	}

	return result
}

func abspath(path string) string {
	path = filepath.Clean(path)
	if strings.Contains(path, "~") {
		if home, err := os.UserHomeDir(); err == nil {
			path = strings.Replace(path, "~", home, 1)
		}
	}

	path, _ = filepath.Abs(path)
	return path
}

func die(msg string) {
	os.Stderr.WriteString(msg)
	os.Exit(1)
}

