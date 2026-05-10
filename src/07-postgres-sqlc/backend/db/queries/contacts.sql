-- ex01/db/queries/contacts.sql
--
-- Anotações sqlc:
--   :many → SELECT que retorna 0..N linhas → ([]Contact, error)
--   :one  → SELECT/INSERT/UPDATE com RETURNING que retorna 1 linha → (Contact, error)
--   :exec → DELETE/UPDATE sem retorno → error
--
-- Cada bloco abaixo gera UMA função no Go. O nome após "name:" vira o nome
-- do método em internal/db/contacts.sql.go. Os "$1", "$2" são placeholders
-- posicionais do PostgreSQL — sqlc gera parâmetros tipados a partir deles.

-- name: ListContacts :many
SELECT * FROM contacts ORDER BY created_at DESC;

-- name: GetContact :one
SELECT * FROM contacts WHERE id = $1;

-- name: CreateContact :one
INSERT INTO contacts (name, email)
VALUES ($1, $2)
RETURNING *;

-- name: DeleteContact :exec
DELETE FROM contacts WHERE id = $1;
