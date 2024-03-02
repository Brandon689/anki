# Anki Go Package

**Anki Go Package** is a Go library for reading Anki .apkg archives and converting them into a structured JSON format.

## Installation

To use this package in your Go project, you can use the `go get` command:

```bash
go get github.com/Brandon689/anki
```

### Usage

Import the anki package in your Go code and use its functions to read Anki .apkg archives.

```go

package main

import (
	"fmt"
	"github.com/Brandon689/anki"
)

func main() {

  srcFile := "./testsuite/Hiragana_Test_Deck.apkg"
	destFolder := "./testsuite/Hiragana_Test_Deck"

	_, err := ExtractAPKGFile(srcFile, destFolder)
	if err != nil {
		panic(err)
	}
	colTable, _ := sqlColTable("./testsuite/test_collection.anki2")
	notesTable, _ := sqlNotesTable("./testsuite/test_collection.anki2")
	r, _ := Convert(colTable, notesTable)
  fmt.Println(r.Name)
}
```
### Features

Read Anki .apkg archives.
Convert Anki card data to JSON format.

### Contributing

If you'd like to contribute to this project, feel free to submit issues or pull requests. Please follow the contribution guidelines.

### License

This project is licensed under the MIT License.
