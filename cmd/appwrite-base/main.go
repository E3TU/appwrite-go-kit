// Command appwrite-base is the reference binary. It is the only place that
// reads environment variables: env.go parses them into base.Config and the
// library packages stay env-free.
package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	if err := dispatch(os.Args[1:]); err != nil {
		if errors.Is(err, errNoArgs) {
			usage()
			os.Exit(2)
		}
		fmt.Fprintln(os.Stderr, "appwrite-base:", err)
		os.Exit(1)
	}
}

var errNoArgs = errors.New("no subcommand given")

// dispatch routes subcommands to serve and verify.
func dispatch(args []string) error {
	if len(args) == 0 {
		return errNoArgs
	}

	switch args[0] {
	case "serve":
		return serve()
	case "verify":
		return verify()
	case "help", "-h", "--help":
		usage()
		return nil
	default:
		return fmt.Errorf("unknown subcommand %q", args[0])
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `usage: appwrite-base <command>

commands:
  serve    run the example Gin server
  verify   check connectivity and configured resource ids
  help     show this help
`)
}
