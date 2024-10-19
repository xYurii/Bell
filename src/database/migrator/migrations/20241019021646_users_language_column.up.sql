SET statement_timeout = 0;

ALTER TABLE users
  ADD COLUMN language VARCHAR(6) NOT NULL DEFAULT 'pt-BR';