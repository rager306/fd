package middleware

import (
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"fd-api/embed"
	"fd-api/handlers"
	"fd-api/internal/envutil"

	"github.com/gin-gonic/gin"
)

const (
	defaultIPRateLimitRPM   = 100
	defaultUserRateLimitRPM = 1000
	rateLimitWindowSeconds  = 60
	maxRateLimitKeys        = 10_000

	headerRateLimitLimit     = "X-RateLimit-Limit"
	headerRateLimitRemaining = "X-RateLimit-Remaining"
	headerRateLimitReset     = "X-RateLimit-Reset"
)

type rateBucket struct {
	tokens float64
	last   time.Time
}

// RateLimiter is a concurrency-safe per-key token bucket limiter.
type RateLimiter struct {
	limit int
	now   func() time.Time
	mu    sync.Mutex
	items map[string]*rateBucket
}

// NewRateLimiter constructs a token bucket limiter with rpm capacity/refill.
func NewRateLimiter(rpm int) *RateLimiter {
	if rpm <= 0 {
		rpm = defaultIPRateLimitRPM
	}
	return &RateLimiter{limit: rpm, now: time.Now, items: make(map[string]*rateBucket)}
}

// IPRateLimitFromEnv returns per-IP rate limiting middleware from env.
func IPRateLimitFromEnv() gin.HandlerFunc {
	return IPRateLimit(rateLimitEnabledFromEnv(), envutil.PositiveInt("FD_RATE_LIMIT_IP_RPM", defaultIPRateLimitRPM))
}

// UserRateLimitFromEnv returns per-user rate limiting middleware from env.
func UserRateLimitFromEnv() gin.HandlerFunc {
	return UserRateLimit(rateLimitEnabledFromEnv(), envutil.PositiveInt("FD_RATE_LIMIT_USER_RPM", defaultUserRateLimitRPM))
}

// IPRateLimit limits requests per client IP when enabled.
func IPRateLimit(enabled bool, rpm int) gin.HandlerFunc {
	limiter := NewRateLimiter(rpm)
	return func(c *gin.Context) {
		if !enabled {
			c.Next()
			return
		}
		if !limiter.allow(c, "ip:"+c.ClientIP()) {
			return
		}
		c.Next()
	}
}

// UserRateLimit limits requests by request.user when enabled and present.
func UserRateLimit(enabled bool, rpm int) gin.HandlerFunc {
	limiter := NewRateLimiter(rpm)
	return func(c *gin.Context) {
		if !enabled {
			c.Next()
			return
		}
		user := userFromValidatedRequest(c)
		if user == "" {
			c.Next()
			return
		}
		if !limiter.allow(c, "user:"+user) {
			return
		}
		c.Next()
	}
}

func (l *RateLimiter) allow(c *gin.Context, key string) bool {
	l.mu.Lock()
	allowed, remaining, reset := l.take(key)
	l.mu.Unlock()

	setRateLimitHeaders(c, l.limit, remaining, reset)
	if allowed {
		return true
	}
	handlers.WriteErrorWithRetryAfter(c, handlers.CodeRateLimitExceeded, "rate_limit", "rate limit exceeded", strconv.Itoa(rateLimitWindowSeconds))
	c.Abort()
	return false
}

func (l *RateLimiter) take(key string) (allowed bool, remaining, reset int) {
	now := l.now()
	bucket := l.items[key]
	if bucket == nil {
		l.makeRoomForNewKey(now)
		bucket = &rateBucket{tokens: float64(l.limit), last: now}
		l.items[key] = bucket
	}
	refillRate := float64(l.limit) / rateLimitWindowSeconds
	elapsed := now.Sub(bucket.last).Seconds()
	bucket.tokens = math.Min(float64(l.limit), bucket.tokens+elapsed*refillRate)
	bucket.last = now

	if bucket.tokens < 1 {
		return false, 0, rateLimitWindowSeconds
	}
	bucket.tokens--
	remaining = int(math.Floor(bucket.tokens))
	reset = int(math.Ceil((float64(l.limit) - bucket.tokens) / refillRate))
	return true, remaining, reset
}

func (l *RateLimiter) makeRoomForNewKey(now time.Time) {
	if len(l.items) < maxRateLimitKeys {
		return
	}
	l.pruneExpired(now)
	if len(l.items) < maxRateLimitKeys {
		return
	}
	l.evictOldest()
}

func (l *RateLimiter) pruneExpired(now time.Time) {
	expiry := time.Duration(rateLimitWindowSeconds) * time.Second
	for key, bucket := range l.items {
		if now.Sub(bucket.last) >= expiry {
			delete(l.items, key)
		}
	}
}

func (l *RateLimiter) evictOldest() {
	var oldestKey string
	var oldestTime time.Time
	first := true
	for key, bucket := range l.items {
		if first || bucket.last.Before(oldestTime) {
			oldestKey = key
			oldestTime = bucket.last
			first = false
		}
	}
	if oldestKey != "" {
		delete(l.items, oldestKey)
	}
}

func setRateLimitHeaders(c *gin.Context, limit, remaining, reset int) {
	c.Header(headerRateLimitLimit, strconv.Itoa(limit))
	c.Header(headerRateLimitRemaining, strconv.Itoa(remaining))
	c.Header(headerRateLimitReset, strconv.Itoa(reset))
}

func userFromValidatedRequest(c *gin.Context) string {
	v, ok := c.Get(handlers.ContextKeyValidatedRequest)
	if !ok {
		return ""
	}
	req, ok := v.(*embed.EmbeddingsRequest)
	if !ok || req.User == nil {
		return ""
	}
	return *req.User
}

func rateLimitEnabledFromEnv() bool {
	return strings.EqualFold(os.Getenv("FD_RATE_LIMIT_ENABLED"), "true")
}
