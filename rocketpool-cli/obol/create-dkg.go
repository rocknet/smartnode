package obol

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/rocket-pool/smartnode/addons/obol"
	"github.com/rocket-pool/smartnode/shared/services/rocketpool"
	cfgtypes "github.com/rocket-pool/smartnode/shared/types/config"
	"github.com/urfave/cli"
)

func createDKG(c *cli.Context) error {

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

	// Validate cluster name
	clusterName := obolCfg.ClusterName.Value.(string)
	if clusterName == "" {
		return fmt.Errorf("cluster name is not set. Please configure it with 'rocketpool service config'")
	}

	// Validate number of validators
	numValidators := obolCfg.NumValidators.Value.(uint64)
	if numValidators == 0 {
		return fmt.Errorf("number of validators must be greater than 0. Please configure it with 'rocketpool service config'")
	}

	// Validate number of operators
	numOperatorsStr := obolCfg.NumOperators.Value.(string)
	numOperators, err := strconv.Atoi(numOperatorsStr)
	if err != nil {
		return fmt.Errorf("invalid number of operators: %w", err)
	}

	// Check if ENR exists for local node
	obolDataDir := fmt.Sprintf("%s/addons/obol", cfg.RocketPoolDirectory)
	enrPrivateKeyPath := filepath.Join(obolDataDir, "charon-enr-private-key")
	if _, err := os.Stat(enrPrivateKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("local ENR not found. Please create it first with 'rocketpool obol create-enr'")
	}

	// Collect operator ENRs (we need numOperators - 1 since the local node is one of them)
	requiredENRs := numOperators - 1
	operatorENRs := make([]string, 0, requiredENRs)

	// Get ENR values from config
	enrParams := []string{
		obolCfg.OperatorENR1.Value.(string),
		obolCfg.OperatorENR2.Value.(string),
		obolCfg.OperatorENR3.Value.(string),
		obolCfg.OperatorENR4.Value.(string),
		obolCfg.OperatorENR5.Value.(string),
		obolCfg.OperatorENR6.Value.(string),
		obolCfg.OperatorENR7.Value.(string),
		obolCfg.OperatorENR8.Value.(string),
		obolCfg.OperatorENR9.Value.(string),
	}

	for i := 0; i < int(numOperators-1); i++ {
		enr := enrParams[i]
		if enr == "" {
			return fmt.Errorf("member ENR %d is required but not configured. Please configure it with 'rocketpool service config'", i+1)
		}
		// Validate the ENR format and signature
		if err := validateENR(enr); err != nil {
			return fmt.Errorf("member ENR %d is invalid: %w", i+1, err)
		}
		operatorENRs = append(operatorENRs, enr)
	}

	// Get node status to check megapool deployment and fee recipient
	nodeStatus, err := rp.NodeStatus()
	if err != nil {
		return fmt.Errorf("error getting node status: %w", err)
	}

	// Check if megapool is deployed
	if !nodeStatus.MegapoolDeployed {
		return fmt.Errorf("megapool is not deployed. Please deploy your megapool first with 'rocketpool megapool deploy'")
	}

	megapoolAddress := nodeStatus.MegapoolAddress

	// Determine fee recipient address
	var feeRecipientAddress string
	if nodeStatus.FeeRecipientInfo.IsInSmoothingPool {
		// Use smoothing pool address
		feeRecipientAddress = nodeStatus.FeeRecipientInfo.SmoothingPoolAddress.Hex()
		fmt.Printf("Using smoothing pool as fee recipient: %s\n", feeRecipientAddress)
	} else {
		// Use fee distributor address
		feeRecipientAddress = nodeStatus.FeeRecipientInfo.FeeDistributorAddress.Hex()
		fmt.Printf("Using fee distributor as fee recipient: %s\n", feeRecipientAddress)
	}

	// Get local ENR by running show-enr command
	fmt.Println("Reading local ENR...")
	containerTag := cfg.Obol.GetContainerTag()

	// Get the full container name for the running container check
	prefix, err := rp.GetContainerPrefix()
	if err != nil {
		return fmt.Errorf("error getting container prefix: %w", err)
	}
	fullContainerName := fmt.Sprintf("%s_%s", prefix, cfg.Obol.GetContainerName())

	// Get local ENR using the shared function
	localENR, err := GetLocalENR(rp, containerTag, obolDataDir, fullContainerName)
	if err != nil {
		return fmt.Errorf("error reading local ENR: %w", err)
	}

	// Validate the local ENR
	if err := validateENR(localENR); err != nil {
		return fmt.Errorf("local ENR is invalid: %w", err)
	}

	// Add local ENR to the operator list
	allOperatorENRs := append([]string{localENR}, operatorENRs...)
	operatorENRString := strings.Join(allOperatorENRs, ",")

	// Display summary
	fmt.Println()
	fmt.Println("=== Cluster Configuration ===")
	fmt.Printf("Cluster Name: %s\n", clusterName)
	fmt.Printf("Number of Validators: %d\n", numValidators)
	fmt.Printf("Number of Operators: %d\n", numOperators)
	fmt.Printf("Megapool Address: %s\n", megapoolAddress.Hex())
	fmt.Printf("Fee Recipient: %s\n", feeRecipientAddress)
	fmt.Println()
	fmt.Println("Operator ENRs:")
	fmt.Printf("  1. %s (local)\n", localENR)
	for i, enr := range operatorENRs {
		fmt.Printf("  %d. %s\n", i+2, enr)
	}
	fmt.Println()

	// Get the network name for Obol
	network := cfg.Smartnode.Network.Value.(cfgtypes.Network)
	obolNetwork := GetCharonNetwork(network)

	// Build the docker command with --publish flag to get a shareable URL
	cmd := fmt.Sprintf("docker run --rm -v %s:/opt/charon/.charon %s create dkg --name=%s --num-validators=%d --withdrawal-addresses=%s --fee-recipient-addresses=%s --operator-enrs=%s --network=%s --publish",
		obolDataDir,
		containerTag,
		shellescape.Quote(clusterName),
		numValidators,
		shellescape.Quote(megapoolAddress.Hex()),
		shellescape.Quote(feeRecipientAddress),
		shellescape.Quote(operatorENRString),
		obolNetwork,
	)

	fmt.Println("Creating DKG cluster definition and publishing to Obol...")
	fmt.Println()

	// Execute the docker command
	err = rp.RunDockerCommand(cmd)
	if err != nil {
		return fmt.Errorf("error creating cluster: %w", err)
	}

	fmt.Println()
	fmt.Printf("Cluster definition created successfully!\n")
	fmt.Printf("The cluster definition has been saved to %s/cluster-definition.json\n", obolDataDir)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("1. Share the cluster invite URL (displayed above) with all cluster members")
	fmt.Println("2. Each operator should accept the invite and run the DKG ceremony with 'rocketpool obol dkg'")
	fmt.Println("3. After DKG completes, import the validator keys")

	return nil
}

// validateENR validates an ENR string by parsing it and verifying its cryptographic signature
func validateENR(enrStr string) error {
	// Parse and validate the ENR (this also verifies the signature)
	// enode.ValidSchemes includes all supported identity schemes (secp256k1-keccak, etc.)
	_, err := enode.Parse(enode.ValidSchemes, enrStr)
	if err != nil {
		return fmt.Errorf("invalid ENR: %w", err)
	}

	return nil
}
