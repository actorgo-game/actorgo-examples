package main

import (
	"fmt"
	"os"

	cherryConst "github.com/actorgo-game/actorgo/const"
	"github.com/actorgo-game/examples/demo_cluster/nodes/center"
	"github.com/actorgo-game/examples/demo_cluster/nodes/game"
	"github.com/actorgo-game/examples/demo_cluster/nodes/gate"
	"github.com/actorgo-game/examples/demo_cluster/nodes/master"
	"github.com/actorgo-game/examples/demo_cluster/nodes/web"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "game cluster node",
		Description: "game cluster node examples",
		Commands: []*cli.Command{
			versionCommand(),
			masterCommand(),
			centerCommand(),
			webCommand(),
			gateCommand(),
			gameCommand(),
		},
	}

	_ = app.Run(os.Args)
}

func versionCommand() *cli.Command {
	return &cli.Command{
		Name:      "version",
		Aliases:   []string{"ver", "v"},
		Usage:     "view version",
		UsageText: "game cluster node version",
		Action: func(c *cli.Context) error {
			fmt.Println(cherryConst.Version())
			return nil
		},
	}
}

func masterCommand() *cli.Command {
	return &cli.Command{
		Name:      "1",
		Usage:     "run 1 node",
		UsageText: "node 1 --path=../../config/demo-cluster.json --node=0.0.1.1",
		Flags:     getFlag(),
		Action: func(c *cli.Context) error {
			path, node := getParameters(c)
			master.Run(path, node)
			return nil
		},
	}
}

func centerCommand() *cli.Command {
	return &cli.Command{
		Name:      "2",
		Usage:     "run 2 node",
		UsageText: "node 2 --path=../../config/demo-cluster.json --node=0.0.2.1",
		Flags:     getFlag(),
		Action: func(c *cli.Context) error {
			path, node := getParameters(c)
			center.Run(path, node)
			return nil
		},
	}
}

func webCommand() *cli.Command {
	return &cli.Command{
		Name:      "3",
		Usage:     "run 3 node",
		UsageText: "node 3 --path=../../config/demo-cluster.json --node=0.0.3.1",
		Flags:     getFlag(),
		Action: func(c *cli.Context) error {
			path, node := getParameters(c)
			web.Run(path, node)
			return nil
		},
	}
}

func gateCommand() *cli.Command {
	return &cli.Command{
		Name:      "4",
		Usage:     "run 4 node",
		UsageText: "node 4 --path=../../config/demo-cluster.json --node=0.0.4.1",
		Flags:     getFlag(),
		Action: func(c *cli.Context) error {
			path, node := getParameters(c)
			gate.Run(path, node)
			return nil
		},
	}
}

func gameCommand() *cli.Command {
	return &cli.Command{
		Name:      "5",
		Usage:     "run 5 node",
		UsageText: "node 5 --path=../../config/demo-cluster.json --node=1.1.5.1",
		Flags:     getFlag(),
		Action: func(c *cli.Context) error {
			path, node := getParameters(c)
			game.Run(path, node)
			return nil
		},
	}
}

func getParameters(c *cli.Context) (path, node string) {
	path = c.String("path")
	node = c.String("node")
	return path, node
}

func getFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "path",
			Usage:    "profile config path",
			Required: false,
			Value:    "../../config/demo-cluster.json",
		},
		&cli.StringFlag{
			Name:     "node",
			Usage:    "node id name",
			Required: true,
			Value:    "",
		},
	}
}
