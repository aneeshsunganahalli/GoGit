package internal

import (
	"fmt"
)

func PrintTrie(node *TrieNode, indent string) {
    for name, child := range node.Children {
        fmt.Printf("%s%s (Mode: %d, IsFile: %v)\n", indent, name, child.Mode, child.IsFile)
        if !child.IsFile {
            PrintTrie(child, indent + "  ")
        }
    }
}