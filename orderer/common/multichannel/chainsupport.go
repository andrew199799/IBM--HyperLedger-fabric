/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package multichannel

import (
	"github.com/hyperledger/fabric/common/crypto"
	"github.com/hyperledger/fabric/orderer/common/blockcutter"
	"github.com/hyperledger/fabric/orderer/common/ledger"
	"github.com/hyperledger/fabric/orderer/common/msgprocessor"
	"github.com/hyperledger/fabric/orderer/consensus"
	cb "github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/utils"
)

// ChainSupport holds the resources for a particular channel.
type ChainSupport struct {
	*ledgerResources
	msgprocessor.Processor
	*BlockWriter
	consensus.Chain
	cutter blockcutter.Receiver
	crypto.LocalSigner
}

func newChainSupport(
	registrar *Registrar,
	ledgerResources *ledgerResources,
	consenters map[string]consensus.Consenter,
	signer crypto.LocalSigner,
) *ChainSupport {
	// Read in the last block and metadata for the channel
	lastBlock := ledger.GetBlock(ledgerResources, ledgerResources.Height()-1)

	metadata, err := utils.GetMetadataFromBlock(lastBlock, cb.BlockMetadataIndex_ORDERER)
	// Assuming a block created with cb.NewBlock(), this should not
	// error even if the orderer metadata is an empty byte slice
	if err != nil {
		logger.Fatalf("[channel: %s] Error extracting orderer metadata: %s", ledgerResources.ConfigtxManager().ChainID(), err)
	}

	// Construct limited support needed as a parameter for additional support
	cs := &ChainSupport{
		ledgerResources: ledgerResources,
		LocalSigner:     signer,
		cutter:          blockcutter.NewReceiverImpl(ledgerResources.SharedConfig()),
	}

	// Set up the msgprocessor
	cs.Processor = msgprocessor.NewStandardChannel(cs, msgprocessor.CreateStandardChannelFilters(cs))

	// Set up the block writer
	cs.BlockWriter = newBlockWriter(lastBlock, registrar, cs)

	// Set up the consenter
	consenterType := ledgerResources.SharedConfig().ConsensusType()
	consenter, ok := consenters[consenterType]
	if !ok {
		logger.Panicf("Error retrieving consenter of type: %s", consenterType)
	}

	cs.Chain, err = consenter.HandleChain(cs, metadata)
	if err != nil {
		logger.Panicf("[channel: %s] Error creating consenter: %s", cs.ChainID(), err)
	}

	logger.Debugf("[channel: %s] Done creating channel support resources", cs.ChainID())

	return cs
}

func (cs *ChainSupport) Reader() ledger.Reader {
	return cs
}

// Signer returns the crypto.Localsigner for this channel.
func (cs *ChainSupport) Signer() crypto.LocalSigner {
	return cs
}

func (cs *ChainSupport) start() {
	cs.Chain.Start()
}

// BlockCutter returns the blockcutter.Receiver instance for this channel.
func (cs *ChainSupport) BlockCutter() blockcutter.Receiver {
	return cs.cutter
}

// Validate passes through to the underlying configtxapi.Manager
func (cs *ChainSupport) Validate(configEnv *cb.ConfigEnvelope) error {
	return cs.ConfigtxManager().Validate(configEnv)
}

// ProposeConfigUpdate passes through to the underlying configtxapi.Manager
func (cs *ChainSupport) ProposeConfigUpdate(configtx *cb.Envelope) (*cb.ConfigEnvelope, error) {
	env, err := cs.ConfigtxManager().ProposeConfigUpdate(configtx)
	if err != nil {
		return nil, err
	}
	bundle, err := cs.CreateBundle(cs.ChainID(), env.Config)
	if err != nil {
		return nil, err
	}
	return env, cs.ValidateNew(bundle)
}

// ChainID passes through to the underlying configtxapi.Manager
func (cs *ChainSupport) ChainID() string {
	return cs.ConfigtxManager().ChainID()
}

// ConfigEnvelope passes through to the underlying configtxapi.Manager
func (cs *ChainSupport) ConfigEnvelope() *cb.ConfigEnvelope {
	return cs.ConfigtxManager().ConfigEnvelope()
}

// Sequence passes through to the underlying configtxapi.Manager
func (cs *ChainSupport) Sequence() uint64 {
	return cs.ConfigtxManager().Sequence()
}
