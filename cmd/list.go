package cmd

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/BruhaBruh/bruhautomate/internal/command"
	"github.com/BruhaBruh/bruhautomate/internal/config"
	"github.com/BruhaBruh/bruhautomate/internal/flag"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all available commands",
	Aliases: []string{"t"},
	Run: func(cmd *cobra.Command, args []string) {
		cfg = config.New(configFile)
		if err := cfg.LoadOrCreate(); err != nil {
			log.Fatalln(err)
		}

		listCommands()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listCommands() {
	commands := cfg.CommandsSlice()

	slices.SortFunc(commands, func(a, b command.Command) int {
		return strings.Compare(
			strings.ToLower(a.Name),
			strings.ToLower(b.Name),
		)
	})

	color.Blue("List of available commands:")
	for _, cmd := range commands {
		var sb strings.Builder
		aliases := make([]string, 0, 1+len(cmd.Aliases))
		aliases = append(aliases, cmd.Name)
		aliases = append(aliases, cmd.Aliases...)
		sb.WriteString("\t")
		sb.WriteString(strings.Join(aliases, ", "))
		if len(cmd.Description) > 0 {
			sb.WriteString(color.WhiteString(" - %s", cmd.Description))
		}
		flags := cmd.ComposeFlags().Flags()
		slices.SortFunc(flags, func(a, b flag.Flag) int {
			return strings.Compare(
				strings.ToLower(a.Name),
				strings.ToLower(b.Name),
			)
		})
		for _, flag := range flags {
			sb.WriteString(color.WhiteString("\n\t\t--%s", flag.Name))
			if len(flag.Shortcut) > 0 {
				sb.WriteString(color.WhiteString(", -%s", flag.Shortcut))
			}
			if flag.HasValue {
				sb.WriteString(color.WhiteString(" (value)"))
			}
			if len(flag.Description) > 0 {
				sb.WriteString(color.WhiteString(" - %s", flag.Description))
			}
		}
		fmt.Println(sb.String())
	}
}
