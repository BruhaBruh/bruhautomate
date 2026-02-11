package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/BruhaBruh/bruhautomate/config"
	"github.com/spf13/cobra"
)

var (
	configFile string
)

var rootCmd = &cobra.Command{
	Use:     "bruhautomate",
	Aliases: []string{"bam"},
	Short:   "CLI utility for build and run commands by json files",
	Example: "bam vrcon -- --flag -f",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.LoadOrCreate(configFile)
		if err != nil {
			log.Fatalln(err.Error())
		}

		if len(args) < 1 {
			log.Fatalln("Command not found")
		}

		err = config.Execute(args[0], args[1:])
		if err != nil {
			log.Fatalf("Fail to execute command \"%s\": %s", strings.Join(args, " "), err.Error())
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file (default is $HOME/.bam.json)")
}
