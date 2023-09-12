package helper

import (
	"os"
	"path/filepath"
)

func OpenFileForAppend(filename string) (*os.File, error) {
	// Change this path to the directory where you want to save the file.
	// Make sure the directory exists and you have the necessary permissions.
	uploadDir := "./uploads"

	// Create the directory if it doesn't exist.
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, err
	}

	// Open the file for append mode (creates the file if it doesn't exist).
	destination, err := os.OpenFile(filepath.Join(uploadDir, filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return destination, nil
}

func createDestinationFile(filename string) (*os.File, error) {
	// Change this path to the directory where you want to save the file.
	// Make sure the directory exists and you have the necessary permissions.
	uploadDir := "./uploads/"

	// Create the directory if it doesn't exist.
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, err
	}

	// Create a new file in the specified directory.
	destination, err := os.Create(filepath.Join(uploadDir, filename))
	if err != nil {
		return nil, err
	}

	return destination, nil
}
