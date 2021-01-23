package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "bot",
	Short: "bot info",
	Long:  `A CLI tool to run discord bot `,
}
