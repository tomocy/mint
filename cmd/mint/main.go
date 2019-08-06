package main

import (
	"fmt"
	"os"

	"github.com/tomocy/mint/cmd/mint/client"
	"github.com/urfave/cli"
)

func main() {
	a := newApp()
	if err := a.run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run app: %s\n", err)
		os.Exit(1)
	}
}

func newApp() *app {
	a := new(app)
	a.setUp()

	return a
}

type app struct {
	driver *cli.App
}

func (a *app) setUp() {
	a.driver = cli.NewApp()
	a.setBasic()
	a.setCommands()
}

func (a *app) setBasic() {
	a.driver.Name = "mint"
}

func (a *app) setCommands() {
	a.driver.Commands = []cli.Command{
		cli.Command{
			Name:   "cli",
			Action: a.runCLI,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name: "p",
				},
			},
		},
	}
}

func (a *app) runCLI(ctx *cli.Context) error {
	client := new(client.CLI)
	isPoll := ctx.Bool("p")
	if isPoll {
		return client.PoleHomeTweets()
	}

	return client.FetchHomeTweets()
}

func (a *app) run(args []string) error {
	return a.driver.Run(args)
}
