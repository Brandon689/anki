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
	apkgFilePath := "./path/to/your/file.apkg"
	deckData, err := anki.ReadAPKGFile(apkgFilePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	
	anki.WriteJson(deckData, fmt.Sprintf("./%s.json", deckData.Name))
	
	html := anki.RenderHTMLCard(deck, 6, true, false, 1)
	anki.WriteFile("./apkg/"+deck.Name+".html", html)
	
}
```
### Features

    Read Anki .apkg archives.
    Convert Anki card data to JSON format.

### Contributing

If you'd like to contribute to this project, feel free to submit issues or pull requests. Please follow the contribution guidelines.

### License

This project is licensed under the MIT License.