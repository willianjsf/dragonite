-- migrations/000008_add_evento_unsigned_filed.down.sql
ALTER TABLE Evento DROP COLUMN unsigned;
