package downloader

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/klauspost/compress/zstd"
	_ "github.com/mattn/go-sqlite3"
)

const (
	puzzleURL = "https://database.lichess.org/lichess_db_puzzle.csv.zst"
	dbName    = "../puzzler/puzzler.db"
)

func DownloadAndStorePuzzles() (int, error) {
	// Download and decompress the file
	resp, err := http.Get(puzzleURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	zr, err := zstd.NewReader(resp.Body)
	if err != nil {
		return 0, err
	}
	defer zr.Close()

	// Create and connect to SQLite database
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	// Prepare insert statement
	stmt, err := db.Prepare(`INSERT INTO puzzle (
		id, fen, moves, rating, rating_deviation, popularity, nb_plays, themes, game_url, opening_tags
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	// Parse CSV and insert into database
	reader := csv.NewReader(zr)
	reader.Comma = ','
	reader.LazyQuotes = true

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return 0, err
	}

	count := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading record: %v", err)
			continue
		}

		// Insert record into database
		_, err = stmt.Exec(
			record[0], record[1], record[2], record[3], record[4],
			record[5], record[6], record[7], record[8], record[9],
		)
		if err != nil {
			log.Printf("Error inserting record: %v", err)
			continue
		}

		count++
		if count%10000 == 0 {
			fmt.Printf("Processed %d puzzles\n", count)
		}
	}

	return count, nil
}

