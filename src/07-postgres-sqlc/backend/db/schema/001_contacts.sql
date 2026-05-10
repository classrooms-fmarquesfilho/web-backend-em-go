-- ex01/db/schema/001_contacts.sql
--
-- Tabela de contatos. Mantemos o mesmo modelo da Sprint 1, mas agora persistido.
--   id          → SERIAL (gerado pelo banco; autoincremento)
--   name, email → texto não-nulo; email com UNIQUE
--   created_at  → TIMESTAMPTZ com default NOW() (preenchido automaticamente)
--
-- Por que `id SERIAL` em vez de `string`?
--   sqlc mapeia SERIAL para int32 — é o tipo idiomático do Postgres para PKs
--   numéricas e simplifica o código gerado.

CREATE TABLE IF NOT EXISTS contacts (
    id         SERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    email      TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
