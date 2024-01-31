package metrics

type Metrics interface {
	IncApprovalCount()
	IncApprovalErrorCount()
	ObserveApprovalDuration(time float64)
	ObserveMerkleProofVerificationDuration(time float64)
	ObserveGetProofDataDuration(time float64)
}
