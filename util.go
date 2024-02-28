package main

import (
	"log"
	"os"
)

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
