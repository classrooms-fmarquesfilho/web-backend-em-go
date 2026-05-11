-- db/queries/contacts.sql
--
-- Anotações sqlc:
--   :many → SELECT que retorna 0..N linhas → ([]Contact, error)
--   :one  → SELECT/INSERT/UPDATE com RETURNING que retorna 1 linha → (Contact, error)
--   :exec → DELETE/UPDATE sem retorno → error
--
-- Cada bloco abaixo gera UMA função no Go. O nome após "name:" vira o nome
-- do método em internal/db/contacts.sql.go. Os "$1", "$2" são placeholders
-- posicionais do PostgreSQL — sqlc gera parâmetros tipados a partir deles.

-- ── Contatos ────────────────────────────────────────────────────────────────

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

-- ── Telefones ───────────────────────────────────────────────────────────────

-- name: ListPhonesByContact :many
SELECT * FROM phones
WHERE contact_id = $1
ORDER BY id ASC;

-- name: GetPhone :one
SELECT * FROM phones WHERE id = $1;

-- name: CreatePhone :one
INSERT INTO phones (contact_id, label, number)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeletePhone :exec
DELETE FROM phones WHERE id = $1;

-- ── JOIN: contatos com seus telefones aninhados ─────────────────────────────
--
-- Esta query é bem importante. LEFT JOIN porque queremos TODOS os contatos
-- (mesmo os que não têm telefones). Resultado é "achatado": um contato com 3
-- telefones vira 3 linhas; um contato sem telefones vira 1 linha com p.* NULL.
--
-- Aliases (AS phone_id) evitam colisão com c.id no struct gerado pelo sqlc.
-- ORDER BY c.id garante que linhas do mesmo contato fiquem adjacentes —
-- essencial para a agregação em Go ser O(N).

-- name: ListContactsWithPhones :many
SELECT
    c.id, c.name, c.email, c.created_at,
    p.id     AS phone_id,
    p.label,
    p.number
FROM contacts c
LEFT JOIN phones p ON p.contact_id = c.id
ORDER BY c.id ASC, p.id ASC;
