package internal

import (
	"encoding/hex"
	"fmt"
	"strings"

)

func PrintTrie(node *TrieNode, indent string) {
    for name, child := range node.Children {
        fmt.Printf("%s%s (Mode: %d, IsFile: %v)\n", indent, name, child.Mode, child.IsFile)
        if !child.IsFile {
            PrintTrie(child, indent + "  ")
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