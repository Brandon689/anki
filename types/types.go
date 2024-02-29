package types

// Anki Data Abbreviation.

type AnkiDeck struct {
	Name         string
	TemplateName string
	Cards        []AnkiCard
	CSS          string
	HTMLFormats  []AnkiHTMLFormat
}

type AnkiCard struct {
	Fields []AnkiFieldKeyValue
}

type AnkiFieldKeyValue struct {
	Key   string
	Value string
	Font  string
	Order int
}

type AnkiHTMLFormat struct {
	Name                        string
	Order                       int
	QuestionFormatHTMLTemplate  string
	AnswerFormatHTMLTemplate    string
	QuestionBFormatHTMLTemplate string
	AnswerBFormatHTMLTemplate   string
}

// end

type Note struct {
	ID    int
	GUID  string
	MID   int
	Mod   int
	USN   int
	Tags  string
	Flds  string
	Sfld  string
	Csum  int
	Flags int
	Data  string
}

type Col struct {
	ID     int    `json:"id"`
	Crt    int    `json:"crt"`
	Mod    int    `json:"mod"`
	Scm    int    `json:"scm"`
	Ver    int    `json:"ver"`
	Dty    int    `json:"dty"`
	Usn    int    `json:"usn"`
	Ls     int    `json:"ls"`
	Conf   string `json:"conf"`
	Models string `json:"models"`
	Decks  string `json:"decks"`
	Dconf  string `json:"dconf"`
	Tags   string `json:"tags"`
}

type AnkiMedia struct {
	ID   string `json:"id"`
	File string `json:"file"`
}

type AnkiModel struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Type  int    `json:"type"`
	Mod   int    `json:"mod"`
	Usn   int    `json:"usn"`
	Sortf int    `json:"sortf"`
	//Did   string `json:"did"`
	Tmpls []struct {
		Name  string      `json:"name"`
		Ord   int         `json:"ord"`
		Qfmt  string      `json:"qfmt"`
		Afmt  string      `json:"afmt"`
		Bqfmt string      `json:"bqfmt"`
		Bafmt string      `json:"bafmt"`
		Did   interface{} `json:"did"`
		Bfont string      `json:"bfont"`
		Bsize int         `json:"bsize"`
	} `json:"tmpls"`
	Flds []struct {
		Name                               string `json:"name"`
		Ord                                int    `json:"ord"`
		Sticky                             bool   `json:"sticky"`
		Rtl                                bool   `json:"rtl"`
		Font                               string `json:"font"`
		Size                               int    `json:"size"`
		Description                        string `json:"description"`
		PlainText                          bool   `json:"plainText"`
		Collapsed                          bool   `json:"collapsed"`
		ExcludeFromSearch                  bool   `json:"excludeFromSearch"`
		CollapsibleFieldsCollapseByDefault bool   `json:"collapsibleFieldsCollapseByDefault,omitempty"`
	} `json:"flds"`
	CSS       string          `json:"css"`
	LatexPre  string          `json:"latexPre"`
	LatexPost string          `json:"latexPost"`
	Latexsvg  bool            `json:"latexsvg"`
	Req       [][]interface{} `json:"req"`
}
