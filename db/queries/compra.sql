-- name: CreateCompra :execresult
INSERT INTO compra (idMetodoPago, idUsuario, idProveedor, fecha, total)
VALUES (?, ?, ?, ?, ?);

-- name: GetAllCompras :many
SELECT 
    c.idCompra,
    c.idMetodoPago,
    c.idUsuario,
    u.nombre as nombreUsuario,
    c.idProveedor,
    p.nombre as nombreProveedor,
    c.fecha,
    c.total
FROM compra c
INNER JOIN usuario u ON c.idUsuario = u.idUsuario
INNER JOIN proveedor p ON c.idProveedor = p.idProveedor
ORDER BY c.fecha DESC;

-- name: GetCompraById :one
SELECT 
    c.idCompra,
    c.idMetodoPago,
    c.idUsuario,
    u.nombre as nombreUsuario,
    c.idProveedor,
    p.nombre as nombreProveedor,
    c.fecha,
    c.total
FROM compra c
INNER JOIN usuario u ON c.idUsuario = u.idUsuario
INNER JOIN proveedor p ON c.idProveedor = p.idProveedor
WHERE c.idCompra = ?;

-- name: GetComprasByProveedor :many
SELECT 
    c.idCompra,
    c.idMetodoPago,
    c.idUsuario,
    u.nombre as nombreUsuario,
    c.idProveedor,
    p.nombre as nombreProveedor,
    c.fecha,
    c.total
FROM compra c
INNER JOIN usuario u ON c.idUsuario = u.idUsuario
INNER JOIN proveedor p ON c.idProveedor = p.idProveedor
WHERE c.idProveedor = ?
ORDER BY c.fecha DESC;

-- name: GetComprasByUsuario :many
SELECT 
    c.idCompra,
    c.idMetodoPago,
    c.idUsuario,
    u.nombre as nombreUsuario,
    c.idProveedor,
    p.nombre as nombreProveedor,
    c.fecha,
    c.total
FROM compra c
INNER JOIN usuario u ON c.idUsuario = u.idUsuario
INNER JOIN proveedor p ON c.idProveedor = p.idProveedor
WHERE c.idUsuario = ?
ORDER BY c.fecha DESC;

-- name: GetComprasByFecha :many
SELECT 
    c.idCompra,
    c.idMetodoPago,
    c.idUsuario,
    u.nombre as nombreUsuario,
    c.idProveedor,
    p.nombre as nombreProveedor,
    c.fecha,
    c.total
FROM compra c
INNER JOIN usuario u ON c.idUsuario = u.idUsuario
INNER JOIN proveedor p ON c.idProveedor = p.idProveedor
WHERE c.fecha = ?
ORDER BY c.fecha DESC;

-- name: GetComprasByFechaRange :many
SELECT 
    c.idCompra,
    c.idMetodoPago,
    c.idUsuario,
    u.nombre as nombreUsuario,
    c.idProveedor,
    p.nombre as nombreProveedor,
    c.fecha,
    c.total
FROM compra c
INNER JOIN usuario u ON c.idUsuario = u.idUsuario
INNER JOIN proveedor p ON c.idProveedor = p.idProveedor
WHERE c.fecha BETWEEN ? AND ?
ORDER BY c.fecha DESC;

-- name: GetComprasByProveedorNombre :many
SELECT 
    c.idCompra,
    c.idMetodoPago,
    c.idUsuario,
    u.nombre as nombreUsuario,
    c.idProveedor,
    p.nombre as nombreProveedor,
    c.fecha,
    c.total
FROM compra c
INNER JOIN usuario u ON c.idUsuario = u.idUsuario
INNER JOIN proveedor p ON c.idProveedor = p.idProveedor
WHERE p.nombre LIKE ?
ORDER BY c.fecha DESC;
