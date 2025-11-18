package main

import (
	"fmt"

	"github.com/chunkgar/gokit/app"
	cliflag "github.com/marmotedu/component-base/pkg/cli/flag"
)

type SubOptions struct {
	Msg   string
	Times int
}

func (o *SubOptions) Flags() (fss cliflag.NamedFlagSets) {
	fs := fss.FlagSet("subcommand")
	fs.StringVar(&o.Msg, "msg", o.Msg, "message to print")
	fs.IntVar(&o.Times, "times", o.Times, "number of times to print")
	return fss
}

func (o *SubOptions) Validate() []error { return nil }

func main() {
	opts := &SubOptions{Msg: "Hello, World from subcommand!", Times: 1}

	subcmd := app.NewCommand(
		"subcommand",
		"subcommand",
		app.WithCommandOptions(opts),
		app.WithCommandRunFunc(func(args []string) error {
			for i := 0; i < opts.Times; i++ {
				fmt.Println(opts.Msg)
			}
			return nil
		}),
	)

	app.NewApp(
		"test name",
		"test",
		app.WithNoConfig(),
		app.WithRunFunc(func(basename string) error {
			fmt.Println("Hello, World!")
			return nil
		}),
		app.WithSubCommand(subcmd),
	).Run()
}
