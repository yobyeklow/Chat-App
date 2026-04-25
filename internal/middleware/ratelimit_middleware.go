package middleware

import (
	"net/http"
	"sync"
	"time"
	"web_socket/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

type Client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mu      sync.Mutex
	clients = make(map[string]*Client)
)

func getClientIp(ctx *gin.Context) string {
	ip := ctx.ClientIP()
	if ip == "" {
		ip = ctx.Request.RemoteAddr
	}
	return ip
}

func getRateLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	client, exists := clients[ip]
	if !exists {
		requestSec := utils.GetEnvInt("RATE_LIMIT_REQUEST_SEC", 5)
		requestBrust := utils.GetEnvInt("RATE_LIMIT_REQUEST_BURST", 10)
		limiter := rate.NewLimiter(rate.Limit(requestSec), requestBrust)
		newClient := &Client{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		clients[ip] = newClient
		return limiter
	}
	client.lastSeen = time.Now()
	return client.limiter
}

func CleanupClients() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()

		for ip, client := range clients {
			if time.Since(client.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

var rateLimitLogCache = sync.Map{}

const rateLimitLogTTL = 20 * time.Second

func CheckClientRequestTime(ip string) bool {
	now := time.Now()
	val, ok := rateLimitLogCache.Load(ip)
	if ok {
		time, ok := val.(time.Time)
		if ok && now.Sub(time) < rateLimitLogTTL {
			return false
		}
	}
	rateLimitLogCache.Store(ip, now)
	return true
}

func RateLimiterMiddleware(rateLimitLogger *zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := getClientIp(ctx)
		limiter := getRateLimiter(ip)

		if !limiter.Allow() {
			if CheckClientRequestTime(ip) {
				rateLimitLogger.Warn().
					Str("method", ctx.Request.Method).
					Str("path", ctx.Request.URL.Path).
					Str("query", ctx.Request.URL.Path).
					Str("client_ip", ctx.ClientIP()).
					Str("user_agent", ctx.Request.UserAgent()).
					Str("referer", ctx.Request.Referer()).
					Str("protocol", ctx.Request.Proto).
					Str("host", ctx.Request.Host).
					Str("remote_addr", ctx.Request.RemoteAddr).
					Str("request_uri", ctx.Request.RequestURI).
					Interface("headers", ctx.Request.Header).
					Msg("Rate limiter exceeded")
			}
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"Error":   "Too many request",
				"Message": "Wait few minutes and try it again",
			})
			return
		}
		ctx.Next()
	}
}
