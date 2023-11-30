package util

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func TestVerifyMerkleProof(t *testing.T) {
	root := MustDecodeHexToBytes("0xad78b6dbdb34cb9a5c0e44fbfcc3bd52d6e9d519eefbfc39ba5c3b232849a064")
	proof := MustDecodeHexArrayToBytes([]string{"0x3262127e4ff0bce1bb67e569baa034637806b4519b19d3ad9dbae7f5ad31fa18"})
	leaf := MustDecodeHexToBytes("0x3c02e13b6c59d7eff8bb7b4f23a8d3bf24d60880f77c3128b0f8503606d54628")

	if !VerifyMerkleProof(root, proof, leaf) {
		t.Error("VerifyMerkleProof failed")
		t.Logf("root: %s", hexutil.Encode(root))
		hash := leaf
		for _, proofElement := range proof {
			hash = hashPair(hash, proofElement)
		}
		t.Logf("hash: %s", hexutil.Encode(hash))
	} else {
		t.Log("VerifyMerkleProof success")
	}
}
