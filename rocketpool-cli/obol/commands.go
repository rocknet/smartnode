package obol

import (
	"github.com/urfave/cli"

	cliutils "github.com/rocket-pool/smartnode/shared/utils/cli"
)

// Register commands
func RegisterCommands(app *cli.App, name string, aliases []string) {
	app.Commands = append(app.Commands, cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Manage Obol Distributed Validator",
		Subcommands: []cli.Command{

			{
				Name:      "create-enr",
				Usage:     "Create an Ethereum Node Record (ENR) for this node",
				UsageText: "rocketpool obol create-enr",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return createEnr(c)

				},
			},

			{
				Name:      "create-dkg",
				Usage:     "Create a DKG cluster definition",
				UsageText: "rocketpool obol create-dkg",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return createDKG(c)

				},
			},

			{
				Name:      "run-dkg",
				Aliases:   []string{"dkg"},
				Usage:     "Run the Distributed Key Generation ceremony",
				UsageText: "rocketpool obol run-dkg",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return runDKG(c)

				},
			},

			{
				Name:      "show-enr",
				Aliases:   []string{"enr"},
				Usage:     "Display this node's ENR (Ethereum Node Record)",
				UsageText: "rocketpool obol show-enr",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return showEnr(c)

				},
			},

			{
				Name:      "status",
				Aliases:   []string{"s"},
				Usage:     "Get node and megapool status for debugging",
				UsageText: "rocketpool obol status",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return getStatus(c)

				},
			},
		},
	})
}
