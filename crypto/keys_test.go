package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T)  {
	privKey := GeneratePrivateKey()
	assert.Equal(t, len(privKey.Bytes()), privKeyLen);

	pubKey := privKey.Public()
	assert.Equal(t, len(pubKey.Bytes()), pubKeyLen)
}

func TestPrivateKeyFromString(t *testing.T) {
	var (
		seed = "44ba732b93d3f41c8cf5400c8eb22285d41236545b5fec0dbdebfd03bf07a82b"
		privKey = NewPrivateKeyFromString(seed)
		addressStr = "6647fa7ae144122acd2a59f5812bb12752c8a3c6"
	)

	assert.Equal(t, privKeyLen, len(privKey.Bytes()))
	address := privKey.Public().Address()
	assert.Equal(t, address.String(), addressStr)
}

func TestPrivateKeySignature(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.Public()
	msg := []byte("test signature")

	sig := privKey.Sign(msg)
	assert.True(t, sig.Verify(pubKey, msg), "Signature should verify with the public key")

	// Test with a different message
	invalidMsg := []byte("different message")
	assert.False(t, sig.Verify(pubKey, invalidMsg), "Signature should not verify with a different message")

	// Test with a different public key
	invalidPubKey := GeneratePrivateKey().Public()
	assert.False(t, sig.Verify(invalidPubKey, msg), "Signature should not verify with a different public key")
}

func TestPublicKeyToAddress(t *testing.T)  {
	privKey := GeneratePrivateKey()
	pubKey := privKey.Public()
	address := pubKey.Address()
	assert.Equal(t, len(address.value), addressLen, "Address length should match addressLen")
}
