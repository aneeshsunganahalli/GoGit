package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

const objectFolder = ".gogit/objects/"

// Writes the object into the .gogit/objects/ folder in the format: object/sd/j8k4... for storage
func WriteObject(objectType string, content string) string {
	
	store := GenerateStore(objectType, content)

	hashStr := GenerateHash(objectType, content)
	fmt.Println(hashStr)
	compressedContent, err := ZlibCompresser(store)

	// Directory Creation
	dir := hashStr[:2]
	file := hashStr[2:]

	path := filepath.Join(objectFolder, dir)

	err = os.MkdirAll(path, 0755)
	if err != nil {
		panic(fmt.Sprintf("Failed to create directory at %s: %v", path, err))
	}

	// File Creation
	fileName := path + "/" + file
	objFile, err := os.Create(fileName)
	if err != nil {
		panic(fmt.Sprintf("Failed to write object at %s: %v", path, err))
	}

	defer objFile.Close()

	os.WriteFile(fileName, compressedContent, 0755)

	return hashStr // You'll need this to keep track of what you have saved
}

func Add() {

}

