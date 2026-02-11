package command

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func FindCommand(commands []Command, cmd string) *Command {
	for _, command := range commands {
		if command.Is(cmd) {
			return &command
		}
	}
	return nil
}

// $*
// $n ${n} ${n=value}
// -f --flag
// $flag ${flag} ${flag=value}
// ${-f} ${--flag}
// ${-f:--flag} ${--flag:-f}
// ${-f=value} ${--flag=value}
// ${-f:--flag=value} ${--flag:-f=value}
func formatInsturction(instruction string, args []string, flags map[string]string) string {
	// $*
	result := strings.ReplaceAll(instruction, "$*", strings.Join(args, " "))
	// ${n}
	result = regexp.MustCompile(`\$\{\d+\}`).ReplaceAllStringFunc(result, func(s string) string {
		num, err := strconv.ParseInt(s[2:len(s)-1], 10, 32)
		if err != nil {
			return s
		}
		if len(args) >= int(num) {
			return args[num-1]
		}
		return ""
	})
	// ${n=value}
	result = regexp.MustCompile(`\$\{\d+=[^\}]+\}`).ReplaceAllStringFunc(result, func(s string) string {
		split := strings.SplitN(s[2:len(s)-1], "=", 2)
		fmt.Printf("%v\n", split)
		num, err := strconv.ParseInt(split[0], 10, 32)
		if err != nil {
			return s
		}
		if len(args) >= int(num) {
			return args[num-1]
		}
		return split[1]
	})
	// ${-f:--flag=value} ${--flag:-f=value}
	result = regexp.MustCompile(`\$\{--?[a-zA-Z0-9_-]+:--?[a-zA-Z0-9_-]+=[^\}]+\}`).ReplaceAllStringFunc(result, func(s string) string {
		values := s[2 : len(s)-1]
		split := strings.SplitN(values, ":", 2)
		first := strings.TrimLeft(split[0], "-")
		secondWithDefaultValue := strings.TrimLeft(split[1], "-")
		splitSecond := strings.SplitN(secondWithDefaultValue, "=", 2)
		second := splitSecond[0]
		value, ok := flags[first]
		if ok {
			return value
		}
		value, ok = flags[second]
		if ok {
			return value
		}
		return splitSecond[1]
	})
	// ${-f=value} ${--flag=value}
	result = regexp.MustCompile(`\$\{--?[a-zA-Z0-9_-]+=[^\}]+\}`).ReplaceAllStringFunc(result, func(s string) string {
		values := s[2 : len(s)-1]
		split := strings.SplitN(values, "=", 2)
		flag := strings.TrimLeft(split[0], "-")
		value, ok := flags[flag]
		if ok {
			return value
		}
		return split[1]
	})
	// ${-f:--flag} ${--flag:-f}
	result = regexp.MustCompile(`\$\{--?[a-zA-Z0-9_-]+:--?[a-zA-Z0-9_-]+\}`).ReplaceAllStringFunc(result, func(s string) string {
		values := s[2 : len(s)-1]
		split := strings.SplitN(values, ":", 2)
		first := strings.TrimLeft(split[0], "-")
		second := strings.TrimLeft(split[1], "-")
		value, ok := flags[first]
		if ok {
			return value
		}
		value, ok = flags[second]
		if ok {
			return value
		}
		return "false"
	})
	// ${-f} ${--flag}
	result = regexp.MustCompile(`\$\{--?[a-zA-Z0-9_-]+\}`).ReplaceAllStringFunc(result, func(s string) string {
		flag := strings.TrimLeft(s[2:len(s)-1], "-")
		value, ok := flags[flag]
		if ok {
			return value
		}
		return "false"
	})
	// -f --flag
	result = regexp.MustCompile(`^--?[a-zA-Z0-9_-]+`).ReplaceAllStringFunc(result, func(s string) string {
		flag := strings.TrimLeft(s, "-")
		value, ok := flags[flag]
		if ok {
			return value
		}
		return "false"
	})
	result = regexp.MustCompile(` --?[a-zA-Z0-9_-]+`).ReplaceAllStringFunc(result, func(s string) string {
		flag := strings.TrimLeft(s[1:], "-")
		value, ok := flags[flag]
		if ok {
			return fmt.Sprintf(" %s", value)
		}
		return " false"
	})
	// $flag
	result = regexp.MustCompile(`\$[a-zA-Z0-9_-]+`).ReplaceAllStringFunc(result, func(s string) string {
		flag := s[1:]
		value, ok := flags[flag]
		if ok {
			return value
		}
		if isUpperCased(flag) {
			if env := os.Getenv(flag); len(env) != 0 {
				return env
			}
		}
		return "false"
	})
	// ${flag}
	result = regexp.MustCompile(`\$\{[a-zA-Z0-9_-]+\}`).ReplaceAllStringFunc(result, func(s string) string {
		flag := s[2 : len(s)-1]
		value, ok := flags[flag]
		if ok {
			return value
		}
		if isUpperCased(flag) {
			if env := os.Getenv(flag); len(env) != 0 {
				return env
			}
		}
		return "false"
	})
	// ${flag=value}
	result = regexp.MustCompile(`\$\{[a-zA-Z0-9_-]+=[^\}]+\}`).ReplaceAllStringFunc(result, func(s string) string {
		split := strings.SplitN(s[2:len(s)-1], "=", 2)
		value, ok := flags[split[0]]
		if ok {
			return value
		}
		return split[1]
	})

	return result
}

func parseArgsAndFlags(input []string) ([]string, map[string]string) {
	args := make([]string, 0, len(input))
	flags := make(map[string]string, len(input))

	valuedStringFlagRe := regexp.MustCompile(`^--?[A-Za-z0-9_-]+=.+$`)
	valuedFlagRe := regexp.MustCompile(`^--?[A-Za-z0-9_-]+=\S+$`)
	booleanFlagRe := regexp.MustCompile(`^--?[A-Za-z0-9_-]+$`)
	hasDoubleDash := false
	for _, arg := range input {
		if hasDoubleDash {
			args = append(args, arg)
		} else if arg == "--" {
			hasDoubleDash = true
		} else if valuedStringFlagRe.MatchString(arg) {
			split := strings.SplitN(arg, "=", 2)
			flag := strings.TrimLeft(split[0], "-")
			flags[flag] = split[1]
		} else if valuedFlagRe.MatchString(arg) {
			split := strings.SplitN(arg, "=", 2)
			flag := strings.TrimLeft(split[0], "-")
			flags[flag] = split[1]
		} else if booleanFlagRe.MatchString(arg) {
			flag := strings.TrimLeft(arg, "-")
			flags[flag] = "true"
		} else {
			args = append(args, arg)
		}
	}
	return args, flags
}

func isUpperCased(value string) bool {
	for _, rune := range value {
		if unicode.IsLetter(rune) && !unicode.IsUpper(rune) {
			return false
		}
	}
	return true
}
