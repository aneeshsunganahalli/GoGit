package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)


func AddChanges() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Testing adding changes",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Add Tester")
			return nil
		},
	}
	return cmd
}