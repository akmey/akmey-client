package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"regexp"
	//"strings"
	"time"

	"github.com/spf13/cobra"
)

// resetCmd represents the remove command
var resetCmd = &cobra.Command{
	Use:     "reset",
	Aliases: []string{"ua"},
	Short:   "Remove every keys",
	// TODO: add a long description
	Long: `Remove every keys from the keyfile`,
	Run: func(cmd *cobra.Command, args []string) {
		// starts... the spinner
		spinner := spinner.New(spinner.CharSets[14], 50*time.Millisecond)
		spinner.Start()
		re := regexp.MustCompile("#-- Akmey START --\n((?:.|\n)+)\n#-- Akmey STOP --")
		db, err := initFileDB(getStoragePath(), keyfile)
		defer db.Close()
		tx, err := db.Begin()
		cfe(err)
		tx.Exec("delete from users")
		cfe(err)
		tx.Exec("delete from keys")
		cfe(err)
		// Step 1 : Fetch installed keys
		dat, err := ioutil.ReadFile(keyfile)
		newContent := ""
		cfe(err)
		match := re.FindStringSubmatch(string(dat))
		if match == nil {
			fmt.Println("Akmey is not present in this file")
			os.Exit(0)
		}
		/* for _, torm := range toberemoved {
			if newContent == "" {
				newContent = strings.Replace(string(dat), match[1], torm, -1)
			} else {
				newContent = strings.Replace(newContent, match[1], torm, -1)
			}
		} */
		err = ioutil.WriteFile(keyfile, []byte(newContent), 0)
		cfe(err)
		tx.Commit()
		fmt.Println("\n")
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// resetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// resetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
