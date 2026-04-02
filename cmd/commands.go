package cmd

import (
	// "fmt"

	"fmt"

	"github.com/aneeshsunganahalli/GoGit/internal"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command {
	Use: "init",
	Short: "Initialize .git",
	Run: internal.Initialize,
	// Run: internal.Hashing,
}

var addCmd = &cobra.Command{
	Use:   "add [path]",
	Short: "Testing adding changes",
	Args: cobra.ExactArgs(1) ,
	RunE: func(cmd *cobra.Command, args []string) error {
			// path := args[0]
		// internal.WriteObject()
		return nil
	},
}

func init(){
	rootCmd.AddCommand(initCmd, addCmd)
}

// helloooo worlssdsdsjfhfuahfuf

func random(){
	fmt.Println("Everything")
}