package approval

import (
	"errors"
	"math/big"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/bnb-chain/node/app"
	"github.com/bnb-chain/node/plugins/airdrop"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"

	"github.com/bnb-chain/airdrop-service/internal/config"
	"github.com/bnb-chain/airdrop-service/internal/store"
	"github.com/bnb-chain/airdrop-service/pkg/keymanager"
)

const (
	RequestTypeClaim = "claim"
)

type ApprovalService struct {
	config           *config.Config
	merkleRoot       []byte
	km               keymanager.KeyManager
	store            store.Store
	accountWhiteList map[string]struct{}

	logger *zerolog.Logger
}

func NewApprovalService(config *config.Config, km keymanager.KeyManager, store store.Store, logger *zerolog.Logger) (*ApprovalService, error) {
	accountWhiteList := make(map[string]struct{})
	for _, addr := range config.AccountWhiteList {
		accountWhiteList[addr] = struct{}{}
	}
	merkleRoot, err := hexutil.Decode(config.MerkleRoot)
	if err != nil {
		return nil, err
	}
	return &ApprovalService{km: km, store: store, config: config, merkleRoot: merkleRoot, accountWhiteList: accountWhiteList, logger: logger}, nil
}

func (svc *ApprovalService) checkWhiteList(acc types.AccAddress) bool {
	if len(svc.accountWhiteList) == 0 {
		return true
	}

	_, ok := svc.accountWhiteList[acc.String()]
	return ok
}

func (svc *ApprovalService) GetClaimApproval(req *GetClaimApprovalRequest) (resp *GetClaimApprovalResponse, err error) {
	ownerPubKeyBytes, err := hexutil.Decode(req.OwnerPubKey)
	if err != nil {
		return nil, err
	}
	ownerAddr, err := svc.getAddressFromPubKey(ownerPubKeyBytes)
	if err != nil {
		return nil, err
	}
	ownerSignature, err := hexutil.Decode(req.OwnerSignature)
	if err != nil {
		return nil, err
	}

	svc.logger.Info().Str("address", ownerAddr.String()).Msg("GetClaimApproval")
	// Check While List
	if !svc.checkWhiteList(ownerAddr) {
		return nil, errors.New("address is not in while list")
	}

	// Get Merkle Proofs and Node
	proofs, err := svc.store.GetAccountAssetProof(ownerAddr, req.TokenSymbol)
	if err != nil {
		return nil, err
	}
	account, err := svc.store.GetAccountByAddress(ownerAddr)
	if err != nil {
		return nil, err
	}
	svc.logger.Debug().Interface("account", account).Msg("GetAccountByAddress")
	svc.logger.Debug().Interface("proofs", proofs).Msg("GetAccountAssetProofs")
	// Check if token amount is zero
	if account.SummaryCoins.AmountOf(req.TokenSymbol) == 0 {
		return nil, errors.New("token amount is zero")
	}
	// Verify user signature
	approvalMsg := airdrop.NewAirdropApprovalMsg(req.TokenSymbol, uint64(account.SummaryCoins.AmountOf(req.TokenSymbol)), req.ClaimAddress.Hex())
	msgBytes, err := svc.getStdMsgBytes(approvalMsg)
	if err != nil {
		return nil, err
	}
	svc.logger.Debug().Str("msg", string(msgBytes)).Msg("GetStdMsgBytes")
	err = svc.verifyTmSignature(ownerPubKeyBytes, ownerSignature, msgBytes)
	if err != nil {
		return nil, err
	}

	nodeBytes, err := account.Serialize(req.TokenSymbol)
	if err != nil {
		return nil, err
	}

	var tokenSymbolBytes [32]byte
	copy(tokenSymbolBytes[:], []byte(req.TokenSymbol))

	signData := make([][]byte, 0, len(proofs)+5)
	signData = append(signData, [][]byte{
		[]byte(svc.config.ChainID), req.ClaimAddress[:], ownerSignature, nodeBytes,
		svc.merkleRoot,
	}...)
	signData = append(signData, proofs...)

	approvalSignature, err := svc.km.Sign(crypto.Keccak256(signData...))
	if err != nil {
		return nil, err
	}
	svc.logger.Debug().Bytes("approval_signature", approvalSignature).Msg("Signed ApprovalSignature")
	return &GetClaimApprovalResponse{
		Amount:            big.NewInt(account.SummaryCoins.AmountOf(req.TokenSymbol)),
		Proofs:            proofs,
		ApprovalSignature: approvalSignature,
	}, nil
}

func (svc *ApprovalService) getStdMsgBytes(msg types.Msg) ([]byte, error) {
	cdc := app.Codec
	builder := authtxb.NewTxBuilderFromCLI().WithCodec(cdc).WithChainID(svc.config.ChainID)
	stdMsg, err := builder.Build([]types.Msg{msg})
	if err != nil {
		return nil, err
	}

	return stdMsg.Bytes(), nil
}

func (svc *ApprovalService) verifyTmSignature(pubKeyBytes, signatureBytes, msgBytes []byte) error {
	pubKey := secp256k1.PubKeySecp256k1(pubKeyBytes)

	ok := pubKey.VerifyBytes(msgBytes, signatureBytes)
	if !ok {
		return errors.New("verify signature failed")
	}
	return nil
}

func (svc ApprovalService) getAddressFromPubKey(pubKeyBytes []byte) (types.AccAddress, error) {
	pubKey := secp256k1.PubKeySecp256k1(pubKeyBytes)
	return types.AccAddress(pubKey.Address()), nil
}
