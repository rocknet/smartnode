package obol

import (
	"fmt"

	"github.com/rocket-pool/smartnode/shared/services/rocketpool"
	"github.com/urfave/cli"
)

func showEnr(c *cli.Context) error {

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

	// Get the container name
	containerName := cfg.Obol.GetContainerName()
	prefix, err := rp.GetContainerPrefix()
	if err != nil {
		return fmt.Errorf("error getting container prefix: %w", err)
	}
	fullContainerName := fmt.Sprintf("%s_%s", prefix, containerName)

	fmt.Printf("Attempting to retrieve Obol ENR (Ethereum Node Record) from running container %s...\n", fullContainerName)
	fmt.Println()

	// Get the container tag and data directory
	containerTag := cfg.Obol.GetContainerTag()
	obolDataDir := fmt.Sprintf("%s/addons/obol", cfg.RocketPoolDirectory)

	// Get the ENR using the shared function
	enr, err := GetLocalENR(rp, containerTag, obolDataDir, fullContainerName)
	if err != nil {
		return err
	}

	fmt.Println(enr)
	fmt.Println()
	fmt.Println("Share this ENR with the other operators in your DV cluster.")

	return nil
}
