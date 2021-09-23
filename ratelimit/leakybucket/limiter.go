package leakybucket

import (
	"math"
	"sync"
	"time"

	"github.com/lunarhalos/pkg/ratelimit"
)

type leakyBucket struct {
	rate       float64 //固定每秒出水速率
	capacity   float64 //桶的容量
	water      float64 //桶中当前水量
	lastLeakMs int64   //桶上次漏水时间戳 ms

	lock sync.Mutex
}

// New 生成漏桶限流器
func New(rate, cap float64) ratelimit.Limiter {
	return &leakyBucket{
		rate:       rate,
		capacity:   cap,
		lastLeakMs: time.Now().UnixNano() / 1e6,
	}
}

func (l *leakyBucket) Allow() bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	now := time.Now().UnixNano() / 1e6
	eclipse := float64(now-l.lastLeakMs) * l.rate / 1000 //先执行漏水
	l.water = l.water - eclipse                          //计算剩余水量
	l.water = math.Max(0, l.water)                       //桶干了
	l.lastLeakMs = now
	if l.water >= l.capacity {
		// 水满，拒绝加水
		return false
	}
	// 尝试加水,并且水还未满
	l.water++
	return true
}
