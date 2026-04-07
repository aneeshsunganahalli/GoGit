package cmd

import (

	"github.com/aneeshsunganahalli/GoGit/internal"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command {
	Use: "init",
	Short: "Initialize .gogit",
	Run: internal.GoGitInit,
}

var addCmd = &cobra.Command{
	Use:   "add [path]",
	Short: "Adds modified, created, removes deleted files",
	Args: cobra.MinimumNArgs(1) ,
	RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			internal.GoGitAdd(path)
		return nil
	},
}

var commitCmd = &cobra.Command{
	Use: "commit [message]",
	Short: "Saves the snapshot at that moment.",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		message := args[0]
		internal.GoGitCommit(message)
		return nil
	},
}

func init(){
	rootCmd.AddCommand(initCmd, addCmd, commitCmd)
}
