package anki

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func WriteJson(deck AnkiDeck, outputFile string) {
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
