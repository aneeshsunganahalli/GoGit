package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Commit should look like this
// tree [40-char-hex-root-tree-hash]
// parent [40-char-hex-parent-commit-hash]
// author Aneesh <aneesh@example.com> 1775151513 +0530
// committer Aneesh <aneesh@example.com> 1775151513 +0530

// [Your commit message here]

func (c *CommitObject) ConvertToString() string {
	var output strings.Builder

	output.WriteString(fmt.Sprintf("tree %s\n", c.TreeHash))
	fmt.Println(c.TreeHash)



	// Could be first commit, so no parent hash
	if c.ParentHash != "" {
		output.WriteString(fmt.Sprintf("parent %s\n", c.ParentHash))
	}

	output.WriteString(fmt.Sprintf("author %s %d %s\n", c.Author, c.Timestamp, c.Timezone))
	output.WriteString(fmt.Sprintf("committer %s %d %s\n", c.Author, c.Timestamp, c.Timezone))
	output.WriteString(fmt.Sprintf("\n%s\n", c.Message))

	return output.String()
}

func CreateAndStoreCommit(treeHash string, parentHash string, message string) string {

	name, email, _ := GetAuthorDetails()
	authorInfo := fmt.Sprintf("%s <%s>", name, email)

	timestamp := time.Now().Unix()
	_, offset := time.Now().Zone()
	timezone := fmt.Sprintf("%+03d%02d", offset/3600, (offset%3600)/60)

	c := CommitObject{
		TreeHash: treeHash,
		ParentHash: parentHash,
		Author: authorInfo,
		Committer: authorInfo,
		Timestamp: timestamp,
		Timezone: timezone,
		Message: message,
	}

	commit := c.ConvertToString()
	hash, _ := writeObject("commit", int64(len(commit)), strings.NewReader(commit))

	return hash
}

func Commit(message string) {
	indexPath := ".gogit/index.json"
	index := LoadIndex(indexPath)

	root := &TrieNode{Children: make(map[string]*TrieNode), Mode: 40000, IsDirty: true}

	for path, entry := range index {
		root.LoadPath(path, entry)
	}

	PrintTrie(root, "")
	rootHash := root.WriteMerkleTree()

	parentHash := GetParentHash()

	commitHash := CreateAndStoreCommit(string(rootHash), parentHash, message)

	refDir := filepath.Join(".gogit", "refs", "heads") 
	refPath := filepath.Join(refDir, "main")

		err := os.MkdirAll(refDir, 0755)
		if err != nil {
			fmt.Println("Error creating refs/heads directory")
			return 
		}

		if err = os.WriteFile(refPath, []byte(commitHash), 0755); err != nil {
			fmt.Println("Error writing to refs/heads/main")
			return
		}

}

// Obtains parent commit hash, using the file HEAD points to, usually refs/head/main
func GetParentHash() string {

	headContent, err := os.ReadFile(".gogit/HEAD")
	if err != nil {
		fmt.Println("Cannot read HEAD")
		return ""
	}

	content := strings.TrimSpace(string(headContent))
	var refPath string

	if strings.HasPrefix(content, "refs:") {
		relPath := strings.TrimPrefix(content, "refs:")
		refPath = filepath.Join(".gogit", relPath)
		} else {
			// Detached HEAD contains hash
			return content
		}

	if _, err := os.Stat(refPath); os.IsNotExist(err) {
			fmt.Println("No parent hash, this is the first commit")
			return ""
	}

	hash, err := os.ReadFile(refPath)
	if err != nil {
		return ""
	}

	return string(hash)
}


func GetCommitHash() 