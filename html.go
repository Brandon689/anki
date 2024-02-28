package main

import (
	"log"
	"os"
	"strings"
)

func renderHTMLCard(deck AnkiDeck, index int, question bool, BSide bool, formatIndex int) string {
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
		html = strings.Replace(html, "{{"+card.Fields[i].Key+"}}", card.Fields[i].Value, -1)
	}
	html += "<style>\n" + deck.CSS + "\n</style>"
	return html
}

func htmlQA(model AnkiModel, note Note, q bool) {
	css := model.CSS
	var str string

	if q == true {
		str = model.Tmpls[0].Qfmt
	} else {
		str = model.Tmpls[0].Afmt
	}
	str = strings.Replace(str, "furigana:", "", -1)
	for i := 0; i < len(model.Flds); i++ {
		fields := strings.Split(note.Flds, "\u001F")
		str = strings.Replace(str, "{{"+model.Flds[i].Name+"}}", fields[i], -1)
	}
	str += "<style>\n" + css + "\n</style>"
	file, err := os.Create("./apkg/flashcard.html")
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	_, err = file.WriteString(str)
	if err != nil {
		log.Fatal(err)
	}
}
