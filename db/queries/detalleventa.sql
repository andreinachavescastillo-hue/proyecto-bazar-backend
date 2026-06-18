-- name: CreateDetalleVenta :execresult
INSERT INTO detalleventa(
    idVenta,
    idProducto,
    cantidad,
    precioUnitario,
    subTotal
)
VALUES (?, ?, ?, ?, ?);

-- name: GetAllDetalleVenta :many
SELECT * FROM detalleventa;

-- name: GetDetalleVentaById :one
SELECT * FROM detalleventa
WHERE idDetalleVenta = ?;

-- name: GetDetalleVentaByVenta :many
SELECT * FROM detalleventa
WHERE idVenta = ?;

-- name: UpdateDetalleVenta :execresult
UPDATE detalleventa
SET
    idProducto = ?,
    cantidad = ?,
    precioUnitario = ?,
    subTotal = ?
WHERE idDetalleVenta = ?;

-- name: DeleteDetalleVenta :execresult
DELETE FROM detalleventa
WHERE idDetalleVenta = ?;