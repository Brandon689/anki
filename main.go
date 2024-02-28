package main

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func sqlNotesTable(filePath string) []Note {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		log.Fatal(err)
	}

	// Query data
	rows, err := db.Query("SELECT * FROM notes")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var notes []Note
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

func sqlColTable(filePath string) AnkiModel {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		log.Fatal(err)
	}

	// Query data
	rows, err := db.Query("SELECT * FROM col LIMIT 1")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate through the result set
	for rows.Next() {
		var col Col
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
		substring := col.Models[1:]
		startIndex := strings.Index(substring, "{")
		innerJSON := substring[startIndex:]
		innerJSON = innerJSON[:len(innerJSON)-1]

		var models AnkiModel

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
	return AnkiModel{}
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
	var myDataSlice []AnkiMedia
	for id, file := range dataMap {
		myDataSlice = append(myDataSlice, AnkiMedia{
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
			fmt.Printf("Renamed %s to %s\n", oldPath, newPath)
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
		rc.Close()
		return false
	}

	_, err = io.Copy(file, rc)
	file.Close()
	rc.Close()

	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func Html(model AnkiModel, notes []Note) {
	css := model.CSS
	var str string

	str = model.Tmpls[0].Afmt
	str = strings.Replace(str, "furigana:", "", -1)
	for i := 0; i < len(model.Flds); i++ {
		fmt.Println(model.Flds[i].Name + ":")

		fields := strings.Split(notes[i].Flds, "\u001F")
		fmt.Println(fields[i] + "\n")

		str = strings.Replace(str, "{{"+model.Flds[i].Name+"}}", fields[i], -1)
	}

	str += "<style>\n" + css + "\n</style>"

	file, err := os.Create("./apkg/flashcard.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	_, err = file.WriteString(str)
	if err != nil {
		log.Fatal(err)
	}
}

func readAPKGFile(apkgFilePath string) {
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
	os.MkdirAll(destPath, 0755)

	// Iterate through the files in the zip archive
	for _, f := range zipReader.File {
		unzip(f, destPath)
	}

	// Make human readable
	rename(destPath)
}

func main() {

	apkgFilePath := "./apkg/Japanese_Core_2000_2k_-_Sorted_w_Audio.apkg"

	readAPKGFile(apkgFilePath)

	folder := strings.TrimSuffix(apkgFilePath, filepath.Ext(apkgFilePath))
	dbFile := filepath.Join(folder, "collection.anki2")
	form := sqlColTable(dbFile)
	cards := sqlNotesTable(dbFile)

	for i := 0; i < len(form.Flds); i++ {
		fmt.Println(form.Flds[i].Name + ":")

		fields := strings.Split(cards[i].Flds, "\u001F")
		fmt.Println(fields[i] + "\n")
	}

	Html(form, cards)
}
