package api

import (
	"database/sql"
	"net/http"
	"rest/dto"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateVentaRequest struct {
	IdMetodoPago int32  `json:"idMetodoPago" binding:"required"`
	IdUsuario    int32  `json:"idUsuario" binding:"required"`
	IdCliente    int32  `json:"idCliente" binding:"required"`
	Fecha        string `json:"fecha" binding:"required"`
	SubTotal     string `json:"subTotal" binding:"required"`
	Descuento    string `json:"descuento"`
	Iva          string `json:"iva" binding:"required"`
	Total        string `json:"total" binding:"required"`
}

type UpdateVentaRequest struct {
	IdMetodoPago int32  `json:"idMetodoPago" binding:"required"`
	IdUsuario    int32  `json:"idUsuario" binding:"required"`
	IdCliente    int32  `json:"idCliente" binding:"required"`
	Fecha        string `json:"fecha" binding:"required"`
	SubTotal     string `json:"subTotal" binding:"required"`
	Descuento    string `json:"descuento"`
	Iva          string `json:"iva" binding:"required"`
	Total        string `json:"total" binding:"required"`
}

func (server *Server) CreateVenta(ctx *gin.Context) {
	var req CreateVentaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	loc, _ := time.LoadLocation("America/Costa_Rica")
	fecha, err := time.ParseInLocation("2006-01-02", req.Fecha, loc)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido"})
		return
	}

	args := dto.CreateVentaParams{
		Idmetodopago: req.IdMetodoPago,
		Idusuario:    req.IdUsuario,
		Idcliente:    req.IdCliente,
		Fecha:        fecha,
		Subtotal:     req.SubTotal,
		Descuento:    sql.NullString{String: req.Descuento, Valid: req.Descuento != ""},
		Iva:          req.Iva,
		Total:        req.Total,
	}

	result, err := server.dbtx.CreateVenta(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	id, _ := result.LastInsertId()
	venta, err := server.dbtx.GetVentaById(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"mensaje": "Venta creada exitosamente",
		"venta":   formatVentaResponse(venta),
	})
}

func (server *Server) GetVentaById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	venta, err := server.dbtx.GetVentaById(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Venta no encontrada"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, formatVentaResponse(venta))
}
func (server *Server) GetAllVentas(ctx *gin.Context) {
	ventas, err := server.dbtx.GetAllVentas(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var response []gin.H
	for _, v := range ventas {
		response = append(response, formatVentaResponse(v))
	}

	ctx.JSON(http.StatusOK, response)
}

func (server *Server) UpdateVenta(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err = server.dbtx.GetVentaById(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Venta no encontrada"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	var req UpdateVentaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	fecha, err := time.Parse("2006-01-02", req.Fecha)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido"})
		return
	}
	args := dto.UpdateVentaParams{
		Idventa:      int32(id),
		Idmetodopago: req.IdMetodoPago,
		Idusuario:    req.IdUsuario,
		Idcliente:    req.IdCliente,
		Fecha:        fecha,
		Subtotal:     req.SubTotal,
		Descuento:    sql.NullString{String: req.Descuento, Valid: req.Descuento != ""},
		Iva:          req.Iva,
		Total:        req.Total,
	}

	_, err = server.dbtx.UpdateVenta(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	venta, _ := server.dbtx.GetVentaById(ctx, int32(id))

	ctx.JSON(http.StatusOK, gin.H{
		"mensaje": "Venta actualizada exitosamente",
		"venta":   formatVentaResponse(venta),
	})
}

func (server *Server) DeleteVenta(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err = server.dbtx.GetVentaById(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Venta no encontrada"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.dbtx.DeleteVenta(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"mensaje": "Venta eliminada correctamente",
	})
}

func formatVentaResponse(venta dto.Ventum) gin.H {
	return gin.H{
		"idVenta":      venta.Idventa,
		"idMetodoPago": venta.Idmetodopago,
		"idUsuario":    venta.Idusuario,
		"idCliente":    venta.Idcliente,
		"fecha":        venta.Fecha.Format("2006-01-02"),
		"subTotal":     venta.Subtotal,
		"descuento":    venta.Descuento.String,
		"iva":          venta.Iva,
		"total":        venta.Total,
	}
}
