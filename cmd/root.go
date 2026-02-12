package cmd

import (
	"log"
	"os"

	"github.com/BruhaBruh/bruhautomate/internal/config"
	"github.com/BruhaBruh/bruhautomate/internal/parser"
	"github.com/spf13/cobra"
)

var (
	configFile string
	cfg        *config.Config
)

var rootCmd = &cobra.Command{
	Use:                "bruhautomate",
	Aliases:            []string{"bam"},
	Short:              "CLI utility for build and run commands by json files",
	Example:            "bam vrcon --flag -f or bam -c bam.yaml vrcon --flag -f",
	DisableFlagParsing: true,
	Args:               cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		args = loadConfig(args)
		if len(args) == 0 {
			listCommands()
			return
		}

		runCommand(args[0], args[1:])
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

func runCommand(cmd string, args []string) {
	command, err := cfg.Command(cmd)
	if err != nil {
		log.Fatalln(err)
	}

	commandFlags := command.ComposeFlags()
	args, flags := parser.ParseArgsAndFlags(commandFlags, args)

	if err := command.Execute(args, flags, commandFlags); err != nil {
		log.Fatalln(err)
	}
}

func loadConfig(args []string) []string {
	defer func() {
		cfg = config.New(configFile)
		if err := cfg.LoadOrCreate(); err != nil {
			log.Fatalln(err)
		}
	}()

	if len(args) == 0 {
		return args
	}

	if args[0] == "--config" || args[0] == "-c" {
		if len(args) > 1 {
			configFile = args[1]
			return args[2:]
		}
		return []string{}
	}

	return args
}
