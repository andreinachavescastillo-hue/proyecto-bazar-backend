-- name: CreateProducto :execresult
INSERT INTO producto (
    idCategoriaProducto,
    idProveedor,
    nombre,
    descripcion,
    precioCompra,
    precioVenta,
    stock,
    imagenUrl
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetAllProductos :many
SELECT p.*, c.nombre as nombreCategoria
FROM producto p
INNER JOIN categoriaproducto c ON p.idCategoriaProducto = c.idCategoriaProducto;

-- name: GetProductoById :one
SELECT * FROM producto
WHERE idProducto = ?;

-- name: GetAllProductosWithDetails :many
SELECT 
    p.*,
    c.nombre as nombreCategoria,
    prov.nombre as nombreProveedor
FROM producto p
INNER JOIN categoriaproducto c ON p.idCategoriaProducto = c.idCategoriaProducto
INNER JOIN proveedor prov ON p.idProveedor = prov.idProveedor;

-- name: UpdateProducto :execresult
UPDATE producto
SET
    idCategoriaProducto = ?,
    idProveedor = ?,
    nombre = ?,
    descripcion = ?,
    precioCompra = ?,
    precioVenta = ?,
    stock = ?,
    imagenUrl = ?
WHERE idProducto = ?;

-- name: DeleteProducto :execresult
DELETE FROM producto
WHERE idProducto = ?;

-- name: GetProductosByCategoria :many
SELECT * FROM producto WHERE idCategoriaProducto = ?;

-- name: GetProductosByProveedor :many
SELECT * FROM producto WHERE idProveedor = ?;

-- name: UpdateProductoStock :execresult
UPDATE producto SET stock = ? WHERE idProducto = ?;

-- name: UpdateProductoStockCompra :execresult
UPDATE producto SET stock = stock + ? WHERE idProducto = ?;

-- name: SearchProductoByName :many
SELECT * FROM producto WHERE nombre LIKE ? OR descripcion LIKE ?;

-- name: SearchProductoById :one
SELECT * FROM producto WHERE idProducto = ?;