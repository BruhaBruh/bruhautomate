package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/BruhaBruh/bruhautomate/internal/config"
	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:                "__complete",
	Hidden:             true,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		toComplete := ""
		if len(args) > 0 {
			toComplete = args[len(args)-1]
		}
		fmt.Fprintf(os.Stderr, "toComplete bytes: %v\n", []byte(toComplete))

		c := config.New(configFile)
		if err := c.LoadOrCreate(); err != nil {
			return
		}

		// если первый аргумент — команда из конфига, completim её флаги
		if len(args) >= 1 {
			commandName := args[0]
			command, err := c.Command(commandName)
			if err == nil {
				flags := command.ComposeFlags().Flags()
				for _, f := range flags {
					if strings.HasPrefix("--"+f.Name, toComplete) {
						if len(f.Description) > 0 {
							fmt.Printf("--%s\t%s\n", f.Name, f.Description)
						} else {
							fmt.Printf("--%s\n", f.Name)
						}
					}
					if len(f.Shortcut) > 0 && strings.HasPrefix("-"+f.Shortcut, toComplete) {
						if len(f.Description) > 0 {
							fmt.Printf("-%s\t%s\n", f.Shortcut, f.Description)
						} else {
							fmt.Printf("-%s\n", f.Shortcut)
						}
					}
				}
				return
			}
		}

		// иначе — completim команды (cobra subcommands + конфиг)
		for _, sub := range cmd.Root().Commands() {
			if sub.Hidden {
				continue
			}
			if strings.HasPrefix(sub.Name(), toComplete) {
				fmt.Printf("%s\t%s\n", sub.Name(), sub.Short)
			}
			for _, alias := range sub.Aliases {
				if strings.HasPrefix(alias, toComplete) {
					fmt.Printf("%s\t%s\n", alias, sub.Short)
				}
			}
		}

		for _, command := range c.CommandsSlice() {
			if strings.HasPrefix(command.Name, toComplete) {
				if len(command.Description) > 0 {
					fmt.Printf("%s\t%s\n", command.Name, command.Description)
				} else {
					fmt.Println(command.Name)
				}
			}
			for _, alias := range command.Aliases {
				if strings.HasPrefix(alias, toComplete) {
					if len(command.Description) > 0 {
						fmt.Printf("%s\t%s\n", alias, command.Description)
					} else {
						fmt.Println(alias)
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)
}
