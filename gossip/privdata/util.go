/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privdata

import (
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/ledger"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/ledger/rwset"
	"github.com/hyperledger/fabric/protos/ledger/rwset/kvrwset"
	"github.com/hyperledger/fabric/protos/peer"
)

type blockFactory struct {
	channelID    string
	transactions [][]byte
}

func (bf *blockFactory) AddTxn(txID string, nsName string, hash []byte, collections ...string) *blockFactory {
	txn := &peer.Transaction{
		Actions: []*peer.TransactionAction{
			{},
		},
	}
	txrws := rwsetutil.TxRwSet{
		NsRwSets: []*rwsetutil.NsRwSet{sampleNsRwSet(nsName, hash, collections...)},
	}

	b, err := txrws.ToProtoBytes()
	if err != nil {
		panic(err)
	}
	ccAction := &peer.ChaincodeAction{
		Results: b,
	}

	ccActionBytes, err := proto.Marshal(ccAction)
	if err != nil {
		panic(err)
	}
	pRespPayload := &peer.ProposalResponsePayload{
		Extension: ccActionBytes,
	}

	respPayloadBytes, err := proto.Marshal(pRespPayload)
	if err != nil {
		panic(err)
	}

	ccPayload := &peer.ChaincodeActionPayload{
		Action: &peer.ChaincodeEndorsedAction{
			ProposalResponsePayload: respPayloadBytes,
		},
	}

	ccPayloadBytes, err := proto.Marshal(ccPayload)
	if err != nil {
		panic(err)
	}

	txn.Actions[0].Payload = ccPayloadBytes
	txBytes, _ := proto.Marshal(txn)

	cHdr := &common.ChannelHeader{
		TxId:      txID,
		Type:      int32(common.HeaderType_ENDORSER_TRANSACTION),
		ChannelId: bf.channelID,
	}
	cHdrBytes, _ := proto.Marshal(cHdr)
	commonPayload := &common.Payload{
		Header: &common.Header{
			ChannelHeader: cHdrBytes,
		},
		Data: txBytes,
	}

	payloadBytes, _ := proto.Marshal(commonPayload)
	envp := &common.Envelope{
		Payload: payloadBytes,
	}
	envelopeBytes, _ := proto.Marshal(envp)

	bf.transactions = append(bf.transactions, envelopeBytes)
	return bf
}

func (bf *blockFactory) create() *common.Block {
	defer func() {
		bf.transactions = nil
	}()
	return &common.Block{
		Header: &common.BlockHeader{
			Number: 1,
		},
		Data: &common.BlockData{
			Data: bf.transactions,
		},
	}
}

func sampleNsRwSet(ns string, hash []byte, collections ...string) *rwsetutil.NsRwSet {
	nsRwSet := &rwsetutil.NsRwSet{NameSpace: ns,
		KvRwSet: sampleKvRwSet(),
	}
	for _, col := range collections {
		nsRwSet.CollHashedRwSets = append(nsRwSet.CollHashedRwSets, sampleCollHashedRwSet(col, hash))
	}
	return nsRwSet
}

func sampleKvRwSet() *kvrwset.KVRWSet {
	rqi1 := &kvrwset.RangeQueryInfo{StartKey: "k0", EndKey: "k9", ItrExhausted: true}
	rqi1.SetRawReads([]*kvrwset.KVRead{
		{Key: "k1", Version: &kvrwset.Version{BlockNum: 1, TxNum: 1}},
		{Key: "k2", Version: &kvrwset.Version{BlockNum: 1, TxNum: 2}},
	})

	rqi2 := &kvrwset.RangeQueryInfo{StartKey: "k00", EndKey: "k90", ItrExhausted: true}
	rqi2.SetMerkelSummary(&kvrwset.QueryReadsMerkleSummary{MaxDegree: 5, MaxLevel: 4, MaxLevelHashes: [][]byte{[]byte("Hash-1"), []byte("Hash-2")}})
	return &kvrwset.KVRWSet{
		Reads:            []*kvrwset.KVRead{{Key: "key1", Version: &kvrwset.Version{BlockNum: 1, TxNum: 1}}},
		RangeQueriesInfo: []*kvrwset.RangeQueryInfo{rqi1},
		Writes:           []*kvrwset.KVWrite{{Key: "key2", IsDelete: false, Value: []byte("value2")}},
	}
}

func sampleCollHashedRwSet(collectionName string, hash []byte) *rwsetutil.CollHashedRwSet {
	collHashedRwSet := &rwsetutil.CollHashedRwSet{
		CollectionName: collectionName,
		HashedRwSet: &kvrwset.HashedRWSet{
			HashedReads: []*kvrwset.KVReadHash{
				{KeyHash: []byte("Key-1-hash"), Version: &kvrwset.Version{1, 2}},
				{KeyHash: []byte("Key-2-hash"), Version: &kvrwset.Version{2, 3}},
			},
			HashedWrites: []*kvrwset.KVWriteHash{
				{KeyHash: []byte("Key-3-hash"), ValueHash: []byte("value-3-hash"), IsDelete: false},
				{KeyHash: []byte("Key-4-hash"), ValueHash: []byte("value-4-hash"), IsDelete: true},
			},
		},
		PvtRwSetHash: hash,
	}
	return collHashedRwSet
}

type pvtDataFactory struct {
	data []*ledger.TxPvtData
}

func (df *pvtDataFactory) addRWSet() *pvtDataFactory {
	seqInBlock := uint64(len(df.data))
	df.data = append(df.data, &ledger.TxPvtData{
		SeqInBlock: seqInBlock,
		WriteSet:   &rwset.TxPvtReadWriteSet{},
	})
	return df
}

func (df *pvtDataFactory) addNSRWSet(namespace string, collections ...string) *pvtDataFactory {
	nsrws := &rwset.NsPvtReadWriteSet{
		Namespace: namespace,
	}
	for _, col := range collections {
		nsrws.CollectionPvtRwset = append(nsrws.CollectionPvtRwset, &rwset.CollectionPvtReadWriteSet{
			CollectionName: col,
			Rwset:          []byte("rws-pre-image"),
		})
	}
	df.data[len(df.data)-1].WriteSet.NsPvtRwset = append(df.data[len(df.data)-1].WriteSet.NsPvtRwset, nsrws)
	return df
}

func (df *pvtDataFactory) create() []*ledger.TxPvtData {
	defer func() {
		df.data = nil
	}()
	return df.data
}
