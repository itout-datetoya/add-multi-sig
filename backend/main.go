package main

import (
	"multisigservice/db"
	"multisigservice/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// DB初期化：PostgreSQLへ接続し、テーブルを自動マイグレーション
	db.InitDB()

	router := gin.Default()

	api := router.Group("/api")
	{
		// 認証関連エンドポイント
		api.GET("/auth/challenge", handlers.ChallengeHandler)
		api.POST("/auth/login", handlers.LoginHandler)
		api.POST("/auth/registerPubkey", handlers.RegisterPubkeyHandler)

		// マルチシグ関連エンドポイント
		api.POST("/multisig/create", handlers.CreateMultiSigHandler)
		api.GET("/multisig/list", handlers.GetMultiSigListHandler)
		api.GET("/multisig/:id/data", handlers.GetMultiSigDataHandler)
		api.POST("/multisig/:id/data", handlers.UpdateMultiSigDataHandler)
	}

	router.Run(":8080")
}
