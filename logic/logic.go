package logic

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Brandon689/anki/types"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func sqlNotesTable(filePath string) []types.Note {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		log.Fatal(err)
	}

	// Query data
	rows, err := db.Query("SELECT * FROM notes")
	if err != nil {
		log.Fatal(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
	var notes []types.Note
	// Iterate through the result set
	for rows.Next() {
		var note types.Note
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
		notes = append(notes, note)
		// Process the retrieved data
		//fmt.Printf("%+v\n", note.Flds)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return notes
}

func sqlColTable(filePath string) types.AnkiModel {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		log.Fatal(err)
	}

	// Query data
	rows, err := db.Query("SELECT * FROM col LIMIT 1")
	if err != nil {
		log.Fatal(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	// Iterate through the result set
	for rows.Next() {
		var col types.Col
		if err := rows.Scan(
			&col.ID,
			&col.Crt,
			&col.Mod,
			&col.Scm,
			&col.Ver,
			&col.Dty,
			&col.Usn,
			&col.Ls,
			&col.Conf,
			&col.Models,
			&col.Decks,
			&col.Dconf,
			&col.Tags,
		); err != nil {
			log.Fatal(err)
		}
		// not feasible when there is gr8 than 1 len object
		//substring := col.Models[1:]
		//startIndex := strings.Index(substring, "{")
		//innerJSON := substring[startIndex:]
		//innerJSON = innerJSON[:len(innerJSON)-1]

		var innerJSON string
		var result map[string]json.RawMessage
		err := json.Unmarshal([]byte(col.Models), &result)

		if err != nil {
			log.Fatal(err)
		}

		for _, value := range result {
			//fmt.Printf("Key -> %s, Value -> %s\n", key, string(value))
			val := string(value)
			if len(innerJSON) < len(val) {
				innerJSON = val
			}
		}

		var models types.AnkiModel

		err = json.Unmarshal([]byte(innerJSON), &models)
		if err != nil {
			log.Fatal(err)
		}
		return models
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	//fileContent, err := os.ReadFile(filePath)
	//if err != nil {
	//	panic(err)
	//}
	return types.AnkiModel{}
}

func rename(destPath string) {
	var p = filepath.Join(destPath, "media")
	err := os.Rename(p, p+".json")
	if err != nil {
		panic(err)
	}
	mediaFile := p + ".json"
	fileContent, err := os.ReadFile(mediaFile)
	if err != nil {
		panic(err)
	}
	// Unmarshal JSON data into a map
	var dataMap map[string]string
	if err := json.Unmarshal(fileContent, &dataMap); err != nil {
		panic(err)
	}
	// Convert map to slice of structs
	var myDataSlice []types.AnkiMedia
	for id, file := range dataMap {
		myDataSlice = append(myDataSlice, types.AnkiMedia{
			ID:   id,
			File: file,
		})
	}
	for _, data := range myDataSlice {
		// Assume file path (current directory + ID)
		oldPath := filepath.Join(destPath, data.ID)

		// Check if the file exists
		if _, err := os.Stat(oldPath); err == nil {
			// Construct the new path using the File member
			newPath := filepath.Join(destPath, data.File)

			// Rename the file
			err := os.Rename(oldPath, newPath)
			if err != nil {
				panic(err)
			}
			//fmt.Printf("Renamed %s to %s\n", oldPath, newPath)
		} else if os.IsNotExist(err) {
			fmt.Printf("File %s does not exist, skipping...\n", oldPath)
		} else {
			panic(err)
		}
	}
}

func unzip(f *zip.File, destPath string) bool {
	rc, err := f.Open()
	if err != nil {
		fmt.Println(err)
		return false
	}
	filePath := filepath.Join(destPath, f.Name)

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		fmt.Println(err)
		err := rc.Close()
		if err != nil {
			return false
		}
		return false
	}

	_, err = io.Copy(file, rc)
	err = file.Close()
	if err != nil {
		return false
	}
	err = rc.Close()
	if err != nil {
		return false
	}
	return true
}

func ReadAPKGFile(apkgFilePath string) string {
	apkgData, err := os.ReadFile(apkgFilePath)
	if err != nil {
		log.Fatal(err)
	}
	// Open the zip archive
	zipReader, err := zip.NewReader(bytes.NewReader(apkgData), int64(len(apkgData)))
	if err != nil {
		log.Fatal(err)
	}
	zipDir := filepath.Dir(apkgFilePath)
	folderName := strings.TrimSuffix(filepath.Base(apkgFilePath), filepath.Ext(apkgFilePath))
	destPath := filepath.Join(zipDir, folderName)
	err = os.MkdirAll(destPath, 0755)
	if err != nil {
		return ""
	}

	// Iterate through the files in the zip archive
	for _, f := range zipReader.File {
		unzip(f, destPath)
	}

	// Make human readable
	rename(destPath)
	return destPath
}
