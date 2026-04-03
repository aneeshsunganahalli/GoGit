package internal

import (
	"encoding/json"
	"fmt"

	"os"
	"path/filepath"
	"strings"
)

// const indexPath = ".gogit/index.json"

// Loads the JSON from the .gogit/index.json file
func LoadIndex(indexPath string) map[string]IndexEntry {

	data, err := os.ReadFile(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]IndexEntry)
		}
		fmt.Println("Error loading the index from index.json: ", err)
		return make(map[string]IndexEntry)
	}

	index := make(map[string]IndexEntry)

	if err := json.Unmarshal(data, &index); err != nil {
		fmt.Println("Failed to unmarshal the json:", err)
		return make(map[string]IndexEntry)
	}

	// prettyJSON, err := json.MarshalIndent(index, "", "  ")

	// fmt.Println(string(prettyJSON))
	return index
}

// Writes the JSON to .gogit/index.json
func writeIndex(indexPath string, index map[string]IndexEntry) error {

	// ensure parent dir exists
	if err := os.MkdirAll(filepath.Dir(indexPath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(index, "", "")
	if err != nil {
		fmt.Println("Error writing index into index.json")
	}

	return os.WriteFile(indexPath, data, 0644)
}


func updateIndex(targetPath string, index map[string]IndexEntry, root *TrieNode) (map[string]bool, error) {
	seenFiles := make(map[string]bool)

	err := filepath.WalkDir(targetPath, func(path string, d os.DirEntry, err error) error {

		if err != nil || d.IsDir() {
			return err
		}
		cleanedPath := filepath.ToSlash(path)

		if cleanedPath == ".gogit" || strings.HasPrefix(cleanedPath, ".gogit/") {
			return nil
		}

		seenFiles[cleanedPath] = true
		info, _ := os.Stat(path)

		existingEntry, exists := index[cleanedPath]

		if exists && info.Size() == existingEntry.Size && info.ModTime().Unix() == existingEntry.Mtime {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		newHash, err := writeObject("blob", info.Size(), file)
		if err != nil {
			return err
		}

		mode := 100644
		if info.Mode()&0111 != 0 {
			mode = 100755
		}

		newEntry := IndexEntry{
			Filename: filepath.Base(path),
			Size:     info.Size(),
			Mtime:    info.ModTime().Unix(),
			Hash:     newHash,
			Mode:     mode,
		}

		index[cleanedPath] = newEntry
		root.MarkPath(cleanedPath, newEntry)

		return nil
	})
	return seenFiles, err
}



