package ethsecp256k1

import (
	"encoding/hex"
	"errors"
)

func RecoveryFromPrivateKey(privateKey string) ([32]byte, error) {
	var keyBytesArray [32]byte

	priBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return keyBytesArray, err
	}

	if len(priBytes) != 32 {
		return keyBytesArray, errors.New("len of keybytes is not equal to 32")
	}

	copy(keyBytesArray[:], priBytes[:32])
	return keyBytesArray, nil
}
