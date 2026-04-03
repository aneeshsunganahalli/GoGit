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
		fmt.Println("Error loading the index from index.json")
	}

	index := make(map[string]IndexEntry)

	err = json.Unmarshal(data, &index)
	if err != nil {
		fmt.Println("Failed to unmarshal the json")
	}

	// prettyJSON, err := json.MarshalIndent(index, "", "  ")

	// fmt.Println(string(prettyJSON))
	return index
}

func UpdateIndexFromPath(targetPath string) {

	indexPath := ".gogit/index.json"
	index := LoadIndex(indexPath)
	seenFiles := make(map[string]bool)

	err := filepath.WalkDir(targetPath, func(path string, d os.DirEntry, err error) error {

		if err != nil || d.IsDir() {
			return err
		}
		cleanedPath := filepath.ToSlash(path)
		seenFiles[cleanedPath] = true
		info, _ := os.Stat(path)

		existingEntry, exists := index[path]

		if exists && existingEntry.Size == info.Size() && existingEntry.Mtime == info.ModTime().Unix() {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Error reading content from %s", path)
			return nil
		}
		newHash := GenerateHash("blob", string(content))

		mode := 100644
		if info.Mode()&0111 != 0 {
			mode = 100755
		}

		if !ObjectExistsInStorage(newHash) {
			// fmt.Printf("Creating new object: %s\n", newHash)
			WriteObject("blob", string(content))
		}

		index[path] = IndexEntry{
			Filename: filepath.Base(path),
			Size:     info.Size(),
			Mtime:    info.ModTime().Unix(),
			Hash:     newHash,
			Mode:     mode,
		}

		return nil
	})

	for path := range index {
		if seenFiles[path] == false {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				delete(index, path)
			}
		}
	}

	writeIndex(".gogit/index.json", index)

	PrintTrie(BuildTrie(index), "")

	if err != nil {
		fmt.Println(err)
	}

}

// Writes the JSON to .gogit/index.json
func writeIndex(indexPath string, index map[string]IndexEntry) error {
	data, err := json.MarshalIndent(index, "", "")
	if err != nil {
		fmt.Println("Error writing index into index.json")
	}

	return os.WriteFile(indexPath, data, 0755)
}

// Checks if object exists in storage, but I am programming above function to be in storage when we encounter a hash, so this is redudant
func ObjectExistsInStorage(hash string) bool {

	if len(hash) < 40 { // Git SHA-1 is 40 chars
        return false
    }
		
	dir, file := hash[:2], hash[2:]

	objectPath := filepath.Join(objectFolder, dir, file)
	// fmt.Println(objectPath)

	info, err := os.Stat(objectPath)
	if err == nil && !info.IsDir() {
		return true
	}

	return false
}

// Builds a temp trie
func BuildTrie(index map[string]IndexEntry) *TrieNode {

	root := &TrieNode{
		Children: make(map[string]*TrieNode),
		Mode:     40000,
	}

	for path, entry := range index {
		parts := strings.Split(path, "/")
		current := root

		for idx, part := range parts {

			// It's a file
			if idx == len(parts)-1 {
				current.Children[part] = &TrieNode{
					Hash:   entry.Hash,
					Mode:   entry.Mode,
					IsFile: true,
				}
			} else { // It's a directory
				if _, exists := current.Children[part]; !exists {
					current.Children[part] = &TrieNode{
						Children: make(map[string]*TrieNode),
						IsFile:   false,
						Mode:     40000,
					}
				}
			}
			current = current.Children[part]
		}

	}

	return root
}


