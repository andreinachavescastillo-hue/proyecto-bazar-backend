-- se crea la tabla del rol de usuario
-- --------------------------------------------------------------------------------
CREATE TABLE rol (
	idRol INT PRIMARY KEY AUTO_INCREMENT,
	nombre varchar(45) UNIQUE NOT NULL
) ENGINE=InnoDB;

-- se crea la tabla del usuario con su relación hacia rol
-- --------------------------------------------------------------------------------
CREATE TABLE usuario(
	idUsuario int PRIMARY KEY AUTO_INCREMENT,
	idRol INT NOT NULL,
	nombre varchar(45) NOT NULL,
	contraseña varchar(45) NOT NULL,
	email varchar(100) UNIQUE,	

CONSTRAINT FK_USUARIO_ROL
    FOREIGN KEY (idRol)
    REFERENCES rol(idRol)
)ENGINE=InnoDB;

-- se crea la tabla de cliente
-- --------------------------------------------------------------------------------
CREATE TABLE cliente(
	idCliente int PRIMARY KEY AUTO_INCREMENT,
	nombre varchar(45) NOT NULL,
	cedula varchar(20) UNIQUE NOT NULL,
	telefono varchar(20) NOT NULL,
	email varchar(100) UNIQUE
)ENGINE=InnoDB;


-- se crea la tabla del proveedor 
-- --------------------------------------------------------------------------------
CREATE TABLE proveedor(
	idProveedor INT PRIMARY KEY AUTO_INCREMENT,
	nombre varchar(45) NOT NULL,
	cedJuridica varchar(45) UNIQUE NOT NULL,
	correo varchar(100) UNIQUE NOT NULL,
	telefono varchar(45) NOT NULL,
	telefonoContacto varchar(45),
	nombreContacto varchar(45)
)ENGINE=InnoDB;

-- tabla de categoría de productos 
-- --------------------------------------------------------------------------------
 CREATE TABLE categoriaproducto(
	idCategoriaProducto INT PRIMARY KEY AUTO_INCREMENT,
	nombre varchar(45) NOT NULL,
	descripcion varchar(100)
)ENGINE=InnoDB;

-- tabla de productos con relación a categoría y proveedor
-- --------------------------------------------------------------------------------
CREATE TABLE producto(
	idProducto INT PRIMARY KEY AUTO_INCREMENT,
	idCategoriaProducto INT NOT NULL,
	idProveedor INT NOT NULL,
	nombre varchar(45) NOT NULL,
	descripcion varchar(100) NOT NULL,
	precioCompra DECIMAL(10,2) NOT NULL,
	precioVenta DECIMAL(10,2) NOT NULL,
	stock INT NOT NULL,
	imagenUrl VARCHAR(155),

CONSTRAINT FK_PRODUCTO_CATEGORIAPRODUCTO
    FOREIGN KEY (idCategoriaProducto)
    REFERENCES categoriaproducto(idCategoriaProducto),

CONSTRAINT FK_PRODUCTO_PROVEEDOR
    FOREIGN KEY (idProveedor)
    REFERENCES proveedor(idProveedor)
)ENGINE=InnoDB;


-- se crea la tabla de metodopago 
-- --------------------------------------------------------------------------------
CREATE TABLE metodopago(
	idMetodoPago INT PRIMARY KEY AUTO_INCREMENT,
	nombre varchar(45) UNIQUE NOT NULL,
	descripcion varchar(100) NOT NULL
)ENGINE=InnoDB;

-- se crea la tabla de compra con relación a usuario 
CREATE TABLE compra(
	idCompra INT PRIMARY KEY AUTO_INCREMENT,
	idMetodoPago INT NOT NULL,
	idUsuario INT NOT NULL,
	idProveedor INT NOT NULL,
	fecha DATE NOT NULL,
	total DECIMAL(10,2) NOT NULL,
	
CONSTRAINT FK_COMPRA_METODOPAGO
    FOREIGN KEY (idMetodoPago)
    REFERENCES metodopago(idMetodoPago),

CONSTRAINT FK_COMPRA_USUARIO
    FOREIGN KEY (idUsuario)
    REFERENCES usuario(idUsuario),

CONSTRAINT FK_COMPRA_PROVEEDOR
    FOREIGN KEY (idProveedor)
    REFERENCES proveedor(idProveedor)
	
)ENGINE=InnoDB;

-- se crea la tabla de detalles de compra 
-- --------------------------------------------------------------------------------

CREATE TABLE detallecompra(
	idDetalleCompra INT PRIMARY KEY AUTO_INCREMENT,
	idCompra INT NOT NULL,
	idProducto INT NOT NULL,
	cantidad INT NOT NULL,
	precioCompra DECIMAL(10,2) NOT NULL,
	subtotal DECIMAL(10,2) NOT NULL,

CONSTRAINT FK_DETALLECOMPRA_COMPRA
    FOREIGN KEY (idCompra)
    REFERENCES compra(idCompra),

CONSTRAINT FK_DETALLECOMPRA_PRODUCTO
    FOREIGN KEY (idProducto)
    REFERENCES producto(idProducto)
)ENGINE=InnoDB;

-- se crea la tabla de ventas
-- -------------------------------------------------------------------------------- 

CREATE TABLE venta(
	idVenta INT PRIMARY KEY AUTO_INCREMENT,
	idMetodoPago INT NOT NULL,
	idUsuario INT NOT NULL,
	idCliente INT NOT NULL,
	fecha DATE NOT NULL,
	subTotal DECIMAL(10,2) NOT NULL,
	descuento DECIMAL(10,2),
	IVA DECIMAL(10,2) NOT NULL,
	total DECIMAL(10,2) NOT NULL,
	

CONSTRAINT FK_VENTA_METODOPAGO
    FOREIGN KEY (idMetodoPago)
    REFERENCES metodopago(idMetodoPago),

CONSTRAINT FK_VENTA_USUARIO
    FOREIGN KEY (idUsuario)
    REFERENCES usuario(idUsuario),

CONSTRAINT FK_VENTA_CLIENTE
    FOREIGN KEY (idCliente)
    REFERENCES cliente(idCliente)

)ENGINE=InnoDB;

-- se crea la tabla de detalle venta
-- --------------------------------------------------------------------------------
CREATE TABLE detalleventa(
	idDetalleVenta INT PRIMARY KEY AUTO_INCREMENT,
	idVenta INT NOT NULL,
	idProducto INT NOT NULL,
	cantidad INT NOT NULL,
	precioUnitario DECIMAL(10,2) NOT NULL,
	subTotal DECIMAL(10,2) NOT NULL,

CONSTRAINT FK_DETALLEVENTA_VENTA
	FOREIGN KEY (idVenta)
	REFERENCES venta(idVenta),

CONSTRAINT FK_DETALLEVENTA_PRODUCTO
	FOREIGN KEY (idProducto)
	REFERENCES producto(idProducto)
)ENGINE=InnoDB;