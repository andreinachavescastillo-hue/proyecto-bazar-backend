-- name: CreateProveedores :execresult
INSERT INTO proveedor (nombre, cedJuridica, correo, telefono, telefonoContacto, nombreContacto) VALUES (?, ?, ?, ?, ?, ?);

-- name: GetAllProveedores :many
SELECT * FROM proveedor;

-- name: GetProveedoresById :one
SELECT * FROM proveedor WHERE idProveedor = ?;

-- name: GetProveedoresByNombre :many
SELECT * FROM proveedor WHERE nombre LIKE ?;
-- name: GetProveedoresByCedJuridica :one
SELECT * FROM proveedor WHERE cedJuridica = ?;

-- name: UpdateProveedores :execresult
UPDATE proveedor SET nombre = ?, cedJuridica = ?, correo = ?, telefono = ?, telefonoContacto = ?, nombreContacto = ? WHERE idProveedor = ?;

-- name: DeleteProveedores :execresult
DELETE FROM proveedor WHERE idProveedor = ?;