package middleware

import "github.com/gin-gonic/gin"

func SecureHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Xss-Protection", "1; mode=block")
		c.Header("X-Dns-Prefetch-Control", "off")
		c.Header("X-Download-Options", "noopen")
		c.Header("Strict-Transport-Security", "max-age=15552000; includeSubDomains")
		c.Header("Referrer-Policy", "no-referrer")
		c.Header("Content-Security-Policy", "self")
		c.Header("Feature-Policy", "none")
		c.Next()
	}
}
