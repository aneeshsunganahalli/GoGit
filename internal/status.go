package internal

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ResolvePointer(input string) (string, error) {

	branchPath := fmt.Sprintf(".gogit/refs/heads/%s", input)
	data, err := os.ReadFile(branchPath)

	if err == nil {
		return strings.TrimSpace(string(data)), nil
	}

	// If it's not a branch, we check if it's valid hash
	isHex := regexp.MustCompile(`^[0-9a-fA-F]{7,40}$`).MatchString(input)
	if isHex {
		return input, nil
	}

	return "", fmt.Errorf("ref %s not found", input)
}

// Checks for unstaged changes, checking index against working directory
func GetUnstagedChanges(index map[string]IndexEntry) []string {
	var unstaged []string

	filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		cleanedPath := filepath.ToSlash(path)
		if cleanedPath == ".gogit" || strings.HasPrefix(cleanedPath, ".gogit/") || cleanedPath == ".git" || strings.HasPrefix(cleanedPath, ".git") {
			return nil
		}

		info, _ := os.Stat(path)
		entry, exists := index[cleanedPath]

		if !exists {
			unstaged = append(unstaged, fmt.Sprint("Untracked: "+cleanedPath))
		} else if info.Size() != entry.Size && info.ModTime().Unix() != entry.Mtime {
			unstaged = append(unstaged, fmt.Sprint("Modified: "+cleanedPath))
		}

		return nil
	})

	return unstaged
}

func GetStatus(root *TrieNode) (bool, string) {

	rootHash := hex.EncodeToString(root.WriteMerkleTree())
	

	parentTreeHash, err := GetHeadTreeHash()


	if err != nil {
		return true, rootHash
	}

	isDifferent := rootHash == parentTreeHash

	return isDifferent, rootHash

}

func ShowStatus(index map[string]IndexEntry, root *TrieNode) {
	unstaged := GetUnstagedChanges(index)

	isStaged, _ := GetStatus(root)

	if len(unstaged) == 0 && !isStaged {
		fmt.Println("nothing to commit, working tree clean")
		return
	}

	if isStaged {
		fmt.Println("Changes to be committed")
	}

	if len(unstaged) > 0 {
		fmt.Println("Changes not staged for commit:")
		for _, msg := range unstaged {
			fmt.Println("  ", msg)
		}
	}
}
