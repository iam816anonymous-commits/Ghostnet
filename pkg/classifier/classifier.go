package classifier

import (
	"ghostnet/pkg/common"
	"strings"
	"sync"
)

type Policy struct {
	MergeThreshold int
	SplitThreshold int
}

type Classifier struct {
	mu           sync.RWMutex
	policy       Policy
	queueVolumes map[common.TrafficClass]int
}

func NewClassifier() *Classifier {
	return &Classifier{
		policy: Policy{
			MergeThreshold: 10,
			SplitThreshold: 100,
		},
		queueVolumes: make(map[common.TrafficClass]int),
	}
}

func (c *Classifier) Classify(p *common.Packet) common.TrafficClass {
	// Simple heuristic-based classification for simulation
	size := p.Size
	payload := string(p.Payload)

	var class common.TrafficClass
	switch {
	case strings.Contains(payload, "GET") && size < 5000:
		class = common.ClassSmallWeb
	case strings.Contains(payload, "API"):
		class = common.ClassAPI
	case size > 1000000:
		class = common.ClassVideo
	case size > 100000:
		class = common.ClassImage // Use ClassImage for medium-large
	case size > 10000:
		class = common.ClassFileUpload
	case strings.Contains(payload, "DNS"):
		class = common.ClassDNS
	case size < 1000:
		class = common.ClassBackground
	default:
		class = common.ClassInteractive
	}

	c.mu.Lock()
	c.queueVolumes[class]++
	c.mu.Unlock()

	return class
}

func (c *Classifier) UpdatePolicy(newPolicy Policy) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.policy = newPolicy
}

func (c *Classifier) GetQueues() []common.TrafficClass {
	c.mu.Lock()
	defer c.mu.Unlock()

	var activeClasses []common.TrafficClass
	for class, vol := range c.queueVolumes {
		// Split logic: if volume is very high, keep separate
		if vol > c.policy.SplitThreshold {
			activeClasses = append(activeClasses, class)
			continue
		}

		// Merge logic: if volume is low, merge into Unknown/General pool to increase anonymity set
		if vol < c.policy.MergeThreshold && vol > 0 {
			c.queueVolumes[common.ClassUnknown] += vol
			c.queueVolumes[class] = 0
		}

		if vol > 0 {
			activeClasses = append(activeClasses, class)
		}
	}

	if c.queueVolumes[common.ClassUnknown] > 0 {
		activeClasses = append(activeClasses, common.ClassUnknown)
	}

	return activeClasses
}
