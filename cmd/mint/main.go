package main

import "github.com/urfave/cli"

func main() {}

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
