-- name: CreateUsuario :execresult
INSERT INTO usuario (idRol, nombre, contraseña, email)
VALUES (?, ?, ?, ?);

-- name: GetAllUsuarios :many
SELECT * FROM usuario;

-- name: GetUsuarioById :one
SELECT * FROM usuario
WHERE idUsuario = ?;

-- name: UpdateUsuario :execresult
UPDATE usuario
SET idRol = ?, nombre = ?, contraseña = ?, email = ?
WHERE idUsuario = ?;

-- name: DeleteUsuario :execresult
DELETE FROM usuario
WHERE idUsuario = ?;

-- name: GetUsuarioByEmail :one
SELECT * FROM usuario
WHERE email = ?;

-- name: SearchUsuarioByName :many
SELECT * FROM usuario WHERE email LIKE ?;

-- name: SearchUsuarioByEmail :many
SELECT * FROM usuario WHERE email LIKE ?;