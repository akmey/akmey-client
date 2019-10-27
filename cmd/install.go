package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/briandowns/spinner"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

var key string

type Team struct {
	ID    float64
	Users []User
	Bio   string
}

type SSHKey struct {
	ID       float64
	Key      string
	Comment  string
	User     User
	LastEdit float64
}

// User is a user representation used for the API
type User struct {
	ID    float64
	Name  string
	Email string
	Keys  []SSHKey
}

// cfe panic in case of an error
/* func cfe(err error) bool {
        if err != nil {
                log.Panicln(err)
                return false
        }
        return true
} */

// fetchUser return User named `user` on the `server`
func fetchUser(user string, server string) (User, error) {
	resp, err := resty.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Accept", "application/json").
		Get(server + "/api/user/match/" + user)
	cfe(err)
	var f User
	err = json.Unmarshal(resp.Body(), &f)
	return f, err
}

// fetchUserSpecificKey returns User named `user` on the `server`, only with a specific key
func fetchUserSpecificKey(user string, key string, server string) (User, error) {
	resp, err := resty.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Accept", "application/json").
		// api path + user/email + filter to search + comment
		Get(server + "/api/user/match/" + user + "?filter=" + key)
	cfe(err)
	var f User
	err = json.Unmarshal(resp.Body(), &f)
	return f, err
}

func fetchTeam(team string, server string) (Team, error) {
	resp, err := resty.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Accept", "application/json").
		Get(server + "/api/team/match/" + team)
	cfe(err)
	var j Team
	err = json.Unmarshal(resp.Body(), &j)
	return j, err
}

// blatantly stolen from https://stackoverflow.com/questions/37145935/checking-if-a-value-exists-in-sqlite-db-with-go
func checkIfUserExists(db *sql.DB, username string) bool {
	sqlStmt := `SELECT name FROM users WHERE name = ? COLLATE NOCASE`
	err := db.QueryRow(sqlStmt, username).Scan(&username)
	if err != nil {
		if err != sql.ErrNoRows {
			cfe(err)
		}

		return false
	}

	return true
}

func getStoragePath() string {
	// we can't just homedir.Expand("~/.ssh/authorized_e=keys") because it will fail if the file doesn't exist, so we basically just get user's home directory and add "/.ssh" at it
	home, err := homedir.Expand("~/")
	cfe(err)
	sshfolder := home + "/.ssh"
	// create the dir (w/ correct permissions) and ignores errors, according to stackoverflow. It's not that good but hey, it works ¯\_(ツ)_/¯
	_ = os.Mkdir(sshfolder, 755)
	// create the file (w/ corrects permissions) if it doesn't already exist, a bit better than for the ssh dir
	_, err = os.OpenFile(keyfile, os.O_RDONLY|os.O_CREATE, 0644)
	cfe(err)
	home, err = homedir.Expand("~")
	cfe(err)
	storagepath := home + "/.akmey"
	return storagepath
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:     "install",
	Aliases: []string{"i", "in", "inst", "insta", "instal"},
	Short:   "Install a user's key",
	// TODO: add a long description
	Long: `Install a user's key'`,
	Run: func(cmd *cobra.Command, args []string) {
		var finalmsg string
		if len(args) < 1 {
			fmt.Println("Please enter someone's name.\nExample: akmey install Luclu7")
			return
		}
		// starts... the spinner
		spinner := spinner.New(spinner.CharSets[14], 50*time.Millisecond)
		spinner.Start()
		re := regexp.MustCompile("#-- Akmey START --\n((?:.|\n)+)\n#-- Akmey STOP --")
		// same stuff as usual
		db, err := initFileDB(getStoragePath(), keyfile)
		fmt.Println(keyfile)
		defer db.Close()
		tx, err := db.Begin()
		cfe(err)
		// TODO: check if someone's keys are already installed
		//checkstmt, err := tx.Prepare(`select name from users where email = "?" or name = "?" collate nocase`)
		//check := "select name from users where email = \"" + args[0] + "\" or name = \"" + args[0] + "\" collate nocase"
		check := checkIfUserExists(db, args[0])
		if check == true {
			finalmsg := "👍  " + args[0] + " is already installed\n"
			spinner.FinalMSG = finalmsg
			spinner.Stop()
			os.Exit(0)
		}
		stmt, err := tx.Prepare("insert into users(id, name, email) values(?, ?, ?)")
		cfe(err)
		// id = key id on server's side, value = the key itself, comment = key name, userid = user's id
		stmt2, err := tx.Prepare("insert into keys(id, value, comment, user_id) values(?, ?, ?, ?)")
		cfe(err)
		defer stmt.Close()
		defer stmt2.Close()
		var tobeinserted string
		// check if --key is used
		if len(key) < 1 {
			// api
			user, err := fetchUser(args[0], server)
			cfe(err)
			for _, key := range user.Keys {
				_, err = stmt2.Exec(key.ID, key.Key, key.Comment, user.ID)
				cfe(err)
				tobeinserted += key.Key + " " + key.Comment + "\n"
			}
			// add user to sqlite db
			_, _ = stmt.Exec(user.ID, user.Name, user.Email)
		} else {
			// api
			user, err := fetchUserSpecificKey(args[0], key, server)
			cfe(err)
			for _, key := range user.Keys {
				_, _ = stmt2.Exec(key.ID, key.Key, key.Comment, user.ID)
				tobeinserted += key.Key + " " + key.Comment + "\n"
			}
			// add user to sqlite db
			_, err = stmt.Exec(user.ID, user.Name, user.Email)
			cfe(err)
		}

		if tobeinserted == "" {
			finalmsg = "❌ " + args[0] + " has not been installed\n"
			spinner.FinalMSG = finalmsg
			fmt.Println("This user either does not exist or doesn't have any keys registered.")
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
		finalmsg = "✅ " + args[0] + " is now installed\n"
		spinner.FinalMSG = finalmsg
		spinner.Stop()

		return
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	installCmd.Flags().StringVarP(&key, "key", "k", "", "key you want")
}
