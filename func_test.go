package main

import (
	"testing"
)

func TestRename(t *testing.T) {
	srcFolder := "./testsuite/Hiragana_Test_Deck"
	destFolder := "./testsuite/Hiragana_Test_Deck_Temp"

	if err := copyDir(srcFolder, destFolder); err != nil {
		t.Error(err)
	}
	err := Rename("./testsuite/Hiragana_Test_Deck_Temp")
	if err != nil {
		t.Errorf("Expected %d", err)
	}
}

func TestUnZip(t *testing.T) {
	srcFile := "./testsuite/Hiragana_Test_Deck.apkg"
	destFolder := "./testsuite/Hiragana_Test_Deck"

	_, err := ExtractAPKGFile(srcFile, destFolder)
	if err != nil {
		t.Errorf("Expected %d", err)
	}
}

func TestColTable(t *testing.T) {
	table, err := sqlColTable("./testsuite/test_collection.anki2")
	if err != nil {
		t.Errorf("Expected %d", err)
	}
	if len(table.Flds) == 0 {
		t.Errorf("no results")
	}
}

func TestNotesTable(t *testing.T) {
	table, err := sqlNotesTable("./testsuite/test_collection.anki2")
	if err != nil {
		t.Errorf("Expected %d", err)
	}
	if len(table) == 0 {
		t.Errorf("no results")
	}
}

func TestConvert(t *testing.T) {
	colTable, _ := sqlColTable("./testsuite/test_collection.anki2")
	notesTable, _ := sqlNotesTable("./testsuite/test_collection.anki2")
	r, _ := Convert(colTable, notesTable)
	if len(r.Cards) == 0 {
		t.Errorf("cards len 0")
	}
}

func TestReadConvert(t *testing.T) {
	r, _ := ReadConvert("C:\\Users\\Brandon\\Videos\\ni\\Hiragana_Test_Deck")
	if r.Name == "" {
		t.Errorf("no deck name")
	}
	if len(r.Cards) == 0 {
		t.Errorf("cards len 0")
	}
}

func TestJsonSerialization(t *testing.T) {
	colTable, _ := sqlColTable("./testsuite/test_collection.anki2")
	notesTable, _ := sqlNotesTable("./testsuite/test_collection.anki2")
	r, _ := Convert(colTable, notesTable)
	err := WriteJson(r, "./testsuite/test_write_deck.json")
	if err != nil {
		t.Error(err)
	}
}
