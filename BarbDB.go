package barbdb

import (
	"encoding/base64" // To encode and decode strings in the database
	"errors"          // To return custom errors
	"os"              // To perform CRUD operations on files
	"strings"         // To split and merge strings
)

// Storage struct represents the database.
type Storage struct {
	file os.File
	keys map[string]bool
}

// New opens a database at the given path.
func New(path string) (*Storage, error) {
	// Open the file
	file, fileError := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600)
	if fileError != nil {
		return nil, fileError
	}

	// Return the database struct
	return &Storage{
		file: *file,
		keys: make(map[string]bool),
	}, nil
}

// Helper function to read the file.
func (db Storage) readFile() ([]string, error) {
	// Read the file
	data, readError := os.ReadFile(db.file.Name())
	if readError != nil {
		return nil, readError
	}

	rows := strings.Split(string(data), "\n")

	for i := 0; i < len(rows); i++ {
		var k = strings.Split(rows[i], "=")[0]
		db.keys[k] = true
	}

	// Split the data into lines and return it
	return rows, nil
}

// Get returns the value of the given key.
func (db Storage) Get(key string) (string, error) {
	// Base64 encode the key
	encodedKey := base64.RawStdEncoding.EncodeToString([]byte(key))

	// Read the file
	fileContent, fileReadError := db.readFile()
	if fileReadError != nil {
		return "", fileReadError
	}

	// Checks in the db keys map, if absent returns error.
	if !db.keys[encodedKey] {
		return "", errors.New("key not found in the keys map")
	}

	// Loop over the lines in the file
	for i := 0; i < len(fileContent); i++ {
		splitString := strings.Split(fileContent[i], "=")
		if splitString[0] == encodedKey {
			// Decode the value and return it
			toReturn, base64DecodeError := base64.RawStdEncoding.DecodeString(splitString[1])
			if base64DecodeError != nil {
				return "", base64DecodeError
			}
			return string(toReturn), nil
		}
	}

	// Return an error if the key doesn't exist
	return "", errors.New("key not found")
}

// Set sets the value of the given key.
func (db Storage) Set(key string, value string) error {
	// Base64 encode the key and value
	encodedKey := base64.RawStdEncoding.EncodeToString([]byte(key))
	encodedValue := base64.RawStdEncoding.EncodeToString([]byte(value))

	// Read the file
	fileContent, fileReadError := db.readFile()
	if fileReadError != nil {
		return fileReadError
	}

	// Check if the key already exists
	for i := 0; i < len(fileContent); i++ {
		splitString := strings.Split(fileContent[i], "=")
		if splitString[0] == encodedKey {
			// Delete the key if it does
			deleteError := db.Delete(key)
			if deleteError != nil {
				return deleteError
			}
		}
	}

	// Write the key and value to the file
	_, fileWriteError := db.file.WriteString(encodedKey + "=" + encodedValue + "\n")
	if fileWriteError != nil {
		return fileWriteError
	}

	// adds key to db keys map.
	db.keys[encodedKey] = true

	return db.file.Sync()
}

// Delete deletes the given key from the database.
func (db Storage) Delete(key string) error {
	// Base64 encode the key
	encodedKey := base64.RawStdEncoding.EncodeToString([]byte(key))

	// deletes key from db keys map.
	delete(db.keys, encodedKey)

	// Read the file
	fileContent, fileReadError := db.readFile()
	if fileReadError != nil {
		return fileReadError
	}

	// Loop over the lines in the file
	for i := 0; i < len(fileContent); i++ {
		splitString := strings.Split(fileContent[i], "=")
		if splitString[0] == encodedKey {
			// Delete the key
			fileContent = append(fileContent[:i], fileContent[i+1:]...)
		}
	}

	// Remove the key and value from the file
	fileWriteError := os.WriteFile(db.file.Name(), []byte(strings.Join(fileContent, "\n")), 0600)
	if fileWriteError != nil {
		return fileWriteError
	}
	return db.file.Sync()
}

// Close closes the database.
func (db Storage) Close() error {
	return db.file.Close()
}
