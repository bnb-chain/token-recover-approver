package approval

import (
	"math/big"
	"path"
	"reflect"
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"

	"github.com/bnb-chain/airdrop-service/internal/config"
	"github.com/bnb-chain/airdrop-service/internal/store"
	"github.com/bnb-chain/airdrop-service/internal/store/memory"
	"github.com/bnb-chain/airdrop-service/pkg/keymanager/local"
)

const (
	approvalPrivKey  = "afc2986f283cf5f9d17e04c6a12ccf8fa46149fc37d48e11abef15a46ae34eb7"
	mockDataBasePath = "../../../example/store"
)

func makeMockStore() (store.Store, error) {
	initSDK()
	return memory.NewMemoryStore(
		path.Join(mockDataBasePath, "state_root.json"),
		path.Join(mockDataBasePath, "assets.json"),
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
		ChainID: "Binance-Chain-Ganges",
	}, km, mockStore, &zerolog.Logger{}), nil
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
					TokenIndex:     0,
					TokenSymbol:    "BNB",
					OwnerPubKey:    "0278caa4d6321aa856d6341dd3e8bcdfe0b55901548871c63c3f5cec43c2ae88a9",
					OwnerSignature: "a7af9a82c98d14de0e7d071b67b740e6a4e265cf787e97f21d628840d2a2d33a41e1d550bd57009b153a0c9ce228a7e0e02b3a8f53e18699b2f13fa13bb809c7",
					ClaimAddress:   common.HexToAddress("0x5B38Da6a701c568545dCfcB03FcB875f56beddC4"),
				},
			},
			wantResp: &GetClaimApprovalResponse{
				Amount:            big.NewInt(1000000000),
				Proofs:            []string{"0x5bb1c3643cde99e00a7f32707bad4b06c8ec5b1e4097a97aa3cf7fabfc0b92f7"},
				ApprovalSignature: "e309cf41cdc9197e636b5f52eb800fc21817443a331179517ca58a4c69f9139a",
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
