package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// teamInstallCmd represents the teamInstall command
var teamInstallCmd = &cobra.Command{
	Use:   "teamInstall",
	Aliases: []string{"ti"},
	Short: "Install a team's members' keys",
	// TODO: add a long description
	Long: `Install team's members' keys`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("teamInstall called")
	},
}

func init() {
	rootCmd.AddCommand(teamInstallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// teamInstallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// teamInstallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
