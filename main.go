package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Payment struct {
  ID          uint           `gorm:"primary_key"`
  Money       int            `gorm:"size:255"`
  Payment     string         `gorm:"size:255"`
  CreatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
  UpdatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
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
  
  db.AutoMigrate(&Payment{})
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
    var payments []Payment
    db.Find(&payments)
    c.JSON(http.StatusOK, payments)
  })
  
  r.GET("/payments/:id", func(c *gin.Context) {
    var payment Payment
    id := c.Param("id")

    if err := db.First(&payment, id).Error; err != nil {
      c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
      return
    }

    db.Find(&payment)
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
    token := c.Param("token")
    if token != token {

    }
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

    if err := c.ShouldBindJSON(&payment); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
    }

    db.Delete(&payment)
    c.JSON(http.StatusOK, payment)
  })

	// 8080ポートでサーバーを起動
	r.Run(":8080")
}