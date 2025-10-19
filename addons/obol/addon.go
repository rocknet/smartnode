package obol

import (
	"github.com/rocket-pool/smartnode/shared/types/addons"
	cfgtypes "github.com/rocket-pool/smartnode/shared/types/config"
)

type Obol struct {
	cfg *ObolConfig `yaml:"config,omitempty"`
}

func NewObol() addons.SmartnodeAddon {
	return &Obol{
		cfg: NewConfig(),
	}
}

func (o *Obol) GetName() string {
	return "Obol DVT"
}

func (o *Obol) GetDescription() string {
	return `Obol Distributed Validator Technology

Obol enables you to run Ethereum validators as a cluster with multiple operators, providing increased resilience and security through distributed key generation and threshold signing.

For more information, see https://obol.org`
}

func (o *Obol) GetConfig() cfgtypes.Config {
	return o.cfg
}

func (o *Obol) GetContainerName() string {
	return CharonContainerName
}

func (o *Obol) GetContainerTag() string {
	return o.cfg.ContainerTag.Value.(string)
}

func (o *Obol) GetEnabledParameter() *cfgtypes.Parameter {
	return &o.cfg.Enabled
}

func (o *Obol) GetP2PPort() uint16 {
	return o.cfg.P2PPort.Value.(uint16)
}
