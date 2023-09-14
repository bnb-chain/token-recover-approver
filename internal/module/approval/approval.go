package approval

import (
	"github.com/bnb-chain/airdrop-service/internal/store"
	"github.com/bnb-chain/airdrop-service/pkg/keymanager"
)

type ApprovalService struct {
	km    keymanager.KeyManager
	store store.Store
}

func NewApprovalService(km keymanager.KeyManager, store store.Store) *ApprovalService {
	return &ApprovalService{km: km, store: store}
}

func (svc *ApprovalService) GetClaimApproval(request *GetClaimApprovalRequest) (response *GetClaimApprovalResponse, err error) {
	// TODO: Implement
	return
}

func (svc *ApprovalService) GetRegisterTokenApproval(request *GetRegisterTokenApprovalRequest) (response *GetRegisterTokenApprovalResponse, err error) {
	// TODO: Implement
	return
}
