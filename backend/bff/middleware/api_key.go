package middleware

import (
	"database/sql"

	"cinema.com/demo/bff/repository"
	"cinema.com/demo/bff/utils"
	"github.com/gin-gonic/gin"
)

type ApiKeyInfo struct {
	ID        string
	RateLimit int
	WindowSec int
}

func ApiKeyMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("X-API-Key")
		if apiKey == "" {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "missing api key"})
			return
		}

		hash := utils.HashApiKey(apiKey)

		apiRepo := repo.NewApiKeyRepo(db)

		key, err := apiRepo.FindByHash(hash)
		if err != nil || !key.IsActive {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "invalid api key"})
			return
		}

		ctx.Set("api_key", ApiKeyInfo{
			ID:        key.ClientID,
			RateLimit: key.RateLimit,
			WindowSec: key.RateWindowSec,
		})

		ctx.Next()
	}
}
