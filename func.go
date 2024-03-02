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

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func ReadConvert(apkgFolder string) (AnkiDeck, error) {
	dbFile := filepath.Join(apkgFolder, "collection.anki21")
	if !fileExists(dbFile) {
		dbFile = filepath.Join(apkgFolder, "collection.anki2")
	}
	form, err := sqlColTable(dbFile)
	if err != nil {
		panic(err)
	}
	cards, err := sqlNotesTable(dbFile)
	if err != nil {
		panic(err)
	}
	deck, err := Convert(form, cards)
	if err != nil {
		panic(err)
	}
	deck.Name = filepath.Base(apkgFolder)
	return deck, nil
}

func Convert(model AnkiModel, notes []Note) (AnkiDeck, error) {
	var ad AnkiDeck
	ad.TemplateName = model.Name
	ad.CSS = model.CSS
	ad.Cards = []AnkiCard{}
	ad.HTMLFormats = []AnkiHTMLFormat{}
	if len(notes) == 0 {
		log.Fatal("no notes")
	}
	for i := 0; i < len(notes); i++ {
		fields := strings.Split(notes[i].Flds, "\u001F")

		ct := AnkiCard{}
		ct.Fields = []AnkiFieldKeyValue{}

		for j := 0; j < len(fields); j++ {
			a := AnkiFieldKeyValue{}
			a.Key = model.Flds[j].Name
			a.Value = fields[j]
			a.Font = model.Flds[j].Font
			a.Order = model.Flds[j].Ord
			ct.Fields = append(ct.Fields, a)
		}
		ad.Cards = append(ad.Cards, ct)
	}
	for i := 0; i < len(model.Tmpls); i++ {
		ahf := AnkiHTMLFormat{}
		ahf.Name = model.Tmpls[i].Name
		ahf.Order = model.Tmpls[i].Ord
		ahf.QuestionFormatHTMLTemplate = model.Tmpls[i].Qfmt
		ahf.AnswerFormatHTMLTemplate = model.Tmpls[i].Afmt
		ahf.QuestionBFormatHTMLTemplate = model.Tmpls[i].Bqfmt
		ahf.AnswerBFormatHTMLTemplate = model.Tmpls[i].Bafmt
		ad.HTMLFormats = append(ad.HTMLFormats, ahf)
	}
	return ad, nil
}

func WriteJson(deck AnkiDeck, outputFile string) error {
	jsonData, err := json.MarshalIndent(deck, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		panic(err)
	}
	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		panic(err)
	}
	return nil
}

func WriteFile(destPath string, contents string) {
	file, err := os.Create(destPath)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)
	_, err = file.WriteString(contents)
	if err != nil {
		log.Fatal(err)
	}
}

func sqlNotesTable(filePath string) ([]Note, error) {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT * FROM notes")
	if err != nil {
		log.Fatal(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)
	var notes []Note
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
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return notes, err
}

func sqlColTable(filePath string) (AnkiModel, error) {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT * FROM col LIMIT 1")
	if err != nil {
		log.Fatal(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
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
			val := string(value)
			if len(innerJSON) < len(val) {
				innerJSON = val
			}
		}
		var models AnkiModel
		err = json.Unmarshal([]byte(innerJSON), &models)
		if err != nil {
			log.Fatal(err)
		}
		return models, err
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return AnkiModel{}, err
}

func Rename(destPath string) error {
	var d = filepath.Join(destPath, "collection.anki2")
	err := os.Rename(d, destPath+"/deck.db")
	if err != nil {
		panic(err)
	}
	var p = filepath.Join(destPath, "media")
	err = os.Rename(p, p+".json")
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
	for id, file := range dataMap {
		oldPath := filepath.Join(destPath, id)

		// Check if the file exists
		if _, err := os.Stat(oldPath); err == nil {
			newPath := filepath.Join(destPath, file)

			// Rename the file
			err := os.Rename(oldPath, newPath)
			if err != nil {
				panic(err)
			}
		} else if os.IsNotExist(err) {
			fmt.Printf("File %s does not exist, skipping...\n", oldPath)
		} else {
			panic(err)
		}
	}
	return nil
}

func extractFile(f *zip.File, destPath string) error {
	rc, err := f.Open()
	if err != nil {
		panic(err)
	}
	filePath := filepath.Join(destPath, f.Name)

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		fmt.Println(err)
		err := rc.Close()
		panic(err)
	}
	_, err = io.Copy(file, rc)
	err = file.Close()
	if err != nil {
		panic(err)
	}
	err = rc.Close()
	if err != nil {
		panic(err)
	}
	return nil
}

func ExtractAPKGFile(apkgFilePath string, destDir string) (string, error) {
	apkgData, err := os.ReadFile(apkgFilePath)
	if err != nil {
		log.Fatal(err)
	}
	zipReader, err := zip.NewReader(bytes.NewReader(apkgData), int64(len(apkgData)))
	if err != nil {
		log.Fatal(err)
	}
	zipDir := filepath.Dir(apkgFilePath)
	folderName := strings.TrimSuffix(filepath.Base(apkgFilePath), filepath.Ext(apkgFilePath))
	destPath := destDir
	if destPath == "" {
		destPath = filepath.Join(zipDir, folderName)
	}
	err = os.MkdirAll(destPath, 0755)
	for _, f := range zipReader.File {
		extractFile(f, destPath)
	}
	return destPath, err
}

func RenderHTMLCard(deck AnkiDeck, index int, question bool, BSide bool, formatIndex int) string {
	var html string
	if question {
		if BSide {
			html = deck.HTMLFormats[formatIndex].QuestionBFormatHTMLTemplate
		} else {
			html = deck.HTMLFormats[formatIndex].QuestionFormatHTMLTemplate
		}
	} else {
		if BSide {
			html = deck.HTMLFormats[formatIndex].AnswerBFormatHTMLTemplate
		} else {
			html = deck.HTMLFormats[formatIndex].AnswerFormatHTMLTemplate
		}
	}
	html = strings.Replace(html, "furigana:", "", -1)

	card := deck.Cards[index]
	for i := 0; i < len(card.Fields); i++ {
		if card.Fields[i].Key == "Image" {
			html = strings.Replace(html, "{{"+card.Fields[i].Key+"}}", card.Fields[i].Value, -1)
			//html = strings.Replace(html, "{{"+card.Fields[i].Key+"}}", deck.Name+"/"+card.Fields[i].Value, -1)
		} else {
			html = strings.Replace(html, "{{"+card.Fields[i].Key+"}}", card.Fields[i].Value, -1)
		}
	}
	html = strings.Replace(html, "<img src=\"", "<img src=\""+deck.Name+"\\", 1)
	html += "<style>\n" + deck.CSS + "\n</style>"
	return html
}
