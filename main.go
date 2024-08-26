package main

import (
	"fmt"
	"log"

	"github.com/yourusername/lichess-puzzle-downloader/downloader"
)

func main() {
	puzzleCount, err := downloader.DownloadAndStorePuzzles()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total puzzles processed: %d\n", puzzleCount)
}
