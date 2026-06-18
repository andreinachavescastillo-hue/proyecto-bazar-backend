package api

import (
	"database/sql"
	"net/http"
	"rest/dto"
	"rest/security"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateUsuarioRequest struct {
	IdRol      int32  `json:"idRol" binding:"required"`
	Nombre     string `json:"nombre" binding:"required"`
	Contraseña string `json:"contraseña" binding:"required"`
	Email      string `json:"email" binding:"required"`
}

type UpdateUsuarioRequest struct {
	IdRol      int32  `json:"idRol" binding:"required"`
	Nombre     string `json:"nombre" binding:"required"`
	Contraseña string `json:"contraseña"`
	Email      string `json:"email" binding:"required"`
}

func (server *Server) CreateUsuario(ctx *gin.Context) {
	authPayload, exists := ctx.Get("authorization_payload")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	payload := authPayload.(*security.Payload)

	// aca es para que sepa que solo rol 1 o admin pueda añadir usuarios
	if payload.IDRol != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "No tienes permisos para crear usuarios"})
		return
	}

	var req CreateUsuarioRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Validar que el email no esté registrado
	existingUser, err := server.dbtx.GetUsuarioByEmail(ctx, sql.NullString{
		String: req.Email,
		Valid:  true,
	})
	if err == nil && existingUser.Idusuario > 0 {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "El email ya está registrado",
		})
		return
	}

	args := dto.CreateUsuarioParams{
		Idrol:      req.IdRol,
		Nombre:     req.Nombre,
		Contraseña: req.Contraseña,
		Email:      sql.NullString{String: req.Email, Valid: true},
	}

	result, err := server.dbtx.CreateUsuario(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	id, _ := result.LastInsertId()
	usuario, err := server.dbtx.GetUsuarioById(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"mensaje": "Usuario creado exitosamente",
		"usuario": gin.H{
			"idUsuario": usuario.Idusuario,
			"nombre":    usuario.Nombre,
			"email":     usuario.Email.String,
			"idRol":     usuario.Idrol,
		},
	})
}

func (server *Server) GetAllUsuarios(ctx *gin.Context) {
	usuarios, err := server.dbtx.GetAllUsuarios(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	var response []gin.H
	for _, u := range usuarios {
		response = append(response, gin.H{
			"idUsuario": u.Idusuario,
			"nombre":    u.Nombre,
			"email":     u.Email.String,
			"idRol":     u.Idrol,
		})
	}

	ctx.JSON(http.StatusOK, response)
}

func (server *Server) GetUsuarioById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	usuario, err := server.dbtx.GetUsuarioById(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"idUsuario": usuario.Idusuario,
		"nombre":    usuario.Nombre,
		"email":     usuario.Email.String,
		"idRol":     usuario.Idrol,
	})
}
func (server *Server) GetUsuarioByName(ctx *gin.Context) {
	name := ctx.Param("nombre")

	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Nombre no proporcionado"})
		return
	}
	searchParam := sql.NullString{
		String: "%" + name + "%",
		Valid:  true,
	}

	usuarios, err := server.dbtx.SearchUsuarioByName(ctx, searchParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var response []gin.H
	for _, u := range usuarios {
		response = append(response, gin.H{
			"idUsuario": u.Idusuario,
			"nombre":    u.Nombre,
			"email":     u.Email.String,
			"idRol":     u.Idrol,
		})
	}

	ctx.JSON(http.StatusOK, response)
}
func (server *Server) SearchUsuarioByEmail(ctx *gin.Context) {
	email := ctx.Param("email")

	if email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email no proporcionado"})
		return
	}
	searchParam := sql.NullString{
		String: "%" + email + "%",
		Valid:  true,
	}

	usuarios, err := server.dbtx.SearchUsuarioByEmail(ctx, searchParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var response []gin.H
	for _, u := range usuarios {
		response = append(response, gin.H{
			"idUsuario": u.Idusuario,
			"nombre":    u.Nombre,
			"email":     u.Email.String,
			"idRol":     u.Idrol,
		})
	}

	ctx.JSON(http.StatusOK, response)
}

func (server *Server) UpdateUsuario(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req UpdateUsuarioRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = server.dbtx.GetUsuarioById(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	args := dto.UpdateUsuarioParams{
		Idusuario: int32(id),
		Idrol:     req.IdRol,
		Nombre:    req.Nombre,
		Email:     sql.NullString{String: req.Email, Valid: true},
	}

	if req.Contraseña != "" {
		args.Contraseña = req.Contraseña
	} else {
		usuarioActual, _ := server.dbtx.GetUsuarioById(ctx, int32(id))
		args.Contraseña = usuarioActual.Contraseña
	}

	_, err = server.dbtx.UpdateUsuario(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"mensaje": "Usuario actualizado exitosamente",
	})
}

func (server *Server) DeleteUsuario(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err = server.dbtx.GetUsuarioById(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	_, err = server.dbtx.DeleteUsuario(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"mensaje": "Usuario eliminado exitosamente",
	})
}
func (server *Server) GetUsuarioByEmail(ctx *gin.Context) {
	email := ctx.Param("email")

	usuario, err := server.dbtx.GetUsuarioByEmail(ctx, sql.NullString{
		String: email,
		Valid:  true,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"idUsuario": usuario.Idusuario,
		"nombre":    usuario.Nombre,
		"email":     usuario.Email.String,
		"idRol":     usuario.Idrol,
	})
}

/*func (server *Server) GetUsuarioByName(ctx *gin.Context) {
	name := ctx.Param("nombre")

	usuarios, err := server.dbtx.GetAllUsuarios(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var response []gin.H
	for _, u := range usuarios {
		if strings.Contains(strings.ToLower(u.Nombre), strings.ToLower(name)) {
			response = append(response, gin.H{
				"idUsuario": u.Idusuario,
				"nombre":    u.Nombre,
				"email":     u.Email.String,
				"idRol":     u.Idrol,
			})
		}
	}

	ctx.JSON(http.StatusOK, response)
}*/
