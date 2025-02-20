package backend

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	rpctypes "github.com/evmos/ethermint/rpc/types"
)

// ChainID is the chain id for the current chain config.
func (b *Backend) ChainID() (*hexutil.Big, error) {
	chainID, err := sdk.ParseChainID(b.clientCtx.ChainID)
	if err != nil {
		panic(err)
	}

	return (*hexutil.Big)(chainID), nil
}

// CurrentHeader returns the latest block header
func (b *Backend) CurrentHeader() *ethtypes.Header {
	header, _ := b.HeaderByNumber(rpctypes.EthLatestBlockNumber)
	return header
}
