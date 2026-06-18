-- name: GetAllCategoriaProducto :many
SELECT * FROM categoriaproducto;

-- name: GetCategoriaProductoByID :one
SELECT * FROM categoriaproducto WHERE idCategoriaProducto = ?;

-- name: UpdateCategoriaProducto :execresult
UPDATE categoriaproducto SET nombre = ?, descripcion = ? WHERE idCategoriaProducto = ?;

-- name: DeleteCategory :execresult
DELETE FROM categoriaproducto WHERE idCategoriaProducto = ?;

-- name: CountProductosByCategoria :one
SELECT COUNT(*) FROM producto WHERE idCategoriaProducto = ?;

-- name: SearchCategoriaProductoByName :many
SELECT * FROM categoriaproducto WHERE nombre LIKE ? OR descripcion LIKE ?;