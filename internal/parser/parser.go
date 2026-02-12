package parser

import (
	"log"
	"regexp"

	"github.com/BruhaBruh/bruhautomate/internal/flag"
)

var (
	flagRegexp         = regexp.MustCompile("^--([A-Za-z0-9_-]+)$")
	flagShortcutRegexp = regexp.MustCompile("^-([A-Za-z0-9_-]+)$")
)

func ParseArgsAndFlags(commandFlags *flag.Flags, input []string) (args []string, flags *flag.Flags) {
	args = make([]string, 0, len(input))
	flags = flag.NewFlags()

	for i := 0; i < len(input); i++ {
		arg := input[i]
		if flagRegexp.MatchString(arg) {
			match := flagRegexp.FindStringSubmatch(arg)
			flag, err := commandFlags.Find(match[1])
			if err != nil {
				args = append(args, arg)
				continue
			}
			if flag.HasValue {
				if len(input) == i+1 {
					log.Fatalf("flag %s must has value", match[0])
				}
				value := input[i+1]
				flags.AddFlag(flag, value)
				i += 1
			} else {
				flags.AddFlag(flag)
			}
		} else if flagShortcutRegexp.MatchString(arg) {
			match := flagShortcutRegexp.FindStringSubmatch(arg)
			flag, err := commandFlags.Find(match[1])
			if err != nil {
				args = append(args, arg)
				continue
			}
			if flag.HasValue {
				if len(input) == i+1 {
					log.Fatalf("flag %s must has value", match[0])
				}
				value := input[i+1]
				flags.AddFlag(flag, value)
				i += 1
			} else {
				flags.AddFlag(flag)
			}
		} else {
			args = append(args, arg)
		}
	}

	return args, flags
}
