package api

import (
	"database/sql"
	"net/http"
	"rest/dto"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UpdateCategoriaProductoRequest struct {
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
}

type getCategoriaProductoRequest struct {
	ID int32 `uri:"id" binding:"required"`
}

func (server *Server) getAllCategoriaProducto(ctx *gin.Context) {
	categories, err := server.dbtx.GetAllCategoriaProducto(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	var response []gin.H
	for _, c := range categories {
		response = append(response, gin.H{
			"idCategoriaProducto": c.Idcategoriaproducto,
			"nombre":              c.Nombre,
			"descripcion":         c.Descripcion.String,
		})
	}

	ctx.JSON(http.StatusOK, response)
}

func (server *Server) getCategoriaProductoById(ctx *gin.Context) {
	var req getCategoriaProductoRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	category, err := server.dbtx.GetCategoriaProductoByID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Categoría no encontrada"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"idCategoriaProducto": category.Idcategoriaproducto,
		"nombre":              category.Nombre,
		"descripcion":         category.Descripcion.String,
	})
}

// busca por la descripcion  y nombre con un like
func (server *Server) searchCategoriaProductoByName(ctx *gin.Context) {
	nombre := ctx.Param("nombre")

	if nombre == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Nombre no proporcionado"})
		return
	}

	searchTerm := "%" + nombre + "%"
	args := dto.SearchCategoriaProductoByNameParams{
		Nombre: searchTerm,
		Descripcion: sql.NullString{
			String: searchTerm,
			Valid:  true,
		},
	}
	categorias, err := server.dbtx.SearchCategoriaProductoByName(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var response []gin.H
	for _, c := range categorias {
		response = append(response, gin.H{
			"idCategoriaProducto": c.Idcategoriaproducto,
			"nombre":              c.Nombre,
			"descripcion":         c.Descripcion.String,
		})
	}

	ctx.JSON(http.StatusOK, response)
}

func (server *Server) updateCategoriaProducto(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req UpdateCategoriaProductoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err = server.dbtx.GetCategoriaProductoByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Categoría no encontrada"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	args := dto.UpdateCategoriaProductoParams{
		Idcategoriaproducto: int32(id),
		Nombre:              req.Nombre,
		Descripcion:         sql.NullString{String: req.Descripcion, Valid: req.Descripcion != ""},
	}

	_, err = server.dbtx.UpdateCategoriaProducto(ctx.Request.Context(), args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"mensaje": "La categoría de producto se actualizó exitosamente",
	})
}

func (server *Server) deleteCategoriaProducto(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// verfica si la categoria existe
	_, err = server.dbtx.GetCategoriaProductoByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Categoría no encontrada"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// se cuentan cuantos productos estan asociados a la categoria
	count, err := server.dbtx.CountProductosByCategoria(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":      "No se puede eliminar la categoría",
			"mensaje":    "La categoría tiene " + strconv.Itoa(int(count)) + " productos asociados",
			"productos":  count,
			"sugerencia": "Primero elimine o reasigne los productos de esta categoría",
		})
		return
	}

	//Si no tiene productos, eliminar
	_, err = server.dbtx.DeleteCategory(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"mensaje": "Categoría eliminada exitosamente",
	})
}
