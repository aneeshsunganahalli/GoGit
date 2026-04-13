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

// Obtains latest commit hash, using the file HEAD points to, usually refs/head/main
func GetHeadHash() (string, error) {

	headContent, err := os.ReadFile(".gogit/HEAD")
	if err != nil {
		fmt.Println("Cannot read HEAD")
		return "", err
	}

	content := strings.TrimSpace(string(headContent))
	if content == "" {
		return "", nil
	}
	var refPath string

	if strings.HasPrefix(content, "refs:") {
		relPath := strings.TrimPrefix(content, "refs:")
		refPath = filepath.Join(".gogit", relPath)
		} else {
			// Detached HEAD contains hash
			return content, nil
		}

	if _, err := os.Stat(refPath); os.IsNotExist(err) {
			fmt.Println("No parent hash, this is the first commit")
			return "", nil
	}

	hash, err := os.ReadFile(refPath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(hash)), nil
}

// Obtains the tree hash from the latest commit, so we can avoid empty commits from happening
func GetHeadTreeHash() (string, error) {
	headHash, err := GetHeadHash()
	if err != nil {
		fmt.Println("Error reading hash from HEAD")
		return "", err
	}

	if headHash == "" {
    return "", nil
  }

	_ , content, err := readObject(headHash)
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

	return "", fmt.Errorf("Tree hash not found in commit")
} 

func ParseCommit(content []byte) CommitObject {
    lines := strings.Split(string(content), "\n")
    c := CommitObject{}
    for _, line := range lines {
        if strings.HasPrefix(line, "tree ") {
            c.TreeHash = strings.TrimPrefix(line, "tree ")
        } else if strings.HasPrefix(line, "parent ") {
            c.ParentHash = strings.TrimPrefix(line, "parent ")
        } else if strings.HasPrefix(line, "author ") {
					c.Author = strings.TrimPrefix(line, "author ")
					c.Committer = c.Author
				} else if strings.HasPrefix(line, "author ") {
					c.Author = strings.TrimPrefix(line, "author ")
				}
    }
    return c
}


