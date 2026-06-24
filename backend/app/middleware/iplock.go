package middleware

import (
	"strings"
	"sync"
	"time"

	"rustdesk-api-server-pro/config"

	"github.com/kataras/iris/v12"
)

type ipEntry struct {
	mu          sync.Mutex
	failedCount int
	lockedUntil time.Time
}

var (
	ipRecords     sync.Map
	ipCleanupOnce sync.Once
)

func GetRealIP(ctx iris.Context) string {
	xff := ctx.GetHeader("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	xri := ctx.GetHeader("X-Real-IP")
	if xri != "" {
		return xri
	}
	return ctx.RemoteAddr()
}

func RecordFailedAttempt(ip string) {
	cfg := config.GetServerConfig().SecurityConfig
	threshold := cfg.IpLockThreshold
	if threshold <= 0 {
		threshold = 5
	}
	lockDuration := time.Duration(cfg.IpLockMinutes) * time.Minute
	if lockDuration <= 0 {
		lockDuration = 15 * time.Minute
	}

	actual, _ := ipRecords.LoadOrStore(ip, &ipEntry{})
	entry := actual.(*ipEntry)

	entry.mu.Lock()
	entry.failedCount++
	if entry.failedCount >= threshold {
		entry.lockedUntil = time.Now().Add(lockDuration)
	}
	entry.mu.Unlock()
}

func RecordSuccess(ip string) {
	ipRecords.Delete(ip)
}

func IsIpLocked(ip string) bool {
	actual, ok := ipRecords.Load(ip)
	if !ok {
		return false
	}
	entry := actual.(*ipEntry)

	entry.mu.Lock()
	defer entry.mu.Unlock()

	if entry.lockedUntil.IsZero() {
		return false
	}
	if time.Now().Before(entry.lockedUntil) {
		return true
	}
	ipRecords.Delete(ip)
	return false
}

func startIpCleanup() {
	ipCleanupOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(10 * time.Minute)
			defer ticker.Stop()
			for range ticker.C {
				ipRecords.Range(func(key, value interface{}) bool {
					entry := value.(*ipEntry)
					entry.mu.Lock()
					if !entry.lockedUntil.IsZero() && time.Now().After(entry.lockedUntil) {
						ipRecords.Delete(key)
					} else if !entry.lockedUntil.IsZero() && time.Since(entry.lockedUntil) > 24*time.Hour {
						ipRecords.Delete(key)
					} else if entry.lockedUntil.IsZero() && entry.failedCount == 0 {
						ipRecords.Delete(key)
					}
					entry.mu.Unlock()
					return true
				})
			}
		}()
	})
}
