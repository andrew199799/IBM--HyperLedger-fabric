/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privdata

import (
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"testing"

	util2 "github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/core/ledger"
	"github.com/hyperledger/fabric/core/transientstore"
	"github.com/hyperledger/fabric/gossip/util"
	"github.com/hyperledger/fabric/protos/common"
	proto "github.com/hyperledger/fabric/protos/gossip"
	"github.com/hyperledger/fabric/protos/ledger/rwset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type persistCall struct {
	*mock.Call
	store *mockTransientStore
}

func (pc *persistCall) expectRWSet(namespace string, collection string, rws []byte) *persistCall {
	if pc.store.persists == nil {
		pc.store.persists = make(map[rwsTriplet]struct{})
	}
	pc.store.persists[rwsTriplet{
		namespace:  namespace,
		collection: collection,
		rwset:      hex.EncodeToString(rws),
	}] = struct{}{}
	return pc
}

type mockTransientStore struct {
	t *testing.T
	mock.Mock
	persists      map[rwsTriplet]struct{}
	lastReqTxID   string
	lastReqFilter map[string]ledger.PvtCollFilter
}

func (store *mockTransientStore) On(methodName string, arguments ...interface{}) *persistCall {
	return &persistCall{
		store: store,
		Call:  store.Mock.On(methodName, arguments...),
	}
}

func (store *mockTransientStore) Persist(txid string, endorsementBlkHt uint64, res *rwset.TxPvtReadWriteSet) error {
	key := rwsTriplet{
		namespace:  res.NsPvtRwset[0].Namespace,
		collection: res.NsPvtRwset[0].CollectionPvtRwset[0].CollectionName,
		rwset:      hex.EncodeToString(res.NsPvtRwset[0].CollectionPvtRwset[0].Rwset)}
	if _, exists := store.persists[key]; !exists {
		store.t.Fatal("Shouldn't have persisted", res)
	}
	delete(store.persists, key)
	store.Called(txid, res)
	return nil
}

func (store *mockTransientStore) GetTxPvtRWSetByTxid(txid string, filter ledger.PvtNsCollFilter) (transientstore.RWSetScanner, error) {
	store.lastReqTxID = txid
	store.lastReqFilter = filter
	args := store.Called(txid, filter)
	if args.Get(1) == nil {
		return args.Get(0).(transientstore.RWSetScanner), nil
	}
	return nil, args.Get(1).(error)
}

type mockRWSetScanner struct {
	err     error
	results []*transientstore.EndorserPvtSimulationResults
}

func (scanner *mockRWSetScanner) withRWSet(ns string, col string) *mockRWSetScanner {
	scanner.results = append(scanner.results, &transientstore.EndorserPvtSimulationResults{
		PvtSimulationResults: &rwset.TxPvtReadWriteSet{
			DataModel: rwset.TxReadWriteSet_KV,
			NsPvtRwset: []*rwset.NsPvtReadWriteSet{
				{
					Namespace: ns,
					CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{
						{
							CollectionName: col,
							Rwset:          []byte("rws-pre-image"),
						},
					},
				},
			},
		},
	})
	return scanner
}

func (scanner *mockRWSetScanner) Next() (*transientstore.EndorserPvtSimulationResults, error) {
	if scanner.err != nil {
		return nil, scanner.err
	}
	var res *transientstore.EndorserPvtSimulationResults
	if len(scanner.results) == 0 {
		return nil, nil
	}
	res, scanner.results = scanner.results[len(scanner.results)-1], scanner.results[:len(scanner.results)-1]
	return res, nil
}

func (*mockRWSetScanner) Close() {
}

type committerMock struct {
	mock.Mock
}

func (mock *committerMock) CommitWithPvtData(blockAndPvtData *ledger.BlockAndPvtData) error {
	args := mock.Called(blockAndPvtData)
	return args.Error(0)
}

func (mock *committerMock) GetPvtDataAndBlockByNum(seqNum uint64) (*ledger.BlockAndPvtData, error) {
	args := mock.Called(seqNum)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ledger.BlockAndPvtData), args.Error(1)
}

func (mock *committerMock) Commit(block *common.Block) error {
	args := mock.Called(block)
	return args.Error(0)
}

func (mock *committerMock) LedgerHeight() (uint64, error) {
	args := mock.Called()
	return args.Get(0).(uint64), args.Error(1)
}

func (mock *committerMock) GetBlocks(blockSeqs []uint64) []*common.Block {
	args := mock.Called(blockSeqs)
	seqs := args.Get(0)
	if seqs == nil {
		return nil
	}
	return seqs.([]*common.Block)
}

func (mock *committerMock) Close() {
	mock.Called()
}

type digests []*proto.PvtDataDigest

func (d digests) Equal(other digests) bool {
	flatten := func(d digests) map[proto.PvtDataDigest]struct{} {
		m := map[proto.PvtDataDigest]struct{}{}
		for _, dig := range d {
			m[*dig] = struct{}{}
		}
		return m
	}
	return reflect.DeepEqual(flatten(d), flatten(other))
}

type fetchCall struct {
	fetcher *fetcherMock
	*mock.Call
}

func (fc *fetchCall) expectingReq(req *proto.RemotePvtDataRequest) *fetchCall {
	fc.fetcher.expectedReq = req
	return fc
}

func (fc *fetchCall) Return(returnArguments ...interface{}) *mock.Call {

	return fc.Call.Return(returnArguments...)
}

type fetcherMock struct {
	t *testing.T
	mock.Mock
	expectedReq *proto.RemotePvtDataRequest
}

func (f *fetcherMock) On(methodName string, arguments ...interface{}) *fetchCall {
	return &fetchCall{
		fetcher: f,
		Call:    f.Mock.On(methodName, arguments...),
	}
}

func (f *fetcherMock) fetch(req *proto.RemotePvtDataRequest) ([]*proto.PvtDataElement, error) {
	assert.True(f.t, digests(req.Digests).Equal(digests(f.expectedReq.Digests)))
	args := f.Called(req)
	if args.Get(1) == nil {
		return args.Get(0).([]*proto.PvtDataElement), nil
	}
	return nil, args.Get(1).(error)
}

func TestPvtDataCollections_FailOnEmptyPayload(t *testing.T) {
	collection := &util.PvtDataCollections{
		&ledger.TxPvtData{
			SeqInBlock: uint64(1),
			WriteSet: &rwset.TxPvtReadWriteSet{
				DataModel: rwset.TxReadWriteSet_KV,
				NsPvtRwset: []*rwset.NsPvtReadWriteSet{
					{
						Namespace: "ns1",
						CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{
							{
								CollectionName: "secretCollection",
								Rwset:          []byte{1, 2, 3, 4, 5, 6, 7},
							},
						},
					},
				},
			},
		},

		nil,
	}

	_, err := collection.Marshal()
	assertion := assert.New(t)
	assertion.Error(err, "Expected to fail since second item has nil payload")
	assertion.Equal("Mallformed private data payload, rwset index 1 is nil", fmt.Sprintf("%s", err))
}

func TestPvtDataCollections_FailMarshalingWriteSet(t *testing.T) {
	collection := &util.PvtDataCollections{
		&ledger.TxPvtData{
			SeqInBlock: uint64(1),
			WriteSet:   nil,
		},
	}

	_, err := collection.Marshal()
	assertion := assert.New(t)
	assertion.Error(err, "Expected to fail since first item has nil writeset")
	assertion.Contains(fmt.Sprintf("%s", err), "Could not marshal private rwset index 0")
}

func TestPvtDataCollections_Marshal(t *testing.T) {
	collection := &util.PvtDataCollections{
		&ledger.TxPvtData{
			SeqInBlock: uint64(1),
			WriteSet: &rwset.TxPvtReadWriteSet{
				DataModel: rwset.TxReadWriteSet_KV,
				NsPvtRwset: []*rwset.NsPvtReadWriteSet{
					{
						Namespace: "ns1",
						CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{
							{
								CollectionName: "secretCollection",
								Rwset:          []byte{1, 2, 3, 4, 5, 6, 7},
							},
						},
					},
				},
			},
		},

		&ledger.TxPvtData{
			SeqInBlock: uint64(2),
			WriteSet: &rwset.TxPvtReadWriteSet{
				DataModel: rwset.TxReadWriteSet_KV,
				NsPvtRwset: []*rwset.NsPvtReadWriteSet{
					{
						Namespace: "ns1",
						CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{
							{
								CollectionName: "secretCollection",
								Rwset:          []byte{42, 42, 42, 42, 42, 42, 42},
							},
						},
					},
					{
						Namespace: "ns2",
						CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{
							{
								CollectionName: "otherCollection",
								Rwset:          []byte{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
							},
						},
					},
				},
			},
		},
	}

	bytes, err := collection.Marshal()

	assertion := assert.New(t)
	assertion.NoError(err)
	assertion.NotNil(bytes)
	assertion.Equal(2, len(bytes))
}

func TestPvtDataCollections_Unmarshal(t *testing.T) {
	collection := util.PvtDataCollections{
		&ledger.TxPvtData{
			SeqInBlock: uint64(1),
			WriteSet: &rwset.TxPvtReadWriteSet{
				DataModel: rwset.TxReadWriteSet_KV,
				NsPvtRwset: []*rwset.NsPvtReadWriteSet{
					{
						Namespace: "ns1",
						CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{
							{
								CollectionName: "secretCollection",
								Rwset:          []byte{1, 2, 3, 4, 5, 6, 7},
							},
						},
					},
				},
			},
		},
	}

	bytes, err := collection.Marshal()

	assertion := assert.New(t)
	assertion.NoError(err)
	assertion.NotNil(bytes)
	assertion.Equal(1, len(bytes))

	var newCol util.PvtDataCollections

	err = newCol.Unmarshal(bytes)
	assertion.NoError(err)
	assertion.Equal(newCol, collection)
}

type rwsTriplet struct {
	namespace  string
	collection string
	rwset      string
}

type privateData map[uint64]*ledger.TxPvtData

func (pd privateData) Equal(that privateData) bool {
	return reflect.DeepEqual(pd.flatten(), that.flatten())
}

func (pd privateData) flatten() map[uint64]map[rwsTriplet]struct{} {
	m := make(map[uint64]map[rwsTriplet]struct{})
	for seqInBlock, namespaces := range pd {
		triplets := make(map[rwsTriplet]struct{})
		for _, namespace := range namespaces.WriteSet.NsPvtRwset {
			for _, col := range namespace.CollectionPvtRwset {
				triplets[rwsTriplet{
					namespace:  namespace.Namespace,
					collection: col.CollectionName,
					rwset:      hex.EncodeToString(col.Rwset),
				}] = struct{}{}
			}
		}
		m[seqInBlock] = triplets
	}
	return m
}

var expectedCommittedPrivateData1 = map[uint64]*ledger.TxPvtData{
	0: {SeqInBlock: 0, WriteSet: &rwset.TxPvtReadWriteSet{
		DataModel: rwset.TxReadWriteSet_KV,
		NsPvtRwset: []*rwset.NsPvtReadWriteSet{
			{
				Namespace: "ns1",
				CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{
					{
						CollectionName: "c1",
						Rwset:          []byte("rws-pre-image"),
					},
					{
						CollectionName: "c2",
						Rwset:          []byte("rws-pre-image"),
					},
				},
			},
		},
	}},
	1: {SeqInBlock: 1, WriteSet: &rwset.TxPvtReadWriteSet{
		DataModel: rwset.TxReadWriteSet_KV,
		NsPvtRwset: []*rwset.NsPvtReadWriteSet{
			{
				Namespace: "ns2",
				CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{
					{
						CollectionName: "c1",
						Rwset:          []byte("rws-pre-image"),
					},
				},
			},
		},
	}},
}

var expectedCommittedPrivateData2 = map[uint64]*ledger.TxPvtData{
	0: {SeqInBlock: 0, WriteSet: &rwset.TxPvtReadWriteSet{
		DataModel: rwset.TxReadWriteSet_KV,
		NsPvtRwset: []*rwset.NsPvtReadWriteSet{
			{
				Namespace: "ns3",
				CollectionPvtRwset: []*rwset.CollectionPvtReadWriteSet{
					{
						CollectionName: "c3",
						Rwset:          []byte("rws-pre-image"),
					},
				},
			},
		},
	}},
}

func TestCoordinatorStoreBlock(t *testing.T) {
	// Green path test, all private data should be obtained successfully
	var commitHappened bool
	assertCommitHappened := func() {
		assert.True(t, commitHappened)
		commitHappened = false
	}
	committer := &committerMock{}
	committer.On("CommitWithPvtData", mock.Anything).Run(func(args mock.Arguments) {
		var privateDataPassed2Ledger privateData = args.Get(0).(*ledger.BlockAndPvtData).BlockPvtData
		assert.True(t, privateDataPassed2Ledger.Equal(expectedCommittedPrivateData1))
		commitHappened = true
	}).Return(nil)
	store := &mockTransientStore{t: t}
	fetcher := &fetcherMock{t: t}

	hash := util2.ComputeSHA256([]byte("rws-pre-image"))
	pdFactory := &pvtDataFactory{}
	bf := &blockFactory{
		channelID: "test",
	}
	block := bf.AddTxn("tx1", "ns1", hash, "c1", "c2").AddTxn("tx2", "ns2", hash, "c1").create()

	// Scenario I: Block we got has sufficient private data alongside it.
	// If the coordinator tries fetching from the transientstore, or peers it would result in panic,
	// because we didn't define yet the "On(...)" invocation of the transient store or other peers.
	pvtData := pdFactory.addRWSet().addNSRWSet("ns1", "c1", "c2").addRWSet().addNSRWSet("ns2", "c1").create()
	coordinator := NewCoordinator(committer, store, fetcher)
	err := coordinator.StoreBlock(block, pvtData)
	assert.NoError(t, err)
	assertCommitHappened()

	// Scenario II: Block we got doesn't have sufficient private data alongside it,
	// it is missing ns1: c2, but the data exists in the transient store
	store.On("GetTxPvtRWSetByTxid", "tx1", mock.Anything).Return((&mockRWSetScanner{}).withRWSet("ns1", "c2"), nil)
	store.On("GetTxPvtRWSetByTxid", "tx2", mock.Anything).Return(&mockRWSetScanner{}, nil)
	pvtData = pdFactory.addRWSet().addNSRWSet("ns1", "c1").addRWSet().addNSRWSet("ns2", "c1").create()
	err = coordinator.StoreBlock(block, pvtData)
	assert.NoError(t, err)
	assertCommitHappened()
	assert.Equal(t, "tx1", store.lastReqTxID)
	assert.Equal(t, map[string]ledger.PvtCollFilter{
		"ns1": map[string]bool{
			"c2": true,
		},
	}, store.lastReqFilter)

	// Scenario III: Block doesn't have sufficient private data alongside it,
	// it is missing ns1: c2, and the data exists in the transient store,
	// but it is also missing ns2: c1, and that data doesn't exist in the transient store - but in a peer
	fetcher.On("fetch", mock.Anything).expectingReq(&proto.RemotePvtDataRequest{
		Digests: []*proto.PvtDataDigest{
			{
				TxId: "tx1", Namespace: "ns1", Collection: "c2", BlockSeq: 1,
			},
			{
				TxId: "tx2", Namespace: "ns2", Collection: "c1", BlockSeq: 1, SeqInBlock: 1,
			},
		},
	}).Return([]*proto.PvtDataElement{
		{
			Digest: &proto.PvtDataDigest{
				BlockSeq:   1,
				Collection: "c2",
				Namespace:  "ns1",
				TxId:       "tx1",
			},
			Payload: [][]byte{[]byte("rws-pre-image")},
		},
		{
			Digest: &proto.PvtDataDigest{
				SeqInBlock: 1,
				BlockSeq:   1,
				Collection: "c1",
				Namespace:  "ns2",
				TxId:       "tx2",
			},
			Payload: [][]byte{[]byte("rws-pre-image")},
		},
	}, nil)
	store.On("Persist", mock.Anything, mock.Anything).
		expectRWSet("ns1", "c2", []byte("rws-pre-image")).
		expectRWSet("ns2", "c1", []byte("rws-pre-image")).Return(nil)
	pvtData = pdFactory.addRWSet().addNSRWSet("ns1", "c1").create()
	err = coordinator.StoreBlock(block, pvtData)
	assert.NoError(t, err)
	assertCommitHappened()

	// Scenario IV: Block came with more than sufficient private data alongside it, some of it is redundant.
	pvtData = pdFactory.addRWSet().addNSRWSet("ns1", "c1", "c2", "c3").
		addRWSet().addNSRWSet("ns2", "c1", "c3").addRWSet().addNSRWSet("ns1", "c4").create()
	err = coordinator.StoreBlock(block, pvtData)
	assert.NoError(t, err)
	assertCommitHappened()

	// Scenario V: Block didn't get with any private data alongside it, and the transient store
	// has some problem.
	// In this case, we should try to fetch data from peers.
	block = bf.AddTxn("tx3", "ns3", hash, "c3").create()
	fetcher = &fetcherMock{t: t}
	fetcher.On("fetch", mock.Anything).expectingReq(&proto.RemotePvtDataRequest{
		Digests: []*proto.PvtDataDigest{
			{
				TxId: "tx3", Namespace: "ns3", Collection: "c3", BlockSeq: 1,
			},
		},
	}).Return([]*proto.PvtDataElement{
		{
			Digest: &proto.PvtDataDigest{
				BlockSeq:   1,
				Collection: "c3",
				Namespace:  "ns3",
				TxId:       "tx3",
			},
			Payload: [][]byte{[]byte("rws-pre-image")},
		},
	}, nil)
	store = &mockTransientStore{t: t}
	store.On("GetTxPvtRWSetByTxid", "tx3", mock.Anything).Return(&mockRWSetScanner{err: errors.New("uh oh")}, nil)
	store.On("Persist", mock.Anything, mock.Anything).expectRWSet("ns3", "c3", []byte("rws-pre-image"))
	committer = &committerMock{}
	committer.On("CommitWithPvtData", mock.Anything).Run(func(args mock.Arguments) {
		var privateDataPassed2Ledger privateData = args.Get(0).(*ledger.BlockAndPvtData).BlockPvtData
		assert.True(t, privateDataPassed2Ledger.Equal(expectedCommittedPrivateData2))
		commitHappened = true
	}).Return(nil)
	coordinator = NewCoordinator(committer, store, fetcher)
	err = coordinator.StoreBlock(block, nil)
	assert.NoError(t, err)
	assertCommitHappened()
}

func TestCoordinatorGetBlocks(t *testing.T) {
	committer := &committerMock{}
	store := &mockTransientStore{t: t}
	fetcher := &fetcherMock{t: t}
	coordinator := NewCoordinator(committer, store, fetcher)

	// Bad path: block is not returned
	committer.On("GetBlocks", mock.Anything).Return([]*common.Block{})
	block, err := coordinator.GetBlockByNum(1)
	assert.Contains(t, err.Error(), "cannot retrieve block number 1")
	assert.Error(t, err)
	assert.Nil(t, block)

	// Green path: Block is returned
	committer.Mock = mock.Mock{}
	committer.On("GetBlocks", mock.Anything).Return([]*common.Block{
		{
			Header: &common.BlockHeader{
				Number: 1,
			},
		},
	})
	block, err = coordinator.GetBlockByNum(1)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), block.Header.Number)

	committer.On("GetPvtDataAndBlockByNum", mock.Anything).Return(&ledger.BlockAndPvtData{
		Block: block,
		BlockPvtData: map[uint64]*ledger.TxPvtData{
			0: {
				SeqInBlock: 0,
			},
		},
	}, nil)

	// Green path - block and private data is returned
	block2, privData, err := coordinator.GetPvtDataAndBlockByNum(1)
	assert.NoError(t, err)
	assert.Equal(t, block, block2)
	assert.Equal(t, util.PvtDataCollections{{SeqInBlock: 0}}, privData)

	// Bad path - error occurs when trying to retrieve the block and private data
	committer.Mock = mock.Mock{}
	committer.On("GetPvtDataAndBlockByNum", mock.Anything).Return(nil, errors.New("uh oh"))
	block2, privData, err = coordinator.GetPvtDataAndBlockByNum(1)
	assert.Nil(t, block2)
	assert.Empty(t, privData)
	assert.Error(t, err)
}
