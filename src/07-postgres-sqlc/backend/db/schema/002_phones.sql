-- 002_phones.sql
--
-- Tabela de telefones. Cada contato pode ter N telefones (1:N).
--
--   contact_id          → FK para contacts.id; cascata para apagar junto
--   label               → 'casa', 'trabalho', 'celular' etc. (não vamos validar)
--   number              → texto livre; formatação fica com a aplicação
--
-- ON DELETE CASCADE: ao deletar um contato, seus telefones somem junto.

CREATE TABLE IF NOT EXISTS phones (
    id         SERIAL PRIMARY KEY,
    contact_id INTEGER NOT NULL REFERENCES contacts(id) ON DELETE CASCADE,
    label      TEXT NOT NULL,
    number     TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS phones_contact_id_idx ON phones(contact_id);
