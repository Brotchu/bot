package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Brotchu/bot/db"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure your bot id",
	Long:  `Configure your bot with your secret bot id`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("configure called")

		home, _ := homedir.Dir()
		dbPath := filepath.Join(home, "botid.db")

		must(db.Init(dbPath))

		must(db.SetBotID(args[0]))
		fmt.Println("Bot Token Configured")
	},
}

func init() {
	RootCmd.AddCommand(configureCmd)
}

func must(err error) {
	if err != nil {
		fmt.Println("[Err] ", err.Error())
		os.Exit(1)
	}
}
