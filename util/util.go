package util

import (
	"encoding/json"
	"fmt"
	"github.com/Brandon689/anki/types"
	"log"
	"os"
)

func WriteJson(deck types.AnkiDeck, outputFile string) {
	jsonData, err := json.MarshalIndent(deck, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}
	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

func WriteFile(destPath string, contents string) {
	file, err := os.Create(destPath)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	_, err = file.WriteString(contents)
	if err != nil {
		log.Fatal(err)
	}
}
