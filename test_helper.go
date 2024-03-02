package main

import (
	"io"
	"os"
	"path/filepath"
)

func copyDir(src, dest string) error {
	// Create the destination directory
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return err
	}

	// Open the source directory
	srcDir, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcDir.Close()

	// Read the contents of the source directory
	files, err := srcDir.Readdir(-1)
	if err != nil {
		return err
	}

	// Copy each file or subdirectory
	for _, file := range files {
		srcPath := filepath.Join(src, file.Name())
		destPath := filepath.Join(dest, file.Name())

		if file.IsDir() {
			// Recursive call for subdirectories
			if err := copyDir(srcPath, destPath); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := copyFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}
