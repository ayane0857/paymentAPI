package handlers

import (
	"net/http"
	"os"

	"paymentAPI/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)
func GetBalance(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var balance models.Balance
		var totalMoney int64

		if err := db.First(&balance).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Balance not found"})
			return
		}

		db.Model(&models.Payment{}).Select("COALESCE(SUM(money), 0)").Scan(&totalMoney)
		response := models.BalanceResponse{
		Balance:    float64(balance.Balance),
		TotalMoney: float64(totalMoney),
		UpdatedAt:  balance.UpdatedAt,
		}
		c.JSON(http.StatusOK, response)
	}
}
func PutBalance(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var balance models.Balance
		var requestBalance models.Balance
		
		// ヘッダーからtokenを取得
		token := c.GetHeader("X-API-Token")
		
		if token != os.Getenv("API_PASSWORD") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		
		// 既存のbalanceを取得
		if err := db.First(&balance).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Balance not found"})
			return
		}
		
		// リクエストボディから新しいbalance値を取得
		if err := c.ShouldBindJSON(&requestBalance); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		// balanceを更新
		balance.Balance = requestBalance.Balance
		db.Save(&balance)
		c.JSON(http.StatusOK, balance)
	}
}