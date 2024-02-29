package main

import (
	"fmt"
	"github.com/Brandon689/anki/logic"
	"github.com/Brandon689/anki/util"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	apkgFilePath := "./apkg/Japanese_course_based_on_Tae_Kims_grammar_guide__anime.apkg"

	dir := logic.ReadAPKGFile(apkgFilePath)

	deck := logic.ReadConvert(dir)
	util.WriteJson(deck, fmt.Sprintf("./apkg/%s.json", deck.Name))

	//folder := strings.TrimSuffix(apkgFilePath, filepath.Ext(apkgFilePath))
	//dbFile := filepath.Join(folder, "collection.anki21")
	//form := sqlColTable(dbFile)
	//cards := sqlNotesTable(dbFile)
	//deck := convert(form, cards)
	//WriteJson(deck, fmt.Sprintf("./apkg/%s.json", deck.Name))
	html := logic.RenderHTMLCard(deck, 6, true, false, 1)
	util.WriteFile("./apkg/"+deck.Name+".html", html)
}
