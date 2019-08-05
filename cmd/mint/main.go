package main

import (
	"fmt"
	"os"

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
	a.driver.Commands = []cli.Command{}
}

func (a *app) run(args []string) error {
	return a.driver.Run(args)
}
