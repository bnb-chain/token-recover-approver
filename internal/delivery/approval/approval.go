package approval

import "github.com/bnb-chain/airdrop-service/pkg/keymanager"

type ApprovalService struct {
	km keymanager.KeyManager
}

func NewApprovalService(km keymanager.KeyManager) *ApprovalService {
	return &ApprovalService{km: km}
}

func (svc *ApprovalService) GetClaimApproval(request *GetClaimApprovalRequest) (response *GetClaimApprovalResponse, err error) {
	// TODO: Implement
	return
}

func (svc *ApprovalService) GetRegisterTokenApproval(request *GetRegisterTokenApprovalRequest) (response *GetRegisterTokenApprovalResponse, err error) {
	// TODO: Implement
	return
}
