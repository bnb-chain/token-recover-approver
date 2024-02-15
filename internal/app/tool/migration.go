package tool

import (
	"errors"

	"github.com/bnb-chain/token-recover-approver/internal/store"
	gormStore "github.com/bnb-chain/token-recover-approver/internal/store/gorm"
	"github.com/bnb-chain/token-recover-approver/internal/store/memory"
	"github.com/bnb-chain/token-recover-approver/pkg/util"
)

var (
	ErrInvalidStore = errors.New("invalid store type")
)

func (tool *Tool) MigrateDataFromLocalToSQL(
	proofsPath string,
) error {
	errChan := make(chan error, 1)
	defer close(errChan)
	// load proofs
	stream := util.NewJSONStream(func() any {
		return &memory.Proof{}
	})

	sqlStore, ok := tool.store.(*gormStore.SQLStore)
	if !ok {
		return ErrInvalidStore
	}

	go func() {
		for data := range stream.Watch() {
			if data.Error != nil {
				errChan <- data.Error
				return
			}
			proof := data.Data.(*memory.Proof)

			dbProof := &store.Proof{
				Address: proof.Address,
				Denom:   proof.Coin.Denom,
				Amount:  proof.Coin.Amount,
				Proof:   proof.Proof,
			}

			err := sqlStore.InsertAccountAssetProof(dbProof)
			if err != nil {
				errChan <- err
				return
			}
		}
		errChan <- nil
	}()
	stream.Start(proofsPath)
	err := <-errChan
	if err != nil {
		return err
	}

	return nil
}
