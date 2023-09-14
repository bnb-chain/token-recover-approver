package memory

import "errors"

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrAssetNotFound   = errors.New("asset not found")
	ErrProofNotFound   = errors.New("proof not found")
)
