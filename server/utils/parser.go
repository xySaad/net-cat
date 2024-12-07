package utils

import (
	"os"
)

func ParseFlags() (int, []string) {
	if len(os.Args) == 1 {
		return 0, os.Args
	}
	result := 0
	args := []string{}
	seen := make(map[string]bool)
	for _, arg := range os.Args {
		if arg[0] == '-' && !seen[arg] {
			seen[arg] = true
			switch arg {
			case "-c":
				result += 1
			case "-u":
				result += 2
			default:
				return -1, nil
			}
		} else if arg[0] != '-' {
			args = append(args, arg)
		}
	}
	return result, args
}
