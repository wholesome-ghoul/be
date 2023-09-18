package middleware

import (
	"github/be/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := common.ValidateJWT(ctx); err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
