package internal

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func GoGitInit(cmd *cobra.Command, args []string) {

	// Creates the .gogit directory
	dir := ".gogit"
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Println("failed to create .gogit:", err)
		return
	}

	subfolders := []string{
		"hooks",
		"refs",
		"info",
		"objects",
	}

	for _, subf := range subfolders {
		err := os.MkdirAll(filepath.Join(dir, subf), 0775)
		if err != nil {
			fmt.Println("Error creating subfolders in .gogit")
			return
		}
	}

	// Creates the HEAD file
	headPath := dir + "/HEAD"
	err := os.WriteFile(headPath, []byte("refs:refs/heads/main\n"), 0644)
	if err != nil {
		fmt.Println("Failed to write HEAD: ", err)
	}

	indexPath := filepath.Join(dir, "index.json")
	err = os.WriteFile(indexPath, []byte("{}\n"), 0644)
	if err != nil {
		fmt.Println("Failed to write HEAD: ", err)
	}

	GetAuthorDetails()

	fmt.Println("Initialized empty repository")
}

func GoGitAdd(targetPath string) {

	indexPath := ".gogit/index.json"
	index := LoadIndex(indexPath)

	seenFiles, err := updateIndex(targetPath, index)

	for path := range index {
		if seenFiles[path] == false {
			if _, err := os.Stat(filepath.FromSlash(path)); os.IsNotExist(err) {
				delete(index, path)
				// root.RemovePath(path)
			}
		}
	}

	writeIndex(".gogit/index.json", index)

	if err != nil {
		fmt.Println(err)
	}

}

func GoGitCommit(message string) {
	indexPath := ".gogit/index.json"
	index := LoadIndex(indexPath)

	root := &TrieNode{Children: make(map[string]*TrieNode), Mode: 40000, IsDirty: true}

	for path, entry := range index {
		root.LoadPath(path, entry)
	}

	// PrintTrie(root, "")
	rootHash := hex.EncodeToString(root.WriteMerkleTree())
	fmt.Println(rootHash)

	parentTreeHash, err := GetHeadTreeHash()
	fmt.Println(parentTreeHash)

	if err == nil && parentTreeHash != "" {
		if rootHash == parentTreeHash {
			fmt.Println("On branch main")
			fmt.Println("nothing to commit, working tree clean")
			return // <--- CRITICAL: Stop here!
		}
	}

	parentHash, err := GetHeadHash()

	commitHash := CreateAndStoreCommit(rootHash, parentHash, message)

	currentBranch := "main" // fallback
	headContent, err := os.ReadFile(".gogit/HEAD")
	if err == nil {
		headStr := strings.TrimSpace(string(headContent))
		
		if strings.HasPrefix(headStr, "refs:refs/heads/") {
			currentBranch = strings.TrimPrefix(headStr, "refs:refs/heads/")
		}
	}

	refDir := filepath.Join(".gogit", "refs", "heads")
	refPath := filepath.Join(refDir, currentBranch)

	err = os.MkdirAll(refDir, 0755)
	if err != nil {
		fmt.Println("Error creating refs/heads directory")
		return
	}

	if err = os.WriteFile(refPath, []byte(commitHash), 0755); err != nil {
		fmt.Printf("Error writing to refs/heads/%s", currentBranch)
		return
	}

}

func GoGitLog() error {
	// Get latest commit hash from HEAD
	currentHash, err := GetHeadHash()
	if err != nil {
		fmt.Println("Error reading hash from HEAD")
		return err
	}

	var history []LogData

	for currentHash != "" {
		_, content, err := readObject(currentHash)
		if err != nil {
			fmt.Println(err)
			return err
		}

		lines := strings.Split(string(content), "\n")
		var c LogData
		nextHash := ""
		isMessage := false

		for _, line := range lines {

			if isMessage {
				c.message = append(c.message, line)
				continue
			}
			// If line is blank, remaining part is the message
			if line == "" {
				isMessage = true
				continue
			}

			if strings.HasPrefix(line, "parent ") {
				nextHash = strings.TrimPrefix(line, "parent ")
			} else if strings.HasPrefix(line, "author ") {
				c.authorLine = strings.TrimPrefix(line, "author ")
			}
		}
		c.hash = currentHash
		history = append(history, c)
		currentHash = nextHash
	}

	for i := len(history) - 1; i >= 0; i-- {
		commit := history[i]
		FormatLog(commit.hash, commit.authorLine, commit.message, i == len(history)-1)
	}

	return nil
}

func GoGitStatus() {
	indexPath := ".gogit/index.json"
	index := LoadIndex(indexPath)

	root := &TrieNode{Children: make(map[string]*TrieNode), Mode: 40000, IsDirty: true}

	ShowStatus(index, root)
}

func GoGitCheckout(branch string) error {

	commitHash, err := ResolvePointer(branch)
	if err != nil || commitHash == "" {
        return fmt.Errorf("error: branch or commit '%s' not found", branch)
    }
	treeHash, err := GetTreeHashFromCommit(commitHash)
	if err != nil {
        return fmt.Errorf("error: could not find tree for commit %s: %v", commitHash, err)
    }

	oldIndex := LoadIndex(".gogit/index.json")

	newIndex := make(map[string]IndexEntry)

	RestoreFromTree(treeHash, ".", newIndex)

	cleanWorkingDirectory(oldIndex, newIndex)

	writeIndex(".gogit/index.json", newIndex)
	updateHEAD(branch, commitHash)

	fmt.Printf("Switched to branch '%s'\n", branch)

	return nil
}

