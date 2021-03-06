package main

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
)

const (
	AESKEY1   = "01234567890123456789012345678901"
	AESKEY2   = "01234567890123456789012345678902"
	ECDSAKEY1 = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIH4Uv66F9kZMdOQxwNegkGm8c3AB3nGPOtxNKi6wb/ZooAoGCCqGSM49
AwEHoUQDQgAEEPE+VLOh+e4NpwIjI/b/fKYHi4weU7r9OTEYPiAJiJBQY6TZnvF5
oRMvwO4MCYxFtpIRO4UxIgcZBj4NCBxKqQ==
-----END EC PRIVATE KEY-----`
	ECDSAKEY2 = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIE8Seyc9TXx+yQfnGPuzjkuEfMbkq203IYdfyvMd0r3OoAoGCCqGSM49
AwEHoUQDQgAE4dcGMMroH2LagI/s5i/Bx4t4ggGDoJPNVkKBDBlIaMYjJFYD1obk
JOWqAZxKKsBxBC5Ssu+fS26VPfdNWxDsFQ==
-----END EC PRIVATE KEY-----`
)

func TestInit(t *testing.T) {
	scc := new(EncCC)
	stub := shim.NewMockStub("enccc", scc)
	stub.MockTransactionStart("a")
	res := scc.Init(stub)
	stub.MockTransactionEnd("a")
	assert.Equal(t, res.Status, int32(shim.OK))
}

// unfortunately we can't tese this cc though invoke since the
// mock shim doesn't support transient. We test failure scenarios
// and the tests below focus on the functionality by invoking
// functions as opposed to cc
func TestInvoke(t *testing.T) {
	scc := &EncCC{factory.GetDefault()}
	stub := shim.NewMockStub("enccc", scc)

	res := stub.MockInvoke("tx", [][]byte{[]byte("barf")})
	assert.NotEqual(t, res.Status, int32(shim.OK))
	res = stub.MockInvoke("tx", [][]byte{[]byte("ENC")})
	assert.NotEqual(t, res.Status, int32(shim.OK))
	res = stub.MockInvoke("tx", [][]byte{[]byte("SIG")})
	assert.NotEqual(t, res.Status, int32(shim.OK))
}

func TestEnc(t *testing.T) {
	factory.InitFactories(nil)

	scc := &EncCC{factory.GetDefault()}
	stub := shim.NewMockStub("enccc", scc)

	// success
	stub.MockTransactionStart("a")
	res := scc.Encrypter(stub, "PUT", []string{"key", "value"}, []byte(AESKEY1))
	stub.MockTransactionEnd("a")
	assert.Equal(t, res.Status, int32(shim.OK))

	// fail - bad key
	stub.MockTransactionStart("a")
	res = scc.Encrypter(stub, "PUT", []string{"key", "value"}, []byte("badkey"))
	stub.MockTransactionEnd("a")
	assert.NotEqual(t, res.Status, int32(shim.OK))

	// fail - not enough args
	stub.MockTransactionStart("a")
	res = scc.Encrypter(stub, "PUT", []string{"key"}, []byte(AESKEY1))
	stub.MockTransactionEnd("a")
	assert.NotEqual(t, res.Status, int32(shim.OK))

	// success
	stub.MockTransactionStart("a")
	res = scc.Encrypter(stub, "GET", []string{"key"}, []byte(AESKEY1))
	stub.MockTransactionEnd("a")
	assert.Equal(t, res.Status, int32(shim.OK))
	assert.True(t, bytes.Equal(res.Payload, []byte("value")))

	// fail - not enough args
	stub.MockTransactionStart("a")
	res = scc.Encrypter(stub, "GET", []string{}, []byte(AESKEY1))
	stub.MockTransactionEnd("a")
	assert.NotEqual(t, res.Status, int32(shim.OK))

	// fail - wrong function name
	stub.MockTransactionStart("a")
	res = scc.Encrypter(stub, "barf", []string{}, []byte(AESKEY1))
	stub.MockTransactionEnd("a")
	assert.NotEqual(t, res.Status, int32(shim.OK))

	// fail - bad kvs key
	stub.MockTransactionStart("a")
	res = scc.Encrypter(stub, "GET", []string{"badkey"}, []byte(AESKEY1))
	stub.MockTransactionEnd("a")
	assert.NotEqual(t, res.Status, int32(shim.OK))

	// fail - bad key
	stub.MockTransactionStart("a")
	res = scc.Encrypter(stub, "GET", []string{"key"}, []byte(AESKEY2))
	stub.MockTransactionEnd("a")
	assert.NotEqual(t, res.Status, int32(shim.OK))
}

func TestSig(t *testing.T) {
	factory.InitFactories(nil)

	scc := &EncCC{factory.GetDefault()}
	stub := shim.NewMockStub("enccc", scc)

	// success
	stub.MockTransactionStart("a")
	res := scc.EncrypterSigner(stub, "PUT", []string{"key", "value"}, []byte(AESKEY1), []byte(ECDSAKEY1))
	stub.MockTransactionEnd("a")
	assert.Equal(t, res.Status, int32(shim.OK))

	// fail - bad key
	stub.MockTransactionStart("a")
	res = scc.EncrypterSigner(stub, "PUT", []string{"key", "value"}, []byte(AESKEY1), []byte("barf"))
	stub.MockTransactionEnd("a")
	assert.NotEqual(t, res.Status, int32(shim.OK))

	// fail - bad args
	stub.MockTransactionStart("a")
	res = scc.EncrypterSigner(stub, "PUT", []string{"key"}, []byte(AESKEY1), []byte("barf"))
	stub.MockTransactionEnd("a")
	assert.NotEqual(t, res.Status, int32(shim.OK))

	// fail - bad signing key
	stub.MockTransactionStart("a")
	res = scc.EncrypterSigner(stub, "GET", []string{"key"}, []byte(AESKEY1), []byte(ECDSAKEY2))
	stub.MockTransactionEnd("a")
	assert.NotEqual(t, res.Status, int32(shim.OK))

	// fail - bad args
	stub.MockTransactionStart("a")
	res = scc.EncrypterSigner(stub, "GET", []string{}, []byte(AESKEY1), []byte(ECDSAKEY1))
	stub.MockTransactionEnd("a")
	assert.NotEqual(t, res.Status, int32(shim.OK))

	// fail - bad kvs key
	stub.MockTransactionStart("a")
	res = scc.EncrypterSigner(stub, "GET", []string{"badkey"}, []byte(AESKEY1), []byte(ECDSAKEY1))
	stub.MockTransactionEnd("a")
	assert.NotEqual(t, res.Status, int32(shim.OK))

	// fail - wrong function name
	stub.MockTransactionStart("a")
	res = scc.EncrypterSigner(stub, "barf", []string{"badkey"}, []byte(AESKEY1), []byte(ECDSAKEY1))
	stub.MockTransactionEnd("a")
	assert.NotEqual(t, res.Status, int32(shim.OK))

	// success
	stub.MockTransactionStart("a")
	res = scc.EncrypterSigner(stub, "PUT", []string{"key", "value"}, []byte(AESKEY1), []byte(ECDSAKEY1))
	stub.MockTransactionEnd("a")
	assert.Equal(t, res.Status, int32(shim.OK))

	// success
	stub.MockTransactionStart("a")
	res = scc.EncrypterSigner(stub, "GET", []string{"key"}, []byte(AESKEY1), []byte(ECDSAKEY1))
	stub.MockTransactionEnd("a")
	assert.Equal(t, res.Status, int32(shim.OK))
	assert.True(t, bytes.Equal(res.Payload, []byte("value")))
}

func TestEncCC_RangeDecrypter(t *testing.T) {
	factory.InitFactories(nil)

	scc := &EncCC{factory.GetDefault()}
	stub := shim.NewMockStub("enccc", scc)

	stub.MockTransactionStart("a")
	res := scc.Encrypter(stub, "PUT", []string{"key1", "value1"}, []byte(AESKEY1))
	stub.MockTransactionEnd("a")
	assert.Equal(t, res.Status, int32(shim.OK))

	stub.MockTransactionStart("a")
	res = scc.Encrypter(stub, "PUT", []string{"key2", "value2"}, []byte(AESKEY1))
	stub.MockTransactionEnd("a")
	assert.Equal(t, res.Status, int32(shim.OK))

	stub.MockTransactionStart("a")
	res = scc.Encrypter(stub, "PUT", []string{"key3", "value3"}, []byte(AESKEY1))
	stub.MockTransactionEnd("a")
	assert.Equal(t, res.Status, int32(shim.OK))

	// success range query
	stub.MockTransactionStart("a")
	res = scc.RangeDecrypter(stub, []byte(AESKEY1))
	stub.MockTransactionEnd("a")
	assert.Equal(t, res.Status, int32(shim.OK))
	keys := []keyValuePair{}
	err := json.Unmarshal(res.Payload, &keys)
	assert.NoError(t, err)
	assert.Equal(t, keys[0].Key, "key1")
	assert.Equal(t, string(keys[0].Value), "value1")
	assert.Equal(t, keys[1].Key, "key2")
	assert.Equal(t, string(keys[1].Value), "value2")
	assert.Equal(t, keys[2].Key, "key3")
	assert.Equal(t, string(keys[2].Value), "value3")
}
