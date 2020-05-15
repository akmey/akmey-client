package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

// teamInstallCmd represents the teamInstall command
var teamInstallCmd = &cobra.Command{
	Use:     "teamInstall",
	Aliases: []string{"ti", "tin", "tinst"},
	Short:   "Install a team's keys",
	// TODO: add a long description
	Long: `Install a team's keys'`,
	Run: func(cmd *cobra.Command, args []string) {
		var finalmsg string
		if len(args) < 1 {
			fmt.Println("Please enter a team's name.\nExample: akmey install Luclu7")
			return
		}
		// starts... the spinner
		team, err := fetchTeam(args[0], server)
		cfe(err)
		for _, user := range team.Users {
			spinner := spinner.New(spinner.CharSets[14], 50*time.Millisecond)
			spinner.Start()
			re := regexp.MustCompile("#-- Akmey START --\n((?:.|\n)+)\n#-- Akmey STOP --")
			// same stuff as usual
			db, err := initFileDB(getStoragePath(), keyfile)
			defer db.Close()
			tx, err := db.Begin()
			cfe(err)
			// TODO: check if someone's keys are already installed
			//checkstmt, err := tx.Prepare(`select name from users where email = "?" or name = "?" collate nocase`)
			//check := "select name from users where email = \"" + args[0] + "\" or name = \"" + args[0] + "\" collate nocase"
			check := checkIfUserExists(db, user.Name)
			if check {
				finalmsg := "👍  " + user.Name + " is already installed\n"
				spinner.FinalMSG = finalmsg
				spinner.Stop()
			}
			stmt, err := tx.Prepare("insert into users(id, name, email) values(?, ?, ?)")
			cfe(err)
			// id = key id on server's side, value = the key itself, comment = key name, userid = user's id
			stmt2, err := tx.Prepare("insert into keys(id, value, comment, user_id) values(?, ?, ?, ?)")
			cfe(err)
			defer stmt.Close()
			defer stmt2.Close()
			var tobeinserted string
			// api
			specificUser, err := fetchUserSpecificKey(user.Name, key, server)
			cfe(err)
			for _, key := range user.Keys {
				_, _ = stmt2.Exec(key.ID, key.Key, key.Comment, user.ID)
				tobeinserted += key.Key + " " + key.Comment + "\n"
			}
			// add user to sqlite db
			_, err = stmt.Exec(specificUser.ID, specificUser.Name, specificUser.Email)
			cfe(err)
			if tobeinserted == "" {
				finalmsg = "❌ " + user.Name + " has not been installed\n"
				spinner.FinalMSG = finalmsg
				//fmt.Println("This user either does not exist or doesn't have any keys registered.")
				spinner.Stop()
				os.Exit(1)
			}
			dat, err := ioutil.ReadFile(keyfile)
			cfe(err)
			match := re.FindStringSubmatch(string(dat))
			// insert keys into authorized_keys
			if match == nil {
				tobeinserted = "\n#-- Akmey START --\n" + tobeinserted
				tobeinserted += "#-- Akmey STOP --\n"
				f, err := os.OpenFile(keyfile, os.O_APPEND|os.O_WRONLY, 0600)
				cfe(err)
				defer f.Close()

				_, err = f.WriteString(tobeinserted)
				cfe(err)
			} else {
				tobeinserted = match[1] + tobeinserted
				newContent := strings.Replace(string(dat), match[1], tobeinserted, -1)
				err = ioutil.WriteFile(keyfile, []byte(newContent), 0)
				cfe(err)
			}
			err = tx.Commit()
			cfe(err)
			finalmsg = "✅ " + user.Name + " is now installed\n"
			spinner.FinalMSG = finalmsg
			spinner.Stop()
		}

		return
	},
}

func init() {
	rootCmd.AddCommand(teamInstallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
