package internal

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"

	// "encoding/hex"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// const indexPath = ".gogit/index.json"

// Loads the JSON from the .gogit/index.json file
func LoadIndex(indexPath string) map[string]IndexEntry {

	data, err := os.ReadFile(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]IndexEntry)
		}
		fmt.Println("Error loading the index from index.json: ", err)
		return make(map[string]IndexEntry)
	}

	index := make(map[string]IndexEntry)

	if err := json.Unmarshal(data, &index); err != nil {
		fmt.Println("Failed to unmarshal the json:", err)
		return make(map[string]IndexEntry)
	}

	// prettyJSON, err := json.MarshalIndent(index, "", "  ")

	// fmt.Println(string(prettyJSON))
	return index
}

func UpdateIndexFromPath(targetPath string) {

	indexPath := ".gogit/index.json"
	index := LoadIndex(indexPath)
	seenFiles := make(map[string]bool)

	root := &TrieNode{Children: make(map[string]*TrieNode), Mode: 40000, IsDirty: true}

	for path, entry := range index {
		root.LoadPath(path, entry)
	}

	err := filepath.WalkDir(targetPath, func(path string, d os.DirEntry, err error) error {

		if err != nil || d.IsDir() {
			return err
		}
		cleanedPath := filepath.ToSlash(path)

		if cleanedPath == ".gogit" || strings.HasPrefix(cleanedPath, ".gogit/") {
			return nil
		}

		seenFiles[cleanedPath] = true
		info, _ := os.Stat(path)

		existingEntry, exists := index[cleanedPath]

		if exists && info.Size() == existingEntry.Size && info.ModTime().Unix() == existingEntry.Mtime {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		newHash, err := writeObject("blob", info.Size(), file)
		if err != nil {
			return err
		}

		mode := 100644
		if info.Mode()&0111 != 0 {
			mode = 100755
		}

		newEntry := IndexEntry{
			Filename: filepath.Base(path),
			Size:     info.Size(),
			Mtime:    info.ModTime().Unix(),
			Hash:     newHash,
			Mode:     mode,
		}

		index[cleanedPath] = newEntry
		root.MarkPath(cleanedPath, newEntry)

		return nil
	})

	for path := range index {
		if seenFiles[path] == false {
			if _, err := os.Stat(filepath.FromSlash(path)); os.IsNotExist(err) {
				delete(index, path)
				root.RemovePath(path)
			}
		}
	}

	writeIndex(".gogit/index.json", index)

	// root := BuildTrie(index)
	PrintTrie(root, "")

	root.WriteMerkleTree()

	if err != nil {
		fmt.Println(err)
	}

}

// Writes the JSON to .gogit/index.json
func writeIndex(indexPath string, index map[string]IndexEntry) error {

	// ensure parent dir exists
	if err := os.MkdirAll(filepath.Dir(indexPath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(index, "", "")
	if err != nil {
		fmt.Println("Error writing index into index.json")
	}

	return os.WriteFile(indexPath, data, 0644)
}

// Checks if object exists in storage, but I am programming above function to be in storage when we encounter a hash, so this is redudant
func ObjectExistsInStorage(hash string) bool {

	if len(hash) < 40 { // Git SHA-1 is 40 chars
		return false
	}

	dir, file := hash[:2], hash[2:]

	objectPath := filepath.Join(objectFolder, dir, file)
	// fmt.Println(objectPath)

	info, err := os.Stat(objectPath)
	if err == nil && !info.IsDir() {
		return true
	}

	return false
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

func (n *TrieNode) WriteMerkleTree() []byte {

	if len(n.Hash) == 20 && !n.IsDirty {
		return n.Hash
	}

	// Base Case: We have a file, so we get its hash
	if n.IsFile {
		n.IsDirty = false
		return n.Hash
	}

	// For Merkle Trees, you have to sort the children alphabetically work it to work properly
	keys := make([]string, 0, len(n.Children))
	for name := range n.Children {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	var treeBuffer bytes.Buffer

	for _, name := range keys {
		child := n.Children[name]

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

	n.Hash, _ = hex.DecodeString(hashStr)
	n.IsDirty = false
	return n.Hash
}

func writeObject(objectType string, size int64, r io.Reader) (string, error) {

	tempDir := filepath.Join(".gogit", "tmp")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", err
	}
	tempFile, err := os.CreateTemp(tempDir, "gogit-obj-*")
	if err != nil {
		return "", err
	}

	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	hasher := sha1.New()
	header := fmt.Sprintf("%s %d\x00", objectType, size)

	zlibWriter := zlib.NewWriter(tempFile)
	multiWriter := io.MultiWriter(hasher, zlibWriter)

	multiWriter.Write([]byte(header))

	if _, err := io.Copy(multiWriter, r); err != nil {
		return "", err
	}

	zlibWriter.Close()

	hashStr := fmt.Sprintf("%x", hasher.Sum(nil))

	dir, file := hashStr[:2], hashStr[2:]
	finalDir := filepath.Join(".gogit", "objects", dir)
	finalPath := filepath.Join(finalDir, file)

	if _, err := os.Stat(finalPath); err == nil {
		return hashStr, nil
	}

	if err := os.MkdirAll(finalDir, 0755); err != nil {
		return "", err
	}

	if err = os.Rename(tempFile.Name(), finalPath); err != nil {
		fmt.Println("File already exists")
	}

	return hashStr, err
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

func CleanDirtyFlags(root *TrieNode) {
	root.IsDirty = false
	for _, child := range root.Children {
		CleanDirtyFlags(child)
	}
}

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
