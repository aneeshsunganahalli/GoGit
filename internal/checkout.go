package internal

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func GetTreeHashFromCommit(commitHash string)(string, error) {

	_ , content, err := readObject(commitHash)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "tree "){
			return strings.TrimPrefix(line, "tree "), nil
		}
	}

	return "", fmt.Errorf("Commit Object not found")
}

func RestoreFromTree(treeHash string, currentPath string, index map[string]IndexEntry) error {

	_, data, _ := readObject(treeHash)

	for len(data) > 0 {
		name, node, n, _ := ParseTreeEntry(data)

		data = data[n:]

		path := filepath.Join(currentPath, name)
		hashHex := hex.EncodeToString(node.Hash)

		if !node.IsFile {
			os.MkdirAll(path, 0755)
			if err := RestoreFromTree(hashHex, path, index); err != nil {
				return err
			}

		} else {
			_, content, _ := readObject(hashHex)
			os.WriteFile(path, content, os.FileMode(node.Mode))

			cleanedPath := filepath.ToSlash(path)

			index[cleanedPath] = IndexEntry{
				Filename: name,
				Hash:     hashHex,
				Size:     int64(len(content)),
				Mode:     node.Mode,
				Mtime:    time.Now().Unix(),
			}
		}
	}

	return nil
}

func updateHEAD(input string, hash string) {
    
    isHex := regexp.MustCompile(`^[0-9a-fA-F]{7,40}$`).MatchString(input)

    var content string
    if isHex {
        
        content = hash + "\n"
    } else {

        // Format must be "refs:refs/heads/<name>"
        content = fmt.Sprintf("refs:refs/heads/%s\n", input)
    }

    err := os.WriteFile(".gogit/HEAD", []byte(content), 0644)
    if err != nil {
        fmt.Printf("Error updating HEAD: %v\n", err)
    }
}

