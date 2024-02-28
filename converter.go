package main

import (
	"strings"
)

func Convert(model AnkiModel, notes []Note) AnkiDeck {
	var ad AnkiDeck
	ad.Name = model.Name
	ad.CSS = model.CSS
	ad.Cards = []AnkiCard{}
	ad.HTMLFormats = []AnkiHTMLFormat{}

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
	return ad
}
