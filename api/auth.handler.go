package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email      string `json:"email" binding:"required"`
	Contraseña string `json:"contraseña" binding:"required"`
}

type LoginResponse struct {
	Token   string      `json:"token"`
	Usuario UsuarioInfo `json:"usuario"`
}

type UsuarioInfo struct {
	IDUsuario int32  `json:"idUsuario"`
	Nombre    string `json:"nombre"`
	Email     string `json:"email"`
	IdRol     int32  `json:"idRol"`
}

func (server *Server) Login(ctx *gin.Context) {
	var req LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Buscar usuario por email
	usuario, err := server.dbtx.GetUsuarioByEmail(ctx, sql.NullString{
		String: req.Email,
		Valid:  true,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Validar contraseña
	if usuario.Contraseña != req.Contraseña {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
		return
	}

	// Crear token usando PASETO
	token, err := server.tokenBuilder.CreateToken(
		int(usuario.Idrol),   // idRol
		usuario.Nombre,       // nombre
		usuario.Email.String, // email
		time.Hour*24,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Respuesta exitosa
	ctx.JSON(http.StatusOK, LoginResponse{
		Token: token,
		Usuario: UsuarioInfo{
			IDUsuario: usuario.Idusuario,
			Nombre:    usuario.Nombre,
			Email:     usuario.Email.String,
			IdRol:     usuario.Idrol,
		},
	})
}
