package internal

import (
	"fmt"
	"os"
	"path/filepath"

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
	rootHash := string(root.WriteMerkleTree())

	parentTreeHash, err := GetHeadTreeHash()

	if rootHash == parentTreeHash {
		fmt.Println("On branch main")
		fmt.Println("nothing to commit, working tree clean")
	}

	parentHash, err := GetParentHash()

	commitHash := CreateAndStoreCommit(rootHash, parentHash, message)

	refDir := filepath.Join(".gogit", "refs", "heads") 
	refPath := filepath.Join(refDir, "main")

		err = os.MkdirAll(refDir, 0755)
		if err != nil {
			fmt.Println("Error creating refs/heads directory")
			return 
		}

		if err = os.WriteFile(refPath, []byte(commitHash), 0755); err != nil {
			fmt.Println("Error writing to refs/heads/main")
			return
		}

}

// Archive, previous implementation of writeObject()

// Writes the object into the .gogit/objects/ folder in the format: object/sd/j8k4... for storage
// func WriteObject(objectType string, content string) string {

// 	store := GenerateStore(objectType, content)

// 	hashStr := GenerateHash(objectType, content)
// 	fmt.Println(hashStr)
// 	compressedContent, err := ZlibCompresser(store)

// 	// Directory Creation
// 	dir := hashStr[:2]
// 	file := hashStr[2:]

// 	path := filepath.Join(objectFolder, dir)

// 	err = os.MkdirAll(path, 0644)
// 	if err != nil {
// 		panic(fmt.Sprintf("Failed to create directory at %s: %v", path, err))
// 	}

// 	// File Creation
// 	fileName := filepath.Join(path, file)

// 	os.WriteFile(fileName, compressedContent, 0644)

// 	return hashStr // You'll need this to keep track of what you have saved
// }

// func writeObject(objectType string, size int64, r io.Reader) (string, error) {
// 	tempFile, err := os.CreateTemp("", "gogit-obj-*")
// 	if err != nil {
// 		return "", err
// 	}

// 	defer os.Remove(tempFile.Name())
// 	defer tempFile.Close()

// 	hasher := sha1.New()
// 	header := fmt.Sprintf("%s %d\x00", objectType, size)

// 	zlibWriter := zlib.NewWriter(tempFile)
// 	multiWriter := io.MultiWriter(hasher, zlibWriter)

// 	multiWriter.Write([]byte(header))

// 	if _ ,err := io.Copy(multiWriter, r); err != nil {
// 		return "", err
// 	}

// 	zlibWriter.Close()

// 	hashStr := fmt.Sprintf("%x", hasher.Sum(nil))

// 	dir, file := hashStr[:2], hashStr[2:]
// 	finalDir := filepath.Join(".gogit", "objects", dir)
// 	finalPath := filepath.Join(finalDir, file)

// 	if _, err := os.Stat(finalPath); err == nil {
// 		return hashStr, nil
// 	}

// 	if err := os.MkdirAll(finalDir, 0755); err != nil {
// 		return "", err
// 	}

// 	return hashStr, os.Rename(tempFile.Name(), finalPath)
// }


func main() {
	
}