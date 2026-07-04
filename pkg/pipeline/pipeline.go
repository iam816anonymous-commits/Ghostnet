package pipeline

import (
	"fmt"
	"ghostnet/pkg/classifier"
	"ghostnet/pkg/common"
	"ghostnet/pkg/cover"
	"ghostnet/pkg/identity"
	"ghostnet/pkg/normalizer"
	"ghostnet/pkg/queue"
	"ghostnet/pkg/routing"
	"time"
)

type Pipeline struct {
	IdentityMgr   *identity.Manager
	Classifier    *classifier.Classifier
	QueueMgr      *queue.Manager
	CoverEngine   *cover.Engine
	Normalizer    *normalizer.Normalizer
	RoutingEngine *routing.Engine

	ProcessedPackets chan *common.Packet
	FinalOut         chan *common.Packet
}

func NewPipeline() *Pipeline {
	finalOut := make(chan *common.Packet, 1000)
	queueBatchOut := make(chan []*common.Packet, 100)

	p := &Pipeline{
		IdentityMgr:      identity.NewManager(),
		Classifier:       classifier.NewClassifier(),
		QueueMgr:         queue.NewManager(queueBatchOut),
		CoverEngine:      cover.NewEngine(make(chan *common.Packet, 100)),
		Normalizer:       normalizer.NewNormalizer(),
		RoutingEngine:    routing.NewEngine(),
		ProcessedPackets: make(chan *common.Packet, 1000),
		FinalOut:         finalOut,
	}

	// Connect Cover Engine to Queue Manager
	go func() {
		coverOut := p.CoverEngine.GetOutChan() // Need to add this getter
		for pkt := range coverOut {
			p.QueueMgr.Enqueue(pkt)
		}
	}()

	// Connect Queue Manager output to Normalizer then Final Out
	go func() {
		for batch := range queueBatchOut {
			for _, pkt := range batch {
				// Normalize
				fragments := p.Normalizer.Normalize(pkt)
				for _, frag := range fragments {
					// Routing
					if frag.Metadata.Route == nil {
						frag.Metadata.Route = p.RoutingEngine.SelectRoute()
					}
					p.FinalOut <- frag
				}
			}
		}
	}()

	return p
}

func (p *Pipeline) ProcessRequest(site string, data []byte) {
	id := p.IdentityMgr.GetIdentity(site)

	pkt := &common.Packet{
		ID:      fmt.Sprintf("pkt-%d", time.Now().UnixNano()),
		Payload: data,
		Size:    len(data),
		Metadata: common.Metadata{
			IdentityID:  id.ID,
			Destination: site,
		},
		Timestamp: time.Now(),
	}

	pkt.Type = p.Classifier.Classify(pkt)
	p.QueueMgr.Enqueue(pkt)
}
