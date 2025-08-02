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