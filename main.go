package main

import (
	"fmt"
	"os"

	"github.com/rancher/config-modifier/pkg/config"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	Version          = "v0.0.1"
	GitCommit        = "HEAD"
	hostPath         string
	configPath       string
	nodeLabels       cli.StringSlice
	preservedEntries cli.StringSlice
)

func main() {
	app := cli.NewApp()
	app.Name = "config modifier"
	app.Version = fmt.Sprintf("%s (%s)", Version, GitCommit)
	app.Usage = "Modify config file in k3s or rke2 nodes"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "host-path",
			EnvVar:      "HOST_PATH",
			Destination: &hostPath,
		},
		cli.StringFlag{
			Name:        "config-path",
			EnvVar:      "CONFIG_PATH",
			Destination: &configPath,
		},
		cli.StringSliceFlag{
			Name:   "node-labels",
			EnvVar: "NODE_LABELS",
			Value:  &nodeLabels,
		},
		cli.StringSliceFlag{
			Name:   "preserved-entries",
			EnvVar: "PRESERVED_ENTRIES",
			Value:  &preservedEntries,
		},
	}
	app.Action = run

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	err := config.PlaceConfigFile(hostPath, configPath, nodeLabels, preservedEntries)
	return err
}
