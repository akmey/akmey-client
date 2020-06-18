package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"

	"github.com/briandowns/spinner"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"r", "u"},
	Short:   "Removes a user's keys",
	Long: `Removes everything from
	this user installed by Akmey`,
	Run: func(cmd *cobra.Command, args []string) {
		// starts... the spinner
		spinner := spinner.New(spinner.CharSets[14], 50*time.Millisecond)
		spinner.Start()
		re := regexp.MustCompile("#-- Akmey START --\n((?:.|\n)+)\n#-- Akmey STOP --")
		db, err := initFileDB(getStoragePath(), keyfile)
		defer db.Close()
		tx, err := db.Begin()
		cfe(err)
		checkstmt, err := tx.Prepare("select id from users where email = ? or name = ? collate nocase")
		cfe(err)
		var check string
		err = checkstmt.QueryRow(args[0], args[0]).Scan(&check)
		if check == "" {
			finalmsg := "👍  " + args[0] + " is not installed\n"
			spinner.FinalMSG = finalmsg
			spinner.Stop()
			os.Exit(0)
		}
		err = nil
		// dear @leonekmi, please comment your code next time
		stmt, err := tx.Prepare("delete from users where email = ? or name = ? collate nocase")
		cfe(err)
		stmt2, err := tx.Prepare("delete from keys where value = ? collate nocase")
		cfe(err)
		stmt3, err := tx.Prepare("select * from keys where user_id = ? collate nocase")
		cfe(err)
		defer checkstmt.Close()
		defer stmt.Close()
		defer stmt2.Close()
		defer stmt3.Close()
		// Step 1 : Fetch installed keys
		rows, err := stmt3.Query(check)
		cfe(err)
		defer rows.Close()
		toberemoved := map[int]string{}
		// Step 2 : Parse the keys in a beautiful map
		dat, err := ioutil.ReadFile(keyfile)
		fileWithKeyRemoved := dat
		for rows.Next() {
			var id int
			var value string
			var comment string
			var user_id string

			err = rows.Scan(&id, &comment, &value, &user_id)

			stmt2.Exec(value)
			toberemoved[id] = "\n" + value + " " + comment
			// creates a temporary slice to convert the key + comment from string to... a slice
			keyByte := []byte(value + " " + comment + "\n")
			// removes the said key, one by one, from the keyfile
			// TODO: only remove in the akmey section of the keyfile
			fileWithKeyRemoved = bytes.Replace(fileWithKeyRemoved, keyByte, []byte(""), 1)
		}

		err = rows.Err()
		cfe(err)
		if len(toberemoved) == 0 {
			fmt.Println("\nThis user does not exist or doesn't have keys registered.")
			os.Exit(1)
		}
		stmt.Exec(args[0], args[0])
		cfe(err)
		match := re.FindStringSubmatch(string(dat))
		if match == nil {
			fmt.Println("Akmey is not present in this file")
			os.Exit(0)
		}

		err = ioutil.WriteFile(keyfile, fileWithKeyRemoved, 0)
		cfe(err)
		tx.Commit()
		finalmsg := "✅ " + args[0] + "'s keys are now removed\n"
		spinner.FinalMSG = finalmsg
		spinner.Stop()
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
