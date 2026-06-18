-- name: CreateDetalleCompra :execresult
INSERT INTO detallecompra (idCompra, idProducto, cantidad, precioCompra, subtotal)
VALUES (?, ?, ?, ?, ?);

-- name: GetDetalleCompraById :one
SELECT * FROM detallecompra WHERE idDetalleCompra = ?;

-- name: GetDetallesByCompraId :many
SELECT 
    dc.*,
    p.nombre as nombreProducto
FROM detallecompra dc
INNER JOIN producto p ON dc.idProducto = p.idProducto
WHERE dc.idCompra = ?;

-- name: UpdateDetalleCompra :execresult
UPDATE detallecompra 
SET idProducto = ?, cantidad = ?, precioCompra = ?, subtotal = ?
WHERE idDetalleCompra = ?;

-- name: DeleteDetalleCompra :execresult
DELETE FROM detallecompra WHERE idDetalleCompra = ?;

