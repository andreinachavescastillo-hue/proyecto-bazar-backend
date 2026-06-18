package api

import (
	"database/sql"
	"net/http"
	"rest/dto"
	"strconv"

	"github.com/gin-gonic/gin"
)

// para recibir los datos para registrar un detalle de venta
type CreateDetalleVentaRequest struct {
	IdVenta        int32  `json:"idVenta" binding:"required"`
	IdProducto     int32  `json:"idProducto" binding:"required"`
	Cantidad       int32  `json:"cantidad" binding:"required"`
	PrecioUnitario string `json:"precioUnitario" binding:"required"`
	SubTotal       string `json:"subTotal" binding:"required"`
}

// para actualizar un detalle de venta existente
type UpdateDetalleVentaRequest struct {
	IdProducto     int32  `json:"idProducto" binding:"required"`
	Cantidad       int32  `json:"cantidad" binding:"required"`
	PrecioUnitario string `json:"precioUnitario" binding:"required"`
	SubTotal       string `json:"subTotal" binding:"required"`
}

// ========== 1. CREAR DETALLE VENTA (actualiza stock) ==========
func (server *Server) CreateDetalleVenta(ctx *gin.Context) {
	var req CreateDetalleVentaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Verificar si la venta existe
	_, err := server.dbtx.GetVentaById(ctx, req.IdVenta)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Venta no encontrada"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Verificar si el producto existe y tiene stock suficiente
	producto, err := server.dbtx.GetProductoById(ctx, req.IdProducto)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if producto.Stock < req.Cantidad {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Stock insuficiente. Disponible: " + strconv.Itoa(int(producto.Stock)),
		})
		return
	}

	// Crear detalle
	args := dto.CreateDetalleVentaParams{
		Idventa:        req.IdVenta,
		Idproducto:     req.IdProducto,
		Cantidad:       req.Cantidad,
		Preciounitario: req.PrecioUnitario,
		Subtotal:       req.SubTotal,
	}

	result, err := server.dbtx.CreateDetalleVenta(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Actualizar stock (restar cantidad vendida)
	newStock := producto.Stock - req.Cantidad
	_, err = server.dbtx.UpdateProductoStock(ctx, dto.UpdateProductoStockParams{
		Idproducto: req.IdProducto,
		Stock:      newStock,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	id, _ := result.LastInsertId()
	detalle, _ := server.dbtx.GetDetalleVentaById(ctx, int32(id))

	ctx.JSON(http.StatusCreated, gin.H{
		"mensaje": "Detalle de venta creado exitosamente",
		"detalle": detalle,
	})
}

// ========== 2. BUSCAR DETALLE POR ID ==========
func (server *Server) GetDetalleVentaById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	detalle, err := server.dbtx.GetDetalleVentaById(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"mensaje": "Detalle encontrado",
		"detalle": detalle,
	})
}

// ========== 3. OBTENER TODOS LOS DETALLES ==========
func (server *Server) GetAllDetalleVenta(ctx *gin.Context) {
	detalles, err := server.dbtx.GetAllDetalleVenta(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, detalles)
}

// ========== 4. OBTENER DETALLES POR VENTA ==========
func (server *Server) GetDetalleVentaByVenta(ctx *gin.Context) {
	idVenta, err := strconv.Atoi(ctx.Param("idVenta"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Verificar si la venta existe
	_, err = server.dbtx.GetVentaById(ctx, int32(idVenta))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Venta no encontrada"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	detalles, err := server.dbtx.GetDetalleVentaByVenta(ctx, int32(idVenta))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, detalles)
}

// ========== 5. ACTUALIZAR DETALLE VENTA (actualiza stock) ==========
func (server *Server) UpdateDetalleVenta(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Obtener detalle actual para conocer la cantidad anterior
	detalleActual, err := server.dbtx.GetDetalleVentaById(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Detalle no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var req UpdateDetalleVentaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Verificar el producto
	producto, err := server.dbtx.GetProductoById(ctx, req.IdProducto)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Ajustar stock: sumar la cantidad anterior (devuelve al stock) y restar la nueva
	stockAjustado := producto.Stock + detalleActual.Cantidad
	if stockAjustado < req.Cantidad {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Stock insuficiente. Disponible: " + strconv.Itoa(int(stockAjustado)),
		})
		return
	}

	args := dto.UpdateDetalleVentaParams{
		Iddetalleventa: int32(id),
		Idproducto:     req.IdProducto,
		Cantidad:       req.Cantidad,
		Preciounitario: req.PrecioUnitario,
		Subtotal:       req.SubTotal,
	}

	_, err = server.dbtx.UpdateDetalleVenta(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Actualizar stock con el nuevo valor
	newStock := stockAjustado - req.Cantidad
	_, err = server.dbtx.UpdateProductoStock(ctx, dto.UpdateProductoStockParams{
		Idproducto: req.IdProducto,
		Stock:      newStock,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	detalle, _ := server.dbtx.GetDetalleVentaById(ctx, int32(id))

	ctx.JSON(http.StatusOK, gin.H{
		"mensaje": "Detalle actualizado exitosamente",
		"detalle": detalle,
	})
}

// ========== 6. ELIMINAR DETALLE VENTA (restaura stock) ==========
func (server *Server) DeleteDetalleVenta(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Obtener detalle antes de eliminar para restaurar stock
	detalle, err := server.dbtx.GetDetalleVentaById(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Detalle no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Restaurar stock (sumar la cantidad que se había vendido)
	producto, err := server.dbtx.GetProductoById(ctx, detalle.Idproducto)
	if err == nil {
		newStock := producto.Stock + detalle.Cantidad
		_, err = server.dbtx.UpdateProductoStock(ctx, dto.UpdateProductoStockParams{
			Idproducto: detalle.Idproducto,
			Stock:      newStock,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	_, err = server.dbtx.DeleteDetalleVenta(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"mensaje": "Detalle eliminado correctamente",
	})
}
