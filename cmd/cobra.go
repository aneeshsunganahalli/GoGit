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

var logCmd = &cobra.Command{
	Use: "log",
	Short: "Shows all commits till HEAD",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := internal.GoGitLog()
		return err
	},
}

var branchCmd = &cobra.Command{
	Use: "branch [branch-name]",
	Short: "Creates a new branch at the current commit",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		branchName := args[0]
		err := internal.CreateBranch(branchName)
		return err
	},
}

func init(){
	rootCmd.AddCommand(initCmd, addCmd, commitCmd, logCmd, branchCmd)
}
