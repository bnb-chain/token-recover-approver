package approval

import (
	"math/big"
	"path"
	"reflect"
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"

	"github.com/bnb-chain/token-recover-approver/internal/config"
	"github.com/bnb-chain/token-recover-approver/internal/store"
	"github.com/bnb-chain/token-recover-approver/internal/store/memory"
	"github.com/bnb-chain/token-recover-approver/pkg/keymanager/local"
	"github.com/bnb-chain/token-recover-approver/pkg/util"
)

const (
	approvalPrivKey  = "afc2986f283cf5f9d17e04c6a12ccf8fa46149fc37d48e11abef15a46ae34eb7"
	mockDataBasePath = "../../../example/store"
	mockMerkleRoot   = "0x5bd43c1c0929f259349cdf93b3b28673e3c98882ae607098a65798efbebfe39c"
)

func makeMockStore() (store.Store, error) {
	initSDK()
	return memory.NewMemoryStore(
		path.Join(mockDataBasePath, "accounts.json"),
		path.Join(mockDataBasePath, "merkle_proofs.json"),
	)
}

func makeMockSvc() (*ApprovalService, error) {
	km, err := local.NewLocalKeyManager(approvalPrivKey)
	if err != nil {
		return nil, err
	}
	mockStore, err := makeMockStore()
	if err != nil {
		return nil, err
	}

	return NewApprovalService(&config.Config{
		ChainID:    "Binance-Chain-Ganges",
		MerkleRoot: mockMerkleRoot,
	}, km, mockStore, &zerolog.Logger{})
}

func initSDK() {
	sdkConfig := types.GetConfig()
	sdkConfig.SetBech32PrefixForAccount("tbnb", "bnbp")
}

func TestApprovalService_GetClaimApproval(t *testing.T) {
	type args struct {
		req *GetClaimApprovalRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *GetClaimApprovalResponse
		wantErr  bool
	}{
		{
			name: "test case 1",
			args: args{
				req: &GetClaimApprovalRequest{
					TokenSymbol:    "BNB",
					OwnerPubKey:    "0x0278caa4d6321aa856d6341dd3e8bcdfe0b55901548871c63c3f5cec43c2ae88a9",
					OwnerSignature: "0x61f2662fbf581f8cd34ee16e025fab50fbb2d481dbdab4bbe49f6f5ee47fba4c0f54cd0d4994a31573aef15fcaaf3a91eaad95a10183678d563e15128af3e173",
					ClaimAddress:   common.HexToAddress("0x5B38Da6a701c568545dCfcB03FcB875f56beddC4"),
				},
			},
			wantResp: &GetClaimApprovalResponse{
				Amount:            big.NewInt(1000000000),
				Proofs:            util.MustDecodeHexArrayToBytes([]string{"0x679c555951fde6f1e516549283ef67bd4f32c9058f72e41e3cacdfc337410f3e"}),
				ApprovalSignature: util.MustDecodeHexToBytes("0x207e4b31476bedf9b73ccf01460a138dd12ca1b1ab65f34d7748ba49aa7d3f7c6a19e1d26c04506eba2fb9e29d91845ad4a39f447d6f959ac56628c6370b714200"),
			},
			wantErr: false,
		},
	}
	svc, err := makeMockSvc()
	if err != nil {
		t.Fatal(err)
		return
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := svc.GetClaimApproval(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApprovalService.GetClaimApproval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("ApprovalService.GetClaimApproval() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
