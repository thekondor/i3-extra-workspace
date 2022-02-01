package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type args struct {
	ExtensionPrefix string
	WsAction        string

	// to keep testing possible
	fatalfFn func(string, ...interface{})
}

func (a args) fatalf(fmt string, args ...interface{}) {
	fn := a.fatalfFn
	if fn == nil {
		fn = log.Fatalf
	}

	fn(fmt, args...)
}

func (a *args) mustParse(args []string) {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-ws-prefix prefix:string] <%s|%s>\n", args[0], toggleWsExtensionAction, flipContainerWsExtensionAction)
		a.fatalf("No argument(s) provided")
		return
	}

	if err := a.parse(args[1:]); err != nil {
		a.fatalf("Failed to parsed arguments: %s", err)
		return
	}
}

func (a *args) parse(args []string) error {
	flagSet := flag.NewFlagSet("", flag.ExitOnError)
	prefixArg := flagSet.String("ws-prefix", defaultExtensionPrefix, "Prefix for extra workspace name")
	if err := flagSet.Parse(args); err != nil {
		return err
	}

	restArgs := flagSet.Args()
	if len(restArgs) != 1 {
		return fmt.Errorf("Action must be provided as a single latest argument")
	}

	a.WsAction = restArgs[0]
	a.ExtensionPrefix = *prefixArg

	return nil
}
