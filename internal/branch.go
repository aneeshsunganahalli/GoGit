package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateBranch(branchName string) error {

	commitHash, err := GetHeadHash()
	if err != nil {
		fmt.Println("fatal: could not obtain commit hash")
		return err
	}

	refDir := filepath.Join(".gogit", "refs", "heads")
	refPath := filepath.Join(refDir, branchName)

	if _, err = os.Stat(refPath); err == nil {
		fmt.Printf("\nBranch %s already exists\n", branchName)
		return nil
	}

	err = os.MkdirAll(refDir, 0755)
	if err != nil {
		fmt.Println("Error creating refs/heads directory")
		return err
	}

	if err = os.WriteFile(refPath, []byte(commitHash), 0755); err != nil {
		fmt.Println("Error writing to refs/heads/", branchName)
		return err
	}

	fmt.Println("\nSuccessfully created branch", branchName)
	return nil
}