package adapter

import (
	"math/big"

	"github.com/EVRICE/tgeth_alpha/common"
	"github.com/EVRICE/tgeth_alpha/consensus"
	"github.com/EVRICE/tgeth_alpha/core/rawdb"
	"github.com/EVRICE/tgeth_alpha/core/state"
	"github.com/EVRICE/tgeth_alpha/core/types"
	"github.com/EVRICE/tgeth_alpha/ethdb"
	"github.com/EVRICE/tgeth_alpha/params"
	"github.com/EVRICE/tgeth_alpha/rpc"
)

type chainContext struct {
	db ethdb.Database
}

func NewChainContext(db ethdb.Database) *chainContext {
	return &chainContext{
		db: db,
	}
}

type powEngine struct {
}

func (c *powEngine) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {

	panic("must not be called")
}
func (c *powEngine) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (func(), <-chan error) {
	panic("must not be called")
}
func (c *powEngine) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	panic("must not be called")
}
func (c *powEngine) VerifySeal(chain consensus.ChainHeaderReader, header *types.Header) error {
	panic("must not be called")
}
func (c *powEngine) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	panic("must not be called")
}
func (c *powEngine) Finalize(chainConfig *params.ChainConfig, header *types.Header, state *state.IntraBlockState, txs []*types.Transaction, uncles []*types.Header) {
	panic("must not be called")
}
func (c *powEngine) FinalizeAndAssemble(chainConfig *params.ChainConfig, header *types.Header, state *state.IntraBlockState, txs []*types.Transaction,
	uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	panic("must not be called")
}
func (c *powEngine) Seal(_ consensus.Cancel, chain consensus.ChainHeaderReader, block *types.Block, results chan<- consensus.ResultWithContext, stop <-chan struct{}) error {
	panic("must not be called")
}
func (c *powEngine) SealHash(header *types.Header) common.Hash {
	panic("must not be called")
}
func (c *powEngine) CalcDifficulty(chain consensus.ChainHeaderReader, time, parentTime uint64, parentDifficulty, parentNumber *big.Int, parentHash, parentUncleHash common.Hash) *big.Int {
	panic("must not be called")
}
func (c *powEngine) APIs(chain consensus.ChainHeaderReader) []rpc.API {
	panic("must not be called")
}

func (c *powEngine) Close() error {
	panic("must not be called")
}

func (c *powEngine) Author(header *types.Header) (common.Address, error) {
	return header.Coinbase, nil
}

func (c *chainContext) GetHeader(hash common.Hash, number uint64) *types.Header {
	return rawdb.ReadHeader(c.db, hash, number)
}

func (c *chainContext) Engine() consensus.Engine {
	return &powEngine{}
}
