-- name: CreateVenta :execresult
INSERT INTO venta(
    idMetodoPago,
    idUsuario,
    idCliente,
    fecha,
    subTotal,
    descuento,
    IVA,
    total
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetAllVentas :many
SELECT * FROM venta;

-- name: GetVentaById :one
SELECT * FROM venta
WHERE idVenta = ?;

-- name: UpdateVenta :execresult
UPDATE venta
SET
    idMetodoPago = ?,
    idUsuario = ?,
    idCliente = ?,
    fecha = ?,
    subTotal = ?,
    descuento = ?,
    IVA = ?,
    total = ?
WHERE idVenta = ?;

-- name: DeleteVenta :execresult
DELETE FROM venta
WHERE idVenta = ?;