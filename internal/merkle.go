package internal

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"path/filepath"
	"sort"
	"strings"

)


func (root *TrieNode) LoadPath(path string, entry IndexEntry) {	
	parts := strings.Split(filepath.ToSlash(path), "/")
	current := root

	for i, part := range parts {
		if _, exists := current.Children[part]; !exists {
			current.Children[part] = &TrieNode{
				Children: make(map[string]*TrieNode),
				IsDirty:  false,
			}

		}
		child := current.Children[part]

		if i == len(parts)-1 {
			child.Hash, _ = hex.DecodeString(entry.Hash)
			child.IsFile = true
			child.Mode = entry.Mode
		} else {
			child.Mode = 40000
			current = child

		}
	}

}

func (root *TrieNode) MarkPath(path string, entry IndexEntry) {
	parts := strings.Split(filepath.ToSlash(path), "/")
	current := root
	current.IsDirty = true

	for i, part := range parts {

		child, exists := current.Children[part]
		if !exists {
			child = &TrieNode{
				Children: make(map[string]*TrieNode),
				Mode:     40000,
			}
			current.Children[part] = child
		}
		child.IsDirty = true

		isLast := i == len(parts)-1

		if isLast {
			child.Hash, _ = hex.DecodeString(entry.Hash)
			child.IsFile = true
			child.Mode = entry.Mode
		} else {
			current = child
		}
	}
}


func (root *TrieNode) RemovePath(path string) {
	parts := strings.Split(filepath.ToSlash(path), "/")
	current := root
	current.IsDirty = true

	for i, part := range parts {
		child, exists := current.Children[part]
		if !exists {
			return // Can't remove something that doesn't exist
		}

		if i == len(parts)-1 {
			delete(current.Children, part)
		} else {
			child.IsDirty = true // deletions are also dirty
			current = child
		}
	}
}

func (root *TrieNode) WriteMerkleTree() []byte {

	if len(root.Hash) == 20 && !root.IsDirty {
		return root.Hash
	}

	// Base Case: We have a file, so we get its hash
	if root.IsFile {
		root.IsDirty = false
		return root.Hash
	}

	// For Merkle Trees, you have to sort the children alphabetically work it to work properly
	keys := make([]string, 0, len(root.Children))
	for name := range root.Children {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	var treeBuffer bytes.Buffer

	for _, name := range keys {
		child := root.Children[name]

		// Need to obtain the child's hash
		childHash := child.WriteMerkleTree()

		// [Mode] [Name][Null Byte][Binary Hash]
		treeBuffer.WriteString(fmt.Sprintf("%o %s\x00", child.Mode, name))
		treeBuffer.Write(childHash)
	}

	// Write the tree object into the .gogit/objects folder
	content := treeBuffer.Bytes()
	hashStr, err := writeObject("tree", int64(len(content)), bytes.NewReader(content))
	if err != nil {
		fmt.Println("Error in writing tree object: ", err)
	}

	root.Hash, _ = hex.DecodeString(hashStr)
	root.IsDirty = false
	return root.Hash
}

func CleanDirtyFlags(root *TrieNode) {
	root.IsDirty = false
	for _, child := range root.Children {
		CleanDirtyFlags(child)
	}
}
