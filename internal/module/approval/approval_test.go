package approval

import (
	"fmt"
	"math/big"
	"path"
	"reflect"
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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
	mockMerkleRoot   = "0x59bb94f7047904a8fdaec42e4785295167f7fd63742b309afeb84bd71f8e6554"
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
			name: "test case 1 - recover BNB",
			args: args{
				req: &GetTokenRecoverApprovalRequest{
					TokenSymbol:    "BNB",
					OwnerPubKey:    "0x036d5d41cd7da2e96d39bcbd0390bfed461a86382f7a2923436ff16c65cabc7719",
					OwnerSignature: "0x5f5391ba7f2b002b4746025f7e803a43e57a397ea66f3939d05302eb7851bbbc0773cda87aae0fbb1e2a29367b606209ed47dc5cba6d1a83f6b79cb70e56efdb",
					ClaimAddress:   common.HexToAddress("0x2e9247B67ae885a8dcfBf77Eb6d0e93A32bea24C"),
				},
			},
			wantResp: &GetTokenRecoverApprovalResponse{
				Amount: big.NewInt(14188000000),
				Proofs: util.MustDecodeHexArrayToBytes([]string{"0x03719d7863e4aba727d7030e7a1916b9be2245d447eb71fc683d3ac0ded5eecd",
					"0x7f9aa9d8246251cbab3cc642416dec81d074d39a85be6ca8326a05ac422e74ab",
					"0x6debec5a4272951843cf24f74c30d5ccf1afec9aafbfc45d0b50cb4eb6f89c09",
					"0x5cb2e4d880e2387764df4de9ce49cbabc41b6e4a07b1c2e1d9fc98957b6643d2",
					"0x88c6195b4444035bef3212847f38822c0d509d811de8c9154e7f5f8ec3778b67",
					"0x27c985cced25522043ded2fc8103baa24edc21b6c9f95c5bfff635ab36bdb29d",
					"0x39a0fbfba925ebd0cf4f5fe5ab4c69eb18317fd1bd4373647a53dc339fb764a9",
					"0x61300a7a7fe0932760c1e1edfa4d4450cc378d9b5c538dcb24ffbbc18f249fe5",
					"0x4d49fcf8a1e0b72b535921dea8e02baac18df614e7f7c462749a2b14ee2737ef",
					"0xc10261d3337346f921c4fef13ba1bcb46a531e947ce41c81e54404e970deaaf5",
					"0x3536a24678835b0f7adeae1f27dae7d6bb22598fb8f8578ec0eef5ea5146f85b",
					"0x925aab793d8080c4f8ea5034e195938c5550f7ba80acf7d7e7d8468f5b5dd70a",
					"0xdef2b6210654ac4f48b4556e24907e027e66729045d0c669a53c75a880477b48",
					"0x4bb1aab890245e6a9e1e969ae3f6f0315ea073606fd6fabe9f3d7514c84fee98",
					"0xe096d4b3669b1c7cd8fcff26b2b00029c09c0f38a34ae632b022622fb46ad69a",
					"0x05e63b558cba63f5add60201151f96ff8f5370d2b8280a96b4fa8fd2d519ab9f",
					"0xa2d456e52facaa953bfbc79a5a6ed7647dda59872b9b35c20183887eeb4640eb"}),
				ApprovalSignature: util.MustDecodeHexToBytes("0x52a0a5ca80beb068d82413cac31c1df0540dc6a61eddec9f31b94419e60b6c586e5342552f4c8034a00c876d640abea8c5ba9c4d72145d0e562fedd09fe1e00a01"),
			},
			wantErr: false,
		},
		{
			name: "test case 2 - recover DYTT991-49A",
			args: args{
				req: &GetTokenRecoverApprovalRequest{
					TokenSymbol:    "DYTT991-49A",
					OwnerPubKey:    "0x036d5d41cd7da2e96d39bcbd0390bfed461a86382f7a2923436ff16c65cabc7719",
					OwnerSignature: "0xfe5cb16008d7afd2723cdaf16649bbbd2635dbfc2764c985847497408485f782562e7efeb7911986f4b8a74a347e49fac4c780dde80bbd4be542d51f7680cf9b",
					ClaimAddress:   common.HexToAddress("0x2e9247B67ae885a8dcfBf77Eb6d0e93A32bea24C"),
				},
			},
			wantResp: &GetTokenRecoverApprovalResponse{
				Amount: big.NewInt(10000000000000000),
				Proofs: util.MustDecodeHexArrayToBytes([]string{"0x061680518f3f97c075a62df766fa55c90b0c415140f737c0d1f7ace5ad2bfee6",
					"0x366f06cef0f1668d848819cb7b5a07b0093ad997da496e60060db2fee754857b",
					"0x6debec5a4272951843cf24f74c30d5ccf1afec9aafbfc45d0b50cb4eb6f89c09",
					"0x5cb2e4d880e2387764df4de9ce49cbabc41b6e4a07b1c2e1d9fc98957b6643d2",
					"0x88c6195b4444035bef3212847f38822c0d509d811de8c9154e7f5f8ec3778b67",
					"0x27c985cced25522043ded2fc8103baa24edc21b6c9f95c5bfff635ab36bdb29d",
					"0x39a0fbfba925ebd0cf4f5fe5ab4c69eb18317fd1bd4373647a53dc339fb764a9",
					"0x61300a7a7fe0932760c1e1edfa4d4450cc378d9b5c538dcb24ffbbc18f249fe5",
					"0x4d49fcf8a1e0b72b535921dea8e02baac18df614e7f7c462749a2b14ee2737ef",
					"0xc10261d3337346f921c4fef13ba1bcb46a531e947ce41c81e54404e970deaaf5",
					"0x3536a24678835b0f7adeae1f27dae7d6bb22598fb8f8578ec0eef5ea5146f85b",
					"0x925aab793d8080c4f8ea5034e195938c5550f7ba80acf7d7e7d8468f5b5dd70a",
					"0xdef2b6210654ac4f48b4556e24907e027e66729045d0c669a53c75a880477b48",
					"0x4bb1aab890245e6a9e1e969ae3f6f0315ea073606fd6fabe9f3d7514c84fee98",
					"0xe096d4b3669b1c7cd8fcff26b2b00029c09c0f38a34ae632b022622fb46ad69a",
					"0x05e63b558cba63f5add60201151f96ff8f5370d2b8280a96b4fa8fd2d519ab9f",
					"0xa2d456e52facaa953bfbc79a5a6ed7647dda59872b9b35c20183887eeb4640eb"}),
				ApprovalSignature: util.MustDecodeHexToBytes("0x693ff2e0458dae7e34a8a1e9929ec122ab3f8224b3bdc0ada23d142ec5e191496fb36c172db406941c55ce715f2f32f54faffcf8551a3dfaf1be045e64a084e900"),
			},
			wantErr: false,
		},
		{
			name: "test case 3 - wrong symbol",
			args: args{
				req: &GetTokenRecoverApprovalRequest{
					TokenSymbol:    "BNBP",
					OwnerPubKey:    "0x036d5d41cd7da2e96d39bcbd0390bfed461a86382f7a2923436ff16c65cabc7719",
					OwnerSignature: "0x5f5391ba7f2b002b4746025f7e803a43e57a397ea66f3939d05302eb7851bbbc0773cda87aae0fbb1e2a29367b606209ed47dc5cba6d1a83f6b79cb70e56efdb",
					ClaimAddress:   common.HexToAddress("0x2e9247B67ae885a8dcfBf77Eb6d0e93A32bea24C"),
				},
			},
			wantResp: &GetTokenRecoverApprovalResponse{
				Amount: big.NewInt(14188000000),
				Proofs: util.MustDecodeHexArrayToBytes([]string{"0x03719d7863e4aba727d7030e7a1916b9be2245d447eb71fc683d3ac0ded5eecd",
					"0x7f9aa9d8246251cbab3cc642416dec81d074d39a85be6ca8326a05ac422e74ab",
					"0x6debec5a4272951843cf24f74c30d5ccf1afec9aafbfc45d0b50cb4eb6f89c09",
					"0x5cb2e4d880e2387764df4de9ce49cbabc41b6e4a07b1c2e1d9fc98957b6643d2",
					"0x88c6195b4444035bef3212847f38822c0d509d811de8c9154e7f5f8ec3778b67",
					"0x27c985cced25522043ded2fc8103baa24edc21b6c9f95c5bfff635ab36bdb29d",
					"0x39a0fbfba925ebd0cf4f5fe5ab4c69eb18317fd1bd4373647a53dc339fb764a9",
					"0x61300a7a7fe0932760c1e1edfa4d4450cc378d9b5c538dcb24ffbbc18f249fe5",
					"0x4d49fcf8a1e0b72b535921dea8e02baac18df614e7f7c462749a2b14ee2737ef",
					"0xc10261d3337346f921c4fef13ba1bcb46a531e947ce41c81e54404e970deaaf5",
					"0x3536a24678835b0f7adeae1f27dae7d6bb22598fb8f8578ec0eef5ea5146f85b",
					"0x925aab793d8080c4f8ea5034e195938c5550f7ba80acf7d7e7d8468f5b5dd70a",
					"0xdef2b6210654ac4f48b4556e24907e027e66729045d0c669a53c75a880477b48",
					"0x4bb1aab890245e6a9e1e969ae3f6f0315ea073606fd6fabe9f3d7514c84fee98",
					"0xe096d4b3669b1c7cd8fcff26b2b00029c09c0f38a34ae632b022622fb46ad69a",
					"0x05e63b558cba63f5add60201151f96ff8f5370d2b8280a96b4fa8fd2d519ab9f",
					"0xa2d456e52facaa953bfbc79a5a6ed7647dda59872b9b35c20183887eeb4640eb"}),
				ApprovalSignature: util.MustDecodeHexToBytes("0x52a0a5ca80beb068d82413cac31c1df0540dc6a61eddec9f31b94419e60b6c586e5342552f4c8034a00c876d640abea8c5ba9c4d72145d0e562fedd09fe1e00a01"),
			},
			wantErr: true,
		},
		{
			name: "test case 4 - wrong owner pubkey",
			args: args{
				req: &GetTokenRecoverApprovalRequest{
					TokenSymbol:    "BNB",
					OwnerPubKey:    "0x11115d41cd7da2e96d39bcbd0390bfed461a86382f7a2923436ff16c65cabc7719",
					OwnerSignature: "0x5f5391ba7f2b002b4746025f7e803a43e57a397ea66f3939d05302eb7851bbbc0773cda87aae0fbb1e2a29367b606209ed47dc5cba6d1a83f6b79cb70e56efdb",
					ClaimAddress:   common.HexToAddress("0x2e9247B67ae885a8dcfBf77Eb6d0e93A32bea24C"),
				},
			},
			wantResp: &GetTokenRecoverApprovalResponse{
				Amount: big.NewInt(14188000000),
				Proofs: util.MustDecodeHexArrayToBytes([]string{"0x03719d7863e4aba727d7030e7a1916b9be2245d447eb71fc683d3ac0ded5eecd",
					"0x7f9aa9d8246251cbab3cc642416dec81d074d39a85be6ca8326a05ac422e74ab",
					"0x6debec5a4272951843cf24f74c30d5ccf1afec9aafbfc45d0b50cb4eb6f89c09",
					"0x5cb2e4d880e2387764df4de9ce49cbabc41b6e4a07b1c2e1d9fc98957b6643d2",
					"0x88c6195b4444035bef3212847f38822c0d509d811de8c9154e7f5f8ec3778b67",
					"0x27c985cced25522043ded2fc8103baa24edc21b6c9f95c5bfff635ab36bdb29d",
					"0x39a0fbfba925ebd0cf4f5fe5ab4c69eb18317fd1bd4373647a53dc339fb764a9",
					"0x61300a7a7fe0932760c1e1edfa4d4450cc378d9b5c538dcb24ffbbc18f249fe5",
					"0x4d49fcf8a1e0b72b535921dea8e02baac18df614e7f7c462749a2b14ee2737ef",
					"0xc10261d3337346f921c4fef13ba1bcb46a531e947ce41c81e54404e970deaaf5",
					"0x3536a24678835b0f7adeae1f27dae7d6bb22598fb8f8578ec0eef5ea5146f85b",
					"0x925aab793d8080c4f8ea5034e195938c5550f7ba80acf7d7e7d8468f5b5dd70a",
					"0xdef2b6210654ac4f48b4556e24907e027e66729045d0c669a53c75a880477b48",
					"0x4bb1aab890245e6a9e1e969ae3f6f0315ea073606fd6fabe9f3d7514c84fee98",
					"0xe096d4b3669b1c7cd8fcff26b2b00029c09c0f38a34ae632b022622fb46ad69a",
					"0x05e63b558cba63f5add60201151f96ff8f5370d2b8280a96b4fa8fd2d519ab9f",
					"0xa2d456e52facaa953bfbc79a5a6ed7647dda59872b9b35c20183887eeb4640eb"}),
				ApprovalSignature: util.MustDecodeHexToBytes("0x52a0a5ca80beb068d82413cac31c1df0540dc6a61eddec9f31b94419e60b6c586e5342552f4c8034a00c876d640abea8c5ba9c4d72145d0e562fedd09fe1e00a01"),
			},
			wantErr: true,
		},
		{
			name: "test case 5 - wrong owner signature",
			args: args{
				req: &GetTokenRecoverApprovalRequest{
					TokenSymbol:    "BNB",
					OwnerPubKey:    "0x036d5d41cd7da2e96d39bcbd0390bfed461a86382f7a2923436ff16c65cabc7719",
					OwnerSignature: "0x111191ba7f2b002b4746025f7e803a43e57a397ea66f3939d05302eb7851bbbc0773cda87aae0fbb1e2a29367b606209ed47dc5cba6d1a83f6b79cb70e56efdb",
					ClaimAddress:   common.HexToAddress("0x2e9247B67ae885a8dcfBf77Eb6d0e93A32bea24C"),
				},
			},
			wantResp: &GetTokenRecoverApprovalResponse{
				Amount: big.NewInt(14188000000),
				Proofs: util.MustDecodeHexArrayToBytes([]string{"0x03719d7863e4aba727d7030e7a1916b9be2245d447eb71fc683d3ac0ded5eecd",
					"0x7f9aa9d8246251cbab3cc642416dec81d074d39a85be6ca8326a05ac422e74ab",
					"0x6debec5a4272951843cf24f74c30d5ccf1afec9aafbfc45d0b50cb4eb6f89c09",
					"0x5cb2e4d880e2387764df4de9ce49cbabc41b6e4a07b1c2e1d9fc98957b6643d2",
					"0x88c6195b4444035bef3212847f38822c0d509d811de8c9154e7f5f8ec3778b67",
					"0x27c985cced25522043ded2fc8103baa24edc21b6c9f95c5bfff635ab36bdb29d",
					"0x39a0fbfba925ebd0cf4f5fe5ab4c69eb18317fd1bd4373647a53dc339fb764a9",
					"0x61300a7a7fe0932760c1e1edfa4d4450cc378d9b5c538dcb24ffbbc18f249fe5",
					"0x4d49fcf8a1e0b72b535921dea8e02baac18df614e7f7c462749a2b14ee2737ef",
					"0xc10261d3337346f921c4fef13ba1bcb46a531e947ce41c81e54404e970deaaf5",
					"0x3536a24678835b0f7adeae1f27dae7d6bb22598fb8f8578ec0eef5ea5146f85b",
					"0x925aab793d8080c4f8ea5034e195938c5550f7ba80acf7d7e7d8468f5b5dd70a",
					"0xdef2b6210654ac4f48b4556e24907e027e66729045d0c669a53c75a880477b48",
					"0x4bb1aab890245e6a9e1e969ae3f6f0315ea073606fd6fabe9f3d7514c84fee98",
					"0xe096d4b3669b1c7cd8fcff26b2b00029c09c0f38a34ae632b022622fb46ad69a",
					"0x05e63b558cba63f5add60201151f96ff8f5370d2b8280a96b4fa8fd2d519ab9f",
					"0xa2d456e52facaa953bfbc79a5a6ed7647dda59872b9b35c20183887eeb4640eb"}),
				ApprovalSignature: util.MustDecodeHexToBytes("0x52a0a5ca80beb068d82413cac31c1df0540dc6a61eddec9f31b94419e60b6c586e5342552f4c8034a00c876d640abea8c5ba9c4d72145d0e562fedd09fe1e00a01"),
			},
			wantErr: true,
		},
		{
			name: "test case 6 - wrong claim address",
			args: args{
				req: &GetTokenRecoverApprovalRequest{
					TokenSymbol:    "BNB",
					OwnerPubKey:    "0x036d5d41cd7da2e96d39bcbd0390bfed461a86382f7a2923436ff16c65cabc7719",
					OwnerSignature: "0xfe5cb16008d7afd2723cdaf16649bbbd2635dbfc2764c985847497408485f782562e7efeb7911986f4b8a74a347e49fac4c780dde80bbd4be542d51f7680cf9b",
					ClaimAddress:   common.HexToAddress("0x561319e67357fa3d2b51E58d011a80EB6268A0f5"),
				},
			},
			wantResp: &GetTokenRecoverApprovalResponse{
				Amount: big.NewInt(14188000000),
				Proofs: util.MustDecodeHexArrayToBytes([]string{"0x03719d7863e4aba727d7030e7a1916b9be2245d447eb71fc683d3ac0ded5eecd",
					"0x7f9aa9d8246251cbab3cc642416dec81d074d39a85be6ca8326a05ac422e74ab",
					"0x6debec5a4272951843cf24f74c30d5ccf1afec9aafbfc45d0b50cb4eb6f89c09",
					"0x5cb2e4d880e2387764df4de9ce49cbabc41b6e4a07b1c2e1d9fc98957b6643d2",
					"0x88c6195b4444035bef3212847f38822c0d509d811de8c9154e7f5f8ec3778b67",
					"0x27c985cced25522043ded2fc8103baa24edc21b6c9f95c5bfff635ab36bdb29d",
					"0x39a0fbfba925ebd0cf4f5fe5ab4c69eb18317fd1bd4373647a53dc339fb764a9",
					"0x61300a7a7fe0932760c1e1edfa4d4450cc378d9b5c538dcb24ffbbc18f249fe5",
					"0x4d49fcf8a1e0b72b535921dea8e02baac18df614e7f7c462749a2b14ee2737ef",
					"0xc10261d3337346f921c4fef13ba1bcb46a531e947ce41c81e54404e970deaaf5",
					"0x3536a24678835b0f7adeae1f27dae7d6bb22598fb8f8578ec0eef5ea5146f85b",
					"0x925aab793d8080c4f8ea5034e195938c5550f7ba80acf7d7e7d8468f5b5dd70a",
					"0xdef2b6210654ac4f48b4556e24907e027e66729045d0c669a53c75a880477b48",
					"0x4bb1aab890245e6a9e1e969ae3f6f0315ea073606fd6fabe9f3d7514c84fee98",
					"0xe096d4b3669b1c7cd8fcff26b2b00029c09c0f38a34ae632b022622fb46ad69a",
					"0x05e63b558cba63f5add60201151f96ff8f5370d2b8280a96b4fa8fd2d519ab9f",
					"0xa2d456e52facaa953bfbc79a5a6ed7647dda59872b9b35c20183887eeb4640eb"}),
				ApprovalSignature: util.MustDecodeHexToBytes("0x52a0a5ca80beb068d82413cac31c1df0540dc6a61eddec9f31b94419e60b6c586e5342552f4c8034a00c876d640abea8c5ba9c4d72145d0e562fedd09fe1e00a01"),
			},
			wantErr: true,
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
			if !tt.wantErr && !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("ApprovalService.GetTokenRecoverApproval() = %v, want %v", gotResp, tt.wantResp)
				return
			}
			// skip merkle proof verification if error
			if err != nil {
				return
			}

			// verify merkle proof
			ownerAddr, err := svc.getAddressFromPubKey(util.MustDecodeHexToBytes(tt.args.req.OwnerPubKey))
			if !tt.wantErr && err != nil {
				t.Errorf("ApprovalService.GetTokenRecoverApproval() error = %v", err)
				return
			}
			leaf := store.Proof{
				Address: ownerAddr,
				Denom:   tt.args.req.TokenSymbol,
				Amount:  tt.wantResp.Amount.Int64(),
			}

			leafBytes, err := leaf.Serialize()
			if !tt.wantErr && err != nil {
				t.Errorf("ApprovalService.GetTokenRecoverApproval() error = %v", err)
				return
			}

			if !tt.wantErr && !util.VerifyMerkleProof(util.MustDecodeHexToBytes(mockMerkleRoot), gotResp.Proofs, leafBytes) {
				t.Errorf("ApprovalService.GetTokenRecoverApproval() error = %v", fmt.Errorf("invalid merkle proof"))
				return
			}

			// verify approval signature is signed by approval address
			signData := make([][]byte, 0, len(gotResp.Proofs)+5)
			signData = append(signData, [][]byte{
				[]byte(svc.config.ChainID), tt.args.req.ClaimAddress[:], util.MustDecodeHexToBytes(tt.args.req.OwnerSignature), leafBytes,
				svc.merkleRoot,
			}...)
			signData = append(signData, gotResp.Proofs...)
			msgHash := crypto.Keccak256(signData...)
			if !tt.wantErr && !svc.km.Verify(msgHash, gotResp.ApprovalSignature) {
				t.Errorf("ApprovalService.GetTokenRecoverApproval() error = %v", fmt.Errorf("invalid approval signature"))
				return
			}
		})
	}
}
