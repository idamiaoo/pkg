package locker

// Locker 分布式锁
type Locker interface {
	Lock() error
	Unlock() error
}
