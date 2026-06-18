package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (server *Server) authMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Obtener la cabecera de autorización
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No se proporcionó token"})
			ctx.Abort()
			return
		}

		// Verificar formato Bearer
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido. Use: Bearer <token>"})
			ctx.Abort()
			return
		}

		tokenString := parts[1]
		payload, err := server.tokenBuilder.VerifyToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o expirado"})
			ctx.Abort()
			return
		}
		ctx.Set("authorization_payload", payload)
		ctx.Set("idUsuario", payload.IDUsuario)
		ctx.Set("idRol", payload.IDRol)
		ctx.Set("nombre", payload.Nombre)
		ctx.Set("email", payload.Email)

		ctx.Next()
	}
}
