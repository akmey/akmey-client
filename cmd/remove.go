package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"r"},
	Short:   "Remove a users key",
	// TODO: add a long description
	Long: `Remove a users key`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("remove called")

	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
