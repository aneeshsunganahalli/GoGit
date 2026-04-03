package internal

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

func GetAuthorDetails() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	if _, err := os.Stat(".gogit/config"); err != nil {
		fmt.Println("Config doesn't exist.")

		fmt.Print("Enter your name: ")
		userName, _ := reader.ReadString('\n')

		fmt.Print("Enter your email: ")
		userEmail, _ := reader.ReadString('\n')

		userName = strings.TrimSpace(userName)
		userEmail = strings.TrimSpace(userEmail)

		WriteUserConfig(userName, userEmail)
		return userName, userEmail, nil
	}

	data, err := os.ReadFile(".gogit/config")
	if err != nil {
		return "", "", fmt.Errorf("Failed to read config")
	}
	content := string(data)

	userName, userEmail := ParseUserConfig(content)

	return userName, userEmail, nil

}

func ParseUserConfig(content string) (string, string) {
	var name, email string

	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line := strings.TrimSpace(line)

		if strings.HasPrefix(line, "name =") {
			name = strings.TrimSpace(strings.TrimPrefix(line, "name ="))
		}

		if strings.HasPrefix(line, "email =") {
			email = strings.TrimSpace(strings.TrimPrefix(line, "email ="))
		}
	}

	return name, email
}

func WriteUserConfig(name, email string) error {

	configContent := fmt.Sprintf(
		`[user]
    name = %s
    email = %s
`, name, email)

	err := os.WriteFile(".gogit/config", []byte(configContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func PrintTrie(node *TrieNode, indent string) {
	for name, child := range node.Children {
		fmt.Printf("%s%s (Mode: %d, IsFile: %v)\n", indent, name, child.Mode, child.IsFile)
		if !child.IsFile {
			PrintTrie(child, indent+"  ")
		}
	}
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
				binaryHash, _ := hex.DecodeString(entry.Hash)
				current.Children[part] = &TrieNode{
					Hash:   binaryHash,
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
