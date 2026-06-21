package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const currentOpenIDKey = "current_openid"

func RequireOpenID() gin.HandlerFunc {
	return func(c *gin.Context) {
		openID := strings.TrimSpace(
			firstNonEmpty(
				c.GetHeader("X-WX-OPENID"),
				c.GetHeader("X-Wx-Openid"),
				c.GetHeader("x-wx-openid"),
				c.Query("openid"),
			),
		)
		if openID == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"code":     -1,
				"errorMsg": "缺少用户身份信息 openid",
				"data":     nil,
			})
			return
		}

		c.Set(currentOpenIDKey, openID)
		c.Next()
	}
}

func CurrentOpenID(c *gin.Context) string {
	value, _ := c.Get(currentOpenIDKey)
	openID, _ := value.(string)
	return openID
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
