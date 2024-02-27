package main

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"os"
	"strings"
)

type Note struct {
	ID    int
	GUID  string
	MID   int
	Mod   int
	USN   int
	Tags  string
	Flds  string
	Sfld  int
	Csum  int
	Flags int
	Data  string
}

func main() {

	apkgFilePath := "C:/Users/Brandon/Downloads/Japanese_Core_2000_2k_-_Sorted_w_Audio.apkg"
	apkgData, err := os.ReadFile(apkgFilePath)
	if err != nil {
		log.Fatal(err)
	}
	// Open the zip archive
	zipReader, err := zip.NewReader(bytes.NewReader(apkgData), int64(len(apkgData)))
	if err != nil {
		log.Fatal(err)
	}

	// Iterate through the files in the zip archive
	for _, file := range zipReader.File {

		// Check if the file is named "collection.anki2"
		if strings.EqualFold(file.Name, "collection.anki2") {
			// Open the file inside the zip archive

			fileReader, err := file.Open()
			if err != nil {
				log.Fatal(err)
			}

			// Read the SQLite database file data
			dbData, err := io.ReadAll(fileReader)
			if err != nil {
				log.Fatal(err)
			}

			tmpfile, err := os.CreateTemp("", "example.*.db")
			if err != nil {
				log.Fatal(err)
			}

			if _, err := tmpfile.Write(dbData); err != nil {
				log.Fatal(err)
			}

			if err := tmpfile.Close(); err != nil {
				log.Fatal(err)
			}

			// Open the SQLite database
			db, err := sql.Open("sqlite3", tmpfile.Name())
			if err != nil {
				log.Fatal(err)
			}

			//// Query data
			rows, err := db.Query("SELECT * FROM notes")
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()

			// Iterate through the result set
			for rows.Next() {
				var note Note
				if err := rows.Scan(
					&note.ID,
					&note.GUID,
					&note.MID,
					&note.Mod,
					&note.USN,
					&note.Tags,
					&note.Flds,
					&note.Sfld,
					&note.Csum,
					&note.Flags,
					&note.Data,
				); err != nil {
					log.Fatal(err)
				}
				// Process the retrieved data
				fmt.Printf("%+v\n", note.Flds)
			}

			// Check for errors from iterating over rows
			if err := rows.Err(); err != nil {
				log.Fatal(err)
			}

			defer os.Remove(tmpfile.Name()) // clean up
			defer db.Close()
		}
	}
}
