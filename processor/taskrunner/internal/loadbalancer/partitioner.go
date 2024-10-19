package loadbalancer

import (
	"errors"
	"math"
	"sort"
	"time"
	"ylem_taskrunner/services/ylem_statistics"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	lru "github.com/hashicorp/golang-lru"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

const (
	HeaderPipelineId     = "x-pipeline-id"
	HeaderOrganizationId = "x-organization-id"
	HeaderPipelineRunId  = "x-pipeline-run-id"

	Rate                      = 1000                 // bucket fill rate, milliseconds
	SlowPipelineCapacity      = 30 * 1000            // if exceeded, the pipeline is considered "slow" and de-prioritized
	DefaultPipelineCapacity   = SlowPipelineCapacity // if there is no duration stats yet, this value is used as required cap, milliseconds
	FastLaneReservePercentage = 0.4                  // 40% of partitions are reserved for fast pipelines, slow pipelines get dispatched into the rest 60% of partitions
	Burst                     = 3600 * 1000          // max total capacity of all pipelines in a partition, in milliseconds. if exceeded, partitioner returns an error
)

type LBPartitioner struct {
	wfStatsProvider       PipelineStatsProvider
	wfrPartitions         *lru.Cache
	partitionRateLimiters map[int32]*rate.Limiter
}

func (p *LBPartitioner) Partition(message *sarama.ProducerMessage, numPartitions int32) (int32, error) {
	p.ensurePartitionRateLimiters(numPartitions)

	// get cached partition number by pipeline run id
	var wfrId, wfId, orgId uuid.UUID
	var err error
	for _, h := range message.Headers {
		if string(h.Key) == HeaderPipelineRunId {
			wfrId, err = uuid.FromBytes(h.Value)
			if err != nil {
				return -1, err
			}
		}

		if string(h.Key) == HeaderPipelineId {
			wfId, err = uuid.FromBytes(h.Value)
			if err != nil {
				return -1, err
			}
		}

		if string(h.Key) == HeaderOrganizationId {
			orgId, _ = uuid.FromBytes(h.Value)
		}
	}

	log.Debugf("smart lb, pipeline %s, run %s", wfId.String(), wfrId.String())

	cpn, ok := p.wfrPartitions.Get(wfrId)
	if ok {
		// if found, return
		log.Debugf("smart lb, pipeline %s, run %s, partition found in cache: %d", wfId.String(), wfrId.String(), cpn)

		return cpn.(int32), nil
	}

	// if not,
	// get approximate required capacity by pipeline id
	reqCap, err := p.getPipelineRequiredCapacity(wfId)
	if err != nil {
		return -1, err
	}

	log.Debugf("smart lb, pipeline %s, run %s, required capacity: %d", wfId.String(), wfrId.String(), reqCap)

	// get candidate partition with the most tokens available atm
	partition, err := p.getCandidatePartition(numPartitions, reqCap >= SlowPipelineCapacity, orgId)
	if err != nil {
		return -1, err
	}

	log.Debugf("smart lb, pipeline %s, run %s, candidate partition: %d", wfId.String(), wfrId.String(), partition)

	// reserve min(required tokens, burst) tokens for the candidate. this will guarantee that we always reserve tokens, see https://pkg.go.dev/golang.org/x/time/rate#Reservation.OK
	// min(required tokens, burst) will also preserve from bugs in telemetry stats, should it return required tokens > burst
	reservation := p.reserve(partition, reqCap)
	if !reservation.OK() {
		return -1, errors.New("unable to reserve capacity")
	}

	// store the candidate in cache for the pipeline run id
	p.wfrPartitions.Add(wfrId, partition)

	return partition, nil
}

func (p *LBPartitioner) getPipelineRequiredCapacity(wfId uuid.UUID) (int, error) {
	cap, err := p.wfStatsProvider.GetApproximatePipelineExecutionTime(wfId)
	if err != nil {
		log.Error(err)
		return DefaultPipelineCapacity, nil
	}

	if cap <= 0 {
		return DefaultPipelineCapacity, nil
	}

	return min(cap, Burst), nil
}

type LimiterSnapshot struct {
	partition int32
	rl        *rate.Limiter
	tokens    float64
}

type LimiterSnapshots []LimiterSnapshot

func (v LimiterSnapshots) Len() int           { return len(v) }
func (v LimiterSnapshots) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v LimiterSnapshots) Less(i, j int) bool { return v[i].tokens < v[j].tokens }

func (p *LBPartitioner) getCandidatePartition(numPartitions int32, isSlowPipeline bool, orgId uuid.UUID) (int32, error) {
	snapshots := LimiterSnapshots{}
	for partition, rl := range p.partitionRateLimiters {
		s := LimiterSnapshot{
			partition: partition,
			rl:        rl,
			tokens:    rl.Tokens(),
		}
		snapshots = append(snapshots, s)
	}

	sort.Sort(sort.Reverse(&snapshots))

	for _, s := range snapshots {
		log.Debugf("partition %d current capacity: %f", s.partition, s.tokens)
	}

	if isSlowPipeline {
		firstPartition := int32(math.Round(FastLaneReservePercentage * float64(numPartitions)))
		snapshots = snapshots[firstPartition:]

		log.Debugf("slow pipeline, reserving %d partitions for fast lane", firstPartition)
	}

	if len(snapshots) == 0 {
		return -1, errors.New("fast lane reserve percentage is too big, no partitions left for slow pipelines")
	}

	return snapshots[0].partition, nil
}

func (p *LBPartitioner) ensurePartitionRateLimiters(numPartitions int32) {
	for i := int32(0); i < numPartitions; i++ {
		if _, ok := p.partitionRateLimiters[i]; !ok {
			p.partitionRateLimiters[i] = rate.NewLimiter(Rate, Burst)
		}
	}
}

func (p *LBPartitioner) reserve(partition int32, reqCap int) *rate.Reservation {
	return p.partitionRateLimiters[partition].ReserveN(time.Now(), int(reqCap))
}

func (p *LBPartitioner) RequiresConsistency() bool {
	return true
}

type PipelineStatsProvider struct {
	client *ylem_statistics.Client
}

// returns milliseconds
func (sp PipelineStatsProvider) GetApproximatePipelineExecutionTime(wfId uuid.UUID) (int, error) {
	return sp.client.GetApproximatePipelineExecutionTime(wfId)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func NewPartitioner() (sarama.Partitioner, error) {
	cache, err := lru.New(10000)
	if err != nil {
		return nil, err
	}

	p := &LBPartitioner{
		wfStatsProvider: PipelineStatsProvider{
			client: ylem_statistics.NewClient(),
		},
		wfrPartitions:         cache,
		partitionRateLimiters: map[int32]*rate.Limiter{},
	}

	return p, nil
}
