package main

import (
	"flag"
	"fmt"
	"ghostnet/pkg/ai"
	"ghostnet/pkg/common"
	"ghostnet/pkg/dashboard"
	"ghostnet/pkg/performance"
	"ghostnet/pkg/pipeline"
	"ghostnet/pkg/research"
	"math/rand"
	"sync"
	"time"
)

func main() {
	userCount := flag.Int("users", 100, "Number of users to simulate")
	duration := flag.Int("duration", 10, "Duration of simulation in seconds")
	privacyMode := flag.String("mode", "MEDIUM", "Privacy mode (LOW, MEDIUM, HIGH, MAXIMUM)")
	flag.Parse()

	fmt.Printf("Starting GhostNet Simulation with %d users in %s mode...\n", *userCount, *privacyMode)

	p := pipeline.NewPipeline()
	mode := common.PrivacyMode(*privacyMode)
	p.CoverEngine.SetMode(mode)
	p.CoverEngine.SetActiveUsers(*userCount)
	p.CoverEngine.Start()

	aiEngine := ai.NewEngine("balanced")
	perfEngine := performance.NewEngine()
	resMode := &research.Mode{Enabled: true}

	var wg sync.WaitGroup
	stop := make(chan struct{})

	// Statistics
	var totalProcessed int64
	var mu sync.Mutex

	// Start users
	for i := 0; i < *userCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sites := []string{"google.com", "api.github.com", "netflix.com", "dns.google"}
			for {
				select {
				case <-stop:
					return
				case <-time.After(time.Duration(rand.Intn(1000)) * time.Millisecond):
					site := sites[rand.Intn(len(sites))]
					p.ProcessRequest(site, []byte("SIMULATED TRAFFIC DATA"))
					mu.Lock()
					totalProcessed++
					mu.Unlock()
				}
			}
		}(i)
	}

	// Consumer of final output
	go func() {
		for range p.FinalOut {
			// Just consume
		}
	}()

	// Reporting loop
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				metrics := dashboard.CalculateMetrics(*userCount, p.QueueMgr, mode)
				metrics.RiskLevel = aiEngine.EstimateRisk(*userCount, mode)
				metrics.AIStrategy = aiEngine.RecommendStrategy(metrics.RiskLevel)
				metrics.TotalPackets = int(totalProcessed)
				metrics.Display()

				perfEngine.OptimizeQueues(len(p.QueueMgr.GetStats().QueueLengths))
				fmt.Println(resMode.EvaluateStrategy("Adaptive batching v2"))
			}
		}
	}()

	time.Sleep(time.Duration(*duration) * time.Second)
	close(stop)
	wg.Wait()

	fmt.Println("Simulation completed.")
}
