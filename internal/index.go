package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

	prettyJSON, err := json.MarshalIndent(index, "", "  ")

	fmt.Println(string(prettyJSON))
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

		seenFiles[path] = true
		info, _ := os.Stat(path)

		mtime, size := index[path].Mtime, index[path].Size

		if info.Size() != size || info.ModTime().Unix() != mtime {

			content, _ := os.ReadFile(path)
			newHash := GenerateHash("blob", string(content))

			mode := 100644
			if info.Mode()&0111 != 0 {
				mode = 100755
			}

			existingEntry, exists := index[path]


			if !exists {

				WriteObject("blob", string(content))

				fmt.Println("Entry not found in index, so we create one entry.")
				index[path] = IndexEntry{
					Filename: filepath.Base(path),
					Size:     info.Size(),
					Mtime:    info.ModTime().Unix(),
					Hash:     newHash,
					Mode:     mode,
				}
			} else {

				WriteObject("blob", string(content))

				if existingEntry.Hash != newHash {

					entry := index[path]
					entry.Hash = newHash
					index[path] = entry
				}
			}
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

	if err != nil {
		panic(err)
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
	dir, file := hash[:2], hash[2:]

	objectPath := filepath.Join(objectFolder, dir, file)
	fmt.Println(objectPath)

	info, err := os.Stat(objectPath)
	if err == nil && !info.IsDir() {
		return true
	}

	return false
}
