package obol

import (
	"github.com/rocket-pool/smartnode/shared/types/config"
)

const (
	ContainerID_Charon  config.ContainerID = "charon"
	CharonContainerName string             = "addon_charon"
)

// Cluster role enum
type ClusterRole string

const (
	ClusterRole_Creator ClusterRole = "creator"
	ClusterRole_Member  ClusterRole = "member"
)

// Configuration for the Obol add-on
type ObolConfig struct {
	Title string `yaml:"-"`

	// The parent config
	parent config.Config `yaml:"-"`

	// Enabled
	Enabled config.Parameter `yaml:"enabled,omitempty"`

	// Cluster Role
	ClusterRole config.Parameter `yaml:"clusterRole,omitempty"`

	// Container tag
	ContainerTag config.Parameter `yaml:"containerTag,omitempty"`

	// P2P Port
	P2PPort config.Parameter `yaml:"p2pPort,omitempty"`

	// Cluster Creator settings
	ClusterName   config.Parameter `yaml:"clusterName,omitempty"`
	NumValidators config.Parameter `yaml:"numValidators,omitempty"`
	NumOperators  config.Parameter `yaml:"numOperators,omitempty"`
	OperatorENR1  config.Parameter `yaml:"operatorENR1,omitempty"`
	OperatorENR2  config.Parameter `yaml:"operatorENR2,omitempty"`
	OperatorENR3  config.Parameter `yaml:"operatorENR3,omitempty"`
	OperatorENR4  config.Parameter `yaml:"operatorENR4,omitempty"`
	OperatorENR5  config.Parameter `yaml:"operatorENR5,omitempty"`
	OperatorENR6  config.Parameter `yaml:"operatorENR6,omitempty"`
	OperatorENR7  config.Parameter `yaml:"operatorENR7,omitempty"`
	OperatorENR8  config.Parameter `yaml:"operatorENR8,omitempty"`
	OperatorENR9  config.Parameter `yaml:"operatorENR9,omitempty"`

	// Cluster Member settings
	ClusterDefinitionURL config.Parameter `yaml:"clusterDefinitionURL,omitempty"`
}

// Creates a new configuration instance
func NewConfig() *ObolConfig {
	return &ObolConfig{
		Title: "Obol Settings",

		Enabled: config.Parameter{
			ID:                 "enabled",
			Name:               "Enabled",
			Description:        "Enable Obol Distributed Validator Technology (DVT)\n\nObol allows you to run validators as a cluster with multiple operators for increased resilience and security.\n\nVisit obol.org for more information.",
			Type:               config.ParameterType_Bool,
			Default:            map[config.Network]interface{}{config.Network_All: false},
			AffectsContainers:  []config.ContainerID{config.ContainerID_Validator, ContainerID_Charon},
			CanBeBlank:         false,
			OverwriteOnUpgrade: false,
		},

		ClusterRole: config.Parameter{
			ID:                 "clusterRole",
			Name:               "Cluster Role",
			Description:        "Select your role in the Obol cluster. Choose 'Cluster Creator' if you are creating a new cluster definition and coordinating the DKG ceremony. Choose 'Cluster Member Node' if you are joining an existing cluster.",
			Type:               config.ParameterType_Choice,
			Default:            map[config.Network]interface{}{config.Network_All: ClusterRole_Creator},
			AffectsContainers:  []config.ContainerID{config.ContainerID_Validator, ContainerID_Charon},
			CanBeBlank:         false,
			OverwriteOnUpgrade: false,
			Options: []config.ParameterOption{{
				Name:        "Cluster Creator",
				Description: "Create a new cluster and coordinate the DKG ceremony",
				Value:       ClusterRole_Creator,
			}, {
				Name:        "Cluster Member Node",
				Description: "Join an existing cluster as a member operator",
				Value:       ClusterRole_Member,
			}},
		},

		ContainerTag: config.Parameter{
			ID:                 "containerTag",
			Name:               "Container Tag",
			Description:        "The tag name of the Obol Charon container you want to use from Docker Hub.",
			Type:               config.ParameterType_String,
			Default:            map[config.Network]interface{}{config.Network_All: "obolnetwork/charon:v1.7.0"},
			AffectsContainers:  []config.ContainerID{ContainerID_Charon},
			CanBeBlank:         false,
			OverwriteOnUpgrade: true,
		},

		P2PPort: config.Parameter{
			ID:                 "p2pPort",
			Name:               "P2P Port",
			Description:        "The TCP port for Charon peer-to-peer communication. Change this if you have multiple Charon instances on the same machine.",
			Type:               config.ParameterType_Uint16,
			Default:            map[config.Network]interface{}{config.Network_All: uint16(3610)},
			AffectsContainers:  []config.ContainerID{ContainerID_Charon},
			CanBeBlank:         false,
			OverwriteOnUpgrade: false,
		},

		ClusterName: config.Parameter{
			ID:                 "clusterName",
			Name:               "Cluster Name",
			Description:        "A friendly name for your Obol cluster.",
			Type:               config.ParameterType_String,
			Default:            map[config.Network]interface{}{config.Network_All: ""},
			AffectsContainers:  []config.ContainerID{ContainerID_Charon},
			CanBeBlank:         false,
			OverwriteOnUpgrade: false,
			MaxLength:          64,
		},

		NumValidators: config.Parameter{
			ID:                 "numValidators",
			Name:               "Number of Validators",
			Description:        "The number of distributed validators to create in this cluster.",
			Type:               config.ParameterType_Uint,
			Default:            map[config.Network]interface{}{config.Network_All: uint64(1)},
			AffectsContainers:  []config.ContainerID{ContainerID_Charon},
			CanBeBlank:         false,
			OverwriteOnUpgrade: false,
		},

		NumOperators: config.Parameter{
			ID:                 "numOperators",
			Name:               "Number of Operators",
			Description:        "The total number of operators in the cluster (including yourself).",
			Type:               config.ParameterType_Choice,
			Default:            map[config.Network]interface{}{config.Network_All: "3"},
			AffectsContainers:  []config.ContainerID{ContainerID_Charon},
			CanBeBlank:         false,
			OverwriteOnUpgrade: false,
			Options: []config.ParameterOption{
				{Name: "3", Description: "3 operators", Value: "3"},
				{Name: "4", Description: "4 operators", Value: "4"},
				{Name: "5", Description: "5 operators", Value: "5"},
				{Name: "6", Description: "6 operators", Value: "6"},
				{Name: "7", Description: "7 operators", Value: "7"},
				{Name: "8", Description: "8 operators", Value: "8"},
				{Name: "9", Description: "9 operators", Value: "9"},
				{Name: "10", Description: "10 operators", Value: "10"},
			},
		},

		OperatorENR1: config.Parameter{
			ID:                 "operatorENR1",
			Name:               "Member ENR 1",
			Description:        "The Ethereum Node Record (ENR) for member operator 1.",
			Type:               config.ParameterType_String,
			Default:            map[config.Network]interface{}{config.Network_All: ""},
			AffectsContainers:  []config.ContainerID{ContainerID_Charon},
			CanBeBlank:         false,
			OverwriteOnUpgrade: false,
		},

		OperatorENR2: config.Parameter{
			ID:                 "operatorENR2",
			Name:               "Member ENR 2",
			Description:        "The Ethereum Node Record (ENR) for member operator 2.",
			Type:               config.ParameterType_String,
			Default:            map[config.Network]interface{}{config.Network_All: ""},
			AffectsContainers:  []config.ContainerID{ContainerID_Charon},
			CanBeBlank:         false,
			OverwriteOnUpgrade: false,
		},

		OperatorENR3: config.Parameter{
			ID:                 "operatorENR3",
			Name:               "Member ENR 3",
			Description:        "The Ethereum Node Record (ENR) for member operator 3.",
			Type:               config.ParameterType_String,
			Default:            map[config.Network]interface{}{config.Network_All: ""},
			AffectsContainers:  []config.ContainerID{ContainerID_Charon},
			CanBeBlank:         true,
			OverwriteOnUpgrade: false,
		},

		OperatorENR4: config.Parameter{
			ID:                 "operatorENR4",
			Name:               "Member ENR 4",
			Description:        "The Ethereum Node Record (ENR) for member operator 4.",
			Type:               config.ParameterType_String,
			Default:            map[config.Network]interface{}{config.Network_All: ""},
			AffectsContainers:  []config.ContainerID{ContainerID_Charon},
			CanBeBlank:         true,
			OverwriteOnUpgrade: false,
		},

		OperatorENR5: config.Parameter{
			ID:                 "operatorENR5",
			Name:               "Member ENR 5",
			Description:        "The Ethereum Node Record (ENR) for member operator 5.",
			Type:               config.ParameterType_String,
			Default:            map[config.Network]interface{}{config.Network_All: ""},
			AffectsContainers:  []config.ContainerID{ContainerID_Charon},
			CanBeBlank:         true,
			OverwriteOnUpgrade: false,
		},

		OperatorENR6: config.Parameter{
			ID:                 "operatorENR6",
			Name:               "Member ENR 6",
			Description:        "The Ethereum Node Record (ENR) for member operator 6.",
			Type:               config.ParameterType_String,
			Default:            map[config.Network]interface{}{config.Network_All: ""},
			AffectsContainers:  []config.ContainerID{ContainerID_Charon},
			CanBeBlank:         true,
			OverwriteOnUpgrade: false,
		},

		OperatorENR7: config.Parameter{
			ID:                 "operatorENR7",
			Name:               "Member ENR 7",
			Description:        "The Ethereum Node Record (ENR) for member operator 7.",
			Type:               config.ParameterType_String,
			Default:            map[config.Network]interface{}{config.Network_All: ""},
			AffectsContainers:  []config.ContainerID{ContainerID_Charon},
			CanBeBlank:         true,
			OverwriteOnUpgrade: false,
		},

		OperatorENR8: config.Parameter{
			ID:                 "operatorENR8",
			Name:               "Member ENR 8",
			Description:        "The Ethereum Node Record (ENR) for member operator 8.",
			Type:               config.ParameterType_String,
			Default:            map[config.Network]interface{}{config.Network_All: ""},
			AffectsContainers:  []config.ContainerID{ContainerID_Charon},
			CanBeBlank:         true,
			OverwriteOnUpgrade: false,
		},

		OperatorENR9: config.Parameter{
			ID:                 "operatorENR9",
			Name:               "Member ENR 9",
			Description:        "The Ethereum Node Record (ENR) for member operator 9.",
			Type:               config.ParameterType_String,
			Default:            map[config.Network]interface{}{config.Network_All: ""},
			AffectsContainers:  []config.ContainerID{ContainerID_Charon},
			CanBeBlank:         true,
			OverwriteOnUpgrade: false,
		},

		ClusterDefinitionURL: config.Parameter{
			ID:                 "clusterDefinitionURL",
			Name:               "Cluster Definition URL",
			Description:        "The URL to download the cluster definition file. This should be the invite URL provided by the cluster creator.",
			Type:               config.ParameterType_String,
			Default:            map[config.Network]interface{}{config.Network_All: ""},
			AffectsContainers:  []config.ContainerID{ContainerID_Charon},
			CanBeBlank:         false,
			OverwriteOnUpgrade: false,
		},
	}
}

// Get the parameters for this config
func (cfg *ObolConfig) GetParameters() []*config.Parameter {
	return []*config.Parameter{
		&cfg.Enabled,
		&cfg.ClusterRole,
		&cfg.ContainerTag,
		&cfg.P2PPort,
		&cfg.ClusterName,
		&cfg.NumValidators,
		&cfg.NumOperators,
		&cfg.OperatorENR1,
		&cfg.OperatorENR2,
		&cfg.OperatorENR3,
		&cfg.OperatorENR4,
		&cfg.OperatorENR5,
		&cfg.OperatorENR6,
		&cfg.OperatorENR7,
		&cfg.OperatorENR8,
		&cfg.OperatorENR9,
		&cfg.ClusterDefinitionURL,
	}
}

// Get the sections for this config
func (cfg *ObolConfig) GetSubconfigs() map[string]config.Config {
	return map[string]config.Config{}
}

// The title for the config
func (cfg *ObolConfig) GetConfigTitle() string {
	return cfg.Title
}
