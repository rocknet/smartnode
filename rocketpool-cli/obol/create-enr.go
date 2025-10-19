package obol

import (
	"fmt"

	"github.com/rocket-pool/smartnode/shared/services/rocketpool"
	"github.com/urfave/cli"
)

func createEnr(c *cli.Context) error {

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
		return fmt.Errorf("Obol integration is not enabled. Please run 'rocketpool service config' to enable it")
	}

	// Get the container tag
	containerTag := cfg.Obol.GetContainerTag()

	// Get the Obol data directory path
	// This should be ~/.rocketpool/addons/obol based on the docker-compose template
	obolDataDir := fmt.Sprintf("%s/addons/obol", cfg.RocketPoolDirectory)

	fmt.Println("Creating Obol ENR (Ethereum Node Record)...")
	fmt.Printf("Using Charon version: %s\n", containerTag)
	fmt.Printf("Data directory: %s\n", obolDataDir)
	fmt.Println()

	// Docker run command to create the ENR
	// Format: docker run --rm -v <data-dir>:/opt/charon/.charon <image> create enr
	cmd := fmt.Sprintf("docker run --rm -v %s:/opt/charon/.charon %s create enr",
		obolDataDir, containerTag)

	// Execute the docker command
	err = rp.RunDockerCommand(cmd)
	if err != nil {
		return fmt.Errorf("error creating ENR: %w", err)
	}

	fmt.Println()
	fmt.Printf("ENR created successfully! The private key has been saved to %s/charon-enr-private-key\n", obolDataDir)
	fmt.Println("You can view your public ENR with: rocketpool obol show-enr")

	return nil
}
