package tokenbucket

import (
	"sync"
	"time"

	"github.com/katakurin/pkg/ratelimit"
)

// tokenBucket 令牌桶限流器
type tokenBucket struct {
	rate         int64 //固定的token放入速率, r/s
	capacity     int64 //桶的容量
	tokens       int64 //桶中当前token数量
	lastTokenSec int64 //桶上次放token的时间戳 s

	lock sync.Mutex
}

// New 生成令牌桶限流器
func New(rate, cap int64) ratelimit.Limiter {
	return &tokenBucket{
		rate:         rate,
		capacity:     cap,
		lastTokenSec: time.Now().Unix(),
	}
}

// Allow 是否通过
func (l *tokenBucket) Allow() bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	now := time.Now().Unix()
	l.tokens = l.tokens + (now-l.lastTokenSec)*l.rate // 先添加令牌

	if l.tokens > l.capacity {
		l.tokens = l.capacity
	}
	l.lastTokenSec = now
	if l.tokens <= 0 {
		// 没有令牌,则拒绝
		return false
	}
	// 还有令牌，领取令牌
	l.tokens--
	return true
}
