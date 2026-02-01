package main

import (
	"fmt"
	"os"
)

func main() {
	url := "https://ja.wikipedia.org/wiki/%E3%83%92%E3%83%97%E3%83%8E%E3%82%B7%E3%82%B9%E3%83%9E%E3%82%A4%E3%82%AF"

	characters, err := ScrapeCharacters(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scraping: %v\n", err)
		os.Exit(1)
	}

	if err := WriteCSV(os.Stdout, characters); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing CSV: %v\n", err)
		os.Exit(1)
	}
}
