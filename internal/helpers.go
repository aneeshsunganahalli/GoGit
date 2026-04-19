package internal

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
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

func ParseTreeEntry(data []byte) (name string, node *TrieNode, hashEnd int, err error) {
    // 1. Find the space: format is "[Mode] [Name]\x00[Hash]"
    spaceIdx := bytes.IndexByte(data, ' ')
    if spaceIdx == -1 {
        return "", nil, 0, fmt.Errorf("invalid entry: no space found")
    }
    modeStr := string(data[:spaceIdx])

    // 2. Find the Null Byte (\x00) after the name
    nullIdx := bytes.IndexByte(data[spaceIdx:], 0)
    if nullIdx == -1 {
        return "", nil ,0, fmt.Errorf("invalid entry: no null terminator")
    }
    // nullIdx is relative to data[spaceIdx:], so we add spaceIdx
    absoluteNullIdx := spaceIdx + nullIdx
    name = string(data[spaceIdx+1 : absoluteNullIdx])

    // 3. The 20 bytes immediately following the Null Byte
    hashStart := absoluteNullIdx + 1
    hashEnd = hashStart + 20
    
    if hashEnd > len(data) {
        return "", nil, 0, fmt.Errorf("invalid entry: data too short for hash")
    }

    hash := make([]byte, 20)
    copy(hash, data[hashStart:hashEnd])

    // 4. Construct the node
    mode, _ := strconv.ParseInt(modeStr, 8, 32)
    node = &TrieNode{
        Hash:     hash,
        Mode:     int(mode),
        Children: make(map[string]*TrieNode),
        IsFile:   (mode != 040000), // In Git, 040000 is a Directory
    }

    return name, node, hashEnd, nil
}

// func cleanWorkingDirectory(oldIndex, currIndex map[string]IndexEntry) {
// 	for filePath := range oldIndex {
// 		if _, exists := currIndex[filePath]; !exists {
// 			// File is tracked in current branch, but not in the target branch
// 			err := os.Remove(filePath)
// 			if err != nil && !os.IsNotExist(err) {
// 				fmt.Printf("Warning: failed to remove stale file %s: %v\n", filePath, err)
// 			}
// 		}
// 	}
// }

func cleanWorkingDirectory(oldIndex map[string]IndexEntry, newIndex map[string]IndexEntry) {
	fmt.Println("\n--- Cleanup Phase Debug ---")
	
	// Show what we loaded from the current branch
	fmt.Printf("oldIndex has %d items:\n", len(oldIndex))
	for k := range oldIndex {
		fmt.Printf("  - %s\n", k)
	}

	// Show what we built from the target branch
	fmt.Printf("\nnewIndex has %d items:\n", len(newIndex))
	for k := range newIndex {
		fmt.Printf("  - %s\n", k)
	}

	fmt.Println("\nStarting Deletions:")
	for filePath := range oldIndex {
		if _, exists := newIndex[filePath]; !exists {
			fmt.Printf(" -> Target branch doesn't have '%s'. DELETING...\n", filePath)
			err := os.Remove(filePath)
			if err != nil && !os.IsNotExist(err) {
				fmt.Printf("    [!] Failed to remove %s: %v\n", filePath, err)
			}
		} else {
			fmt.Printf(" -> Target branch has '%s'. Keeping.\n", filePath)
		}
	}
	fmt.Println("---------------------------")
}