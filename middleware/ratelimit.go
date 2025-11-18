package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

func IPRateLimit(addr string) gin.HandlerFunc {
	limiter := redis_rate.NewLimiter(redis.NewClient(&redis.Options{
		Addr: addr,
	}))
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if result, err := limiter.Allow(c.Request.Context(), "ipRateLimit:"+ip, redis_rate.PerMinute(60)); err != nil {
			c.AbortWithStatus(500)
			return
		} else {
			c.Header("X-RateLimit-Remaining", strconv.Itoa(int(result.Remaining)))
			c.Header("X-RateLimit-Reset", strconv.Itoa(int(result.ResetAfter.Seconds())))

			if result.Remaining <= 0 {
				c.AbortWithStatus(429)
				return
			}
		}

		c.Next()
	}
}
