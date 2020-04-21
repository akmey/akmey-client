package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"regexp"
	//"strings"
	"github.com/spf13/cobra"
	"time"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"r", "u"},
	Short:   "Remove a users key",
	// TODO: make it work
	Long: `Remove a users key`,
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
		getKey, err := tx.Prepare("select value from keys where user_id = ? collate nocase")
		cfe(err)
		defer checkstmt.Close()
		defer stmt.Close()
		defer stmt2.Close()
		defer stmt3.Close()
		defer getKey.Close()
		// Step 1 : Fetch installed keys
		rows, err := stmt3.Query(check)
		cfe(err)
		defer rows.Close()
		//toberemoved := map[int]string{}
		var toberemoved string
		// Step 2 : Parse the keys in a beautiful map
		dat, err := ioutil.ReadFile(keyfile)
		var content string
		for rows.Next() {
			var id int
			var value string
			var comment string
			var user_id string

			err = rows.Scan(&id, &comment, &value, &user_id)

			stmt2.Exec(value)
			fmt.Println("id: ", id)
			fmt.Println("value: ", value)
			fmt.Println("comment: ", comment)
			fmt.Println("user_id:", user_id)
			fmt.Println(`
			`)
			//toberemoved[id] = "\n" + value + " " + comment
			//content = strings.Replace(string(dat), value+" "+comment, "", 1)
			fmt.Println("Content:", content)
			toberemoved += string(content) + " " + string(comment) + "\n"
		}
		fmt.Println("toberemoved: ", len(toberemoved))
		err = rows.Err()
		cfe(err)
		if len(toberemoved) == 0 {
			fmt.Println("\nThis user does not exist or doesn't have keys registered.")
			os.Exit(1)
		}
		stmt.Exec(args[0], args[0])
		newContent := ""
		cfe(err)
		match := re.FindStringSubmatch(string(dat))
		if match == nil {
			fmt.Println("Akmey is not present in this file")
			os.Exit(0)
		}

		// well, for now match matches everything
		/*for nb, torm := range toberemoved {
			if newContent == "" {
				newContent = strings.Replace(string(dat), match[1], torm, -1)
				fmt.Println(nb, "newContent1: ", newContent)
			} else {
				newContent = strings.Replace(newContent, match[1], torm, -1)
				fmt.Println(nb, "newContent2: ", newContent)
			}
		} */

		fmt.Println("newContent final: ", newContent)

		err = ioutil.WriteFile(keyfile, []byte(newContent), 0)
		cfe(err)
		tx.Commit()
		fmt.Println("\n")
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
