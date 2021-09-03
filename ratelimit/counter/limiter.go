package counter

import (
	"sync"
	"time"

	"github.com/katakurin/pkg/ratelimit"
)

type counter struct {
	rate  int           //计数周期内最多允许的请求数
	begin time.Time     //计数开始时间
	cycle time.Duration //计数周期
	count int           //计数周期内累计收到的请求数

	lock sync.Mutex
}

// New 生成计数器限流器
func New(rate int, cycle time.Duration) ratelimit.Limiter {
	return &counter{
		rate:  rate,
		begin: time.Now(),
		cycle: cycle,
	}
}

func (l *counter) Allow() bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.count == l.rate-1 {
		now := time.Now()
		if now.Sub(l.begin) < l.cycle {
			return false
		}

		//速度允许范围内，重置计数器
		l.reset(now)
		return true
	}
	//没有达到速率限制，计数加1
	l.count++
	return true
}

func (l *counter) reset(t time.Time) {
	l.begin = t
	l.count = 0
}
