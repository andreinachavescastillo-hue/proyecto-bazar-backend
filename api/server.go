package api

import (
	"rest/dto"
	"rest/security"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

type Server struct {
	dbtx         *dto.Queries
	router       *gin.Engine
	tokenBuilder security.Builder
}

func NewServer(dbtx *dto.Queries, secret string) (*Server, error) {
	builder, err := security.NewPasetoBuilder(secret)
	if err != nil {
		return nil, err
	}
	server := &Server{
		dbtx:         dbtx,
		tokenBuilder: builder,
	}
	router := gin.Default()

	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET,POST,PUT,DELETE",
		RequestHeaders:  "Origin,Authorization,Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	}))

	router.Static("/uploads", "./uploads")

	//sin validar
	// Login
	router.POST("/api/v1/login", server.Login)

	// Categorías
	router.GET("/api/v1/categorias", server.getAllCategoriaProducto)
	router.GET("/api/v1/categoria/:id", server.getCategoriaProductoById)
	router.GET("/api/v1/categorias/buscar/:nombre", server.searchCategoriaProductoByName)

	// Proveedores
	router.GET("/api/v1/proveedores", server.GetAllProveedores)
	router.GET("/api/v1/proveedores/:id", server.GetProveedorById)
	router.GET("/api/v1/proveedores/nombre/:nombre", server.GetProveedorByNombre)

	// Usuarioss
	router.GET("/api/v1/usuarios", server.GetAllUsuarios)
	router.GET("/api/v1/usuarios/:id", server.GetUsuarioById)
	router.GET("/api/v1/usuarios/buscar/:nombre", server.GetUsuarioByName)
	router.GET("/api/v1/usuarios/buscar/email/:email", server.SearchUsuarioByEmail)

	// Productos
	router.GET("/api/v1/productos", server.getAllProductos)
	router.GET("/api/v1/productos/:id", server.searchProductoById)
	router.GET("/api/v1/productos/categoria/:idCategoria", server.getProductosByCategoria)
	router.GET("/api/v1/productos/proveedor/:idProveedor", server.getProductosByProveedor)
	router.GET("/api/v1/productos/buscar/:termino", server.searchProductoByName)

	// Ventas
	router.GET("/api/v1/ventas", server.GetAllVentas)
	router.GET("/api/v1/ventas/:id", server.GetVentaById)

	// Detalle Venta (consultas)
	router.GET("/api/v1/detalles/venta/:idVenta", server.GetDetalleVentaByVenta)
	router.GET("/api/v1/detalles/:id", server.GetDetalleVentaById)

	// rutas protegidas
	authRoutes := router.Group("/").Use(server.authMiddleware())
	{
		// Categorías
		authRoutes.PUT("/api/v1/categoria/:id", server.updateCategoriaProducto)
		authRoutes.DELETE("/api/v1/categoria/:id", server.deleteCategoriaProducto)

		// Proveedores
		authRoutes.POST("/api/v1/proveedores", server.createProveedores)
		authRoutes.PUT("/api/v1/proveedores/:id", server.UpdateProveedores)
		authRoutes.DELETE("/api/v1/proveedores/:id", server.DeleteProveedor)

		// Usuarios
		authRoutes.POST("/api/v1/usuarios", server.CreateUsuario)
		authRoutes.PUT("/api/v1/usuarios/:id", server.UpdateUsuario)
		authRoutes.DELETE("/api/v1/usuarios/:id", server.DeleteUsuario)

		// Productos
		authRoutes.POST("/api/v1/productos", server.createProducto)
		authRoutes.PUT("/api/v1/productos/:id", server.updateProducto)
		authRoutes.DELETE("/api/v1/productos/:id", server.deleteProducto)
		authRoutes.POST("/api/v1/productos/:id/imagen", server.subirImagenProducto)

		// Ventas
		authRoutes.POST("/api/v1/ventas", server.CreateVenta)
		authRoutes.PUT("/api/v1/ventas/:id", server.UpdateVenta)
		authRoutes.DELETE("/api/v1/ventas/:id", server.DeleteVenta)

		// Detalle Venta
		authRoutes.POST("/api/v1/detalles", server.CreateDetalleVenta)
		authRoutes.PUT("/api/v1/detalles/:id", server.UpdateDetalleVenta)
		authRoutes.DELETE("/api/v1/detalles/:id", server.DeleteDetalleVenta)

		//compra
		authRoutes.POST("/api/v1/compras", server.CreateCompra)
		authRoutes.GET("/api/v1/compras", server.GetAllCompras)
		authRoutes.GET("/api/v1/compras/:id", server.GetCompraById)
		authRoutes.GET("/api/v1/compras/proveedor/:idProveedor", server.GetComprasByProveedor)
		// Compras y filtros
		authRoutes.GET("/api/v1/compras/usuario/:idUsuario", server.GetComprasByUsuario)
		authRoutes.GET("/api/v1/compras/proveedor/nombre/:nombre", server.GetComprasByProveedorNombre)
		authRoutes.GET("/api/v1/compras/fecha/:fecha", server.GetComprasByFecha)
		authRoutes.GET("/api/v1/compras/fechas", server.GetComprasByFechaRange)
	}

	server.router = router
	return server, nil
}

func (server *Server) Start(url string) error {
	return server.router.Run(url)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
