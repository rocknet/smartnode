package obol

import (
	"fmt"
	"strings"

	"github.com/rocket-pool/smartnode/shared/services/rocketpool"
	cfgtypes "github.com/rocket-pool/smartnode/shared/types/config"
)

// GetLocalENR retrieves the local node's ENR, trying the running container first,
// then falling back to reading from the private key file using a temporary container
func GetLocalENR(rp *rocketpool.Client, containerTag, obolDataDir, containerName string) (string, error) {
	// Try to get ENR from running container first
	cmd := fmt.Sprintf("docker exec %s charon enr", containerName)

	output, err := rp.RunCommandAndReturnOutput(cmd)
	if err != nil {
		// If container is not running, try reading from the private key file
		// Run charon enr in a temporary container with the data directory mounted
		cmd = fmt.Sprintf("docker run --rm -v %s:/opt/charon/.charon %s enr",
			obolDataDir, containerTag)

		output, err = rp.RunCommandAndReturnOutput(cmd)
		if err != nil {
			return "", fmt.Errorf("error retrieving ENR: %w\n\nPlease ensure you have created an ENR first with 'rocketpool obol create-enr'", err)
		}
	}

	return strings.TrimSpace(string(output)), nil
}

// GetCharonNetwork returns the network name that Charon expects based on the Rocket Pool network configuration.
// Devnet and Testnet both map to "hoodi" (the current long-lived testnet).
func GetCharonNetwork(network cfgtypes.Network) string {
	switch network {
	case cfgtypes.Network_Mainnet:
		return "mainnet"
	case cfgtypes.Network_Devnet, cfgtypes.Network_Testnet:
		return "hoodi"
	default:
		return "mainnet"
	}
}
