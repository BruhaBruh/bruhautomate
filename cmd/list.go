package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/BruhaBruh/bruhautomate/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available commands",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.LoadOrCreate(configFile)
		if err != nil {
			log.Fatalln(err.Error())
		}
		fmt.Println("\033[34mList of available commands:\t\033[0m")
		for _, cmd := range config.Commands {
			var sb strings.Builder
			sb.WriteString("\t\033[0m")
			sb.WriteString(cmd.Name)
			if len(cmd.Aliases) > 0 {
				sb.WriteString(" ")
				sb.WriteString(strings.Join(cmd.Aliases, " "))
			}
			if len(cmd.Description) > 0 {
				sb.WriteString("\033[90m: ")
				sb.WriteString(cmd.Description)
			}
			sb.WriteString("\t\033[0m")
			fmt.Println(sb.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
