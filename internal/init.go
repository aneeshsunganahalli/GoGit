package internal

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Initialize(cmd *cobra.Command, args []string) {
	
	// Creates the .gogit directory
	dir := ".gogit"
	err := os.Mkdir(dir, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}

	subfolders := []string{
		"/hooks",
		"/refs",
		"/info",
		"/objects",
	}

	for _, subf := range subfolders {
		err = os.MkdirAll(dir + subf, 0777)
		if err != nil {
			fmt.Println("Error creating subfolders in .gogit")
			return
		}
	}

	

	// Creates the HEAD file
	headPath := dir + "/HEAD"
	head, err := os.Create(headPath)
	if err != nil {
		panic(err)
	}
	defer head.Close()

	head.WriteString("refs:refs/heads/main")

	fmt.Println("Initialized empty repository")	
}
