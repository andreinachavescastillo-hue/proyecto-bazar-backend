package api

import (
	"database/sql"
	"net/http"
	"rest/dto"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ========== REQUESTS ==========

type createProductoRequest struct {
	Nombre              string `json:"nombre" binding:"required"`
	Descripcion         string `json:"descripcion" binding:"required"`
	PrecioCompra        string `json:"precioCompra" binding:"required"`
	Stock               int32  `json:"stock" binding:"required"`
	IdCategoriaProducto int32  `json:"idCategoriaProducto" binding:"required"`
	IdProveedor         int32  `json:"idProveedor" binding:"required"`
}

type updateProductoRequest struct {
	Nombre              string `json:"nombre" binding:"required"`
	Descripcion         string `json:"descripcion" binding:"required"`
	PrecioCompra        string `json:"precioCompra" binding:"required"`
	Stock               int32  `json:"stock" binding:"required"`
	IdCategoriaProducto int32  `json:"idCategoriaProducto" binding:"required"`
	IdProveedor         int32  `json:"idProveedor" binding:"required"`
}

func (server *Server) createProducto(ctx *gin.Context) {
	var req createProductoRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Verificar categoría
	_, err := server.dbtx.GetCategoriaProductoByID(ctx, req.IdCategoriaProducto)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Categoría no encontrada"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Verificar proveedor
	_, err = server.dbtx.GetProveedoresById(ctx, req.IdProveedor)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Proveedor no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Calcular precio de venta
	precioCompraFloat := req.PrecioCompra
	precioVenta := calcularPrecioVenta(precioCompraFloat)

	args := dto.CreateProductoParams{
		Idcategoriaproducto: req.IdCategoriaProducto,
		Idproveedor:         req.IdProveedor,
		Nombre:              req.Nombre,
		Descripcion:         req.Descripcion,
		Preciocompra:        precioCompraFloat,
		Precioventa:         precioVenta,
		Stock:               req.Stock,
		Imagenurl:           sql.NullString{String: "", Valid: false},
	}

	result, err := server.dbtx.CreateProducto(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	id, _ := result.LastInsertId()
	producto, _ := server.dbtx.GetProductoById(ctx, int32(id))

	ctx.JSON(http.StatusCreated, gin.H{
		"mensaje":  "Producto creado exitosamente",
		"producto": formatProductoResponse(producto),
	})
}

func (server *Server) getAllProductos(ctx *gin.Context) {
	productos, err := server.dbtx.GetAllProductosWithDetails(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var response []gin.H
	for _, p := range productos {
		response = append(response, gin.H{
			"idProducto":          p.Idproducto,
			"nombre":              p.Nombre,
			"descripcion":         p.Descripcion,
			"precioCompra":        p.Preciocompra,
			"precioVenta":         p.Precioventa,
			"stock":               p.Stock,
			"imagenUrl":           p.Imagenurl.String,
			"idCategoriaProducto": p.Idcategoriaproducto,
			"idProveedor":         p.Idproveedor,
			"nombreCategoria":     p.Nombrecategoria,
			"nombreProveedor":     p.Nombreproveedor,
		})
	}

	ctx.JSON(http.StatusOK, response)
}

func (server *Server) searchProductoById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	producto, err := server.dbtx.SearchProductoById(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, formatProductoResponse(producto))
}
func (server *Server) searchProductoByName(ctx *gin.Context) {
	term := ctx.Param("termino")

	if term == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Término de búsqueda no proporcionado"})
		return
	}

	searchTerm := "%" + term + "%"
	args := dto.SearchProductoByNameParams{
		Nombre:      searchTerm,
		Descripcion: searchTerm,
	}

	productos, err := server.dbtx.SearchProductoByName(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var response []gin.H
	for _, p := range productos {
		response = append(response, formatProductoResponse(p))
	}

	ctx.JSON(http.StatusOK, response)
}

func (server *Server) updateProducto(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req updateProductoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Verificar si el producto existe
	_, err = server.dbtx.GetProductoById(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Verificar categoría
	_, err = server.dbtx.GetCategoriaProductoByID(ctx, req.IdCategoriaProducto)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Categoría no encontrada"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Verificar proveedor
	_, err = server.dbtx.GetProveedoresById(ctx, req.IdProveedor)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Proveedor no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Calcular precio de venta
	precioVenta := calcularPrecioVenta(req.PrecioCompra)

	args := dto.UpdateProductoParams{
		Idproducto:          int32(id),
		Idcategoriaproducto: req.IdCategoriaProducto,
		Idproveedor:         req.IdProveedor,
		Nombre:              req.Nombre,
		Descripcion:         req.Descripcion,
		Preciocompra:        req.PrecioCompra,
		Precioventa:         precioVenta,
		Stock:               req.Stock,
		Imagenurl:           sql.NullString{String: "", Valid: false},
	}

	_, err = server.dbtx.UpdateProducto(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	producto, _ := server.dbtx.GetProductoById(ctx, int32(id))

	ctx.JSON(http.StatusOK, gin.H{
		"mensaje":  "Producto actualizado exitosamente",
		"producto": formatProductoResponse(producto),
	})
}

func (server *Server) deleteProducto(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Verificar si el producto existe
	_, err = server.dbtx.GetProductoById(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.dbtx.DeleteProducto(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"mensaje": "Producto eliminado exitosamente",
	})
}

func (server *Server) getProductosByCategoria(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("idCategoria"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de categoría inválido"})
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

	productos, err := server.dbtx.GetProductosByCategoria(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var response []gin.H
	for _, p := range productos {
		response = append(response, formatProductoResponse(p))
	}

	ctx.JSON(http.StatusOK, response)
}

func (server *Server) getProductosByProveedor(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("idProveedor"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de proveedor inválido"})
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

	productos, err := server.dbtx.GetProductosByProveedor(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var response []gin.H
	for _, p := range productos {
		response = append(response, formatProductoResponse(p))
	}

	ctx.JSON(http.StatusOK, response)
}

func (server *Server) subirImagenProducto(ctx *gin.Context) {
	// Obtener id del producto
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Verificar que el producto existe
	producto, err := server.dbtx.GetProductoById(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Obtener archivo del form data
	file, err := ctx.FormFile("imagen")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Imagen requerida"})
		return
	}

	// Validar extensión
	ext := file.Filename[len(file.Filename)-4:]
	if ext != ".jpg" && ext != "jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Formato de imagen no soportado"})
		return
	}

	// Validar tamaño
	if file.Size > 5*1024*1024 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "La imagen no puede superar los 5MB"})
		return
	}

	// Generar nombre unico
	nombreUnico := generarNombreUnico(file.Filename)
	filePath := "./uploads/productos/" + nombreUnico
	url := "http://localhost:8080/uploads/productos/" + nombreUnico

	// Guardar archivo físico
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Actualizar producto con la URL de la imagen
	args := dto.UpdateProductoParams{
		Idproducto:          int32(id),
		Idcategoriaproducto: producto.Idcategoriaproducto,
		Idproveedor:         producto.Idproveedor,
		Nombre:              producto.Nombre,
		Descripcion:         producto.Descripcion,
		Preciocompra:        producto.Preciocompra,
		Precioventa:         producto.Precioventa,
		Stock:               producto.Stock,
		Imagenurl:           sql.NullString{String: url, Valid: true},
	}

	_, err = server.dbtx.UpdateProducto(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"mensaje":   "Imagen subida correctamente",
		"imagenUrl": url,
	})
}

// el como se devuelve el producto

func formatProductoResponse(p dto.Producto) gin.H {
	return gin.H{
		"idProducto":          p.Idproducto,
		"nombre":              p.Nombre,
		"descripcion":         p.Descripcion,
		"precioCompra":        p.Preciocompra,
		"precioVenta":         p.Precioventa,
		"stock":               p.Stock,
		"imagenUrl":           p.Imagenurl.String,
		"idCategoriaProducto": p.Idcategoriaproducto,
		"idProveedor":         p.Idproveedor,
	}
}

// aca se hizo el calculo, con 40% de ganancia
func calcularPrecioVenta(precioCompra string) string {
	precio, _ := strconv.ParseFloat(precioCompra, 64)
	precioVenta := precio * 1.4
	return strconv.FormatFloat(precioVenta, 'f', 2, 64)
}

func generarNombreUnico(filename string) string {
	return time.Now().Format("20060102150405") + "_" + filename
}
