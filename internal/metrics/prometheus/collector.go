package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"

	"github.com/bnb-chain/token-recover-approver/internal/metrics"
	"github.com/bnb-chain/token-recover-approver/internal/version"
)

var _ metrics.Metrics = (*Collector)(nil)

func NewCollector(registry *prometheus.Registry) *Collector {
	// Add go runtime metrics and process collectors.
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	// Add custom metrics.
	versionInfo := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "version_info",
		Help: "Version information about this binary",
		ConstLabels: prometheus.Labels{
			"version":    version.AppVersion,
			"git_commit": version.GitCommit,
			"build_date": version.GitCommitDate,
		},
	})
	approvalCount := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "approval_count",
		Help: "The total number of approvals.",
	})
	approvalErrorCount := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "approval_error_count",
		Help: "The total number of approval errors.",
	})
	approvalDuration := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "approval_duration_seconds",
			Help: "A summary of the approval duration of the requests.",
			Objectives: map[float64]float64{
				0:    0.0,
				0.25: 0.05,
				0.5:  0.001,
				0.75: 0.001,
				1:    0.0,
			},
		},
	)
	getProofDataDuration := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "get_proof_data_duration_seconds",
			Help: "A summary of the get proof data duration of the requests.",
			Objectives: map[float64]float64{
				0:    0.0,
				0.25: 0.05,
				0.5:  0.001,
				0.75: 0.001,
				1:    0.0,
			},
		},
	)
	merkleProofVerificationDuration := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "merkle_proof_verification_duration_seconds",
			Help: "A summary of the merkle proof verification duration of the requests.",
			Objectives: map[float64]float64{
				0:    0.0,
				0.25: 0.05,
				0.5:  0.001,
				0.75: 0.001,
				1:    0.0,
			},
		},
	)

	registry.MustRegister(
		versionInfo,
		approvalCount,
		approvalErrorCount,
		approvalDuration,
		getProofDataDuration,
		merkleProofVerificationDuration,
	)

	return &Collector{
		approvalCount:                   approvalCount,
		approvalErrorCount:              approvalErrorCount,
		approvalDuration:                approvalDuration,
		getProofDataDuration:            getProofDataDuration,
		merkleProofVerificationDuration: merkleProofVerificationDuration,
	}
}

type Collector struct {
	approvalCount                   prometheus.Counter
	approvalErrorCount              prometheus.Counter
	approvalDuration                prometheus.Summary
	getProofDataDuration            prometheus.Summary
	merkleProofVerificationDuration prometheus.Summary
}

// ObserveApprovalDuration implements metrics.Metrics.
func (c *Collector) ObserveApprovalDuration(time float64) {
	c.approvalDuration.Observe(time)
}

// ObserveGetProofDataDuration implements metrics.Metrics.
func (c *Collector) ObserveGetProofDataDuration(time float64) {
	c.getProofDataDuration.Observe(time)
}

// ObserveMerkleProofVerificationDuration implements metrics.Metrics.
func (c *Collector) ObserveMerkleProofVerificationDuration(time float64) {
	c.merkleProofVerificationDuration.Observe(time)
}

// IncApprovalCount implements metrics.Metrics.
func (c *Collector) IncApprovalCount() {
	c.approvalCount.Inc()
}

// IncApprovalErrorCount implements metrics.Metrics.
func (c *Collector) IncApprovalErrorCount() {
	c.approvalErrorCount.Inc()
}
