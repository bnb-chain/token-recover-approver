package approval

import (
	"math/big"
	"path"
	"reflect"
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"

	"github.com/bnb-chain/token-recover-approver/internal/config"
	collector "github.com/bnb-chain/token-recover-approver/internal/metrics/prometheus"
	"github.com/bnb-chain/token-recover-approver/internal/store"
	"github.com/bnb-chain/token-recover-approver/internal/store/memory"
	"github.com/bnb-chain/token-recover-approver/pkg/keymanager/local"
	"github.com/bnb-chain/token-recover-approver/pkg/util"
)

const (
	approvalPrivKey  = "afc2986f283cf5f9d17e04c6a12ccf8fa46149fc37d48e11abef15a46ae34eb7"
	mockDataBasePath = "../../../example/store"
	mockMerkleRoot   = "0xad78b6dbdb34cb9a5c0e44fbfcc3bd52d6e9d519eefbfc39ba5c3b232849a064"
)

func makeMockStore() (store.Store, error) {
	initSDK()
	return memory.NewMemoryStore(
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
	}, km, mockStore, collector.NewCollector(prometheus.NewRegistry()), &zerolog.Logger{})
}

func initSDK() {
	sdkConfig := types.GetConfig()
	sdkConfig.SetBech32PrefixForAccount("tbnb", "bnbp")
}

func TestApprovalService_GetTokenRecoverApproval(t *testing.T) {
	type args struct {
		req *GetTokenRecoverApprovalRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *GetTokenRecoverApprovalResponse
		wantErr  bool
	}{
		{
			name: "test case 1",
			args: args{
				req: &GetTokenRecoverApprovalRequest{
					TokenSymbol:    "BNB",
					OwnerPubKey:    "0x02dcd743516b78366a217a1bf2aa562ec5accd07163db3332d924fa48e643875a6",
					OwnerSignature: "0xcd32af98a3cf4b66deaba53dc81c7cf8c810a83eb2fa23bf1a555a718826e2f03d47e3711a10e1ae72031fcd3faabac51325999d74c0cff31d554b4d657dbc64",
					ClaimAddress:   common.HexToAddress("0x5B38Da6a701c568545dCfcB03FcB875f56beddC4"),
				},
			},
			wantResp: &GetTokenRecoverApprovalResponse{
				Amount:            big.NewInt(19999999000000000),
				Proofs:            util.MustDecodeHexArrayToBytes([]string{"0x3262127e4ff0bce1bb67e569baa034637806b4519b19d3ad9dbae7f5ad31fa18"}),
				ApprovalSignature: util.MustDecodeHexToBytes("0xe2b44f23b1e8713419a8c1a881c35df72b2f3e55696a8601bb0356529d02a79113590c6d32e8a0b22b609947b1ab84f439037dd99f20bbe738034438ff301f9800"),
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
			gotResp, err := svc.GetTokenRecoverApproval(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApprovalService.GetTokenRecoverApproval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("ApprovalService.GetTokenRecoverApproval() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
