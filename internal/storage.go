package internal


import (
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Checks if object exists in storage, but I am programming above function to be in storage when we encounter a hash, so this is redudant
func ObjectExistsInStorage(hash string) bool {

	if len(hash) < 40 { // Git SHA-1 is 40 chars
		return false
	}

	dir, file := hash[:2], hash[2:]

	objectPath := filepath.Join(".gogit/objects", dir, file)
	// fmt.Println(objectPath)

	info, err := os.Stat(objectPath)
	if err == nil && !info.IsDir() {
		return true
	}

	return false
}




func writeObject(objectType string, size int64, r io.Reader) (string, error) {

	tempDir := filepath.Join(".gogit", "tmp")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", err
	}
	tempFile, err := os.CreateTemp(tempDir, "gogit-obj-*")
	if err != nil {
		return "", err
	}

	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	hasher := sha1.New()
	header := fmt.Sprintf("%s %d\x00", objectType, size)

	zlibWriter := zlib.NewWriter(tempFile)
	multiWriter := io.MultiWriter(hasher, zlibWriter)

	multiWriter.Write([]byte(header))

	if _, err := io.Copy(multiWriter, r); err != nil {
		return "", err
	}

	zlibWriter.Close()

	hashStr := fmt.Sprintf("%x", hasher.Sum(nil))

	dir, file := hashStr[:2], hashStr[2:]
	finalDir := filepath.Join(".gogit", "objects", dir)
	finalPath := filepath.Join(finalDir, file)

	if _, err := os.Stat(finalPath); err == nil {
		return hashStr, nil
	}

	if err := os.MkdirAll(finalDir, 0755); err != nil {
		return "", err
	}

	if err = os.Rename(tempFile.Name(), finalPath); err != nil {
		fmt.Println("File already exists")
	}

	fmt.Println(hashStr)
	return hashStr, err
}
