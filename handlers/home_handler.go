package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetHome(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// User-Agentを取得してブラウザかどうか判定
		userAgent := c.GetHeader("User-Agent")
		if isBrowser(userAgent) {
			// ブラウザの場合はリダイレクト
			c.Redirect(http.StatusMovedPermanently, "https://docs.ayane0857.net/")
			return
		}

		// ブラウザ以外の場合はシンプルなJSON応答
		c.JSON(http.StatusOK, gin.H{
			"content": "Hello World",
		})
	}
}

// ブラウザかどうかを判定する関数
func isBrowser(userAgent string) bool {
	userAgent = strings.ToLower(userAgent)
	
	// 一般的なブラウザのUser-Agentパターン
	browserPatterns := []string{
		"mozilla",
		"chrome",
		"safari",
		"firefox",
		"edge",
		"opera",
	}
	
	// ブラウザパターンに一致する場合はブラウザ
	for _, pattern := range browserPatterns {
		if strings.Contains(userAgent, pattern) {
			return true
		}
	}
	
	return false
}
