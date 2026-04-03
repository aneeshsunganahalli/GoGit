package internal

import (
	// "compress/zlib"
	// "crypto/sha1"
	"fmt"
	// "io"
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

	err = os.MkdirAll(path, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to create directory at %s: %v", path, err))
	}

	// File Creation
	fileName := filepath.Join(path, file)
	

	os.WriteFile(fileName, compressedContent, 0644)

	return hashStr // You'll need this to keep track of what you have saved
}

// func writeObject(objectType string, size int64, r io.Reader) (string, error) {
// 	tempFile, err := os.CreateTemp("", "gogit-obj-*")
// 	if err != nil {
// 		return "", err
// 	}

// 	defer os.Remove(tempFile.Name())
// 	defer tempFile.Close()
	

// 	hasher := sha1.New()
// 	header := fmt.Sprintf("%s %d\x00", objectType, size)

// 	zlibWriter := zlib.NewWriter(tempFile)
// 	multiWriter := io.MultiWriter(hasher, zlibWriter)

// 	multiWriter.Write([]byte(header))


// 	if _ ,err := io.Copy(multiWriter, r); err != nil {
// 		return "", err
// 	}

// 	zlibWriter.Close()

// 	hashStr := fmt.Sprintf("%x", hasher.Sum(nil))

// 	dir, file := hashStr[:2], hashStr[2:]
// 	finalDir := filepath.Join(".gogit", "objects", dir)
// 	finalPath := filepath.Join(finalDir, file)

// 	if _, err := os.Stat(finalPath); err == nil {
// 		return hashStr, nil
// 	}

// 	if err := os.MkdirAll(finalDir, 0755); err != nil {
// 		return "", err
// 	}

// 	return hashStr, os.Rename(tempFile.Name(), finalPath)
// }

