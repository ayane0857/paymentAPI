package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Payment struct {
  ID               uint           `gorm:"primary_key"`
  Money            int            `gorm:"size:255"`
  PaymentMethod   string         `gorm:"size:255"`
  CreatedAt        time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
  UpdatedAt        time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
}

type Balance struct {
  ID               uint           `gorm:"primary_key"`
  Balance          int            `gorm:"size:255"`
  CreatedAt        time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
  UpdatedAt        time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
}

type BalanceResponse struct {
	Balance    float64   `json:"balance"`     // `json:"balance"`でJSONキー名を指定
	TotalMoney float64   `json:"total_money"` // こちらも同様
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func main() {
  err := godotenv.Load()
  if err != nil {
	log.Fatal("Error loading .env file")
  }

   // 環境変数から接続情報を取得
  dbUser := os.Getenv("POSTGRES_USER")
  dbPassword := os.Getenv("POSTGRES_PASSWORD")
  dbName := os.Getenv("POSTGRES_DB")
  dbHost := "localhost" // または環境変数から取得
  dbPort := "5432"      // または環境変数から取得

  dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo", dbHost, dbUser, dbPassword, dbName, dbPort)

	
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  
  if err != nil {
    log.Fatal("エラー出たよだるいね...", err)
  }
  
  db.AutoMigrate(&Payment{}, &Balance{})
  var count int64
  db.Model(&Balance{}).Count(&count)
  
  if count == 0 {
      // データが存在しない場合のみ作成
      balance := Balance{
          Balance: 0,
      }
      db.Create(&balance)
      println("データを作成しました")
  } else {
      println("データは既に存在します")
  }
  // Ginエンジンのインスタンスを作成
  r := gin.Default()

  // CORSの設定
  r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"*"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
    AllowCredentials: true,
  }))

	// ルートURL ("/") に対するGETリクエストをハンドル
	r.GET("/", func(c *gin.Context) {
		// JSONレスポンスを返す
		c.JSON(200, gin.H{
      "content": "Hello World!!",
		})
	})

  r.GET("/payments", func(c *gin.Context) {
    limitStr := c.DefaultQuery("limit", "10")
    offsetStr := c.DefaultQuery("offset", "0")

    if limitStr == "" || offsetStr == "" {
      c.JSON(http.StatusBadRequest, gin.H{"error": "Limit and offset parameters are required"})
      return
    }

		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}
    if limit <= 0 {
      c.JSON(http.StatusBadRequest, gin.H{"error": "Limit must be greater than 0"})
      return
    }
    if limit > 100 {
      c.JSON(http.StatusBadRequest, gin.H{"error": "Limit must not exceed 100"})
      return
    }

		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
			return
		}
    if offset < 0 {
      c.JSON(http.StatusBadRequest, gin.H{"error": "Offset must be 0 or greater"})
      return
    }

    var payments []Payment
		if err := db.Limit(limit).Offset(offset).Find(&payments).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve payments"})
			return
		}
    if len(payments) == 0 {
      c.JSON(http.StatusNotFound, gin.H{"error": "No payments found"})
      return
    }

    c.JSON(http.StatusOK, payments)
  })
  
  r.GET("/payments/:id", func(c *gin.Context) {
    var payment Payment
    id := c.Param("id")

    if err := db.First(&payment, id).Error; err != nil {
      c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
      return
    }

    c.JSON(http.StatusOK, payment)
  })

  r.POST("/payments", func(c *gin.Context) {
    var payment Payment

    if err := c.ShouldBindJSON(&payment); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
    }

    db.Create(&payment)
    c.JSON(http.StatusOK, payment)
  })

  r.PUT("/payments/:id", func(c *gin.Context) {
    var payment Payment
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
  })

  r.DELETE("/payments/:id", func(c *gin.Context) {
    var payment Payment
    id := c.Param("id")

    if err := db.First(&payment, id).Error; err != nil {
      c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
      return
    }

    db.Delete(&payment)
    c.JSON(http.StatusOK, payment)
  })

  r.GET("/balance", func(c *gin.Context) {
    var balance Balance
    var totalMoney int64

    if err := db.First(&balance).Error; err != nil {
      c.JSON(http.StatusNotFound, gin.H{"error": "Balance not found"})
      return
    }

    db.Model(&Payment{}).Select("COALESCE(SUM(money), 0)").Scan(&totalMoney)
    response := BalanceResponse{
      Balance:    float64(balance.Balance),
      TotalMoney: float64(totalMoney),
      UpdatedAt:  balance.UpdatedAt,
    }
    c.JSON(http.StatusOK, response)
  })

  r.PUT("/balance", func(c *gin.Context) {
      var balance Balance
      var requestBalance Balance
      
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
  })


	// 8080ポートでサーバーを起動
	r.Run(":8080")
}