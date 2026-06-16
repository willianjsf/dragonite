CREATE OR REPLACE FUNCTION notify_matrix_sync_by_room()
RETURNS trigger AS $$
BEGIN
    -- We send the id_canal (Room ID) instead of the user_id
    PERFORM pg_notify('matrix_room_sync_channel', NEW.id_canal::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_notify_new_evento
AFTER INSERT OR UPDATE ON Evento
FOR EACH ROW
EXECUTE FUNCTION notify_matrix_sync_by_room();
