-- migrations/000008_add_evento_filed.up.sql
ALTER TABLE Evento ADD COLUMN unsigned JSONB;
