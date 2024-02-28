package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	apkgFilePath := "./apkg/Japanese_course_based_on_Tae_Kims_grammar_guide__anime.apkg"

	dir := readAPKGFile(apkgFilePath)

	deck := ReadConvert(dir)
	WriteJson(deck, fmt.Sprintf("./apkg/%s.json", deck.Name))

	//folder := strings.TrimSuffix(apkgFilePath, filepath.Ext(apkgFilePath))
	//dbFile := filepath.Join(folder, "collection.anki21")
	//form := sqlColTable(dbFile)
	//cards := sqlNotesTable(dbFile)
	//deck := convert(form, cards)
	//WriteJson(deck, fmt.Sprintf("./apkg/%s.json", deck.Name))
	html := renderHTMLCard(deck, 6, true, false, 1)
	WriteFile("./apkg/"+deck.Name+".html", html)
}
