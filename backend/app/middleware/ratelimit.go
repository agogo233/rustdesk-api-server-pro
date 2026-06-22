package middleware

import (
	"sync"
	"time"

	"github.com/kataras/iris/v12"
)

type rateEntry struct {
	timestamps []time.Time
}

var (
	rateRecords     sync.Map
	rateCleanupOnce sync.Once
)

func RateLimit(maxPerMinute int) iris.Handler {
	startRateCleanup()

	return func(ctx iris.Context) {
		if maxPerMinute <= 0 {
			ctx.Next()
			return
		}

		// 跳过 CAPTCHA 接口的限流
		if ctx.Method() == "GET" {
			ctx.Next()
			return
		}

		ip := GetRealIP(ctx)
		now := time.Now()
		windowStart := now.Add(-1 * time.Minute)

		actual, _ := rateRecords.LoadOrStore(ip, &rateEntry{})
		entry := actual.(*rateEntry)

		var valid []time.Time
		for _, t := range entry.timestamps {
			if t.After(windowStart) {
				valid = append(valid, t)
			}
		}
		entry.timestamps = valid

		if len(valid) >= maxPerMinute {
			ctx.StatusCode(iris.StatusTooManyRequests)
			ctx.JSON(iris.Map{
				"code":    429,
				"message": "RequestTooFrequent",
				"data":    nil,
			})
			return
		}

		entry.timestamps = append(entry.timestamps, now)
		ctx.Next()
	}
}

func startRateCleanup() {
	rateCleanupOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(10 * time.Minute)
			defer ticker.Stop()
			for range ticker.C {
				cutoff := time.Now().Add(-1 * time.Minute)
				rateRecords.Range(func(key, value interface{}) bool {
					entry := value.(*rateEntry)
					var valid []time.Time
					for _, t := range entry.timestamps {
						if t.After(cutoff) {
							valid = append(valid, t)
						}
					}
					if len(valid) == 0 {
						rateRecords.Delete(key)
					} else {
						entry.timestamps = valid
					}
					return true
				})
			}
		}()
	})
}

func init() {
	startIpCleanup()
}