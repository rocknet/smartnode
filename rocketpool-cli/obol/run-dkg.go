package obol

import (
	"fmt"

	"github.com/rocket-pool/smartnode/addons/obol"
	"github.com/rocket-pool/smartnode/shared/services/rocketpool"
	"github.com/urfave/cli"
)

func runDKG(c *cli.Context) error {

	// Get RP client
	rp := rocketpool.NewClientFromCtx(c)
	defer rp.Close()

	// Get the config
	cfg, _, err := rp.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	// Check if Obol is enabled
	if !cfg.Obol.GetEnabledParameter().Value.(bool) {
		return fmt.Errorf("obol integration is not enabled. Please run 'rocketpool service config' to enable it")
	}

	// Get Obol config
	obolCfg := cfg.Obol.GetConfig().(*obol.ObolConfig)

	// Determine the data directory
	obolDataDir := fmt.Sprintf("%s/addons/obol", cfg.RocketPoolDirectory)

	// Get container tag
	containerTag := cfg.Obol.GetContainerTag()

	var cmd string

	// Get the cluster definition URL
	clusterDefURL := obolCfg.ClusterDefinitionURL.Value.(string)
	if clusterDefURL == "" {
		return fmt.Errorf("cluster definition URL is not set. Please configure it with 'rocketpool service config'")
	}

	fmt.Println("=== Running DKG Ceremony (Cluster Member) ===")
	fmt.Println()
	fmt.Printf("Cluster Definition URL: %s\n", clusterDefURL)
	fmt.Println()
	fmt.Println("This will download the cluster definition and run the Distributed Key Generation ceremony.")
	fmt.Println()

	// Build the docker command for cluster member
	cmd = fmt.Sprintf("docker run --rm -v %s:/opt/charon/.charon %s dkg --definition-file=%s",
		obolDataDir,
		containerTag,
		clusterDefURL,
	)

	fmt.Println("Starting DKG ceremony...")
	fmt.Println()

	// Execute the docker command
	err = rp.RunDockerCommand(cmd)
	if err != nil {
		return fmt.Errorf("error running DKG ceremony: %w", err)
	}

	fmt.Println()
	fmt.Println("DKG ceremony completed successfully!")
	fmt.Printf("Validator keys have been generated in %s/validator_keys/\n", obolDataDir)
	fmt.Println()
	fmt.Println("Next steps:")
	// TODO: I think the process here will be the obol deposit command, which will import the keyshares
	// into the validator clients.
	// fmt.Println("1. Import the validator keys with the appropriate import command")
	// fmt.Println("2. Start your validator client to begin attesting")
	// fmt.Println("3. Coordinate with other cluster members to ensure all validators are running")
	fmt.Println("1. Finish the code to add obol deposit, which will load validator key shares into the validator client")
	fmt.Println("2. ...")
	fmt.Println("3. Profit")

	return nil
}
