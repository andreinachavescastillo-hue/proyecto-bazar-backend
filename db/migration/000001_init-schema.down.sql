-- primero se eliminan las tablas hijas o las que tengan relacion
DROP TABLE IF EXISTS detalleventa;
DROP TABLE IF EXISTS detallecompra;
DROP TABLE IF EXISTS venta;
DROP TABLE IF EXISTS compra;
DROP TABLE IF EXISTS producto;

-- Luego las tablas padres que no dependan de otras
DROP TABLE IF EXISTS cliente;
DROP TABLE IF EXISTS metodopago;
DROP TABLE IF EXISTS usuario;
DROP TABLE IF EXISTS proveedor;
DROP TABLE IF EXISTS categoriaproducto;
DROP TABLE IF EXISTS rol; 