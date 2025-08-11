package main

import (
	"fmt"
	"log"
	"os"

	"paymentAPI/handlers"
	"paymentAPI/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
  err := godotenv.Load()
  if err != nil {
    log.Println("Error loading .env file")
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
  
  db.AutoMigrate(&models.Payment{}, &models.Balance{})
  var count int64
  db.Model(&models.Balance{}).Count(&count)
  
  if count == 0 {
      // データが存在しない場合のみ作成
      balance := models.Balance{
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

  r.GET("/payments", handlers.GetPayments(db))
  r.GET("/payments/:id", handlers.GetPayment(db))

  r.POST("/payments", handlers.PostPayment(db))

  r.PUT("/payments/:id", handlers.PutPayment(db))

  r.DELETE("/payments/:id", handlers.DeletePayment(db))

  r.GET("/balance", handlers.GetBalance(db))

  r.PUT("/balance", handlers.PutBalance(db))

  r.Run(":8080")
}