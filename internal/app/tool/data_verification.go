package tool

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/bnb-chain/node/app"
	"github.com/bnb-chain/node/app/config"
	nodetypes "github.com/bnb-chain/node/common/types"
	"github.com/bnb-chain/token-recover-approver/internal/store"
	"github.com/bnb-chain/token-recover-approver/pkg/util"
	"github.com/spf13/viper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog"

	abci "github.com/tendermint/tendermint/abci/types"
	tmCrypto "github.com/tendermint/tendermint/crypto"
	dbm "github.com/tendermint/tendermint/libs/db"
)

const (
	displayProcessInterval = time.Second
)

var (
	ErrEmptyState = errors.New("empty state")
)

func (tool *Tool) VerifyDataFromFullnode(nodeCtx *config.BNBBeaconChainContext, home string, verifyMerkleRoot bool) error {
	emptyState, err := isEmptyState(home)
	if err != nil {
		return err
	}
	if emptyState {
		return ErrEmptyState
	}
	db, err := openDB(home)
	if err != nil {
		return err
	}
	viper.Set("home", home)
	ctx := nodeCtx.ToCosmosServerCtx()
	dapp := app.NewBNBBeaconChain(ctx.Logger, db, io.Discard)
	appCtx := dapp.NewContext(sdk.RunTxModeCheck, abci.Header{})

	totalInStore, err := tool.store.CountAccountAssetProofs()
	if err != nil {
		return err
	}
	escrowAccs := escrowAccs(tool.logger)
	count := int64(0)
	merkleRoot := util.MustDecodeHexToBytes(tool.config.MerkleRoot)
	ticker := time.NewTicker(displayProcessInterval)
	defer ticker.Stop()
	dapp.AccountKeeper.IterateAccounts(appCtx, func(acc sdk.Account) (stop bool) {
		select {
		case <-ticker.C:
			tool.logger.Info().
				Str("process", fmt.Sprintf("%d", count*100/totalInStore)+"%").
				Int64("total", totalInStore).
				Int64("count", count).Msg("verifying accounts")
		default:
		}

		namedAcc := acc.(nodetypes.NamedAccount)
		addr := namedAcc.GetAddress()
		if _, matched := escrowAccs[addr.String()]; matched {
			tool.logger.Info().Msg("skip escrow account: " + addr.String())
			return false
		}

		coins := namedAcc.GetCoins()
		frozenCoins := namedAcc.GetFrozenCoins()
		lockedCoins := namedAcc.GetLockedCoins()

		allCoins := coins.Plus(frozenCoins)
		allCoins = allCoins.Plus(lockedCoins)

		for _, coin := range allCoins {
			if coin.Amount > 0 {
				proof, err := tool.store.GetAccountAssetProof(addr, coin.Denom)
				if err != nil {
					tool.logger.Error().Str("address", addr.String()).Str("symbol", coin.Denom).Msg("proof not found")
					return true
				}

				if coin.Amount != proof.Amount {
					tool.logger.Error().
						Str("address", addr.String()).
						Str("symbol", coin.Denom).
						Int64("expected", coin.Amount).
						Int64("actual", proof.Amount).
						Msg("amount mismatch")
					return true
				}

				if verifyMerkleRoot {
					// verify merkle proof
					leaf := store.Proof{Address: addr, Denom: coin.Denom, Amount: coin.Amount}
					leafHash, err := leaf.Serialize()
					if err != nil {
						tool.logger.Error().
							Str("address", addr.String()).
							Str("symbol", coin.Denom).
							Int64("amount", coin.Amount).
							Msg("merkle proof serialization failed")
						return true
					}
					if !util.VerifyMerkleProof(merkleRoot, proof.Proof, leafHash) {
						tool.logger.Error().
							Str("address", addr.String()).
							Str("symbol", coin.Denom).
							Int64("amount", coin.Amount).
							Msg("merkle proof verification failed")
						return true
					}
				}
				count++
			}
		}

		return false
	})

	if count != totalInStore {
		return fmt.Errorf("account mismatch: %d != %d", count, totalInStore)
	}

	return nil
}

// Escrow Accounts
func escrowAccs(logger *zerolog.Logger) map[string]struct{} {
	escrowAccs := make(map[string]struct{})
	// bnb prefix address: bnb1vu5max8wqn997ayhrrys0drpll2rlz4dh39s3h
	// tbnb prefix address: tbnb1vu5max8wqn997ayhrrys0drpll2rlz4deyv53x
	depositedCoinsAccAddr := sdk.AccAddress(tmCrypto.AddressHash([]byte("BinanceChainDepositedCoins")))
	// bnb prefix address: bnb1j725qk29cv4kwpers4addy9x93ukhw7czfkjaj
	// tbnb prefix address: tbnb1j725qk29cv4kwpers4addy9x93ukhw7cvulkar
	delegationAccAddr := sdk.AccAddress(tmCrypto.AddressHash([]byte("BinanceChainStakeDelegation")))
	// bnb prefix address: bnb1v8vkkymvhe2sf7gd2092ujc6hweta38xadu2pj
	// tbnb prefix address: tbnb1v8vkkymvhe2sf7gd2092ujc6hweta38xnc4wpr
	pegAccount := sdk.AccAddress(tmCrypto.AddressHash([]byte("BinanceChainPegAccount")))
	// bnb prefix address: bnb1wxeplyw7x8aahy93w96yhwm7xcq3ke4f8ge93u
	// tbnb prefix address: tbnb1wxeplyw7x8aahy93w96yhwm7xcq3ke4ffasp3d
	atomicSwapCoinsAccAddr := sdk.AccAddress(tmCrypto.AddressHash([]byte("BinanceChainAtomicSwapCoins")))
	// bnb prefix address: bnb1hn8ym9xht925jkncjpf7lhjnax6z8nv24fv2yq
	// tbnb prefix address: tbnb1hn8ym9xht925jkncjpf7lhjnax6z8nv2mu9wy3
	timeLockCoinsAccAddr := sdk.AccAddress(tmCrypto.AddressHash([]byte("BinanceChainTimeLockCoins")))
	// nil address
	emptyAccAddr := sdk.AccAddress(tmCrypto.AddressHash([]byte(nil)))
	// 0x0000... address
	zeroAccAddr, err := sdk.AccAddressFromHex("0000000000000000000000000000000000000000")
	if err != nil {
		panic(err)
	}

	logger.Info().Str("depositedCoinsAccAddr:", depositedCoinsAccAddr.String()).
		Str("delegationAccAddr:", delegationAccAddr.String()).
		Str("pegAccount:", pegAccount.String()).
		Str("atomicSwapCoinsAccAddr:", atomicSwapCoinsAccAddr.String()).
		Str("timeLockCoinsAccAddr:", timeLockCoinsAccAddr.String()).
		Str("emptyAccAddr:", emptyAccAddr.String()).
		Str("zeroAccAddr:", zeroAccAddr.String()).
		Msg("escrow accounts")

	escrowAccs[depositedCoinsAccAddr.String()] = struct{}{}
	escrowAccs[delegationAccAddr.String()] = struct{}{}
	escrowAccs[pegAccount.String()] = struct{}{}
	escrowAccs[atomicSwapCoinsAccAddr.String()] = struct{}{}
	escrowAccs[timeLockCoinsAccAddr.String()] = struct{}{}
	escrowAccs[emptyAccAddr.String()] = struct{}{}
	escrowAccs[zeroAccAddr.String()] = struct{}{}

	return escrowAccs
}

func isEmptyState(home string) (bool, error) {
	files, err := os.ReadDir(path.Join(home, "data"))
	if err != nil {
		return false, err
	}

	// only priv_validator_state.json is created
	return len(files) == 1 && files[0].Name() == "priv_validator_state.json", nil
}

func openDB(rootDir string) (dbm.DB, error) {
	dataDir := filepath.Join(rootDir, "data")
	db, err := dbm.NewGoLevelDB("application", dataDir)
	return db, err
}
