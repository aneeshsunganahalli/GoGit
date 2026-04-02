package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// const indexPath = ".gogit/index.json"

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

func RecursiveWalk(start string) {

	// indexPath := ""
	index := LoadIndex(".gogit/index.json")
	seenFiles := make(map[string]bool)

	err := filepath.WalkDir(start, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}


		if !d.IsDir() {
			// fileName := filepath.Base(path)
			seenFiles[path] = true

			info, err := os.Stat(path)
			if err != nil {
				fmt.Println("Error with getting stats")
			}

			mtime, size := index[path].Mtime, index[path].Size
			if info.Size() != size || info.ModTime().Unix() != mtime {

				content, err := os.ReadFile(path)
				if err != nil {
					fmt.Println("Error reading content from file")
				}

				newHash := GenerateHash("blob", string(content))

				existingEntry, ok := index[path]
				if !ok {
					fmt.Println("Entry not found in index, so we create one entry.")
					index[path] = IndexEntry{
						Filename: filepath.Base(path),
						Size: info.Size(),
						Mtime: info.ModTime().Unix(),
						Hash: newHash,
						Mode: 120000,
					}
				}

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
				delete(index, path)
			}
		}

		writeIndex(".gogit/index.json", index)

	if err != nil {
		panic(err)
	}

}

func writeIndex(indexPath string, index map[string]IndexEntry) error {
	data, err := json.MarshalIndent(index, "", "")
	if err != nil {
		fmt.Println("Error writing index into index.json")
	}

	return os.WriteFile(indexPath, data, 0755)
}

