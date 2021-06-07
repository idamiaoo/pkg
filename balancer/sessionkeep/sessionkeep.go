package weightedroundrobin

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/pescaria/pkg/metadata"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

// Name is the name of session_keep balancer.
const Name = "session_keep"

// NewSessionKeepBuilder creates a new session_keep balancer builder.
func NewSessionKeepBuilder(rds redis.Cmdable) balancer.Builder {
	return base.NewBalancerBuilderV2(Name, &sessionKeepBuilder{rds: rds}, base.Config{HealthCheck: true})
}

type subConn struct {
	conn balancer.SubConn
	addr resolver.Address
	meta metadata.MD

	// effective weight
	ewt int64
	// current weight
	cwt int64
	// last score
	score float64
}

type Config struct {
	RedisConfig *redis.Options
}

type sessionKeepBuilder struct {
	rds      redis.Cmdable
	keepFunc func(ctx context.Context)
}

func (*sessionKeepBuilder) Build(info base.PickerBuildInfo) balancer.V2Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPickerV2(balancer.ErrNoSubConnAvailable)
	}
	p := &wrrPicker{}
	for sc, sci := range info.ReadySCs {

		md, ok := metadata.FromAttributes(sci.Address.Attributes)
		if !ok {
			md = metadata.MD{
				Weight: 1,
			}
		}
		if md.Weight == 0 {
			md.Weight = 1
		}

		subc := &subConn{
			conn:  sc,
			addr:  sci.Address,
			meta:  md,
			ewt:   md.Weight,
			score: -1,
		}
		p.subConns = append(p.subConns, subc)
	}
	return p
}

type wrrPicker struct {
	subConns []*subConn
	updateAt int64
	mu       sync.Mutex
}

func (p *wrrPicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	var (
		conn        *subConn
		totalWeight int64
	)
	p.mu.Lock()

	for _, sc := range p.subConns {
		totalWeight += sc.ewt
		sc.cwt += sc.ewt
		if conn == nil || conn.cwt < sc.cwt {
			conn = sc
		}
	}
	conn.cwt -= totalWeight
	p.mu.Unlock()
	return balancer.PickResult{SubConn: conn.conn}, nil
}
