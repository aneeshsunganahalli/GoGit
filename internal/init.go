package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func Initialize(cmd *cobra.Command, args []string) {

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

	fmt.Println("Initialized empty repository")
}
