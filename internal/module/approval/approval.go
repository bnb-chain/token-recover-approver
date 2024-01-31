package approval

import (
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/bnb-chain/node/app"
	"github.com/bnb-chain/node/plugins/recover"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"

	"github.com/bnb-chain/token-recover-approver/internal/config"
	"github.com/bnb-chain/token-recover-approver/internal/metrics"
	"github.com/bnb-chain/token-recover-approver/internal/store"
	"github.com/bnb-chain/token-recover-approver/pkg/keymanager"
	"github.com/bnb-chain/token-recover-approver/pkg/util"
)

type ApprovalService struct {
	config           *config.Config
	merkleRoot       []byte
	km               keymanager.KeyManager
	store            store.Store
	accountWhiteList map[string]struct{}

	metrics metrics.Metrics
	logger  *zerolog.Logger
}

func NewApprovalService(config *config.Config, km keymanager.KeyManager, store store.Store, metrics metrics.Metrics, logger *zerolog.Logger) (*ApprovalService, error) {
	accountWhiteList := make(map[string]struct{})
	for _, addr := range config.AccountWhiteList {
		accountWhiteList[addr] = struct{}{}
	}
	merkleRoot, err := hexutil.Decode(config.MerkleRoot)
	if err != nil {
		return nil, err
	}
	return &ApprovalService{km: km, store: store, config: config, merkleRoot: merkleRoot, accountWhiteList: accountWhiteList, metrics: metrics, logger: logger}, nil
}

func (svc *ApprovalService) checkWhiteList(acc types.AccAddress) bool {
	if len(svc.accountWhiteList) == 0 {
		return true
	}

	_, ok := svc.accountWhiteList[acc.String()]
	return ok
}

func (svc *ApprovalService) GetTokenRecoverApproval(req *GetTokenRecoverApprovalRequest) (*GetTokenRecoverApprovalResponse, error) {
	approvalStartTime := time.Now()
	ownerPubKeyBytes, err := hexutil.Decode(req.OwnerPubKey)
	if err != nil {
		svc.metrics.IncApprovalErrorCount()
		return nil, err
	}
	ownerAddr, err := svc.getAddressFromPubKey(ownerPubKeyBytes)
	if err != nil {
		svc.metrics.IncApprovalErrorCount()
		return nil, err
	}
	ownerSignature, err := hexutil.Decode(req.OwnerSignature)
	if err != nil {
		svc.metrics.IncApprovalErrorCount()
		return nil, err
	}

	svc.logger.Info().Str("address", ownerAddr.String()).Msg("GetTokenRecoverApproval")
	// Check While List
	if !svc.checkWhiteList(ownerAddr) {
		svc.metrics.IncApprovalErrorCount()
		return nil, errors.New("address is not in while list")
	}

	// Get Merkle Proofs and Node
	getProofsStartTime := time.Now()
	proof, err := svc.store.GetAccountAssetProof(ownerAddr, req.TokenSymbol)
	if err != nil {
		svc.metrics.IncApprovalErrorCount()
		return nil, err
	}
	svc.metrics.ObserveGetProofDataDuration(float64(time.Since(getProofsStartTime).Seconds()))
	svc.logger.Debug().Interface("proof", proof).Msg("GetAccountAssetProofs")

	// Check if token amount is zero
	if proof.Amount == 0 {
		svc.metrics.IncApprovalErrorCount()
		return nil, errors.New("token amount is zero")
	}
	// Verify user signature
	tokenRecoverRequestMsg := recover.NewTokenRecoverRequestMsg(req.TokenSymbol, uint64(proof.Amount), strings.ToLower(req.ClaimAddress.Hex()))
	msgBytes, err := svc.getStdMsgBytes(tokenRecoverRequestMsg)
	if err != nil {
		svc.metrics.IncApprovalErrorCount()
		return nil, err
	}
	svc.logger.Debug().Str("msg", string(msgBytes)).Msg("GetStdMsgBytes")
	err = svc.verifyTmSignature(ownerPubKeyBytes, ownerSignature, msgBytes)
	if err != nil {
		svc.metrics.IncApprovalErrorCount()
		return nil, err
	}

	nodeBytes, err := proof.Serialize()
	if err != nil {
		svc.metrics.IncApprovalErrorCount()
		return nil, err
	}
	svc.logger.Debug().Str("leaf", hexutil.Encode(nodeBytes)).Msg("LeafNodeBytes")

	// Verify merkle proof
	merkleProofVerificationStartTime := time.Now()
	if !util.VerifyMerkleProof(svc.merkleRoot, proof.Proof, nodeBytes) {
		svc.metrics.IncApprovalErrorCount()
		return nil, errors.New("verify merkle proof failed")
	}
	svc.metrics.ObserveMerkleProofVerificationDuration(float64(time.Since(merkleProofVerificationStartTime).Seconds()))

	// Sign ApprovalSignature
	var tokenSymbolBytes [32]byte
	copy(tokenSymbolBytes[:], []byte(req.TokenSymbol))

	signData := make([][]byte, 0, len(proof.Proof)+5)
	signData = append(signData, [][]byte{
		[]byte(svc.config.ChainID), req.ClaimAddress[:], ownerSignature, nodeBytes,
		svc.merkleRoot,
	}...)
	signData = append(signData, proof.Proof...)

	approvalSignature, err := svc.km.Sign(crypto.Keccak256(signData...))
	if err != nil {
		svc.metrics.IncApprovalErrorCount()
		return nil, err
	}
	svc.logger.Debug().Str("approval_signature", hexutil.Encode(approvalSignature)).Msg("Signed ApprovalSignature")

	svc.metrics.IncApprovalCount()
	svc.metrics.ObserveApprovalDuration(float64(time.Since(approvalStartTime).Seconds()))
	return &GetTokenRecoverApprovalResponse{
		Amount:            big.NewInt(proof.Amount),
		Proofs:            proof.Proof,
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

func (svc *ApprovalService) getAddressFromPubKey(pubKeyBytes []byte) (types.AccAddress, error) {
	pubKey := secp256k1.PubKeySecp256k1(pubKeyBytes)
	return types.AccAddress(pubKey.Address()), nil
}
