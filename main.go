package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	apkgFilePath := "./apkg/Kanji_in_Context_Revised_Edition_2024_Edit.apkg"

	dir := readAPKGFile(apkgFilePath)

	deck := ReadConvert(dir)
	WriteJson(deck, fmt.Sprintf("./apkg/%s.json", deck.Name))

	//folder := strings.TrimSuffix(apkgFilePath, filepath.Ext(apkgFilePath))
	//dbFile := filepath.Join(folder, "collection.anki21")
	//form := sqlColTable(dbFile)
	//cards := sqlNotesTable(dbFile)
	//
	//deck := convert(form, cards)
	//WriteJson(deck, fmt.Sprintf("./apkg/%s.json", deck.Name))
	//html := renderHTMLCard(deck, 0, false, false, 0)
	//WriteFile("./apkg/flash.html", html)
}
