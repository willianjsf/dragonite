CREATE TABLE IF NOT EXISTS VersaoBackup (
    id_versao BIGSERIAL PRIMARY KEY,
    id_usuario VARCHAR(512) NOT NULL REFERENCES Usuario(id_usuario) ON DELETE CASCADE,
    algorithm VARCHAR(255) NOT NULL,
    auth_data JSONB NOT NULL,
    count INTEGER NOT NULL DEFAULT 0,
    etag VARCHAR(255) NOT NULL DEFAULT '0',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_versaobackup_usuario ON VersaoBackup (id_usuario, id_versao DESC);

CREATE TABLE IF NOT EXISTS ChaveBackup (
    id_versao BIGINT NOT NULL REFERENCES VersaoBackup(id_versao) ON DELETE CASCADE,
    id_canal VARCHAR(255) NOT NULL,
    id_sessao VARCHAR(255) NOT NULL,
    first_message_index BIGINT NOT NULL,
    forwarded_count BIGINT NOT NULL DEFAULT 0,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    session_data JSONB NOT NULL,
    PRIMARY KEY (id_versao, id_canal, id_sessao)
);