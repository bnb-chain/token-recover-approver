package ethsecp256k1

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestPrivKey(t *testing.T) {
	addr := common.HexToAddress("0xb26859a7321AB7B2025E5E6a425D697e2eacbFB1")
	pk, err := NewPrivKey("afc2986f283cf5f9d17e04c6a12ccf8fa46149fc37d48e11abef15a46ae34eb7")
	if err != nil {
		t.Fatal(err)
	}
	if pk.PubKey().Address() != addr {
		t.Fatalf("expected address %s, got %s", addr, pk.PubKey().Address())
	}
	message := crypto.Keccak256([]byte("hello"))
	signature, err := pk.Sign(message)
	if err != nil {
		t.Fatal(err)
	}
	valid := pk.PubKey().Verify(message, signature)
	if !valid {
		t.Fatal("expected valid signature")
	}
}
