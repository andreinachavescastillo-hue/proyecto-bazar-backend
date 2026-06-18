package api

import (
	"database/sql"
	"net/http"
	"rest/dto"
	"strconv"

	"github.com/gin-gonic/gin"
)

type createProveedoresRequest struct {
	Nombre           string `json:"nombre" binding:"required"`
	CedJuridica      string `json:"cedJuridica" binding:"required"`
	Correo           string `json:"correo" binding:"required"`
	Telefono         string `json:"telefono" binding:"required"`
	TelefonoContacto string `json:"telefonoContacto"`
	NombreContacto   string `json:"nombreContacto"`
}

type UpdateProveedoresRequest struct {
	Nombre           string `json:"nombre" binding:"required"`
	CedJuridica      string `json:"cedJuridica" binding:"required"`
	Correo           string `json:"correo" binding:"required"`
	Telefono         string `json:"telefono" binding:"required"`
	TelefonoContacto string `json:"telefonoContacto"`
	NombreContacto   string `json:"nombreContacto"`
}

func (server *Server) createProveedores(ctx *gin.Context) {
	var req createProveedoresRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err := server.dbtx.GetProveedoresByCedJuridica(ctx, req.CedJuridica)
	if err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Ya existe un proveedor con esa cédula jurídica"})
		return
	}

	args := dto.CreateProveedoresParams{
		Nombre:           req.Nombre,
		Cedjuridica:      req.CedJuridica,
		Correo:           req.Correo,
		Telefono:         req.Telefono,
		Telefonocontacto: sql.NullString{String: req.TelefonoContacto, Valid: req.TelefonoContacto != ""},
		Nombrecontacto:   sql.NullString{String: req.NombreContacto, Valid: req.NombreContacto != ""},
	}

	result, err := server.dbtx.CreateProveedores(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	id, _ := result.LastInsertId()
	proveedor, err := server.dbtx.GetProveedoresById(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"mensaje":   "Proveedor creado exitosamente",
		"proveedor": formatProveedorResponse(proveedor),
	})
}

func (server *Server) GetAllProveedores(ctx *gin.Context) {
	proveedores, err := server.dbtx.GetAllProveedores(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	var response []gin.H
	for _, p := range proveedores {
		response = append(response, formatProveedorResponse(p))
	}

	ctx.JSON(http.StatusOK, response)
}

func (server *Server) GetProveedorById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	proveedor, err := server.dbtx.GetProveedoresById(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Proveedor no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, formatProveedorResponse(proveedor))
}
func (server *Server) GetProveedorByNombre(ctx *gin.Context) {
	nombre := ctx.Param("nombre")
	if nombre == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Nombre no proporcionado"})
		return
	}

	proveedores, err := server.dbtx.GetProveedoresByNombre(ctx, "%"+nombre+"%")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var response []gin.H
	for _, p := range proveedores {
		response = append(response, formatProveedorResponse(p))
	}

	ctx.JSON(http.StatusOK, response)
}
func (server *Server) UpdateProveedores(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Verificar si el proveedor existe
	_, err = server.dbtx.GetProveedoresById(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Proveedor no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var req UpdateProveedoresRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := dto.UpdateProveedoresParams{
		Idproveedor:      int32(id),
		Nombre:           req.Nombre,
		Cedjuridica:      req.CedJuridica,
		Correo:           req.Correo,
		Telefono:         req.Telefono,
		Telefonocontacto: sql.NullString{String: req.TelefonoContacto, Valid: req.TelefonoContacto != ""},
		Nombrecontacto:   sql.NullString{String: req.NombreContacto, Valid: req.NombreContacto != ""},
	}

	_, err = server.dbtx.UpdateProveedores(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Obtener el proveedor actualizado
	proveedor, err := server.dbtx.GetProveedoresById(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"mensaje":   "Proveedor actualizado exitosamente",
		"proveedor": formatProveedorResponse(proveedor),
	})
}
func (server *Server) DeleteProveedor(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = server.dbtx.GetProveedoresById(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Proveedor no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.dbtx.DeleteProveedores(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"mensaje": "Proveedor eliminado correctamente",
	})
}

func formatProveedorResponse(p dto.Proveedor) gin.H {
	return gin.H{
		"idProveedor":      p.Idproveedor,
		"nombre":           p.Nombre,
		"cedJuridica":      p.Cedjuridica,
		"correo":           p.Correo,
		"telefono":         p.Telefono,
		"telefonoContacto": p.Telefonocontacto.String,
		"nombreContacto":   p.Nombrecontacto.String,
	}
}
