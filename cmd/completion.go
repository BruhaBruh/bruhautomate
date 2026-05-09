package cmd

import (
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/cobra"
)

type completionData struct {
	Name    string
	Aliases []string
}

var completionCmd = &cobra.Command{
	Use:       "completion [bash|zsh|fish|powershell]",
	Short:     "Generate shell completion script",
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	RunE: func(cmd *cobra.Command, args []string) error {
		data := completionData{
			Name:    cmd.Root().Name(),
			Aliases: cmd.Root().Aliases,
		}
		switch args[0] {
		case "bash":
			return bashCompletion(data)
		case "zsh":
			return zshCompletion(data)
		case "fish":
			return fishCompletion(data)
		default:
			return fmt.Errorf("unsupported shell: %s", args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}

func bashCompletion(data completionData) error {
	t := template.Must(template.New("").Parse(`_{{.Name}}_completion() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local IFS=$'\n'
    COMPREPLY=( $({{.Name}} __complete "${COMP_WORDS[@]:1}" 2>/dev/null | awk -F'\t' '{print $1}') )
}
complete -F _{{.Name}}_completion {{.Name}}{{range .Aliases}} {{.}}{{end}}`))
	return t.Execute(os.Stdout, data)
}

func zshCompletion(data completionData) error {
	t := template.Must(template.New("").Parse(`#compdef {{.Name}}{{range .Aliases}} {{.}}{{end}}

_{{.Name}}() {
    local -a completions
    local IFS=$'\n'
    completions=("${(@f)$({{.Name}} __complete "${words[@]:1}" 2>/dev/null)}")
    _describe '{{.Name}} commands' completions
}

_{{.Name}}`))
	return t.Execute(os.Stdout, data)
}

func fishCompletion(data completionData) error {
	t := template.Must(template.New("").Parse(`function __{{.Name}}_complete
    set -l tokens (commandline -opc)
    set -l current (commandline -ct)
    {{.Name}} __complete $tokens[2..] $current 2>/dev/null
end

complete -c {{.Name}} -f -a '(__{{.Name}}_complete)'{{range .Aliases}}
complete -c {{.}} -f -a '(__{{$.Name}}_complete)'{{end}}`))
	return t.Execute(os.Stdout, data)
}
