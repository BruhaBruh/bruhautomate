package instruction

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/BruhaBruh/bruhautomate/internal/flag"
)

var (
	allArgumentsRegexp                   = regexp.MustCompile(`\$@`)
	argumentRegexp                       = regexp.MustCompile(`\$(\d+)`)
	argumentInBracketsRegexp             = regexp.MustCompile(`\${(\d+)}`)
	flagRegexp                           = regexp.MustCompile(`\$([A-Za-z0-9_-]+)`)
	flagInBracketsRegexp                 = regexp.MustCompile(`\${([A-Za-z0-9_-]+)}`)
	flagInBracketsWithDefaultValueRegexp = regexp.MustCompile(`\${([A-Za-z0-9_-]+)=([^}]+)}`)
)

// $@ - all arguments. Blank if empty
// $n, ${n} - n argument. Blank if empty
// $[A-Za-z0-9_-]+, ${[A-Za-z0-9_-]+} - flag. If non valued returns boolean by existing otherwise returns value or blank. As fallback use environment variable if name of flag is uppercased.
// ${[A-Za-z0-9_-]+=[^}]+} - flag w/ default value
func Build(
	raw string,
	args []string,
	flags *flag.Flags,
	commandFlags *flag.Flags,
) string {
	result := raw
	// $@
	result = allArgumentsRegexp.ReplaceAllStringFunc(result, func(s string) string {
		return strings.Join(args, " ")
	})
	// ${n}
	result = argumentInBracketsRegexp.ReplaceAllStringFunc(result, func(s string) string {
		match := argumentInBracketsRegexp.FindStringSubmatch(s)
		if len(match[1]) > 0 {
			num, err := strconv.ParseInt(match[1], 10, 32)
			if err != nil {
				return s
			}
			if len(args) >= int(num) {
				return args[num-1]
			}
			return ""
		}
		return s
	})
	// $n
	result = argumentRegexp.ReplaceAllStringFunc(result, func(s string) string {
		match := argumentRegexp.FindStringSubmatch(s)
		if len(match[1]) > 0 {
			num, err := strconv.ParseInt(match[1], 10, 32)
			if err != nil {
				return s
			}
			if len(args) >= int(num) {
				return args[num-1]
			}
			return ""
		}
		return s
	})
	// ${[A-Za-z0-9_-]+=[^}]+}
	result = flagInBracketsWithDefaultValueRegexp.ReplaceAllStringFunc(result, func(s string) string {
		match := flagInBracketsWithDefaultValueRegexp.FindStringSubmatch(s)
		flag, err := commandFlags.Find(match[1])
		if err == nil {
			if flag.HasValue {
				valuedFlag, err := flags.Find(match[1])
				if err != nil {
					return match[2]
				}
				return valuedFlag.Value
			} else {
				return fmt.Sprintf("%t", flags.HasFlag(match[1]))
			}
		} else if isUpperCased(match[1]) {
			if env := os.Getenv(match[1]); len(env) != 0 {
				return env
			}
		}
		return match[2]
	})
	// ${[A-Za-z0-9_-]+}
	result = flagInBracketsRegexp.ReplaceAllStringFunc(result, func(s string) string {
		match := flagInBracketsRegexp.FindStringSubmatch(s)
		flag, err := commandFlags.Find(match[1])
		if err == nil {
			if flag.HasValue {
				valuedFlag, err := flags.Find(match[1])
				if err != nil {
					return ""
				}
				return valuedFlag.Value
			} else {
				return fmt.Sprintf("%t", flags.HasFlag(match[1]))
			}
		} else if isUpperCased(match[1]) {
			if env := os.Getenv(match[1]); len(env) != 0 {
				return env
			}
		}
		return s
	})
	// $[A-Za-z0-9_-]+
	result = flagRegexp.ReplaceAllStringFunc(result, func(s string) string {
		match := flagRegexp.FindStringSubmatch(s)
		flag, err := commandFlags.Find(match[1])
		if err == nil {
			if flag.HasValue {
				valuedFlag, err := flags.Find(match[1])
				if err != nil {
					return ""
				}
				return valuedFlag.Value
			} else {
				return fmt.Sprintf("%t", flags.HasFlag(match[1]))
			}
		} else if isUpperCased(match[1]) {
			if env := os.Getenv(match[1]); len(env) != 0 {
				return env
			}
		}
		return s
	})

	return result
}

func isUpperCased(value string) bool {
	for _, rune := range value {
		if unicode.IsLetter(rune) && !unicode.IsUpper(rune) {
			return false
		}
	}
	return true
}
