package handlers

import (
	"net/http"
	"os"
	"strconv"

	"paymentAPI/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPayments(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        limitStr := c.DefaultQuery("limit", "10")
        offsetStr := c.DefaultQuery("offset", "0")

        limit, err := strconv.Atoi(limitStr)
        if err != nil || limit <= 0 || limit > 100 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
            return
        }
        offset, err := strconv.Atoi(offsetStr)
        if err != nil || offset < 0 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
            return
        }

        var payments []models.Payment
        if err := db.Limit(limit).Offset(offset).Find(&payments).Error; err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve payments"})
            return
        }
        if len(payments) == 0 {
            c.JSON(http.StatusNotFound, gin.H{"error": "No payments found"})
            return
        }
        c.JSON(http.StatusOK, payments)
    }
}

func GetPayment(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var payment models.Payment
        id := c.Param("id")

        if err := db.First(&payment, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
        return
        }

        c.JSON(http.StatusOK, payment)
    }
}

func PostPayment(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req models.PaymentRequest

        // ヘッダーからtokenを取得
		token := c.GetHeader("X-API-Token")
		
		if token != os.Getenv("API_PASSWORD") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
            return
        }

        if req.Money == nil || req.PaymentMethod == nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Money and Payment method are required"})
            return
        }
        if *req.Money <= 0 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Money must be greater than 0"})
            return
        }
        if *req.PaymentMethod == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Payment method is required"})
            return
        }

        payment := models.Payment{
            Money:         *req.Money,
            PaymentMethod: *req.PaymentMethod,
        }
        if err := db.Create(&payment).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment"})
            return
        }
        c.JSON(http.StatusOK, payment)
    }
}

func PutPayment(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var payment models.Payment

		// ヘッダーからtokenを取得
		token := c.GetHeader("X-API-Token")
		
		if token != os.Getenv("API_PASSWORD") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

        id := c.Param("id")

        if err := db.First(&payment, id).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
            return
        }

        if err := c.ShouldBindJSON(&payment); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        db.Save(&payment)
        c.JSON(http.StatusOK, payment)
    }
}

func DeletePayment(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var payment models.Payment

		// ヘッダーからtokenを取得
		token := c.GetHeader("X-API-Token")
		
		if token != os.Getenv("API_PASSWORD") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

        id := c.Param("id")

        if err := db.First(&payment, id).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
            return
        }

        db.Delete(&payment)
        c.JSON(http.StatusOK, payment)
    }
}