package commands

import (
	"context"

	"github.com/EVRICE/tgeth_alpha/common"
	"github.com/EVRICE/tgeth_alpha/common/hexutil"
	"github.com/EVRICE/tgeth_alpha/crypto"
	"github.com/EVRICE/tgeth_alpha/params"
)

// Web3API provides interfaces for the web3_ RPC commands
type Web3API interface {
	ClientVersion(_ context.Context) (string, error)
	Sha3(_ context.Context, input hexutil.Bytes) hexutil.Bytes
}

type Web3APIImpl struct {
	*BaseAPI
}

// NewWeb3APIImpl returns Web3APIImpl instance
func NewWeb3APIImpl() *Web3APIImpl {
	return &Web3APIImpl{
		BaseAPI: &BaseAPI{},
	}
}

// ClientVersion implements web3_clientVersion. Returns the current client version.
func (api *Web3APIImpl) ClientVersion(_ context.Context) (string, error) {
	return common.MakeName("TurboGeth", params.VersionWithCommit(gitCommit, "")), nil
}

// Sha3 implements web3_sha3. Returns Keccak-256 (not the standardized SHA3-256) of the given data.
func (api *Web3APIImpl) Sha3(_ context.Context, input hexutil.Bytes) hexutil.Bytes {
	return crypto.Keccak256(input)
}

var (
	gitCommit string
)
