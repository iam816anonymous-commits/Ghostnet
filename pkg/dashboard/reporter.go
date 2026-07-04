package dashboard

import (
	"fmt"
	"ghostnet/pkg/common"
	"ghostnet/pkg/queue"
	"time"
)

type Metrics struct {
	PrivacyScore     float64
	RiskLevel        float64
	AIStrategy       string
	Latency          time.Duration
	AnonymitySetSize int
	CoverTraffic     int
	TotalPackets     int
	QueueStats       queue.Stats
}

func (m Metrics) Display() {
	fmt.Println("==========================================")
	fmt.Println("GHOSTNET PRIVACY DASHBOARD")
	fmt.Println("==========================================")
	fmt.Printf("Privacy Score:       %.2f/10.0\n", m.PrivacyScore)
	fmt.Printf("Risk Level:          %.2f/10.0\n", m.RiskLevel)
	fmt.Printf("AI Recommendation:   %s\n", m.AIStrategy)
	fmt.Printf("Latency:             %v\n", m.Latency)
	fmt.Printf("Anonymity Set:       %d users\n", m.AnonymitySetSize)
	fmt.Printf("Cover Traffic:       %d pkts\n", m.CoverTraffic)
	fmt.Printf("Total Packets:       %d\n", m.TotalPackets)
	fmt.Printf("Queue Status:        %v\n", m.QueueStats)
	fmt.Println("==========================================")
}

func CalculateMetrics(activeUsers int, qMgr *queue.Manager, privacyMode common.PrivacyMode) Metrics {
	stats := qMgr.GetStats()

	score := 5.0
	if privacyMode == common.ModeMaximumPrivacy {
		score += 4.0
	} else if privacyMode == common.ModeHigh {
		score += 2.0
	}

	if activeUsers > 1000 {
		score += 1.0
	}

	return Metrics{
		PrivacyScore:     score,
		Latency:          250 * time.Millisecond,
		AnonymitySetSize: activeUsers,
		QueueStats:       stats,
	}
}
