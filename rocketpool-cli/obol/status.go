package obol

import (
	"encoding/json"
	"fmt"

	"github.com/rocket-pool/smartnode/shared/services/rocketpool"
	"github.com/urfave/cli"
)

func getStatus(c *cli.Context) error {

	// Get RP client
	rp := rocketpool.NewClientFromCtx(c)
	defer rp.Close()

	fmt.Println("=== Node Status ===")
	nodeStatus, err := rp.NodeStatus()
	if err != nil {
		fmt.Printf("Error getting node status: %v\n", err)
	} else {
		// Pretty print the entire struct as JSON
		jsonBytes, err := json.MarshalIndent(nodeStatus, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling node status: %v\n", err)
		} else {
			fmt.Println(string(jsonBytes))
		}
	}

	fmt.Println()
	fmt.Println("=== Megapool Status ===")
	megapoolStatus, err := rp.MegapoolStatus()
	if err != nil {
		fmt.Printf("Error getting megapool status: %v\n", err)
	} else {
		// Pretty print the entire struct as JSON
		jsonBytes, err := json.MarshalIndent(megapoolStatus, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling megapool status: %v\n", err)
		} else {
			fmt.Println(string(jsonBytes))
		}
	}

	return nil
}
