package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Song struct {
	Name      string    `json:"master_metadata_track_name"`
	Artist    string    `json:"master_metadata_album_artist_name"`
	Album     string    `json:"master_metadata_album_album_name"`
	Timestamp time.Time `json:"ts"`
	MsPlayed  int       `json:"ms_played"`
}

func getFiles(dirPath string) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var files []string
	re := regexp.MustCompile(`Audio.*\.json$`)
	for _, entry := range entries {
		if !entry.IsDir() && re.MatchString(entry.Name()) {
			fullPath := filepath.Join(dirPath, entry.Name())
			files = append(files, fullPath)
		}
	}

	if len(files) == 0 {
		return nil, errors.New("No files read")
	}

	return files, nil
}

func parseJson(files []string) ([]Song, error) {
	var allSongs []Song

	for _, f := range files {
		file, err := os.Open(f)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		var songs []Song
		dec := json.NewDecoder(file)
		err = dec.Decode(&songs)
		if err != nil {
			return nil, err
		}

		allSongs = append(allSongs, songs...)
	}

	allSongs = slices.DeleteFunc(allSongs, func(s Song) bool {
		return s.Name == "" || s.Artist == "" || s.Album == ""
	})

	return allSongs, nil
}

func insertSongs(db *sql.DB, songs []Song) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
        INSERT INTO songs (name, artist, album, timestamp, ms_played)
        VALUES (?, ?, ?, ?, ?)
    `)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, s := range songs {
		_, err := stmt.Exec(
			s.Name,
			s.Artist,
			s.Album,
			s.Timestamp,
			s.MsPlayed,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func writeToDb(dbName string, songs []Song) error {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return err
	}
	defer db.Close()

	db.Exec(`PRAGMA journal_mode = WAL;`)
	db.Exec(`PRAGMA synchronous = NORMAL;`)

	_, err = db.Exec(`
		DROP TABLE IF EXISTS songs;
		CREATE TABLE songs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			artist TEXT NOT NULL,
			album TEXT NOT NULL,
			timestamp DATETIME NOT NULL,
			ms_played INTEGER NOT NULL
		);
	`)

	if err != nil {
		return err
	}

	err = insertSongs(db, songs)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Provide path to Spotify listening data")
		os.Exit(1)
	}

	root := os.Args[1]
	files, err := getFiles(root)
	if err != nil {
		log.Fatal(err)
	}

	songs, err := parseJson(files)
	if err != nil {
		log.Fatal(err)
	}

	dbName := "spotify.db"
	err = writeToDb(dbName, songs)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d song listens written to %s\n", len(songs), dbName)
}
