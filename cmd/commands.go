package cmd

import (
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
	Use:   "add",
	Short: "Testing adding changes",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Add Tester")
		return nil
	},
}

func init(){
	rootCmd.AddCommand(initCmd, addCmd)
}
