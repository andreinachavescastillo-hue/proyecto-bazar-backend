package api

import (
	"database/sql"
	"log"
	"net/http"
	"rest/dto"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateDetalleCompraRequest struct {
	IdProducto   int32  `json:"idProducto" binding:"required"`
	Cantidad     int32  `json:"cantidad" binding:"required"`
	PrecioCompra string `json:"precioCompra" binding:"required"`
	Subtotal     string `json:"subtotal" binding:"required"`
}

type CreateCompraRequest struct {
	IdMetodoPago int32                        `json:"idMetodoPago" binding:"required"`
	IdUsuario    int32                        `json:"idUsuario" binding:"required"`
	IdProveedor  int32                        `json:"idProveedor" binding:"required"`
	Detalles     []CreateDetalleCompraRequest `json:"detalles" binding:"required"`
}

// crea la compra
func (server *Server) CreateCompra(ctx *gin.Context) {
	var req CreateCompraRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Calcula total
	var total float64 = 0
	for _, d := range req.Detalles {
		subtotal, _ := strconv.ParseFloat(d.Subtotal, 64)
		total += subtotal
	}

	// Crear compra
	args := dto.CreateCompraParams{
		Idmetodopago: req.IdMetodoPago,
		Idusuario:    req.IdUsuario,
		Idproveedor:  req.IdProveedor,
		Fecha:        time.Now(),
		Total:        formatFloat(total),
	}
	result, err := server.dbtx.CreateCompra(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	compraID, _ := result.LastInsertId()

	for _, d := range req.Detalles {
		// Crear detalle
		_, err := server.dbtx.CreateDetalleCompra(ctx, dto.CreateDetalleCompraParams{
			Idcompra:     int32(compraID),
			Idproducto:   d.IdProducto,
			Cantidad:     d.Cantidad,
			Preciocompra: d.PrecioCompra,
			Subtotal:     d.Subtotal,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		_, err = server.dbtx.UpdateProductoStockCompra(ctx, dto.UpdateProductoStockCompraParams{
			Idproducto: d.IdProducto,
			Stock:      d.Cantidad,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"mensaje":  "Compra realizada exitosamente",
		"idCompra": compraID,
		"total":    formatFloat(total),
	})
}

// todo
func (server *Server) GetAllCompras(ctx *gin.Context) {
	compras, err := server.dbtx.GetAllCompras(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	log.Printf("Compras: %+v", compras)

	ctx.JSON(http.StatusOK, compras)
}

// se obtiene por id con los detalles incluidos
func (server *Server) GetCompraById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	compra, err := server.dbtx.GetCompraById(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Compra no encontrada"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	detalles, err := server.dbtx.GetDetallesByCompraId(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"compra":   compra,
		"detalles": detalles,
	})
}

// por proveedor
func (server *Server) GetComprasByProveedor(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("idProveedor"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	compras, err := server.dbtx.GetComprasByProveedor(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, compras)
}
func (server *Server) GetComprasByProveedorNombre(ctx *gin.Context) {
	nombre := ctx.Param("nombre")

	if nombre == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Nombre no proporcionado"})
		return
	}

	compras, err := server.dbtx.GetComprasByProveedorNombre(ctx, "%"+nombre+"%")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, compras)
}

func formatFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}

// obtiene compras por ID de usuario
func (server *Server) GetComprasByUsuario(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("idUsuario"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	compras, err := server.dbtx.GetComprasByUsuario(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, compras)
}

// obtiene compras por una fecha
func (server *Server) GetComprasByFecha(ctx *gin.Context) {
	fechaParam := ctx.Param("fecha")
	fecha, err := time.Parse("2006-01-02", fechaParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido. Use YYYY-MM-DD"})
		return
	}

	compras, err := server.dbtx.GetComprasByFecha(ctx, fecha)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, compras)
}

// obtiene compras en un rango de fechas
func (server *Server) GetComprasByFechaRange(ctx *gin.Context) {
	fechaInicio := ctx.Query("fechaInicio")
	fechaFin := ctx.Query("fechaFin")

	if fechaInicio == "" || fechaFin == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Se requieren fechaInicio y fechaFin"})
		return
	}
	inicio, err := time.Parse("2006-01-02", fechaInicio)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fechaInicio inválido"})
		return
	}
	fin, err := time.Parse("2006-01-02", fechaFin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fechaFin inválido"})
		return
	}
	compras, err := server.dbtx.GetComprasByFechaRange(ctx, dto.GetComprasByFechaRangeParams{
		FromFecha: inicio,
		ToFecha:   fin,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, compras)
}
