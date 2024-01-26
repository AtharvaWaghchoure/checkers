package app

import (
	appparams "github.com/alice/checkers/app/params"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/ignite-hq/cli/ignite/pkg/cosmoscmd"
)

func MakeTestEncodingConfig() cosmoscmd.EncodingConfig {
	encodingConfig := appparams.MakeTestEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}
