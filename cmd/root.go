package cmd

import (
	"fmt"
	"os"
	"log"
	"database/sql"
//	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func initFileDB(storagepath string, keyfilepath string) (*sql.DB, error) {
        var dbpath string
	home, err := homedir.Expand("~")
	cfe(err)
	storagepath = home + "/.akmey"
	if _, err := os.Stat(storagepath); os.IsNotExist(err) {
		err = os.MkdirAll(storagepath, 0755)
		cfe(err)
	}

        //fullfilepath, err := filepath.Abs(keyfilepath)
        dbpath = storagepath + "/keys.db"
        db, err := sql.Open("sqlite3", "file:"+dbpath+"?cache=shared&mode=rwc")
        cfe(err)
        sqlStmt := `
        create table if not exists users (id integer not null, name text, email text);
        create table if not exists keys (id integer not null, comment text, value text, user_id integer not null);
        create table if not exists teams (id integer not null, name text, bio name, keys text)
	`
        _, err = db.Exec(sqlStmt)
        return db, err
}


func cfe(err error) bool {
        if err != nil {
                log.Panicln(err)
                return false
        }
        return true
}


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "akmey",
	Short: "Add/Remove SSH keys to grant access to your friends, coworkers, etc...",
	Version: "v0.2-alpha",
	  TraverseChildren: true,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// we can't just homedir.Expand("~/.ssh/authorized_e=keys") because it will fail if the file doesn't exist, so we basically just get user's home directory and add "/.ssh" to it
	var dest string
	var server string
	home, err := homedir.Expand("~/")
	cfe(err)
	sshfolder := home + "/.ssh"
	_ = os.Mkdir(sshfolder, 755) // create the dir (w/ correct permissions) and ignores errors, according to stackoverflow. It's not that good but hey, it works ¯\_(ツ)_/¯
	keyfile := sshfolder + "/authorized_keys"
	os.OpenFile(keyfile, os.O_RDONLY|os.O_CREATE, 0755) // create the file (w/ corrects permissions) if it doesn't already exist, a bit better than for the ssh dir

	server = "https://akmey.leonekmi.fr"

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.e

	rootCmd.PersistentFlags().StringVar(&dest, "dest", keyfile, "Where Akmey should act (your authorized_keys file)")
	rootCmd.PersistentFlags().StringVar(&server, "server", server, "Specify a custom Akmey server here")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".akmey" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".akmey")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
