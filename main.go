package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"
	"strings"
)

func main() {

	apkgFilePath := "./apkg/Japanese_Core_2000_2k_-_Sorted_w_Audio.apkg"

	//
	//readAPKGFile(apkgFilePath)

	folder := strings.TrimSuffix(apkgFilePath, filepath.Ext(apkgFilePath))
	dbFile := filepath.Join(folder, "collection.anki2")
	form := sqlColTable(dbFile)
	cards := sqlNotesTable(dbFile)

	for i := 0; i < len(form.Flds); i++ {
		fmt.Println(form.Flds[i].Name + ":")

		fields := strings.Split(cards[i].Flds, "\u001F")
		fmt.Println(fields[i] + "\n")
	}

	deck := Convert(form, cards)
	html := renderHTMLCard(deck, 0, false, false, 0)
	WriteFile("./apkg/flash.html", html)
}
