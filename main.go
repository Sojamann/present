package main

import (
	"bufio"
	"os"
	"strings"
)

func splitAt(sep string) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		 // Return nothing if at end of file and no data passed
        if atEOF && len(data) == 0 {
            return 0, nil, nil
        }

        // Find the index of the input of the separator sep
        if i := strings.Index(string(data), sep); i >= 0 {
            return i + len(sep), data[0:i], nil
        }

        // If at end of file with data return the data
        if atEOF {
            return len(data), data, nil
        }

		return 0, nil, nil
	}
}

func main() {
	file, _ := os.Open(os.Args[1])
	scanner := bufio.NewScanner(file)
	scanner.Split(splitAt("---"))

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		println("Slide: ", text)
	}
}

